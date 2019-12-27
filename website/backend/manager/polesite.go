package manager

import (
	"encoding/json"
	"fmt"
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
		ClientByName:  m.genGetClient(),
		ActorById:     m.genActorById(),
		ShowTeam:      !m.CurrentUser.Permissions["Review"],
	}
	polesiteStats, err := statContext.CalcStats(m.Polesites, m.visibleItemizableSiteByClientFilter(), m.CurrentUser.Permissions["Invoice"])
	if err != nil {
		return err
	}
	return json.NewEncoder(writer).Encode(polesiteStats)
}
