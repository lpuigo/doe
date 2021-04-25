package main

import (
	"crypto/tls"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/lpuig/ewin/doe/website/backend/config"
	"github.com/lpuig/ewin/doe/website/backend/logger"
	"github.com/lpuig/ewin/doe/website/backend/manager"
	"github.com/lpuig/ewin/doe/website/backend/route"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
	_ "net/http/pprof"
	"os/exec"
	"sync"
	"time"
)

type Conf struct {
	manager.ManagerConfig

	LogFile     string
	ServiceHost string
	ServicePort string
	AssetsDir   string
	AssetsRoot  string
	RootDir     string

	LaunchWebBrowser    bool
	InProduction        bool
	RedirectHTTPToHTTPS bool
}

const (
	AssetsDir  = `../../WebAssets`
	AssetsRoot = `/Assets/`
	RootDir    = `./Dist`

	ServiceHost = "vps642354.ovh.net"
	ServicePort = ":8080"
	SessionKey  = "SECRET_KEY"

	WorksitesDir   = `C:\Users\Laurent\Golang\src\github.com\lpuig\ewin\doe\Ressources\Worksites`
	RipsitesDir    = `C:\Users\Laurent\Golang\src\github.com\lpuig\ewin\doe\Ressources\Ripsites`
	PolesitesDir   = `C:\Users\Laurent\Golang\src\github.com\lpuig\ewin\doe\Ressources\Polesites`
	FoasitesDir    = `C:\Users\Laurent\Golang\src\github.com\lpuig\ewin\doe\Ressources\Foasites`
	UsersDir       = `C:\Users\Laurent\Golang\src\github.com\lpuig\ewin\doe\Ressources\Users`
	ActorsDir      = `C:\Users\Laurent\Golang\src\github.com\lpuig\ewin\doe\Ressources\Actors`
	ActorInfosDir  = `C:\Users\Laurent\Golang\src\github.com\lpuig\ewin\doe\Ressources\Actorinfos`
	TimeSheetsDir  = `C:\Users\Laurent\Golang\src\github.com\lpuig\ewin\doe\Ressources\Timesheets`
	CalendarFile   = `C:\Users\Laurent\Golang\src\github.com\lpuig\ewin\doe\Ressources\Calendar\holidays.json`
	ClientsDir     = `C:\Users\Laurent\Golang\src\github.com\lpuig\ewin\doe\Ressources\Clients`
	GroupsDir      = `C:\Users\Laurent\Golang\src\github.com\lpuig\ewin\doe\Ressources\Groups`
	VehiculesDir   = `C:\Users\Laurent\Golang\src\github.com\lpuig\ewin\doe\Ressources\Vehicules`
	TemplatesDir   = `C:\Users\Laurent\Golang\src\github.com\lpuig\ewin\doe\Ressources\DocTemplates`
	SaveArchiveDir = `C:\Users\Laurent\Golang\src\github.com\lpuig\ewin\doe\SaveArchive`

	LaunchWebBrowser    = true
	InProduction        = false
	RedirectHTTPToHTTPS = false

	ConfigFile = `./config.json`
	LogFile    = `./server.log`
)

func LaunchPageInBrowser(c *Conf) error {
	if !c.LaunchWebBrowser {
		return nil
	}
	cmd := exec.Command("cmd", "/c", "start", "http://localhost"+c.ServicePort)
	return cmd.Start()
}

// createEmptyRouter sets a default router
func createEmptyRouter(conf *Conf) http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello Secure World")
	})
	return router
}

