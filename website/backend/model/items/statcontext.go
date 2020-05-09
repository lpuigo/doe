package items

import (
	"fmt"
	"github.com/lpuig/ewin/doe/website/backend/model/bpu"
	"github.com/lpuig/ewin/doe/website/backend/model/clients"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	"github.com/lpuig/ewin/doe/website/frontend/model/ripsite"
)

type StatContext struct {
	DayIncr       int
	MaxVal        int
	StartDate     string
	DateFor       date.DateAggreg
	IsTeamVisible clients.IsTeamVisible
	ClientByName  clients.ClientByName
	ActorById     clients.ActorById
	ShowTeam      bool

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
	Itemize(*bpu.Bpu) ([]*Item, error)
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
		startDate = date.GetMonth(date.GetDateAfter(startDate, (1-maxVal)*30))
	default:
		return nil, fmt.Errorf("unsupported stat period '%s'", freq)
	}

	return &StatContext{
		DayIncr:   dayIncr,
		MaxVal:    maxVal,
		StartDate: startDate,
		DateFor:   dateFor,
	}, nil
}

func (sc *StatContext) SetSerieTeamSiteConf() {
	sc.Data1 = func(s StatKey) string { return s.Serie } // Bars Family
	sc.Data2 = func(s StatKey) string { return s.Team }  // Graphs
	sc.Data3 = func(s StatKey) string { return s.Site }  // side block
	sc.Filter1 = KeepAll
	sc.Filter2 = KeepAll
	sc.Filter3 = KeepAll
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
	addValue := func(date, serie string, actors []string, value float64) {
		stats.AddStatValue(site.GetRef(), site.GetClient(), date, "", serie, value)
		if sc.ShowTeam && len(actors) > 0 {
			value /= float64(len(actors))
			for _, actName := range actors {
				stats.AddStatValue(site.GetRef(), site.GetClient()+" : "+actName, date, "", serie, value)
			}
		}
	}

	calcItems, err := site.Itemize(currentBpu)
	if err != nil {
		return fmt.Errorf("error on %s stat itemize for '%s':%s", site.GetType(), site.GetRef(), err.Error())
	}
	for _, item := range calcItems {
		if !item.Done {
			continue
		}
		dateItem := sc.DateFor(item.Date)
		if dateItem < sc.StartDate {
			continue
		}
		actorsName := make([]string, len(item.Actors))
		for i, actId := range item.Actors {
			actorsName[i] = sc.ActorById(actId)
		}
		addValue(dateItem, StatSerieWork, actorsName, item.Work())
		if showprice {
			addValue(dateItem, StatSeriePrice, actorsName, item.Price())
		}
	}
	return nil
}
