package invoiceupdate

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/tronconstatustag"
	"github.com/lpuig/ewin/doe/website/frontend/comp/worksitestatustag"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Registration

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("invoice-update", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		tronconstatustag.RegisterComponent(),
		worksitestatustag.RegisterComponent(),
		hvue.Template(template),
		hvue.Props("worksite", "user"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewInvoiceUpdateModel(vm)
		}),
		hvue.MethodsOf(&InvoiceUpdateModel{}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type InvoiceUpdateModel struct {
	*js.Object

	Worksite          *fm.Worksite `js:"worksite"`
	ReferenceWorksite *fm.Worksite `js:"refWorksite"`
	User              *fm.User     `js:"user"`
	Filter            string       `js:"filter"`

	VM *hvue.VM `js:"VM"`
}

func NewInvoiceUpdateModel(vm *hvue.VM) *InvoiceUpdateModel {
	wum := &InvoiceUpdateModel{Object: tools.O()}
	wum.VM = vm
	wum.Worksite = nil
	wum.ReferenceWorksite = nil
	wum.User = nil
	wum.Filter = ""
	return wum
}

func InvoiceUpdateModelFromJS(o *js.Object) *InvoiceUpdateModel {
	return &InvoiceUpdateModel{Object: o}
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Component Methods

func (ium *InvoiceUpdateModel) IsDisabled(vm *hvue.VM, info string) bool {
	ium = &InvoiceUpdateModel{Object: vm.Object}
	return ium.Worksite.IsInfoDisabled(info)
}
