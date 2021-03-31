package actorscalendar

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/model/actor"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/date"
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

type GroupActor struct {
	*js.Object
	Actor      *actor.Actor    `js:"Actor"`
	GroupName  string          `js:"GroupName"`
	Assignment *date.DateRange `js:"Assignment"`
}

func NewGroupActor(actor *actor.Actor, group string, assign *date.DateRange) *GroupActor {
	ga := &GroupActor{Object: tools.O()}
	ga.Actor = actor
	ga.GroupName = group
	ga.Assignment = assign

	return ga
}

func GroupActorFromJS(o *js.Object) *GroupActor {
	return &GroupActor{Object: o}
}
