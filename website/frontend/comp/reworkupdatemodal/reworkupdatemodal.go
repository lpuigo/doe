package reworkupdatemodal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/reworkupdate"
	wem "github.com/lpuig/ewin/doe/website/frontend/comp/worksiteeditmodal"
	"github.com/lpuig/ewin/doe/website/frontend/comp/worksiteinfo"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
)

type ReworkUpdateModalModel struct {
	*wem.WorksiteEditModalModel
}

func NewReworkUpdateModalModel(vm *hvue.VM) *ReworkUpdateModalModel {
	rumm := &ReworkUpdateModalModel{WorksiteEditModalModel: wem.NewWorksiteEditModalModel(vm)}
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
		reworkupdate.RegisterComponent(),
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
