package teamproductivitymodal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/modal"
	"github.com/lpuig/ewin/doe/website/frontend/comp/ripteamproductivitychart"
	"github.com/lpuig/ewin/doe/website/frontend/comp/teamproductivitychart"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	rs "github.com/lpuig/ewin/doe/website/frontend/model/ripsite"
	"github.com/lpuig/ewin/doe/website/frontend/model/worksite"
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
		ripteamproductivitychart.RegisterComponent(),
		teamproductivitychart.RegisterComponent(),
		hvue.Template(template),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewTeamProductivityModalModel(vm)
		}),
		hvue.MethodsOf(&TeamProductivityModalModel{}),
		hvue.Computed("GetSites", func(vm *hvue.VM) interface{} {
			//tpmm := TeamProductivityModalModelFromJS(vm.Object)
			return []string{"site 1", "site 2", "site 3", "site 4"}
		}),
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Modal Methods

type TeamProductivityModalModel struct {
	*modal.ModalModel

	User     *fm.User `js:"user"`
	SiteMode string   `js:"SiteMode"`

	Stats     *worksite.WorksiteStats `js:"Stats"`
	TeamStats []*worksite.TeamStats   `js:"TeamStats"`

	RipStats     *rs.RipsiteStats `js:"RipStats"`
	RipTeamStats []*rs.TeamStats  `js:"RipTeamStats"`
	//AllSites         bool                                  `js:"AllSites"`
	//FewSitesSelected bool                                  `js:"FewSitesSelected"`
	SelectedSites map[string]bool                       `js:"SelectedSites"`
	SiteColors    ripteamproductivitychart.SiteColorMap `js:"SiteColors"`

	ActiveMode string `js:"ActiveMode"`
	GroupMode  string `js:"GroupMode"`
}

func NewTeamProductivityModalModel(vm *hvue.VM) *TeamProductivityModalModel {
	tpmm := &TeamProductivityModalModel{
		ModalModel: modal.NewModalModel(vm),
	}
	tpmm.User = fm.NewUser()
	tpmm.SiteMode = ""

	tpmm.Stats = worksite.NewWorksiteStats()
	tpmm.TeamStats = []*worksite.TeamStats{}

	tpmm.RipStats = rs.NewRipsiteStats()
	tpmm.RipTeamStats = []*rs.TeamStats{}
	//tpmm.AllSites = true
	//tpmm.FewSitesSelected = false
	tpmm.SelectedSites = map[string]bool{}
	tpmm.SiteColors = ripteamproductivitychart.SiteColorMap{}

	tpmm.ActiveMode = "week"
	tpmm.GroupMode = "activity"
	return tpmm
}

func TeamProductivityModalModelFromJS(o *js.Object) *TeamProductivityModalModel {
	tpmm := &TeamProductivityModalModel{
		ModalModel: &modal.ModalModel{Object: o},
	}
	return tpmm
}

func (tpmm *TeamProductivityModalModel) Show(user *fm.User, siteMode string) {
	tpmm.Stats = worksite.NewWorksiteStats()
	tpmm.RipStats = rs.NewRipsiteStats()
	tpmm.TeamStats = []*worksite.TeamStats{}
	tpmm.RipTeamStats = []*rs.TeamStats{}
	tpmm.User = user
	tpmm.SiteMode = siteMode
	tpmm.RefreshStat()
	tpmm.ModalModel.Show()
}

func (tpmm *TeamProductivityModalModel) RefreshStat() {
	tpmm.Loading = true
	if tpmm.SiteMode == "Rip" {
		go tpmm.callGetRipsitesStats("/api/ripsites/stat/" + tpmm.GroupMode + "/" + tpmm.ActiveMode)
		return
	}
	if tpmm.SiteMode == "Poles" {
		go tpmm.callGetRipsitesStats("/api/polesites/stat/" + tpmm.ActiveMode)
		return
	}
	go tpmm.callGetWorksitesStats()
}

func (tpmm *TeamProductivityModalModel) initSitesColors() {
	workCM := ripteamproductivitychart.ColorMap{
		HueStart:   160,
		HueEnd:     360,
		Light:      40,
		Saturation: 60,
	}
	priceCM := ripteamproductivitychart.ColorMap{
		HueStart:   160,
		HueEnd:     360,
		Light:      55,
		Saturation: 70,
	}
	tpmm.SiteColors = ripteamproductivitychart.SetColor(tpmm.RipStats.Sites, workCM, priceCM)
}

//func (tpmm *TeamProductivityModalModel) CheckAllSitesChange() {
//	tpmm.FewSitesSelected = false
//	if tpmm.AllSites {
//		for site, _ := range tpmm.SelectedSites {
//			tpmm.Object.Get("SelectedSites").Set(site, true)
//		}
//		return
//	}
//	for site, _ := range tpmm.SelectedSites {
//		tpmm.Object.Get("SelectedSites").Set(site, false)
//	}
//	tpmm.RipTeamStats = tpmm.RipStats.CreateTeamStats(tpmm.SelectedSites)
//}

func (tpmm *TeamProductivityModalModel) CheckSitesChange() {
	//allFalse, allTrue := true, true
	//for _, value := range tpmm.SelectedSites {
	//	allFalse = allFalse && !value
	//	allTrue = allTrue && value
	//}
	//tpmm.FewSitesSelected = allFalse == allTrue
	//tpmm.AllSites = allTrue
	tpmm.RipTeamStats = tpmm.RipStats.CreateTeamStats(tpmm.SelectedSites)
}

func (tpmm *TeamProductivityModalModel) SiteCircleStyle(site string) string {
	return "color: " + tpmm.SiteColors.GetWorkColor(site)
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// WS call Methods

func (tpmm *TeamProductivityModalModel) callGetWorksitesStats() {
	defer func() { tpmm.Loading = false }()
	url := "/api/worksites/stat/"
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
	tpmm.Stats = worksite.WorksiteStatsFromJs(req.Response)
	tpmm.TeamStats = tpmm.Stats.CreateTeamStats()
	return
}

func (tpmm *TeamProductivityModalModel) callGetRipsitesStats(url string) {
	defer func() { tpmm.Loading = false }()
	req := xhr.NewRequest("GET", url)
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
	tpmm.RipStats = rs.RipsiteStatsFromJs(req.Response)
	tpmm.initSitesColors()
	tpmm.SelectedSites = tpmm.RipStats.Sites
	tpmm.CheckSitesChange()
	//tpmm.RipTeamStats = tpmm.RipStats.CreateTeamStats(tpmm.SelectedSites)
}
