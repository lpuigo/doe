package manager

import (
	"encoding/json"
	"fmt"
	"github.com/lpuig/ewin/doe/model"
	"github.com/lpuig/ewin/doe/website/backend/logger"
	"github.com/lpuig/ewin/doe/website/backend/model/actors"
	"github.com/lpuig/ewin/doe/website/backend/model/clients"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	doc "github.com/lpuig/ewin/doe/website/backend/model/doctemplate"
	ps "github.com/lpuig/ewin/doe/website/backend/model/polesites"
	rs "github.com/lpuig/ewin/doe/website/backend/model/ripsites"
	"github.com/lpuig/ewin/doe/website/backend/model/session"
	"github.com/lpuig/ewin/doe/website/backend/model/users"
	ws "github.com/lpuig/ewin/doe/website/backend/model/worksites"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"io"
)

type Manager struct {
	Worksites      *ws.WorkSitesPersister
	Ripsites       *rs.SitesPersister
	Polesites      *ps.PoleSitesPersister
	Users          *users.UsersPersister
	Actors         *actors.ActorsPersister
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
	logger.Entry("Server").LogInfo(fmt.Sprintf("loaded %d Worksites", wsp.NbWorsites()))

	// Init RipSites persister
	rsp, err := rs.NewSitesPersit(conf.RipsitesDir)
	if err != nil {
		return nil, fmt.Errorf("could not create Ripsites persister: %s", err.Error())
	}
	err = rsp.LoadDirectory()
	if err != nil {
		return nil, fmt.Errorf("could not populate Ripsites: %s", err.Error())
	}
	logger.Entry("Server").LogInfo(fmt.Sprintf("loaded %d Ripsites", rsp.NbSites()))

	// Init PoleSites persister
	psp, err := ps.NewPoleSitesPersist(conf.PolesitesDir)
	if err != nil {
		return nil, fmt.Errorf("could not create Polesites persister: %s", err.Error())
	}
	err = psp.LoadDirectory()
	if err != nil {
		return nil, fmt.Errorf("could not populate Polesites: %s", err.Error())
	}
	logger.Entry("Server").LogInfo(fmt.Sprintf("loaded %d Polesites", psp.NbSites()))

	// Init Users persister
	up, err := users.NewUsersPersister(conf.UsersDir)
	if err != nil {
		return nil, fmt.Errorf("could not create users: %s", err.Error())
	}
	err = up.LoadDirectory()
	if err != nil {
		return nil, fmt.Errorf("could not populate user: %s", err.Error())
	}
	logger.Entry("Server").LogInfo(fmt.Sprintf("loaded %d Users", up.NbUsers()))

	// Init Actors persister
	ap, err := actors.NewActorsPersister(conf.ActorsDir)
	if err != nil {
		return nil, fmt.Errorf("could not create actors: %s", err.Error())
	}
	err = ap.LoadDirectory()
	if err != nil {
		return nil, fmt.Errorf("could not populate actor: %s", err.Error())
	}
	logger.Entry("Server").LogInfo(fmt.Sprintf("loaded %d Actors", ap.NbActors()))

	// Init Clients persister
	cp, err := clients.NewClientsPersister(conf.ClientsDir)
	if err != nil {
		return nil, fmt.Errorf("could not create clients: %s", err.Error())
	}
	err = cp.LoadDirectory()
	if err != nil {
		return nil, fmt.Errorf("could not populate client: %s", err.Error())
	}
	logger.Entry("Server").LogInfo(fmt.Sprintf("loaded %d Clients", cp.NbClients()))

	// Init DocTemplate engine
	te, err := doc.NewDocTemplateEngine(conf.TemplatesDir)
	if err != nil {
		return nil, fmt.Errorf("could not create doc template engine", err.Error())
	}

	// Init manager
	m := &Manager{
		Worksites:      wsp,
		Ripsites:       rsp,
		Polesites:      psp,
		Users:          up,
		Actors:         ap,
		Clients:        cp,
		TemplateEngine: te,
		SessionStore:   session.NewSessionStore(conf.SessionKey),
		//CurrentUser: is set during session control transaction
	}

	return m, nil
}

func (m Manager) Clone() *Manager {
	return &m
}

// =====================================================================================================================
// User related methods
//

