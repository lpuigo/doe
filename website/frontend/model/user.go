package model

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

type User struct {
	*js.Object

	Name        string          `js:"Name"`
	Pwd         string          `js:"Pwd"`
	Connected   bool            `js:"Connected"`
	Clients     []string        `js:"Clients"`
	Permissions map[string]bool `js:"Permissions"`
	Teams       []string        `js:"Teams"`
}

func NewUser() *User {
	user := &User{Object: tools.O()}
	user.Name = ""
	user.Pwd = ""
	user.Connected = false
	user.Clients = []string{}
	user.Permissions = make(map[string]bool)
	user.Teams = []string{}
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
	u.Teams = ou.Teams
}
