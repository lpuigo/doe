package manager

import (
	"fmt"

	"github.com/lpuig/ewin/doe/website/backend/logger"
	"github.com/lpuig/ewin/doe/website/backend/model/actorinfos"
	"github.com/lpuig/ewin/doe/website/backend/model/actors"
	"github.com/lpuig/ewin/doe/website/backend/model/clients"
	doc "github.com/lpuig/ewin/doe/website/backend/model/doctemplate"
	fs "github.com/lpuig/ewin/doe/website/backend/model/foasites"
	"github.com/lpuig/ewin/doe/website/backend/model/groups"
	ps "github.com/lpuig/ewin/doe/website/backend/model/polesites"
	rs "github.com/lpuig/ewin/doe/website/backend/model/ripsites"
	"github.com/lpuig/ewin/doe/website/backend/model/session"
	"github.com/lpuig/ewin/doe/website/backend/model/timesheets"
	"github.com/lpuig/ewin/doe/website/backend/model/users"
	"github.com/lpuig/ewin/doe/website/backend/model/vehicules"
	ws "github.com/lpuig/ewin/doe/website/backend/model/worksites"
)

type Manager struct {
	Worksites      *ws.WorkSitesPersister
	Ripsites       *rs.SitesPersister
	Polesites      *ps.PoleSitesPersister
	Foasites       *fs.FoaSitesPersister
	Users          *users.UsersPersister
	Actors         *actors.ActorsPersister
	ActorInfos     *actorinfos.ActorInfosPersister
	TimeSheets     *timesheets.TimeSheetsPersister
	DaysOff        *Calendar
	Clients        *clients.ClientsPersister
	Groups         *groups.GroupsPersister
	Vehicules      *vehicules.VehiculesPersister
	TemplateEngine *doc.DocTemplateEngine
	SessionStore   *session.SessionStore
	CurrentUser    *users.UserRecord
	Config         ManagerConfig
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

	// Init FoaSites persister
	fsp, err := fs.NewFoaSitesPersist(conf.FoasitesDir)
	if err != nil {
		return nil, fmt.Errorf("could not create Foasites persister: %s", err.Error())
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

	// Init ActorInfos persister
	aip, err := actorinfos.NewActorInfosPersister(conf.ActorInfosDir)
	if err != nil {
		return nil, fmt.Errorf("could not create actorinfos: %s", err.Error())
	}

	// Init TimeSheets persister
	tsp, err := timesheets.NewTimeSheetsPersister(conf.TimeSheetsDir)
	if err != nil {
		return nil, fmt.Errorf("could not create timesheets: %s", err.Error())
	}

	// Init Clients persister
	cp, err := clients.NewClientsPersister(conf.ClientsDir)
	if err != nil {
		return nil, fmt.Errorf("could not create clients: %s", err.Error())
	}

	// Init Groups persister
	gp, err := groups.NewGroupsPersister(conf.GroupsDir)
	if err != nil {
		return nil, fmt.Errorf("could not create groups: %s", err.Error())
	}

	// Init Vehicules persister
	vp, err := vehicules.NewVehiculesPersister(conf.VehiculesDir)
	if err != nil {
		return nil, fmt.Errorf("could not create vehicules: %s", err.Error())
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
		Foasites:       fsp,
		Users:          up,
		Actors:         ap,
		ActorInfos:     aip,
		TimeSheets:     tsp,
		DaysOff:        NewCalendar(conf.CalendarFile),
		Clients:        cp,
		Groups:         gp,
		Vehicules:      vp,
		TemplateEngine: te,
		SessionStore:   session.NewSessionStore(conf.SessionKey),
		//CurrentUser: is set during session control transaction
		Config: conf,
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
	err := m.DaysOff.Reload()
	if err != nil {
		return fmt.Errorf("could not set holiday calendar: %s", err.Error())
	}
	logger.Entry("Server").LogInfo(fmt.Sprintf("loaded %d public holidays dates", m.DaysOff.NbDays()))

	if m.Config.IsWorksitesActive {
		err = m.Worksites.LoadDirectory()
		if err != nil {
			return fmt.Errorf("could not populate worksites: %s", err.Error())
		}
	}
	logger.Entry("Server").LogInfo(fmt.Sprintf("loaded %d Worksites", m.Worksites.NbWorksites()))

	if m.Config.IsRipsitesActive {
		err = m.Ripsites.LoadDirectory()
		if err != nil {
			return fmt.Errorf("could not populate ripsites: %s", err.Error())
		}
	}
	logger.Entry("Server").LogInfo(fmt.Sprintf("loaded %d Ripsites", m.Ripsites.NbSites()))

	if m.Config.IsPolesitesActive {
		err = m.Polesites.LoadDirectory()
		if err != nil {
			return fmt.Errorf("could not populate polesites: %s", err.Error())
		}
	}
	logger.Entry("Server").LogInfo(fmt.Sprintf("loaded %d Polesites", m.Polesites.NbSites()))

	if m.Config.IsFoasitesActive {
		err = m.Foasites.LoadDirectory()
		if err != nil {
			return fmt.Errorf("could not populate foasites: %s", err.Error())
		}
	}
	logger.Entry("Server").LogInfo(fmt.Sprintf("loaded %d Foasites", m.Foasites.NbSites()))

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

	err = m.ActorInfos.LoadDirectory()
	if err != nil {
		return fmt.Errorf("could not populate actorinfo: %s", err.Error())
	}
	logger.Entry("Server").LogInfo(fmt.Sprintf("loaded %d ActorInfos", m.ActorInfos.NbActorInfos()))

	err = m.TimeSheets.LoadDirectory()
	if err != nil {
		return fmt.Errorf("could not populate timesheet: %s", err.Error())
	}
	logger.Entry("Server").LogInfo(fmt.Sprintf("loaded %d TimeSheets", m.TimeSheets.NbTimeSheets()))

	err = m.Clients.LoadDirectory()
	if err != nil {
		return fmt.Errorf("could not populate client: %s", err.Error())
	}
	logger.Entry("Server").LogInfo(fmt.Sprintf("loaded %d Clients", m.Clients.NbClients()))

	err = m.Groups.LoadDirectory()
	if err != nil {
		return fmt.Errorf("could not populate group: %s", err.Error())
	}
	logger.Entry("Server").LogInfo(fmt.Sprintf("loaded %d Groups", m.Groups.NbGroups()))

	err = m.Vehicules.LoadDirectory()
	if err != nil {
		return fmt.Errorf("could not populate vehicule: %s", err.Error())
	}
	logger.Entry("Server").LogInfo(fmt.Sprintf("loaded %d Vehicules", m.Vehicules.NbVehicules()))

	return nil
}
