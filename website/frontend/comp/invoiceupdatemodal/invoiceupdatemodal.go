package invoiceupdatemodal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/invoiceupdate"
	wem "github.com/lpuig/ewin/doe/website/frontend/comp/worksiteeditmodal"
	"github.com/lpuig/ewin/doe/website/frontend/comp/worksiteinfo"
)

type InvoiceUpdateModalModel struct {
	*wem.WorksiteEditModalModel
}

func NewInvoiceUpdateModalModel(vm *hvue.VM) *InvoiceUpdateModalModel {
	wumm := &InvoiceUpdateModalModel{WorksiteEditModalModel: wem.NewWorksiteEditModalModel(vm)}
	return wumm
}

func InvoiceUpdateModalModelFromJS(o *js.Object) *InvoiceUpdateModalModel {
	wumm := &InvoiceUpdateModalModel{WorksiteEditModalModel: &wem.WorksiteEditModalModel{Object: o}}
	return wumm
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Component Methods

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("invoice-update-modal", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		worksiteinfo.RegisterComponent(),
		invoiceupdate.RegisterComponent(),
		hvue.Template(template),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewInvoiceUpdateModalModel(vm)
		}),
		hvue.MethodsOf(&InvoiceUpdateModalModel{}),
		hvue.Computed("hasChanged", func(vm *hvue.VM) interface{} {
			m := InvoiceUpdateModalModelFromJS(vm.Object)
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
// Component Methods
