package model

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

type User struct {
	*js.Object

	Name      string `js:"Name"`
	Pwd       string `js:"Pwd"`
	Connected bool   `js:"Connected"`
}

func NewUser() *User {
	user := &User{Object: tools.O()}
	user.Name = ""
	user.Pwd = ""
	user.Connected = false
	return user
}

func NewUserFromJS(o *js.Object) *User {
	return &User{Object: o}
}

func (u *User) Copy(ou *User) {
	u.Name = ou.Name
	u.Connected = ou.Connected
}
