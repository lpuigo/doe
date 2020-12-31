package manager

import (
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	"github.com/lpuig/ewin/doe/website/backend/model/items"
	rs "github.com/lpuig/ewin/doe/website/backend/model/ripsites"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"

	"encoding/json"
	"fmt"
	"io"
)

// visibleRipsiteFilter returns a filtering function on CurrentUser.Clients visibility
func (m *Manager) visibleRipsiteFilter() rs.IsSiteVisible {
	clts, err := m.GetCurrentUserClients()
	if err != nil {
		return func(*rs.Site) bool { return false }
	}
	isVisible := make(map[string]bool)
	for _, client := range clts {
		isVisible[client.Name] = true
	}
	return func(s *rs.Site) bool {
		return isVisible[s.Client]
	}
}

// GetRipsitesInfo returns array of RipsiteInfos (JSON in writer) visibles by current user
func (m Manager) GetRipsitesInfo(writer io.Writer) error {
	clientByName := m.genGetClient()
	rsis := []*fm.RipsiteInfo{}
	for _, rsr := range m.Ripsites.GetAll(m.visibleItemizableSiteByClientFilter()) {
		rsis = append(rsis, rsr.Site.GetInfo(clientByName))
	}

	return json.NewEncoder(writer).Encode(rsis)
}

func (m Manager) GetRipsitesStats(writer io.Writer, freq, groupBy string) error {
	statContext, err := m.NewStatContext(freq)
	if err != nil {
		return err
	}

	switch groupBy {
	case "activity":
		statContext.StackedSerieName = func(item *items.Item) string {
			return item.Activity
		}
	case "actor":
		statContext.SetGraphNameByActor()
	case "site":

	default:
		return fmt.Errorf("unsupported group type '%s'", groupBy)
	}

	ripsiteStats, err := statContext.CalcStats(m.Ripsites, m.visibleItemizableSiteByClientFilter(), m.CurrentUser.Permissions["Invoice"])
	if err != nil {
		return err
	}
	return json.NewEncoder(writer).Encode(ripsiteStats)
}

func (m Manager) GetRipsitesActorsActivity(writer io.Writer, freq string) error {

	var dateFor date.DateAggreg
	var firstDate string
	switch freq {
	case "week":
		dateFor = func(d string) string {
			return date.GetMonday(d)
		}
	case "month":
		dateFor = func(d string) string {
			return date.GetMonth(d)
		}
	default:
		return fmt.Errorf("unsupported stat period '%s'", freq)
	}

	// set firstDate according to freq choice, in order to have at least a full month of data
	// month : last and current month
	// week : 5 last weeks and current
	firstDate = dateFor(date.Today().AddDays(-32).String())

	itms, err := m.Ripsites.GetAllItems(firstDate, dateFor, m.visibleRipsiteFilter(), m.genGetClient())
	if err != nil {
		return err
	}
	return m.TemplateEngine.GetItemsXLSAttachement(writer, itms, m.genActorNameById())
}
