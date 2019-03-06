package teamproductivitymodal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/modal"
	"github.com/lpuig/ewin/doe/website/frontend/comp/teamproductivitychart"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements/message"
	"honnef.co/go/js/xhr"
	"strconv"
)

type TeamProductivityModalModel struct {
	*modal.ModalModel

	User      *fm.User          `js:"user"`
	Stats     *fm.WorksiteStats `js:"Stats"`
	TeamStats []*fm.TeamStats   `js:"TeamStats"`
}

func NewTeamProductivityModalModel(vm *hvue.VM) *TeamProductivityModalModel {
	tpmm := &TeamProductivityModalModel{
		ModalModel: modal.NewModalModel(vm),
	}
	tpmm.Stats = fm.NewWorksiteStats()
	tpmm.User = fm.NewUser()
	tpmm.TeamStats = []*fm.TeamStats{}
	return tpmm
}

func TeamProductivityModalModelFromJS(o *js.Object) *TeamProductivityModalModel {
	tpmm := &TeamProductivityModalModel{
		ModalModel: &modal.ModalModel{Object: o},
	}
	return tpmm
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Component Methods

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("team-productivity-modal", componentOption()...)
}

func componentOption() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		teamproductivitychart.RegisterComponent(),
		hvue.Template(template),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewTeamProductivityModalModel(vm)
		}),
		hvue.MethodsOf(&TeamProductivityModalModel{}),
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Modal Methods

func (tpmm *TeamProductivityModalModel) Show(user *fm.User) {
	tpmm.Stats = fm.NewWorksiteStats()
	tpmm.TeamStats = []*fm.TeamStats{}
	tpmm.User = user
	tpmm.Loading = true
	go tpmm.callGetWorksitesStats()

	tpmm.ModalModel.Show()
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// WS call Methods

func (tpmm *TeamProductivityModalModel) errorMessage(req *xhr.Request) {
	message.SetDuration(tools.WarningMsgDuration)
	msg := "Quelque chose c'est mal passé !\n"
	msg += "Le server retourne un code " + strconv.Itoa(req.Status) + "\n"
	message.ErrorMsgStr(tpmm.VM, msg, req.Response, true)
}

func (tpmm *TeamProductivityModalModel) callGetWorksitesStats() {
	defer func() { tpmm.Loading = false }()
	req := xhr.NewRequest("GET", "/api/worksites/stat")
	req.Timeout = tools.TimeOut
	req.ResponseType = xhr.JSON
	err := req.Send(nil)
	if err != nil {
		message.ErrorStr(tpmm.VM, "Oups! "+err.Error(), true)
		tpmm.Hide()
		return
	}
	if req.Status != tools.HttpOK {
		tpmm.errorMessage(req)
		tpmm.Hide()
		return
	}
	tpmm.Stats = fm.WorksiteStatsFromJs(req.Response)
	tpmm.TeamStats = tpmm.Stats.CreateTeamStats()
	return
}
