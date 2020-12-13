package worksites

import (
	"fmt"
	"github.com/lpuig/ewin/doe/model"
	"github.com/lpuig/ewin/doe/website/backend/model/archives"
	"github.com/lpuig/ewin/doe/website/backend/model/clients"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	"github.com/lpuig/ewin/doe/website/backend/persist"
	"github.com/lpuig/ewin/doe/website/frontend/model/worksite"
	"path/filepath"
	"sort"
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
		persister: persist.NewPersister("Worksites", dir),
	}
	err := wsp.persister.CheckDirectory()
	if err != nil {
		return nil, err
	}
	wsp.persister.SetPersistDelay(1 * time.Second)
	return wsp, nil
}

func (wsp WorkSitesPersister) NbWorksites() int {
	return len(wsp.workSites)
}

// LoadDirectory loads all persisted Worksite Records
func (wsp *WorkSitesPersister) LoadDirectory() error {
	wsp.Lock()
	defer wsp.Unlock()

	wsp.persister.Reinit()
	wsp.workSites = []*WorkSiteRecord{}

	files, err := wsp.persister.GetFilesList("deleted")
	if err != nil {
		return fmt.Errorf("could not get files from worksites persister: %v", err)
	}

	for _, file := range files {
		wsr, err := NewWorkSiteRecordFromFile(file)
		if err != nil {
			return fmt.Errorf("could not create worksite from '%s': %v", filepath.Base(file), err)
		}
		err = wsp.persister.Load(wsr)
		if err != nil {
			return fmt.Errorf("error while loading %s: %s", file, err.Error())
		}
		wsp.workSites = append(wsp.workSites, wsr)
	}
	return nil
}

// GetAll returns all contained WorkSiteRecords for which keep(wsr.Worksite) == true
func (wsp WorkSitesPersister) GetAll(keep model.IsWSVisible) []*WorkSiteRecord {
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

	// Set the Update Date
	wsr.Worksite.UpdateDate = date.Today().String()

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
	owsr.Worksite.UpdateDate = date.Today().String()
	wsp.persister.MarkDirty(owsr)
	return nil
}

// Remove removes the given WorkSiteRecord from the WorkSitesPersister (pertaining file is moved to deleted dir)
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
func (wsp *WorkSitesPersister) GetById(id int) *WorkSiteRecord {
	wsp.RLock()
	defer wsp.RUnlock()

	for _, wsr := range wsp.workSites {
		if wsr.Id == id {
			return wsr
		}
	}
	return nil
}

func (wsp *WorkSitesPersister) getVisibleWorksitesStats(dateFor date.DateAggreg, isWSVisible model.IsWSVisible, isTeamVisible clients.IsTeamVisible, clientByName clients.ClientByName, showTeam, calcToDo bool) map[model.StatKey]int {
	nbEls := make(map[model.StatKey]int)
	for _, wsr := range wsp.workSites {
		if isWSVisible(wsr.Worksite) {
			client := clientByName(wsr.Worksite.Client)
			if client == nil {
				continue
			}
			wsr.AddStat(nbEls, dateFor, isTeamVisible, client.GenTeamNameByMember(), showTeam, calcToDo)
		}
	}
	return nbEls
}

// GetStats returns all Stats about all contained WorkSiteRecords visible with isWSVisible = true and IsTeamVisible = true
//
// maxVal most recent dates are returned if maxVal > 0 (otherwise, all dates a returned)
func (wsp *WorkSitesPersister) GetStats(maxVal int, dateFor date.DateAggreg, isWSVisible model.IsWSVisible, isTeamVisible clients.IsTeamVisible, clientByName clients.ClientByName, showTeam, calcToDo bool) *worksite.WorksiteStats {
	wsp.RLock()
	defer wsp.RUnlock()

	// calc Nb installed ELs per Team/date/mesurement
	nbEls := wsp.getVisibleWorksitesStats(dateFor, isWSVisible, isTeamVisible, clientByName, showTeam, calcToDo)

	ws := worksite.NewBEWorksiteStats()

	//create client, team, measurments & dates Lists
	end := date.Today()
	start := end.String()
	teamset := make(map[string]int)
	messet := make(map[string]int)
	for key, _ := range nbEls {
		teamset[key.Team] = 1
		messet[key.Mes] = 1
		if key.Date < start {
			start = key.Date
		}
	}
	teams := []string{}
	for t, _ := range teamset {
		teams = append(teams, t)
	}
	sort.Strings(teams)

	measurements := []string{}
	for m, _ := range messet {
		measurements = append(measurements, m)
	}
	sort.Strings(measurements)

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
	dates := []string{}
	for d, _ := range dateset {
		dates = append(dates, d)
	}
	sort.Strings(dates)
	// keep maxVal newest data
	if maxVal > 0 && len(dates) > maxVal {
		dates = dates[len(dates)-maxVal:]
	}

	ws.Values = make(map[string][][]int)
	ws.Dates = dates

	for _, teamName := range teams {
		teamActivity := 0
		values := make(map[string][]int)
		for _, meas := range measurements {
			values[meas] = make([]int, len(dates))
			for dateNum, d := range dates {
				nbEl := nbEls[model.StatKey{Team: teamName, Date: d, Mes: meas}]
				teamActivity += nbEl
				values[meas][dateNum] = nbEl
			}
		}
		if teamActivity == 0 {
			// current team as no activity on the time laps, skip it
			continue
		}
		ws.Teams = append(ws.Teams, teamName)
		for _, meas := range measurements {
			ws.Values[meas] = append(ws.Values[meas], values[meas])
		}
	}
	return ws
}

