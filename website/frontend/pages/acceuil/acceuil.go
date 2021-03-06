package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/adminmodal"
	"github.com/lpuig/ewin/doe/website/frontend/comp/foasitetable"
	"github.com/lpuig/ewin/doe/website/frontend/comp/invoicetable"
	"github.com/lpuig/ewin/doe/website/frontend/comp/invoiceupdatemodal"
	"github.com/lpuig/ewin/doe/website/frontend/comp/polesitetable"
	"github.com/lpuig/ewin/doe/website/frontend/comp/reworkeditmodal"
	"github.com/lpuig/ewin/doe/website/frontend/comp/reworkupdatemodal"
	"github.com/lpuig/ewin/doe/website/frontend/comp/ripsitetable"
	"github.com/lpuig/ewin/doe/website/frontend/comp/teamproductivitymodal"
	"github.com/lpuig/ewin/doe/website/frontend/comp/userloginmodal"
	"github.com/lpuig/ewin/doe/website/frontend/comp/worksiteeditmodal"
	"github.com/lpuig/ewin/doe/website/frontend/comp/worksitetable"
	"github.com/lpuig/ewin/doe/website/frontend/comp/worksiteupdatemodal"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/model/worksite"
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
		//ripsiteupdatemodal.RegisterComponent(),
		worksitetable.RegisterComponent(),
		ripsitetable.RegisterComponent(),
		foasitetable.RegisterComponent(),
		polesitetable.RegisterComponent(),
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
		hvue.Computed("SiteModeLabel", func(vm *hvue.VM) interface{} {
			mpm := &MainPageModel{Object: vm.Object}
			return mpm.SiteModeLabel()
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
	FoasiteInfos  []*fm.FoaSiteInfo  `js:"foasiteInfos"`
	PolesiteInfos []*fm.PolesiteInfo `js:"polesiteInfos"`
}

func NewMainPageModel() *MainPageModel {
	mpm := &MainPageModel{Object: tools.O()}
	mpm.User = fm.NewUser()
	mpm.ClearSiteInfos()
	mpm.ClearModes()
	//mpm.SetMode()

	return mpm
}

func (m *MainPageModel) ClearModes() {
	m.SiteMode = ""
	m.ActiveMode = ""
}

func (m *MainPageModel) ClearSiteInfos() {
	m.WorksiteInfos = []*fm.WorksiteInfo{}
	m.RipsiteInfos = []*fm.RipsiteInfo{}
	m.FoasiteInfos = []*fm.FoaSiteInfo{}
	m.PolesiteInfos = []*fm.PolesiteInfo{}
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Action Methods

func (m *MainPageModel) CheckSiteMode(mode string) {
	m.SiteMode = mode
	if m.SiteMode == "Orange" && m.User.Permissions["Create"] {
		m.ActiveMode = "Create"
	} else {
		m.ActiveMode = "Update"
	}
}

func (m *MainPageModel) SetMode() {
	// SiteMode setting if not choosen yet
	if m.SiteMode == "" {
		if len(m.WorksiteInfos) > 0 {
			m.SiteMode = "Orange"
		} else if len(m.RipsiteInfos) > 0 {
			m.SiteMode = "Rip"
		} else if len(m.FoasiteInfos) > 0 {
			m.SiteMode = "Foa"
		} else if len(m.PolesiteInfos) > 0 {
			m.SiteMode = "Poles"
		}
	}

	// Set ActiveMode if not set yet
	if m.ActiveMode != "" {
		return
	}

	switch {
	case m.User.Name == "":
		m.ActiveMode = ""
	case m.User.Permissions["Update"]:
		m.ActiveMode = "Update"
	case m.User.Permissions["Create"]:
		m.ActiveMode = "Create"
	case m.User.Permissions["Invoice"]:
		m.ActiveMode = "Invoice"
		m.SiteMode = "Orange"
	case m.User.Permissions["Review"]:
		m.ActiveMode = "Review"
	default:
		m.ActiveMode = ""
	}
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
	go m.callGetFoaSiteInfos()
	go m.callGetPoleSiteInfos()
}

func (m *MainPageModel) GetActiveSiteInfos() {
	switch m.SiteMode {
	case "Orange":
		go m.callGetWorkSiteInfos()
	case "Rip":
		go m.callGetRipSiteInfos()
	case "Foa":
		go m.callGetFoaSiteInfos()
	case "Poles":
		go m.callGetPoleSiteInfos()
	}
}

func (m *MainPageModel) GetWorkSiteInfos() {
	go m.callGetWorkSiteInfos()
}

func (m *MainPageModel) GetRipSiteInfos() {
	go m.callGetRipSiteInfos()
}

func (m *MainPageModel) GetFoaSiteInfos() {
	go m.callGetFoaSiteInfos()
}

func (m *MainPageModel) GetPoleSiteInfos() {
	go m.callGetPoleSiteInfos()
}

func (m *MainPageModel) EditWorksite(id int) {
	m.VM.Refs("WorksiteEditModal").Call("Show", id, m.User)
}

func (m *MainPageModel) UpdateWorksite(id int) {
	m.VM.Refs("WorksiteUpdateModal").Call("Show", id, m.User)
}

//func (m *MainPageModel) UpdateRipsite(id int) {
//	m.VM.Refs("RipsiteUpdateModal").Call("Show", id, m.User)
//}
//
func (m *MainPageModel) OpenRipsite(id int) {
	js.Global.Get("window").Call("open", "ripsite.html?v=1.0&rsid="+strconv.Itoa(id))
}

func (m *MainPageModel) OpenActors() {
	js.Global.Get("window").Call("open", "actor.html")
}

func (m *MainPageModel) OpenVehicules() {
	js.Global.Get("window").Call("open", "vehicule.html")
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
	m.VM.Refs("TeamProductivityModal").Call("Show", m.User, m.SiteMode)
}

func (m *MainPageModel) ShowAdmin() {
	m.VM.Refs("AdminModal").Call("Show", m.User)
}

func (m *MainPageModel) GetUpdatableWorsiteInfos() []*fm.WorksiteInfo {
	res := []*fm.WorksiteInfo{}
	for _, wsi := range m.WorksiteInfos {
		if worksite.WorksiteIsUpdatable(wsi.Status) || wsi.NeedRework() {
			res = append(res, wsi)
		}
	}
	return res
}

func (m *MainPageModel) GetReviewableWorsiteInfos() []*fm.WorksiteInfo {
	res := []*fm.WorksiteInfo{}
	for _, wsi := range m.WorksiteInfos {
		if worksite.WorksiteIsReviewable(wsi.Status) {
			res = append(res, wsi)
		}
	}
	return res
}

func (m *MainPageModel) GetUpdatableWorsiteNb() int {
	res := 0
	if m.SiteMode == "Orange" {
		for _, wsi := range m.WorksiteInfos {
			if worksite.WorksiteIsUpdatable(wsi.Status) {
				res += 1
			}
		}
		return res
	}
	return len(m.RipsiteInfos)
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
		if worksite.WorksiteIsBillable(wsi.Status) {
			res = append(res, wsi)
		}
	}
	return res
}

func (m *MainPageModel) GetBillableWorksiteNb() int {
	res := 0
	for _, wsi := range m.WorksiteInfos {
		if worksite.WorksiteIsBillable(wsi.Status) {
			res += 1
		}
	}
	return res
}

func (m *MainPageModel) SiteModeLabel() string {
	var res string
	switch m.SiteMode {
	case "Orange":
		res = "Orange : " + strconv.Itoa(len(m.WorksiteInfos))
	case "Rip":
		res = "Rip : " + strconv.Itoa(len(m.RipsiteInfos))
	case "Foa":
		res = "Foa : " + strconv.Itoa(len(m.FoasiteInfos))
	case "Poles":
		res = "Poteaux : " + strconv.Itoa(len(m.PolesiteInfos))
	}
	return res
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// WS call Methods

//func (m *MainPageModel) errorMessage(req *xhr.Request) {
//	message.SetDuration(tools.WarningMsgDuration)
//	msg := "Quelquechose c'est mal passé !\n"
//	msg += "Le server retourne un code " + strconv.Itoa(req.Status) + "\n"
//	message.ErrorMsgStr(m.VM, msg, req.Response, true)
//}
//
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
		message.ErrorRequestMessage(m.VM, req)
		return
	}
	m.User.Copy(fm.UserFromJS(req.Response))
	if m.User.Name != "" {
		m.User.Connected = true
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
		message.ErrorRequestMessage(m.VM, req)
		return
	}
	m.User = fm.NewUser()
	m.User.Connected = false
	m.ClearSiteInfos()
	m.ClearModes()
}

