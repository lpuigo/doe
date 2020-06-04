package group

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

// Type Group reflects ewin/doe/website/backend/model/groups.Group
type Group struct {
	*js.Object

	Id               int     `js:"Id"`
	Name             string  `js:"Name"`
	ActorDailyWork   float64 `js:"ActorDailyWork"`
	ActorDailyIncome float64 `js:"ActorDailyIncome"`
}

func NewGroup() *Group {
	g := &Group{Object: tools.O()}
	g.Id = -1
	g.Name = ""
	g.ActorDailyIncome = 0
	g.ActorDailyWork = 0
	return g
}

func GroupFromJS(obj *js.Object) *Group {
	return &Group{Object: obj}
}
