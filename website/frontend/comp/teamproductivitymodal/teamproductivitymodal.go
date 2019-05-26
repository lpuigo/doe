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
)

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

type TeamProductivityModalModel struct {
	*modal.ModalModel

	User       *fm.User          `js:"user"`
	SiteMode   string            `js:"SiteMode"`
	Stats      *fm.WorksiteStats `js:"Stats"`
	TeamStats  []*fm.TeamStats   `js:"TeamStats"`
	ActiveMode string            `js:"ActiveMode"`
}

func NewTeamProductivityModalModel(vm *hvue.VM) *TeamProductivityModalModel {
	tpmm := &TeamProductivityModalModel{
		ModalModel: modal.NewModalModel(vm),
	}
	tpmm.Stats = fm.NewWorksiteStats()
	tpmm.User = fm.NewUser()
	tpmm.SiteMode = ""
	tpmm.TeamStats = []*fm.TeamStats{}
	tpmm.ActiveMode = "week"
	return tpmm
}

func TeamProductivityModalModelFromJS(o *js.Object) *TeamProductivityModalModel {
	tpmm := &TeamProductivityModalModel{
		ModalModel: &modal.ModalModel{Object: o},
	}
	return tpmm
}

func (tpmm *TeamProductivityModalModel) Show(user *fm.User, siteMode string) {
	tpmm.Stats = fm.NewWorksiteStats()
	tpmm.TeamStats = []*fm.TeamStats{}
	tpmm.User = user
	tpmm.SiteMode = siteMode
	tpmm.RefreshStat()
	tpmm.ModalModel.Show()
}

func (tpmm *TeamProductivityModalModel) RefreshStat() {
	tpmm.Loading = true
	go tpmm.callGetWorksitesStats()
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// WS call Methods

func (tpmm *TeamProductivityModalModel) callGetWorksitesStats() {
	defer func() { tpmm.Loading = false }()
	url := ""
	switch tpmm.SiteMode {
	case "Rip":
		url = "/api/ripsites/stat/"
	default:
		url = "/api/worksites/stat/"
	}
	req := xhr.NewRequest("GET", url+tpmm.ActiveMode)
	req.Timeout = tools.TimeOut
	req.ResponseType = xhr.JSON
	err := req.Send(nil)
	if err != nil {
		message.ErrorStr(tpmm.VM, "Oups! "+err.Error(), true)
		tpmm.Hide()
		return
	}
	if req.Status != tools.HttpOK {
		message.ErrorRequestMessage(tpmm.VM, req)
		tpmm.Hide()
		return
	}
	tpmm.Stats = fm.WorksiteStatsFromJs(req.Response)
	tpmm.TeamStats = tpmm.Stats.CreateTeamStats()
	return
}
