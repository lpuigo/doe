package ripsite

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/date"
	"strings"
)

type RipsiteStats struct {
	*js.Object

	Dates  []string                          `js:"Dates"`
	Teams  []string                          `js:"Teams"`
	Sites  map[string]bool                   `js:"Sites"`
	Values map[string][]map[string][]float64 `js:"Values"` // map[measurement][teamNum][site][dateNum]float64{}
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

func (rs *RipsiteStats) CreateTeamStats(showSite map[string]bool) []*TeamStats {
	res := []*TeamStats{}
	totStats := NewTeamStats()
	totStats.Team = "Total"
	totStats.Dates = date.ConvertDates(rs.Dates)
	totStats.IsClientTeam = true
	totValues := make(map[string]map[string][]float64)
	for serie, _ := range rs.Values {
		vals := make([]float64, len(rs.Dates))
		svals := map[string][]float64{}
		svals["Total"] = vals
		totValues[serie] = svals
	}
	currentTeam := ""
	prevTS := NewTeamStats()
	for numTeam, team := range rs.Teams {
		isTeamClient := false
		if !(currentTeam != "" && strings.HasPrefix(team, currentTeam)) {
			currentTeam = team
			isTeamClient = true
		} else {
			prevTS.HasTeams = true
		}
		teamHasData := false
		ts := NewTeamStats()
		if isTeamClient {
			prevTS = ts
		}
		ts.Team = team
		ts.Dates = totStats.Dates
		ts.IsClientTeam = isTeamClient
		for serie, _ := range rs.Values {
			//ts.Values[serie] = rs.Values[serie][numTeam]
			datas := map[string][]float64{}
			totVals := totValues[serie]["Total"]
			for site, data := range rs.Values[serie][numTeam] {
				if showSite[site] {
					teamHasData = true
					datas[site] = data
					if isTeamClient {
						for i, val := range data {
							totVals[i] += val
						}
					}
				}
			}
			ts.Get("Values").Set(serie, datas)
			totValues[serie]["Total"] = totVals
		}
		if teamHasData {
			res = append(res, ts)
		}
	}
	totStats.Values = totValues
	res = append([]*TeamStats{totStats}, res...)
	return res
}
