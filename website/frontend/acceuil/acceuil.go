package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/worksiteeditmodal"
	"github.com/lpuig/ewin/doe/website/frontend/comp/worksitetable"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements/message"
	"github.com/lpuig/ewin/doe/website/frontend/tools/json"
	"honnef.co/go/js/xhr"
	"strconv"
)

//go:generate bash ./makejs.sh

func main() {
	mpm := NewMainPageModel()

	hvue.NewVM(
		hvue.El("#app"),
		worksiteeditmodal.RegisterComponent(),
		worksitetable.RegisterComponent(),
		hvue.DataS(mpm),
		hvue.MethodsOf(mpm),
		hvue.Mounted(func(vm *hvue.VM) {
			mpm := &MainPageModel{Object: vm.Object}
			mpm.GetWorkSites()
		}),
	)
	//js.Global.Get("Vue").Call("use", "ELEMENT.lang.fr")

	// TODO to remove after debug
	js.Global.Set("mpm", mpm)
}

type MainPageModel struct {
	*js.Object

	VM *hvue.VM `js:"VM"`

	Worksites      []*fm.Worksite `js:"worksites"`
	EditedWorksite *fm.Worksite   `js:"editedWorksite"`
}

func NewMainPageModel() *MainPageModel {
	mpm := &MainPageModel{Object: tools.O()}
	mpm.Worksites = []*fm.Worksite{}
	mpm.EditedWorksite = nil
	return mpm
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Action Methods

func (m *MainPageModel) GetWorkSites() {
	go m.callGetWorkSites()
}

func (m *MainPageModel) EditWorksite(ws *fm.Worksite) {
	print("current_worksite", ws.Object)
	m.EditedWorksite = ws
	m.VM.Refs("WorksiteEditModal").Call("Show", ws)
	//message.SetDuration(tools.WarningMsgDuration)
	//message.InfoStr(m.VM, "Selected Worksite : "+ws.Ref, false)
}

//func (m *MainPageModel) CreateNewProject() {
//	p := fm.NewProject()
//	p.Status = business.DefaultStatus()
//	p.Type = business.DefaultType()
//	p.Risk = business.DefaultRisk()
//	m.EditProject(p)
//}
//
func (m *MainPageModel) ProcessEditedWorksite(uws *fm.Worksite) {
	print("ProcessEditedWorkSite on", uws.Id, uws.Ref)
	if uws.Id >= 0 {
		go m.callUpdateWorksite(uws)
	} else {
		go m.callCreateWorksite(uws)
	}
}

func (m *MainPageModel) ProcessDeleteWorksite(uws *fm.Worksite) {
	m.EditedWorksite = uws
	if m.EditedWorksite.Id >= 0 {
		go m.callDeleteWorksite(m.EditedWorksite)
	}
}

//func (m *MainPageModel) ShowProjectStat(p *fm.Project) {
//	m.EditedWorksite = p
//	m.VM.Refs("ProjectStat").Call("Show", p)
//}
//
//func (m *MainPageModel) ShowProjectAudit(p *fm.Project) {
//	infos := "Audit for " + p.Client + " - " + p.Name + ":\n"
//	for _, a := range p.Audits {
//		infos += a.Priority + " " + a.Title + "\n"
//	}
//	message.SetDuration(tools.WarningMsgDuration)
//	message.InfoStr(m.VM, infos, false)
//}
//
//func (m *MainPageModel) ShowJiraStat() {
//	m.VM.Refs("JiraStat").Call("Show")
//}
//
//func (m *MainPageModel) ShowWorkloadSchedule() {
//	m.VM.Refs("WorkloadSchedule").Call("Show")
//}
//
//func (m *MainPageModel) ShowTimeLine() {
//	m.VM.Refs("TimeLine").Call("Show", timeline_modal.NewInfos(m.Projects))
//}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// WS call Methods

func (m *MainPageModel) errorMessage(req *xhr.Request) {
	message.SetDuration(tools.WarningMsgDuration)
	msg := "Quelquechose c'est mal passé !\n"
	msg += "Le server retourne un code " + strconv.Itoa(req.Status) + "\n"
	message.ErrorMsgStr(m.VM, msg, req.Response, true)
}

func (m *MainPageModel) callGetWorkSites() {
	req := xhr.NewRequest("GET", "/api/worksites")
	req.Timeout = tools.LongTimeOut
	req.ResponseType = xhr.JSON
	//m.DispPrj = false
	err := req.Send(nil)
	if err != nil {
		message.ErrorStr(m.VM, "Oups! "+err.Error(), true)
		return
	}
	if req.Status != tools.HttpOK {
		m.errorMessage(req)
		return
	}
	worksites := []*fm.Worksite{}
	req.Response.Call("forEach", func(item *js.Object) {
		ws := fm.WorksiteFromJS(item)
		//p.SetAuditResult(m.auditer.Audit(p))
		worksites = append(worksites, ws)
	})
	m.Worksites = worksites
}

func (m *MainPageModel) callUpdateWorksite(uws *fm.Worksite) {
	req := xhr.NewRequest("PUT", "/api/worksites/"+strconv.Itoa(uws.Id))
	req.Timeout = tools.TimeOut
	req.ResponseType = xhr.JSON
	err := req.Send(json.Stringify(uws))
	if err != nil {
		message.ErrorStr(m.VM, "Oups! "+err.Error(), true)
		return
	}
	if req.Status == 200 {
		uws.Dirty = false
		message.SuccesStr(m.VM, "Chantier sauvegardé")
	} else {
		m.errorMessage(req)
	}
}

func (m *MainPageModel) callCreateWorksite(uws *fm.Worksite) {
	req := xhr.NewRequest("POST", "/api/worksites")
	req.Timeout = tools.TimeOut
	req.ResponseType = xhr.JSON
	err := req.Send(json.Stringify(uws))
	if err != nil {
		message.ErrorStr(m.VM, "Oups! "+err.Error(), true)
		return
	}
	if req.Status == 201 {
		uws.Dirty = false
		uws.Copy(fm.WorksiteFromJS(req.Response))
		message.SetDuration(tools.SuccessMsgDuration)
		message.SuccesStr(m.VM, "Nouveau chantier sauvegardé")
	} else {
		m.errorMessage(req)
	}
}

func (m *MainPageModel) callDeleteWorksite(dws *fm.Worksite) {
	req := xhr.NewRequest("DELETE", "/api/worksites/"+strconv.Itoa(dws.Id))
	req.Timeout = tools.TimeOut
	req.ResponseType = xhr.JSON
	err := req.Send(nil)
	if err != nil {
		message.ErrorStr(m.VM, "Oups! "+err.Error(), true)
		return
	}
	if req.Status == 200 {
		m.deletePrj(dws)
		message.SetDuration(tools.SuccessMsgDuration)
		message.SuccesStr(m.VM, "Chantier supprimé !")
	} else {
		m.errorMessage(req)
	}
}

func (m *MainPageModel) deletePrj(dws *fm.Worksite) {
	for i, ws := range m.Worksites {
		if ws.Id == dws.Id {
			m.EditedWorksite = nil
			m.Worksites = append(m.Worksites[:i], m.Worksites[i+1:]...)
			break
		}
	}
}
