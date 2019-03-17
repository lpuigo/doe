package model

import "github.com/gopherjs/gopherjs/js"

type Client struct {
	*js.Object

	Name  string  `js:"Name"`
	Teams []*Team `js:"Teams"`
}
