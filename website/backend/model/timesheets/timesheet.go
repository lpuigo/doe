package timesheets

import "github.com/lpuig/ewin/doe/website/frontend/tools/date"

type ActorsTime struct {
	Hours []int
}

func NewActorsTime() *ActorsTime {
	return &ActorsTime{
		Hours: make([]int, 6),
	}
}

func NewMonthlyActorsTime() *ActorsTime {
	return &ActorsTime{
		Hours: make([]int, 31),
	}
}

func (at *ActorsTime) Clone() *ActorsTime {
	return &ActorsTime{
		Hours: at.Hours[:],
	}
}

type TimeSheet struct {
	Id          int
	WeekDate    string
	ActorsTimes map[int]*ActorsTime
}

func NewTimeSheet(weekdate string) *TimeSheet {
	ts := &TimeSheet{
		Id:          -1,
		WeekDate:    weekdate,
		ActorsTimes: make(map[int]*ActorsTime),
	}
	return ts
}

// NewTimeSheetForActorsIds returns a new TimeSheet with empty ActorsTimes for each given Ids
func NewTimeSheetForActorsIds(weekdate string, ids []int) *TimeSheet {
	ts := NewTimeSheet(weekdate)
	for _, id := range ids {
		ts.ActorsTimes[id] = NewActorsTime()
	}
	return ts
}

// NewMonthlyTimeSheetForActorsIds returns a new TimeSheet with empty Monthly ActorsTimes for each given Ids
func NewMonthlyTimeSheetForActorsIds(monthDate string, ids []int) *TimeSheet {
	ts := NewTimeSheet(monthDate)
	for _, id := range ids {
		ts.ActorsTimes[id] = NewMonthlyActorsTime()
	}
	return ts
}

// CloneForActorIds returns a TimeSheet with the same date, and filled with given ids ActorsTime
//
// for already registered ids, ActorsTimes are copied, for others (unregistered yet) Ids new ActorsTimes are created
func (ts *TimeSheet) CloneForActorIds(ids []int) *TimeSheet {
	nts := NewTimeSheet(ts.WeekDate)
	for _, id := range ids {
		at, found := ts.ActorsTimes[id]
		if !found {
			nts.ActorsTimes[id] = NewActorsTime()
			continue
		}
		nts.ActorsTimes[id] = at.Clone()
	}
	return nts
}

func (ts *TimeSheet) UpdateActorsTimesFrom(uts *TimeSheet) {
	for id, at := range uts.ActorsTimes {
		ts.ActorsTimes[id] = at
	}
}

// Merge places ots.actorsTime hours in ts.actortime (ids from ots not defined in ts are skipped)
func (ts *TimeSheet) Merge(ots *TimeSheet) {
	offsetDay := int(date.NbDaysBetween(ts.WeekDate, ots.WeekDate))
	for id, at := range ots.ActorsTimes {
		currentActorTime, exist := ts.ActorsTimes[id]
		if !exist {
			continue
		}
		//ts.ActorsTimes[id].MergeAtPos(numDay, at.Hours)
		for oNumDay, hours := range at.Hours {
			i := oNumDay + offsetDay
			if !(i >= 0 && i < len(currentActorTime.Hours)) {
				continue
			}
			currentActorTime.Hours[i] = hours
		}
	}
}
