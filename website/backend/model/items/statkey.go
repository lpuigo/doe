package items

import (
	"github.com/lpuig/ewin/doe/website/backend/model/clients"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	rs "github.com/lpuig/ewin/doe/website/frontend/model/ripsite"
)

const (
	StatSerieWork  string = "Work"
	StatSeriePrice string = "Price"
)

type StatContext struct {
	MaxVal        int
	DateFor       date.DateAggreg
	IsTeamVisible clients.IsTeamVisible
	ShowTeam      bool
}

type StatKey struct {
	Team    string
	Date    string
	Site    string
	Article string
	Serie   string
}

type Stats map[StatKey]float64

func NewStats() Stats {
	return make(Stats)
}

func (s Stats) AddStatValue(site, team, date, article, serie string, value float64) {
	s[StatKey{
		Team:    team,
		Date:    date,
		Site:    site,
		Article: article,
		Serie:   serie,
	}] += value
}

// Aggregate returns struct with ws.Values : map{D1}[#D2]{D3}[#date]float64
//
// Standard usage is ws.Values : map{series}[#team]{sites}[#date]float64
func (s Stats) Aggregate(sc StatContext, d1, d2, d3 func(StatKey) string, f1, f2, f3 func(string) bool) *rs.RipsiteStats {

	//create client-team, sites, Series & dates Lists
	serieset := make(ElemSet) // d1
	teamset := make(ElemSet)  // d2
	siteset := make(ElemSet)  // d3
	dateset := make(ElemSet)

	end := date.Today()
	start := end.String()

	agrValues := make(Stats) // values

	for key, val := range s {
		k1 := d1(key)
		k2 := d2(key)
		k3 := d3(key)
		serieset[k1] = 1 // d1
		teamset[k2] = 1  // d2
		siteset[k3] = 1  // d3
		if key.Date < start {
			start = key.Date
		}
		agrValues[StatKey{
			Serie:   k1, // d1
			Team:    k2, // d2
			Site:    k3, // d3
			Article: "", // d4 always aggregated
			Date:    key.Date,
		}] += val
	}

	curStringDate := sc.DateFor(date.DateFrom(start).String())
	curDate := date.DateFrom(curStringDate)
	endStringDate := sc.DateFor(end.String())
	endReached := false
	for !endReached {
		dateset[curStringDate] = 1
		curDate = curDate.AddDays(7)
		curStringDate = sc.DateFor(curDate.String())
		endReached = curStringDate > endStringDate
	}

	series := serieset.SortedKeys(f1)    // d1
	teams := teamset.SortedKeys(f2)      // d2
	sites := siteset.SortedKeys(f3)      // d3
	dates := dateset.SortedKeys(KeepAll) // d4

	// keep maxVal newest data
	if len(dates) > sc.MaxVal {
		dates = dates[len(dates)-sc.MaxVal:]
	}

	// ws.Values : map{D1}[#D2]{D3}[#D4]float64
	ws := rs.NewBERipsiteStats()
	ws.Dates = dates // d4
	sitesmap := map[string]bool{}
	for _, site := range sites {
		sitesmap[site] = true
	}
	ws.Sites = sitesmap // d3

	for _, teamName := range teams {
		teamActivity := 0.0
		values := make(map[string]map[string][]float64)
		for _, serie := range series {
			values[serie] = make(map[string][]float64)
			for _, site := range sites {
				siteData := make([]float64, len(dates))
				siteActivity := 0.0
				for dateNum, d := range dates {
					val := agrValues[StatKey{
						Team:    teamName,
						Date:    d,
						Site:    site,
						Article: "",
						Serie:   serie,
					}]
					teamActivity += val
					siteActivity += val
					siteData[dateNum] = val
				}
				if siteActivity > 0 {
					values[serie][site] = siteData
				}
			}
		}
		if teamActivity == 0 {
			// current team as no activity on the time laps, skip it // d2
			continue
		}
		ws.Teams = append(ws.Teams, teamName) // d2
		for _, serie := range series {
			ws.Values[serie] = append(ws.Values[serie], values[serie]) // d1
		}
	}

	return ws
}
