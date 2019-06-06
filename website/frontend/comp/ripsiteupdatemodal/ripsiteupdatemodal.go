package ripsiteupdatemodal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/modal"
	"github.com/lpuig/ewin/doe/website/frontend/comp/ripjunctionupdate"
	"github.com/lpuig/ewin/doe/website/frontend/comp/ripmeasurementupdate"
	"github.com/lpuig/ewin/doe/website/frontend/comp/rippullingupdate"
	"github.com/lpuig/ewin/doe/website/frontend/comp/ripsiteinfo"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	fmrip "github.com/lpuig/ewin/doe/website/frontend/model/ripsite"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements/message"
	"github.com/lpuig/ewin/doe/website/frontend/tools/json"
	"honnef.co/go/js/xhr"
	"strconv"
)

type RipsiteUpdateModalModel struct {
	*modal.ModalModel

	ActivityMode   string         `js:"ActivityMode"`
	User           *fm.User       `js:"User"`
	Filter         string         `js:"filter"`
	EditedRipsite  *fmrip.Ripsite `js:"edited_ripsite"`
	CurrentRipsite *fmrip.Ripsite `js:"current_ripsite"`

	Saving            bool `js:"saving"`
	ShowConfirmDelete bool `js:"showconfirmdelete"`
}

func NewRipsiteUpdateModalModel(vm *hvue.VM) *RipsiteUpdateModalModel {
	rsumm := &RipsiteUpdateModalModel{
		ModalModel: modal.NewModalModel(vm),
	}

	rsumm.ActivityMode = "Pulling"
	rsumm.User = fm.NewUser()
	rsumm.Filter = ""
	rsumm.initSites()
	//rsumm.EditedRipsite = fmrip.NewRisite()
	//rsumm.CurrentRipsite = fmrip.NewRisite()

	rsumm.Saving = false
	rsumm.ShowConfirmDelete = false

	return rsumm
}

func RipsiteUpdateModalModelFromJS(o *js.Object) *RipsiteUpdateModalModel {
	rsumm := &RipsiteUpdateModalModel{
		ModalModel: &modal.ModalModel{Object: o},
	}
	return rsumm
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Component Methods

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("ripsite-update-modal", componentOption()...)
}

func componentOption() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		ripsiteinfo.RegisterComponent(),
		rippullingupdate.RegisterComponent(),
		ripjunctionupdate.RegisterComponent(),
		ripmeasurementupdate.RegisterComponent(),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewRipsiteUpdateModalModel(vm)
		}),
		hvue.MethodsOf(&RipsiteUpdateModalModel{}),
		hvue.Computed("isNewRipsite", func(vm *hvue.VM) interface{} {
			m := RipsiteUpdateModalModelFromJS(vm.Object)
			return m.CurrentRipsite.Id == -1
		}),
		hvue.Computed("hasChanged", func(vm *hvue.VM) interface{} {
			m := RipsiteUpdateModalModelFromJS(vm.Object)
			return m.HasChanged()
		}),
		hvue.Computed("hasWarning", func(vm *hvue.VM) interface{} {
			//m := RipsiteUpdateModalModelFromJS(vm.Object)
			//if len(m.CurrentProject.Audits) > 0 {
			//	return "warning"
			//}
			return "success"
		}),
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Modal Methods

func (rsumm *RipsiteUpdateModalModel) initSites() {
	rsumm.EditedRipsite = fmrip.NewRisite()
	rsumm.CurrentRipsite = fmrip.NewRisite()
}

func (rsumm *RipsiteUpdateModalModel) HasChanged() bool {
	if rsumm.EditedRipsite.Object == js.Undefined {
		return true
	}
	return rsumm.CurrentRipsite.SearchInString() != rsumm.EditedRipsite.SearchInString()
}

func (rsumm *RipsiteUpdateModalModel) Show(id int, user *fm.User) {
	rsumm.EditedRipsite = fmrip.NewRisite()
	if id < 0 {
		rsumm.CurrentRipsite = rsumm.EditedRipsite.Clone()
	}
	rsumm.User = user
	rsumm.ShowConfirmDelete = false
	if id >= 0 {
		rsumm.Loading = false
		go rsumm.callGetRipsite(id)
	}

	rsumm.ModalModel.Show()
}

func (rsumm *RipsiteUpdateModalModel) HideWithControl() {
	if rsumm.HasChanged() {
		message.ConfirmWarning(rsumm.VM, "OK pour perdre les changements effectués ?", rsumm.Hide)
		return
	}
	rsumm.Hide()
}

