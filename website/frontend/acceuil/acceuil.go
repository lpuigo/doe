package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/userloginmodal"
	"github.com/lpuig/ewin/doe/website/frontend/comp/worksiteeditmodal"
	"github.com/lpuig/ewin/doe/website/frontend/comp/worksitetable"
	"github.com/lpuig/ewin/doe/website/frontend/comp/worksiteupdatemodal"
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
		worksiteupdatemodal.RegisterComponent(),
		worksitetable.RegisterComponent(),
		hvue.DataS(mpm),
		hvue.MethodsOf(mpm),
		hvue.Mounted(func(vm *hvue.VM) {
			mpm := &MainPageModel{Object: vm.Object}
			mpm.GetUserSession()
		}),
		hvue.Computed("LoggedUser", func(vm *hvue.VM) interface{} {
			mpm := &MainPageModel{Object: vm.Object}
			if mpm.User.Name == "" {
				return "Non connecté"
			}
			return mpm.User.Name
		}),
		hvue.Computed("UpdatableWorksiteInfos", func(vm *hvue.VM) interface{} {
			mpm := &MainPageModel{Object: vm.Object}
			return mpm.GetUpdatableWorsiteInfos()
		}),
		hvue.Computed("ReworkWorksiteInfos", func(vm *hvue.VM) interface{} {
			mpm := &MainPageModel{Object: vm.Object}
			return mpm.GetReworkWorksiteInfos()
		}),
		hvue.Computed("NbUpdate", func(vm *hvue.VM) interface{} {
			mpm := &MainPageModel{Object: vm.Object}
			return mpm.GetUpdatableWorsiteNb()
		}),
		hvue.Computed("NbRework", func(vm *hvue.VM) interface{} {
			mpm := &MainPageModel{Object: vm.Object}
			return mpm.GetReworkWorsiteNb()
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

	ActiveMode string `js:"ActiveMode"`

	WorksiteInfos []*fm.WorksiteInfo `js:"worksiteInfos"`
	//EditedWorksite int      `js:"editedWorksite"`
}

func NewMainPageModel() *MainPageModel {
	mpm := &MainPageModel{Object: tools.O()}
	mpm.User = fm.NewUser()
	mpm.ActiveMode = "Update"
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

func (m *MainPageModel) UpdateWorksite(id int) {
	//m.EditedWorksite = id
	m.VM.Refs("WorksiteUpdateModal").Call("Show", id)
}

func (m *MainPageModel) CreateNewWorksite() {
	m.EditWorksite(-1)
}

func (m *MainPageModel) GetUpdatableWorsiteInfos() []*fm.WorksiteInfo {
	res := []*fm.WorksiteInfo{}
	for _, wsi := range m.WorksiteInfos {
		if fm.WorksiteIsUpdatable(wsi.Status) {
			res = append(res, wsi)
		}
	}
	return res
}

func (m *MainPageModel) GetUpdatableWorsiteNb() int {
	res := 0
	for _, wsi := range m.WorksiteInfos {
		if fm.WorksiteIsUpdatable(wsi.Status) {
			res += 1
		}
	}
	return res

}

func (m *MainPageModel) GetReworkWorksiteInfos() []*fm.WorksiteInfo {
	res := []*fm.WorksiteInfo{}
	for _, wsi := range m.WorksiteInfos {
		if fm.WorksiteHasRework(wsi.Status) {
			res = append(res, wsi)
		}
	}
	return res
}

func (m *MainPageModel) GetReworkWorsiteNb() int {
	res := 0
	for _, wsi := range m.WorksiteInfos {
		if fm.WorksiteHasRework(wsi.Status) {
			res += 1
		}
	}
	return res

}

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
		m.GetWorkSiteInfos()
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
	m.WorksiteInfos = []*fm.WorksiteInfo{}
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
