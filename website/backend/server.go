package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/lpuig/ewin/doe/website/backend/config"
	"github.com/lpuig/ewin/doe/website/backend/logger"
	"github.com/lpuig/ewin/doe/website/backend/manager"
	"github.com/lpuig/ewin/doe/website/backend/route"
	"net/http"
	"os/exec"
)

type Conf struct {
	manager.ManagerConfig

	LogFile     string
	ServicePort string
	AssetsDir   string
	AssetsRoot  string
	RootDir     string

	LaunchWebBrowser bool
}

const (
	AssetsDir  = `../../WebAssets`
	AssetsRoot = `/Assets/`
	RootDir    = `./Dist`

	ServicePort = ":8080"
	SessionKey  = "SECRET_KEY"

	WorksitesDir  = `C:\Users\Laurent\Golang\src\github.com\lpuig\ewin\doe\Ressources\Worksites`
	RipsitesDir   = `C:\Users\Laurent\Golang\src\github.com\lpuig\ewin\doe\Ressources\Ripsites`
	PolesitesDir  = `C:\Users\Laurent\Golang\src\github.com\lpuig\ewin\doe\Ressources\Polesites`
	FoasitesDir   = `C:\Users\Laurent\Golang\src\github.com\lpuig\ewin\doe\Ressources\Foasites`
	UsersDir      = `C:\Users\Laurent\Golang\src\github.com\lpuig\ewin\doe\Ressources\Users`
	ActorsDir     = `C:\Users\Laurent\Golang\src\github.com\lpuig\ewin\doe\Ressources\Actors`
	TimeSheetsDir = `C:\Users\Laurent\Golang\src\github.com\lpuig\ewin\doe\Ressources\Timesheets`
	CalendarFile  = `C:\Users\Laurent\Golang\src\github.com\lpuig\ewin\doe\Ressources\Calendar\holidays.json`
	ClientsDir    = `C:\Users\Laurent\Golang\src\github.com\lpuig\ewin\doe\Ressources\Clients`
	TemplatesDir  = `C:\Users\Laurent\Golang\src\github.com\lpuig\ewin\doe\Ressources\DocTemplates`

	LaunchWebBrowser = true

	ConfigFile = `./config.json`
	LogFile    = `./server.log`
)

