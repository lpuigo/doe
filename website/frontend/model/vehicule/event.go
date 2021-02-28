package vehicule

import "github.com/gopherjs/gopherjs/js"

// Type Event reflects ewin/doe/website/backend/model/vehicules.Event
type Event struct {
	*js.Object

	StartDate string `js:"StartDate"`
	EndDate   string `js:"EndDate"`
	Type      string `js:"Type"`
	Comment   string `js:"Comment"`
}
