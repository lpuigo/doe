package persist

import (
	"fmt"
	"gopkg.in/src-d/go-vitess.v1/vt/log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Recorder interface {
	Id() int
	SetId(id int)
	Dirty()
	Persist(path string) error
	Remove(path string) error
}

const (
	DefaultPersistDelay = 2 * time.Second
	ParallelPersister   = 10
)

type Persister struct {
	directory string
	delay     time.Duration
	records   map[int]Recorder
	nextId    int

	mut          sync.RWMutex
	dirtyIds     []int
	persistTimer *time.Timer
}

func NewPersister(dir string) *Persister {
	return &Persister{
		directory: dir,
		records:   make(map[int]Recorder),
		delay:     DefaultPersistDelay,
	}
}

// SetPersistDelay sets the Pesistance Delay of the Persister
func (p *Persister) SetPersistDelay(persistDelay time.Duration) {
	p.delay = persistDelay
}

// CheckDirectory checks if Persister directory exists (if ok, return err is nil)
func (p Persister) CheckDirectory() error {
	fi, err := os.Stat(p.directory)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return fmt.Errorf("not a proper directory: %s\n", p.directory)
	}
	return nil
}

// GetFilesList returns all the record files contained in persister directory (User class is responsible to Load the record)
func (p Persister) GetFilesList() (list []string, err error) {
	err = filepath.Walk(p.directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		list = append(list, path)
		return nil
	})
	if err != nil {
		return
	}
	return
}

// Add adds the given Record to the Persister, triggers Persit mechanism and returns its (new) id
func (p *Persister) Add(r Recorder) int {
	p.mut.Lock()
	defer p.mut.Unlock()
	defer func() { p.nextId++ }()

	r.SetId(p.nextId)
	p.records[p.nextId] = r
	p.markDirty(r)

	return p.nextId
}

// Load adds the given Record to the Persister
func (p *Persister) Load(r Recorder) {
	if _, ok := p.records[r.Id()]; ok {
		panic(fmt.Sprintf("persister already contains given record with Id %d", r.Id()))
	}
	p.mut.Lock()
	defer p.mut.Unlock()
	p.records[r.Id()] = r
	if p.nextId <= r.Id() {
		p.nextId = r.Id() + 1
	}
}

// markDirty marks the given recorder as dirty and triggers the persistence mechanism
func (p *Persister) MarkDirty(r Recorder) {
	p.mut.Lock()
	defer p.mut.Unlock()
	p.markDirty(r)
}

func (p *Persister) markDirty(r Recorder) {
	if _, ok := p.records[r.Id()]; !ok {
		panic("persister does not contains given record")
	}
	r.Dirty()
	p.dirtyIds = append(p.dirtyIds, r.Id())
	p.triggerPersist()
}

// Remove removes the given recorder from the persister (pertaining persisted file is deleted)
func (p *Persister) Remove(r Recorder) {
	id := r.Id()
	if _, ok := p.records[id]; !ok {
		panic("persister does not contains given record")
	}
	p.mut.Lock()
	defer p.mut.Unlock()
	go func(dr Recorder) {
		err := dr.Remove(p.directory)
		if err != nil {
			log.Errorf("error removing record ID %d: %v\n", id, err)
		}
	}(r)
	delete(p.records, id)
}

// PersistAll immediatly persist all contained recorder(persistance delay is ignored)
func (p *Persister) PersistAll(r Recorder) {
	p.mut.Lock()
	defer p.mut.Unlock()
	if p.persistTimer != nil {
		p.persistTimer.Stop()
		p.persistTimer = nil
	}

	token := make(chan struct{}, ParallelPersister)
	defer close(token)
	for _, r := range p.records {
		token <- struct{}{}
		go func(pr Recorder) {
			err := r.Persist(p.directory)
			if err != nil {
				log.Errorf("error persisting record ID %d: %v\n", r.Id(), err)
			}
			_ = <-token
		}(r)
	}

	for i := 0; i < ParallelPersister; i++ {
		token <- struct{}{}
	}
}

func (p *Persister) triggerPersist() {
	if p.delay == 0 {
		if p.persistTimer != nil {
			p.persistTimer.Stop()
			p.persistTimer = nil
		}
		p.persist()
		return
	}
	if p.persistTimer != nil {
		return
	}
	p.persistTimer = time.AfterFunc(p.delay, func() {
		p.mut.Lock()
		defer p.mut.Unlock()
		p.persistTimer = nil
		p.persist()
	})
}

func (p *Persister) persist() {
	token := make(chan struct{}, ParallelPersister)
	for _, id := range p.dirtyIds {
		r, found := p.records[id]
		if !found { // can happen if record was remove before persistence was triggered
			continue
		}
		token <- struct{}{}
		go func(pr Recorder) {
			err := pr.Persist(p.directory)
			if err != nil {
				log.Errorf("error persisting record ID %d: %v\n", pr.Id(), err)
			}
			_ = <-token
		}(r)
	}
	for i := 0; i < ParallelPersister; i++ {
		token <- struct{}{}
	}
	p.dirtyIds = []int{}
}
