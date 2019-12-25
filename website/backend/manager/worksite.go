package manager

import (
	"encoding/json"
	"fmt"
	"github.com/lpuig/ewin/doe/model"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"

	"io"
)

// visibleWorksiteFilter returns a filtering function on CurrentUser.Clients visibility
func (m *Manager) visibleWorksiteFilter() model.IsWSVisible {
	if len(m.CurrentUser.Clients) == 0 {
		return func(ws *model.Worksite) bool { return true }
	}
	isVisible := make(map[string]bool)
	for _, client := range m.CurrentUser.Clients {
		isVisible[client] = true
	}
	return func(ws *model.Worksite) bool {
		return isVisible[ws.Client]
	}
}

// GetWorkSites returns array of WorksiteInfos (JSON in writer) visibles by current user
func (m Manager) GetWorksitesInfo(writer io.Writer) error {
	priceByClientArticle := m.Clients.CalcPriceByClientArticleGetter()

	wsis := []*fm.WorksiteInfo{}
	for _, wsr := range m.Worksites.GetAll(m.visibleWorksiteFilter()) {
		wsis = append(wsis, wsr.Worksite.GetInfo(priceByClientArticle))
	}

	return json.NewEncoder(writer).Encode(wsis)
}

// GetWorksitesStats returns  visibles by current user Worksites Stats per Freq (week or month) as JSON in writer
func (m Manager) GetWorksitesStats(writer io.Writer, info, freq string) error {
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

	isTeamVisible, err := m.genIsTeamVisible()
	if err != nil {
		return err
	}
	switch info {
	case "prod":
		return json.NewEncoder(writer).Encode(m.Worksites.GetStats(maxVal, dateFor, m.visibleWorksiteFilter(), isTeamVisible, m.genGetClient(), !m.CurrentUser.Permissions["Review"], false))
	case "stock":
		return json.NewEncoder(writer).Encode(m.Worksites.GetStockStats(maxVal, dateFor, m.visibleWorksiteFilter(), isTeamVisible, m.genGetClient()))
	default:
		return fmt.Errorf("unsupported info '%s'", info)
	}
}

func (m Manager) GetWorksiteXLSAttachement(writer io.Writer, ws *model.Worksite) error {
	return m.TemplateEngine.GetWorksiteXLSAttachment(writer, ws, m.genGetClient())
}
