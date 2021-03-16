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
	var container archives.ArchivableRecordContainer
	switch recordType {
	case "worksites":
		container = mgr.Worksites
	case "ripsites":
		container = mgr.Ripsites
	case "polesites":
		container = mgr.Polesites
	case "foasites":
		container = mgr.Foasites
	case "clients":
		container = mgr.Clients
	case "actors":
		container = mgr.Actors
	case "actorinfos":
		container = mgr.ActorInfos
	case "timesheet":
		container = mgr.TimeSheets
	case "groups":
		container = mgr.Groups
	case "vehicules":
		container = mgr.Vehicules
	case "users":
		container = mgr.Users
	default:
		AddError(w, logmsg, "unsupported archive type '"+recordType+"'", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", archives.ArchiveName(container)))
	w.Header().Set("Content-Type", "application/zip")

	err := archives.CreateRecordsArchive(w, container)
	if err != nil {
		AddError(w, logmsg, err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.AddInfoResponse(fmt.Sprintf("%s archive produced", recordType), http.StatusOK)
}
