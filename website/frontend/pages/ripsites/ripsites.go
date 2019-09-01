package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/ripinfoupdate"
	"github.com/lpuig/ewin/doe/website/frontend/comp/ripjunctionupdate"
	"github.com/lpuig/ewin/doe/website/frontend/comp/ripmeasurementupdate"
	"github.com/lpuig/ewin/doe/website/frontend/comp/rippullingupdate"
	"github.com/lpuig/ewin/doe/website/frontend/comp/ripsiteinfo"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/model/ripsite"
	fmrip "github.com/lpuig/ewin/doe/website/frontend/model/ripsite"
	"github.com/lpuig/ewin/doe/website/frontend/model/ripsite/ripconst"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
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
		ripsiteinfo.RegisterComponent(),
		ripinfoupdate.RegisterComponent(),
		rippullingupdate.RegisterComponent(),
		ripjunctionupdate.RegisterComponent(),
		ripmeasurementupdate.RegisterComponent(),
		hvue.El("#ripsites_app"),
		hvue.DataS(mpm),
		hvue.MethodsOf(mpm),
		hvue.Mounted(func(vm *hvue.VM) {
			mpm := &MainPageModel{Object: vm.Object}
			tools.BeforeUnloadConfirmation(mpm.PreventLeave)
			mpm.GetUserSession(func() {
				mpm.LoadRipsite(false)
			})
		}),
		hvue.Computed("Title", func(vm *hvue.VM) interface{} {
			mpm := &MainPageModel{Object: vm.Object}
			if mpm.Ripsite.Object == js.Undefined {
				return ""
			}
			if tools.Empty(mpm.Ripsite.Client) && tools.Empty(mpm.Ripsite.Ref) {
				return ""
			}
			return mpm.Ripsite.Client + " / " + mpm.Ripsite.Ref
		}),
		hvue.Computed("IsDirty", func(vm *hvue.VM) interface{} {
			mpm := &MainPageModel{Object: vm.Object}
			mpm.Dirty = (mpm.Reference != json.Stringify(mpm.Ripsite))
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

	ActivityMode string           `js:"ActivityMode"`
	Filter       string           `js:"Filter"`
	FilterType   string           `js:"FilterType"`
	Ripsite      *ripsite.Ripsite `js:"Ripsite"`
	Reference    string           `js:"Reference"`
	Dirty        bool             `js:"Dirty"`
}

func NewMainPageModel() *MainPageModel {
	mpm := &MainPageModel{Object: tools.O()}
	mpm.VM = nil
	mpm.User = fm.NewUser()
	mpm.ActivityMode = "Info"
	mpm.Filter = ""
	mpm.FilterType = ripconst.FilterValueAll
	mpm.Ripsite = fmrip.NewRisite()
	mpm.Reference = ""
	mpm.Dirty = false
	return mpm
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// User Management Methods

func (mpm *MainPageModel) GetUserSession(callback func()) {
	go mpm.callGetUser(callback)
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Action Methods

func (mpm *MainPageModel) PreventLeave() bool {
	return mpm.Dirty
}

func (mpm *MainPageModel) LoadRipsite(update bool) {
	orsid := tools.GetURLSearchParam("rsid")
	if orsid == nil {
		print("rsid undefined")
		return
	}
	if orsid.String() == "" {
		print("rsid empty")
		return
	}

	callback := mpm.SetActivityMode
	go mpm.callGetRipsite(orsid.Int(), callback)
}

func (mpm *MainPageModel) SaveRipsite(vm *hvue.VM) {
	mpm = &MainPageModel{Object: vm.Object}
	go mpm.callUpdateRipsite(mpm.Ripsite)
}

//// UpdateRip updates current Ripsite component
//func (mpm *MainPageModel) UpdateRip() {
//}

// SwitchActiveMode handles ActiveMode change
func (mpm *MainPageModel) SwitchActiveMode(vm *hvue.VM) {
	// TODO Map Display None
	//mpm = &MainPageModel{Object: vm.Object}
	//switch mpm.ActiveMode {
	//case "Map":
	//	mpm.UpdateRip()
	//default:
	//}
}

func (mpm *MainPageModel) GetFilterType(am string) []*elements.ValueLabel {
	switch am {
	//case "Pulling":
	//	return fmrip.GetPullingFilterTypeValueLabel()
	case "Junction":
		return fmrip.GetJunctionFilterTypeValueLabel()
	case "Measurement":
		return fmrip.GetMeasurementFilterTypeValueLabel()
	default: // Pulling
		return fmrip.GetPullingFilterTypeValueLabel()
	}
}

//
func (mpm *MainPageModel) ClearFilter(vm *hvue.VM) {
	mpm = &MainPageModel{Object: vm.Object}
	mpm.Filter = ""
	mpm.FilterType = ripconst.FilterValueAll
	mpm.ApplyFilter(vm)
}

//
func (mpm *MainPageModel) ApplyFilter(vm *hvue.VM) {
	mpm = &MainPageModel{Object: vm.Object}

}

//
func (mpm *MainPageModel) SetActivityMode() {
	mpm.ActivityMode = "Info"
	return
	if len(mpm.Ripsite.Pullings) > 0 {
		mpm.ActivityMode = "Pulling"
		return
	}
	if len(mpm.Ripsite.Junctions) > 0 {
		mpm.ActivityMode = "Junction"
		return
	}
	if len(mpm.Ripsite.Pullings) > 0 {
		mpm.ActivityMode = "Measurement"
		return
	}
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

func (mpm *MainPageModel) callGetRipsite(rsid int, callback func()) {
	req := xhr.NewRequest("GET", "/api/ripsites/"+strconv.Itoa(rsid))
	req.Timeout = tools.LongTimeOut
	req.ResponseType = xhr.JSON
	defer func() {
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
	mpm.Ripsite = fmrip.RipsiteFromJS(req.Response)
	mpm.Reference = json.Stringify(req.Response)
	newTitle := mpm.Ripsite.Ref + " - " + mpm.Ripsite.Client
	js.Global.Get("document").Set("title", newTitle)
}

func (mpm *MainPageModel) callUpdateRipsite(urs *fmrip.Ripsite) {
	//defer func() {}()
	req := xhr.NewRequest("PUT", "/api/ripsites/"+strconv.Itoa(urs.Id))
	req.Timeout = tools.TimeOut
	req.ResponseType = xhr.JSON
	err := req.Send(json.Stringify(urs))
	if err != nil {
		message.ErrorStr(mpm.VM, "Oups! "+err.Error(), true)
		return
	}
	if req.Status != tools.HttpOK {
		message.ErrorRequestMessage(mpm.VM, req)
		return
	}
	message.SuccesStr(mpm.VM, "Modification sauvegardée")
	mpm.Reference = json.Stringify(urs)
}
