package items

import (
	"github.com/lpuig/ewin/doe/website/backend/model/clients"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	rs "github.com/lpuig/ewin/doe/website/frontend/model/ripsite"
	"strings"
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

// CalcTeamMean adds in given RipsiteStats mean values for each teams (D2 dimension) : map{D1}[#D2]{D3}[#date]float64
//
// Standard usage is ws.Values : map{series}[#team]{sites}[#date]float64
func CalcTeamMean(aggrStat *rs.RipsiteStats, threshold float64) *rs.RipsiteStats {
	// team index of main team and pertaining role teams
	mainTeamIndex := map[string]int{}            // index of main team
	roleTeamIndex := map[string]int{}            // index of role team
	actorsIndex := map[string]map[string][]int{} // index of mainteam > roleteams > actors

	// Build main Team, role Team and Actors Team Indexes
	for index, teamName := range aggrStat.Teams {
		if !strings.Contains(teamName, ":") {
			mainTeamIndex[teamName] = index
			continue
		}
		if !strings.Contains(teamName, "/") {
			roleTeamIndex[teamName] = index
			continue
		}
		mainteam := strings.TrimSpace(teamName[0:strings.Index(teamName, " : ")])
		roleteam := strings.TrimSpace(teamName[0:strings.Index(teamName, " / ")])
		if actorsIndex[mainteam] == nil {
			actorsIndex[mainteam] = map[string][]int{}
		}
		if actorsIndex[mainteam][roleteam] == nil {
			actorsIndex[mainteam][roleteam] = []int{}
		}
		actorsIndex[mainteam][roleteam] = append(actorsIndex[mainteam][roleteam], index)
	}

	// Calc Nb Actors
	const (
		StatSerieNbActor      string = "NbActorsWork"
		StatSerieRoleMeanData string = "MoyenneWork"
		StatSerieMainMeanData string = "MoyenneWork"
	)

	workData := aggrStat.Values[StatSerieWork]
	nbActors := make([]map[string][]float64, len(aggrStat.Teams))
	mainTeamNbAct := map[string][]float64{}
	roleTeamRole := map[string]string{}

	// for each main teams
	for mainTeamName, mainTeamPos := range mainTeamIndex {
		nbActors[mainTeamPos] = make(map[string][]float64)
		mainNbAct := make([]float64, len(aggrStat.Dates)) // nb of actor for role team
		// for each role team within main team
		for roleTeamName, roleTeamActorsPos := range actorsIndex[mainTeamName] {
			roleNbAct := make([]float64, len(aggrStat.Dates)) // nb of actor for role team
			// for each actors of current roleTeam (having role)
			var roleTeamRoleName string
			for _, actPos := range roleTeamActorsPos {
				nbActors[actPos] = map[string][]float64{}
				for roleName, actorRoleData := range workData[actPos] {
					roleTeamRoleName = roleName
					actNbAct := make([]float64, len(aggrStat.Dates)) // nb of actor for actor
					for i := 0; i < len(actNbAct); i++ {
						if actorRoleData[i] > threshold {
							actNbAct[i] = 1
							roleNbAct[i] += 1
							mainNbAct[i] += 1
						}
					}
					// Set actor's > Role > NbActors Data
					nbActors[actPos][roleName] = actNbAct
				}
			}
			nbActors[roleTeamIndex[roleTeamName]] = map[string][]float64{roleTeamRoleName: roleNbAct}
			nbActors[mainTeamPos][roleTeamRoleName] = roleNbAct
			roleTeamRole[roleTeamName] = roleTeamRoleName
		}
		mainTeamNbAct[mainTeamName] = mainNbAct
	}
	aggrStat.Values[StatSerieNbActor] = nbActors

	//Calc Mean Values
	zeros := make([]float64, len(aggrStat.Dates))
	roleMeanData := make([]map[string][]float64, len(aggrStat.Teams))
	//mainMeanData := make([]map[string][]float64, len(aggrStat.Teams))

	// for each main teams
	for mainTeamName, mainTeamPos := range mainTeamIndex {
		mainTeamData := make([]float64, len(aggrStat.Dates))
		mainTeamNbActorData := make([]float64, len(aggrStat.Dates))
		roleMeanData[mainTeamPos] = map[string][]float64{}

		// for each role team within main team
		for roleTeamName, roleTeamActorsPos := range actorsIndex[mainTeamName] {
			// calc roleTeam Mean work
			roleTeamData := make([]float64, len(aggrStat.Dates))
			roleTeamNbActorData := nbActors[roleTeamIndex[roleTeamName]][roleTeamRole[roleTeamName]]
			var roleTeamRoleName string
			// for each dates
			for i := 0; i < len(roleTeamData); i++ {
				// for each actor in role team
				for _, actPos := range roleTeamActorsPos {
					// get pertaining role and sum work
					for roleName, actorRoleData := range workData[actPos] {
						roleTeamRoleName = roleName
						roleTeamData[i] += actorRoleData[i]
					}
				}
				mainTeamData[i] += roleTeamData[i]
				mainTeamNbActorData[i] += roleTeamNbActorData[i]

				// then calc mean
				if roleTeamNbActorData[i] > 0 {
					roleTeamData[i] /= roleTeamNbActorData[i]
				}
			}

			// set to zero work data for mainTeam & roleTeam
			workData[roleTeamIndex[roleTeamName]][roleTeamRoleName] = zeros
			workData[mainTeamPos][roleTeamRoleName] = zeros
			// set role mean data ...
			// ... for mainTeam
			roleMeanData[mainTeamPos][roleTeamRoleName] = zeros
			// ... for roleTeam
			roleMeanData[roleTeamIndex[roleTeamName]] = map[string][]float64{roleTeamRoleName: roleTeamData}
			// ... and for each pertaining actors
			for _, actPos := range roleTeamActorsPos {
				nbAct := nbActors[actPos][roleTeamRoleName]
				actData := make([]float64, len(aggrStat.Dates))
				for i := 0; i < len(nbAct); i++ {
					if nbAct[i] > 0 {
						actData[i] = roleTeamData[i]
					}
				}
				roleMeanData[actPos] = map[string][]float64{roleTeamRoleName: actData}
			}

		}

	}
	aggrStat.Values[StatSerieRoleMeanData] = roleMeanData

	return aggrStat
}
