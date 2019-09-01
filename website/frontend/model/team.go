package model

import "github.com/gopherjs/gopherjs/js"

type Team struct {
	*js.Object

	Name     string `js:"Name"`
	Members  string `js:"Members"`
	IsActive bool   `js:"IsActive"`
}

type Actor struct {
	*js.Object

	Id        int    `js:"Id"`
	LastName  string `js:"LastName"`
	FirstName string `js:"FirstName"`
	Role      string `js:"Role"`
	Active    bool   `js:"Active"`
}

func (a *Actor) GetRef() string {
	ext := ""
	if !a.Active {
		ext = " (parti)"
	} else {
		ext = " (" + a.Role + ")"
	}
	return a.LastName + " " + a.FirstName + ext
}
