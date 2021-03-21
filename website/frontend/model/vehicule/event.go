package vehicule

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/model/vehicule/vehiculeconst"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/date"
	"github.com/lpuig/ewin/doe/website/frontend/tools/json"
)

// Type Event reflects ewin/doe/website/backend/model/vehicules.Event
type Event struct {
	*js.Object

	StartDate string `js:"StartDate"`
	EndDate   string `js:"EndDate"`
	Type      string `js:"Type"`
	Comment   string `js:"Comment"`
}

func EventFromJS(obj *js.Object) *Event {
	return &Event{Object: obj}
}

func NewEvent() *Event {
	e := &Event{Object: tools.O()}
	e.StartDate = date.TodayAfter(0)
	e.EndDate = ""
	e.Type = vehiculeconst.EventTypeMisc
	e.Comment = ""
	return e
}

func (e *Event) Copy() *Event {
	return EventFromJS(json.Parse(json.Stringify(e.Object)))
}

func CompareEventDate(a, b *Event) int {
	if a.StartDate > b.StartDate {
		return -1
	}
	if a.StartDate == b.StartDate {
		return 0
	}
	return 1
}
