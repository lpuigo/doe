package reworkeditmodal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/tronconstatustag"
	wem "github.com/lpuig/ewin/doe/website/frontend/comp/worksiteeditmodal"
	"github.com/lpuig/ewin/doe/website/frontend/comp/worksiteinfo"
	"github.com/lpuig/ewin/doe/website/frontend/comp/worksitestatustag"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
)

type ReworkEditModalModel struct {
	*wem.WorksiteEditModalModel
}

func NewReworkEditModalModel(vm *hvue.VM) *ReworkEditModalModel {
	remm := &ReworkEditModalModel{WorksiteEditModalModel: wem.NewWorksiteEditModalModel(vm)}

	return remm
}

func NewReworkEditModalModelFromJS(o *js.Object) *ReworkEditModalModel {
	remm := &ReworkEditModalModel{WorksiteEditModalModel: &wem.WorksiteEditModalModel{Object: o}}
	return remm
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Component Methods

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("rework-edit-modal", ComponentOptions()...)
}

func ComponentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		worksiteinfo.RegisterComponent(),
		tronconstatustag.RegisterComponent(),
		worksitestatustag.RegisterComponent(),
		hvue.Template(template),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewReworkEditModalModel(vm)
		}),
		hvue.MethodsOf(&ReworkEditModalModel{}),
		hvue.Computed("HasRework", func(vm *hvue.VM) interface{} {
			m := NewReworkEditModalModelFromJS(vm.Object)
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
			m := NewReworkEditModalModelFromJS(vm.Object)
			return m.GetReworks()
		}),
		hvue.Computed("hasChanged", func(vm *hvue.VM) interface{} {
			m := NewReworkEditModalModelFromJS(vm.Object)
			return m.HasChanged()
		}),
		hvue.Computed("hasWarning", func(vm *hvue.VM) interface{} {
			//m := &WorksiteEditModalModel{Object: vm.Object}
			//if len(m.CurrentProject.Audits) > 0 {
			//	return "warning"
			//}
			return "success"
		}),
		//hvue.Filter("FormatTronconRef", func(vm *hvue.VM, value *js.Object, args ...*js.Object) interface{} {
		//	m := NewReworkEditModalModelFromJS(vm.Object)
		//	t := &fm.Troncon{Object: value}
		//	return m.GetFormatTronconRef(t)
		//}),
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Modal Methods

func (remm *ReworkEditModalModel) TableRowClassName(rowInfo *js.Object) string {
	//wsi := &fm.WorksiteInfo{Object: rowInfo.Get("row")}
	return ""
}

func (remm *ReworkEditModalModel) GetReworks() []*fm.Defect {
	return remm.CurrentWorksite.Rework.Defects
}

func (remm *ReworkEditModalModel) GetPTs() []*elements.ValueLabel {
	res := []*elements.ValueLabel{}
	for _, o := range remm.CurrentWorksite.Orders {
		for _, t := range o.Troncons {
			label := t.Pb.Ref + " / " + t.Pb.RefPt + " (" + t.Ref + ")"
			res = append(res, elements.NewValueLabel(t.Pb.RefPt, label))
		}
	}
	return res
}

func (remm *ReworkEditModalModel) AddDefect() {
	//m := NewReworkEditModalModelFromJS(vm.Object)
	r := remm.CurrentWorksite.Rework
	d := fm.NewDefect()
	d.SubmissionDate = r.SubmissionDate
	r.Defects = append(r.Defects, d)
}

func (remm *ReworkEditModalModel) RemoveDefect(i int) {
	remm.CurrentWorksite.Rework.Object.Get("Defects").Call("splice", i, 1)
}
