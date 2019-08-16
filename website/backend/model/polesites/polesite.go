package polesites

import (
	"fmt"
	"github.com/lpuig/ewin/doe/website/backend/model/bpu"
	"github.com/lpuig/ewin/doe/website/backend/model/clients"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	"github.com/lpuig/ewin/doe/website/backend/model/items"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/model/polesite/poleconst"
	"strings"
)

type PoleSite struct {
	Id         int
	Client     string
	Ref        string
	Manager    string
	OrderDate  string
	UpdateDate string
	Status     string
	Comment    string

	Poles []*Pole
}

func (ps *PoleSite) GetInfo() *fm.PolesiteInfo {
	psi := fm.NewBEPolesiteInfo()

	psi.Id = ps.Id
	psi.Client = ps.Client
	psi.Ref = ps.Ref
	psi.Manager = ps.Manager
	psi.OrderDate = ps.OrderDate
	psi.UpdateDate = ps.UpdateDate
	psi.Status = ps.Status
	psi.Comment = ps.Comment

	psi.NbPole, psi.NbPoleBlocked, psi.NbPoleDone = ps.GetPolesNumbers()

	var searchBuilder strings.Builder
	fmt.Fprintf(&searchBuilder, "%s:%s,", "Client", strings.ToUpper(ps.Client))
	fmt.Fprintf(&searchBuilder, "%s:%s,", "Ref", strings.ToUpper(ps.Ref))
	fmt.Fprintf(&searchBuilder, "%s:%s,", "Manager", strings.ToUpper(ps.Manager))
	fmt.Fprintf(&searchBuilder, "%s:%s,", "OrderDate", strings.ToUpper(ps.OrderDate))
	fmt.Fprintf(&searchBuilder, "%s:%s,", "Comment", strings.ToUpper(ps.Comment))
	for _, pole := range ps.Poles {
		fmt.Fprintf(&searchBuilder, "%s,", pole.SearchString())
	}
	psi.Search = searchBuilder.String()

	return psi
}

// GetPolesNumbers returns total, blocked and done number of Pullings
func (ps *PoleSite) GetPolesNumbers() (total, blocked, done int) {
	for _, p := range ps.Poles {
		switch p.State {
		//case poleconst.StateNotSubmitted:
		case poleconst.StateNoGo:
			total++
			blocked++
		case poleconst.StateToDo:
			total++
		case poleconst.StateHoleDone:
			total++
		case poleconst.StateIncident:
			total++
			blocked++
		case poleconst.StateDone:
			total++
			done++
			//case poleconst.StateCancelled:
		}
	}
	return
}

type IsPolesiteVisible func(s *PoleSite) bool

// Itemize returns slice of item pertaining to polesite poles list
func (ps *PoleSite) Itemize(currentBpu *bpu.Bpu, actorById clients.ActorById) ([]*items.Item, error) {
	res := []*items.Item{}

	for _, pole := range ps.Poles {
		items, err := pole.Itemize(currentBpu, actorById)
		if err != nil {
			return nil, err
		}
		res = append(res, items...)
	}
	return res, nil
}

// AddStat adds Stats into values for given Polesite
func (ps *PoleSite) AddStat(stats items.Stats, dateFor date.DateAggreg,
	isActorVisible clients.IsTeamVisible, actorById clients.ActorById,
	currentBpu *bpu.Bpu, teamName clients.TeamNameByMember, showprice bool) error {

	addValue := func(client, site, team, date, article, serie string, val float64) {
		teamInfo := "Eq. " + teamName(team)
		stats.AddStatValue(site, client+" : "+teamInfo, dateFor(date), article, serie, val)
		//values[items.StatKey{
		//	Team:    client + " : " + teamInfo,
		//	Date:    dateFor(date),
		//	Site:    site,
		//	Article: article,
		//	Serie:   serie,
		//}] += val
		stats.AddStatValue(site, client, dateFor(date), article, serie, val)
		//values[items.StatKey{
		//	Team:    client,
		//	Date:    dateFor(date),
		//	Site:    site,
		//	Article: article,
		//	Serie:   serie,
		//}] += val
	}

	calcItems, err := ps.Itemize(currentBpu, actorById)
	if err != nil {
		return fmt.Errorf("error on polesite stat itemize for '%s':%s", ps.Ref, err.Error())
	}
	for _, item := range calcItems {
		if !item.Done {
			continue
		}
		work := item.Work()
		addValue(ps.Client, ps.Ref, item.Team, item.Date, item.Article.Name, items.StatSerieWork, work)
		if showprice {
			price := item.Price()
			addValue(ps.Client, ps.Ref, item.Team, item.Date, item.Article.Name, items.StatSeriePrice, price)
		}
	}
	return nil
}
