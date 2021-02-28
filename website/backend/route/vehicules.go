package route

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lpuig/ewin/doe/website/backend/logger"
	mgr "github.com/lpuig/ewin/doe/website/backend/manager"
	"github.com/lpuig/ewin/doe/website/backend/model/vehicules"
)

func GetVehicules(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("GetVehicules").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	w.Header().Set("Content-Type", "application/json")

	err := mgr.GetVehicules(w)
	if err != nil {
		AddError(w, logmsg, err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.Response = http.StatusOK
}

func UpdateVehicules(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	logmsg := logger.TimedEntry("Route").AddRequest("UpdateVehicules").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	w.Header().Set("Content-Type", "application/json")
	if r.Body == nil {
		AddError(w, logmsg, "request Vehicules missing", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	updatedVehicules := []*vehicules.Vehicule{}
	err := json.NewDecoder(r.Body).Decode(&updatedVehicules)
	if err != nil {
		AddError(w, logmsg, fmt.Sprintf("misformatted request body :%v", err.Error()), http.StatusBadRequest)
		return
	}

	err = mgr.UpdateVehicules(updatedVehicules)
	if err != nil {
		AddError(w, logmsg, fmt.Sprintf("error updating vehicules:%v", err.Error()), http.StatusInternalServerError)
		return
	}

	logmsg.Response = http.StatusOK
}
