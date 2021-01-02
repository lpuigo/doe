package beuser

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

// type BeUSer reflects backend/model/users.user struct
type BeUser struct {
	*js.Object
	Id          int             `js:"Id"`
	Name        string          `js:"Name"`
	Password    string          `js:"Password"`
	Clients     []string        `js:"Clients"`
	Groups      []int           `js:"Groups"`
	Permissions map[string]bool `js:"Permissions"`
}

func BeUserFromJS(obj *js.Object) *BeUser {
	bu := &BeUser{Object: obj}
	if bu.Get("Groups") == nil {
		bu.Groups = []int{}
	}
	return bu
}

func NewBeUser() *BeUser {
	usr := &BeUser{Object: tools.O()}
	usr.Id = -1
	usr.Name = ""
	usr.Password = ""
	usr.Clients = []string{}
	usr.Groups = []int{}
	usr.Permissions = make(map[string]bool)
	return usr
}

func (bu *BeUser) SortClients() {
	//sort.Strings(bu.Clients)
	bu.Object.Get("Clients").Call("sort")
}

func (bu *BeUser) RemoveGroupId(grpId int) {
	//sort.Strings(bu.Clients)
	for i, gId := range bu.Groups {
		if gId == grpId {
			bu.Object.Call("splice", i, 1)
			return
		}
	}
}
