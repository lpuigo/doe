package ripsites

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

type SitesPersister struct {
	sync.RWMutex
	persister *persist.Persister

	sites []*SiteRecord
}

func NewSitesPersit(dir string) (*SitesPersister, error) {
	sp := &SitesPersister{
		persister: persist.NewPersister("Ripsites", dir),
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

	sp.persister.Reinit()
	sp.sites = []*SiteRecord{}

	files, err := sp.persister.GetFilesList("deleted")
	if err != nil {
		return fmt.Errorf("could not get files from sites persister: %v", err)
	}

	for _, file := range files {
		wsr, err := NewSiteRecordFromFile(file)
		if err != nil {
			return fmt.Errorf("could not create site from '%s': %v", filepath.Base(file), err)
		}
		err = sp.persister.Load(wsr)
		if err != nil {
			return fmt.Errorf("error while loading %s: %s", file, err.Error())
		}
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

//func sortedSetKeys(set map[string]int) []string {
//	res := []string{}
//	for key, _ := range set {
//		res = append(res, key)
//	}
//	sort.Strings(res)
//	return res
//}
//

// getSitesItems returns items for all visible ripsites
func (sp *SitesPersister) getSitesItems(isRSVisible IsSiteVisible, clientByName clients.ClientByName) ([]*items.Item, error) {
	items := []*items.Item{}
	for _, sr := range sp.sites {
		if isRSVisible(sr.Site) {
			client := clientByName(sr.Site.Client)
			if client == nil {
				continue
			}
			siteItems, err := sr.Itemize(client.Bpu)
			if err != nil {
				return nil, fmt.Errorf("error on ripsite stat itemize for '%s':%s", sr.Ref, err.Error())
			}
			items = append(items, siteItems...)
		}
	}
	return items, nil
}

// GetProdStats returns all Stats about all contained RipsiteRecords visible with isWSVisible = true and IsTeamVisible = true
func (sp *SitesPersister) GetProdStats(sc items.StatContext, isRSVisible IsSiteVisible, clientByName clients.ClientByName, actorById clients.ActorById, showprice bool, groupBy string) (*rs.RipsiteStats, error) {
	sp.RLock()
	defer sp.RUnlock()

	// Build Item List
	sitesItems, err := sp.getSitesItems(isRSVisible, clientByName)
	if err != nil {
		return nil, err
	}

	// create Prod Stats
	stats := items.NewStats()

	addValue := func(site, team, date, serie string, actors []string, value float64) {
		stats.AddStatValue(site, team, date, "", serie, value)
		if sc.ShowTeam && len(actors) > 0 {
			value /= float64(len(actors))
			for _, actName := range actors {
				stats.AddStatValue(site, team+" : "+actName, date, "", serie, value)
			}
		}
	}

	for _, item := range sitesItems {
		if !item.Done {
			continue
		}
		actorsName := make([]string, len(item.Actors))
		for i, actId := range item.Actors {
			actorsName[i] = actorById(actId)
		}
		dateFor := sc.DateFor(item.Date)
		switch groupBy {
		case "activity":
			addValue(item.Activity, item.Client, dateFor, items.StatSerieWork, actorsName, item.Work())
			if showprice {
				addValue(item.Activity, item.Client, dateFor, items.StatSeriePrice, actorsName, item.Price())
			}
		case "site":
			addValue(item.Site, item.Client, dateFor, items.StatSerieWork, actorsName, item.Work())
			if showprice {
				addValue(item.Site, item.Client, dateFor, items.StatSeriePrice, actorsName, item.Price())
			}
		case "mean":
			addValue(item.Activity, item.Client, dateFor, items.StatSerieWork, actorsName, item.Work())
		default:
			return nil, fmt.Errorf("unsupported groupBy value '%s'", groupBy)
		}
	}

	// Aggregate Stats
	d1 := func(s items.StatKey) string { return s.Serie }
	d2 := func(s items.StatKey) string { return s.Team }
	d3 := func(s items.StatKey) string { return s.Site }
	f1 := items.KeepAll
	f2 := items.KeepAll
	//f2 := func(e string) bool { return !(!sc.ShowTeam && strings.Contains(e, " : ")) }
	f3 := items.KeepAll
	aggrStats := stats.Aggregate(sc, d1, d2, d3, f1, f2, f3)

	return aggrStats, nil
}

// GetProdStats returns all Stats about all contained RipsiteRecords visible with isWSVisible = true and IsTeamVisible = true
func (sp *SitesPersister) GetMeanProdStats(sc items.StatContext, isRSVisible IsSiteVisible, clientByName clients.ClientByName, actorInfoById clients.ActorInfoById) (*rs.RipsiteStats, error) {
	sp.RLock()
	defer sp.RUnlock()

	// Build Item List
	sitesItems, err := sp.getSitesItems(isRSVisible, clientByName)
	if err != nil {
		return nil, err
	}

	// create Prod Stats
	stats := items.NewStats()

	for _, item := range sitesItems {
		if !item.Done {
			continue
		}
		dateFor := sc.DateFor(item.Date)

		nbActors := float64(len(item.Actors))
		var work float64
		if nbActors > 0 {
			work = item.Work() / nbActors
		}
		for _, actId := range item.Actors {
			actInfos := actorInfoById(actId)
			if len(actInfos) < 2 {
				continue
			}
			actRole, actName := actInfos[0], actInfos[1]

			client := item.Client
			clientRole := item.Client + " : " + actRole
			clientRoleName := item.Client + " : " + actRole + " / " + actName

			stats.AddStatValue(actRole, client, dateFor, "", items.StatSerieWork, work)
			stats.AddStatValue(actRole, clientRole, dateFor, "", items.StatSerieWork, work)
			stats.AddStatValue(actRole, clientRoleName, dateFor, "", items.StatSerieWork, work)
		}
	}

	// Aggregate Stats
	d1 := func(s items.StatKey) string { return s.Serie }
	d2 := func(s items.StatKey) string { return s.Team }
	d3 := func(s items.StatKey) string { return s.Site }
	f1 := items.KeepAll
	f2 := items.KeepAll
	//f2 := func(e string) bool { return !(!sc.ShowTeam && strings.Contains(e, " : ")) }
	f3 := items.KeepAll
	aggrStats := stats.Aggregate(sc, d1, d2, d3, f1, f2, f3)

	return aggrStats, nil
}

func (sp *SitesPersister) GetAllItems(firstDate string, dateFor date.DateAggreg, isRSVisible IsSiteVisible, clientByName clients.ClientByName) ([]*items.Item, error) {

	// Build Item List
	sitesItems, err := sp.getSitesItems(isRSVisible, clientByName)
	if err != nil {
		return nil, err
	}

	res := []*items.Item{}
	for _, itm := range sitesItems {
		if !(itm.Todo && itm.Done) {
			continue
		}
		if dateFor(itm.Date) < firstDate {
			continue
		}
		itm.Date = dateFor(itm.Date)
		res = append(res, itm.SplitByActors()...)
	}
	return res, nil
}

func (sp *SitesPersister) GetAllSites() []archives.ArchivableRecord {
	sp.RLock()
	defer sp.RUnlock()

	archivableSites := make([]archives.ArchivableRecord, len(sp.sites))
	for i, site := range sp.sites {
		archivableSites[i] = site
	}
	return archivableSites
}

func (sp *SitesPersister) GetName() string {
	return "Ripsites"
}
