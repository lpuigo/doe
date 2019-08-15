package adminmodal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/modal"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements/message"
	"honnef.co/go/js/xhr"
)

type AdminModalModel struct {
	*modal.ModalModel

	User *fm.User `js:"user"`
}

func NewAdminModalModel(vm *hvue.VM) *AdminModalModel {
	tpmm := &AdminModalModel{
		ModalModel: modal.NewModalModel(vm),
	}
	tpmm.User = fm.NewUser()
	return tpmm
}

func AdminModalModelFromJS(o *js.Object) *AdminModalModel {
	tpmm := &AdminModalModel{
		ModalModel: &modal.ModalModel{Object: o},
	}
	return tpmm
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Component Methods

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("admin-modal", componentOption()...)
}

func componentOption() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewAdminModalModel(vm)
		}),
		hvue.MethodsOf(&AdminModalModel{}),
	}
}

func (amm *AdminModalModel) ReloadData() {
	go amm.callReloadData()
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Modal Methods

func (amm *AdminModalModel) Show(user *fm.User) {
	amm.User = user
	amm.Loading = false
	amm.ModalModel.Show()
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// WS call Methods

func (amm *AdminModalModel) callReloadData() {
	defer func() { amm.Loading = false }()
	req := xhr.NewRequest("GET", "/api/admin/reload")
	req.Timeout = tools.TimeOut
	req.ResponseType = xhr.JSON
	err := req.Send(nil)
	if err != nil {
		message.ErrorStr(amm.VM, "Oups! "+err.Error(), true)
		amm.Hide()
		return
	}
	if req.Status != tools.HttpOK {
		message.ErrorRequestMessage(amm.VM, req)
		amm.Hide()
		return
	}
	message.SuccesStr(amm.VM, "Rechargement des données effectué")
	amm.VM.Emit("reload")
	return
}
