package worksitedetail

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/orderedit"
	"github.com/lpuig/ewin/doe/website/frontend/comp/ptedit"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Registration

func Register() {
	hvue.NewComponent("worksite-detail",
		ComponentOptions()...,
	)
}

func ComponentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Component("pt-edit", ptedit.ComponentOptions()...),
		hvue.Component("order-edit", orderedit.ComponentOptions()...),
		hvue.Template(template),
		hvue.Props("worksite", "readonly"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewWorksiteDetailModel(vm)
		}),
		hvue.MethodsOf(&WorksiteDetailModel{}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type WorksiteDetailModel struct {
	*js.Object

	Worksite       *fm.Worksite `js:"worksite"`
	EditedWorksite *fm.Worksite `js:"editedWorksite"`
	ReadOnly       bool         `js:"readonly"`

	VM *hvue.VM `js:"VM"`
}

func NewWorksiteDetailModel(vm *hvue.VM) *WorksiteDetailModel {
	wdm := &WorksiteDetailModel{Object: tools.O()}
	wdm.VM = vm
	wdm.Worksite = nil
	wdm.EditedWorksite = nil
	wdm.ReadOnly = false
	return wdm
}

func (wdm *WorksiteDetailModel) DeleteOrder(vm *hvue.VM, i int) {
	wdm = &WorksiteDetailModel{Object: vm.Object}
	wdm.Worksite.DeleteOrder(i)
}

func (wdm *WorksiteDetailModel) AddOrder(vm *hvue.VM) {
	wdm = &WorksiteDetailModel{Object: vm.Object}
	wdm.Worksite.AddOrder()
}

func (wdm *WorksiteDetailModel) Save(vm *hvue.VM) {
	wdm = &WorksiteDetailModel{Object: vm.Object}
	vm.Emit("save_worksite", wdm.Worksite)
}
