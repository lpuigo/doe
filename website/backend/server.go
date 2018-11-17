package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/lpuig/ewin/doe/website/backend/config"
	"github.com/lpuig/ewin/doe/website/backend/logger"
	"github.com/lpuig/ewin/doe/website/backend/manager"
	"github.com/lpuig/ewin/doe/website/backend/route"
	"log"
	"net/http"
	"os/exec"
)

type Conf struct {
	WorksitesDir string

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

	WorksitesDir = `C:\Users\Laurent\Golang\src\github.com\lpuig\ewin\doe\Ressources\Worksites`

	LaunchWebBrowser = true

	ConfigFile = `./config.json`
	LogFile    = `./server.log`
)

func main() {
	conf := &Conf{
		WorksitesDir:     WorksitesDir,
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
		log.Fatal(err)
	}
	log.Println("Server Started =============================================================================")

	mgr, err := manager.NewManager(conf.WorksitesDir)
	if err != nil {
		log.Fatal(err)
	}
	withManager := func(hf route.MgrHandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			hf(mgr, w, r)
		}
	}

	router := mux.NewRouter()
	// Manager Routes
	//TODO Remove
	router.HandleFunc("/ptf", withManager(route.GetWorkSites)).Methods("GET")
	router.HandleFunc("/ptf", withManager(route.CreatePrj)).Methods("POST")

	// Static Files serving
	router.PathPrefix(conf.AssetsRoot).Handler(http.StripPrefix(conf.AssetsRoot, http.FileServer(http.Dir(conf.AssetsDir))))
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(conf.RootDir)))

	gzipedrouter := handlers.CompressHandler(router)
	//gzipedrouter := router

	LaunchPageInBrowser(conf)
	log.Print("Listening on ", ServicePort)
	log.Fatal(http.ListenAndServe(ServicePort, gzipedrouter))
}

func LaunchPageInBrowser(c *Conf) error {
	if !c.LaunchWebBrowser {
		return nil
	}
	cmd := exec.Command("cmd", "/c", "start", "http://localhost"+c.ServicePort)
	return cmd.Start()
}
