package model

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

type TeamStats struct {
	*js.Object

	Team      string `js:"Team"`
	StartDate string `js:"StartDate"`
	NbEls     []int  `js:"NbEls"`
}

func NewTeamStats() *TeamStats {
	ts := &TeamStats{Object: tools.O()}
	ts.Team = ""
	ts.StartDate = ""
	ts.NbEls = []int{}
	return ts
}