// GetStockStats returns all Stock Stats about all contained WorkSiteRecords visible with isWSVisible = true and IsTeamVisible = true
func (wsp *WorkSitesPersister) GetStockStats(maxVal int, dateFor date.DateAggreg, isWSVisible model.IsWSVisible, isTeamVisible clients.IsTeamVisible, clientByName clients.ClientByName) *worksite.WorksiteStats {
	prodStats := wsp.GetStats(-1, dateFor, isWSVisible, isTeamVisible, clientByName, false, true)
	dates := prodStats.Dates
	// keep maxVal newest data
	if len(dates) > maxVal {
		dates = dates[len(dates)-maxVal:]
	}
	stockStats := &worksite.WorksiteStats{
		Dates:  dates,
		Teams:  prodStats.Teams,
		Values: make(map[string][][]int),
	}
	for _, value := range []string{worksite.NbElsToInstall, worksite.NbElsToMeasure, worksite.NbElsToDOE, worksite.NbElsToBill} {
		stockStats.Values[value] = make([][]int, len(prodStats.Teams))
	}

	addSlice := func(slA, slB []int) []int {
		res := make([]int, len(slA))
		for i, a := range slA {
			res[i] = a + slB[i]
		}
		return res
	}

	subSlice := func(slA, slB []int) []int {
		res := make([]int, len(slA))
		for i, a := range slA {
			res[i] = a - slB[i]
		}
		return res
	}

	cumSlice := func(sl []int) []int {
		res := make([]int, len(sl))
		for i, a := range sl {
			v := a
			if i > 0 {
				v += res[i-1]
			}
			res[i] = v
		}
		// keep maxVal newest data
		if len(res) > maxVal {
			res = res[len(res)-maxVal:]
		}
		return res
	}

	for teamNum, _ := range prodStats.Teams {
		submitted := prodStats.Values[worksite.NbElsSumitted][teamNum]
		installed := prodStats.Values[worksite.NbElsInstalled][teamNum]
		done := addSlice(installed, prodStats.Values[worksite.NbElsBlocked][teamNum])
		measured := prodStats.Values[worksite.NbElsMeasured][teamNum]
		DOEed := prodStats.Values[worksite.NbElsDOE][teamNum]
		Billed := prodStats.Values[worksite.NbElsBilled][teamNum]

		stockStats.Values[worksite.NbElsToInstall][teamNum] = cumSlice(subSlice(submitted, done))
		stockStats.Values[worksite.NbElsToMeasure][teamNum] = cumSlice(subSlice(installed, measured))
		stockStats.Values[worksite.NbElsToDOE][teamNum] = cumSlice(subSlice(measured, DOEed))
		stockStats.Values[worksite.NbElsToBill][teamNum] = cumSlice(subSlice(DOEed, Billed))
	}

	return stockStats
}

func (wsp *WorkSitesPersister) GetAllSites() []archives.ArchivableRecord {
	wsp.RLock()
	defer wsp.RUnlock()

	archivableSites := make([]archives.ArchivableRecord, len(wsp.workSites))
	for i, site := range wsp.workSites {
		archivableSites[i] = site
	}
	return archivableSites
}

func (wsp *WorkSitesPersister) GetName() string {
	return "Worksites"
}
