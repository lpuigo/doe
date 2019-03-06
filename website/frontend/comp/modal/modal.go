package modal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

type ModalModel struct {
	*js.Object

	Visible bool     `js:"visible"`
	VM      *hvue.VM `js:"VM"`
	Loading bool     `js:"loading"`
}

func NewModalModel(vm *hvue.VM) *ModalModel {
	mm := &ModalModel{Object: tools.O()}
	mm.Visible = false
	mm.VM = vm

	mm.Loading = false

	return mm
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Component Methods

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("modal", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewModalModel(vm)
		}),
		hvue.MethodsOf(&ModalModel{}),
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Modal Methods

func (mm *ModalModel) Show() {
	mm.Visible = true
}

func (mm *ModalModel) Hide() {
	mm.Visible = false
}
