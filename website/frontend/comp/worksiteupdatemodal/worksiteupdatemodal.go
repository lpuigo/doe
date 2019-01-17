package worksiteupdatemodal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	wem "github.com/lpuig/ewin/doe/website/frontend/comp/worksiteeditmodal"
	"github.com/lpuig/ewin/doe/website/frontend/comp/worksiteinfo"
)

type WorksiteUpdateModalModel struct {
	*wem.WorksiteEditModalModel
}

func NewWorksiteUpdateModalModel(vm *hvue.VM) *WorksiteUpdateModalModel {
	wumm := &WorksiteUpdateModalModel{WorksiteEditModalModel: wem.NewWorksiteEditModalModel(vm)}
	return wumm
}

func NewWorksiteUpdateModalModelFromJS(o *js.Object) *WorksiteUpdateModalModel {
	wumm := &WorksiteUpdateModalModel{WorksiteEditModalModel: &wem.WorksiteEditModalModel{Object: o}}
	return wumm
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Component Methods

func Register() {
	hvue.NewComponent("worksite-update-modal",
		ComponentOptions()...,
	)
}

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("worksite-update-modal", ComponentOptions()...)
}

func ComponentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		worksiteinfo.RegisterComponent(),
		hvue.Template(template),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewWorksiteUpdateModalModel(vm)
		}),
		hvue.MethodsOf(&WorksiteUpdateModalModel{}),

		//hvue.Computed("isNewWorksite", func(vm *hvue.VM) interface{} {
		//	m := &WorksiteEditModalModel{Object: vm.Object}
		//	return m.CurrentWorksite.Id == -1
		//}),
		hvue.Computed("hasChanged", func(vm *hvue.VM) interface{} {
			m := NewWorksiteUpdateModalModel(vm)
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
