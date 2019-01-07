package worksiteeditmodal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/worksiteedit"
	"github.com/lpuig/ewin/doe/website/frontend/comp/worksiteinfo"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements/message"
)

type WorksiteEditModalModel struct {
	*js.Object

	Visible bool     `js:"visible"`
	VM      *hvue.VM `js:"VM"`

	//ActiveTabName string `js:"activeTabName"`

	EditedWorksite  *fm.Worksite `js:"edited_worksite"`
	CurrentWorksite *fm.Worksite `js:"current_worksite"`

	ShowConfirmDelete bool `js:"showconfirmdelete"`
}

func NewWorksiteEditModalModel(vm *hvue.VM) *WorksiteEditModalModel {
	wemm := &WorksiteEditModalModel{Object: tools.O()}
	wemm.Visible = false
	wemm.VM = vm

	//wemm.ActiveTabName = "project"

	wemm.EditedWorksite = fm.NewWorkSite()
	wemm.CurrentWorksite = fm.NewWorkSite()
	wemm.ShowConfirmDelete = false

	return wemm
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Component Methods

func Register() {
	hvue.NewComponent("worksite-edit-modal",
		ComponentOptions()...,
	)
}

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("worksite-edit-modal", ComponentOptions()...)
}

func ComponentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		worksiteedit.RegisterComponent(),
		worksiteinfo.RegisterComponentWorksite(),
		hvue.Template(template),
		hvue.Props("edited_worksite"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewWorksiteEditModalModel(vm)
		}),
		hvue.MethodsOf(&WorksiteEditModalModel{}),

		hvue.Computed("isNewWorksite", func(vm *hvue.VM) interface{} {
			m := &WorksiteEditModalModel{Object: vm.Object}
			return m.CurrentWorksite.Id == -1
		}),
		hvue.Computed("hasChanged", func(vm *hvue.VM) interface{} {
			m := &WorksiteEditModalModel{Object: vm.Object}
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

func (wemm *WorksiteEditModalModel) HasChanged() bool {
	if wemm.EditedWorksite.Object == nil {
		return true
	}
	return wemm.CurrentWorksite.SearchInString() != wemm.EditedWorksite.SearchInString()
}

func (wemm *WorksiteEditModalModel) Show(ws *fm.Worksite) {
	wemm.EditedWorksite = ws
	wemm.CurrentWorksite = ws.Clone()
	wemm.ShowConfirmDelete = false
	//wemm.ActiveTabName = "project" //force Project Tab active
	wemm.Visible = true
}

func (wemm *WorksiteEditModalModel) HideWithControl() {
	if wemm.HasChanged() {
		message.ConfirmWarning(wemm.VM, "OK pour perdre les changements effectués ?", wemm.Hide)
		return
	}
	wemm.Hide()
}

func (wemm *WorksiteEditModalModel) Hide() {
	wemm.Visible = false
	wemm.ShowConfirmDelete = false
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Action Button Methods

func (wemm *WorksiteEditModalModel) ConfirmChange() {
	wemm.EditedWorksite.Copy(wemm.CurrentWorksite)
	wemm.VM.Emit("update:edited_worksite", wemm.EditedWorksite)
	wemm.Hide()
}

func (wemm *WorksiteEditModalModel) UndoChange() {
	wemm.CurrentWorksite.Copy(wemm.EditedWorksite)
}

func (wemm *WorksiteEditModalModel) DeleteWorksite() {
	wemm.VM.Emit("delete:edited_worksite", wemm.EditedWorksite)
	wemm.Hide()
}

func (wemm *WorksiteEditModalModel) Duplicate() {
	wemm.EditedWorksite = wemm.CurrentWorksite
	wemm.CurrentWorksite.Ref += " (Copy)"
	wemm.CurrentWorksite.Id = -1
}
