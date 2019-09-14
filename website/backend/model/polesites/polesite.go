package polesites

import (
	"fmt"
	"github.com/lpuig/ewin/doe/website/backend/model/bpu"
	"github.com/lpuig/ewin/doe/website/backend/model/clients"
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
		case poleconst.StateDictToDo:
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
func (ps *PoleSite) Itemize(currentBpu *bpu.Bpu) ([]*items.Item, error) {
	res := []*items.Item{}

	for _, pole := range ps.Poles {
		items, err := pole.Itemize(ps.Client, ps.Ref, currentBpu)
		if err != nil {
			return nil, err
		}
		res = append(res, items...)
	}
	return res, nil
}

// AddStat adds Stats into values for given Polesite
func (ps *PoleSite) AddStat(stats items.Stats, sc items.StatContext,
	actorById clients.ActorById, currentBpu *bpu.Bpu, showprice bool) error {

	addValue := func(date, serie string, actors []string, value float64) {
		stats.AddStatValue(ps.Ref, ps.Client, date, "", serie, value)
		if sc.ShowTeam && len(actors) > 0 {
			value /= float64(len(actors))
			for _, actName := range actors {
				stats.AddStatValue(ps.Ref, ps.Client+" : "+actName, date, "", serie, value)
			}
		}
	}

	calcItems, err := ps.Itemize(currentBpu)
	if err != nil {
		return fmt.Errorf("error on polesite stat itemize for '%s':%s", ps.Ref, err.Error())
	}
	for _, item := range calcItems {
		if !item.Done {
			continue
		}
		actorsName := make([]string, len(item.Actors))
		for i, actId := range item.Actors {
			actorsName[i] = actorById(actId)
		}
		addValue(sc.DateFor(item.Date), items.StatSerieWork, actorsName, item.Work())
		if showprice {
			addValue(sc.DateFor(item.Date), items.StatSeriePrice, actorsName, item.Price())
		}
	}
	return nil
}
