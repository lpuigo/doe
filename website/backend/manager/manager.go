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
	"github.com/lpuig/ewin/doe/website/backend/model/items"
	ps "github.com/lpuig/ewin/doe/website/backend/model/polesites"
	rs "github.com/lpuig/ewin/doe/website/backend/model/ripsites"
	"github.com/lpuig/ewin/doe/website/backend/model/session"
	"github.com/lpuig/ewin/doe/website/backend/model/users"
	ws "github.com/lpuig/ewin/doe/website/backend/model/worksites"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"io"
	"sort"
	"strconv"
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

	// Init RipSites persister
	rsp, err := rs.NewSitesPersit(conf.RipsitesDir)
	if err != nil {
		return nil, fmt.Errorf("could not create Ripsites persister: %s", err.Error())
	}

	// Init PoleSites persister
	psp, err := ps.NewPoleSitesPersist(conf.PolesitesDir)
	if err != nil {
		return nil, fmt.Errorf("could not create Polesites persister: %s", err.Error())
	}

	// Init Users persister
	up, err := users.NewUsersPersister(conf.UsersDir)
	if err != nil {
		return nil, fmt.Errorf("could not create users: %s", err.Error())
	}

	// Init Actors persister
	ap, err := actors.NewActorsPersister(conf.ActorsDir)
	if err != nil {
		return nil, fmt.Errorf("could not create actors: %s", err.Error())
	}

	// Init Clients persister
	cp, err := clients.NewClientsPersister(conf.ClientsDir)
	if err != nil {
		return nil, fmt.Errorf("could not create clients: %s", err.Error())
	}

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

	err = m.Reload()
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (m Manager) Clone() *Manager {
	return &m
}

