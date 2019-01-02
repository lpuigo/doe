package model

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

type User struct {
	*js.Object

	Name string `js:"Name"`
	Pwd  string `js:"Pwd"`
}

func NewUser() *User {
	user := &User{Object: tools.O()}
	user.Name = ""
	user.Pwd = ""
	return user
}

func (u *User) Copy(ou *User) {
	u.Name = ou.Name
	u.Pwd = ou.Pwd
}
