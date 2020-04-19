package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/actorscalendar"
	"github.com/lpuig/ewin/doe/website/frontend/comp/actorsstatsmodal"
	"github.com/lpuig/ewin/doe/website/frontend/comp/actorstable"
	"github.com/lpuig/ewin/doe/website/frontend/comp/actorstimesheet"
	"github.com/lpuig/ewin/doe/website/frontend/comp/actorupdatemodal"
	"github.com/lpuig/ewin/doe/website/frontend/comp/actorvacancyeditmodal"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/model/actor"
	"github.com/lpuig/ewin/doe/website/frontend/model/actor/actorconst"
	"github.com/lpuig/ewin/doe/website/frontend/model/actorinfo"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/date"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements/message"
	"github.com/lpuig/ewin/doe/website/frontend/tools/json"
	"honnef.co/go/js/xhr"
	"strconv"
)

//go:generate bash ./makejs.sh

func main() {
	mpm := NewMainPageModel()

	hvue.NewVM(
		hvue.El("#actor_app"),
		actorupdatemodal.RegisterComponent(),
		actorvacancyeditmodal.RegisterComponent(),
		actorsstatsmodal.RegisterComponent(),
		actorscalendar.RegisterComponent(),
		actorstimesheet.RegisterComponent(),
		actorstable.RegisterComponent(),
		hvue.DataS(mpm),
		hvue.MethodsOf(mpm),
		hvue.Mounted(func(vm *hvue.VM) {
			mpm := &MainPageModel{Object: vm.Object}
			tools.BeforeUnloadConfirmation(mpm.PreventLeave)
			mpm.GetUserSession(func() {
				mpm.LoadActors(true)
			})
		}),
		hvue.Computed("IsDirty", func(vm *hvue.VM) interface{} {
			mpm := &MainPageModel{Object: vm.Object}
			mpm.Dirty = mpm.CheckReference()
			mpm.ActorInfosDirty = mpm.CheckActorInfosReference()
			return mpm.Dirty || mpm.ActorInfosDirty
		}),
	)

	js.Global.Set("mpm", mpm)
}

type MainPageModel struct {
	*js.Object

	VM   *hvue.VM `js:"VM"`
	User *fm.User `js:"User"`

	ActiveMode          string                       `js:"ActiveMode"`
	Filter              string                       `js:"Filter"`
	FilterType          string                       `js:"FilterType"`
	Actors              []*actor.Actor               `js:"Actors"`
	Reference           string                       `js:"Reference"`
	Dirty               bool                         `js:"Dirty"`
	ActorInfos          []*actorinfo.ActorInfo       `js:"ActorInfos"`
	ActorInfosByActorId map[int]*actorinfo.ActorInfo `js:"ActorInfosByActorId"`
	ActorInfosReference string                       `js:"ActorInfosReference"`
	ActorInfosDirty     bool                         `js:"ActorInfosDirty"`

	CraVisible bool   `js:"craVisible"`
	CraMonth   string `js:"craMonth"`
}

