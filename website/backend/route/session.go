package route

import (
	"encoding/json"
	"github.com/lpuig/ewin/doe/website/backend/logger"
	mgr "github.com/lpuig/ewin/doe/website/backend/manager"
	"github.com/lpuig/ewin/doe/website/backend/model/users"
	"net/http"
)

type authentUser struct {
	Name string
}

func (au *authentUser) SetFrom(ur *users.UserRecord) {
	au.Name = ur.Name
}

func GetUser(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("GetUser")
	defer logmsg.Log()

	w.Header().Set("Content-Type", "application/json")

	user := authentUser{Name: ""}
	// check for session cookie
	if mgr.CheckSessionUser(r) {
		// found a correct one, set user
		user.SetFrom(mgr.CurrentUser)
		logmsg.AddUser(user.Name)
		logmsg.Info = "authenticated"
	} else {
		// not found or improper, remove it first
		err := mgr.SessionStore.RemoveSessionCookie(w, r)
		if err != nil {
			AddError(w, logmsg, "could not remove session info", http.StatusInternalServerError)
			return
		}
		logmsg.Info = "not authenticated"
	}
	err := json.NewEncoder(w).Encode(user)
	if err != nil {
		AddError(w, logmsg, "could not encode authent user", http.StatusInternalServerError)
		return
	}
	logmsg.Response = http.StatusOK
}

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
	logmsg.AddUser(user)

	//TODO Improve Login/pwd and authorization here
	u := mgr.Users.GetByName(user)
	if u == nil {
		AddError(w, logmsg, "user not authorized", http.StatusUnauthorized)
		return
	}
	if u.Password != pwd {
		AddError(w, logmsg, "wrong password", http.StatusUnauthorized)
		return
	}

	if err := mgr.SessionStore.AddSessionCookie(u, w, r); err != nil {
		AddError(w, logmsg, err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.Response = http.StatusOK
	logmsg.Info = "logged in"
}
