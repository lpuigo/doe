package manager

import (
	"encoding/json"
	"github.com/lpuig/ewin/doe/website/backend/model/clients"
	"github.com/lpuig/ewin/doe/website/backend/model/items"
	"io"
)

// visibleItemizableSiteByClientFilter returns a filtering function on CurrentUser.Clients visibility
func (m *Manager) visibleItemizableSiteByClientFilter() items.IsItemizableSiteVisible {
	clts, err := m.GetCurrentUserClients()
	if err != nil {
		return func(site items.ItemizableSite) bool { return false }
	}
	isVisible := make(map[string]bool)
	for _, client := range clts {
		isVisible[client.Name] = true
	}
	return func(site items.ItemizableSite) bool {
		return isVisible[site.GetClient()]
	}
}

func (m Manager) GetClients(writer io.Writer) error {
	clients := m.Clients.GetAllClients()
	return json.NewEncoder(writer).Encode(clients)
}

func (m Manager) UpdateClients(updatedClients []*clients.Client) error {
	return m.Clients.UpdateClients(updatedClients)
}
