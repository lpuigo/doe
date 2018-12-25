package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
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
		hvue.Component("worksite-table", worksitetable.ComponentOptions()...),
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

	Worksites      []*fm.Worksite `js:"worksites"`
	Filter         string         `js:"filter"`
	EditedWorksite *fm.Worksite   `js:"editedWorksite"`
}

func NewMainPageModel() *MainPageModel {
	mpm := &MainPageModel{Object: tools.O()}
	mpm.Worksites = []*fm.Worksite{}
	mpm.Filter = ""
	return mpm
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Action Methods

func (m *MainPageModel) GetWorkSites() {
	go m.callGetWorkSites()
}

func (m *MainPageModel) EditWorksite(ws *fm.Worksite) {
	m.EditedWorksite = ws
	//m.VM.Refs("WorksiteEdit").Call("Show", ws)
	message.SetDuration(tools.WarningMsgDuration)
	message.InfoStr(m.VM, "Selected Worksite : "+ws.Ref, false)
}

//func (m *MainPageModel) CreateNewProject() {
//	p := fm.NewProject()
//	p.Status = business.DefaultStatus()
//	p.Type = business.DefaultType()
//	p.Risk = business.DefaultRisk()
//	m.EditProject(p)
//}
//
func (m *MainPageModel) ProcessEditedWorkSite(uws *fm.Worksite) {
	print("ProcessEditedWorkSite on", uws.Id, uws.Ref)
	if uws.Id >= 0 {
		go m.callUpdateWorksite(uws)
	} else {
		go m.callCreateWorksite(uws)
	}
}

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
		nws := fm.WorksiteFromJS(req.Response)
		uws.Copy(nws)
	} else {
		message.SetDuration(tools.WarningMsgDuration)
		msg := "Something went wrong!\n"
		msg += "Server returned code " + strconv.Itoa(req.Status) + "\n"
		message.ErrorMsgStr(m.VM, msg, req.Response, true)
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
		uws.Copy(fm.WorksiteFromJS(req.Response))
	} else {
		message.SetDuration(tools.WarningMsgDuration)
		msg := "Something went wrong!\n"
		msg += req.Response.String()
		msg += "Server returned code " + strconv.Itoa(req.Status) + "\n"
		message.ErrorStr(m.VM, msg, true)
	}
}

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
