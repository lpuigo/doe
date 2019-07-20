package route

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/lpuig/ewin/doe/website/backend/logger"
	mgr "github.com/lpuig/ewin/doe/website/backend/manager"
	"net/http"
	"strconv"
)

func GetPolesitesInfo(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("GetPolesitesInfo").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	w.Header().Set("Content-Type", "application/json")

	err := mgr.GetPolesitesInfo(w)
	if err != nil {
		AddError(w, logmsg, err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.Response = http.StatusOK
}

func GetPolesite(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("GetPolesite").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	psid, err := strconv.Atoi(vars["psid"])
	if err != nil {
		AddError(w, logmsg, "mis-formatted Poleite id '"+vars["psid"]+"'", http.StatusBadRequest)
		return
	}
	psr := mgr.Polesites.GetById(psid)
	if psr == nil {
		AddError(w, logmsg, fmt.Sprintf("poleSite with id %d does not exist", psid), http.StatusNotFound)
		return
	}
	err = psr.Marshall(w)
	if err != nil {
		AddError(w, logmsg, "could not marshall Polesite. "+err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.AddInfoResponse(fmt.Sprintf("Polesite Id %d (%s) returned", psr.Id, psr.Ref), http.StatusOK)
}
