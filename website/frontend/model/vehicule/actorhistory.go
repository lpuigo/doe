package vehicule

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/date"
)

// Type ActorHistory reflects ewin/doe/website/backend/model/vehicules.ActorHistory
type ActorHistory struct {
	*js.Object

	Date    string `js:"Date"`
	ActorId int    `js:"ActorId"`
}

func NewActorHistory() *ActorHistory {
	nah := &ActorHistory{Object: tools.O()}
	nah.Date = date.TodayAfter(0)
	nah.ActorId = -1
	return nah
}

func (ah ActorHistory) Copy() *ActorHistory {
	nah := NewActorHistory()
	nah.Date = ah.Date
	nah.ActorId = ah.ActorId
	return nah
}

func CompareActorHistory(a, b ActorHistory) int {
	if a.Date > b.Date {
		return -1
	}
	if a.Date == b.Date {
		return 0
	}
	return 1
}
