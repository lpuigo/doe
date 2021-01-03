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
	Assigned  bool
}

// authentActorsFrom returns a []authentActor populated with all Actor in given visibleActors map.
//
// authentActor.Assigned is set to true if given Actor is also in given assignedActors map
//
// Result Slice is sorted by Actors Properties : Active (true first) then Assigned (true first) then full name
func authentActorsFrom(mgr *mgr.Manager, visibleActors, assignedActors map[int]*actors.Actor) []authentActor {
	today := date.Today().String()
	res := make([]authentActor, len(visibleActors))
	i := 0
	for _, actor := range visibleActors {
		active := actor.IsActiveOn(today)
		res[i] = authentActor{
			Id:        actor.Id,
			LastName:  actor.LastName,
			FirstName: actor.FirstName,
			Role:      actor.Role,
			Active:    active,
			Assigned:  active && assignedActors[actor.Id] != nil,
		}
		i++
	}
	sort.Slice(res, func(i, j int) bool {
		if res[i].Active != res[j].Active {
			return res[i].Active
		}
		if res[i].Assigned != res[j].Assigned {
			return res[i].Assigned
		}
		si := res[i].LastName + " " + res[i].FirstName
		sj := res[j].LastName + " " + res[j].FirstName
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

// getAuthentClientFrom returns a slice of authentClient based on given clients list
//
// authentClient are populated with visible actors, based on client visibility via actor's groups
func getAuthentClientFrom(mgr *mgr.Manager, clients []*clients.Client) []authentClient {
	res := []authentClient{}
	getActorsByClientName := mgr.GenActorsByClientName()
	getAssignedActorsByGroupId := mgr.GenCurrentlyAssignedActorsByGroupId()
	clientActors := make(map[string]map[int]*actors.Actor)         // client's eligible actors
	clientAssignedActors := make(map[string]map[int]*actors.Actor) // client's assigned actors

	// actors visible by client name
	for _, client := range clients {
		clientActors[client.Name] = getActorsByClientName(client.Name)
	}

	// actors assigned by client name
	for _, group := range mgr.GetCurrentUserVisibleGroups() {
		groupAssignActors := getAssignedActorsByGroupId(group.Id)
		for _, clientName := range group.Clients {
			assignedActs := clientAssignedActors[clientName]
			if assignedActs == nil {
				assignedActs = make(map[int]*actors.Actor)
			}
			for _, actor := range groupAssignActors {
				assignedActs[actor.Id] = actor
			}
			clientAssignedActors[clientName] = assignedActs
		}
	}

	// create authentClient list
	for _, client := range clients {
		authentActors := authentActorsFrom(mgr, clientActors[client.Name], clientAssignedActors[client.Name])
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
	DaysOff     map[string]string
}

func newAuthentUser() authentUser {
	return authentUser{
		Name:        "",
		Clients:     []authentClient{},
		Permissions: make(map[string]bool),
		DaysOff:     make(map[string]string),
	}
}

func (au *authentUser) SetFrom(mgr *mgr.Manager, clients []*clients.Client) {
	au.Name = mgr.CurrentUser.Name
	au.Clients = getAuthentClientFrom(mgr, clients)
	au.Permissions = mgr.CurrentUser.Permissions
	au.DaysOff = mgr.DaysOff.GetDays()
}

// GetUser checks for session cookie, and returns pertaining user
//
// If user is not authenticated, the session is removed
func GetUser(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("GetUser")
	defer logmsg.Log()

	w.Header().Set("Content-Type", "application/json")

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
	}

	// found a correct one, set user
	logmsg.AddUser(mgr.CurrentUser.Name)
	clts, err := mgr.GetCurrentUserClients()
	if err != nil {
		AddError(w, logmsg, err.Error(), http.StatusInternalServerError)
		return
	}
	user := newAuthentUser()
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
