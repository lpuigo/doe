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
	SetId(id int)
	Dirty()
	Persist(path string) error
}

type Persister struct {
	directory string
	delay     time.Duration
	index     map[int]int
	records   []Recorder
	nextId    int

	mut          sync.RWMutex
	dirtyIds     []int
	persistTimer *time.Timer
}

func NewPersister(dir string, persistDelay time.Duration) *Persister {
	return &Persister{
		directory: dir,
		index:     make(map[int]int),
		delay:     persistDelay,
	}
}

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

// Add adds the given Record to the Persister, triggers Persit mechanism and return its (new) id
func (p *Persister) Add(r Recorder) int {
	p.mut.Lock()
	defer p.mut.Unlock()
	defer func() { p.nextId++ }()

	r.SetId(p.nextId)
	p.index[p.nextId] = len(p.records)
	p.records = append(p.records, r)

	p.setDirty(p.nextId)

	return p.nextId
}

func (p *Persister) triggerPersist() {
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

func (p *Persister) setDirty(id int) {
	i, found := p.index[id]
	if !found {
		return
	}
	p.records[i].Dirty()
	p.dirtyIds = append(p.dirtyIds, id)
	p.triggerPersist()
}

func (p *Persister) persist() {
	for _, id := range p.dirtyIds {
		r := p.records[p.index[id]]
		err := r.Persist(p.directory)
		if err != nil {
			log.Errorf("error persisting record ID %d: %v\n")
			continue
		}
	}
	p.dirtyIds = []int{}
}
