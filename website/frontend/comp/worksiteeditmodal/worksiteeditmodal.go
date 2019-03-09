package worksiteeditmodal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/worksiteedit"
	"github.com/lpuig/ewin/doe/website/frontend/comp/worksiteinfo"
	"github.com/lpuig/ewin/doe/website/frontend/comp/worksiteupdate"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements/message"
	"github.com/lpuig/ewin/doe/website/frontend/tools/json"
	"honnef.co/go/js/xhr"
	"strconv"
)

type WorksiteEditModalModel struct {
	*js.Object

	Visible bool     `js:"visible"`
	VM      *hvue.VM `js:"VM"`

	ActiveTabName string `js:"activeTabName"`

	User            *fm.User     `js:"user"`
	EditedWorksite  *fm.Worksite `js:"edited_worksite"`
	CurrentWorksite *fm.Worksite `js:"current_worksite"`

	Loading           bool `js:"loading"`
	Saving            bool `js:"saving"`
	ShowConfirmDelete bool `js:"showconfirmdelete"`
}

func NewWorksiteEditModalModel(vm *hvue.VM) *WorksiteEditModalModel {
	wemm := &WorksiteEditModalModel{Object: tools.O()}
	wemm.Visible = false
	wemm.VM = vm

	wemm.ActiveTabName = "Create"

	wemm.User = fm.NewUser()
	wemm.EditedWorksite = fm.NewWorkSite()
	wemm.CurrentWorksite = fm.NewWorkSite()
	wemm.Loading = false
	wemm.Saving = false
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
		worksiteupdate.RegisterComponent(),
		worksiteinfo.RegisterComponent(),
		hvue.Template(template),
		//hvue.Props("edited_worksite"),
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
	if wemm.EditedWorksite.Object == js.Undefined {
		return true
	}
	return wemm.CurrentWorksite.SearchInString() != wemm.EditedWorksite.SearchInString()
}

func (wemm *WorksiteEditModalModel) Show(id int, user *fm.User) {
	wemm.EditedWorksite = fm.NewWorkSite()
	if id < 0 {
		wemm.EditedWorksite.AddOrder()
	}
	wemm.CurrentWorksite = wemm.EditedWorksite.Clone()
	wemm.User = user
	wemm.SetActiveTab()
	wemm.ShowConfirmDelete = false
	if id >= 0 {
		wemm.Loading = true
		go wemm.callGetWorksite(id)
	}
	//wemm.ActiveTabName = "project" //force Project Tab active
	wemm.Visible = true
}

func (wemm *WorksiteEditModalModel) SetActiveTab() {
	wemm.ActiveTabName = "Update"
	if wemm.User.Permissions["Create"] {
		wemm.ActiveTabName = "Create"
	}
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

func (wemm *WorksiteEditModalModel) Attachment() string {
	url := "/api/worksites/" + strconv.Itoa(wemm.CurrentWorksite.Id) + "/attach"
	return url
}

func (wemm *WorksiteEditModalModel) ConfirmChange() {
	wemm.Saving = true
	if wemm.CurrentWorksite.Id >= 0 {
		go wemm.callUpdateWorksite(wemm.CurrentWorksite)
	} else {
		go wemm.callCreateWorksite(wemm.CurrentWorksite)
	}
}

func (wemm *WorksiteEditModalModel) UndoChange() {
	wemm.CurrentWorksite.Copy(wemm.EditedWorksite)
}

func (wemm *WorksiteEditModalModel) DeleteWorksite() {
	wemm.Saving = true
	go wemm.callDeleteWorksite(wemm.CurrentWorksite)
}

func (wemm *WorksiteEditModalModel) Duplicate() {
	wemm.EditedWorksite = wemm.CurrentWorksite
	wemm.CurrentWorksite.Ref += " (Copy)"
	wemm.CurrentWorksite.Id = -1
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// WS call Methods

func (wemm *WorksiteEditModalModel) errorMessage(req *xhr.Request) {
	message.SetDuration(tools.WarningMsgDuration)
	msg := "Quelquechose c'est mal passé !\n"
	msg += "Le server retourne un code " + strconv.Itoa(req.Status) + "\n"
	message.ErrorMsgStr(wemm.VM, msg, req.Response, true)
}

func (wemm *WorksiteEditModalModel) callGetWorksite(id int) {
	defer func() { wemm.Loading = false }()
	req := xhr.NewRequest("GET", "/api/worksites/"+strconv.Itoa(id))
	req.Timeout = tools.TimeOut
	req.ResponseType = xhr.JSON
	err := req.Send(nil)
	if err != nil {
		message.ErrorStr(wemm.VM, "Oups! "+err.Error(), true)
		return
	}
	if req.Status != tools.HttpOK {
		wemm.errorMessage(req)
		wemm.Hide()
		return
	}
	wemm.EditedWorksite = fm.WorksiteFromJS(req.Response)
	wemm.CurrentWorksite.Copy(wemm.EditedWorksite)
	return
}

func (wemm *WorksiteEditModalModel) callUpdateWorksite(uws *fm.Worksite) {
	defer func() { wemm.Saving = false }()
	req := xhr.NewRequest("PUT", "/api/worksites/"+strconv.Itoa(uws.Id))
	req.Timeout = tools.TimeOut
	req.ResponseType = xhr.JSON
	err := req.Send(json.Stringify(uws))
	if err != nil {
		message.ErrorStr(wemm.VM, "Oups! "+err.Error(), true)
		return
	}
	if req.Status != tools.HttpOK {
		wemm.errorMessage(req)
		return
	}
	wemm.VM.Emit("update_worksite")
	message.SuccesStr(wemm.VM, "Chantier sauvegardé")
	wemm.Hide()
}

func (wemm *WorksiteEditModalModel) callCreateWorksite(uws *fm.Worksite) {
	defer func() { wemm.Saving = false }()
	req := xhr.NewRequest("POST", "/api/worksites")
	req.Timeout = tools.TimeOut
	req.ResponseType = xhr.JSON
	err := req.Send(json.Stringify(uws))
	if err != nil {
		message.ErrorStr(wemm.VM, "Oups! "+err.Error(), true)
		return
	}
	if req.Status != tools.HttpCreated {
		wemm.errorMessage(req)
		return
	}
	wemm.VM.Emit("update_worksite")
	message.SuccesStr(wemm.VM, "Nouveau chantier sauvegardé")
	wemm.Hide()
}

func (wemm *WorksiteEditModalModel) callDeleteWorksite(dws *fm.Worksite) {
	defer func() { wemm.Saving = false }()
	req := xhr.NewRequest("DELETE", "/api/worksites/"+strconv.Itoa(dws.Id))
	req.Timeout = tools.TimeOut
	req.ResponseType = xhr.JSON
	err := req.Send(nil)
	if err != nil {
		message.ErrorStr(wemm.VM, "Oups! "+err.Error(), true)
		return
	}
	if req.Status != tools.HttpOK {
		wemm.errorMessage(req)
	}
	wemm.VM.Emit("update_worksite")
	message.SuccesStr(wemm.VM, "Chantier supprimé !")
	wemm.Hide()
}
