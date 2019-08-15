package polesites

import (
	"archive/zip"
	"fmt"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	"github.com/lpuig/ewin/doe/website/backend/persist"
	"io"
	"path/filepath"
	"sync"
	"time"
)

type PoleSitesPersister struct {
	sync.RWMutex
	persister *persist.Persister

	polesites []*PoleSiteRecord
}

func NewPoleSitesPersist(dir string) (*PoleSitesPersister, error) {
	psp := &PoleSitesPersister{
		persister: persist.NewPersister(dir),
	}
	err := psp.persister.CheckDirectory()
	if err != nil {
		return nil, err
	}
	psp.persister.SetPersistDelay(1 * time.Second)
	return psp, nil
}

func (psp PoleSitesPersister) NbSites() int {
	return len(psp.polesites)
}

// LoadDirectory loads all persisted PoleSiteRecords
func (psp *PoleSitesPersister) LoadDirectory() error {
	psp.Lock()
	defer psp.Unlock()

	psp.persister.Reinit()
	psp.polesites = []*PoleSiteRecord{}

	files, err := psp.persister.GetFilesList("deleted")
	if err != nil {
		return fmt.Errorf("could not get files from polesites persister: %v", err)
	}

	for _, file := range files {
		psr, err := NewPoleSiteRecordFromFile(file)
		if err != nil {
			return fmt.Errorf("could not instantiate polesite from '%s': %v", filepath.Base(file), err)
		}
		err = psp.persister.Load(psr)
		if err != nil {
			return fmt.Errorf("error while loading %s: %s", file, err.Error())
		}
		psp.polesites = append(psp.polesites, psr)
	}
	return nil
}

// GetAll returns all contained PoleSiteRecord for which keep(sr.Site) == true
func (psp PoleSitesPersister) GetAll(keep func(s *PoleSite) bool) []*PoleSiteRecord {
	psp.RLock()
	defer psp.RUnlock()

	srs := []*PoleSiteRecord{}
	for _, psr := range psp.polesites {
		if keep(psr.PoleSite) {
			srs = append(srs, psr)
		}
	}
	return srs
}

// GetById returns the PoleSiteRecord with given Id (or nil if Id not found)
func (psp *PoleSitesPersister) GetById(id int) *PoleSiteRecord {
	psp.RLock()
	defer psp.RUnlock()

	for _, wsr := range psp.polesites {
		if wsr.Id == id {
			return wsr
		}
	}
	return nil
}

// Add adds the given PoleSiteRecord to the PoleSitesPersister and return its (new) SiteRecord
func (psp *PoleSitesPersister) Add(psr *PoleSiteRecord) *PoleSiteRecord {
	psp.Lock()
	defer psp.Unlock()

	// Set the Update Date
	psr.PoleSite.UpdateDate = date.Today().String()

	// give the record its new ID
	psp.persister.Add(psr)
	psr.Id = psr.GetId()
	psp.polesites = append(psp.polesites, psr)
	return psr
}

// Update updates the given PoleSiteRecord
func (psp *PoleSitesPersister) Update(usr *PoleSiteRecord) error {
	psp.RLock()
	defer psp.RUnlock()

	osr := psp.GetById(usr.Id)
	if osr == nil {
		return fmt.Errorf("id not found")
	}
	osr.PoleSite = usr.PoleSite
	osr.PoleSite.UpdateDate = date.Today().String()
	psp.persister.MarkDirty(osr)
	return nil
}

// Remove removes the given PoleSiteRecord from the PoleSitesPersister (pertaining file is moved to deleted dir)
func (psp *PoleSitesPersister) Remove(psr *PoleSiteRecord) error {
	psp.Lock()
	defer psp.Unlock()

	err := psp.persister.Remove(psr)
	if err != nil {
		return err
	}

	i := psp.findIndex(psr)
	copy(psp.polesites[i:], psp.polesites[i+1:])
	psp.polesites[len(psp.polesites)-1] = nil // or the zero value of T
	psp.polesites = psp.polesites[:len(psp.polesites)-1]
	return nil
}

func (psp PoleSitesPersister) findIndex(sr *PoleSiteRecord) int {
	for i, rec := range psp.polesites {
		if rec.GetId() == sr.GetId() {
			return i
		}
	}
	return -1
}

// ArchiveName returns the PoleSiteArchive file name with today's date
func (psp PoleSitesPersister) ArchiveName() string {
	return fmt.Sprintf("Polesites %s.zip", date.Today().String())
}

// CreateArchive writes a zipped archive of all contained Polesites files to the given writer
func (psp *PoleSitesPersister) CreateArchive(writer io.Writer) error {
	psp.RLock()
	defer psp.RUnlock()

	zw := zip.NewWriter(writer)

	for _, sr := range psp.polesites {
		wfw, err := zw.Create(sr.GetFileName())
		if err != nil {
			return fmt.Errorf("could not create zip entry for polesite %d", sr.Id)
		}
		err = sr.Marshall(wfw)
		if err != nil {
			return fmt.Errorf("could not write zip entry for polesite %d", sr.Id)
		}
	}

	return zw.Close()
}
