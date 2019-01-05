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
	case fm.WsStatusNew:
		return "info"
	case fm.WsStatusFormInProgress:
		return "warning"
	case fm.WsStatusInProgress:
		return "warning"
	case fm.WsStatusDOE:
		return ""
	case fm.WsStatusAttachment:
		return "success"
	case fm.WsStatusPayment:
		return "success"
	case fm.WsStatusDone:
		return "success"
	case fm.WsStatusRework:
		return "danger"
	}
	return "danger"
}

//func (wdm *WorksiteDetailModel) CheckDoeDate(vm *hvue.VM) {
//	wdm = &WorksiteDetailModel{Object: vm.Object}
//	if tools.Empty(wdm.Worksite.DoeDate) {
//		wdm.Worksite.Status = fm.WsStatusDOE
//		return
//	}
//	wdm.Worksite.Status = fm.WsStatusDone
//}

func (wdm *WorksiteDetailModel) CalcWorksiteStatus() string {
	// New if Worksite base info not completed
	ws := wdm.Worksite
	if !ws.IsDefined() {
		return fm.WsStatusNew
	}
	if !ws.IsFilledIn() {
		return fm.WsStatusFormInProgress
	}
	if !ws.OrdersCompleted() {
		ws.DoeDate = ""
		return fm.WsStatusInProgress
	}
	if ws.IsBlocked() {
		return fm.WsStatusDone
	}
	if tools.Empty(ws.DoeDate) {
		return fm.WsStatusDOE
	}
	if tools.Empty(ws.AttachmentDate) {
		return fm.WsStatusAttachment
	}
	if tools.Empty(ws.PaymentDate) {
		return fm.WsStatusPayment
	}
	return fm.WsStatusDone
}
