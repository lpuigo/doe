package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/actorstable"
	"github.com/lpuig/ewin/doe/website/frontend/comp/actorupdatemodal"
	"github.com/lpuig/ewin/doe/website/frontend/comp/actorvacancyeditmodal"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/model/actor"
	"github.com/lpuig/ewin/doe/website/frontend/model/actor/actorconst"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements/message"
	"github.com/lpuig/ewin/doe/website/frontend/tools/json"
	"honnef.co/go/js/xhr"
)

//go:generate bash ./makejs.sh

func main() {
	mpm := NewMainPageModel()

	hvue.NewVM(
		hvue.El("#actor_app"),
		actorupdatemodal.RegisterComponent(),
		actorvacancyeditmodal.RegisterComponent(),
		actorstable.RegisterComponent(),
		hvue.DataS(mpm),
		hvue.MethodsOf(mpm),
		hvue.Mounted(func(vm *hvue.VM) {
			mpm := &MainPageModel{Object: vm.Object}
			tools.BeforeUnloadConfirmation(mpm.PreventLeave)
			mpm.GetUserSession(func() {
				mpm.LoadActors(false)
			})
		}),
		//hvue.Computed("Title", func(vm *hvue.VM) interface{} {
		//	//mpm := &MainPageModel{Object: vm.Object}
		//	return "To Refactor"
		//}),
		hvue.Computed("IsDirty", func(vm *hvue.VM) interface{} {
			mpm := &MainPageModel{Object: vm.Object}
			mpm.Dirty = (mpm.Reference != json.Stringify(mpm.Actors))
			return mpm.Dirty
		}),
		//hvue.Computed("ShowTable", func(vm *hvue.VM) interface{} {
		//	mpm := &MainPageModel{Object: vm.Object}
		//	if mpm.ActiveMode != "Table" {
		//		return "display: none;"
		//	}
		//	return ""
		//}),
	)

	js.Global.Set("mpm", mpm)
}

type MainPageModel struct {
	*js.Object

	VM   *hvue.VM `js:"VM"`
	User *fm.User `js:"User"`

	ActiveMode string         `js:"ActiveMode"`
	Filter     string         `js:"Filter"`
	FilterType string         `js:"FilterType"`
	Actors     []*actor.Actor `js:"Actors"`
	Reference  string         `js:"Reference"`
	Dirty      bool           `js:"Dirty"`
}

