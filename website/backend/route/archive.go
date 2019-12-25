package route

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/lpuig/ewin/doe/website/backend/logger"
	mgr "github.com/lpuig/ewin/doe/website/backend/manager"
	"github.com/lpuig/ewin/doe/website/backend/model/archives"
	"net/http"
)

func GetRecordsArchive(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("GetRecordsArchive").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	vars := mux.Vars(r)
	recordType := vars["recordtype"]
	var sites archives.ArchivableRecordContainer
	switch recordType {
	case "worksites":
		sites = mgr.Worksites
	case "ripsites":
		sites = mgr.Ripsites
	case "polesites":
		sites = mgr.Polesites
	case "foasites":
		sites = mgr.Foasites
	case "clients":
		sites = mgr.Clients
	case "actors":
		sites = mgr.Actors
	case "timesheet":
		sites = mgr.TimeSheets
	default:
		AddError(w, logmsg, "unsupported archive type '"+recordType+"'", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", archives.ArchiveName(sites)))
	w.Header().Set("Content-Type", "application/zip")

	err := archives.CreateRecordsArchive(w, sites)
	if err != nil {
		AddError(w, logmsg, err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.AddInfoResponse(fmt.Sprintf("%s archive produced", recordType), http.StatusOK)
}
