package polesites

import (
	"fmt"
	"github.com/lpuig/ewin/doe/website/backend/model/archives"
	"github.com/lpuig/ewin/doe/website/backend/model/clients"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	"github.com/lpuig/ewin/doe/website/backend/model/items"
	"github.com/lpuig/ewin/doe/website/backend/persist"
	rs "github.com/lpuig/ewin/doe/website/frontend/model/ripsite"
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

// GetStats returns all Stats about all contained RipsiteRecords visible with isWSVisible = true and IsTeamVisible = true
func (psp *PoleSitesPersister) GetStats(sc items.StatContext, isPSVisible IsPolesiteVisible, clientByName clients.ClientByName, actorById clients.ActorById, showprice bool) (*rs.RipsiteStats, error) {
	psp.RLock()
	defer psp.RUnlock()

	// calc per Team/date/indicator values
	calcValues := make(items.Stats)
	for _, pr := range psp.polesites {
		if isPSVisible(pr.PoleSite) {
			client := clientByName(pr.PoleSite.Client)
			if client == nil {
				continue
			}
			err := pr.AddStat(calcValues, sc, actorById, client.Bpu, showprice)
			if err != nil {
				return nil, err
			}
		}
	}

	d1 := func(s items.StatKey) string { return s.Serie } // Bars Family
	d2 := func(s items.StatKey) string { return s.Team }  // Graphs
	d3 := func(s items.StatKey) string { return s.Site }  // side block
	f1 := items.KeepAll
	//f2 := func(e string) bool { return !(!sc.ShowTeam && strings.Contains(e, " : ")) }
	f2 := items.KeepAll
	f3 := items.KeepAll
	return calcValues.Aggregate(sc, d1, d2, d3, f1, f2, f3), nil
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
