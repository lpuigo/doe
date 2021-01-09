package beclient

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/model/ref"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements/message"
	"github.com/lpuig/ewin/doe/website/frontend/tools/json"
	"honnef.co/go/js/xhr"
	"sort"
)

type Store struct {
	*js.Object

	Clients []*BeClient `js:"Clients"`
	Ref     *ref.Ref    `js:"Ref"`
}

func NewStore() *Store {
	bcs := &Store{Object: tools.O()}
	bcs.Clients = []*BeClient{}
	bcs.Ref = ref.NewRef(func() string {
		return json.Stringify(bcs.Clients)
	})
	return bcs
}

// Functional Methods

// AddNewClient sets the given BeClient a new negative ID, and adds it to the receiver Store
func (bcs *Store) AddNewClient(client *BeClient) {
	nextNewClientId := -1
	if len(bcs.Clients) > 1 && bcs.Clients[0].Id <= 0 {
		nextNewClientId = bcs.Clients[0].Id - 1
	}
	client.Id = nextNewClientId
	bcs.Clients = append([]*BeClient{client}, bcs.Clients...)
}

func (bcs *Store) GetReferenceClients() []*BeClient {
	refClients := []*BeClient{}
	json.Parse(bcs.Ref.Reference).Call("forEach", func(item *js.Object) {
		grp := BeClientFromJS(item)
		refClients = append(refClients, grp)
	})
	return refClients
}

func (bcs *Store) GetClientById(id int) *BeClient {
	for _, beClient := range bcs.Clients {
		if beClient.Id == id {
			return beClient
		}
	}
	return nil
}

// API Methods

func (bcs *Store) CallGetClients(vm *hvue.VM, onSuccess func()) {
	go bcs.callGetClients(vm, onSuccess)
}

func (bcs *Store) callGetClients(vm *hvue.VM, onSuccess func()) {
	req := xhr.NewRequest("GET", "/api/clients")
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
	loadedClients := []*BeClient{}
	req.Response.Call("forEach", func(item *js.Object) {
		bclt := BeClientFromJS(item)
		loadedClients = append(loadedClients, bclt)
	})
	sort.Slice(loadedClients, func(i, j int) bool {
		return loadedClients[i].Name < loadedClients[j].Name
	})
	bcs.Clients = loadedClients
	bcs.Ref.SetReference()
	onSuccess()
}

func (bcs *Store) CallUpdateClients(vm *hvue.VM, onSuccess func()) {
	go bcs.callUpdateClients(vm, onSuccess)
}

func (bcs *Store) callUpdateClients(vm *hvue.VM, onSuccess func()) {
	updatedClients := bcs.getUpdatedClients()
	if len(updatedClients) == 0 {
		message.ErrorStr(vm, "Could not find any updated client", false)
		return
	}

	req := xhr.NewRequest("PUT", "/api/clients")
	req.Timeout = tools.TimeOut
	req.ResponseType = xhr.JSON
	err := req.Send(json.Stringify(updatedClients))
	if err != nil {
		message.ErrorStr(vm, "Oups! "+err.Error(), true)
		return
	}
	if req.Status != tools.HttpOK {
		message.ErrorRequestMessage(vm, req)
		return
	}
	message.NotifySuccess(vm, "Clients", "Modifications sauvegardées")
	bcs.Ref.SetReference()
	onSuccess()
}

func (bcs *Store) getUpdatedClients() []*BeClient {
	updClients := []*BeClient{}
	refClients := bcs.GetReferenceClients()
	ClientById := map[int]*BeClient{}
	for _, usr := range refClients {
		ClientById[usr.Id] = usr
	}
	for _, clt := range bcs.Clients {
		refClt := ClientById[clt.Id]
		if !(refClt != nil && json.Stringify(clt) == json.Stringify(refClt)) {
			updClients = append(updClients, clt)
		}
	}
	return updClients
}
