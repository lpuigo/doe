package items

import (
	"fmt"
	"github.com/lpuig/ewin/doe/website/backend/model/bpu"
	"github.com/lpuig/ewin/doe/website/backend/model/clients"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	"github.com/lpuig/ewin/doe/website/frontend/model/ripsite"
)

type NamePct struct {
	Name string
	Pct  float64
}

type StatContext struct {
	DayIncr          int
	MaxVal           int
	StartDate        string
	DateFor          date.DateAggreg
	IsTeamVisible    clients.IsTeamVisible
	ClientByName     clients.ClientByName
	ActorNameById    clients.ActorNameById
	GraphName        func(item *Item) []NamePct
	StackedSerieName func(item *Item) string

	ShowTeam bool

	Data1   func(s StatKey) string
	Data2   func(s StatKey) string
	Data3   func(s StatKey) string
	Filter1 func(string) bool
	Filter2 func(string) bool
	Filter3 func(string) bool
}

type IsItemizableSiteVisible func(site ItemizableSite) bool

type ItemizableContainer interface {
	GetItemizableSites(IsItemizableSiteVisible) []ItemizableSite
	GetItemizableSiteById(int) ItemizableSite
}

type ItemizableSite interface {
	GetRef() string
	GetClient() string
	GetType() string
	GetUpdateDate() string
	Itemize(currentBpu *bpu.Bpu, doneOnly bool) ([]*Item, error)
}

func NewStatContext(freq string) (*StatContext, error) {
	maxVal := 12
	dayIncr := 7
	startDate := date.Today().String()

	var dateFor date.DateAggreg
	switch freq {
	case "day":
		dayIncr = 1
		maxVal = 15
		dateFor = func(d string) string {
			return d
		}
		startDate = date.Today().AddDays(1 - maxVal).String()
	case "week":
		dateFor = func(d string) string {
			return date.GetMonday(d)
		}
		startDate = date.GetDateAfter(date.GetMonday(startDate), (1-maxVal)*7)
	case "month":
		dateFor = func(d string) string {
			return date.GetMonth(d)
		}
		startDate = dateFor(date.GetDateAfter(dateFor(startDate), (1-maxVal)*30))
	default:
		return nil, fmt.Errorf("unsupported stat period '%s'", freq)
	}

	sc := &StatContext{
		DayIncr:   dayIncr,
		MaxVal:    maxVal,
		StartDate: startDate,
		DateFor:   dateFor,
		StackedSerieName: func(item *Item) string {
			return item.Site
		},
	}
	sc.GraphName = func(item *Item) []NamePct {
		mainName := item.Client
		res := []NamePct{}
		globPct := 1.0
		if sc.ShowTeam && len(item.Actors) > 0 {
			pct := 1.0 / float64(len(item.Actors))
			globPct = 0.0
			for _, actId := range item.Actors {
				actorName := sc.ActorNameById(actId)
				if actorName == "" { // Skip unknown or not visible Actors
					continue
				}
				res = append(res, NamePct{
					Name: mainName + " : " + actorName,
					Pct:  pct,
				})
				globPct += pct
			}
		}
		res = append(res, NamePct{
			Name: mainName,
			Pct:  globPct,
		})
		return res
	}
	return sc, nil
}

func (sc *StatContext) SetSerieTeamSiteConf() {
	sc.Data1 = func(s StatKey) string { return s.Serie }        // Bars Family
	sc.Data2 = func(s StatKey) string { return s.Graph }        // Graphs
	sc.Data3 = func(s StatKey) string { return s.StackedSerie } // side block
	sc.Filter1 = KeepAll
	sc.Filter2 = KeepAll
	sc.Filter3 = KeepAll
}

func (sc *StatContext) SetGraphNameByActor() {
	sc.GraphName = func(item *Item) []NamePct {
		res := []NamePct{}
		if len(item.Actors) > 0 {
			pct := 1.0 / float64(len(item.Actors))
			for _, actId := range item.Actors {
				actName := sc.ActorNameById(actId)
				if actName == "" { // Skip unknown or not visible Actors
					continue
				}
				res = append(res, NamePct{
					Name: actName,
					Pct:  pct,
				})
			}
		} else {
			res = []NamePct{{
				Name: "Pas d'acteur",
				Pct:  1.0,
			}}
		}
		return res
	}
}

func (sc StatContext) CalcStats(sites ItemizableContainer, isSiteVisible IsItemizableSiteVisible, showprice bool) (*ripsite.RipsiteStats, error) {
	calcValues := NewStats()
	for _, site := range sites.GetItemizableSites(isSiteVisible) {
		if site.GetUpdateDate() < sc.StartDate {
			continue
		}
		client := sc.ClientByName(site.GetClient())
		if client == nil {
			continue
		}
		err := sc.addStat(calcValues, site, client.Bpu, showprice)
		if err != nil {
			return nil, err
		}
	}

	return calcValues.Aggregate(sc), nil
}

func (sc StatContext) addStat(stats Stats, site ItemizableSite, currentBpu *bpu.Bpu, showprice bool) error {
	addValue := func(item *Item, statDate, serie string, value float64) {
		stackedSerie := sc.StackedSerieName(item)
		for _, gNP := range sc.GraphName(item) {
			stats.AddStatValue(stackedSerie, gNP.Name, statDate, "", serie, value*gNP.Pct)
		}
	}

	calcItems, err := site.Itemize(currentBpu, true)
	if err != nil {
		return fmt.Errorf("error on %s stat itemize for '%s':%s", site.GetType(), site.GetRef(), err.Error())
	}
	for _, item := range calcItems {
		statDate := sc.DateFor(item.Date)
		if statDate < sc.StartDate {
			continue
		}
		addValue(item, statDate, StatSerieWork, item.Work())
		if showprice {
			addValue(item, statDate, StatSeriePrice, item.Price())
		}
	}
	return nil
}
