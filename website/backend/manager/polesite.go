package manager

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/lpuig/ewin/doe/website/backend/model/date"
	"github.com/lpuig/ewin/doe/website/backend/model/items"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
)

// GetPolesitesInfo returns array of PolesiteInfos (JSON in writer) visibles by current user
func (m Manager) GetPolesitesInfo(writer io.Writer) error {
	psis := []*fm.PolesiteInfo{}
	for _, psr := range m.Polesites.GetAll(m.visibleItemizableSiteByClientFilter()) {
		psis = append(psis, psr.PoleSite.GetInfo())
	}

	return json.NewEncoder(writer).Encode(psis)
}

func (m Manager) GetPolesitesStats(writer io.Writer, freq, groupBy string) error {
	statContext, err := m.NewStatContext(freq)
	if err != nil {
		return err
	}

	switch groupBy {
	case "actor":
		statContext.SetGraphNameByActor()
	case "client":

	default:
		return fmt.Errorf("unsupported group type '%s'", groupBy)
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
	// Extract all month's days
	statContext.MaxVal = date.NbDaysBetween(month, date.GetMonth(date.GetDateAfter(month, 32)))
	groupNameByClient := make(map[string]string)
	for _, group := range m.Groups.GetGroups() {
		for _, client := range group.Clients {
			groupNameByClient[client] = group.Name
		}
	}
	defaultGroup := m.Groups.GetById(0).Name
	getActorById := m.genActorById()
	getGroupById := m.GenGroupById()
	// GraphName returns group name, based on item actors
	statContext.GraphName = func(item *items.Item) []items.NamePct {
		var mainName string
		groupName, found := groupNameByClient[item.Client]
		if !found {
			mainName = defaultGroup
		}
		mainName = groupName
		res := []items.NamePct{}
		globPct := 1.0
		if statContext.ShowTeam && len(item.Actors) > 0 {
			pct := 1.0 / float64(len(item.Actors))
			globPct = 0.0
			for _, actId := range item.Actors {
				act := getActorById(actId)
				if act == nil { // Skip unknown or not visible Actors
					continue
				}
				grp := getGroupById(act.Groups.ActiveGroupOnDate(item.Date))
				if grp == nil {
					grp = getGroupById(0)
				}
				res = append(res, items.NamePct{
					Name: grp.Name + " : " + act.GetLabel(),
					Pct:  pct,
				})
				globPct += pct
			}
		}
		res = append(res, items.NamePct{
			Name: mainName,
			Pct:  globPct,
		})
		return res
	}
	statContext.Data3 = func(s items.StatKey) string { return items.StatSiteProgress }
	//statContext.Data3 = func(s items.StatKey) string { return s.Team}

	polesiteStats, err := statContext.CalcStats(m.Polesites, m.visibleItemizableSiteByClientFilter(), m.CurrentUser.Permissions["Invoice"])
	if err != nil {
		return err
	}
	return json.NewEncoder(writer).Encode(items.CalcProgress(polesiteStats, m.GenGroupByName(), m.GroupSizePerDays(polesiteStats.Dates)))
}
