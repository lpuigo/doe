package route

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/lpuig/ewin/doe/website/backend/logger"
	mgr "github.com/lpuig/ewin/doe/website/backend/manager"
	"github.com/lpuig/ewin/doe/website/backend/model/worksites"
	"net/http"
	"strconv"
)

func GetWorkSites(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("GetWorkSites")
	defer logmsg.Log()

	w.Header().Set("Content-Type", "application/json")
	err := mgr.GetWorkSites(w)
	if err != nil {
		addError(w, logmsg, err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.Response = http.StatusOK
}

func GetWorkSite(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("GetWorkSite")
	defer logmsg.Log()

	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	wsrid, err := strconv.Atoi(vars["wsid"])
	if err != nil {
		addError(w, logmsg, "mis-formatted WorkSite id '"+vars["wsid"]+"'", http.StatusBadRequest)
		return
	}
	wsr := mgr.Worksites.GetById(wsrid)
	if wsr == nil {
		addError(w, logmsg, fmt.Sprintf("workSite with id %d does not exist", wsrid), http.StatusNotFound)
		return
	}
	err = wsr.Marshall(w)
	if err != nil {
		addError(w, logmsg, "could not marshall WorkSite. "+err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.AddInfoResponse(fmt.Sprintf("workSite Id %d returned", wsr.Id), http.StatusOK)
}

func CreateWorkSite(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("CreateWorkSite")
	defer logmsg.Log()

	if r.Body == nil {
		addError(w, logmsg, "request WorkSite missing", http.StatusBadRequest)
		return
	}
	wsr, err := worksites.NewWorkSiteRecordFrom(r.Body)
	if err != nil {
		addError(w, logmsg, "malformed WorkSite: "+err.Error(), http.StatusBadRequest)
		return
	}
	wsr = mgr.Worksites.Add(wsr)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = wsr.Marshall(w)
	if err != nil {
		addError(w, logmsg, "could not marshall WorkSite. "+err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.AddInfoResponse(fmt.Sprintf("New WorkSite Id %d added", wsr.Id), http.StatusCreated)
}

func UpdateWorkSite(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("UpdateWorkSite")
	defer logmsg.Log()

	if r.Body == nil {
		addError(w, logmsg, "request WorkSite missing", http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	wsrid, err := strconv.Atoi(vars["wsid"])
	if err != nil {
		addError(w, logmsg, "mis-formatted WorkSite id '"+vars["wsid"]+"'", http.StatusBadRequest)
		return
	}
	wsr := mgr.Worksites.GetById(wsrid)
	if wsr == nil {
		addError(w, logmsg, fmt.Sprintf("workSite with id %d does not exist", wsrid), http.StatusNotFound)
		return
	}
	wsr, err = worksites.NewWorkSiteRecordFrom(r.Body)
	if err != nil {
		addError(w, logmsg, "malformed WorkSite: "+err.Error(), http.StatusBadRequest)
		return
	}
	if wsr.Id != wsrid {
		addError(w, logmsg, fmt.Sprintf("inconsitent WorkSite id between request (%d) and body (%d)", wsrid, wsr.Id), http.StatusBadRequest)
		return
	}
	err = mgr.Worksites.Update(wsr)
	if err != nil {
		addError(w, logmsg, fmt.Sprintf("could not update WorkSite with id %d: %v", wsrid, err), http.StatusInternalServerError)
		return
	}
	logmsg.AddInfoResponse(fmt.Sprintf("WorkSite with id %d updated", wsrid), http.StatusOK)
}

func DeleteWorkSite(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("DeleteWorkSite")
	defer logmsg.Log()

	vars := mux.Vars(r)
	wsrid, err := strconv.Atoi(vars["wsid"])
	if err != nil {
		addError(w, logmsg, "mis-formatted WorkSite id '"+vars["wsid"]+"'", http.StatusBadRequest)
		return
	}
	wsr := mgr.Worksites.GetById(wsrid)
	if wsr == nil {
		addError(w, logmsg, fmt.Sprintf("workSite with id %d does not exist", wsrid), http.StatusNoContent)
		return
	}
	err = mgr.Worksites.Remove(wsr)
	if err != nil {
		addError(w, logmsg, fmt.Sprintf("could not delete WorkSite with id %d: %v", wsrid, err), http.StatusInternalServerError)
		return
	}
	logmsg.AddInfoResponse(fmt.Sprintf("WorkSite with id %d deleted", wsrid), http.StatusOK)
}
