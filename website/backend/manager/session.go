package manager

import (
	"net/http"
)

func (m *Manager) AddSessionCookie(w http.ResponseWriter, r *http.Request) error {
	// Get a session. Get() always returns a session, even if empty.
	session, err := m.SessionStore.Get(r, "EWin-Session")
	if err != nil {
		return err
	}
	// Set some session values.
	session.Values["user"] = "test"
	// Save it before we write to the response/return from the handler.
	return session.Save(r, w)
}
