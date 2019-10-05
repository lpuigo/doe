package route

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/lpuig/ewin/doe/website/backend/logger"
	mgr "github.com/lpuig/ewin/doe/website/backend/manager"
	"github.com/lpuig/ewin/doe/website/backend/model/actors"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	"net/http"
)

func GetActors(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("GetActors").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	w.Header().Set("Content-Type", "application/json")

	err := mgr.GetActors(w)
	if err != nil {
		AddError(w, logmsg, err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.Response = http.StatusOK
}

func UpdateActors(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	logmsg := logger.TimedEntry("Route").AddRequest("UpdateActors").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	w.Header().Set("Content-Type", "application/json")
	if r.Body == nil {
		AddError(w, logmsg, "request Actors missing", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	updatedActors := []*actors.Actor{}
	err := json.NewDecoder(r.Body).Decode(&updatedActors)
	if err != nil {
		AddError(w, logmsg, fmt.Sprintf("misformatted request body :%v", err.Error()), http.StatusBadRequest)
		return
	}

	err = mgr.Actors.UpdateActors(updatedActors)
	if err != nil {
		AddError(w, logmsg, fmt.Sprintf("error updating actors:%v", err.Error()), http.StatusInternalServerError)
		return
	}

	logmsg.Response = http.StatusOK
}

func GetActorsArchive(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("GetActorsArchive").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", mgr.ActorsArchiveName()))
	w.Header().Set("Content-Type", "application/zip")

	err := mgr.CreateActorsArchive(w)
	if err != nil {
		AddError(w, logmsg, err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.Response = http.StatusOK
}

func GetActorsWorkingHoursRecord(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("GetActorsWorkingHoursRecord").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	monthDate := mux.Vars(r)["month"]
	_, err := date.ParseDate(monthDate)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		AddError(w, logmsg, fmt.Sprintf("misformated date '%s'", monthDate), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", mgr.GetActorsWorkingHoursRecordXLSName(monthDate)))
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	err = mgr.GetActorsWorkingHoursRecordXLS(w, monthDate)
	if err != nil {
		AddError(w, logmsg, err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.AddInfoResponse(fmt.Sprintf("Actors Working Hours Record for %s produced", monthDate), http.StatusOK)
}
