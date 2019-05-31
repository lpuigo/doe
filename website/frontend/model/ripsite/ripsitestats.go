package ripsite

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

type RipsiteStats struct {
	*js.Object

	Dates  []string                          `js:"Dates"`
	Teams  []string                          `js:"Teams"`
	Values map[string][]map[string][]float64 `js:"Values"` // map[measurement][teamNum][site][dateNum]int{}
}

func NewRipsiteStats() *RipsiteStats {
	ws := &RipsiteStats{Object: tools.O()}
	ws.Dates = []string{}
	ws.Teams = []string{}
	ws.Values = map[string][]map[string][]float64{}
	return ws
}

func NewBERipsiteStats() *RipsiteStats {
	return &RipsiteStats{
		Object: nil,
		Dates:  nil,
		Teams:  nil,
		Values: map[string][]map[string][]float64{},
	}
}

func RipsiteStatsFromJs(o *js.Object) *RipsiteStats {
	return &RipsiteStats{Object: o}
}

func (rs *RipsiteStats) CreateTeamStats() []*TeamStats {
	res := []*TeamStats{}
	for i, team := range rs.Teams {
		ts := NewTeamStats()
		ts.Team = team
		ts.Dates = rs.Dates
		for mes, _ := range rs.Values {
			//ts.Values[mes] = rs.Values[mes][i]
			ts.Get("Values").Set(mes, rs.Values[mes][i])
		}
		res = append(res, ts)
	}
	return res
}
