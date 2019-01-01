package route

import (
	"fmt"
	"github.com/lpuig/ewin/doe/website/backend/logger"
	mgr "github.com/lpuig/ewin/doe/website/backend/manager"
	"net/http"
)

func Login(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("Login")
	defer logmsg.Log()

	if r.ParseMultipartForm(1024) != nil {
		AddError(w, logmsg, "user info missing", http.StatusBadRequest)
		return
	}

	for key, value := range r.MultipartForm.Value {
		fmt.Printf("%s = %s\n", key, value)
	}
	logmsg.Response = http.StatusOK
	logmsg.Info = "user " + r.MultipartForm.Value["user"][0] + " connected"
}
