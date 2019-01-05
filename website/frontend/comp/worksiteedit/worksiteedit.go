package worksiteedit

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
	hvue.NewComponent("worksite-edit",
		ComponentOptions()...,
	)
}

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("worksite-edit", ComponentOptions()...)
}

func ComponentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		ptedit.RegisterComponent(),
		orderedit.RegisterComponent(),
		hvue.Template(template),
		hvue.Props("worksite", "readonly"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewWorksiteDetailModel(vm)
		}),
		hvue.Computed("HasChanged", func(vm *hvue.VM) interface{} {
			wdm := &WorksiteDetailModel{Object: vm.Object}
			if wdm.ReferenceWorksite.Object == nil {
				wdm.ReferenceWorksite = wdm.Worksite.Clone()
				return wdm.Worksite.Dirty
			}
			s1 := wdm.Worksite.SearchInString()
			s2 := wdm.ReferenceWorksite.SearchInString()
			wdm.Worksite.Dirty = s1 != s2
			return wdm.Worksite.Dirty
		}),
		hvue.Computed("StatusType", func(vm *hvue.VM) interface{} {
			wdm := &WorksiteDetailModel{Object: vm.Object}
			wdm.Worksite.Status = wdm.CalcWorksiteStatus()
			return wdm.WorksiteStatusType()
		}),
		hvue.Filter("FormatStatus", func(vm *hvue.VM, value *js.Object, args ...*js.Object) interface{} {
			status := value.String()
			return fm.WorksiteStatusLabel(status)
		}),
		hvue.MethodsOf(&WorksiteDetailModel{}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type WorksiteDetailModel struct {
	*js.Object

	Worksite          *fm.Worksite `js:"worksite"`
	ReferenceWorksite *fm.Worksite `js:"refWorksite"`
	ReadOnly          bool         `js:"readonly"`

	VM *hvue.VM `js:"VM"`
}

func NewWorksiteDetailModel(vm *hvue.VM) *WorksiteDetailModel {
	wdm := &WorksiteDetailModel{Object: tools.O()}
	wdm.VM = vm
	wdm.Worksite = nil
	wdm.ReferenceWorksite = nil
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

func (wdm *WorksiteDetailModel) Undo(vm *hvue.VM) {
	wdm = &WorksiteDetailModel{Object: vm.Object}
	wdm.Worksite.Copy(wdm.ReferenceWorksite)
}

//func (wdm *WorksiteDetailModel) WorksiteStatusValTexts() []*elements.ValText {
//	res := []*elements.ValText{}
//	for _, v := range []string{"New", "InProgress", "DOE", "Done", "Rework"} {
//		res = append(res, elements.NewValText(v, fm.WorksiteStatusLabel(v)))
//	}
//	return res
//}

func (wdm *WorksiteDetailModel) WorksiteStatusType() string {
	switch wdm.Worksite.Status {
	case "00 New":
		return "info"
	case "10 FormInProgress":
		return "warning"
	case "20 InProgress":
		return "warning"
	case "30 DOE":
		return ""
	case "40 Attachment":
		return "success"
	case "50 Payment":
		return "success"
	case "99 Done":
		return "success"
	case "80 Rework":
		return "danger"
	}
	return "danger"
}

//func (wdm *WorksiteDetailModel) CheckDoeDate(vm *hvue.VM) {
//	wdm = &WorksiteDetailModel{Object: vm.Object}
//	if tools.Empty(wdm.Worksite.DoeDate) {
//		wdm.Worksite.Status = "30 DOE"
//		return
//	}
//	wdm.Worksite.Status = "99 Done"
//}

func (wdm *WorksiteDetailModel) CalcWorksiteStatus() string {
	// New if Worksite base info not completed
	ws := wdm.Worksite
	if !ws.IsDefined() {
		return "00 New"
	}
	if !ws.IsFilledIn() {
		return "10 FormInProgress"
	}
	if !ws.OrdersCompleted() {
		ws.DoeDate = ""
		return "20 InProgress"
	}
	if ws.IsBlocked() {
		return "99 Done"
	}
	if tools.Empty(ws.DoeDate) {
		return "30 DOE"
	}
	if tools.Empty(ws.AttachmentDate) {
		return "40 Attachment"
	}
	if tools.Empty(ws.PaymentDate) {
		return "50 Payment"
	}
	return "99 Done"
}