func main() {
	conf := &Conf{
		ManagerConfig: manager.ManagerConfig{
			WorksitesDir:  WorksitesDir,
			RipsitesDir:   RipsitesDir,
			PolesitesDir:  PolesitesDir,
			FoasitesDir:   FoasitesDir,
			UsersDir:      UsersDir,
			ActorsDir:     ActorsDir,
			TimeSheetsDir: TimeSheetsDir,
			CalendarFile:  CalendarFile,
			ClientsDir:    ClientsDir,
			TemplatesDir:  TemplatesDir,
			SessionKey:    SessionKey,
		},
		LogFile:          LogFile,
		ServicePort:      ServicePort,
		AssetsDir:        AssetsDir,
		AssetsRoot:       AssetsRoot,
		RootDir:          RootDir,
		LaunchWebBrowser: LaunchWebBrowser,
	}

	if err := config.SetFromFile(ConfigFile, conf); err != nil {
		logger.Entry("Server").Fatal(err)
	}

	logFile := logger.StartLog(conf.LogFile)
	defer logFile.Close()
	logger.Entry("Server").LogInfo("============================= SERVER STARTING ==================================")

	mgr, err := manager.NewManager(conf.ManagerConfig)
	if err != nil {
		logger.Entry("Server").Fatal(err)
	}

	withManager := func(hf route.MgrHandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			hf(mgr.Clone(), w, r)
		}
	}

	withUserManager := func(request string, hf route.MgrHandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			m := mgr.Clone()
			if !m.CheckSessionUser(r) {
				logmsg := logger.Entry("Route").AddRequest(request)
				route.AddError(w, logmsg, "User not connected or not authorized", http.StatusUnauthorized)
				logmsg.Log()
				return
			}
			hf(m, w, r)
		}
	}

	router := mux.NewRouter()
	// session management
	router.HandleFunc("/api/login", withManager(route.GetUser)).Methods("GET")
	router.HandleFunc("/api/login", withUserManager("Logout", route.Logout)).Methods("DELETE")
	router.HandleFunc("/api/login", withManager(route.Login)).Methods("POST")

	// Worksites methods
	router.HandleFunc("/api/worksites", withUserManager("GetWorksitesInfo", route.GetWorksitesInfo)).Methods("GET")
	router.HandleFunc("/api/worksites", withUserManager("CreateWorkSite", route.CreateWorkSite)).Methods("POST")
	router.HandleFunc("/api/worksites/archive", withUserManager("GetWorksitesArchive", route.GetWorksitesArchive)).Methods("GET")
	router.HandleFunc("/api/worksites/stat/{info}/{freq}", withUserManager("GetWorksitesStats", route.GetWorksitesStats)).Methods("GET")
	router.HandleFunc("/api/worksites/{wsid:[0-9]+}", withUserManager("GetWorkSite", route.GetWorkSite)).Methods("GET")
	router.HandleFunc("/api/worksites/{wsid:[0-9]+}/attach", withUserManager("GetWorkSiteAttachement", route.GetWorkSiteAttachement)).Methods("GET")
	router.HandleFunc("/api/worksites/{wsid:[0-9]+}/zip", withUserManager("GetWorkSiteDOEArchive", route.GetWorkSiteDOEArchive)).Methods("GET")
	router.HandleFunc("/api/worksites/{wsid:[0-9]+}", withUserManager("UpdateWorkSite", route.UpdateWorkSite)).Methods("PUT")
	router.HandleFunc("/api/worksites/{wsid:[0-9]+}", withUserManager("DeleteWorkSite", route.DeleteWorkSite)).Methods("DELETE")

	// Ripsites methods
	router.HandleFunc("/api/ripsites", withUserManager("GetRipsitesInfo", route.GetRipsitesInfo)).Methods("GET")
	//router.HandleFunc("/api/ripsites", withUserManager("CreateRipSite", route.CreateRipSite)).Methods("POST")
	router.HandleFunc("/api/ripsites/archive", withUserManager("GetRipsitesArchive", route.GetRipsitesArchive)).Methods("GET")
	router.HandleFunc("/api/ripsites/stat/{groupby}/{freq}", withUserManager("GetRipsitesStats", route.GetRipsitesStats)).Methods("GET")
	router.HandleFunc("/api/ripsites/{rsid:[0-9]+}", withUserManager("GetRipSite", route.GetRipSite)).Methods("GET")
	router.HandleFunc("/api/ripsites/{rsid:[0-9]+}/attach", withUserManager("GetRipSiteAttachement", route.GetRipSiteAttachement)).Methods("GET")
	router.HandleFunc("/api/ripsites/{rsid:[0-9]+}", withUserManager("UpdateRipSite", route.UpdateRipSite)).Methods("PUT")
	router.HandleFunc("/api/ripsites/{rsid:[0-9]+}", withUserManager("DeleteRipSite", route.DeleteRipSite)).Methods("DELETE")
	router.HandleFunc("/api/ripsites/measurement", withUserManager("MeasurementRipSite", route.MeasurementRipSite)).Methods("POST")
	router.HandleFunc("/api/ripsites/actors/{freq}", withUserManager("GetRipsitesActorsActivity", route.GetRipsitesActorsActivity)).Methods("GET")

	// Polesites methods
	router.HandleFunc("/api/polesites", withUserManager("GetPolesitesInfo", route.GetPolesitesInfo)).Methods("GET")
	router.HandleFunc("/api/polesites/archive", withUserManager("GetPolesitesArchive", route.GetPolesitesArchive)).Methods("GET")
	router.HandleFunc("/api/polesites/stat/{freq}", withUserManager("GetPolesitesStats", route.GetPolesitesStats)).Methods("GET")
	router.HandleFunc("/api/polesites/{psid:[0-9]+}", withUserManager("GetPolesite", route.GetPolesite)).Methods("GET")
	router.HandleFunc("/api/polesites/{psid:[0-9]+}/attach", withUserManager("GetPolesiteAttachement", route.GetPolesiteAttachement)).Methods("GET")
	router.HandleFunc("/api/polesites/{psid:[0-9]+}/dictzip", withUserManager("DictZip", route.DictZip)).Methods("GET")
	router.HandleFunc("/api/polesites/{psid:[0-9]+}/export", withUserManager("GetPolesiteExport", route.GetPolesiteExport)).Methods("GET")
	router.HandleFunc("/api/polesites/{psid:[0-9]+}", withUserManager("UpdatePolesite", route.UpdatePolesite)).Methods("PUT")

	// Foasites methods
	router.HandleFunc("/api/foasites", withUserManager("GetFoaSitesInfo", route.GetFoaSitesInfo)).Methods("GET")
	router.HandleFunc("/api/foasites/archive", withUserManager("GetFoaSitesArchive", route.GetFoaSitesArchive)).Methods("GET")
	router.HandleFunc("/api/foasites/stat/{freq}", withUserManager("GetFoaSitesStats", route.GetFoaSitesStats)).Methods("GET")
	router.HandleFunc("/api/foasites/{fsid:[0-9]+}", withUserManager("GetFoaSite", route.GetFoaSite)).Methods("GET")
	router.HandleFunc("/api/foasites/{fsid:[0-9]+}/attach", withUserManager("GetFoaSiteAttachement", route.GetFoaSiteAttachement)).Methods("GET")
	router.HandleFunc("/api/foasites/{fsid:[0-9]+}", withUserManager("UpdateFoaSite", route.UpdateFoaSite)).Methods("PUT")

	// Clients methods
	router.HandleFunc("/api/clients/archive", withUserManager("GetClientsArchive", route.GetClientsArchive)).Methods("GET")

	// Actors methods
	router.HandleFunc("/api/actors", withUserManager("GetActors", route.GetActors)).Methods("GET")
	router.HandleFunc("/api/actors", withUserManager("UpdateActors", route.UpdateActors)).Methods("PUT")
	router.HandleFunc("/api/actors/archive", withUserManager("GetActorsArchive", route.GetActorsArchive)).Methods("GET")
	router.HandleFunc("/api/actors/whrecord/{month:[0-9]{4}-[0-9]{2}-[0-9]{2}}", withUserManager("GetActorsWorkingHoursRecord", route.GetActorsWorkingHoursRecord)).Methods("GET")

	// TimeSheets methods
	router.HandleFunc("/api/timesheet/{week:[0-9]{4}-[0-9]{2}-[0-9]{2}}", withUserManager("GetTimeSheet", route.GetTimeSheet)).Methods("GET")
	router.HandleFunc("/api/timesheet/{week:[0-9]{4}-[0-9]{2}-[0-9]{2}}", withUserManager("UpdateTimeSheet", route.UpdateTimeSheet)).Methods("PUT")
	router.HandleFunc("/api/timesheet/archive", withUserManager("GetTimeSheetsArchive", route.GetTimeSheetsArchive)).Methods("GET")

	// Administration methods
	router.HandleFunc("/api/admin/reload", withUserManager("ReloadPersister", route.ReloadPersister)).Methods("GET")

	// Static Files serving
	router.PathPrefix(conf.AssetsRoot).Handler(http.StripPrefix(conf.AssetsRoot, http.FileServer(http.Dir(conf.AssetsDir))))
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(conf.RootDir)))

	gzipedrouter := handlers.CompressHandler(router)
	//gzipedrouter := router

	LaunchPageInBrowser(conf)
	logger.Entry("Server").LogInfo("listening on " + conf.ServicePort)
	logger.Entry("Server").LogInfo("============================== SERVER READY ====================================")
	logger.Entry("Server").Fatal(http.ListenAndServe(conf.ServicePort, gzipedrouter))
}

func LaunchPageInBrowser(c *Conf) error {
	if !c.LaunchWebBrowser {
		return nil
	}
	cmd := exec.Command("cmd", "/c", "start", "http://localhost"+c.ServicePort)
	return cmd.Start()
}