// createRouter sets a router with all functional route  using given configuration
func createRouter(conf *Conf) http.Handler {
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
	// attach pprof route from defaultServeMux
	router.PathPrefix("/debug/pprof/").Handler(http.DefaultServeMux)
	// session management
	router.HandleFunc("/api/login", withManager(route.GetUser)).Methods("GET")
	router.HandleFunc("/api/login", withUserManager("Logout", route.Logout)).Methods("DELETE")
	router.HandleFunc("/api/login", withManager(route.Login)).Methods("POST")

	// Worksites methods
	router.HandleFunc("/api/worksites", withUserManager("GetWorksitesInfo", route.GetWorksitesInfo)).Methods("GET")
	router.HandleFunc("/api/worksites", withUserManager("CreateWorkSite", route.CreateWorkSite)).Methods("POST")
	router.HandleFunc("/api/worksites/stat/{info}/{freq}", withUserManager("GetWorksitesStats", route.GetWorksitesStats)).Methods("GET")
	router.HandleFunc("/api/worksites/{wsid:[0-9]+}", withUserManager("GetWorkSite", route.GetWorkSite)).Methods("GET")
	router.HandleFunc("/api/worksites/{wsid:[0-9]+}/attach", withUserManager("GetWorkSiteAttachement", route.GetWorkSiteAttachement)).Methods("GET")
	router.HandleFunc("/api/worksites/{wsid:[0-9]+}/zip", withUserManager("GetWorkSiteDOEArchive", route.GetWorkSiteDOEArchive)).Methods("GET")
	router.HandleFunc("/api/worksites/{wsid:[0-9]+}", withUserManager("UpdateWorkSite", route.UpdateWorkSite)).Methods("PUT")
	router.HandleFunc("/api/worksites/{wsid:[0-9]+}", withUserManager("DeleteWorkSite", route.DeleteWorkSite)).Methods("DELETE")

	// Ripsites methods
	router.HandleFunc("/api/ripsites", withUserManager("GetRipsitesInfo", route.GetRipsitesInfo)).Methods("GET")
	//router.HandleFunc("/api/ripsites", withUserManager("CreateRipSite", route.CreateRipSite)).Methods("POST")
	router.HandleFunc("/api/ripsites/stat/{groupby}/{freq}", withUserManager("GetRipsitesStats", route.GetRipsitesStats)).Methods("GET")
	router.HandleFunc("/api/ripsites/progress/{month}", withUserManager("GetRipsitesProgress", route.GetRipsitesProgress)).Methods("GET")
	router.HandleFunc("/api/ripsites/{rsid:[0-9]+}", withUserManager("GetRipSite", route.GetRipSite)).Methods("GET")
	router.HandleFunc("/api/ripsites/{rsid:[0-9]+}", withUserManager("UpdateRipSite", route.UpdateRipSite)).Methods("PUT")
	router.HandleFunc("/api/ripsites/{rsid:[0-9]+}", withUserManager("DeleteRipSite", route.DeleteRipSite)).Methods("DELETE")
	router.HandleFunc("/api/ripsites/measurement", withUserManager("MeasurementRipSite", route.MeasurementRipSite)).Methods("POST")
	router.HandleFunc("/api/ripsites/actors/{freq}", withUserManager("GetRipsitesActorsActivity", route.GetRipsitesActorsActivity)).Methods("GET")

	// Polesites methods
	router.HandleFunc("/api/polesites", withUserManager("GetPolesitesInfo", route.GetPolesitesInfo)).Methods("GET")
	router.HandleFunc("/api/polesites/stat/{groupby}/{freq}", withUserManager("GetPolesitesStats", route.GetPolesitesStats)).Methods("GET")
	router.HandleFunc("/api/polesites/progress/{month}", withUserManager("GetPolesitesProgress", route.GetPolesitesProgress)).Methods("GET")
	router.HandleFunc("/api/polesites/{psid:[0-9]+}", withUserManager("GetPolesite", route.GetPolesite)).Methods("GET")
	router.HandleFunc("/api/polesites/{psid:[0-9]+}", withUserManager("UpdatePolesite", route.UpdatePolesite)).Methods("PUT")
	router.HandleFunc("/api/polesites/{psid:[0-9]+}/archivecompleted", withUserManager("ArchiveCompletedPoleRefs", route.ArchiveCompletedPoleRefs)).Methods("GET")
	router.HandleFunc("/api/polesites/{psid:[0-9]+}/dictzip", withUserManager("DictZip", route.DictZip)).Methods("GET")
	router.HandleFunc("/api/polesites/{psid:[0-9]+}/export", withUserManager("GetPolesiteExport", route.GetPolesiteExport)).Methods("GET")
	router.HandleFunc("/api/polesites/{psid:[0-9]+}/progress", withUserManager("GetPolesiteProgress", route.GetPolesiteProgress)).Methods("GET")
	router.HandleFunc("/api/polesites/{psid:[0-9]+}/planning", withUserManager("GetPolesitePlanning", route.GetPolesitePlanning)).Methods("GET")
	router.HandleFunc("/api/polesites/{psid:[0-9]+}/refexport", withUserManager("GetPolesiteRefExport", route.GetPolesiteRefExport)).Methods("GET")
	router.HandleFunc("/api/polesites/{psid:[0-9]+}/kizeo", withUserManager("GetPolesiteRefKizeo", route.GetPolesiteRefKizeo)).Methods("POST")
	router.HandleFunc("/api/polesites/{psid:[0-9]+}/import", withUserManager("PostPolesiteImport", route.PostPolesiteImport)).Methods("POST")

	// Foasites methods
	router.HandleFunc("/api/foasites", withUserManager("GetFoaSitesInfo", route.GetFoaSitesInfo)).Methods("GET")
	router.HandleFunc("/api/foasites/stat/{freq}", withUserManager("GetFoaSitesStats", route.GetFoaSitesStats)).Methods("GET")
	router.HandleFunc("/api/foasites/{fsid:[0-9]+}", withUserManager("GetFoaSite", route.GetFoaSite)).Methods("GET")
	router.HandleFunc("/api/foasites/{fsid:[0-9]+}", withUserManager("UpdateFoaSite", route.UpdateFoaSite)).Methods("PUT")

	// Archives methods
	router.HandleFunc("/api/{recordtype}/archive", withUserManager("GetRecordsArchive", route.GetRecordsArchive)).Methods("GET")
	router.HandleFunc("/api/archive", withUserManager("GetSaveArchive", route.GetSaveArchive)).Methods("GET")

	// Attachements methods
	router.HandleFunc("/api/{sitetype}/{id:[0-9]+}/attach", withUserManager("GetSiteAttachement", route.GetItemizableSiteAttachement)).Methods("GET")

	// Groups methods
	router.HandleFunc("/api/groups", withUserManager("GetGroups", route.GetGroups)).Methods("GET")
	router.HandleFunc("/api/groups", withUserManager("UpdateGroups", route.UpdateGroups)).Methods("PUT")

	// Vehicules methods
	router.HandleFunc("/api/vehicules", withUserManager("GetVehicules", route.GetVehicules)).Methods("GET")
	router.HandleFunc("/api/vehicules", withUserManager("UpdateVehicules", route.UpdateVehicules)).Methods("PUT")

	// Clients methods
	router.HandleFunc("/api/clients", withUserManager("GetClients", route.GetClients)).Methods("GET")
	router.HandleFunc("/api/clients", withUserManager("UpdateClients", route.UpdateClients)).Methods("PUT")

	// Users methods
	router.HandleFunc("/api/users", withUserManager("GetUsers", route.GetUsers)).Methods("GET")
	router.HandleFunc("/api/users", withUserManager("UpdateUsers", route.UpdateUsers)).Methods("PUT")

	// Actors methods
	router.HandleFunc("/api/actors", withUserManager("GetActors", route.GetActors)).Methods("GET")
	router.HandleFunc("/api/actors", withUserManager("UpdateActors", route.UpdateActors)).Methods("PUT")
	router.HandleFunc("/api/actors/whrecord/{month:[0-9]{4}-[0-9]{2}-[0-9]{2}}", withUserManager("GetActorsWorkingHoursRecord", route.GetActorsWorkingHoursRecord)).Methods("GET")

	//// ActorInfos methods
	//router.HandleFunc("/api/actorinfos", withUserManager("GetActorInfos", route.GetActorInfos)).Methods("GET")
	//router.HandleFunc("/api/actorinfos", withUserManager("UpdateActorInfos", route.UpdateActorInfos)).Methods("PUT")

	// TimeSheets methods
	router.HandleFunc("/api/timesheet/{week:[0-9]{4}-[0-9]{2}-[0-9]{2}}", withUserManager("GetTimeSheet", route.GetTimeSheet)).Methods("GET")
	router.HandleFunc("/api/timesheet/{week:[0-9]{4}-[0-9]{2}-[0-9]{2}}", withUserManager("UpdateTimeSheet", route.UpdateTimeSheet)).Methods("PUT")

	// Administration methods
	router.HandleFunc("/api/admin/reload", withUserManager("ReloadPersister", route.ReloadPersister)).Methods("GET")

	// Static Files serving
	router.PathPrefix(conf.AssetsRoot).Handler(http.StripPrefix(conf.AssetsRoot, http.FileServer(http.Dir(conf.AssetsDir))))
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(conf.RootDir)))

	gzipedrouter := handlers.CompressHandler(router)

	return gzipedrouter
}

