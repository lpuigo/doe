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

	WorksitesDir = `C:\Users\Laurent\Golang\src\github.com\lpuig\ewin\doe\Ressources\Worksites`
	UsersDir     = `C:\Users\Laurent\Golang\src\github.com\lpuig\ewin\doe\Ressources\Users`

	LaunchWebBrowser = true

	ConfigFile = `./config.json`
	LogFile    = `./server.log`
)

func main() {
	conf := &Conf{
		ManagerConfig: manager.ManagerConfig{
			WorksitesDir: WorksitesDir,
			UsersDir:     UsersDir,
			SessionKey:   SessionKey,
		},
		LogFile:          LogFile,
		ServicePort:      ServicePort,
		AssetsDir:        AssetsDir,
		AssetsRoot:       AssetsRoot,
		RootDir:          RootDir,
		LaunchWebBrowser: LaunchWebBrowser,
	}

	logFile := logger.StartLog(conf.LogFile)
	defer logFile.Close()

	if err := config.SetFromFile(ConfigFile, conf); err != nil {
		logger.Entry("Server").Fatal(err)
	}
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
	router.HandleFunc("/api/login", withManager(route.Login)).Methods("POST")
	// Worsite method
	router.HandleFunc("/api/worksites", withUserManager("GetWorkSites", route.GetWorkSites)).Methods("GET")
	router.HandleFunc("/api/worksites", withUserManager("CreateWorkSite", route.CreateWorkSite)).Methods("POST")
	router.HandleFunc("/api/worksites/{wsid:[0-9]+}", withUserManager("GetWorkSite", route.GetWorkSite)).Methods("GET")
	router.HandleFunc("/api/worksites/{wsid:[0-9]+}", withUserManager("UpdateWorkSite", route.UpdateWorkSite)).Methods("PUT")
	router.HandleFunc("/api/worksites/{wsid:[0-9]+}", withUserManager("DeleteWorkSite", route.DeleteWorkSite)).Methods("DELETE")

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
