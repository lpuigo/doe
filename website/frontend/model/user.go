package model

import (
	"sort"

	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
)

type User struct {
	*js.Object

	Name        string            `js:"Name"`
	Pwd         string            `js:"Pwd"`
	Connected   bool              `js:"Connected"`
	Clients     []*Client         `js:"Clients"`
	Permissions map[string]bool   `js:"Permissions"`
	DaysOff     map[string]string `js:"DaysOff"`
}

func NewUser() *User {
	user := &User{Object: tools.O()}
	user.Name = ""
	user.Pwd = ""
	user.Connected = false
	user.Clients = []*Client{}
	user.Permissions = make(map[string]bool)
	user.DaysOff = make(map[string]string)
	return user
}

func UserFromJS(o *js.Object) *User {
	return &User{Object: o}
}

func (u *User) Copy(ou *User) {
	u.Name = ou.Name
	u.Pwd = ou.Pwd
	u.Connected = ou.Connected
	u.Clients = ou.Clients
	u.Permissions = ou.Permissions
	u.DaysOff = ou.DaysOff
}

// GetClientByName returns the client with given name (nil if not found)
func (u *User) GetClientByName(clientName string) *Client {
	for _, c := range u.Clients {
		if c.Name == clientName {
			return c
		}
	}
	return nil
}

// GetClientByName returns the client with given name (nil if not found)
func (u *User) GetSortedClientNames() []string {
	res := make([]string, len(u.Clients))
	for i, c := range u.Clients {
		res[i] = c.Name
	}
	sort.Strings(res)
	return res
}

func (u *User) GetTeamValueLabelsFor(clientName string) []*elements.ValueLabel {
	res := []*elements.ValueLabel{}
	client := u.GetClientByName(clientName)
	if client == nil {
		return nil
	}
	for _, team := range client.Teams {
		if team.IsActive {
			res = append(res, elements.NewValueLabel(team.Members, team.Name+": "+team.Members))
		}
	}
	return res
}

func (u *User) GetArticlesValueLabelsFor(clientName string) []*elements.ValueLabel {
	res := []*elements.ValueLabel{}
	client := u.GetClientByName(clientName)
	if client == nil {
		return nil
	}
	for _, a := range client.Articles {
		res = append(res, elements.NewValueLabel(a, a))
	}
	return res
}

func (u *User) IsDayOff(day string) bool {
	_, exists := u.DaysOff[day]
	return exists
}

func (u *User) HasPermissionInvoice() bool {
	return u.Permissions["Invoice"]
}

func (u *User) HasPermissionHR() bool {
	return u.Permissions["HR"]
}

func (u *User) HasPermissionUpdate() bool {
	return u.Permissions["Update"]
}
