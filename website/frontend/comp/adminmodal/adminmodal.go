package adminmodal

import (
	"sort"
	"strconv"
	"strings"

	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/modal"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/model/beuser"
	"github.com/lpuig/ewin/doe/website/frontend/model/group"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements/message"
	"honnef.co/go/js/xhr"
)

type AdminModalModel struct {
	*modal.ModalModel

	User       *fm.User          `js:"user"`
	UsersStore *beuser.Store     `js:"UsersStore"`
	GroupStore *group.GroupStore `js:"GroupStore"`
}

func NewAdminModalModel(vm *hvue.VM) *AdminModalModel {
	tpmm := &AdminModalModel{
		ModalModel: modal.NewModalModel(vm),
	}
	tpmm.User = fm.NewUser()
	tpmm.UsersStore = beuser.NewStore()
	tpmm.GroupStore = group.NewGroupStore()
	return tpmm
}

func AdminModalModelFromJS(o *js.Object) *AdminModalModel {
	tpmm := &AdminModalModel{
		ModalModel: &modal.ModalModel{Object: o},
	}
	return tpmm
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Component Methods

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("admin-modal", componentOption()...)
}

func componentOption() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewAdminModalModel(vm)
		}),
		hvue.MethodsOf(&AdminModalModel{}),
		hvue.Computed("filteredUsers", func(vm *hvue.VM) interface{} {
			amm := AdminModalModelFromJS(vm.Object)
			return amm.UsersStore.Users
		}),
		hvue.Computed("hasChanged", func(vm *hvue.VM) interface{} {
			amm := AdminModalModelFromJS(vm.Object)
			return amm.UsersStore.Ref.IsDirty()
		}),
	}
}

