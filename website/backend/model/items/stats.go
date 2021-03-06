package items

import (
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	"github.com/lpuig/ewin/doe/website/backend/model/groups"
	rs "github.com/lpuig/ewin/doe/website/frontend/model/ripsite"
	"strings"
)

type Stats map[StatKey]float64

func NewStats() Stats {
	return make(Stats)
}

func (s Stats) AddStatValue(stackedSerie, graphName, date, article, serie string, value float64) {
	s[StatKey{
		Graph:        graphName,
		Date:         date,
		StackedSerie: stackedSerie,
		Article:      article,
		Serie:        serie,
	}] += value
}

// Aggregate returns struct with ws.Values : map{D1}[#D2]{D3}[#date]float64
//
// Standard usage is ws.Values : map{series}[#team]{sites}[#date]float64
func (s Stats) Aggregate(sc StatContext) *rs.RipsiteStats {

	//create client-team, sites, Series & dates Lists
	serieSet := make(ElemSet)        // d1
	graphSet := make(ElemSet)        // d2
	stackedSerieSet := make(ElemSet) // d3
	dateSet := make(ElemSet)

	end := date.Today()
	start := end.String()

	agrValues := make(Stats) // values

	for key, val := range s {
		k1 := sc.Data1(key)
		k2 := sc.Data2(key)
		k3 := sc.Data3(key)
		serieSet[k1] = 1        // d1
		graphSet[k2] = 1        // d2
		stackedSerieSet[k3] = 1 // d3
		if key.Date < start {
			start = key.Date
		}
		agrValues[StatKey{
			Serie:        k1, // d1
			Graph:        k2, // d2
			StackedSerie: k3, // d3
			Article:      "", // d4 always aggregated
			Date:         key.Date,
		}] += val
	}

	curDate := date.DateFrom(sc.StartDate)
	for len(dateSet) < sc.MaxVal {
		curStringDate := sc.DateFor(curDate.String())
		dateSet[curStringDate] = 1
		curDate = curDate.AddDays(sc.DayIncr)
	}

	series := serieSet.SortedKeys(sc.Filter1)               // d1
	graphs := graphSet.SortedKeys(sc.Filter2)               // d2
	stackedSeries := stackedSerieSet.SortedKeys(sc.Filter3) // d3
	dates := dateSet.SortedKeys(KeepAll)                    // d4

	// stats.Values : map{D1}[#D2]{D3}[#D4]float64
	stats := rs.NewBERipsiteStats()
	stats.Dates = dates // d4
	sitesmap := map[string]bool{}

	for _, graphName := range graphs {
		graphActivity := 0.0
		values := make(map[string]map[string][]float64)
		for _, serie := range series {
			values[serie] = make(map[string][]float64)
			for _, stackedSerie := range stackedSeries {
				siteData := make([]float64, len(dates))
				siteActivity := 0.0
				for dateNum, d := range dates {
					val := agrValues[StatKey{
						Graph:        graphName,
						Date:         d,
						StackedSerie: stackedSerie,
						Article:      "",
						Serie:        serie,
					}]
					graphActivity += val
					siteActivity += val
					siteData[dateNum] = val
				}
				if siteActivity > 0 {
					sitesmap[stackedSerie] = true
					values[serie][stackedSerie] = siteData
				}
			}
		}
		if graphActivity == 0 {
			// current graph as no activity on the time laps, skip it // d2
			continue
		}
		stats.Teams = append(stats.Teams, graphName) // d2
		for _, serie := range series {
			stats.Values[serie] = append(stats.Values[serie], values[serie]) // d1
		}
	}
	stats.Sites = sitesmap // d3

	return stats
}

