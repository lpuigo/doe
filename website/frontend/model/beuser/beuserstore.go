package beuser

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/model/ref"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements/message"
	"github.com/lpuig/ewin/doe/website/frontend/tools/json"
	"honnef.co/go/js/xhr"
)

type Store struct {
	*js.Object

	Users []*BeUser `js:"Users"`

	Ref *ref.Ref `js:"Ref"`
}

func NewStore() *Store {
	bus := &Store{Object: tools.O()}
	bus.Users = []*BeUser{}
	bus.Ref = ref.NewRef(func() string {
		return json.Stringify(bus.Users)
	})
	return bus
}

// Functional Methods

// AddNewUser sets the given user a new negative ID, and adds it to the receiver Store
func (bus *Store) AddNewUser(user *BeUser) {
	nextNewUserId := -1
	if len(bus.Users) > 1 && bus.Users[0].Id <= 0 {
		nextNewUserId = bus.Users[0].Id - 1
	}
	user.Id = nextNewUserId
	bus.Users = append([]*BeUser{user}, bus.Users...)
}

// API Methods

func (bus *Store) CallGetUsers(vm *hvue.VM, onSuccess func()) {
	go bus.callGetUsers(vm, onSuccess)
}

func (bus *Store) callGetUsers(vm *hvue.VM, onSuccess func()) {
	req := xhr.NewRequest("GET", "/api/users")
	req.Timeout = tools.LongTimeOut
	req.ResponseType = xhr.JSON

	err := req.Send(nil)
	if err != nil {
		message.ErrorStr(vm, "Oups! "+err.Error(), true)
		return
	}
	if req.Status != tools.HttpOK {
		message.ErrorRequestMessage(vm, req)
		return
	}
	loadedUsers := []*BeUser{}
	req.Response.Call("forEach", func(item *js.Object) {
		beuser := BeUserFromJS(item)
		beuser.SortClients()
		loadedUsers = append(loadedUsers, beuser)
	})
	bus.Users = loadedUsers
	bus.Ref.SetReference()
	onSuccess()
}

func (bus *Store) CallUpdateUsers(vm *hvue.VM, onSuccess func()) {
	go bus.callUpdateGroups(vm, onSuccess)
}

func (bus *Store) callUpdateGroups(vm *hvue.VM, onSuccess func()) {
	updatedUsers := bus.getUpdatedUsers()
	if len(updatedUsers) == 0 {
		message.ErrorStr(vm, "Could not find any updated user", false)
		return
	}

	req := xhr.NewRequest("PUT", "/api/users")
	req.Timeout = tools.TimeOut
	req.ResponseType = xhr.JSON
	err := req.Send(json.Stringify(updatedUsers))
	if err != nil {
		message.ErrorStr(vm, "Oups! "+err.Error(), true)
		return
	}
	if req.Status != tools.HttpOK {
		message.ErrorRequestMessage(vm, req)
		return
	}
	message.NotifySuccess(vm, "Utilisateurs", "Modifications sauvegardées")
	bus.Ref.SetReference()
	onSuccess()
}

func (bus *Store) getUpdatedUsers() []*BeUser {
	updUsers := []*BeUser{}
	refUsers := bus.GetReferenceUsers()
	userById := map[int]*BeUser{}
	for _, usr := range refUsers {
		userById[usr.Id] = usr
	}
	for _, usr := range bus.Users {
		refUsr := userById[usr.Id]
		if !(refUsr != nil && json.Stringify(usr) == json.Stringify(refUsr)) {
			updUsers = append(updUsers, usr)
		}
	}
	return updUsers
}

func (gs *Store) GetReferenceUsers() []*BeUser {
	refUsers := []*BeUser{}
	json.Parse(gs.Ref.Reference).Call("forEach", func(item *js.Object) {
		grp := BeUserFromJS(item)
		refUsers = append(refUsers, grp)
	})
	return refUsers
}
