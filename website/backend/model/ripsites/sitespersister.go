package ripsites

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

type SitesPersister struct {
	sync.RWMutex
	persister *persist.Persister

	sites []*SiteRecord
}

func NewSitesPersit(dir string) (*SitesPersister, error) {
	sp := &SitesPersister{
		persister: persist.NewPersister(dir),
	}
	err := sp.persister.CheckDirectory()
	if err != nil {
		return nil, err
	}
	sp.persister.SetPersistDelay(1 * time.Second)
	return sp, nil
}

func (sp SitesPersister) NbSites() int {
	return len(sp.sites)
}

// LoadDirectory loads all persisted Site Records
func (sp *SitesPersister) LoadDirectory() error {
	sp.Lock()
	defer sp.Unlock()

	files, err := sp.persister.GetFilesList("deleted")
	if err != nil {
		return fmt.Errorf("could not get files from sites persister: %v", err)
	}

	for _, file := range files {
		wsr, err := NewSiteRecordFromFile(file)
		if err != nil {
			return fmt.Errorf("could not create site from '%s': %v", filepath.Base(file), err)
		}
		sp.persister.Load(wsr)
		sp.sites = append(sp.sites, wsr)
	}
	return nil
}

// GetAll returns all contained SiteRecords for which keep(sr.Site) == true
func (sp SitesPersister) GetAll(keep func(s *Site) bool) []*SiteRecord {
	sp.RLock()
	defer sp.RUnlock()

	srs := []*SiteRecord{}
	for _, sr := range sp.sites {
		if keep(sr.Site) {
			srs = append(srs, sr)
		}
	}
	return srs
}

// GetById returns the SiteRecord with given Id (or nil if Id not found)
func (sp *SitesPersister) GetById(id int) *SiteRecord {
	sp.RLock()
	defer sp.RUnlock()

	for _, wsr := range sp.sites {
		if wsr.Id == id {
			return wsr
		}
	}
	return nil
}

// Add adds the given SiteRecord to the SitesPersister and return its (new) SiteRecord
func (sp *SitesPersister) Add(sr *SiteRecord) *SiteRecord {
	sp.Lock()
	defer sp.Unlock()

	// Set the Update Date
	sr.Site.UpdateDate = date.Today().String()

	// give the record its new ID
	sp.persister.Add(sr)
	sr.Id = sr.GetId()
	sp.sites = append(sp.sites, sr)
	return sr
}

// Update updates the given WorkSiteRecord
func (sp *SitesPersister) Update(usr *SiteRecord) error {
	sp.RLock()
	defer sp.RUnlock()

	osr := sp.GetById(usr.Id)
	if osr == nil {
		return fmt.Errorf("id not found")
	}
	osr.Site = usr.Site
	osr.Site.UpdateDate = date.Today().String()
	sp.persister.MarkDirty(osr)
	return nil
}

// Remove removes the given SiteRecord from the SitesPersister (pertaining file is moved to deleted dir)
func (sp *SitesPersister) Remove(sr *SiteRecord) error {
	sp.Lock()
	defer sp.Unlock()

	err := sp.persister.Remove(sr)
	if err != nil {
		return err
	}

	i := sp.findIndex(sr)
	copy(sp.sites[i:], sp.sites[i+1:])
	sp.sites[len(sp.sites)-1] = nil // or the zero value of T
	sp.sites = sp.sites[:len(sp.sites)-1]
	return nil
}

func (sp SitesPersister) findIndex(sr *SiteRecord) int {
	for i, rec := range sp.sites {
		if rec.GetId() == sr.GetId() {
			return i
		}
	}
	return -1
}

// WorksitesArchiveName returns the SiteArchive file name with today's date
func (sp SitesPersister) ArchiveName() string {
	return fmt.Sprintf("Ripsites %s.zip", date.Today().String())
}

// CreateWorksitesArchive writes a zipped archive of all contained Worksites files to the given writer
func (sp *SitesPersister) CreateArchive(writer io.Writer) error {
	sp.RLock()
	defer sp.RUnlock()

	zw := zip.NewWriter(writer)

	for _, sr := range sp.sites {
		wfw, err := zw.Create(sr.GetFileName())
		if err != nil {
			return fmt.Errorf("could not create zip entry for site %d", sr.Id)
		}
		err = sr.Marshall(wfw)
		if err != nil {
			return fmt.Errorf("could not write zip entry for site %d", sr.Id)
		}
	}

	return zw.Close()
}
