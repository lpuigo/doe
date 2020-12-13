package manager

import "github.com/lpuig/ewin/doe/website/backend/model/items"

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
