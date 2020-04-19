package route

import (
	"encoding/json"
	"fmt"
	"github.com/lpuig/ewin/doe/website/backend/logger"
	mgr "github.com/lpuig/ewin/doe/website/backend/manager"
	"github.com/lpuig/ewin/doe/website/backend/model/actorinfos"
	"net/http"
)

func GetActorInfos(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("GetActorInfos").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	w.Header().Set("Content-Type", "application/json")

	err := mgr.GetActorInfos(w)
	if err != nil {
		AddError(w, logmsg, err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.Response = http.StatusOK
}

func UpdateActorInfos(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	logmsg := logger.TimedEntry("Route").AddRequest("UpdateActorInfos").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	w.Header().Set("Content-Type", "application/json")
	if r.Body == nil {
		AddError(w, logmsg, "request Actorinfos missing", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	updatedActorInfos := []*actorinfos.ActorInfo{}
	err := json.NewDecoder(r.Body).Decode(&updatedActorInfos)
	if err != nil {
		AddError(w, logmsg, fmt.Sprintf("misformatted request body :%v", err.Error()), http.StatusBadRequest)
		return
	}

	err = mgr.ActorInfos.UpdateActorInfos(updatedActorInfos)
	if err != nil {
		AddError(w, logmsg, fmt.Sprintf("error updating actorinfos:%v", err.Error()), http.StatusInternalServerError)
		return
	}

	logmsg.Response = http.StatusOK
}
