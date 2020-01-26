package worksite

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"strings"
)

const (
	NbElsSumitted  string = "Submitted"
	NbElsInstalled string = "Installed"
	NbElsBlocked   string = "Blocked"
	NbElsMeasured  string = "Measured"
	NbElsDOE       string = "DOE"
	NbElsBilled    string = "Billed"

	NbElsToInstall string = "ToInstall"
	NbElsToMeasure string = "ToMeasure"
	NbElsToDOE     string = "ToDOE"
	NbElsToBill    string = "ToBill"
)

type WorksiteStats struct {
	*js.Object

	Dates  []string           `js:"Dates"`
	Teams  []string           `js:"Teams"`
	Values map[string][][]int `js:"Values"` // map[measurement][teamNum][dateNum]int{}
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
	res := make([]*TeamStats, len(ws.Teams))
	currentTeam := ""
	prevTS := NewTeamStats()
	for i, team := range ws.Teams {
		isTeamClient := false
		if !(currentTeam != "" && strings.HasPrefix(team, currentTeam)) {
			currentTeam = team
			isTeamClient = true
		} else {
			prevTS.HasTeams = true
		}
		ts := NewTeamStats()
		if isTeamClient {
			prevTS = ts
		}
		ts.Team = team
		ts.Dates = ws.Dates
		ts.IsClientTeam = isTeamClient
		for mes, _ := range ws.Values {
			//ts.Values[mes] = ws.Values[mes][i]
			ts.Get("Values").Set(mes, ws.Values[mes][i])
		}
		res[i] = ts
	}
	return res
}
