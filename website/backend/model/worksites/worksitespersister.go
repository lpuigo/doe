package worksites

import (
	"fmt"
	"github.com/lpuig/ewin/doe/model"
	"github.com/lpuig/ewin/doe/website/backend/persist"
	"path/filepath"
	"sync"
	"time"
)

type WorkSitesPersister struct {
	sync.RWMutex
	persister *persist.Persister

	workSites []*WorkSiteRecord
}

func NewWorkSitesPersist(dir string) (*WorkSitesPersister, error) {
	wsp := &WorkSitesPersister{
		persister: persist.NewPersister(dir),
	}
	err := wsp.persister.CheckDirectory()
	if err != nil {
		return nil, err
	}
	wsp.persister.SetPersistDelay(1 * time.Second)
	return wsp, nil
}

// LoadDirectory loads all persisted Worksite Records
func (wsp *WorkSitesPersister) LoadDirectory() error {
	wsp.Lock()
	defer wsp.Unlock()

	files, err := wsp.persister.GetFilesList()
	if err != nil {
		return fmt.Errorf("could not get files from persister: %v", err)
	}

	for _, file := range files {
		wsr, err := NewWorkSiteRecordFromFile(file)
		if err != nil {
			return fmt.Errorf("could not create worksite from '%s': %v", filepath.Base(file), err)
		}
		wsp.persister.Load(wsr)
		wsp.workSites = append(wsp.workSites, wsr)
	}
	return nil
}

// GetAll returns all contained WorkSiteRecords for which keep(wsr.Worksite) == true
func (wsp WorkSitesPersister) GetAll(keep func(ws *model.Worksite) bool) []*WorkSiteRecord {
	wsp.RLock()
	defer wsp.RUnlock()

	ws := []*WorkSiteRecord{}
	for _, wsr := range wsp.workSites {
		if keep(wsr.Worksite) {
			ws = append(ws, wsr)
		}
	}
	return ws
}

// Add adds the given WorkSiteRecord to the WorkSitesPersister and return its (new) WorkSiteRecord
func (wsp *WorkSitesPersister) Add(wsr *WorkSiteRecord) *WorkSiteRecord {
	wsp.Lock()
	defer wsp.Unlock()

	// give the record its new ID
	wsp.persister.Add(wsr)
	wsr.Id = wsr.GetId()
	wsp.workSites = append(wsp.workSites, wsr)
	return wsr
}

// Update updates the given WorkSiteRecord
func (wsp *WorkSitesPersister) Update(uwsr *WorkSiteRecord) error {
	wsp.RLock()
	defer wsp.RUnlock()

	owsr := wsp.GetById(uwsr.Id)
	if owsr == nil {
		return fmt.Errorf("id not found")
	}
	owsr.Worksite = uwsr.Worksite
	wsp.persister.MarkDirty(owsr)
	return nil
}

// Remove removes the given WorkSiteRecord from the WorkSitesPersister (pertaining file is deleted)
func (wsp *WorkSitesPersister) Remove(wsr *WorkSiteRecord) error {
	wsp.Lock()
	defer wsp.Unlock()

	err := wsp.persister.Remove(wsr)
	if err != nil {
		return err
	}

	i := wsp.findIndex(wsr)
	copy(wsp.workSites[i:], wsp.workSites[i+1:])
	wsp.workSites[len(wsp.workSites)-1] = nil // or the zero value of T
	wsp.workSites = wsp.workSites[:len(wsp.workSites)-1]
	return nil
}

func (wsp WorkSitesPersister) findIndex(wsr *WorkSiteRecord) int {
	for i, rec := range wsp.workSites {
		if rec.GetId() == wsr.GetId() {
			return i
		}
	}
	return -1
}

// GetById returns the WorkSiteRecord with given Id (or nil if Id not found)
func (wsp WorkSitesPersister) GetById(id int) *WorkSiteRecord {
	for _, wsr := range wsp.workSites {
		if wsr.Id == id {
			return wsr
		}
	}
	return nil
}