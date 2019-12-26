package foasites

import (
	"fmt"
	"github.com/lpuig/ewin/doe/website/backend/model/archives"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	"github.com/lpuig/ewin/doe/website/backend/model/items"
	"github.com/lpuig/ewin/doe/website/backend/persist"
	"path/filepath"
	"sync"
	"time"
)

type FoaSitesPersister struct {
	sync.RWMutex
	persister *persist.Persister

	FoaSites []*FoaSiteRecord
}

func NewFoaSitesPersist(dir string) (*FoaSitesPersister, error) {
	fsp := &FoaSitesPersister{
		persister: persist.NewPersister("Foasites", dir),
	}
	err := fsp.persister.CheckDirectory()
	if err != nil {
		return nil, err
	}
	fsp.persister.SetPersistDelay(1 * time.Second)
	return fsp, nil
}

func (fsp FoaSitesPersister) NbSites() int {
	return len(fsp.FoaSites)
}

// LoadDirectory loads all persisted FoaSiteRecords
func (fsp *FoaSitesPersister) LoadDirectory() error {
	fsp.Lock()
	defer fsp.Unlock()

	fsp.persister.Reinit()
	fsp.FoaSites = []*FoaSiteRecord{}

	files, err := fsp.persister.GetFilesList("deleted")
	if err != nil {
		return fmt.Errorf("could not get files from FoaSites persister: %v", err)
	}

	for _, file := range files {
		psr, err := NewFoaSiteRecordFromFile(file)
		if err != nil {
			return fmt.Errorf("could not instantiate FoaSite from '%s': %v", filepath.Base(file), err)
		}
		err = fsp.persister.Load(psr)
		if err != nil {
			return fmt.Errorf("error while loading %s: %s", file, err.Error())
		}
		fsp.FoaSites = append(fsp.FoaSites, psr)
	}
	return nil
}

// GetAll returns all contained FoaSiteRecord for which keep(*FoaSite) == true
func (fsp *FoaSitesPersister) GetAll(isSiteVisible items.IsItemizableSiteVisible) []*FoaSiteRecord {
	fsp.RLock()
	defer fsp.RUnlock()

	foaSiteRecords := []*FoaSiteRecord{}
	for _, psr := range fsp.FoaSites {
		if isSiteVisible(psr.FoaSite) {
			foaSiteRecords = append(foaSiteRecords, psr)
		}
	}
	return foaSiteRecords
}

// GetItemizableSites returns all contained FoaSiteRecord as ItemizableSite for which isSiteVisible(*ItemizableSite) == true
func (fsp *FoaSitesPersister) GetItemizableSites(isSiteVisible items.IsItemizableSiteVisible) []items.ItemizableSite {
	fsp.RLock()
	defer fsp.RUnlock()

	iss := []items.ItemizableSite{}
	for _, psr := range fsp.FoaSites {
		if isSiteVisible(psr) {
			iss = append(iss, psr)
		}
	}
	return iss
}

func (fsp *FoaSitesPersister) GetItemizableSiteById(id int) items.ItemizableSite {
	site := fsp.GetById(id)
	if site == nil {
		return nil // force untyped nil to enable (return == nil) test
	}
	return site
}

// GetById returns the FoaSiteRecord with given Id (or nil if Id not found)
func (fsp *FoaSitesPersister) GetById(id int) *FoaSiteRecord {
	fsp.RLock()
	defer fsp.RUnlock()

	for _, wsr := range fsp.FoaSites {
		if wsr.Id == id {
			return wsr
		}
	}
	return nil
}

// Add adds the given FoaSiteRecord to the FoaSitesPersister and return its (new) FoaSiteRecord
func (fsp *FoaSitesPersister) Add(psr *FoaSiteRecord) *FoaSiteRecord {
	fsp.Lock()
	defer fsp.Unlock()

	// Set the Update Date
	psr.FoaSite.UpdateDate = date.Today().String()

	// give the record its new ID
	fsp.persister.Add(psr)
	psr.Id = psr.GetId()
	fsp.FoaSites = append(fsp.FoaSites, psr)
	return psr
}

// Update updates the given FoaSiteRecord
func (fsp *FoaSitesPersister) Update(usr *FoaSiteRecord) error {
	fsp.RLock()
	defer fsp.RUnlock()

	osr := fsp.GetById(usr.Id)
	if osr == nil {
		return fmt.Errorf("id not found")
	}
	osr.FoaSite = usr.FoaSite
	osr.FoaSite.UpdateDate = date.Today().String()
	fsp.persister.MarkDirty(osr)
	return nil
}

// Remove removes the given FoaSiteRecord from the FoaSitesPersister (pertaining file is moved to deleted dir)
func (fsp *FoaSitesPersister) Remove(psr *FoaSiteRecord) error {
	fsp.Lock()
	defer fsp.Unlock()

	err := fsp.persister.Remove(psr)
	if err != nil {
		return err
	}

	i := fsp.findIndex(psr)
	copy(fsp.FoaSites[i:], fsp.FoaSites[i+1:])
	fsp.FoaSites[len(fsp.FoaSites)-1] = nil // or the zero value of T
	fsp.FoaSites = fsp.FoaSites[:len(fsp.FoaSites)-1]
	return nil
}

func (fsp FoaSitesPersister) findIndex(sr *FoaSiteRecord) int {
	for i, rec := range fsp.FoaSites {
		if rec.GetId() == sr.GetId() {
			return i
		}
	}
	return -1
}

func (fsp *FoaSitesPersister) GetAllSites() []archives.ArchivableRecord {
	fsp.RLock()
	defer fsp.RUnlock()

	archivableSites := make([]archives.ArchivableRecord, len(fsp.FoaSites))
	for i, site := range fsp.FoaSites {
		archivableSites[i] = site
	}
	return archivableSites
}

func (fsp *FoaSitesPersister) GetName() string {
	return "FoaSites"
}
