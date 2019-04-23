package manager

import (
	"encoding/json"
	"fmt"
	"github.com/lpuig/ewin/doe/model"
	"github.com/lpuig/ewin/doe/website/backend/logger"
	"github.com/lpuig/ewin/doe/website/backend/model/clients"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	doc "github.com/lpuig/ewin/doe/website/backend/model/doctemplate"
	"github.com/lpuig/ewin/doe/website/backend/model/session"
	"github.com/lpuig/ewin/doe/website/backend/model/users"
	ws "github.com/lpuig/ewin/doe/website/backend/model/worksites"
	"io"
)

type Manager struct {
	Worksites      *ws.WorkSitesPersister
	Users          *users.UsersPersister
	Clients        *clients.ClientsPersister
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
		return nil, fmt.Errorf("could not populate worksites: %s", err.Error())
	}

	// Init Users persister
	up, err := users.NewUsersPersister(conf.UsersDir)
	if err != nil {
		return nil, fmt.Errorf("could not create users: %s", err.Error())
	}
	err = up.LoadDirectory()
	if err != nil {
		return nil, fmt.Errorf("could not populate user: %s", err.Error())
	}

	// Init Clients persister
	cp, err := clients.NewClientsPersister(conf.ClientsDir)
	if err != nil {
		return nil, fmt.Errorf("could not create clients: %s", err.Error())
	}
	err = cp.LoadDirectory()
	if err != nil {
		return nil, fmt.Errorf("could not populate client: %s", err.Error())
	}

	// Init DocTemplate engine
	te, err := doc.NewDocTemplateEngine(conf.TemplatesDir)
	if err != nil {
		return nil, fmt.Errorf("could not create doc template engine", err.Error())
	}

	// Init manager
	m := &Manager{
		Worksites:      wsp,
		Users:          up,
		Clients:        cp,
		TemplateEngine: te,
	}
	logger.Entry("Server").LogInfo(
		fmt.Sprintf("loaded %d Worsites, %d Clients and %d users",
			wsp.NbWorsites(),
			cp.NbClients(),
			up.NbUsers(),
		))

	m.SessionStore = session.NewSessionStore(conf.SessionKey)

	// m.CurrentUser is set transaction during session control

	return m, nil
}

func (m Manager) Clone() *Manager {
	return &m
}

// visibleWorksiteFilter returns a filtering function on CurrentUser.Clients visibility
func (m *Manager) visibleWorksiteFilter() model.IsWSVisible {
	if len(m.CurrentUser.Clients) == 0 {
		return func(ws *model.Worksite) bool { return true }
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
//func (m Manager) GetWorkSites(writer io.Writer) error {
//	return json.NewEncoder(writer).Encode(m.Worksites.GetAll(func(ws *model.Worksite) bool { return true }))
//}

// GetWorkSites returns array of WorksiteInfos (JSON in writer) visibles by current user
func (m Manager) GetWorksitesInfo(writer io.Writer) error {
	return json.NewEncoder(writer).Encode(m.Worksites.GetAllInfo(m.visibleWorksiteFilter()))
}

// GetWorksitesWeekStats returns Worksites Stats per Week (JSON in writer) visibles by current user
func (m Manager) GetWorksitesWeekStats(writer io.Writer) error {
	df := func(d string) string {
		return date.GetMonday(d)
	}
	return m.getWorksitesStats(writer, 12, df)
}

// GetWorksitesWeekStats returns Worksites Stats per Month (JSON in writer) visibles by current user
func (m Manager) GetWorksitesMonthStats(writer io.Writer) error {
	df := func(d string) string {
		return date.GetMonth(d)
	}
	return m.getWorksitesStats(writer, 12, df)
}

func (m Manager) getWorksitesStats(writer io.Writer, maxVal int, dateFor model.DateAggreg) error {
	var isTeamVisible model.IsTeamVisible
	if len(m.CurrentUser.Clients) > 0 {
		teamVisible := make(map[model.ClientTeam]bool)
		clts, err := m.GetCurrentUserClients()
		if err != nil {
			return err
		}
		for _, client := range clts {
			for _, team := range client.Teams {
				teamVisible[model.ClientTeam{Client: client.Name, Team: team.Members}] = true
			}
		}
		isTeamVisible = func(ct model.ClientTeam) bool {
			return teamVisible[ct]
		}
	} else {
		isTeamVisible = func(model.ClientTeam) bool { return true }
	}
	return json.NewEncoder(writer).Encode(m.Worksites.GetStats(maxVal, dateFor, m.visibleWorksiteFilter(), isTeamVisible, !m.CurrentUser.Permissions["Review"]))
}

func (m Manager) ArchiveName() string {
	return m.Worksites.ArchiveName()
}

func (m Manager) CreateArchive(writer io.Writer) error {
	return m.Worksites.CreateArchive(writer)
}

// GetCurrentUserClients returns Clients visible by current user (if user has no client, returns all clients)
func (m Manager) GetCurrentUserClients() ([]*clients.Client, error) {
	res := []*clients.Client{}
	if m.CurrentUser == nil {
		return nil, nil
	}
	if len(m.CurrentUser.Clients) == 0 {
		return m.Clients.GetAllClients(), nil
	}
	for _, clientName := range m.CurrentUser.Clients {
		client := m.Clients.GetByName(clientName)
		if client == nil {
			return nil, fmt.Errorf("could not retrieve client '%s' info", clientName)
		}
		res = append(res, client.Client)
	}
	return res, nil
}
