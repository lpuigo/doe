package route

import (
	"fmt"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
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
		AddError(w, logmsg, "mis-formatted Polesite id '"+vars["psid"]+"'", http.StatusBadRequest)
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
	logmsg := logger.TimedEntry("Route").AddRequest("UpdatePolesite").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	if r.Body == nil {
		AddError(w, logmsg, "request Polesite missing", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
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
	case "day", "week", "month":
		err = mgr.GetPolesitesStats(w, freq)
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

func GetPolesitesProgress(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("GetPolesitesProgress").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	w.Header().Set("Content-Type", "application/json")
	var err error
	vars := mux.Vars(r)
	month := vars["month"]
	if date.GetMonth(month) != month {
		AddError(w, logmsg, "misformated date '"+month+"'", http.StatusBadRequest)
		return
	}

	err = mgr.GetPolesitesProgress(w, month)
	if err != nil {
		AddError(w, logmsg, err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.AddInfoResponse(fmt.Sprintf("polesite progress produced for %s", month), http.StatusOK)
}

// GetPolesiteExport
func GetPolesiteExport(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("GetPolesiteExport").AddUser(mgr.CurrentUser.Name)
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

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", psr.PoleSite.ExportName()))
	w.Header().Set("Content-Type", "vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	err = psr.PoleSite.XLSExport(w)
	if err != nil {
		AddError(w, logmsg, "could not generate Polesite XLS export file. "+err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.AddInfoResponse(fmt.Sprintf("Export XLS produced for Polesite id %d (%s)", rpsid, psr.PoleSite.Ref), http.StatusOK)
}

// GetPolesiteRefExport
func GetPolesiteRefExport(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("GetPolesiteRefExport").AddUser(mgr.CurrentUser.Name)
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

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", psr.PoleSite.RefExportName()))
	w.Header().Set("Content-Type", "vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	err = psr.PoleSite.XLSRefExport(w)
	if err != nil {
		AddError(w, logmsg, "could not generate Polesite XLS references export file. "+err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.AddInfoResponse(fmt.Sprintf("Reference Export XLS produced for Polesite id %d (%s)", rpsid, psr.PoleSite.Ref), http.StatusOK)
}

// GetPolesiteProgress
func GetPolesiteProgress(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("GetPolesiteProgress").AddUser(mgr.CurrentUser.Name)
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

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", psr.PoleSite.ProgressName()))
	w.Header().Set("Content-Type", "vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	err = psr.PoleSite.XLSProgress(w)
	if err != nil {
		AddError(w, logmsg, "could not generate Polesite XLS progress file. "+err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.AddInfoResponse(fmt.Sprintf("Progress XLS produced for Polesite id %d (%s)", rpsid, psr.PoleSite.Ref), http.StatusOK)
}

// ArchiveCompletedPoleRefs
func ArchiveCompletedPoleRefs(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("ArchiveCompletedPoleRefs").AddUser(mgr.CurrentUser.Name)
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

	err = mgr.Polesites.ArchiveCompletedPoleRefs(psr)
	if err != nil {
		AddError(w, logmsg, "could not archive completed ref poles: "+err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.AddInfoResponse(fmt.Sprintf("completed for Polesite id %d (%s)", rpsid, psr.PoleSite.Ref), http.StatusOK)
}

// DictZip
func DictZip(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("DictZip").AddUser(mgr.CurrentUser.Name)
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

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", psr.PoleSite.DictZipName()))
	w.Header().Set("Content-Type", "application/zip")

	err = psr.PoleSite.DictZipArchive(w)
	if err != nil {
		AddError(w, logmsg, "could not generate Polesite Dict Zip Archive: "+err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.AddInfoResponse(fmt.Sprintf("Dict Zip Archive produced for Polesite id %d (%s)", rpsid, psr.PoleSite.Ref), http.StatusOK)
}
