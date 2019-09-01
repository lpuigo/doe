package reworkeditmodal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/reworkedit"
	"github.com/lpuig/ewin/doe/website/frontend/comp/reworkupdate"
	wem "github.com/lpuig/ewin/doe/website/frontend/comp/worksiteeditmodal"
	"github.com/lpuig/ewin/doe/website/frontend/comp/worksiteinfo"
	"github.com/lpuig/ewin/doe/website/frontend/model/worksite"
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
		reworkedit.RegisterComponent(),
		reworkupdate.RegisterComponent(),
		hvue.Template(template),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewReworkEditModalModel(vm)
		}),
		hvue.MethodsOf(&ReworkEditModalModel{}),
		hvue.Computed("HasRework", func(vm *hvue.VM) interface{} {
			m := NewReworkEditModalModelFromJS(vm.Object)
			if m.Loading || !worksite.WorksiteIsReworkable(m.CurrentWorksite.Status) {
				return false
			}
			if m.CurrentWorksite.Rework != nil && m.CurrentWorksite.Rework.Object != js.Undefined {
				return true
			}
			m.CurrentWorksite.Rework = worksite.NewRework()
			return true
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
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Modal Methods
