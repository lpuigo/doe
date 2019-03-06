package manager

import (
	"encoding/json"
	"fmt"
	"github.com/lpuig/ewin/doe/model"
	"github.com/lpuig/ewin/doe/website/backend/logger"
	doc "github.com/lpuig/ewin/doe/website/backend/model/doctemplate"
	"github.com/lpuig/ewin/doe/website/backend/model/session"
	"github.com/lpuig/ewin/doe/website/backend/model/users"
	ws "github.com/lpuig/ewin/doe/website/backend/model/worksites"
	"io"
)

type Manager struct {
	Worksites      *ws.WorkSitesPersister
	Users          *users.UsersPersister
	TemplateEngine *doc.DocTemplateEngine
	SessionStore   *session.SessionStore
	CurrentUser    *users.UserRecord
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

	// Init DocTemplate engine
	te, err := doc.NewDocTemplateEngine(conf.TemplatesDir)
	if err != nil {
		return nil, fmt.Errorf("could not create doc template engine", err.Error())
	}

	// Init manager
	m := &Manager{Worksites: wsp, Users: up, TemplateEngine: te}
	logger.Entry("Server").LogInfo(fmt.Sprintf("loaded %d Worsites and %d users", wsp.NbWorsites(), up.NbUsers()))

	m.SessionStore = session.NewSessionStore(conf.SessionKey)

	// m.CurrentUser is set transaction during session control

	return m, nil
}

func (m Manager) Clone() *Manager {
	return &m
}

func (m *Manager) visibleWorksiteFilter() func(*model.Worksite) bool {
	if len(m.CurrentUser.Clients) == 0 {
		return func(*model.Worksite) bool { return true }
	}
	isVisible := make(map[string]bool)
	for _, client := range m.CurrentUser.Clients {
		isVisible[client] = true
	}
	return func(ws *model.Worksite) bool {
		return isVisible[ws.Client]
	}
}

// GetWorkSites returns Arrays of Worksites (JSON in writer)
func (m Manager) GetWorkSites(writer io.Writer) error {
	return json.NewEncoder(writer).Encode(m.Worksites.GetAll(func(ws *model.Worksite) bool { return true }))
}

// GetWorkSites returns array of WorksiteInfos (JSON in writer) visibles by current user
func (m Manager) GetWorksitesInfo(writer io.Writer) error {
	return json.NewEncoder(writer).Encode(m.Worksites.GetAllInfo(m.visibleWorksiteFilter()))
}

// GetWorkSitesStats returns Worksites Stats (JSON in writer) visibles by current user
func (m Manager) GetWorkSitesStats(writer io.Writer) error {
	var isTeamVisible func(team string) bool
	if len(m.CurrentUser.Teams) > 0 {
		teamVisible := make(map[string]bool)
		for _, team := range m.CurrentUser.Teams {
			teamVisible[team] = true
		}
		isTeamVisible = func(team string) bool {
			return teamVisible[team]
		}
	} else {
		isTeamVisible = func(team string) bool { return true }
	}
	return json.NewEncoder(writer).Encode(m.Worksites.GetStats(m.visibleWorksiteFilter(), isTeamVisible))
}

func (m Manager) ArchiveName() string {
	return m.Worksites.ArchiveName()
}

func (m Manager) CreateArchive(writer io.Writer) error {
	return m.Worksites.CreateArchive(writer)
}
