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

	d1 := func(s StatKey) string { return s.Serie } // Bars Family
	d2 := func(s StatKey) string { return s.Team }  // Graphs
	d3 := func(s StatKey) string { return s.Site }  // side block
	f1 := KeepAll
	//f2 := func(e string) bool { return !(!sc.ShowTeam && strings.Contains(e, " : ")) }
	f2 := KeepAll
	f3 := KeepAll
	return calcValues.Aggregate(sc, d1, d2, d3, f1, f2, f3), nil
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
