package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/adminmodal"
	"github.com/lpuig/ewin/doe/website/frontend/comp/invoicetable"
	"github.com/lpuig/ewin/doe/website/frontend/comp/invoiceupdatemodal"
	"github.com/lpuig/ewin/doe/website/frontend/comp/reworkeditmodal"
	"github.com/lpuig/ewin/doe/website/frontend/comp/reworkupdatemodal"
	"github.com/lpuig/ewin/doe/website/frontend/comp/ripsitetable"
	"github.com/lpuig/ewin/doe/website/frontend/comp/ripsiteupdatemodal"
	"github.com/lpuig/ewin/doe/website/frontend/comp/teamproductivitymodal"
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
		reworkeditmodal.RegisterComponent(),
		reworkupdatemodal.RegisterComponent(),
		invoiceupdatemodal.RegisterComponent(),
		ripsiteupdatemodal.RegisterComponent(),
		worksitetable.RegisterComponent(),
		ripsitetable.RegisterComponent(),
		invoicetable.RegisterComponent(),
		teamproductivitymodal.RegisterComponent(),
		adminmodal.RegisterComponent(),
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
		hvue.Computed("ReviewWorksiteInfos", func(vm *hvue.VM) interface{} {
			mpm := &MainPageModel{Object: vm.Object}
			return mpm.GetReviewableWorsiteInfos()
		}),
		hvue.Computed("UpdatableWorksiteInfos", func(vm *hvue.VM) interface{} {
			mpm := &MainPageModel{Object: vm.Object}
			return mpm.GetUpdatableWorsiteInfos()
		}),
		hvue.Computed("ReworkWorksiteInfos", func(vm *hvue.VM) interface{} {
			mpm := &MainPageModel{Object: vm.Object}
			return mpm.GetReworkWorksiteInfos()
		}),
		hvue.Computed("BillableWorksiteInfos", func(vm *hvue.VM) interface{} {
			mpm := &MainPageModel{Object: vm.Object}
			return mpm.GetBillableWorksiteInfos()
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

	SiteMode   string `js:"SiteMode"`
	ActiveMode string `js:"ActiveMode"`

	WorksiteInfos []*fm.WorksiteInfo `js:"worksiteInfos"`
	RipsiteInfos  []*fm.RipsiteInfo  `js:"ripsiteInfos"`
	//EditedWorksite int      `js:"editedWorksite"`
}

func NewMainPageModel() *MainPageModel {
	mpm := &MainPageModel{Object: tools.O()}
	mpm.User = fm.NewUser()
	mpm.WorksiteInfos = []*fm.WorksiteInfo{}
	mpm.RipsiteInfos = []*fm.RipsiteInfo{}
	mpm.SetMode()

	return mpm
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Action Methods

func (m *MainPageModel) SetMode() {
	// SiteMode setting
	if len(m.RipsiteInfos) > 0 && len(m.WorksiteInfos) == 0 {
		m.SiteMode = "Rip"
	} else {
		m.SiteMode = "Orange"
	}

	// ActiveMode setting
	switch {
	case m.User.Name == "":
		m.ActiveMode = ""
	case m.User.Permissions["Create"]:
		m.ActiveMode = "Create"
	case m.User.Permissions["Update"]:
		m.ActiveMode = "Update"
	case m.User.Permissions["Invoice"]:
		m.ActiveMode = "Invoice"
	case m.User.Permissions["Review"]:
		m.ActiveMode = "Review"
	default:
		m.ActiveMode = ""
	}
	return
}

func (m *MainPageModel) GetUserSession() {
	go m.callGetUser()
}

func (m *MainPageModel) ShowUserLogin() {
	m.VM.Refs("UserLoginModal").Call("Show", m.User)
}

func (m *MainPageModel) UserLogout() {
	go m.callLogout()
}

func (m *MainPageModel) GetSiteInfos() {
	go m.callGetWorkSiteInfos()
	go m.callGetRipSiteInfos()
}

func (m *MainPageModel) GetWorkSiteInfos() {
	go m.callGetWorkSiteInfos()
}

func (m *MainPageModel) GetRipSiteInfos() {
	go m.callGetRipSiteInfos()
}

func (m *MainPageModel) EditWorksite(id int) {
	m.VM.Refs("WorksiteEditModal").Call("Show", id, m.User)
}

func (m *MainPageModel) UpdateWorksite(id int) {
	m.VM.Refs("WorksiteUpdateModal").Call("Show", id, m.User)
}

func (m *MainPageModel) UpdateRipsite(id int) {
	m.VM.Refs("RipsiteUpdateModal").Call("Show", id, m.User)
}

func (m *MainPageModel) CreateNewWorksite() {
	m.EditWorksite(-1)
}

func (m *MainPageModel) EditRework(id int) {
	m.VM.Refs("ReworkEditModal").Call("Show", id, m.User)
}

func (m *MainPageModel) UpdateRework(id int) {
	m.VM.Refs("ReworkUpdateModal").Call("Show", id, m.User)
}

func (m *MainPageModel) UpdateInvoice(id int) {
	m.VM.Refs("InvoiceUpdateModal").Call("Show", id, m.User)
}

func (m *MainPageModel) ShowTeamProductivity() {
	m.VM.Refs("TeamProductivityModal").Call("Show", m.User)
}

func (m *MainPageModel) ShowAdmin() {
	m.VM.Refs("AdminModal").Call("Show", m.User)
}

func (m *MainPageModel) GetUpdatableWorsiteInfos() []*fm.WorksiteInfo {
	res := []*fm.WorksiteInfo{}
	for _, wsi := range m.WorksiteInfos {
		if fm.WorksiteIsUpdatable(wsi.Status) || wsi.NeedRework() {
			res = append(res, wsi)
		}
	}
	return res
}

func (m *MainPageModel) GetReviewableWorsiteInfos() []*fm.WorksiteInfo {
	res := []*fm.WorksiteInfo{}
	for _, wsi := range m.WorksiteInfos {
		if fm.WorksiteIsReviewable(wsi.Status) {
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
		//if fm.WorksiteMustRework(wsi.Status) {
		if wsi.NeedRework() {
			res = append(res, wsi)
		}
	}
	return res
}

func (m *MainPageModel) GetReworkWorsiteNb() int {
	res := 0
	for _, wsi := range m.WorksiteInfos {
		if wsi.NeedRework() {
			res += 1
		}
	}
	return res
}

func (m *MainPageModel) GetBillableWorksiteInfos() []*fm.WorksiteInfo {
	res := []*fm.WorksiteInfo{}
	for _, wsi := range m.WorksiteInfos {
		if fm.WorksiteIsBillable(wsi.Status) {
			res = append(res, wsi)
		}
	}
	return res
}

func (m *MainPageModel) GetBillableWorksiteNb() int {
	res := 0
	for _, wsi := range m.WorksiteInfos {
		if fm.WorksiteIsBillable(wsi.Status) {
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
	m.User.Copy(fm.UserFromJS(req.Response))
	if m.User.Name != "" {
		m.User.Connected = true
		m.SetMode()
		m.GetSiteInfos()
		return
	}
	m.User = fm.NewUser()
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
	m.User = fm.NewUser()
	m.User.Connected = false
	m.WorksiteInfos = []*fm.WorksiteInfo{}
	m.SetMode()
}

func (m *MainPageModel) callGetWorkSiteInfos() {
	req := xhr.NewRequest("GET", "/api/worksites")
	req.Timeout = tools.LongTimeOut
	req.ResponseType = xhr.JSON
	sites := m.WorksiteInfos
	m.WorksiteInfos = nil
	defer func() {
		m.WorksiteInfos = sites
	}()
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
	sites = wsis
}

func (m *MainPageModel) callGetRipSiteInfos() {
	req := xhr.NewRequest("GET", "/api/ripsites")
	req.Timeout = tools.LongTimeOut
	req.ResponseType = xhr.JSON
	sites := m.RipsiteInfos
	m.RipsiteInfos = nil
	defer func() {
		m.RipsiteInfos = sites
	}()
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
	rsis := []*fm.RipsiteInfo{}
	req.Response.Call("forEach", func(item *js.Object) {
		rs := fm.NewRipsiteInfoFromJS(item)
		rsis = append(rsis, rs)
	})
	sites = rsis
}
