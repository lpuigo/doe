package reworkupdatemodal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	wem "github.com/lpuig/ewin/doe/website/frontend/comp/worksiteeditmodal"
	"github.com/lpuig/ewin/doe/website/frontend/comp/worksiteinfo"
	"github.com/lpuig/ewin/doe/website/frontend/comp/worksitestatustag"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
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
		hvue.Filter("FormatTronconRef", func(vm *hvue.VM, value *js.Object, args ...*js.Object) interface{} {
			m := NewReworkUpdateModalModelFromJS(vm.Object)
			defect := &fm.Defect{Object: value}
			tr := m.Pts[defect.PT]
			if tr == nil {
				return "PT non trouvé"
			}
			return tr.Pb.Ref + " / " + tr.Pb.RefPt
		}),
		hvue.Filter("FormatTronconAddress", func(vm *hvue.VM, value *js.Object, args ...*js.Object) interface{} {
			m := NewReworkUpdateModalModelFromJS(vm.Object)
			defect := &fm.Defect{Object: value}
			tr := m.Pts[defect.PT]
			if tr == nil {
				return ""
			}
			return tr.Pb.Address
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
			rumm.Pts[defect.PT] = rumm.CurrentWorksite.GetPtByName(defect.PT)
		}
	}
	return res
}

//func (rumm *ReworkUpdateModalModel) GetPTs() []*elements.ValueLabel {
//	res := []*elements.ValueLabel{}
//	for _, o := range rumm.CurrentWorksite.Orders {
//		for _, t := range o.Troncons {
//			label := t.Pb.Ref + " / " + t.Pb.RefPt + " (" + t.Ref + ")"
//			res = append(res, elements.NewValueLabel(t.Pb.RefPt, label))
//		}
//	}
//	return res
//}

//func (rumm *ReworkUpdateModalModel) AddDefect() {
//	//m := NewReworkUpdateModalModelFromJS(vm.Object)
//	r := rumm.CurrentWorksite.Rework
//	d := fm.NewDefect()
//	d.SubmissionDate = r.ControlDate
//	r.Defects = append(r.Defects, d)
//}
//
//func (rumm *ReworkUpdateModalModel) RemoveDefect(i int) {
//	rumm.CurrentWorksite.Rework.Object.Get("Defects").Call("splice", i, 1)
//}
