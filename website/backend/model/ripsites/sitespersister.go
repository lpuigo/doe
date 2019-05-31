package ripsites

import (
	"archive/zip"
	"fmt"
	"github.com/lpuig/ewin/doe/website/backend/model/clients"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	"github.com/lpuig/ewin/doe/website/backend/model/items"
	"github.com/lpuig/ewin/doe/website/backend/persist"
	rs "github.com/lpuig/ewin/doe/website/frontend/model/ripsite"
	"io"
	"path/filepath"
	"sort"
	"strings"
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

func sortedSetKeys(set map[string]int) []string {
	res := []string{}
	for key, _ := range set {
		res = append(res, key)
	}
	sort.Strings(res)
	return res
}

// GetStats returns all Stats about all contained RipsiteRecords visible with isWSVisible = true and IsTeamVisible = true
func (sp *SitesPersister) GetStats(maxVal int, dateFor date.DateAggreg, isRSVisible IsSiteVisible, isTeamVisible clients.IsTeamVisible, clientByName clients.ClientByName, showTeam bool, showprice bool) *rs.RipsiteStats {
	sp.RLock()
	defer sp.RUnlock()

	// calc per Team/date/indicator values
	calcValues := make(map[items.StatKey]float64)
	for _, sr := range sp.sites {
		if isRSVisible(sr.Site) {
			client := clientByName(sr.Site.Client)
			if client == nil {
				continue
			}
			sr.AddStat(calcValues, dateFor, isTeamVisible, client.Bpu, client.GenTeamNameByMember(), showprice)
		}
	}

	//create client-team, sites, Series & dates Lists
	end := date.Today()
	start := end.String()
	teamset := make(map[string]int)
	serieset := make(map[string]int)
	siteset := make(map[string]int)
	agrValues := make(map[items.StatKey]float64)
	for key, val := range calcValues {
		teamset[key.Team] = 1
		serieset[key.Serie] = 1
		siteset[key.Site] = 1
		if key.Date < start {
			start = key.Date
		}
		agrValues[items.StatKey{
			Team:    key.Team,
			Date:    key.Date,
			Site:    key.Site,
			Article: "",
			Serie:   key.Serie,
		}] += val
	}
	teams := []string{}
	for t, _ := range teamset {
		// if showTeam is false, only show client sum-up
		if !(!showTeam && strings.Contains(t, " : ")) {
			teams = append(teams, t)
		}
	}
	sort.Strings(teams)

	series := sortedSetKeys(serieset)
	sites := sortedSetKeys(siteset)

	dateset := make(map[string]int)
	curStringDate := dateFor(date.DateFrom(start).String())
	curDate := date.DateFrom(curStringDate)
	endStringDate := dateFor(end.String())
	endReached := false
	for !endReached {
		dateset[curStringDate] = 1
		curDate = curDate.AddDays(7)
		curStringDate = dateFor(curDate.String())
		endReached = curStringDate > endStringDate
	}
	dates := sortedSetKeys(dateset)
	// keep maxVal newest data
	if len(dates) > maxVal {
		dates = dates[len(dates)-maxVal:]
	}

	ws := rs.NewBERipsiteStats()
	//ws.Values : map{series}[#team]{sites}[#date]float64
	ws.Dates = dates

	for _, teamName := range teams {
		teamActivity := 0.0
		values := make(map[string]map[string][]float64)
		for _, serie := range series {
			values[serie] = make(map[string][]float64)
			for _, site := range sites {
				siteData := make([]float64, len(dates))
				siteActivity := 0.0
				for dateNum, d := range dates {
					val := agrValues[items.StatKey{
						Team:    teamName,
						Date:    d,
						Site:    site,
						Article: "",
						Serie:   serie,
					}]
					teamActivity += val
					siteActivity += val
					siteData[dateNum] = val
				}
				if siteActivity > 0 {
					values[serie][site] = siteData
				}
			}
		}
		if teamActivity == 0 {
			// current team as no activity on the time laps, skip it
			continue
		}
		ws.Teams = append(ws.Teams, teamName)
		for _, serie := range series {
			ws.Values[serie] = append(ws.Values[serie], values[serie])
		}
	}
	return ws
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
