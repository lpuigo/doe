package worksiteupdatemodal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	wem "github.com/lpuig/ewin/doe/website/frontend/comp/worksiteeditmodal"
	"github.com/lpuig/ewin/doe/website/frontend/comp/worksiteinfo"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
)

type WorksiteUpdateModalModel struct {
	*wem.WorksiteEditModalModel

	Filter string `js:"filter"`
}

func NewWorksiteUpdateModalModel(vm *hvue.VM) *WorksiteUpdateModalModel {
	wumm := &WorksiteUpdateModalModel{WorksiteEditModalModel: wem.NewWorksiteEditModalModel(vm)}
	wumm.Filter = ""
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

		hvue.Computed("filteredTroncons", func(vm *hvue.VM) interface{} {
			m := NewWorksiteUpdateModalModelFromJS(vm.Object)
			return m.GetTroncons()
		}),
		hvue.Computed("hasChanged", func(vm *hvue.VM) interface{} {
			m := NewWorksiteUpdateModalModelFromJS(vm.Object)
			return m.HasChanged()
		}),
		hvue.Computed("hasWarning", func(vm *hvue.VM) interface{} {
			//m := &WorksiteEditModalModel{Object: vm.Object}
			//if len(m.CurrentProject.Audits) > 0 {
			//	return "warning"
			//}
			return "success"
		}),
		hvue.Filter("FormatTronconRef", func(vm *hvue.VM, value *js.Object, args ...*js.Object) interface{} {
			m := NewWorksiteUpdateModalModelFromJS(vm.Object)
			t := &fm.Troncon{Object: value}
			return m.GetFormatTronconRef(t)
		}),
		hvue.Filter("FormatStatus", func(vm *hvue.VM, value *js.Object, args ...*js.Object) interface{} {
			m := NewWorksiteUpdateModalModelFromJS(vm.Object)
			t := &fm.Troncon{Object: value}
			return m.GetFormatStatus(t)
		}),
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Modal Methods

func (wumm *WorksiteUpdateModalModel) GetTroncons() []*fm.Troncon {
	res := []*fm.Troncon{}
	for _, o := range wumm.CurrentWorksite.Orders {
		for _, t := range o.Troncons {
			res = append(res, t)
		}
	}

	return res
}

func (wumm *WorksiteUpdateModalModel) TableRowClassName(rowInfo *js.Object) string {
	//wsi := &fm.WorksiteInfo{Object: rowInfo.Get("row")}
	return ""
}

func (wumm *WorksiteUpdateModalModel) GetFormatTronconRef(t *fm.Troncon) string {
	return t.Ref
}

func (wumm *WorksiteUpdateModalModel) GetFormatStatus(t *fm.Troncon) string {
	return "TODO"
}
