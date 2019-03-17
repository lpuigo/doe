package model

import "github.com/gopherjs/gopherjs/js"

type Team struct {
	*js.Object

	Name     string `js:"Name"`
	Members  string `js:"Members"`
	IsActive bool   `js:"IsActive"`
}
