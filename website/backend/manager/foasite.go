package manager

import (
	"encoding/json"
	fs "github.com/lpuig/ewin/doe/website/backend/model/foasites"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"io"
)

// visiblePolesiteFilter returns a filtering function on CurrentUser.Clients visibility
func (m *Manager) visibleFoasiteFilter() fs.IsFoasiteVisible {
	if len(m.CurrentUser.Clients) == 0 {
		return func(fs *fs.FoaSite) bool { return true }
	}
	isVisible := make(map[string]bool)
	for _, client := range m.CurrentUser.Clients {
		isVisible[client] = true
	}
	return func(fs *fs.FoaSite) bool {
		return isVisible[fs.Client]
	}
}

// GetFoaSitesInfo returns array of FoaSiteInfos (JSON in writer) visibles by current user
func (m Manager) GetFoaSitesInfo(writer io.Writer) error {
	fsis := []*fm.FoaSiteInfo{}
	for _, fsr := range m.Foasites.GetAll(m.visibleFoasiteFilter()) {
		fsis = append(fsis, fsr.FoaSite.GetInfo())
	}

	return json.NewEncoder(writer).Encode(fsis)
}

func (m Manager) FoaSitesArchiveName() string {
	return m.Foasites.ArchiveName()
}

func (m Manager) CreateFoaSitesArchive(writer io.Writer) error {
	return m.Foasites.CreateArchive(writer)
}

func (m Manager) GetFoaSiteXLSAttachement(writer io.Writer, site *fs.FoaSite) error {
	return m.TemplateEngine.GetFoaSiteXLSAttachement(writer, site, m.genGetClient(), m.genActorById())
}
