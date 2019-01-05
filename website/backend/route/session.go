package route

import (
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

	getValue := func(key string) (string, bool) {
		info, found := r.MultipartForm.Value[key]
		if !found {
			return "", false
		}
		return info[0], true
	}

	user, hasUser := getValue("user")
	pwd, hasPwd := getValue("pwd")
	if !(hasUser && hasPwd) {
		AddError(w, logmsg, "user/password info missing", http.StatusBadRequest)
		return
	}
	//TODO Improve Login/pwd and authorization here
	u := mgr.Users.GetByName(user)
	if u == nil {
		AddError(w, logmsg, "user not authorized", http.StatusUnauthorized)
		return
	}
	if u.Password != pwd {
		AddError(w, logmsg, "user/password not authorized", http.StatusUnauthorized)
		return
	}

	if err := mgr.SessionStore.AddSessionCookie(u, w, r); err != nil {
		AddError(w, logmsg, err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.Response = http.StatusOK
	logmsg.Info = "user '" + user + "' connected with password '" + pwd + "'"
}