func NewMainPageModel() *MainPageModel {
	mpm := &MainPageModel{Object: tools.O()}
	mpm.VM = nil
	mpm.User = fm.NewUser()
	mpm.ActiveMode = "Table"
	mpm.Filter = ""
	mpm.FilterType = ""
	mpm.Actors = []*actor.Actor{}
	mpm.Reference = ""
	mpm.Dirty = false
	return mpm
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Update & Undo Methods

func (mpm *MainPageModel) PreventLeave() bool {
	return mpm.Dirty
}

func (mpm *MainPageModel) GetReference() string {
	return json.Stringify(mpm.Actors)
}

func (mpm *MainPageModel) SetReference() {
	mpm.Reference = mpm.GetReference()
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// User Management Methods

func (mpm *MainPageModel) GetUserSession(callback func()) {
	go mpm.callGetUser(callback)
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Action Methods

func (mpm *MainPageModel) LoadActors(update bool) {
	updateLoadedActors := func() {
		mpm.SetReference()
		for _, act := range mpm.Actors {
			act.UpdateState()
		}
		// IsDirty is set to true if some update are undertaken
	}
	go mpm.callGetActors(updateLoadedActors)
}

func (mpm *MainPageModel) SaveActors(vm *hvue.VM) {
	mpm = &MainPageModel{Object: vm.Object}

	updateActors := func() {
		mpm.LoadActors(false)
	}

	go mpm.callUpdateActors(updateActors)
}

// SwitchActiveMode handles ActiveMode change
func (mpm *MainPageModel) SwitchActiveMode(vm *hvue.VM) {
	message.ErrorStr(vm, "TODO Implement SwitchActiveMode", false)
}

func (mpm *MainPageModel) GetFilterType(vm *hvue.VM, activeMode string) []*elements.ValueLabel {
	return actor.GetFilterTypeValueLabel()
}

//
func (mpm *MainPageModel) ClearFilter(vm *hvue.VM) {
	mpm = &MainPageModel{Object: vm.Object}
	mpm.Filter = ""
	mpm.FilterType = actorconst.FilterValueAll
	mpm.ApplyFilter(vm)
}

//
func (mpm *MainPageModel) ApplyFilter(vm *hvue.VM) {
	// No OP
}

//
func (mpm *MainPageModel) ShowEditActor(vm *hvue.VM, act *actor.Actor) {
	aem := actorupdatemodal.ActorUpdateModalModelFromJS(mpm.VM.Refs("ActorEditModal"))
	aem.Show(act, mpm.User)
}

//
func (mpm *MainPageModel) ShowEditActorVacancy(vm *hvue.VM, act *actor.Actor) {
	aem := actorvacancyeditmodal.ActorVacancyEditModalModelFromJS(mpm.VM.Refs("ActorVacancyEditModal"))
	aem.Show(act, mpm.User)
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// WS call Methods

func (mpm *MainPageModel) callGetUser(callback func()) {
	req := xhr.NewRequest("GET", "/api/login")
	req.Timeout = tools.LongTimeOut
	req.ResponseType = xhr.JSON
	err := req.Send(nil)
	if err != nil {
		message.ErrorStr(mpm.VM, "Oups! "+err.Error(), true)
		return
	}
	if req.Status != tools.HttpOK {
		message.ErrorRequestMessage(mpm.VM, req)
		return
	}
	mpm.User.Copy(fm.UserFromJS(req.Response))
	if mpm.User.Name == "" {
		mpm.User = fm.NewUser()
		return
	}
	mpm.User.Connected = true
	callback()
}

func (mpm *MainPageModel) callGetActors(callback func()) {
	req := xhr.NewRequest("GET", "/api/actors")
	req.Timeout = tools.LongTimeOut
	req.ResponseType = xhr.JSON

	actors := mpm.Actors[:]
	defer func() {
		mpm.Actors = actors
		callback()
	}()

	err := req.Send(nil)
	if err != nil {
		message.ErrorStr(mpm.VM, "Oups! "+err.Error(), true)
		return
	}
	if req.Status != tools.HttpOK {
		message.ErrorRequestMessage(mpm.VM, req)
		return
	}
	loadedActors := []*actor.Actor{}
	req.Response.Call("forEach", func(item *js.Object) {
		act := actor.NewActorFromJS(item)
		loadedActors = append(loadedActors, act)
	})
	actors = loadedActors
}

func (mpm *MainPageModel) callUpdateActors(callback func()) {
	updatedActors := mpm.getUpdatedActors()
	if len(updatedActors) == 0 {
		message.ErrorStr(mpm.VM, "Could not find any updated actors", false)
		return
	}

	defer callback()

	req := xhr.NewRequest("PUT", "/api/actors")
	req.Timeout = tools.TimeOut
	req.ResponseType = xhr.JSON
	err := req.Send(json.Stringify(updatedActors))
	if err != nil {
		message.ErrorStr(mpm.VM, "Oups! "+err.Error(), true)
		return
	}
	if req.Status != tools.HttpOK {
		message.ErrorRequestMessage(mpm.VM, req)
		return
	}
	message.SuccesStr(mpm.VM, "Modification sauvegardée")
}

func (mpm *MainPageModel) getUpdatedActors() []*actor.Actor {
	refActors := []*actor.Actor{}
	json.Parse(mpm.Reference).Call("forEach", func(item *js.Object) {
		act := actor.NewActorFromJS(item)
		refActors = append(refActors, act)
	})
	refDict := makeDictActors(refActors)
	updDict := makeDictActors(mpm.Actors)

	udpActors := []*actor.Actor{}
	for id, act := range updDict {
		refact := refDict[id]
		if !(refact != nil && json.Stringify(act) == json.Stringify(refDict[id])) {
			print("Changed User", act.Id, act.Ref)
		}
	}
	return udpActors
}

func makeDictActors(actors []*actor.Actor) map[int]*actor.Actor {
	res := make(map[int]*actor.Actor)
	for _, act := range actors {
		res[act.Id] = act
	}
	return res
}
