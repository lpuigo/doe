package route

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/lpuig/ewin/doe/website/backend/logger"
	mgr "github.com/lpuig/ewin/doe/website/backend/manager"
	"github.com/lpuig/ewin/doe/website/backend/model/worksites"
	"net/http"
	"strconv"
	"time"
)

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

func GetWorkSite(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.LogRequest("GetWorkSite")
	defer logger.LogService(time.Now(), &logmsg)

	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	wsrid, err := strconv.Atoi(vars["wsid"])
	if err != nil {
		logmsg += addError(w, "mis-formatted WorkSite id '"+vars["wsid"]+"'", http.StatusBadRequest)
		return
	}
	wsr := mgr.Worksites.GetById(wsrid)
	if wsr == nil {
		logmsg += addError(w, fmt.Sprintf("workSite with id %d does not exist", wsrid), http.StatusNotFound)
		return
	}
	err = wsr.Marshall(w)
	if err != nil {
		logmsg += addError(w, "could not marshall WorkSite. "+err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg += logger.LogResponseInfo(fmt.Sprintf("workSite Id %d returned", wsr.Id), http.StatusOK)
}

func CreateWorkSite(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.LogRequest("CreateWorkSite")
	defer logger.LogService(time.Now(), &logmsg)

	if r.Body == nil {
		logmsg += addError(w, "request WorkSite missing", http.StatusBadRequest)
		return
	}
	wsr, err := worksites.NewWorkSiteRecordFrom(r.Body)
	if err != nil {
		logmsg += addError(w, "malformed WorkSite: "+err.Error(), http.StatusBadRequest)
		return
	}
	wsr = mgr.Worksites.Add(wsr)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = wsr.Marshall(w)
	if err != nil {
		logmsg += addError(w, "could not marshall WorkSite. "+err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg += logger.LogResponseInfo(fmt.Sprintf("New WorkSite Id %d added", wsr.Id), http.StatusCreated)
}

func UpdateWorkSite(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.LogRequest("UpdateWorkSite")
	defer logger.LogService(time.Now(), &logmsg)

	if r.Body == nil {
		logmsg += addError(w, "request WorkSite missing", http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	wsrid, err := strconv.Atoi(vars["wsid"])
	if err != nil {
		logmsg += addError(w, "mis-formatted WorkSite id '"+vars["wsid"]+"'", http.StatusBadRequest)
		return
	}
	wsr := mgr.Worksites.GetById(wsrid)
	if wsr == nil {
		logmsg += addError(w, fmt.Sprintf("workSite with id %d does not exist", wsrid), http.StatusNotFound)
		return
	}
	wsr, err = worksites.NewWorkSiteRecordFrom(r.Body)
	if err != nil {
		logmsg += addError(w, "malformed WorkSite: "+err.Error(), http.StatusBadRequest)
		return
	}
	if wsr.Id != wsrid {
		logmsg += addError(w, fmt.Sprintf("inconsitent WorkSite id between request (%d) and body (%d)", wsrid, wsr.Id), http.StatusBadRequest)
		return
	}
	err = mgr.Worksites.Update(wsr)
	if err != nil {
		logmsg += addError(w, fmt.Sprintf("could not update WorkSite with id %d: %v", wsrid, err), http.StatusInternalServerError)
		return
	}
	logmsg += logger.LogResponseInfo(fmt.Sprintf("WorkSite with id %d updated", wsrid), http.StatusOK)
}

func DeleteWorkSite(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.LogRequest("DeleteWorkSite")
	defer logger.LogService(time.Now(), &logmsg)

	vars := mux.Vars(r)
	wsrid, err := strconv.Atoi(vars["wsid"])
	if err != nil {
		logmsg += addError(w, "mis-formatted WorkSite id '"+vars["wsid"]+"'", http.StatusBadRequest)
		return
	}
	wsr := mgr.Worksites.GetById(wsrid)
	if wsr == nil {
		logmsg += addError(w, fmt.Sprintf("workSite with id %d does not exist", wsrid), http.StatusNoContent)
		return
	}
	err = mgr.Worksites.Remove(wsr)
	if err != nil {
		logmsg += addError(w, fmt.Sprintf("could not delete WorkSite with id %d: %v", wsrid, err), http.StatusInternalServerError)
		return
	}
	logmsg += logger.LogResponseInfo(fmt.Sprintf("WorkSite with id %d deleted", wsrid), http.StatusOK)
}
