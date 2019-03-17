package worksites

import (
	"archive/zip"
	"fmt"
	"github.com/lpuig/ewin/doe/model"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	"github.com/lpuig/ewin/doe/website/backend/persist"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"io"
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
		persister: persist.NewPersister(dir),
	}
	err := wsp.persister.CheckDirectory()
	if err != nil {
		return nil, err
	}
	wsp.persister.SetPersistDelay(1 * time.Second)
	return wsp, nil
}

func (wsp WorkSitesPersister) NbWorsites() int {
	return len(wsp.workSites)
}

// LoadDirectory loads all persisted Worksite Records
func (wsp *WorkSitesPersister) LoadDirectory() error {
	wsp.Lock()
	defer wsp.Unlock()

	files, err := wsp.persister.GetFilesList("deleted")
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

// GetAll returns all contained WorkSiteRecords for which keep(wsr.Worksite) == true
func (wsp WorkSitesPersister) GetAllInfo(keep func(ws *model.Worksite) bool) []*fm.WorksiteInfo {
	wsp.RLock()
	defer wsp.RUnlock()

	ws := []*fm.WorksiteInfo{}
	for _, wsr := range wsp.workSites {
		if keep(wsr.Worksite) {
			ws = append(ws, wsr.Worksite.GetInfo())
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

// GetStats returns all Stats about all contained WorkSiteRecords such as keep(wsr.Worksite) == true
func (wsp *WorkSitesPersister) GetStats(isWSVisible model.IsWSVisible, isTeamVisible model.IsTeamVisible) *fm.WorksiteStats {
	wsp.RLock()
	defer wsp.RUnlock()

	// calc Nb installed ELs per Team/date
	nbEls := make(map[model.StatKey]int)
	for _, wsr := range wsp.workSites {
		if isWSVisible(wsr.Worksite) {
			wsr.AddStat(nbEls, isTeamVisible)
		}
	}

	ws := fm.NewBEWorksiteStats()

	//create team List & Dates
	end := date.Today().GetMonday()
	start := end.String()
	teamset := make(map[string]int)
	for key, _ := range nbEls {
		teamset[key.Team] = 1
		if key.Date < start {
			start = key.Date
		}
	}
	teams := []string{}
	for t, _ := range teamset {
		teams = append(teams, t)
	}
	sort.Strings(teams)
	teams = append([]string{"GLOBAL"}, teams...)
	dates := []string{}
	for d := date.DateFrom(start); !d.After(end); d = d.AddDays(7) {
		dates = append(dates, d.String())
	}
	ws.StartDate = start
	ws.Teams = teams

	// calc nbEls per teams/date
	ws.NbEls = make([][]int, len(teams))
	ws.NbEls[0] = make([]int, len(dates))
	for i, t := range teams[1:] {
		ws.NbEls[i+1] = make([]int, len(dates))
		for j, d := range dates {
			nbEl := nbEls[model.StatKey{Team: t, Date: d}]
			ws.NbEls[i+1][j] = nbEl
			ws.NbEls[0][j] += nbEl
		}
	}

	return ws
}

// ArchiveName returns the WorksiteArchive file name with today's date
func (wsp WorkSitesPersister) ArchiveName() string {
	return fmt.Sprintf("Worksites %s.zip", date.Today().String())
}

// CreateArchive writes a zipped archive of all contained Worksites files to the given writer
func (wsp *WorkSitesPersister) CreateArchive(writer io.Writer) error {
	wsp.RLock()
	defer wsp.RUnlock()

	zw := zip.NewWriter(writer)

	for _, wsr := range wsp.workSites {
		wfw, err := zw.Create(wsr.GetFileName())
		if err != nil {
			return fmt.Errorf("could not create zip entry for worksite %d", wsr.Id)
		}
		err = wsr.Marshall(wfw)
		if err != nil {
			return fmt.Errorf("could not write zip entry for worksite %d", wsr.Id)
		}
	}

	return zw.Close()
}
