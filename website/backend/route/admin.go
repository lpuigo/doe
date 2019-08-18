package route

import (
	"fmt"
	"github.com/lpuig/ewin/doe/website/backend/logger"
	mgr "github.com/lpuig/ewin/doe/website/backend/manager"
	"net/http"
)

func ReloadPersister(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("ReloadPersister").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	w.Header().Set("Content-Type", "application/json")

	if !mgr.CurrentUser.Permissions["Admin"] {
		AddError(w, logmsg, fmt.Sprintf("User '%s' not authorized to do that", mgr.CurrentUser.Name), http.StatusUnauthorized)
		return
	}

	err := mgr.Reload()
	if err != nil {
		AddError(w, logmsg, err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.Response = http.StatusOK
}