func (m *Manager) Reload() error {
	err := m.Worksites.LoadDirectory()
	if err != nil {
		return fmt.Errorf("could not populate worksites: %s", err.Error())
	}
	logger.Entry("Server").LogInfo(fmt.Sprintf("loaded %d Worksites", m.Worksites.NbWorsites()))

	err = m.Ripsites.LoadDirectory()
	if err != nil {
		return fmt.Errorf("could not populate Ripsites: %s", err.Error())
	}
	logger.Entry("Server").LogInfo(fmt.Sprintf("loaded %d Ripsites", m.Ripsites.NbSites()))

	err = m.Polesites.LoadDirectory()
	if err != nil {
		return fmt.Errorf("could not populate Polesites: %s", err.Error())
	}
	logger.Entry("Server").LogInfo(fmt.Sprintf("loaded %d Polesites", m.Polesites.NbSites()))

	err = m.Users.LoadDirectory()
	if err != nil {
		return fmt.Errorf("could not populate user: %s", err.Error())
	}
	logger.Entry("Server").LogInfo(fmt.Sprintf("loaded %d Users", m.Users.NbUsers()))

	err = m.Actors.LoadDirectory()
	if err != nil {
		return fmt.Errorf("could not populate actor: %s", err.Error())
	}
	logger.Entry("Server").LogInfo(fmt.Sprintf("loaded %d Actors", m.Actors.NbActors()))

	err = m.Clients.LoadDirectory()
	if err != nil {
		return fmt.Errorf("could not populate client: %s", err.Error())
	}
	logger.Entry("Server").LogInfo(fmt.Sprintf("loaded %d Clients", m.Clients.NbClients()))

	return nil
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

// genActorById returns a ActorById function: func(actorId string) string. Returned string (actor ref) is "" if actorId is not found
func (m *Manager) genActorById() clients.ActorById {
	return func(actorId string) string {
		var ar *actors.ActorRecord
		if actId, err := strconv.Atoi(actorId); err == nil {
			ar = m.Actors.GetById(actId)
		} else {
			ar = m.Actors.GetByRef(actorId)
		}
		if ar == nil {
			return ""
		}
		return "(" + ar.Role + ") " + ar.Actor.Ref
	}
}

// genActorInfoById returns a ActorInfoById function: func(actorId string) []string which returns nil if actorId is not known, or [0] Actor Role [1] Actor Ref
func (m *Manager) genActorInfoById() clients.ActorInfoById {
	return func(actorId string) []string {
		var ar *actors.ActorRecord
		if actId, err := strconv.Atoi(actorId); err == nil {
			ar = m.Actors.GetById(actId)
		} else {
			ar = m.Actors.GetByRef(actorId)
		}
		if ar == nil {
			return nil
		}
		return []string{ar.Role, ar.Actor.Ref}
	}
}

// GetCurrentUserClientsName returns Clients' names visible by current user (if user has no client, returns all clients)
func (m Manager) GetCurrentUserClientsName() []string {
	if m.CurrentUser == nil {
		return nil
	}
	if len(m.CurrentUser.Clients) > 0 {
		return m.CurrentUser.Clients
	}
	clientsNames := []string{}
	for _, client := range m.Clients.GetAllClients() {
		clientsNames = append(clientsNames, client.Name)
	}
	return clientsNames
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
	if len(m.CurrentUser.Clients) == 0 {
		return func(clients.ClientTeam) bool { return true }, nil
	}

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

// genIsActorVisible returns a IsTeamVisible function: func(ClientTeam) bool, which is true when current user is allowed to see clientteam (by actorId) related activity
func (m Manager) genIsActorVisible() (clients.IsTeamVisible, error) {
	if len(m.CurrentUser.Clients) == 0 {
		return func(clients.ClientTeam) bool { return true }, nil
	}

	actorVisible := make(map[clients.ClientTeam]bool)
	clts, err := m.GetCurrentUserClients()
	if err != nil {
		return nil, err
	}
	for _, client := range clts {
		allowedActors := m.Actors.GetActorsByClient(false, client.Name)
		for _, actor := range allowedActors {
			actorVisible[clients.ClientTeam{Client: client.Name, Team: strconv.Itoa(actor.Id)}] = true
			actorVisible[clients.ClientTeam{Client: client.Name, Team: actor.LastName}] = true
		}
	}
	return func(ct clients.ClientTeam) bool {
		return actorVisible[ct]
	}, nil
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

// GetWorksitesStats returns  visibles by current user Worksites Stats per Freq (week or month) as JSON in writer
func (m Manager) GetWorksitesStats(writer io.Writer, info, freq string) error {
	maxVal := 12

	var dateFor date.DateAggreg
	switch freq {
	case "week":
		dateFor = func(d string) string {
			return date.GetMonday(d)
		}
	case "month":
		dateFor = func(d string) string {
			return date.GetMonth(d)
		}
	default:
		return fmt.Errorf("unsupported stat period '%s'", freq)
	}

	isTeamVisible, err := m.genIsTeamVisible()
	if err != nil {
		return err
	}
	switch info {
	case "prod":
		return json.NewEncoder(writer).Encode(m.Worksites.GetStats(maxVal, dateFor, m.visibleWorksiteFilter(), isTeamVisible, m.genGetClient(), !m.CurrentUser.Permissions["Review"], false))
	case "stock":
		return json.NewEncoder(writer).Encode(m.Worksites.GetStockStats(maxVal, dateFor, m.visibleWorksiteFilter(), isTeamVisible, m.genGetClient()))
	default:
		return fmt.Errorf("unsupported info '%s'", info)
	}
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
	clientByName := m.genGetClient()
	rsis := []*fm.RipsiteInfo{}
	for _, rsr := range m.Ripsites.GetAll(m.visibleRipsiteFilter()) {
		rsis = append(rsis, rsr.Site.GetInfo(clientByName))
	}

	return json.NewEncoder(writer).Encode(rsis)
}

func (m Manager) GetRipsiteXLSAttachement(writer io.Writer, rs *rs.Site) error {
	return m.TemplateEngine.GetRipsiteXLSAttachement(writer, rs, m.genGetClient(), m.genActorById())
}

func (m Manager) RipsitesArchiveName() string {
	return m.Ripsites.ArchiveName()
}

func (m Manager) CreateRipsitesArchive(writer io.Writer) error {
	return m.Ripsites.CreateArchive(writer)
}

func (m Manager) GetRipsitesStats(writer io.Writer, freq, groupBy string) error {
	maxVal := 12

	var dateFor date.DateAggreg
	switch freq {
	case "week":
		dateFor = func(d string) string {
			return date.GetMonday(d)
		}
	case "month":
		dateFor = func(d string) string {
			return date.GetMonth(d)
		}
	default:
		return fmt.Errorf("unsupported stat period '%s'", freq)
	}

	isActorVisible, err := m.genIsActorVisible()
	if err != nil {
		return err
	}

	statContext := items.StatContext{
		MaxVal:        maxVal,
		DateFor:       dateFor,
		IsTeamVisible: isActorVisible,
		ShowTeam:      !m.CurrentUser.Permissions["Review"],
	}

	switch groupBy {
	case "activity", "site":
		ripsiteStats, err := m.Ripsites.GetProdStats(statContext, m.visibleRipsiteFilter(), m.genGetClient(), m.genActorById(), m.CurrentUser.Permissions["Invoice"], groupBy)
		if err != nil {
			return err
		}
		return json.NewEncoder(writer).Encode(ripsiteStats)

	case "mean":
		ripsiteStats, err := m.Ripsites.GetMeanProdStats(statContext, m.visibleRipsiteFilter(), m.genGetClient(), m.genActorInfoById())
		if err != nil {
			return err
		}
		meanStats := items.CalcTeamMean(ripsiteStats, 1)
		return json.NewEncoder(writer).Encode(meanStats)

	default:
		return fmt.Errorf("unsupported group type '%s'", groupBy)
	}
}

func (m Manager) GetRipsitesActorsActivity(writer io.Writer, freq string) error {

	var dateFor date.DateAggreg
	var firstDate string
	switch freq {
	case "week":
		dateFor = func(d string) string {
			return date.GetMonday(d)
		}
	case "month":
		dateFor = func(d string) string {
			return date.GetMonth(d)
		}
	default:
		return fmt.Errorf("unsupported stat period '%s'", freq)
	}

	// set firstDate according to freq choice, in order to have at least a full month of data
	// month : last and current month
	// week : 5 last weeks and current
	firstDate = dateFor(date.Today().AddDays(-32).String())

	itms, err := m.Ripsites.GetAllItems(firstDate, dateFor, m.visibleRipsiteFilter(), m.genGetClient())
	if err != nil {
		return err
	}
	return m.TemplateEngine.GetItemsXLSAttachement(writer, itms, m.genActorById())
}

// =====================================================================================================================
// Polesites related methods
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

func (m Manager) GetPolesiteXLSAttachement(writer io.Writer, ps *ps.PoleSite) error {
	return m.TemplateEngine.GetPolesiteXLSAttachement(writer, ps, m.genGetClient(), m.genActorById())
}

func (m Manager) PolesitesArchiveName() string {
	return m.Polesites.ArchiveName()
}

func (m Manager) CreatePolesitesArchive(writer io.Writer) error {
	return m.Polesites.CreateArchive(writer)
}

// GetPolesitesWeekStats returns Polesites Stats per Week (JSON in writer) visibles by current user
func (m Manager) GetPolesitesWeekStats(writer io.Writer) error {
	df := func(d string) string {
		return date.GetMonday(d)
	}
	return m.getPolesitesStats(writer, 12, df)
}

// GetPolesitesMonthStats returns Polesites Stats per Month (JSON in writer) visibles by current user
func (m Manager) GetPolesitesMonthStats(writer io.Writer) error {
	df := func(d string) string {
		return date.GetMonth(d)
	}
	return m.getPolesitesStats(writer, 12, df)
}

func (m Manager) getPolesitesStats(writer io.Writer, maxVal int, dateFor date.DateAggreg) error {
	isActorVisible, err := m.genIsActorVisible()
	if err != nil {
		return err
	}

	statContext := items.StatContext{
		MaxVal:        maxVal,
		DateFor:       dateFor,
		IsTeamVisible: isActorVisible,
		ShowTeam:      !m.CurrentUser.Permissions["Review"],
	}

	polesiteStats, err := m.Polesites.GetStats(statContext, m.visiblePolesiteFilter(), m.genGetClient(), m.genActorById(), m.CurrentUser.Permissions["Invoice"])
	if err != nil {
		return err
	}
	return json.NewEncoder(writer).Encode(polesiteStats)
}

// =====================================================================================================================
// Actors related methods
//

func (m Manager) GetActors(writer io.Writer) error {
	clientsNames := m.GetCurrentUserClientsName()
	actors := m.Actors.GetActorsByClient(false, clientsNames...)
	return json.NewEncoder(writer).Encode(actors)
}

func (m Manager) ActorsArchiveName() string {
	return m.Actors.ArchiveName()
}

func (m Manager) CreateActorsArchive(writer io.Writer) error {
	return m.Actors.CreateArchive(writer)
}

func (m Manager) GetActorsWorkingHoursRecordXLSName(monthDate string) string {
	return fmt.Sprintf("CRA %s.xlsx", monthDate)
}

func (m Manager) GetActorsWorkingHoursRecordXLS(writer io.Writer, date string) error {
	actors := m.Actors.GetAllActors()
	sort.Slice(actors, func(i, j int) bool {
		return actors[i].Ref < actors[j].Ref
	})
	return m.TemplateEngine.GetActorsWorkingHoursRecordXLS(writer, date, actors)
}

// =====================================================================================================================
// Clients related methods
//

func (m Manager) ClientsArchiveName() string {
	return m.Clients.ArchiveName()
}

func (m Manager) CreateClientsArchive(writer io.Writer) error {
	return m.Clients.CreateArchive(writer)
}
