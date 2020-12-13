package manager

import (
	"encoding/json"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	"github.com/lpuig/ewin/doe/website/backend/model/items"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"

	"io"
)

// GetPolesitesInfo returns array of PolesiteInfos (JSON in writer) visibles by current user
func (m Manager) GetPolesitesInfo(writer io.Writer) error {
	psis := []*fm.PolesiteInfo{}
	for _, psr := range m.Polesites.GetAll(m.visibleItemizableSiteByClientFilter()) {
		psis = append(psis, psr.PoleSite.GetInfo())
	}

	return json.NewEncoder(writer).Encode(psis)
}

func (m Manager) GetPolesitesStats(writer io.Writer, freq string) error {
	statContext, err := m.NewStatContext(freq)
	if err != nil {
		return err
	}

	polesiteStats, err := statContext.CalcStats(m.Polesites, m.visibleItemizableSiteByClientFilter(), m.CurrentUser.Permissions["Invoice"])
	if err != nil {
		return err
	}
	return json.NewEncoder(writer).Encode(polesiteStats)
}

func (m Manager) GetPolesitesProgress(writer io.Writer, month string) error {
	statContext, err := m.NewStatContext("day")
	if err != nil {
		return err
	}
	// start on month first
	statContext.StartDate = month
	// Extraxt all month's days
	statContext.MaxVal = date.NbDaysBetween(month, date.GetMonth(date.GetDateAfter(month, 32)))
	groupNameByClient := make(map[string]string)
	for _, group := range m.Groups.GetGroups() {
		for _, client := range group.Clients {
			groupNameByClient[client] = group.Name
		}
	}
	statContext.GraphName = func(site items.ItemizableSite) string {
		groupName, found := groupNameByClient[site.GetClient()]
		if !found {
			return m.Groups.GetById(0).Name
		}
		return groupName
	}
	statContext.Data3 = func(s items.StatKey) string { return items.StatSiteProgress }
	//statContext.Data3 = func(s items.StatKey) string { return s.Team}

	polesiteStats, err := statContext.CalcStats(m.Polesites, m.visibleItemizableSiteByClientFilter(), m.CurrentUser.Permissions["Invoice"])
	if err != nil {
		return err
	}
	return json.NewEncoder(writer).Encode(items.CalcProgress(polesiteStats, m.GenGroupByName(), m.GroupSizeOnMonth(polesiteStats.Dates)))
}
