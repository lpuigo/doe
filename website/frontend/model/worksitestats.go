package model

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

type WorksiteStats struct {
	*js.Object

	StartDate string   `js:"StartDate"`
	Teams     []string `js:"Teams"`
	NbEls     [][]int  `js:"NbEls"`
}

func NewWorksiteStats() *WorksiteStats {
	ws := &WorksiteStats{Object: tools.O()}
	ws.StartDate = ""
	ws.Teams = nil
	ws.NbEls = nil
	return ws
}

func NewBEWorksiteStats() *WorksiteStats {
	return &WorksiteStats{}
}

func WorksiteStatsFromJs(o *js.Object) *WorksiteStats {
	return &WorksiteStats{Object: o}
}

func (ws *WorksiteStats) CreateTeamStats() []*TeamStats {
	res := []*TeamStats{}
	for i, team := range ws.Teams {
		//if i == 0 {
		//	continue // Skip first Team as it is Worksites global
		//}
		ts := NewTeamStats()
		ts.Team = team
		ts.StartDate = ws.StartDate
		ts.NbEls = ws.NbEls[i]
		res = append(res, ts)
	}
	return res
}
