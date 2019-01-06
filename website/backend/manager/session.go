package manager

import (
	"net/http"
)

// CheckSessionUser checks for Session Cookie provided in the request.
//
// Returns true if User authenticated (and place current user in manager) false otherwise
func (m *Manager) CheckSessionUser(r *http.Request) bool {
	userid := m.SessionStore.CheckUser(r)
	if userid == -1 {
		m.CurrentUser = nil
		return false
	}
	ur := m.Users.GetById(userid)
	m.CurrentUser = ur
	if ur == nil {
		return false
	}
	return true
}
