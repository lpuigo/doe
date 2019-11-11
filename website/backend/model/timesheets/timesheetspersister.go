package timesheets

import (
	"fmt"
	"github.com/lpuig/ewin/doe/website/backend/model/actors"
	"github.com/lpuig/ewin/doe/website/backend/persist"
	"path/filepath"
	"sync"
	"time"
)

type TimeSheetsPersister struct {
	sync.RWMutex
	persister *persist.Persister

	timeSheetRecords []*TimeSheetRecord
}

func NewTimeSheetsPersister(dir string) (*TimeSheetsPersister, error) {
	tsp := &TimeSheetsPersister{
		persister:        persist.NewPersister("TimeSheets", dir),
		timeSheetRecords: nil,
	}
	err := tsp.persister.CheckDirectory()
	if err != nil {
		return nil, err
	}
	tsp.persister.SetPersistDelay(1 * time.Second)
	return tsp, nil
}

func (tsp *TimeSheetsPersister) NbTimeSheets() int {
	return len(tsp.timeSheetRecords)
}

// LoadDirectory loads all persisted TimeSheets Records
func (tsp *TimeSheetsPersister) LoadDirectory() error {
	tsp.Lock()
	defer tsp.Unlock()

	tsp.persister.Reinit()
	tsp.timeSheetRecords = []*TimeSheetRecord{}

	files, err := tsp.persister.GetFilesList("deleted")
	if err != nil {
		return fmt.Errorf("could not get files from timesheets persister: %v", err)
	}

	for _, file := range files {
		tsr, err := NewTimeSheetRecordFromFile(file)
		if err != nil {
			return fmt.Errorf("could not instantiate timesheets from '%s': %v", filepath.Base(file), err)
		}
		err = tsp.persister.Load(tsr)
		if err != nil {
			return fmt.Errorf("error while loading %s: %s", file, err.Error())
		}
		tsp.timeSheetRecords = append(tsp.timeSheetRecords, tsr)
	}
	return nil
}

// GetById returns the TimeSheetRecord with given Id (or nil if Id not found)
func (tsp *TimeSheetsPersister) GetById(id int) *TimeSheetRecord {
	tsp.RLock()
	defer tsp.RUnlock()

	for _, tsr := range tsp.timeSheetRecords {
		if tsr.Id == id {
			return tsr
		}
	}
	return nil
}

// GetByWeekDate returns the TimeSheetRecord with given WeekDate (or nil if Id not found)
func (tsp *TimeSheetsPersister) GetByWeekDate(weekDate string) *TimeSheetRecord {
	tsp.RLock()
	defer tsp.RUnlock()

	return tsp.getByWeekDate(weekDate)
}

func (tsp *TimeSheetsPersister) getByWeekDate(weekDate string) *TimeSheetRecord {
	for _, tsr := range tsp.timeSheetRecords {
		if tsr.WeekDate == weekDate {
			return tsr
		}
	}
	return nil
}

// Add adds the given TimeSheetRecord to the TimeSheetsPersister and return its (updated with new id) TimeSheetRecord
func (tsp *TimeSheetsPersister) Add(ntsr *TimeSheetRecord) *TimeSheetRecord {
	tsp.Lock()
	defer tsp.Unlock()

	return tsp.add(ntsr)
}

func (tsp *TimeSheetsPersister) add(ntsr *TimeSheetRecord) *TimeSheetRecord {
	// give the record its new ID
	tsp.persister.Add(ntsr)
	ntsr.Id = ntsr.GetId()
	tsp.timeSheetRecords = append(tsp.timeSheetRecords, ntsr)
	return ntsr
}

// Update updates the given TimeSheetRecord
func (tsp *TimeSheetsPersister) Update(ntsr *TimeSheetRecord) error {
	tsp.RLock()
	defer tsp.RUnlock()

	otsr := tsp.GetById(ntsr.Id)
	if otsr == nil {
		return fmt.Errorf("timesheet id not found")
	}
	otsr.TimeSheet = ntsr.TimeSheet
	tsp.persister.MarkDirty(otsr)
	return nil
}

func (tsp *TimeSheetsPersister) findIndex(tsr *TimeSheetRecord) int {
	for i, rec := range tsp.timeSheetRecords {
		if rec.GetId() == tsr.GetId() {
			return i
		}
	}
	return -1
}

// Remove removes the given TimeSheetRecord from the TimeSheetsPersister (pertaining file is moved to deleted dir)
func (tsp *TimeSheetsPersister) Remove(rtsr *TimeSheetRecord) error {
	tsp.Lock()
	defer tsp.Unlock()

	err := tsp.persister.Remove(rtsr)
	if err != nil {
		return err
	}

	i := tsp.findIndex(rtsr)
	copy(tsp.timeSheetRecords[i:], tsp.timeSheetRecords[i+1:])
	tsp.timeSheetRecords[len(tsp.timeSheetRecords)-1] = nil // or the zero value of T
	tsp.timeSheetRecords = tsp.timeSheetRecords[:len(tsp.timeSheetRecords)-1]
	return nil
}

//GetTimeSheetFor returns timesheet for given week date and actors (only actors active at given date are returned)
func (tsp *TimeSheetsPersister) GetTimeSheetFor(weekdate string, actors []*actors.Actor) (*TimeSheet, error) {
	tsp.RLock()
	defer tsp.RUnlock()

	activeIds := []int{}
	for _, act := range actors {
		if act.IsActiveOnWeek(weekdate) {
			activeIds = append(activeIds, act.Id)
		}
	}
	// seek TimeSheet for weekDate
	tsr := tsp.GetByWeekDate(weekdate)
	if tsr == nil {
		return NewTimeSheetForActorsIds(weekdate, activeIds), nil
	}

	// clone timesheet with given actors ids
	ntsr := tsr.CloneForActorIds(activeIds)
	return ntsr, nil
}

//UpdateTimeSheet updates and persists given timesheet (if no timesheet record for same weekdate if found, creates a new timesheet record)
func (tsp *TimeSheetsPersister) UpdateTimeSheet(uts *TimeSheet) error {
	tsp.Lock()
	defer tsp.Unlock()

	// seek existing TimeSheet for given Week
	tsr := tsp.getByWeekDate(uts.WeekDate)

	if tsr == nil {
		// weekDate not found => persist given timesheet
		tsp.add(NewTimeSheetRecordFromTimeSheet(uts))
		return nil
	}
	// existing timesheet record was found for given weekDate => Update it
	tsr.TimeSheet.UpdateActorsTimesFrom(uts)
	tsp.persister.MarkDirty(tsr)
	return nil
}
