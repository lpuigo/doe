package actor

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

type ActorStore struct {
	*js.Object

	Actors []*Actor `js:"Actors"`

	Ref *ref.Ref `js:"Ref"`
}

func NewActorStore() *ActorStore {
	as := &ActorStore{Object: tools.O()}
	as.Actors = []*Actor{}
	as.Ref = ref.NewRef(func() string {
		return json.Stringify(as.Actors)
	})
	return as
}

func (as *ActorStore) GetActorsSortedByName() []*Actor {
	res := as.Actors[:]
	sort.Slice(res, func(i, j int) bool {
		return res[i].Ref < res[j].Ref
	})
	return res
}

func (as *ActorStore) sortActorsById() {
	sort.Slice(as.Actors, func(i, j int) bool {
		return as.Actors[i].Id < as.Actors[j].Id
	})
}

func (as *ActorStore) CallGetActors(vm *hvue.VM, onSuccess func()) {
	go as.callGetActors(vm, onSuccess)
}

func (as *ActorStore) callGetActors(vm *hvue.VM, onSuccess func()) {
	req := xhr.NewRequest("GET", "/api/actors")
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
	loadedActors := []*Actor{}
	req.Response.Call("forEach", func(item *js.Object) {
		act := ActorFromJS(item)
		act.UpdateState()
		loadedActors = append(loadedActors, act)
	})
	as.Actors = loadedActors
	as.sortActorsById()
	as.Ref.SetReference()
	onSuccess()
}

func (as *ActorStore) GetActorById(id int) *Actor {
	for _, actor := range as.Actors {
		if actor.Id == id {
			return actor
		}
	}
	return nil
}
