package worksitedetail

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/ptedit"
	"github.com/lpuig/ewin/doe/website/frontend/comp/tronconedit"
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
		hvue.Component("troncon-edit", tronconedit.ComponentOptions()...),
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
