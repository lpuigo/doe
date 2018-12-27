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
}

func (m *MainPageModel) CreateNewWorksite() {
	ws := fm.NewWorkSite()
	ws.AddOrder()
	m.EditWorksite(ws)
}

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
	if req.Status != 200 {
		m.errorMessage(req)
		return
	}
	uws.Dirty = false
	message.SuccesStr(m.VM, "Chantier sauvegardé")

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
	if req.Status != 201 {
		m.errorMessage(req)
	}
	uws.Dirty = false
	uws.Copy(fm.WorksiteFromJS(req.Response))
	m.Worksites = append(m.Worksites, uws)
	message.SetDuration(tools.SuccessMsgDuration)
	message.SuccesStr(m.VM, "Nouveau chantier sauvegardé")
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
	if req.Status != 200 {
		m.errorMessage(req)
	}
	m.deletePrj(dws)
	message.SetDuration(tools.SuccessMsgDuration)
	message.SuccesStr(m.VM, "Chantier supprimé !")
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
