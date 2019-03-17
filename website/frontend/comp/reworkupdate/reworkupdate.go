package reworkupdate

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/worksitestatustag"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/autocomplete"
	"github.com/lpuig/ewin/doe/website/frontend/tools/dates"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Registration

func Register() {
	hvue.NewComponent("rework-update",
		componentOptions()...,
	)
}

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("rework-update", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		worksitestatustag.RegisterComponent(),
		hvue.Template(template),
		hvue.Props("worksite", "user"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewReworkUpdateModel(vm)
		}),
		hvue.Computed("filteredReworks", func(vm *hvue.VM) interface{} {
			m := ReworkUpdateModelFromJS(vm.Object)
			return m.GetReworks()
		}),
		hvue.MethodsOf(&ReworkUpdateModel{}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type ReworkUpdateModel struct {
	*js.Object

	Worksite *fm.Worksite `js:"worksite"`
	//ReferenceWorksite *fm.Worksite `js:"refWorksite"`
	User *fm.User               `js:"user"`
	Pts  map[string]*fm.Troncon `js:"Pts"`

	VM *hvue.VM `js:"VM"`
}

func NewReworkUpdateModel(vm *hvue.VM) *ReworkUpdateModel {
	rum := &ReworkUpdateModel{Object: tools.O()}
	rum.VM = vm
	rum.Worksite = fm.NewWorkSite()
	//rum.ReferenceWorksite = nil
	rum.User = fm.NewUser()
	rum.Pts = make(map[string]*fm.Troncon)

	return rum
}

func ReworkUpdateModelFromJS(o *js.Object) *ReworkUpdateModel {
	return &ReworkUpdateModel{Object: o}
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Modal Methods

//func (rem *ReworkUpdateModel) TableRowClassName(rowInfo *js.Object) string {
//	//wsi := &fm.WorksiteInfo{Object: rowInfo.Get("row")}
//	return ""
//}

func (rum *ReworkUpdateModel) GetReworks() []*fm.Defect {
	res := []*fm.Defect{}
	rum.Pts = make(map[string]*fm.Troncon)
	for _, defect := range rum.Worksite.Rework.Defects {
		if defect.ToBeFixed {
			res = append(res, defect)
			//Pts[defect.PT] = rumm.CurrentWorksite.GetPtByName(defect.PT)
			rum.Get("Pts").Set(defect.PT, rum.Worksite.GetPtByName(defect.PT))
		}
	}
	return res
}

func (rum *ReworkUpdateModel) UserSearch(vm *hvue.VM, query string, callback *js.Object) {
	rum = ReworkUpdateModelFromJS(vm.Object)
	res := []*autocomplete.Result{}
	q := strings.ToLower(query)
	client := rum.User.GetClientByName(rum.Worksite.Client)
	if client == nil {
		callback.Invoke(res)
		return
	}
	for _, team := range client.Teams {
		if (q == "" && team.IsActive) || (q != "" && strings.Contains(strings.ToLower(team.Members), q)) {
			res = append(res, autocomplete.NewResult(team.Members))
		}
	}
	callback.Invoke(res)
}

func (rum *ReworkUpdateModel) GetTronconRef(vm *hvue.VM, ptref string) string {
	rum = ReworkUpdateModelFromJS(vm.Object)
	tr := rum.Pts[ptref]
	if tr == nil {
		return "PT non trouvé"
	}
	return tr.Pb.Ref + " / " + tr.Pb.RefPt
}

func (rum *ReworkUpdateModel) GetTronconAddress(vm *hvue.VM, ptref string) string {
	rum = ReworkUpdateModelFromJS(vm.Object)
	tr := rum.Pts[ptref]
	if tr == nil {
		return "PT non trouvé"
	}
	return tr.Pb.Address
}

func (rum *ReworkUpdateModel) GetTronconInstallDate(vm *hvue.VM, ptref string) string {
	rum = ReworkUpdateModelFromJS(vm.Object)
	tr := rum.Pts[ptref]
	if tr == nil {
		return ""
	}
	return "Install.: " + date.DateString(tr.InstallDate)
}

func (rum *ReworkUpdateModel) GetTronconInstallActor(vm *hvue.VM, ptref string) string {
	rum = ReworkUpdateModelFromJS(vm.Object)
	tr := rum.Pts[ptref]
	if tr == nil {
		return ""
	}
	return "par: " + tr.InstallActor
}
