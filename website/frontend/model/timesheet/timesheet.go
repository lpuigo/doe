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

func (at *ActorsTime) SetActiveWeek() {
	at.Get("Hours").Call("splice", 0, 5, 7, 7, 7, 7, 7)
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
