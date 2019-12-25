package manager

import (
	"encoding/json"
	"fmt"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	fs "github.com/lpuig/ewin/doe/website/backend/model/foasites"
	"github.com/lpuig/ewin/doe/website/backend/model/items"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"io"
)

// visiblePolesiteFilter returns a filtering function on CurrentUser.Clients visibility
func (m *Manager) visibleFoaItemizableSiteFilter() items.IsItemizableSiteVisible {
	if len(m.CurrentUser.Clients) == 0 {
		return func(site items.ItemizableSite) bool { return true }
	}
	isVisible := make(map[string]bool)
	for _, client := range m.CurrentUser.Clients {
		isVisible[client] = true
	}
	return func(site items.ItemizableSite) bool {
		return isVisible[site.GetClient()]
	}
}

// GetFoaSitesInfo returns array of FoaSiteInfos (JSON in writer) visibles by current user
func (m Manager) GetFoaSitesInfo(writer io.Writer) error {
	fsis := []*fm.FoaSiteInfo{}
	for _, fsr := range m.Foasites.GetAll(m.visibleFoaItemizableSiteFilter()) {
		fsis = append(fsis, fsr.FoaSite.GetInfo())
	}

	return json.NewEncoder(writer).Encode(fsis)
}

func (m Manager) GetFoaSitesStats(writer io.Writer, freq string) error {
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
		ShowTeam:      false,
	}
	foaStats, err := statContext.CalcStats(m.Foasites, m.visibleFoaItemizableSiteFilter(), m.CurrentUser.Permissions["Invoice"])
	if err != nil {
		return err
	}
	return json.NewEncoder(writer).Encode(foaStats)
}

func (m Manager) GetFoaSiteXLSAttachement(writer io.Writer, site *fs.FoaSite) error {
	return m.TemplateEngine.GetFoaSiteXLSAttachement(writer, site, m.genGetClient(), m.genActorById())
}
