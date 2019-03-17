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
	Clients     []*Client       `js:"Clients"`
	Permissions map[string]bool `js:"Permissions"`
}

func NewUser() *User {
	user := &User{Object: tools.O()}
	user.Name = ""
	user.Pwd = ""
	user.Connected = false
	user.Clients = []*Client{}
	user.Permissions = make(map[string]bool)
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
}

func (u *User) GetClientByName(clientName string) *Client {
	for _, c := range u.Clients {
		if c.Name == clientName {
			return c
		}
	}
	return nil
}
