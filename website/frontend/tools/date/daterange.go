package date

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

// Type DateRange reflects ewin/doe/website/backend/model/date.DateRange
type DateRange struct {
	*js.Object

	Begin string `js:"Begin"`
	End   string `js:"End"`
}

func NewDateRange() *DateRange {
	dr := &DateRange{Object: tools.O()}
	dr.Begin = ""
	dr.End = ""
	return dr
}

func NewDateRangeFrom(beg, end string) *DateRange {
	dr := &DateRange{Object: tools.O()}
	dr.Begin = beg
	dr.End = end
	return dr
}
