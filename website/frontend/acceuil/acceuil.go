package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/userloginmodal"
	"github.com/lpuig/ewin/doe/website/frontend/comp/worksiteeditmodal"
	"github.com/lpuig/ewin/doe/website/frontend/comp/worksitetable"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements/message"
	"honnef.co/go/js/xhr"
	"strconv"
)

//go:generate bash ./makejs.sh

func main() {
	mpm := NewMainPageModel()

	hvue.NewVM(
		hvue.El("#app"),
		userloginmodal.RegisterComponent(),
		worksiteeditmodal.RegisterComponent(),
		worksitetable.RegisterComponent(),
		hvue.DataS(mpm),
		hvue.MethodsOf(mpm),
		hvue.Mounted(func(vm *hvue.VM) {
			mpm := &MainPageModel{Object: vm.Object}
			mpm.GetUserSession()
			mpm.GetWorkSiteInfos()
		}),
		hvue.Computed("LoggedUser", func(vm *hvue.VM) interface{} {
			mpm := &MainPageModel{Object: vm.Object}
			if mpm.User.Name == "" {
				return "Non connecté"
			}
			return mpm.User.Name
		}),
	)
	//js.Global.Get("Vue").Call("use", "ELEMENT.lang.fr")

	// TODO to remove after debug
	js.Global.Set("mpm", mpm)
}

type MainPageModel struct {
	*js.Object

	VM *hvue.VM `js:"VM"`

	User *fm.User `js:"User"`

	WorksiteInfos []*fm.WorksiteInfo `js:"worksiteInfos"`
	//EditedWorksite int      `js:"editedWorksite"`
}

func NewMainPageModel() *MainPageModel {
	mpm := &MainPageModel{Object: tools.O()}
	mpm.User = fm.NewUser()
	mpm.WorksiteInfos = []*fm.WorksiteInfo{}
	//mpm.EditedWorksite = -2
	return mpm
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Action Methods

func (m *MainPageModel) GetUserSession() {
	go m.callGetUser()
}

func (m *MainPageModel) ShowUserLogin() {
	m.VM.Refs("UserLoginModal").Call("Show", m.User)
}

func (m *MainPageModel) UserLogout() {
	go m.callLogout()
}

func (m *MainPageModel) GetWorkSiteInfos() {
	go m.callGetWorkSiteInfos()
}

func (m *MainPageModel) EditWorksite(id int) {
	//m.EditedWorksite = id
	m.VM.Refs("WorksiteEditModal").Call("Show", id)
}

func (m *MainPageModel) CreateNewWorksite() {
	m.EditWorksite(-1)
}

//TODO Move to WorksiteEditModal
//func (m *MainPageModel) ProcessEditedWorksite(uws *fm.Worksite) {
//	if uws.Id >= 0 {
//		go m.callUpdateWorksite(uws)
//	} else {
//		go m.callCreateWorksite(uws)
//	}
//}

//TODO Move to WorksiteEditModal
//func (m *MainPageModel) ProcessDeleteWorksite(uws *fm.Worksite) {
//	//m.EditedWorksite = uws
//	//if m.EditedWorksite.Id >= 0 {
//	//	go m.callDeleteWorksite(m.EditedWorksite)
//	//}
//}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// WS call Methods

func (m *MainPageModel) errorMessage(req *xhr.Request) {
	message.SetDuration(tools.WarningMsgDuration)
	msg := "Quelquechose c'est mal passé !\n"
	msg += "Le server retourne un code " + strconv.Itoa(req.Status) + "\n"
	message.ErrorMsgStr(m.VM, msg, req.Response, true)
}

func (m *MainPageModel) callGetUser() {
	req := xhr.NewRequest("GET", "/api/login")
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
	m.User.Copy(fm.NewUserFromJS(req.Response))
	if m.User.Name != "" {
		m.User.Connected = true
	}
}

func (m *MainPageModel) callLogout() {
	req := xhr.NewRequest("DELETE", "/api/login")
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
	m.User.Name = ""
	m.User.Connected = false
}

func (m *MainPageModel) callGetWorkSiteInfos() {
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
	wsis := []*fm.WorksiteInfo{}
	req.Response.Call("forEach", func(item *js.Object) {
		ws := fm.NewWorksiteInfoFromJs(item)
		wsis = append(wsis, ws)
	})
	m.WorksiteInfos = wsis
}

//TODO Move to WorksiteEditModal
//func (m *MainPageModel) callUpdateWorksite(uws *fm.Worksite) {
//	req := xhr.NewRequest("PUT", "/api/worksites/"+strconv.Itoa(uws.Id))
//	req.Timeout = tools.TimeOut
//	req.ResponseType = xhr.JSON
//	err := req.Send(json.Stringify(uws))
//	if err != nil {
//		message.ErrorStr(m.VM, "Oups! "+err.Error(), true)
//		return
//	}
//	if req.Status != tools.HttpOK {
//		m.errorMessage(req)
//		return
//	}
//	uws.Dirty = false
//	message.SuccesStr(m.VM, "Chantier sauvegardé")
//
//}

//TODO Move to WorksiteEditModal
//func (m *MainPageModel) callCreateWorksite(uws *fm.Worksite) {
//	req := xhr.NewRequest("POST", "/api/worksites")
//	req.Timeout = tools.TimeOut
//	req.ResponseType = xhr.JSON
//	err := req.Send(json.Stringify(uws))
//	if err != nil {
//		message.ErrorStr(m.VM, "Oups! "+err.Error(), true)
//		return
//	}
//	if req.Status != tools.HttpCreated {
//		m.errorMessage(req)
//	}
//	uws.Dirty = false
//	uws.Copy(fm.WorksiteFromJS(req.Response))
//	message.SetDuration(tools.SuccessMsgDuration)
//	message.SuccesStr(m.VM, "Nouveau chantier sauvegardé")
//}

//TODO Move to WorksiteEditModal
//func (m *MainPageModel) callDeleteWorksite(dws *fm.Worksite) {
//	req := xhr.NewRequest("DELETE", "/api/worksites/"+strconv.Itoa(dws.Id))
//	req.Timeout = tools.TimeOut
//	req.ResponseType = xhr.JSON
//	err := req.Send(nil)
//	if err != nil {
//		message.ErrorStr(m.VM, "Oups! "+err.Error(), true)
//		return
//	}
//	if req.Status != tools.HttpOK {
//		m.errorMessage(req)
//	}
//	message.SetDuration(tools.SuccessMsgDuration)
//	message.SuccesStr(m.VM, "Chantier supprimé !")
//}
