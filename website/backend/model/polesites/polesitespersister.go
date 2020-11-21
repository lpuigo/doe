package polesites

import (
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/lpuig/ewin/doe/website/backend/model/archives"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	"github.com/lpuig/ewin/doe/website/backend/model/items"
	"github.com/lpuig/ewin/doe/website/backend/persist"
)

type PoleSitesPersister struct {
	sync.RWMutex
	persister *persist.Persister

	polesites []*PoleSiteRecord
}

func NewPoleSitesPersist(dir string) (*PoleSitesPersister, error) {
	psp := &PoleSitesPersister{
		persister: persist.NewPersister("Polesites", dir),
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
func (psp PoleSitesPersister) GetAll(isSiteVisible items.IsItemizableSiteVisible) []*PoleSiteRecord {
	psp.RLock()
	defer psp.RUnlock()

	srs := []*PoleSiteRecord{}
	for _, psr := range psp.polesites {
		if isSiteVisible(psr.PoleSite) {
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

// Update updates the given PoleSiteRecord with per pole control : outdated poles are ignored (outdated : one original pole' timestamp is newer than the corresponding pole in the updated polesite)
func (psp *PoleSitesPersister) Update(psr *PoleSiteRecord) ([]string, error) {
	psp.RLock()
	defer psp.RUnlock()
	ignoredPolesRefs := []string{}
	osr := psp.GetById(psr.Id)
	if osr == nil {
		return ignoredPolesRefs, fmt.Errorf("id not found")
	}
	//osr.PoleSite = usr.PoleSite
	ignoredPolesRefs = osr.UpdateWith(psr.PoleSite)
	osr.PoleSite.UpdateDate = date.Today().String()
	psp.persister.MarkDirty(osr)
	psr.PoleSite = osr.PoleSite // psr now points to server updated polesite
	return ignoredPolesRefs, nil
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

// ArchiveCompletedPoleRefs extracts all completed ref group from given PoleSiteRecord and move them to an new archive PoleSiteRecord
//
// Completed ref group is identified if all pertaining poles a either in attached or cancelled states
func (psp *PoleSitesPersister) ArchiveCompletedPoleRefs(psr *PoleSiteRecord) error {
	// init archive PoleSiteRecord
	archivePoleSite := NewPoleSiteRecord()
	archivePoleSite.PoleSite = psr.PoleSite.CloneInfoForArchive()

	nbap := psr.MoveCompletedGroupTo(archivePoleSite.PoleSite)
	if nbap == 0 {
		return nil
	}
	_, err := psp.Update(psr)
	if err != nil {
		return err
	}
	psp.Add(archivePoleSite)
	return nil
}

func (psp *PoleSitesPersister) GetItemizableSites(isSiteVisible items.IsItemizableSiteVisible) []items.ItemizableSite {
	psp.RLock()
	defer psp.RUnlock()

	pss := []items.ItemizableSite{}
	for _, psr := range psp.polesites {
		if isSiteVisible(psr) {
			pss = append(pss, psr)
		}
	}
	return pss
}

func (psp *PoleSitesPersister) GetItemizableSiteById(id int) items.ItemizableSite {
	site := psp.GetById(id)
	if site == nil {
		return nil // force untyped nil to enable (return == nil) test
	}
	return site
}

func (psp *PoleSitesPersister) GetAllSites() []archives.ArchivableRecord {
	psp.RLock()
	defer psp.RUnlock()

	archivableSites := make([]archives.ArchivableRecord, len(psp.polesites))
	for i, site := range psp.polesites {
		archivableSites[i] = site
	}
	return archivableSites
}

func (psp *PoleSitesPersister) GetName() string {
	return "Polesites"
}
