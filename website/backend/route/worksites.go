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

func GetWorksites(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("GetWorksites").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	w.Header().Set("Content-Type", "application/json")

	//TODO Manage User Authorization (control on mgr.CurrentUser)

	err := mgr.GetWorkSites(w)
	if err != nil {
		AddError(w, logmsg, err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.Response = http.StatusOK
}

func GetWorksitesInfo(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("GetWorksitesInfo").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	w.Header().Set("Content-Type", "application/json")

	//TODO Manage User Authorization (control on mgr.CurrentUser)

	err := mgr.GetWorksitesInfo(w)
	if err != nil {
		AddError(w, logmsg, err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.Response = http.StatusOK
}

func GetWorkSite(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("GetWorkSite").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	wsrid, err := strconv.Atoi(vars["wsid"])
	if err != nil {
		AddError(w, logmsg, "mis-formatted WorkSite id '"+vars["wsid"]+"'", http.StatusBadRequest)
		return
	}
	wsr := mgr.Worksites.GetById(wsrid)
	if wsr == nil {
		AddError(w, logmsg, fmt.Sprintf("workSite with id %d does not exist", wsrid), http.StatusNotFound)
		return
	}
	err = wsr.Marshall(w)
	if err != nil {
		AddError(w, logmsg, "could not marshall WorkSite. "+err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.AddInfoResponse(fmt.Sprintf("workSite Id %d returned", wsr.Id), http.StatusOK)
}

func CreateWorkSite(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("CreateWorkSite").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	if r.Body == nil {
		AddError(w, logmsg, "request WorkSite missing", http.StatusBadRequest)
		return
	}
	wsr, err := worksites.NewWorkSiteRecordFrom(r.Body)
	if err != nil {
		AddError(w, logmsg, "malformed WorkSite: "+err.Error(), http.StatusBadRequest)
		return
	}
	wsr = mgr.Worksites.Add(wsr)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = wsr.Marshall(w)
	if err != nil {
		AddError(w, logmsg, "could not marshall WorkSite. "+err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.AddInfoResponse(fmt.Sprintf("New WorkSite Id %d added", wsr.Id), http.StatusCreated)
}

func UpdateWorkSite(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("UpdateWorkSite").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	if r.Body == nil {
		AddError(w, logmsg, "request WorkSite missing", http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	wsrid, err := strconv.Atoi(vars["wsid"])
	if err != nil {
		AddError(w, logmsg, "mis-formatted WorkSite id '"+vars["wsid"]+"'", http.StatusBadRequest)
		return
	}
	wsr := mgr.Worksites.GetById(wsrid)
	if wsr == nil {
		AddError(w, logmsg, fmt.Sprintf("workSite with id %d does not exist", wsrid), http.StatusNotFound)
		return
	}
	wsr, err = worksites.NewWorkSiteRecordFrom(r.Body)
	if err != nil {
		AddError(w, logmsg, "malformed WorkSite: "+err.Error(), http.StatusBadRequest)
		return
	}
	if wsr.Id != wsrid {
		AddError(w, logmsg, fmt.Sprintf("inconsitent WorkSite id between request (%d) and body (%d)", wsrid, wsr.Id), http.StatusBadRequest)
		return
	}
	err = mgr.Worksites.Update(wsr)
	if err != nil {
		AddError(w, logmsg, fmt.Sprintf("could not update WorkSite with id %d: %v", wsrid, err), http.StatusInternalServerError)
		return
	}
	logmsg.AddInfoResponse(fmt.Sprintf("WorkSite with id %d updated", wsrid), http.StatusOK)
}

func DeleteWorkSite(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("DeleteWorkSite").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	vars := mux.Vars(r)
	wsrid, err := strconv.Atoi(vars["wsid"])
	if err != nil {
		AddError(w, logmsg, "mis-formatted WorkSite id '"+vars["wsid"]+"'", http.StatusBadRequest)
		return
	}
	wsr := mgr.Worksites.GetById(wsrid)
	if wsr == nil {
		AddError(w, logmsg, fmt.Sprintf("workSite with id %d does not exist", wsrid), http.StatusNoContent)
		return
	}
	err = mgr.Worksites.Remove(wsr)
	if err != nil {
		AddError(w, logmsg, fmt.Sprintf("could not delete WorkSite with id %d: %v", wsrid, err), http.StatusInternalServerError)
		return
	}
	logmsg.AddInfoResponse(fmt.Sprintf("WorkSite with id %d deleted", wsrid), http.StatusOK)
}

func GetWorksitesStats(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("GetWorksitesStats").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	w.Header().Set("Content-Type", "application/json")

	//TODO Manage User Authorization (control on mgr.CurrentUser)

	err := mgr.GetWorkSitesStats(w)
	if err != nil {
		AddError(w, logmsg, err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.Response = http.StatusOK

}
