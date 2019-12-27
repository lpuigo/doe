package manager

import "github.com/lpuig/ewin/doe/website/backend/model/items"

// visiblePolesiteFilter returns a filtering function on CurrentUser.Clients visibility
func (m *Manager) visibleItemizableSiteByClientFilter() items.IsItemizableSiteVisible {
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
