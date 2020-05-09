package manager

import (
	"encoding/json"
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
	statContext, err := m.NewStatContext(freq)
	if err != nil {
		return err
	}

	foaStats, err := statContext.CalcStats(m.Foasites, m.visibleItemizableSiteByClientFilter(), m.CurrentUser.Permissions["Invoice"])
	if err != nil {
		return err
	}
	return json.NewEncoder(writer).Encode(foaStats)
}
