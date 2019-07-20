package worksite

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

type TeamStats struct {
	*js.Object

	Team   string           `js:"Team"`
	Dates  []string         `js:"Dates"`
	Values map[string][]int `js:"Values"`
}

func NewTeamStats() *TeamStats {
	ts := &TeamStats{Object: tools.O()}
	ts.Team = ""
	ts.Dates = []string{}
	ts.Values = make(map[string][]int)
	return ts
}
