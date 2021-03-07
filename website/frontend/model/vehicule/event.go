package vehicule

import (
	"github.com/gopherjs/gopherjs/js"
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

func (e *Event) Copy() *Event {
	return EventFromJS(json.Parse(json.Stringify(e.Object)))
}
