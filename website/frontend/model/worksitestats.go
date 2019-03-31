package model

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

type WorksiteStats struct {
	*js.Object

	Dates  []string           `js:"Dates"`
	Teams  []string           `js:"Teams"`
	Values map[string][][]int `js:"Values"`
}

func NewWorksiteStats() *WorksiteStats {
	ws := &WorksiteStats{Object: tools.O()}
	ws.Dates = []string{}
	ws.Teams = []string{}
	ws.Values = map[string][][]int{}
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
		ts := NewTeamStats()
		ts.Team = team
		ts.Dates = ws.Dates
		for mes, _ := range ws.Values {
			//ts.Values[mes] = ws.Values[mes][i]
			ts.Get("Values").Set(mes, ws.Values[mes][i])
		}
		res = append(res, ts)
	}
	return res
}
