package manager

import (
	"encoding/json"
	"fmt"
	"github.com/lpuig/ewin/doe/model"
	"github.com/lpuig/ewin/doe/website/backend/logger"
	"github.com/lpuig/ewin/doe/website/backend/model/session"
	"github.com/lpuig/ewin/doe/website/backend/model/users"
	ws "github.com/lpuig/ewin/doe/website/backend/model/worksites"
	"io"
)

type Manager struct {
	Worksites    *ws.WorkSitesPersister
	Users        *users.UsersPersister
	SessionStore *session.SessionStore
	CurrentUser  *users.UserRecord
}

func NewManager(conf ManagerConfig) (*Manager, error) {
	// Init Worksites persister
	wsp, err := ws.NewWorkSitesPersist(conf.WorksitesDir)
	if err != nil {
		return nil, fmt.Errorf("could not create worksites: %s", err.Error())
	}
	err = wsp.LoadDirectory()
	if err != nil {
		return nil, fmt.Errorf("could not populate worksites:%s", err.Error())
	}

	// Init Users persister
	up, err := users.NewUsersPersister(conf.UsersDir)
	if err != nil {
		return nil, fmt.Errorf("could not create users: %s", err.Error())
	}
	err = up.LoadDirectory()
	if err != nil {
		return nil, fmt.Errorf("could not populate user:%s", err.Error())
	}

	// Init manager
	m := &Manager{Worksites: wsp, Users: up}
	logger.Entry("Server").LogInfo(fmt.Sprintf("loaded %d Worsites and %d users", wsp.NbWorsites(), up.NbUsers()))

	m.SessionStore = session.NewSessionStore(conf.SessionKey)

	return m, nil
}

// GetWorkSites return Arrays of Worksites (JSON in writer)
func (m Manager) GetWorkSites(writer io.Writer) error {
	return json.NewEncoder(writer).Encode(m.Worksites.GetAll(func(ws *model.Worksite) bool { return true }))
}

// GetWorkSites return Arrays of WorksiteInfos (JSON in writer)
func (m Manager) GetWorksitesInfo(writer io.Writer) error {
	return json.NewEncoder(writer).Encode(m.Worksites.GetAllInfo(func(ws *model.Worksite) bool { return true }))
}

func (m Manager) Clone() *Manager {
	return &m
}
