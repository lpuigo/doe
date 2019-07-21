package polesite

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

// Polesite reflects backend/model/polesites.polesite struct
type Polesite struct {
	*js.Object

	Id         int     `js:"Id"`
	Client     string  `js:"Client"`
	Ref        string  `js:"Ref"`
	Manager    string  `js:"Manager"`
	OrderDate  string  `js:"OrderDate"`
	UpdateDate string  `js:"UpdateDate"`
	Status     string  `js:"Status"`
	Comment    string  `js:"Comment"`
	Poles      []*Pole `js:"Poles"`
}

func NewPolesite() *Polesite {
	return &Polesite{Object: tools.O()}
}

func PolesiteFromJS(o *js.Object) *Polesite {
	return &Polesite{Object: o}
}