// CalcTeamMean adds in given RipsiteStats mean values for each teams (D2 dimension) : map{D1}[#D2]{D3}[#date]float64
//
// Standard usage is ws.Values : map{series}[#team]{sites}[#date]float64
func CalcTeamMean(aggrStat *rs.RipsiteStats, threshold float64) *rs.RipsiteStats {

	mapData := func(data, criteria []float64, threshold float64, ope func(data, criteria float64) float64) []float64 {
		res := make([]float64, len(data))
		for i, val := range data {
			if criteria[i] > threshold {
				res[i] = ope(val, criteria[i])
			}
		}
		return res
	}
	getData := func(data, criteria float64) float64 { return data }
	//divData := func(t, a float64) float64 { return t / a }

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

	const (
		StatSerieRoleMeanData string = "RoleMeanWork"
		StatSerieMainMeanData string = "GlobalMeanWork"
		StatSerieNbActor      string = "NbActorsWork"

		StatSiteGlobal string = " Global"
	)

	workData := aggrStat.Values[StatSerieWork]
	nbActors := make([]map[string][]float64, len(aggrStat.Teams))
	mainTeamNbAct := map[string][]float64{}
	roleTeamRole := map[string]string{}

	// Calc Nb Actors
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
		//mainTeamData := make([]float64, len(aggrStat.Dates))
		mainTeamNbActorData := make([]float64, len(aggrStat.Dates))
		roleMeanData[mainTeamPos] = map[string][]float64{}
		//mainMeanData[mainTeamPos] = map[string][]float64{}

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
						if actorRoleData[i] > threshold {
							roleTeamData[i] += actorRoleData[i]
						}
					}
				}
				//mainTeamData[i] += roleTeamData[i]
				mainTeamNbActorData[i] += roleTeamNbActorData[i]

				// then calc mean
				if roleTeamNbActorData[i] > 0 {
					roleTeamData[i] /= roleTeamNbActorData[i]
				}
			}

			// set to zero work data for mainTeam & roleTeam
			workData[roleTeamIndex[roleTeamName]][roleTeamRoleName] = zeros
			//workData[mainTeamPos][roleTeamRoleName] = zeros

			// set role mean data ...
			// ... for mainTeam
			roleMeanData[mainTeamPos][roleTeamRoleName] = zeros

			// ... for roleTeam
			roleMeanData[roleTeamIndex[roleTeamName]] = map[string][]float64{roleTeamRoleName: roleTeamData}
			//mainMeanData[roleTeamIndex[roleTeamName]] = map[string][]float64{}
			// ... and for each pertaining actors
			for _, actPos := range roleTeamActorsPos {
				actData := mapData(roleTeamData, nbActors[actPos][roleTeamRoleName], 0, getData)
				roleMeanData[actPos] = map[string][]float64{roleTeamRoleName: actData}
				//mainMeanData[actPos] = map[string][]float64{}
			}
		}
	}

	aggrStat.Values[StatSerieRoleMeanData] = roleMeanData
	//aggrStat.Values[StatSerieMainMeanData] = mainMeanData

	//aggrStat.Sites[StatSiteGlobal] = true

	return aggrStat
}

// CalcProgress processes given RipsiteStats, computing cumulative figures and adding forecast result
//
// groupSize contains , per group name, active group's actor number (slice len is month days number)
func CalcProgress(aggrStat *rs.RipsiteStats, groupByName groups.GroupByName, groupSize map[string][]int) *rs.RipsiteStats {
	var serieNameTarget string

	isActiveDate := make([]bool, len(aggrStat.Dates))
	//nbActiveDay := 0
	// dayoff calendar not to be managed
	for dateId, day := range aggrStat.Dates {
		if date.DateFrom(day).IsSaturdaySunday() {
			continue
		}
		isActiveDate[dateId] = true
		//nbActiveDay++
	}

	newSerieTeamDataValues := make(map[string][]map[string][]float64)
	for serieName, teamDataValues := range aggrStat.Values {
		//_ = serieName: Price | Work
		switch serieName {
		case StatSerieWork:
			serieNameTarget = StatSerieWorkTarget
		case StatSeriePrice:
			serieNameTarget = StatSeriePriceTarget
		}
		newTeamDataValues := make([]map[string][]float64, len(teamDataValues))
		for teamIndex, dataValues := range teamDataValues {
			//_ = teamIndex <=> group & actor per graph

			// Replace values by Cumulative values
			for _, values := range dataValues {
				//_ = siteName
				for dateId := 1; dateId < len(values); dateId++ {
					values[dateId] += values[dateId-1]
				}
			}

			// check if teamIndex is group or individual actor
			teamName := aggrStat.Teams[teamIndex]
			groupActor := strings.Split(teamName, " : ")
			groupName := groupActor[0] // retrieve group Name
			//actorLabel := ""
			//if len(groupActor) > 1 {
			//	actorLabel = groupActor[1]
			//}
			incrVal := 0.0
			switch serieName {
			case StatSerieWork:
				incrVal = groupByName(groupName).ActorDailyWork
			case StatSeriePrice:
				incrVal = groupByName(groupName).ActorDailyIncome
			}

			newDataValues := make(map[string][]float64)
			// if team[teamIndex] == main team
			// 		Calc incremental target
			nbDays := len(aggrStat.Dates)
			//name := actorLabel
			//if actorLabel == "" {
			//	// current team is an group
			//	name = groupName
			//}
			activity, found := groupSize[teamName]
			targetVal := 0.0
			targetValues := make([]float64, nbDays)
			if found {
				for dateId := 0; dateId < nbDays; dateId++ {
					if isActiveDate[dateId] {
						targetVal += incrVal * float64(activity[dateId])
					}
					targetValues[dateId] = targetVal
				}
			}
			newDataValues[StatSiteProgressTarget] = targetValues
			newTeamDataValues[teamIndex] = newDataValues
		}
		newSerieTeamDataValues[serieNameTarget] = newTeamDataValues
	}
	// Add newSerieTeamDataValues to existing values
	for serieName, serieValues := range newSerieTeamDataValues {
		aggrStat.Values[serieName] = serieValues
	}
	aggrStat.Sites[StatSiteProgressTarget] = true
	return aggrStat
}
