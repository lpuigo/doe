package route

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/lpuig/ewin/doe/website/backend/logger"
	mgr "github.com/lpuig/ewin/doe/website/backend/manager"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	"github.com/lpuig/ewin/doe/website/backend/model/timesheets"
	"net/http"
)

func GetTimeSheet(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("GetTimeSheet").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	weekDate := mux.Vars(r)["week"]
	_, err := date.ParseDate(weekDate)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		AddError(w, logmsg, fmt.Sprintf("misformated date '%s'", weekDate), http.StatusBadRequest)
		return
	}

	if weekDate != date.GetMonday(weekDate) {
		w.Header().Set("Content-Type", "application/json")
		AddError(w, logmsg, fmt.Sprintf("date '%s' is not a monday", weekDate), http.StatusBadRequest)
		return
	}

	err = mgr.GetTimeSheet(w, weekDate)
	if err != nil {
		AddError(w, logmsg, err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.Response = http.StatusOK
}

func UpdateTimeSheet(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	logmsg := logger.TimedEntry("Route").AddRequest("UpdateTimeSheet").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	weekDate := mux.Vars(r)["week"]
	_, err := date.ParseDate(weekDate)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		AddError(w, logmsg, fmt.Sprintf("misformated date '%s'", weekDate), http.StatusBadRequest)
		return
	}

	if weekDate != date.GetMonday(weekDate) {
		w.Header().Set("Content-Type", "application/json")
		AddError(w, logmsg, fmt.Sprintf("date '%s' is not a monday", weekDate), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	updatedTimeSheet := &timesheets.TimeSheet{}
	err = json.NewDecoder(r.Body).Decode(updatedTimeSheet)
	if err != nil {
		AddError(w, logmsg, fmt.Sprintf("misformatted request body :%v", err.Error()), http.StatusBadRequest)
		return
	}

	if weekDate != updatedTimeSheet.WeekDate {
		w.Header().Set("Content-Type", "application/json")
		AddError(w, logmsg, fmt.Sprintf("inconsistent request body date '%s'", updatedTimeSheet.WeekDate), http.StatusBadRequest)
		return
	}

	err = mgr.TimeSheets.UpdateTimeSheet(updatedTimeSheet)
	if err != nil {
		AddError(w, logmsg, fmt.Sprintf("error updating TimeSheet:%v", err.Error()), http.StatusInternalServerError)
		return
	}

	logmsg.Response = http.StatusOK
	logmsg.AddInfoResponse(fmt.Sprintf("TimeSheet for %s updated", weekDate), http.StatusOK)
}

func GetTimeSheetsArchive(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("GetTimeSheetsArchive").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", mgr.TimeSheetsArchiveName()))
	w.Header().Set("Content-Type", "application/zip")

	err := mgr.CreateTimeSheetsArchive(w)
	if err != nil {
		AddError(w, logmsg, err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.Response = http.StatusOK
}