func (amm *AdminModalModel) ReloadData() {
	go amm.callReloadData()
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Modal Methods

func (amm *AdminModalModel) Show(user *fm.User) {
	amm.User = user
	amm.Loading = false
	amm.UsersStore.CallGetUsers(amm.VM, func() {})
	amm.GroupStore.CallGetGroups(amm.VM, func() {})
	amm.ModalModel.Show()
}

func (amm *AdminModalModel) HideWithControl(user *fm.User) {
	if amm.UsersStore.Ref.Dirty {
		message.ConfirmCancelWarning(amm.VM, "Sauvegarder les modifications apportées ?",
			func() { // confirm
				amm.UsersStore.CallUpdateUsers(amm.VM, func() {
					amm.ModalModel.Hide()
				})
			},
			func() {
				amm.ModalModel.Hide()
			},
		)
	}
}

func (amm *AdminModalModel) UndoChange() {
	amm.UsersStore.Users = amm.UsersStore.GetReferenceUsers()
}

func (amm *AdminModalModel) ConfirmChange() {
	amm.UsersStore.CallUpdateUsers(amm.VM, func() {})
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Users Tabs Methods

func (amm *AdminModalModel) TableRowClassName(vm *hvue.VM) string {
	return ""
}

func (amm *AdminModalModel) ClientList(vm *hvue.VM, usr *beuser.BeUser) string {
	clients := amm.GetVisibleClientForUser(usr)
	if len(clients) == 0 {
		return "Tous"
	}
	return strings.Join(clients, ", ")
}

func (amm *AdminModalModel) GroupList(vm *hvue.VM, usr *beuser.BeUser) string {
	amm = AdminModalModelFromJS(vm.Object)
	groupNames := make([]string, len(usr.Groups))
	for i, grpId := range usr.Groups {
		groupNames[i] = amm.GroupStore.GetGroupNameById(grpId)
	}
	if len(groupNames) == 0 {
		return "Aucun"
	}
	sort.Strings(groupNames)
	return strings.Join(groupNames, ", ")
}

func (amm *AdminModalModel) GetClientList(vm *hvue.VM) []*elements.ValueLabel {
	amm = AdminModalModelFromJS(vm.Object)
	res := []*elements.ValueLabel{}
	for _, clientName := range amm.User.GetSortedClientNames() {
		res = append(res, elements.NewValueLabel(clientName, clientName))
	}
	return res
}

func (amm *AdminModalModel) GetGroupList(vm *hvue.VM) []*elements.IntValueLabel {
	amm = AdminModalModelFromJS(vm.Object)
	res := []*elements.IntValueLabel{}
	for _, grp := range amm.GroupStore.GetGroupsSortedByName() {
		res = append(res, elements.NewIntValueLabel(grp.Id, grp.Name))
	}
	return res
}

func (amm *AdminModalModel) AddNewUser(vm *hvue.VM) {
	amm = AdminModalModelFromJS(vm.Object)
	nuser := beuser.NewBeUser()
	nuser.Name = "Nouvel Utilisateur"
	nuser.Password = "default"
	amm.UsersStore.AddNewUser(nuser)
}

func (amm *AdminModalModel) UpdateUserClients(usr *beuser.BeUser) {
	usr.SortClients()
}

func (amm *AdminModalModel) GetVisibleClientForUser(usr *beuser.BeUser) []string {
	if len(usr.Groups) > 0 {
		clientDict := make(map[string]int)
		for _, grpId := range usr.Groups {
			grp := amm.GroupStore.GetGroupById(grpId)
			if grp == nil {
				// GrpId does not exist, remove it
				usr.RemoveGroupId(grpId)
				continue
			}
			for _, clientName := range grp.Clients {
				clientDict[clientName]++
			}
		}
		res := make([]string, len(clientDict))
		i := 0
		for clientName, _ := range clientDict {
			res[i] = clientName
			i++
		}
		sort.Strings(res)
		return res
	}
	return usr.Clients
}

// Column Filtering Related Methods

func (amm *AdminModalModel) FilterHandler(vm *hvue.VM, value string, p *js.Object, col *js.Object) bool {
	prop := col.Get("property").String()
	grp := group.GroupFromJS(p)
	switch prop {
	case "Client":
		for _, c := range grp.Clients {
			if c == value {
				return true
			}
		}
		return false
	}
	return p.Get(prop).String() == value
}

func (amm *AdminModalModel) FilterList(vm *hvue.VM, prop string) []*elements.ValText {
	amm = AdminModalModelFromJS(vm.Object)
	count := map[string]int{}
	attribs := []string{}

	var translate func(string) string
	//switch prop {
	//case "State":
	//	translate = func(state string) string {
	//		return GetStateLabel(state)
	//	}
	//default:
	translate = func(val string) string { return val }
	//}

	for _, usr := range amm.UsersStore.Users {
		var attrs []string
		switch prop {
		case "Client":
			attrs = usr.Clients
		default:
			attrs = []string{usr.Object.Get(prop).String()}
		}
		for _, a := range attrs {
			if _, exist := count[a]; !exist {
				attribs = append(attribs, a)
			}
			count[a]++
		}
	}
	sort.Strings(attribs)
	res := []*elements.ValText{}
	for _, a := range attribs {
		fa := a
		if fa == "" {
			fa = "Vide"
		}
		res = append(res, elements.NewValText(a, translate(fa)+" ("+strconv.Itoa(count[a])+")"))
	}
	return res
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// WS call Methods

func (amm *AdminModalModel) callReloadData() {
	defer func() { amm.Loading = false }()
	req := xhr.NewRequest("GET", "/api/admin/reload")
	req.Timeout = tools.TimeOut
	req.ResponseType = xhr.JSON
	err := req.Send(nil)
	if err != nil {
		message.ErrorStr(amm.VM, "Oups! "+err.Error(), true)
		amm.Hide()
		return
	}
	if req.Status != tools.HttpOK {
		message.ErrorRequestMessage(amm.VM, req)
		amm.Hide()
		return
	}
	message.SuccesStr(amm.VM, "Rechargement des données effectué")
	amm.VM.Emit("reload")
	return
}
