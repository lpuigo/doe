package polesites

import (
	"fmt"
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
		case poleconst.StateNotSubmitted:
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
		case poleconst.StateCancelled:
		}
	}
	return
}

type IsPolesiteVisible func(s *PoleSite) bool
