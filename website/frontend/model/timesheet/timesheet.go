package timesheet

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

type ActorsTime struct {
	*js.Object

	Hours []int `js:"Hours"`
}

func NewActorTime() *ActorsTime {
	at := &ActorsTime{Object: tools.O()}
	at.Hours = make([]int, 6)
	return at
}

func (at *ActorsTime) SetActiveWeek(active []int) {
	for i := 0; i < 5; i++ {
		if active[i] >= 1 {
			at.Get("Hours").Call("splice", i, 1, 7)
		} else {
			at.Get("Hours").Call("splice", i, 1, 0)
		}
	}
}

func (at *ActorsTime) Equals(oat *ActorsTime) bool {
	for i, value := range at.Hours {
		if oat.Hours[i] != value {
			return false
		}
	}
	return true
}

// NbActiveDays returns the numbers of days with recorded hours (hours > 0) (Saturdays are ignored)
func (at *ActorsTime) NbActiveDays() int {
	ad := 0
	for i, hour := range at.Hours {
		if hour > 0 && i < 5 {
			ad++
		}
	}
	return ad
}

type TimeSheet struct {
	*js.Object

	Id          int                 `js:"Id"`
	WeekDate    string              `js:"WeekDate"`
	ActorsTimes map[int]*ActorsTime `js:"ActorsTimes"`
}

func NewTimeSheet() *TimeSheet {
	ts := &TimeSheet{Object: tools.O()}
	ts.Id = -1
	ts.WeekDate = ""
	ts.ActorsTimes = map[int]*ActorsTime{}
	return ts
}

func TimeSheetFromJS(obj *js.Object) *TimeSheet {
	return &TimeSheet{Object: obj}
}

// Clone returns a copy of given TimeSheet without any ActorsTime
func (ts *TimeSheet) Clone() *TimeSheet {
	nts := NewTimeSheet()
	nts.Id = ts.Id
	nts.WeekDate = ts.WeekDate

	return nts
}

func (ts *TimeSheet) AddUpdatedActorsTimes(refTs, currentTs *TimeSheet) {
	ats := make(map[int]*ActorsTime)
	for id, cat := range currentTs.ActorsTimes {
		rat, found := refTs.ActorsTimes[id]
		if !(found && rat.Equals(cat)) {
			ats[id] = cat
		}
	}
	ts.ActorsTimes = ats
}
