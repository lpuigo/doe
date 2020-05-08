package manager

import (
	"encoding/json"
	"fmt"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	"github.com/lpuig/ewin/doe/website/backend/model/items"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"io"
)

// GetFoaSitesInfo returns array of FoaSiteInfos (JSON in writer) visibles by current user
func (m Manager) GetFoaSitesInfo(writer io.Writer) error {
	fsis := []*fm.FoaSiteInfo{}
	for _, fsr := range m.Foasites.GetAll(m.visibleItemizableSiteByClientFilter()) {
		fsis = append(fsis, fsr.FoaSite.GetInfo())
	}

	return json.NewEncoder(writer).Encode(fsis)
}

func (m Manager) GetFoaSitesStats(writer io.Writer, freq string) error {
	maxVal := 12
	dayIncr := 7
	startDate := date.Today().String()

	var dateFor date.DateAggreg
	switch freq {
	case "day":
		maxVal = 15
		dayIncr = 1
		dateFor = func(d string) string {
			return d
		}
		startDate = date.Today().AddDays(1 - maxVal).String()
	case "week":
		dateFor = func(d string) string {
			return date.GetMonday(d)
		}
		startDate = date.GetDateAfter(date.GetMonday(startDate), (1-maxVal)*7)
	case "month":
		dateFor = func(d string) string {
			return date.GetMonth(d)
		}
		startDate = date.GetMonth(date.GetDateAfter(startDate, -maxVal*30))
	default:
		return fmt.Errorf("unsupported stat period '%s'", freq)
	}

	isActorVisible, err := m.genIsActorVisible()
	if err != nil {
		return err
	}

	statContext := items.StatContext{
		DayIncr:       dayIncr,
		MaxVal:        maxVal,
		StartDate:     startDate,
		DateFor:       dateFor,
		IsTeamVisible: isActorVisible,
		ClientByName:  m.genGetClient(),
		ActorById:     m.genActorById(),
		ShowTeam:      false,
	}
	foaStats, err := statContext.CalcStats(m.Foasites, m.visibleItemizableSiteByClientFilter(), m.CurrentUser.Permissions["Invoice"])
	if err != nil {
		return err
	}
	return json.NewEncoder(writer).Encode(foaStats)
}