func makeServerFromMux(mux http.Handler) *http.Server {
	// set timeouts so that a slow or malicious client doesn't
	// hold resources forever
	return &http.Server{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}
}

func main() {
	conf := &Conf{
		ManagerConfig: manager.ManagerConfig{
			WorksitesDir:      WorksitesDir,
			IsWorksitesActive: true,
			RipsitesDir:       RipsitesDir,
			IsRipsitesActive:  true,
			PolesitesDir:      PolesitesDir,
			IsPolesitesActive: true,
			FoasitesDir:       FoasitesDir,
			IsFoasitesActive:  true,
			UsersDir:          UsersDir,
			ActorsDir:         ActorsDir,
			ActorInfosDir:     ActorInfosDir,
			TimeSheetsDir:     TimeSheetsDir,
			CalendarFile:      CalendarFile,
			ClientsDir:        ClientsDir,
			GroupsDir:         GroupsDir,
			VehiculesDir:      VehiculesDir,
			TemplatesDir:      TemplatesDir,
			SessionKey:        SessionKey,
			SaveArchiveDir:    SaveArchiveDir,
		},
		LogFile:             LogFile,
		ServiceHost:         ServiceHost,
		ServicePort:         ServicePort,
		AssetsDir:           AssetsDir,
		AssetsRoot:          AssetsRoot,
		RootDir:             RootDir,
		LaunchWebBrowser:    LaunchWebBrowser,
		InProduction:        InProduction,
		RedirectHTTPToHTTPS: RedirectHTTPToHTTPS,
	}

	if err := config.SetFromFile(ConfigFile, conf); err != nil {
		logger.Entry("Server").Fatal(err)
	}

	logFile := logger.StartLog(conf.LogFile)
	defer logFile.Close()
	logger.Entry("Server").LogInfo("============================= SERVER STARTING ==================================")

	router := createRouter(conf)

	wg := sync.WaitGroup{}

	if conf.InProduction {
		logger.Entry("Server").LogInfo("Init Production setup")
		//hostPolicy := func(ctx context.Context, host string) error {
		//	if host == conf.ServiceHost {
		//		return nil
		//	}
		//	err := fmt.Errorf("acme/autocert: only %s host is allowed", conf.ServiceHost)
		//	logger.Entry("Server").LogError(err.Error())
		//	return err
		//}

		dataDir := "."
		certManager := &autocert.Manager{
			Prompt: autocert.AcceptTOS,
			//HostPolicy: hostPolicy,
			Cache: autocert.DirCache(dataDir),
		}

		httpsSrv := makeServerFromMux(router)
		httpsSrv.Addr = ":443"
		httpsSrv.TLSConfig = &tls.Config{GetCertificate: certManager.GetCertificate}

		wg.Add(2)
		go func() {
			logger.Entry("Server").LogInfo("listening HTTPS on " + httpsSrv.Addr)
			logger.Entry("Server").LogErr(httpsSrv.ListenAndServeTLS("", ""))
			wg.Done()
		}()
		go func() {
			logger.Entry("Server").LogInfo("listening HTTP on :80 for certification handshake")
			logger.Entry("Server").LogErr(http.ListenAndServe(":80", certManager.HTTPHandler(nil)))
			wg.Done()
		}()
	} else {
		logger.Entry("Server").LogInfo("Init Non Production setup")
		httpSrv := makeServerFromMux(router)
		httpSrv.Addr = conf.ServicePort
		wg.Add(1)
		go func() {
			logger.Entry("Server").LogInfo("listening HTTP on " + httpSrv.Addr)
			logger.Entry("Server").LogErr(httpSrv.ListenAndServe())
			wg.Done()
		}()
	}

	logger.Entry("Server").LogInfo("============================== SERVER READY ====================================")
	LaunchPageInBrowser(conf)
	wg.Wait()
}
