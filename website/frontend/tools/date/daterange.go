package date

import "github.com/gopherjs/gopherjs/js"

// Type DateRange reflects ewin/doe/website/backend/model/date.DateRange
type DateRange struct {
	*js.Object

	Begin string `js:"Begin"`
	End   string `js:"End"`
}
