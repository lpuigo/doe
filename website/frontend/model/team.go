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

	LastName  string `js:"LastName"`
	FirstName string `js:"FirstName"`
	Role      string `js:"Role"`
}
