package adminmodal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/modal"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
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

//////////////////////////////////////////////////////////////////////////////////////////////
// Modal Methods

func (amm *AdminModalModel) Show(user *fm.User) {
	amm.User = user
	amm.Loading = false
	amm.ModalModel.Show()
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// WS call Methods

//func (amm *AdminModalModel) errorMessage(req *xhr.Request) {
//	message.SetDuration(tools.WarningMsgDuration)
//	msg := "Quelque chose c'est mal passé !\n"
//	msg += "Le server retourne un code " + strconv.Itoa(req.Status) + "\n"
//	message.ErrorMsgStr(amm.VM, msg, req.Response, true)
//}
//
//func (amm *AdminModalModel) callGetWorksitesStats() {
//	defer func() { amm.Loading = false }()
//	req := xhr.NewRequest("GET", "/api/worksites/stat")
//	req.Timeout = tools.TimeOut
//	req.ResponseType = xhr.JSON
//	err := req.Send(nil)
//	if err != nil {
//		message.ErrorStr(amm.VM, "Oups! "+err.Error(), true)
//		amm.Hide()
//		return
//	}
//	if req.Status != tools.HttpOK {
//		amm.errorMessage(req)
//		amm.Hide()
//		return
//	}
//	amm.Stats = fm.WorksiteStatsFromJs(req.Response)
//	amm.TeamStats = amm.Stats.CreateTeamStats()
//	return
//}