func NewMainPageModel() *MainPageModel {
	mpm := &MainPageModel{Object: tools.O()}
	mpm.VM = nil
	mpm.User = fm.NewUser()
	mpm.ActiveMode = "Calendar"
	mpm.Filter = ""
	mpm.FilterType = ""

	mpm.Actors = []*actor.Actor{}
	mpm.Reference = ""
	mpm.Dirty = false

	mpm.ActorInfos = []*actorinfo.ActorInfo{}
	mpm.ActorInfosByActorId = make(map[int]*actorinfo.ActorInfo)
	mpm.ActorInfosReference = ""
	mpm.ActorInfosDirty = false

	mpm.CraVisible = false
	mpm.CraMonth = date.GetFirstOfMonth(date.TodayAfter(-7))

	return mpm
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Update & Undo Methods

func (mpm *MainPageModel) PreventLeave() bool {
	return mpm.Dirty
}

// Ref for Actors

func (mpm *MainPageModel) GetReference() string {
	return json.Stringify(mpm.Actors)
}

func (mpm *MainPageModel) SetReference() {
	mpm.Reference = mpm.GetReference()
}

// CheckReference returns true when some Actors has changed
func (mpm *MainPageModel) CheckReference() bool {
	return mpm.Reference != json.Stringify(mpm.Actors)
}

// Ref for ActorInfos

func (mpm *MainPageModel) GetActorInfosReference() string {
	return json.Stringify(mpm.ActorInfos)
}

func (mpm *MainPageModel) SetActorInfosReference() {
	mpm.ActorInfosReference = mpm.GetActorInfosReference()
}

// CheckActorInfosReference returns true when some ActorInfos has changed
func (mpm *MainPageModel) CheckActorInfosReference() bool {
	return mpm.ActorInfosReference != json.Stringify(mpm.ActorInfos)
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// User Management Methods

func (mpm *MainPageModel) GetUserSession(callback func()) {
	go mpm.callGetUser(callback)
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Action Methods

func (mpm *MainPageModel) AddActor() {
	newActor := actor.NewActor()
	mpm.Actors = append(mpm.Actors, newActor)
	mpm.ShowEditActor(mpm.VM, newActor)
}

func (mpm *MainPageModel) LoadActors(init bool) {
	onLoadedActors := func() {
		mpm.SetReference()
		for _, act := range mpm.Actors {
			act.UpdateState()
		}
		// IsDirty is set to true if some update are undertaken
		if init && mpm.CheckReference() {
			mpm.SaveActors(mpm.VM)
		}
		mpm.LoadActorInfos()
	}
	go mpm.callGetActors(onLoadedActors)
}

func (mpm *MainPageModel) SaveActors(vm *hvue.VM) {
	mpm = &MainPageModel{Object: vm.Object}

	if mpm.Dirty {
		onUpdatedActors := func() {
			//mpm.LoadActors(false)
			mpm.SetReference()
		}
		go mpm.callUpdateActors(onUpdatedActors)
	}

	if mpm.ActorInfosDirty {
		mpm.SaveActorInfos(vm)
	}
}

func (mpm *MainPageModel) LoadActorInfos() {
	if !(mpm.User.Connected && mpm.User.HasPermissionHR()) {
		return
	}

	onLoadedActorInfos := func() {
		mpm.SetActorInfosReference()
		mpm.updateActorInfosByActorId()
	}

	go mpm.callGetActorInfos(onLoadedActorInfos)
}

func (mpm *MainPageModel) updateActorInfosByActorId() {
	actorInfoByActorId := make(map[int]*actorinfo.ActorInfo)
	for _, actInf := range mpm.ActorInfos {
		actorInfoByActorId[actInf.ActorId] = actInf
	}
	mpm.ActorInfosByActorId = actorInfoByActorId
}

func (mpm *MainPageModel) SaveActorInfos(vm *hvue.VM) {
	mpm = &MainPageModel{Object: vm.Object}
	if !(mpm.User.Connected && mpm.User.HasPermissionHR()) {
		return
	}

	onUpdatedActorInfos := func() {
		mpm.SetActorInfosReference()
	}

	go mpm.callUpdateActorInfos(onUpdatedActorInfos)
}

// SwitchActiveMode handles ActiveMode change
func (mpm *MainPageModel) SwitchActiveMode(vm *hvue.VM) {
	//message.ErrorStr(vm, "TODO Implement SwitchActiveMode", false)
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

//
func (mpm *MainPageModel) ShowActorsStats(vm *hvue.VM) {
	aem := actorsstatsmodal.ActorsStatsModalModelFromJS(mpm.VM.Refs("ActorsStatsModal"))
	aem.Show(mpm.Actors, mpm.User)
}

//
func (mpm *MainPageModel) GetActorsWorkingHoursRecord(vm *hvue.VM) {
	js.Global.Get("window").Call("open", "/api/actors/whrecord/"+mpm.CraMonth)
	mpm.CraVisible = false
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

// WS Actors

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
	defer callback()
	if len(updatedActors) == 0 {
		message.ErrorStr(mpm.VM, "Could not find any updated actors", false)
		return
	}

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
	message.NotifySuccess(mpm.VM, "Equipes", "Modifications sauvegardées")
}

func (mpm *MainPageModel) getUpdatedActors() []*actor.Actor {
	refActors := []*actor.Actor{}
	json.Parse(mpm.Reference).Call("forEach", func(item *js.Object) {
		act := actor.NewActorFromJS(item)
		refActors = append(refActors, act)
	})
	refDict := makeDictActors(refActors)

	udpActors := []*actor.Actor{}
	for _, act := range mpm.Actors {
		if act.Ref == "" {
			continue
		}
		refact := refDict[act.Id]
		if !(refact != nil && json.Stringify(act) == json.Stringify(refact)) {
			udpActors = append(udpActors, act)
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

// WS ActorInfos

func (mpm *MainPageModel) callGetActorInfos(callback func()) {
	req := xhr.NewRequest("GET", "/api/actorinfos")
	req.Timeout = tools.LongTimeOut
	req.ResponseType = xhr.JSON

	actorinfos := mpm.ActorInfos
	success := false
	defer func() {
		mpm.ActorInfos = actorinfos
		if success {
			callback()
		}
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
	loadedActorInfos := []*actorinfo.ActorInfo{}
	req.Response.Call("forEach", func(item *js.Object) {
		actinf := actorinfo.NewActorInfoFromJS(item)
		loadedActorInfos = append(loadedActorInfos, actinf)
	})
	success = true
	actorinfos = loadedActorInfos
}

func (mpm *MainPageModel) callUpdateActorInfos(callback func()) {
	updatedActorInfos := mpm.getUpdatedActors()
	defer callback()
	if len(updatedActorInfos) == 0 {
		message.ErrorStr(mpm.VM, "Could not find any updated actors", false)
		return
	}

	req := xhr.NewRequest("PUT", "/api/actorinfos")
	req.Timeout = tools.TimeOut
	req.ResponseType = xhr.JSON
	err := req.Send(json.Stringify(updatedActorInfos))
	if err != nil {
		message.ErrorStr(mpm.VM, "Oups! "+err.Error(), true)
		return
	}
	if req.Status != tools.HttpOK {
		message.ErrorRequestMessage(mpm.VM, req)
		return
	}
	nbUpd := len(updatedActorInfos)
	msg := " modification sauvegardée"
	if nbUpd > 1 {
		msg = " modifications sauvegardées"
	}
	message.NotifySuccess(mpm.VM, "Informations Equipes", strconv.Itoa(nbUpd)+msg)
}

func (mpm *MainPageModel) getUpdatedActorInfos() []*actorinfo.ActorInfo {
	refActors := make(map[int]*actorinfo.ActorInfo)
	json.Parse(mpm.ActorInfosReference).Call("forEach", func(item *js.Object) {
		act := actorinfo.NewActorInfoFromJS(item)
		refActors[act.ActorId] = act
	})
	udpActorInfos := []*actorinfo.ActorInfo{}
	for _, actInf := range mpm.ActorInfos {
		refact := refActors[actInf.ActorId]
		if !(refact != nil && json.Stringify(actInf) == json.Stringify(refact)) {
			udpActorInfos = append(udpActorInfos, actInf)
		}
	}
	return udpActorInfos
}
