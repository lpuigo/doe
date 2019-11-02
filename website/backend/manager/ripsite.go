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
	if len(m.CurrentUser.Clients) == 0 {
		return func(*rs.Site) bool { return true }
	}
	isVisible := make(map[string]bool)
	for _, client := range m.CurrentUser.Clients {
		isVisible[client] = true
	}
	return func(s *rs.Site) bool {
		return isVisible[s.Client]
	}
}

// GetRipsitesInfo returns array of RipsiteInfos (JSON in writer) visibles by current user
func (m Manager) GetRipsitesInfo(writer io.Writer) error {
	clientByName := m.genGetClient()
	rsis := []*fm.RipsiteInfo{}
	for _, rsr := range m.Ripsites.GetAll(m.visibleRipsiteFilter()) {
		rsis = append(rsis, rsr.Site.GetInfo(clientByName))
	}

	return json.NewEncoder(writer).Encode(rsis)
}

func (m Manager) GetRipsiteXLSAttachement(writer io.Writer, rs *rs.Site) error {
	return m.TemplateEngine.GetRipsiteXLSAttachement(writer, rs, m.genGetClient(), m.genActorById())
}

func (m Manager) RipsitesArchiveName() string {
	return m.Ripsites.ArchiveName()
}

func (m Manager) CreateRipsitesArchive(writer io.Writer) error {
	return m.Ripsites.CreateArchive(writer)
}

func (m Manager) GetRipsitesStats(writer io.Writer, freq, groupBy string) error {
	maxVal := 12

	var dateFor date.DateAggreg
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

	isActorVisible, err := m.genIsActorVisible()
	if err != nil {
		return err
	}

	statContext := items.StatContext{
		MaxVal:        maxVal,
		DateFor:       dateFor,
		IsTeamVisible: isActorVisible,
		ShowTeam:      !m.CurrentUser.Permissions["Review"],
	}

	switch groupBy {
	case "activity", "site":
		ripsiteStats, err := m.Ripsites.GetProdStats(statContext, m.visibleRipsiteFilter(), m.genGetClient(), m.genActorById(), m.CurrentUser.Permissions["Invoice"], groupBy)
		if err != nil {
			return err
		}
		return json.NewEncoder(writer).Encode(ripsiteStats)

	case "mean":
		ripsiteStats, err := m.Ripsites.GetMeanProdStats(statContext, m.visibleRipsiteFilter(), m.genGetClient(), m.genActorInfoById())
		if err != nil {
			return err
		}
		meanStats := items.CalcTeamMean(ripsiteStats, 1)
		return json.NewEncoder(writer).Encode(meanStats)

	default:
		return fmt.Errorf("unsupported group type '%s'", groupBy)
	}
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
	return m.TemplateEngine.GetItemsXLSAttachement(writer, itms, m.genActorById())
}
