package route

import (
	"fmt"
	"github.com/lpuig/ewin/doe/website/backend/logger"
	mgr "github.com/lpuig/ewin/doe/website/backend/manager"
	"net/http"
	"time"
)

type MgrHandlerFunc func(*mgr.Manager, http.ResponseWriter, *http.Request)

func addError(w http.ResponseWriter, errmsg string, code int) string {
	res := logger.LogResponse(code)
	res += logger.LogInfo(errmsg)
	http.Error(w, errmsg, code)
	return res
}

func GetWorkSites(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.LogRequest("GetWorkSites")
	defer logger.LogService(time.Now(), &logmsg)

	w.Header().Set("Content-Type", "application/json")
	err := mgr.GetWorkSites(w)
	if err != nil {
		logmsg += addError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg += logger.LogResponse(http.StatusOK)
}

func CreatePrj(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.LogRequest("CreatePrj")
	defer logger.LogService(time.Now(), &logmsg)

	//var prj = &fm.Project{}
	//if r.Body == nil {
	//	logmsg += addError(w, "request project missing", http.StatusBadRequest)
	//	return
	//}
	//err := json.NewDecoder(r.Body).Decode(prj)
	//if err != nil {
	//	logmsg += addError(w, "unable to retrieve request project. "+err.Error(), http.StatusBadRequest)
	//	return
	//}
	//ptfPrj, hasStat := mgr.CreateProject(fm.CloneFEProject(prj))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	//json.NewEncoder(w).Encode(fm.CloneBEProject(ptfPrj, hasStat))
	logmsg += logger.LogInfo(fmt.Sprintf("New project Id %d added (%d)", 10, http.StatusCreated))
}
