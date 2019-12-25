package manager

import (
	"encoding/json"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	"github.com/lpuig/ewin/doe/website/backend/model/items"
	ps "github.com/lpuig/ewin/doe/website/backend/model/polesites"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"

	"io"
)

// visiblePolesiteFilter returns a filtering function on CurrentUser.Clients visibility
func (m *Manager) visiblePolesiteFilter() ps.IsPolesiteVisible {
	if len(m.CurrentUser.Clients) == 0 {
		return func(ps *ps.PoleSite) bool { return true }
	}
	isVisible := make(map[string]bool)
	for _, client := range m.CurrentUser.Clients {
		isVisible[client] = true
	}
	return func(ps *ps.PoleSite) bool {
		return isVisible[ps.Client]
	}
}

// GetPolesitesInfo returns array of PolesiteInfos (JSON in writer) visibles by current user
func (m Manager) GetPolesitesInfo(writer io.Writer) error {
	psis := []*fm.PolesiteInfo{}
	for _, psr := range m.Polesites.GetAll(m.visiblePolesiteFilter()) {
		psis = append(psis, psr.PoleSite.GetInfo())
	}

	return json.NewEncoder(writer).Encode(psis)
}

func (m Manager) GetPolesiteXLSAttachement(writer io.Writer, ps *ps.PoleSite) error {
	return m.TemplateEngine.GetPolesiteXLSAttachement(writer, ps, m.genGetClient(), m.genActorById())
}

// GetPolesitesWeekStats returns Polesites Stats per Week (JSON in writer) visibles by current user
func (m Manager) GetPolesitesWeekStats(writer io.Writer) error {
	df := func(d string) string {
		return date.GetMonday(d)
	}
	return m.getPolesitesStats(writer, 12, df)
}

// GetPolesitesMonthStats returns Polesites Stats per Month (JSON in writer) visibles by current user
func (m Manager) GetPolesitesMonthStats(writer io.Writer) error {
	df := func(d string) string {
		return date.GetMonth(d)
	}
	return m.getPolesitesStats(writer, 12, df)
}

func (m Manager) getPolesitesStats(writer io.Writer, maxVal int, dateFor date.DateAggreg) error {
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

	polesiteStats, err := m.Polesites.GetStats(statContext, m.visiblePolesiteFilter(), m.genGetClient(), m.genActorById(), m.CurrentUser.Permissions["Invoice"])
	if err != nil {
		return err
	}
	return json.NewEncoder(writer).Encode(polesiteStats)
}
