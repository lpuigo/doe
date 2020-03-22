package actorscalendar

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

type CalendarDayInfo struct {
	*js.Object

	Class   string `js:"Class"`
	Comment string `js:"Comment"`
}

func NewCalendarDayInfo() *CalendarDayInfo {
	cdi := &CalendarDayInfo{Object: tools.O()}
	cdi.Class = ""
	cdi.Comment = ""
	return cdi
}
