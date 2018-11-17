package manager

import (
	"fmt"
	ws "github.com/lpuig/ewin/doe/website/backend/model/worksites"
	"net/http"
)

type Manager struct {
	Worksites *ws.WorkSitesPersister
}

func NewManager(worsitesDir string) (*Manager, error) {
	wsp := ws.NewWorkSitesPersist(worsitesDir)
	if err := wsp.CheckDirectory(); err != nil {
		return nil, fmt.Errorf("worksites:%s", err.Error())
	}
	m := &Manager{Worksites: wsp}

	return m, nil
}

func (m Manager) GetWorkSites(writer http.ResponseWriter) {
	m.Worksites.GetAll()
}