// genGetClient returns a GetClientByName function: func(clientName string) *clients.Client. Returned client is nil if clientName is not found
func (m *Manager) genGetClient() clients.ClientByName {
	return func(clientName string) *clients.Client {
		cr := m.Clients.GetByName(clientName)
		if cr == nil {
			return nil
		}
		return cr.Client
	}
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

// genIsTeamVisible returns a IsTeamVisible function: func(ClientTeam) bool, which is true when current user is allowed to see clientteam related activity
func (m Manager) genIsTeamVisible() (clients.IsTeamVisible, error) {
	if len(m.CurrentUser.Clients) > 0 {
		teamVisible := make(map[clients.ClientTeam]bool)
		clts, err := m.GetCurrentUserClients()
		if err != nil {
			return nil, err
		}
		for _, client := range clts {
			for _, team := range client.Teams {
				teamVisible[clients.ClientTeam{Client: client.Name, Team: team.Members}] = true
			}
		}
		return func(ct clients.ClientTeam) bool {
			return teamVisible[ct]
		}, nil
	}
	return func(clients.ClientTeam) bool { return true }, nil
}

// =====================================================================================================================
// Worksites related methods
//

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

// GetWorkSites returns array of WorksiteInfos (JSON in writer) visibles by current user
func (m Manager) GetWorksitesInfo(writer io.Writer) error {
	priceByClientArticle := m.Clients.CalcPriceByClientArticleGetter()

	wsis := []*fm.WorksiteInfo{}
	for _, wsr := range m.Worksites.GetAll(m.visibleWorksiteFilter()) {
		wsis = append(wsis, wsr.Worksite.GetInfo(priceByClientArticle))
	}

	return json.NewEncoder(writer).Encode(wsis)
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

func (m Manager) getWorksitesStats(writer io.Writer, maxVal int, dateFor date.DateAggreg) error {
	isTeamVisible, err := m.genIsTeamVisible()
	if err != nil {
		return err
	}
	return json.NewEncoder(writer).Encode(m.Worksites.GetStats(maxVal, dateFor, m.visibleWorksiteFilter(), isTeamVisible, m.genGetClient(), !m.CurrentUser.Permissions["Review"]))
}

func (m Manager) GetWorksiteXLSAttachement(writer io.Writer, ws *model.Worksite) error {
	return m.TemplateEngine.GetWorksiteXLSAttachment(writer, ws, m.genGetClient())
}

func (m Manager) WorksitesArchiveName() string {
	return m.Worksites.ArchiveName()
}

func (m Manager) CreateWorksitesArchive(writer io.Writer) error {
	return m.Worksites.CreateArchive(writer)
}

// =====================================================================================================================
// Ripsites related methods
//

// visibleRipsiteFilter returns a filtering function on CurrentUser.Clients visibility
func (m *Manager) visibleRipsiteFilter() rs.IsSiteVisible {
	if len(m.CurrentUser.Clients) == 0 {
		return func(*rs.Site) bool { return true }
	}
	isVisible := make(map[string]bool)
	for _, client := range m.CurrentUser.Clients {
		isVisible[client] = true
	}
	return func(s *rs.Site) bool {
		return isVisible[s.Client]
	}
}

// GetRipsitesInfo returns array of RipsiteInfos (JSON in writer) visibles by current user
func (m Manager) GetRipsitesInfo(writer io.Writer) error {
	rsis := []*fm.RipsiteInfo{}
	for _, rsr := range m.Ripsites.GetAll(m.visibleRipsiteFilter()) {
		rsis = append(rsis, rsr.Site.GetInfo())
	}

	return json.NewEncoder(writer).Encode(rsis)
}

func (m Manager) GetRipsiteXLSAttachement(writer io.Writer, rs *rs.Site) error {
	return m.TemplateEngine.GetRipsiteXLSAttachement(writer, rs, m.genGetClient())
}

func (m Manager) RipsitesArchiveName() string {
	return m.Ripsites.ArchiveName()
}

func (m Manager) CreateRipsitesArchive(writer io.Writer) error {
	return m.Ripsites.CreateArchive(writer)
}

// GetWorksitesWeekStats returns Worksites Stats per Week (JSON in writer) visibles by current user
func (m Manager) GetRipsitesWeekStats(writer io.Writer) error {
	df := func(d string) string {
		return date.GetMonday(d)
	}
	return m.getRipsitesStats(writer, 12, df)
}

// GetWorksitesWeekStats returns Worksites Stats per Month (JSON in writer) visibles by current user
func (m Manager) GetRipsitesMonthStats(writer io.Writer) error {
	df := func(d string) string {
		return date.GetMonth(d)
	}
	return m.getRipsitesStats(writer, 12, df)
}

func (m Manager) getRipsitesStats(writer io.Writer, maxVal int, dateFor date.DateAggreg) error {
	isTeamVisible, err := m.genIsTeamVisible()
	if err != nil {
		return err
	}
	ripsiteStats, err := m.Ripsites.GetStats(maxVal, dateFor, m.visibleRipsiteFilter(), isTeamVisible, m.genGetClient(), !m.CurrentUser.Permissions["Review"], m.CurrentUser.Permissions["Invoice"])
	if err != nil {
		return err
	}
	return json.NewEncoder(writer).Encode(ripsiteStats)
}

// =====================================================================================================================
// Ripsites related methods
//

// visibleRipsiteFilter returns a filtering function on CurrentUser.Clients visibility
func (m *Manager) visiblePolesiteFilter() ps.IsPolesiteVisible {
	if len(m.CurrentUser.Clients) == 0 {
		return func(ps *ps.PoleSite) bool { return true }
	}
	isVisible := make(map[string]bool)
	for _, client := range m.CurrentUser.Clients {
		isVisible[client] = true
	}
	return func(ps *ps.PoleSite) bool {
		return isVisible[ps.Client]
	}
}

// GetPolesitesInfo returns array of PolesiteInfos (JSON in writer) visibles by current user
func (m Manager) GetPolesitesInfo(writer io.Writer) error {
	psis := []*fm.PolesiteInfo{}
	for _, psr := range m.Polesites.GetAll(m.visiblePolesiteFilter()) {
		psis = append(psis, psr.PoleSite.GetInfo())
	}

	return json.NewEncoder(writer).Encode(psis)
}
