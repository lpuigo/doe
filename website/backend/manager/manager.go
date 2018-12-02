package manager

import (
	"encoding/json"
	"fmt"
	"github.com/lpuig/ewin/doe/model"
	ws "github.com/lpuig/ewin/doe/website/backend/model/worksites"
	"io"
)

type Manager struct {
	Worksites *ws.WorkSitesPersister
}

func NewManager(worsitesDir string) (*Manager, error) {
	wsp, err := ws.NewWorkSitesPersist(worsitesDir)
	if err != nil {
		return nil, fmt.Errorf("could not create worksites: %s", err.Error())
	}
	err = wsp.LoadDirectory()
	if err != nil {
		return nil, fmt.Errorf("could not populate worksites:%s", err.Error())
	}
	m := &Manager{Worksites: wsp}

	return m, nil
}

func (m Manager) GetWorkSites(writer io.Writer) error {
	return json.NewEncoder(writer).Encode(m.Worksites.GetAll(func(ws *model.Worksite) bool { return true }))
}
