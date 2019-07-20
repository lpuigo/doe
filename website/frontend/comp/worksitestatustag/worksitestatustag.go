package worksitestatustag

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/ptedit"
	"github.com/lpuig/ewin/doe/website/frontend/model/worksite"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

const (
	template string = `
<el-tag :type="Status" size="medium" style="width: 100%;text-align: left">{{StatusText}}</el-tag>
`
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Registration

func Register() {
	hvue.NewComponent("worksite-status-tag",
		ComponentOptions()...,
	)
}

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("worksite-status-tag", ComponentOptions()...)
}

func ComponentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		ptedit.RegisterComponent(),
		hvue.Template(template),
		hvue.Props("value"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewWorksiteStatusTagModel(vm)
		}),
		hvue.Computed("Status", func(vm *hvue.VM) interface{} {
			wst := &WorksiteStatusTagModel{Object: vm.Object}
			wst.Worksite.UpdateStatus()
			statusType, statusText := wst.SetStatus()
			wst.StatusText = statusText
			return statusType
		}),
		hvue.MethodsOf(&WorksiteStatusTagModel{}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type WorksiteStatusTagModel struct {
	*js.Object

	Worksite   *worksite.Worksite `js:"value"`
	StatusText string             `js:"StatusText"`

	VM *hvue.VM `js:"VM"`
}

func NewWorksiteStatusTagModel(vm *hvue.VM) *WorksiteStatusTagModel {
	tem := &WorksiteStatusTagModel{Object: tools.O()}
	tem.VM = vm
	tem.Worksite = nil
	tem.StatusText = ""
	return tem
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Actions

func (wst *WorksiteStatusTagModel) SetStatus() (statusType, statusText string) {
	statusText = wst.Worksite.WorksiteStatusLabel()
	switch wst.Worksite.Status {
	case worksite.WsStatusNew:
		statusType = "info"
	case worksite.WsStatusFormInProgress:
		statusType = "warning"
	case worksite.WsStatusInProgress:
		statusType = "warning"
	case worksite.WsStatusDOE:
		statusType = ""
	case worksite.WsStatusAttachment:
		statusType = "success"
	case worksite.WsStatusPayment:
		statusType = "success"
	case worksite.WsStatusCancelled:
		statusType = "info"
	case worksite.WsStatusDone:
		statusType = "success"
	default:
		statusType = "danger"
	}
	return
}