func (rsumm *RipsiteUpdateModalModel) Hide() {
	rsumm.ShowConfirmDelete = false
	rsumm.ModalModel.Hide()
	rsumm.initSites()
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Action Button Methods

func (rsumm *RipsiteUpdateModalModel) Attachment() string {
	url := "/api/ripsites/" + strconv.Itoa(rsumm.CurrentRipsite.Id) + "/attach"
	return url
}

func (rsumm *RipsiteUpdateModalModel) ConfirmChange() {
	rsumm.Saving = true
	if rsumm.CurrentRipsite.Id >= 0 {
		go rsumm.callUpdateRipsite(rsumm.CurrentRipsite)
	} else {
		go rsumm.callCreateRipsite(rsumm.CurrentRipsite)
	}
}

func (rsumm *RipsiteUpdateModalModel) UndoChange() {
	rsumm.CurrentRipsite.Copy(rsumm.EditedRipsite)
}

func (rsumm *RipsiteUpdateModalModel) DeleteRipsite() {
	rsumm.Saving = true
	go rsumm.callDeleteRipsite(rsumm.CurrentRipsite)
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// WS call Methods

func (rsumm *RipsiteUpdateModalModel) SetActivityMode() {
	if len(rsumm.CurrentRipsite.Pullings) > 0 {
		rsumm.ActivityMode = "Pulling"
		return
	}
	if len(rsumm.CurrentRipsite.Junctions) > 0 {
		rsumm.ActivityMode = "Junction"
		return
	}
	if len(rsumm.CurrentRipsite.Pullings) > 0 {
		rsumm.ActivityMode = "Measurement"
		return
	}
}

func (rsumm *RipsiteUpdateModalModel) callGetRipsite(id int) {
	defer func() { rsumm.Loading = false }()
	req := xhr.NewRequest("GET", "/api/ripsites/"+strconv.Itoa(id))
	req.Timeout = tools.TimeOut
	req.ResponseType = xhr.JSON
	err := req.Send(nil)
	if err != nil {
		message.ErrorStr(rsumm.VM, "Oups! "+err.Error(), true)
		rsumm.Hide()
		return
	}
	if req.Status != tools.HttpOK {
		message.ErrorRequestMessage(rsumm.VM, req)
		rsumm.Hide()
		return
	}
	rsumm.EditedRipsite = fmrip.RipsiteFromJS(req.Response)
	rsumm.CurrentRipsite.Copy(rsumm.EditedRipsite)
	rsumm.SetActivityMode()
	return
}

func (rsumm *RipsiteUpdateModalModel) callUpdateRipsite(urs *fmrip.Ripsite) {
	defer func() { rsumm.Saving = false }()
	req := xhr.NewRequest("PUT", "/api/ripsites/"+strconv.Itoa(urs.Id))
	req.Timeout = tools.TimeOut
	req.ResponseType = xhr.JSON
	err := req.Send(json.Stringify(urs))
	if err != nil {
		message.ErrorStr(rsumm.VM, "Oups! "+err.Error(), true)
		return
	}
	if req.Status != tools.HttpOK {
		message.ErrorRequestMessage(rsumm.VM, req)
		return
	}
	rsumm.VM.Emit("update_ripsite")
	message.SuccesStr(rsumm.VM, "Chantier sauvegardé")
	rsumm.Hide()
}

func (rsumm *RipsiteUpdateModalModel) callCreateRipsite(urs *fmrip.Ripsite) {
	defer func() { rsumm.Saving = false }()
	req := xhr.NewRequest("POST", "/api/ripsites")
	req.Timeout = tools.TimeOut
	req.ResponseType = xhr.JSON
	err := req.Send(json.Stringify(urs))
	if err != nil {
		message.ErrorStr(rsumm.VM, "Oups! "+err.Error(), true)
		return
	}
	if req.Status != tools.HttpCreated {
		message.ErrorRequestMessage(rsumm.VM, req)
		return
	}
	rsumm.VM.Emit("update_ripsite")
	message.SuccesStr(rsumm.VM, "Nouveau chantier sauvegardé")
	rsumm.Hide()
}

func (rsumm *RipsiteUpdateModalModel) callDeleteRipsite(drs *fmrip.Ripsite) {
	defer func() { rsumm.Saving = false }()
	req := xhr.NewRequest("DELETE", "/api/ripsites/"+strconv.Itoa(drs.Id))
	req.Timeout = tools.TimeOut
	req.ResponseType = xhr.JSON
	err := req.Send(nil)
	if err != nil {
		message.ErrorStr(rsumm.VM, "Oups! "+err.Error(), true)
		return
	}
	if req.Status != tools.HttpOK {
		message.ErrorRequestMessage(rsumm.VM, req)
	}
	rsumm.VM.Emit("update_ripsite")
	message.SuccesStr(rsumm.VM, "Chantier supprimé !")
	rsumm.Hide()
}
