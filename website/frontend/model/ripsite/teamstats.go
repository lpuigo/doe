package ripsite

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

type TeamStats struct {
	*js.Object

	Team         string                          `js:"Team"`
	Dates        []string                        `js:"Dates"`
	Values       map[string]map[string][]float64 `js:"Values"` // {serie}{site}[#date]float64
	IsClientTeam bool                            `js:"IsClientTeam"`
	HasTeams     bool                            `js:"HasTeams"`
	ShowTeams    bool                            `js:"ShowTeams"`
}

func NewTeamStats() *TeamStats {
	ts := &TeamStats{Object: tools.O()}
	ts.Team = ""
	ts.Dates = []string{}
	ts.Values = make(map[string]map[string][]float64)
	ts.IsClientTeam = false
	ts.HasTeams = false
	ts.ShowTeams = false
	return ts
}
