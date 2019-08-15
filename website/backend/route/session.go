package route

import (
	"encoding/json"
	"github.com/lpuig/ewin/doe/website/backend/logger"
	mgr "github.com/lpuig/ewin/doe/website/backend/manager"
	"github.com/lpuig/ewin/doe/website/backend/model/actors"
	"github.com/lpuig/ewin/doe/website/backend/model/clients"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	"net/http"
	"sort"
)

// Facade structs dedicated to expose User & Client info to FrontEnd
type authentActor struct {
	Id        int
	LastName  string
	FirstName string
	Role      string
	Active    bool
}

func authentActorsFrom(acs []*actors.Actor) []authentActor {
	today := date.Today().String()
	res := make([]authentActor, len(acs))
	for i, actor := range acs {
		res[i] = authentActor{
			Id:        actor.Id,
			LastName:  actor.LastName,
			FirstName: actor.FirstName,
			Role:      actor.Role,
			Active:    actor.IsActiveOn(today),
		}
	}
	sort.Slice(res, func(i, j int) bool {
		si := res[i].Role + " " + res[i].LastName + " " + res[i].FirstName
		sj := res[j].Role + " " + res[j].LastName + " " + res[j].FirstName
		return si < sj
	})
	return res
}

type authentClient struct {
	Name     string
	Teams    []clients.Team
	Actors   []authentActor
	Articles []string
}

func getAuthentClientFrom(mgr *mgr.Manager, clients []*clients.Client) []authentClient {
	res := []authentClient{}
	for _, client := range clients {
		actors := mgr.Actors.GetActorsByClient(client.Name, false)
		authentActors := authentActorsFrom(actors)
		authClient := authentClient{
			Name:     client.Name,
			Teams:    client.Teams,
			Actors:   authentActors,
			Articles: client.GetOrangeArticleNames(),
		}
		res = append(res, authClient)
	}
	return res
}

type authentUser struct {
	Name        string
	Clients     []authentClient
	Permissions map[string]bool
}

func newAuthentUser() authentUser {
	return authentUser{
		Name:        "",
		Clients:     []authentClient{},
		Permissions: make(map[string]bool),
	}
}

func (au *authentUser) SetFrom(mgr *mgr.Manager, clients []*clients.Client) {
	au.Name = mgr.CurrentUser.Name
	au.Clients = getAuthentClientFrom(mgr, clients)
	au.Permissions = mgr.CurrentUser.Permissions
}

// GetUser checks for session cookie, and return pertaining user
//
// If user is not authenticated, the session is removed
func GetUser(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("GetUser")
	defer logmsg.Log()

	w.Header().Set("Content-Type", "application/json")

	user := newAuthentUser()
	// check for session cookie
	if !mgr.CheckSessionUser(r) {
		// user cookie not found or improper, remove it first
		err := mgr.SessionStore.RemoveSessionCookie(w, r)
		if err != nil {
			AddError(w, logmsg, "could not remove session info", http.StatusInternalServerError)
			return
		}
		AddError(w, logmsg, "user not authorized", http.StatusUnauthorized)
		return
		// Todo Exit
	}

	// found a correct one, set user
	logmsg.AddUser(mgr.CurrentUser.Name)
	clts, err := mgr.GetCurrentUserClients()
	if err != nil {
		AddError(w, logmsg, err.Error(), http.StatusInternalServerError)
		return
	}
	user.SetFrom(mgr, clts)

	// refresh session cookie
	err = mgr.SessionStore.RefreshSessionCookie(w, r)
	if err != nil {
		AddError(w, logmsg, "could not refresh session cookie", http.StatusInternalServerError)
		return
	}

	// write response
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		AddError(w, logmsg, "could not encode authent user", http.StatusInternalServerError)
		return
	}
	logmsg.Info = "authenticated"
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

func Logout(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("Logout").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	err := mgr.SessionStore.RemoveSessionCookie(w, r)
	if err != nil {
		AddError(w, logmsg, err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.Info = "logged out"
	logmsg.Response = http.StatusOK
}
