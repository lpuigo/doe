package route

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lpuig/ewin/doe/website/backend/logger"
	mgr "github.com/lpuig/ewin/doe/website/backend/manager"
	"github.com/lpuig/ewin/doe/website/backend/model/polesites"
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

func UpdatePolesite(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("UpdatePolesite").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	if r.Body == nil {
		AddError(w, logmsg, "request Polesite missing", http.StatusBadRequest)
		return
	}
	reqPolesiteId := mux.Vars(r)["psid"]
	rpsid, err := strconv.Atoi(reqPolesiteId)
	if err != nil {
		AddError(w, logmsg, "mis-formatted Polesite id '"+reqPolesiteId+"'", http.StatusBadRequest)
		return
	}
	psr := mgr.Polesites.GetById(rpsid)
	if psr == nil {
		AddError(w, logmsg, fmt.Sprintf("Polesite with id %d does not exist", rpsid), http.StatusNotFound)
		return
	}
	psr, err = polesites.NewPoleSiteRecordFrom(r.Body)
	if err != nil {
		AddError(w, logmsg, "malformed Polesite: "+err.Error(), http.StatusBadRequest)
		return
	}
	if psr.Id != rpsid {
		AddError(w, logmsg, fmt.Sprintf("inconsitent Polesite id between request (%d) and body (%d)", rpsid, psr.Id), http.StatusBadRequest)
		return
	}
	err = mgr.Polesites.Update(psr)
	if err != nil {
		AddError(w, logmsg, fmt.Sprintf("could not update Polesite with id %d: %v", rpsid, err), http.StatusInternalServerError)
		return
	}
	logmsg.AddInfoResponse(fmt.Sprintf("Polesite with id %d (%s) updated", rpsid, psr.Ref), http.StatusOK)

}

func GetPolesitesStats(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("GetPolesitesStats").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	w.Header().Set("Content-Type", "application/json")

	var err error

	vars := mux.Vars(r)
	freq := vars["freq"]
	switch freq {
	case "week":
		err = mgr.GetPolesitesWeekStats(w)
	case "month":
		err = mgr.GetPolesitesMonthStats(w)
	default:
		AddError(w, logmsg, "unsupported stat type '"+freq+"'", http.StatusBadRequest)
		return
	}

	if err != nil {
		AddError(w, logmsg, err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.AddInfoResponse(fmt.Sprintf("%s polesite stats produced", freq), http.StatusOK)
}

func GetPolesitesArchive(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("GetPolesitesArchive").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", mgr.PolesitesArchiveName()))
	w.Header().Set("Content-Type", "application/zip")

	err := mgr.CreatePolesitesArchive(w)
	if err != nil {
		AddError(w, logmsg, err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.Response = http.StatusOK
}

// GetPolesiteAttachement
func GetPolesiteAttachement(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("GetPolesiteAttachement").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	reqPolesiteId := mux.Vars(r)["psid"]
	rpsid, err := strconv.Atoi(reqPolesiteId)
	if err != nil {
		AddError(w, logmsg, "mis-formatted Polesite id '"+reqPolesiteId+"'", http.StatusBadRequest)
		return
	}
	psr := mgr.Polesites.GetById(rpsid)
	if psr == nil {
		AddError(w, logmsg, fmt.Sprintf("Polesite with id %d does not exist", rpsid), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", mgr.TemplateEngine.GetPolesiteXLSAttachementName(psr.PoleSite)))
	w.Header().Set("Content-Type", "vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	err = mgr.GetPolesiteXLSAttachement(w, psr.PoleSite)
	if err != nil {
		AddError(w, logmsg, "could not generate Polesite XLS Attachment file. "+err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.AddInfoResponse(fmt.Sprintf("Attachment XLS produced for Polesite id %d (%s)", rpsid, psr.PoleSite.Ref), http.StatusOK)
}
