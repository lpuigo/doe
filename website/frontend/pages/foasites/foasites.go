package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/foainfoupdate"
	"github.com/lpuig/ewin/doe/website/frontend/comp/foaupdate"
	"github.com/lpuig/ewin/doe/website/frontend/comp/foaupdatemodal"
	"github.com/lpuig/ewin/doe/website/frontend/comp/ripsiteinfo"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	fmfoa "github.com/lpuig/ewin/doe/website/frontend/model/foasite"
	"github.com/lpuig/ewin/doe/website/frontend/model/foasite/foaconst"
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
		ripsiteinfo.RegisterComponent(),
		foainfoupdate.RegisterComponent(),
		foaupdatemodal.RegisterComponent(),
		foaupdate.RegisterComponent(),
		hvue.El("#foasites_app"),
		hvue.DataS(mpm),
		hvue.MethodsOf(mpm),
		hvue.Mounted(func(vm *hvue.VM) {
			mpm := &MainPageModel{Object: vm.Object}
			tools.BeforeUnloadConfirmation(mpm.PreventLeave)
			mpm.GetUserSession(func() {
				mpm.LoadFoaSite(false)
			})
		}),
		hvue.Computed("Title", func(vm *hvue.VM) interface{} {
			mpm := &MainPageModel{Object: vm.Object}
			if mpm.Foasite.Object == js.Undefined {
				return ""
			}
			if tools.Empty(mpm.Foasite.Client) && tools.Empty(mpm.Foasite.Ref) {
				return ""
			}
			return mpm.Foasite.Client + " / " + mpm.Foasite.Ref
		}),
		hvue.Computed("IsDirty", func(vm *hvue.VM) interface{} {
			mpm := &MainPageModel{Object: vm.Object}
			mpm.Dirty = (mpm.Reference != json.Stringify(mpm.Foasite))
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

	ActivityMode string         `js:"ActivityMode"`
	Filter       string         `js:"Filter"`
	FilterType   string         `js:"FilterType"`
	Foasite      *fmfoa.FoaSite `js:"Foasite"`
	Reference    string         `js:"Reference"`
	Dirty        bool           `js:"Dirty"`
}

func NewMainPageModel() *MainPageModel {
	mpm := &MainPageModel{Object: tools.O()}
	mpm.VM = nil
	mpm.User = fm.NewUser()
	mpm.ActivityMode = "Info"
	mpm.Filter = ""
	mpm.FilterType = foaconst.FilterValueAll
	mpm.Foasite = fmfoa.NewFoaSite()
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
// Action & HTML Methods

func (mpm *MainPageModel) SaveFoaSite(vm *hvue.VM) {
	mpm = &MainPageModel{Object: vm.Object}
	go mpm.callUpdateFoasite(mpm.Foasite)
}

func (mpm *MainPageModel) LoadFoaSite(update bool) {
	ofsid := tools.GetURLSearchParam("fsid")
	if ofsid == nil {
		print("fsid undefined")
		return
	}
	if ofsid.String() == "" {
		print("fsid empty")
		return
	}

	callback := mpm.SetActivityMode
	go mpm.callGetFoasite(ofsid.Int(), callback)
}

// SwitchActiveMode handles ActiveMode change
//func (mpm *MainPageModel) SwitchActiveMode(vm *hvue.VM) {
//	// No Op
//}

func (mpm *MainPageModel) GetFilterType() []*elements.ValueLabel {
	return fmfoa.GetFilterType()
}

//
func (mpm *MainPageModel) ClearFilter(vm *hvue.VM) {
	mpm = &MainPageModel{Object: vm.Object}
	mpm.Filter = ""
	mpm.FilterType = foaconst.FilterValueAll
	mpm.ApplyFilter(vm)
}

//
func (mpm *MainPageModel) ApplyFilter(vm *hvue.VM) {
	mpm = &MainPageModel{Object: vm.Object}
}

//
func (mpm *MainPageModel) AddNewFoa() {
	fumm := foaupdatemodal.FoaUpdateModalModelFromJS(mpm.VM.Refs("FoaUpdateModal"))
	newFoa := fmfoa.NewFoa()
	newFoa.State.Date = date.TodayAfter(0)
	newFoa.Ref = "New Foa"
	onApply := func() {
		mpm.Foasite.AddFoa(newFoa)
	}
	fumm.ShowEdit(newFoa, onApply)
}

//
func (mpm *MainPageModel) UpdateState(foas *fmfoa.FoaSite, f *fmfoa.Foa) {
	fumm := foaupdatemodal.FoaUpdateModalModelFromJS(mpm.VM.Refs("FoaUpdateModal"))
	foaUpdate := foaupdate.FoaUpdateModelFromJS(mpm.VM.Refs("foaUpdateComp"))
	onApply := func() {
		foaUpdate.ClearSelection()
	}
	if f != nil && f.Object != nil { // if f is not nil, Edit mode on single foa pointed in foas
		fumm.SetModel(f)
		fumm.ShowEdit(f, onApply)
		return
	}
	fumm.Show(foas, onApply)
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Tool Methods

func (mpm *MainPageModel) PreventLeave() bool {
	return mpm.Dirty
}

//
func (mpm *MainPageModel) SetActivityMode() {
	// No Op
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

func (mpm *MainPageModel) callGetFoasite(fsid int, callback func()) {
	req := xhr.NewRequest("GET", "/api/foasites/"+strconv.Itoa(fsid))
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
	mpm.Foasite = fmfoa.FoaSiteFromJS(req.Response)
	mpm.Reference = json.Stringify(req.Response)
	newTitle := mpm.Foasite.Ref + " - " + mpm.Foasite.Client
	js.Global.Get("document").Set("title", newTitle)
}

func (mpm *MainPageModel) callUpdateFoasite(ufs *fmfoa.FoaSite) {
	//defer func() {}()
	req := xhr.NewRequest("PUT", "/api/foasites/"+strconv.Itoa(ufs.Id))
	req.Timeout = tools.TimeOut
	req.ResponseType = xhr.JSON
	err := req.Send(json.Stringify(ufs))
	if err != nil {
		message.ErrorStr(mpm.VM, "Oups! "+err.Error(), true)
		return
	}
	if req.Status != tools.HttpOK {
		message.ErrorRequestMessage(mpm.VM, req)
		return
	}
	message.NotifySuccess(mpm.VM, "FOA", "Modifications sauvegardées")
	mpm.Reference = json.Stringify(ufs)
}
