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

func (ps *Polesite) getNextId() int {
	// naive algorithm ... something smarter must be possible
	maxid := -1
	for _, pole := range ps.Poles {
		if pole.Id >= maxid {
			maxid = pole.Id + 1
		}
	}
	return maxid
}

// AddPole adds the given pole to polesite, and sets pole's new Id to ensure Id unicity
func (ps *Polesite) AddPole(pole *Pole) {
	pole.Id = ps.getNextId()
	ps.Poles = append(ps.Poles, pole)
}

// DeletePole deletes the given pole and returns true if it was found and deleted, false otherwise (no-op)
func (ps *Polesite) DeletePole(pole *Pole) bool {
	for i, p := range ps.Poles {
		if p.Id == pole.Id {
			// remove the item the JS way, to triggger vueJS observers
			ps.Object.Get("Poles").Call("splice", i, 1)
			return true
		}
	}
	return false
}

func NewPolesite() *Polesite {
	return &Polesite{Object: tools.O()}
}

func PolesiteFromJS(o *js.Object) *Polesite {
	return &Polesite{Object: o}
}
