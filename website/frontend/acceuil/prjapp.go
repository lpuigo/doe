package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/worksitesummary"
	fm "github.com/lpuig/ewin/doe/website/frontend/frontmodel"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/goel/message"
	"honnef.co/go/js/xhr"
	"strconv"
)

//go:generate bash ./makejs.sh

func main() {
	mpm := NewMainPageModel()

	hvue.NewVM(
		hvue.El("#app"),
		hvue.Component("worksite-summary", worksitesummary.ComponentOptions()...),
		hvue.DataS(mpm),
		hvue.MethodsOf(mpm),
		hvue.Mounted(func(vm *hvue.VM) {
			mpm := &MainPageModel{Object: vm.Object}
			mpm.GetWorkSites()
		}),
	)
	js.Global.Get("Vue").Call("use", "ELEMENT.lang.en")

	// TODO to remove after debug
	js.Global.Set("mpm", mpm)
}

type MainPageModel struct {
	*js.Object

	VM *hvue.VM `js:"VM"`

	Worksites []*fm.Worksite `js:"Worksites"`
}

func NewMainPageModel() *MainPageModel {
	mpm := &MainPageModel{Object: tools.O()}
	mpm.Worksites = []*fm.Worksite{}
	return mpm
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Action Methods

func (m *MainPageModel) GetWorkSites() {
	go m.callGetWorkSites()
}

//func (m *MainPageModel) EditProject(p *fm.Project) {
//	m.EditedProject = p
//	m.VM.Refs("ProjectEdit").Call("Show", p)
//}
//
//func (m *MainPageModel) CreateNewProject() {
//	p := fm.NewProject()
//	p.Status = business.DefaultStatus()
//	p.Type = business.DefaultType()
//	p.Risk = business.DefaultRisk()
//	m.EditProject(p)
//}
//
//func (m *MainPageModel) ProcessEditedProject(p *fm.Project) {
//	m.EditedProject = p
//	if p.GetId >= 0 {
//		go m.callUpdatePrj(p)
//	} else {
//		go m.callCreatePrj(p)
//	}
//}
//
//func (m *MainPageModel) ProcessDeleteProject(p *fm.Project) {
//	m.EditedProject = p
//	if m.EditedProject.GetId >= 0 {
//		go m.callDeletePrj(m.EditedProject)
//	}
//}
//
//func (m *MainPageModel) ShowProjectStat(p *fm.Project) {
//	m.EditedProject = p
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
		message.SetDuration(tools.WarningMsgDuration)
		message.WarningStr(m.VM, "Something went wrong!\nServer returned code "+strconv.Itoa(req.Status), true)
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

//func (m *MainPageModel) callUpdatePrj(uprj *fm.Project) {
//	req := xhr.NewRequest("PUT", "/ptf/"+strconv.Itoa(uprj.GetId))
//	req.Timeout = tools.TimeOut
//	req.ResponseType = xhr.JSON
//	err := req.Send(json.Stringify(uprj))
//	if err != nil {
//		message.ErrorStr(m.VM, "Oups! "+err.Error(), true)
//		return
//	}
//	if req.Status == 200 {
//		uprj.Copy(fm.ProjectFromJS(req.Response))
//		uprj.SetAuditResult(m.auditer.Audit(uprj))
//	} else {
//		message.SetDuration(tools.WarningMsgDuration)
//		message.WarningStr(m.VM, "Something went wrong!\nServer returned code "+strconv.Itoa(req.Status), true)
//	}
//}
//
//func (m *MainPageModel) callCreatePrj(uprj *fm.Project) {
//	req := xhr.NewRequest("POST", "/ptf")
//	req.Timeout = tools.TimeOut
//	req.ResponseType = xhr.JSON
//	err := req.Send(json.Stringify(uprj))
//	if err != nil {
//		message.ErrorStr(m.VM, "Oups! "+err.Error(), true)
//		return
//	}
//	if req.Status == 201 {
//		m.EditedProject.Copy(fm.ProjectFromJS(req.Response))
//		m.EditedProject.SetAuditResult(m.auditer.Audit(m.EditedProject))
//		m.Projects = append(m.Projects, m.EditedProject)
//	} else {
//		message.SetDuration(tools.WarningMsgDuration)
//		message.WarningStr(m.VM, "Something went wrong!\nServer returned code "+strconv.Itoa(req.Status), true)
//	}
//}
//
//func (m *MainPageModel) callDeletePrj(dprj *fm.Project) {
//	req := xhr.NewRequest("DELETE", "/ptf/"+strconv.Itoa(dprj.GetId))
//	req.Timeout = tools.TimeOut
//	req.ResponseType = xhr.JSON
//	err := req.Send(nil)
//	if err != nil {
//		message.ErrorStr(m.VM, "Oups! "+err.Error(), true)
//		return
//	}
//	if req.Status == 200 {
//		m.deletePrj(dprj)
//		message.SetDuration(tools.SuccessMsgDuration)
//		message.SuccesStr(m.VM, "Project deleted !", true)
//	} else {
//		message.SetDuration(tools.WarningMsgDuration)
//		message.WarningStr(m.VM, "Something went wrong!\nServer returned code "+strconv.Itoa(req.Status), true)
//	}
//}
//
//func (m *MainPageModel) deletePrj(dprj *fm.Project) {
//	for i, p := range m.Projects {
//		if p.GetId == dprj.GetId {
//			m.EditedProject = nil
//			m.Projects = append(m.Projects[:i], m.Projects[i+1:]...)
//			break
//		}
//	}
//}
