package route

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/lpuig/ewin/doe/website/backend/logger"
	mgr "github.com/lpuig/ewin/doe/website/backend/manager"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
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

	err = mgr.GetTimeSheet(w, weekDate)
	if err != nil {
		AddError(w, logmsg, err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.Response = http.StatusOK
}
