package reworkupdatemodal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	wem "github.com/lpuig/ewin/doe/website/frontend/comp/worksiteeditmodal"
	"github.com/lpuig/ewin/doe/website/frontend/comp/worksiteinfo"
	"github.com/lpuig/ewin/doe/website/frontend/comp/worksitestatustag"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/tools/autocomplete"
	"github.com/lpuig/ewin/doe/website/frontend/tools/dates"
	"strings"
)

type ReworkUpdateModalModel struct {
	*wem.WorksiteEditModalModel

	Pts map[string]*fm.Troncon `js:"Pts"`
}

func NewReworkUpdateModalModel(vm *hvue.VM) *ReworkUpdateModalModel {
	rumm := &ReworkUpdateModalModel{WorksiteEditModalModel: wem.NewWorksiteEditModalModel(vm)}
	rumm.Pts = make(map[string]*fm.Troncon)
	return rumm
}

func NewReworkUpdateModalModelFromJS(o *js.Object) *ReworkUpdateModalModel {
	rumm := &ReworkUpdateModalModel{WorksiteEditModalModel: &wem.WorksiteEditModalModel{Object: o}}
	return rumm
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Component Methods

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("rework-update-modal", ComponentOptions()...)
}

func ComponentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		worksiteinfo.RegisterComponent(),
		worksitestatustag.RegisterComponent(),
		hvue.Template(template),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewReworkUpdateModalModel(vm)
		}),
		hvue.MethodsOf(&ReworkUpdateModalModel{}),
		hvue.Computed("HasRework", func(vm *hvue.VM) interface{} {
			m := NewReworkUpdateModalModelFromJS(vm.Object)
			if m.Loading || !fm.WorksiteIsReworkable(m.CurrentWorksite.Status) {
				return false
			}
			if m.CurrentWorksite.Rework != nil && m.CurrentWorksite.Rework.Object != js.Undefined {
				return true
			}
			m.CurrentWorksite.Rework = fm.NewRework()
			return true
		}),
		hvue.Computed("filteredReworks", func(vm *hvue.VM) interface{} {
			m := NewReworkUpdateModalModelFromJS(vm.Object)
			return m.GetReworks()
		}),
		hvue.Computed("hasChanged", func(vm *hvue.VM) interface{} {
			m := NewReworkUpdateModalModelFromJS(vm.Object)
			return m.HasChanged()
		}),
		hvue.Computed("hasWarning", func(vm *hvue.VM) interface{} {
			//m := NewReworkUpdateModalModelFromJS(vm.Object)
			//if len(m.CurrentProject.Audits) > 0 {
			//	return "warning"
			//}
			return "success"
		}),
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Modal Methods

func (rumm *ReworkUpdateModalModel) TableRowClassName(rowInfo *js.Object) string {
	//wsi := &fm.WorksiteInfo{Object: rowInfo.Get("row")}
	return ""
}

func (rumm *ReworkUpdateModalModel) GetReworks() []*fm.Defect {
	res := []*fm.Defect{}
	rumm.Pts = make(map[string]*fm.Troncon)
	for _, defect := range rumm.CurrentWorksite.Rework.Defects {
		if defect.ToBeFixed {
			res = append(res, defect)
			//Pts[defect.PT] = rumm.CurrentWorksite.GetPtByName(defect.PT)
			rumm.Get("Pts").Set(defect.PT, rumm.CurrentWorksite.GetPtByName(defect.PT))

		}
	}
	return res
}

func (rumm *ReworkUpdateModalModel) UserSearch(vm *hvue.VM, query string, callback *js.Object) {
	users := fm.GetTeamUsers()

	q := strings.ToLower(query)

	res := []*autocomplete.Result{}
	for _, u := range users {
		if q == "" || strings.Contains(strings.ToLower(u), q) {
			res = append(res, autocomplete.NewResult(u))
		}
	}
	callback.Invoke(res)
}

func (rumm *ReworkUpdateModalModel) GetTronconRef(ptref string) string {
	tr := rumm.Pts[ptref]
	if tr == nil {
		return "PT non trouvé"
	}
	return tr.Pb.Ref + " / " + tr.Pb.RefPt
}

func (rumm *ReworkUpdateModalModel) GetTronconAddress(ptref string) string {
	tr := rumm.Pts[ptref]
	if tr == nil {
		return "PT non trouvé"
	}
	return tr.Pb.Address
}

func (rumm *ReworkUpdateModalModel) GetTronconInstallDate(ptref string) string {
	tr := rumm.Pts[ptref]
	if tr == nil {
		return ""
	}
	return "Install.: " + date.DateString(tr.InstallDate)
}

func (rumm *ReworkUpdateModalModel) GetTronconInstallActor(ptref string) string {
	tr := rumm.Pts[ptref]
	if tr == nil {
		return ""
	}
	return "par: " + tr.InstallActor
}
