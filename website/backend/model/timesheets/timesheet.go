package timesheets

type ActorsTime struct {
	Hours []int
}

func NewActorsTime() *ActorsTime {
	return &ActorsTime{
		Hours: make([]int, 6),
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
