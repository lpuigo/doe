package route

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/lpuig/ewin/doe/website/backend/logger"
	mgr "github.com/lpuig/ewin/doe/website/backend/manager"
	"github.com/lpuig/ewin/doe/website/backend/model/foasites"
	"net/http"
	"strconv"
)

func GetFoaSitesInfo(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("GetFoaSitesInfo").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	w.Header().Set("Content-Type", "application/json")

	err := mgr.GetFoaSitesInfo(w)
	if err != nil {
		AddError(w, logmsg, err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.Response = http.StatusOK
}

func GetFoaSitesStats(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("GetFoaSitesStats").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	w.Header().Set("Content-Type", "application/json")

	var err error

	vars := mux.Vars(r)
	freq := vars["freq"]
	switch freq {
	case "day", "week", "month":
		err = mgr.GetFoaSitesStats(w, freq)
	default:
		AddError(w, logmsg, "unsupported stat type '"+freq+"'", http.StatusBadRequest)
		return
	}

	if err != nil {
		AddError(w, logmsg, err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.AddInfoResponse(fmt.Sprintf("%s foasite stats produced", freq), http.StatusOK)
}

func GetFoaSite(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("GetFoaSite").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	fsid, err := strconv.Atoi(vars["fsid"])
	if err != nil {
		AddError(w, logmsg, "mis-formatted foasite id '"+vars["fsid"]+"'", http.StatusBadRequest)
		return
	}
	psr := mgr.Foasites.GetById(fsid)
	if psr == nil {
		AddError(w, logmsg, fmt.Sprintf("foasite with id %d does not exist", fsid), http.StatusNotFound)
		return
	}
	err = psr.Marshall(w)
	if err != nil {
		AddError(w, logmsg, "could not marshall foasite. "+err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.AddInfoResponse(fmt.Sprintf("foasite Id %d (%s) returned", psr.Id, psr.Ref), http.StatusOK)
}

func UpdateFoaSite(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	logmsg := logger.TimedEntry("Route").AddRequest("UpdateFoaSite").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	if r.Body == nil {
		AddError(w, logmsg, "request FoaSite missing", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	reqFoaSiteId := mux.Vars(r)["fsid"]
	rfsid, err := strconv.Atoi(reqFoaSiteId)
	if err != nil {
		AddError(w, logmsg, "mis-formatted FoaSite id '"+reqFoaSiteId+"'", http.StatusBadRequest)
		return
	}
	fsr := mgr.Foasites.GetById(rfsid)
	if fsr == nil {
		AddError(w, logmsg, fmt.Sprintf("FoaSite with id %d does not exist", rfsid), http.StatusNotFound)
		return
	}
	fsr, err = foasites.NewFoaSiteRecordFrom(r.Body)
	if err != nil {
		AddError(w, logmsg, "malformed FoaSite: "+err.Error(), http.StatusBadRequest)
		return
	}
	if fsr.Id != rfsid {
		AddError(w, logmsg, fmt.Sprintf("inconsitent FoaSite id between request (%d) and body (%d)", rfsid, fsr.Id), http.StatusBadRequest)
		return
	}
	err = mgr.Foasites.Update(fsr)
	if err != nil {
		AddError(w, logmsg, fmt.Sprintf("could not update FoaSite with id %d: %v", rfsid, err), http.StatusInternalServerError)
		return
	}
	logmsg.AddInfoResponse(fmt.Sprintf("FoaSite with id %d (%s) updated", rfsid, fsr.Ref), http.StatusOK)

}
