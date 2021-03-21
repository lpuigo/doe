package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/vehiculestable"
	"github.com/lpuig/ewin/doe/website/frontend/comp/vehiculeupdatemodal"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/model/actor"
	"github.com/lpuig/ewin/doe/website/frontend/model/ref"
	"github.com/lpuig/ewin/doe/website/frontend/model/vehicule"
	"github.com/lpuig/ewin/doe/website/frontend/model/vehicule/vehiculeconst"
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
		hvue.El("#vehicule_app"),
		vehiculeupdatemodal.RegisterComponent(),
		vehiculestable.RegisterComponent(),
		hvue.DataS(mpm),
		hvue.MethodsOf(mpm),
		hvue.Mounted(func(vm *hvue.VM) {
			mpm := &MainPageModel{Object: vm.Object}
			tools.BeforeUnloadConfirmation(mpm.PreventLeave)
			mpm.GetUserSession(func() {
				mpm.LoadVehicules(true)
			})
		}),
		hvue.Computed("IsDirty", func(vm *hvue.VM) interface{} {
			mpm := &MainPageModel{Object: vm.Object}
			return mpm.Ref.IsDirty()
		}),
	)

	js.Global.Set("mpm", mpm)
}

type MainPageModel struct {
	*js.Object
	Ref *ref.Ref `js:"Ref"`

	VM       *hvue.VM          `js:"VM"`
	User     *fm.User          `js:"User"`
	ActorStr *actor.ActorStore `js:"ActorStr"`

	//ActiveMode  string            `js:"ActiveMode"`
	Filter         string               `js:"Filter"`
	FilterType     string               `js:"FilterType"`
	Vehicules      []*vehicule.Vehicule `js:"Vehicules"`
	NextVehiculeId int                  `js:"NextVehiculeId"`
}

func NewMainPageModel() *MainPageModel {
	mpm := &MainPageModel{Object: tools.O()}
	mpm.VM = nil
	mpm.User = fm.NewUser()
	mpm.ActorStr = actor.NewActorStore()
	mpm.Filter = ""
	mpm.FilterType = ""

	mpm.Vehicules = []*vehicule.Vehicule{}
	mpm.Ref = ref.NewRef(func() string {
		return json.Stringify(mpm.Vehicules)
	})
	mpm.NextVehiculeId = -2

	return mpm
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Update & Undo Methods

func (mpm *MainPageModel) PreventLeave() bool {
	return mpm.Ref.Dirty
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// User Management Methods

func (mpm *MainPageModel) GetUserSession(callback func()) {
	go mpm.callGetUser(callback)
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Action Methods

func (mpm *MainPageModel) AddVehicule() {
	mpm.ShowEditVehicule(mpm.VM, vehicule.NewVehicule())
}

func (mpm *MainPageModel) ShowAddVehicule() bool {
	return true
}

func (mpm *MainPageModel) HandleEditedVehicule(vehic *vehicule.Vehicule) {
	if vehic.Id == -1 { // edited vehic is a new vehic => add it to the Vehicule array
		vehic.Id = mpm.NextVehiculeId
		//vehic.Info.ActorId = vehic.Id
		mpm.Vehicules = append(mpm.Vehicules, vehic)
		mpm.NextVehiculeId--
	}
}

func (mpm *MainPageModel) LoadVehicules(init bool) {
	onLoadedVehicules := func() {
		mpm.Ref.SetReference()
		// IsDirty is set to true if some update are undertaken
		if init && mpm.Ref.IsDirty() {
			mpm.SaveVehicules(mpm.VM)
		}
	}
	go mpm.callGetVehicules(onLoadedVehicules)
	mpm.ActorStr.CallGetActors(mpm.VM, func() {})
}

func (mpm *MainPageModel) SaveVehicules(vm *hvue.VM) {
	mpm = &MainPageModel{Object: vm.Object}

	if mpm.Ref.Dirty {
		onUpdatedVehicules := func() {
			mpm.LoadVehicules(false)
			//mpm.SetReference()
		}
		go mpm.callUpdateVehicules(onUpdatedVehicules)
	}
}

func (mpm *MainPageModel) GetFilterType(vm *hvue.VM) []*elements.ValueLabel {
	return vehicule.GetFilterTypeValueLabel()
}

//
func (mpm *MainPageModel) ClearFilter(vm *hvue.VM) {
	mpm = &MainPageModel{Object: vm.Object}
	mpm.Filter = ""
	mpm.FilterType = vehiculeconst.FilterValueAll
	mpm.ApplyFilter(vm)
}

//
func (mpm *MainPageModel) ApplyFilter(vm *hvue.VM) {
	// No OP
}

//
func (mpm *MainPageModel) ShowEditVehicule(vm *hvue.VM, vehic *vehicule.Vehicule) {
	vum := vehiculeupdatemodal.VehiculeUpdateModalModelFromJS(mpm.VM.Refs("VehiculeUpdateModal"))
	vum.Show(vehic, mpm.User, mpm.ActorStr, mpm.Vehicules)
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

// WS Vehicules

func (mpm *MainPageModel) callGetVehicules(callback func()) {
	req := xhr.NewRequest("GET", "/api/vehicules")
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
	loadedVehicules := []*vehicule.Vehicule{}
	req.Response.Call("forEach", func(item *js.Object) {
		vehic := vehicule.VehiculeFromJS(item)
		loadedVehicules = append(loadedVehicules, vehic)
	})
	mpm.Vehicules = loadedVehicules
	callback()
}

func (mpm *MainPageModel) callUpdateVehicules(callback func()) {
	updatedVehicules := mpm.getUpdatedVehicules()
	if len(updatedVehicules) == 0 {
		message.ErrorStr(mpm.VM, "Could not find any updated vehicules", false)
		return
	}

	req := xhr.NewRequest("PUT", "/api/vehicules")
	req.Timeout = tools.TimeOut
	req.ResponseType = xhr.JSON
	err := req.Send(json.Stringify(updatedVehicules))
	if err != nil {
		message.ErrorStr(mpm.VM, "Oups! "+err.Error(), true)
		return
	}
	if req.Status != tools.HttpOK {
		message.ErrorRequestMessage(mpm.VM, req)
		return
	}
	message.NotifySuccess(mpm.VM, "Véhicules", "Modifications sauvegardées")
	callback()
}

func (mpm *MainPageModel) getUpdatedVehicules() []*vehicule.Vehicule {
	refVehicules := []*vehicule.Vehicule{}
	json.Parse(mpm.Ref.Reference).Call("forEach", func(item *js.Object) {
		vehic := vehicule.VehiculeFromJS(item)
		refVehicules = append(refVehicules, vehic)
	})
	refDict := makeDictVehicules(refVehicules)

	udpVehicules := []*vehicule.Vehicule{}
	for _, vehic := range mpm.Vehicules {
		if vehic.Immat == "" {
			continue
		}
		refVehic := refDict[vehic.Id]
		if !(refVehic != nil && json.Stringify(vehic) == json.Stringify(refVehic)) {
			udpVehicules = append(udpVehicules, vehic)
		}
	}
	return udpVehicules
}

func makeDictVehicules(vehicules []*vehicule.Vehicule) map[int]*vehicule.Vehicule {
	res := make(map[int]*vehicule.Vehicule)
	for _, vehic := range vehicules {
		res[vehic.Id] = vehic
	}
	return res
}