func (m *MainPageModel) callGetWorkSiteInfos() {
	req := xhr.NewRequest("GET", "/api/worksites")
	req.Timeout = tools.LongTimeOut
	req.ResponseType = xhr.JSON
	sites := m.WorksiteInfos
	//m.WorksiteInfos = nil
	defer func() {
		m.WorksiteInfos = sites
		m.SetMode()
	}()
	//m.DispPrj = false
	err := req.Send(nil)
	if err != nil {
		message.ErrorStr(m.VM, "Oups! "+err.Error(), true)
		return
	}
	if req.Status != tools.HttpOK {
		message.ErrorRequestMessage(m.VM, req)
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
	//m.RipsiteInfos = nil
	defer func() {
		m.RipsiteInfos = sites
		m.SetMode()
	}()
	//m.DispPrj = false
	err := req.Send(nil)
	if err != nil {
		message.ErrorStr(m.VM, "Oups! "+err.Error(), true)
		return
	}
	if req.Status != tools.HttpOK {
		message.ErrorRequestMessage(m.VM, req)
		return
	}
	rsis := []*fm.RipsiteInfo{}
	req.Response.Call("forEach", func(item *js.Object) {
		rs := fm.NewRipsiteInfoFromJS(item)
		rsis = append(rsis, rs)
	})
	sites = rsis
}

func (m *MainPageModel) callGetFoaSiteInfos() {
	req := xhr.NewRequest("GET", "/api/foasites")
	req.Timeout = tools.LongTimeOut
	req.ResponseType = xhr.JSON
	sites := m.FoasiteInfos
	//m.RipsiteInfos = nil
	defer func() {
		m.FoasiteInfos = sites
		m.SetMode()
	}()
	//m.DispPrj = false
	err := req.Send(nil)
	if err != nil {
		message.ErrorStr(m.VM, "Oups! "+err.Error(), true)
		return
	}
	if req.Status != tools.HttpOK {
		message.ErrorRequestMessage(m.VM, req)
		return
	}
	fsis := []*fm.FoaSiteInfo{}
	req.Response.Call("forEach", func(item *js.Object) {
		rs := fm.NewFoaSiteInfoFromJS(item)
		fsis = append(fsis, rs)
	})
	sites = fsis
}

func (m *MainPageModel) callGetPoleSiteInfos() {
	req := xhr.NewRequest("GET", "/api/polesites")
	req.Timeout = tools.LongTimeOut
	req.ResponseType = xhr.JSON
	sites := m.PolesiteInfos
	//m.RipsiteInfos = nil
	defer func() {
		m.PolesiteInfos = sites
		m.SetMode()
	}()
	//m.DispPrj = false
	err := req.Send(nil)
	if err != nil {
		message.ErrorStr(m.VM, "Oups! "+err.Error(), true)
		return
	}
	if req.Status != tools.HttpOK {
		message.ErrorRequestMessage(m.VM, req)
		return
	}
	rsis := []*fm.PolesiteInfo{}
	req.Response.Call("forEach", func(item *js.Object) {
		rs := fm.NewPolesiteInfoFromJS(item)
		rsis = append(rsis, rs)
	})
	sites = rsis
}
