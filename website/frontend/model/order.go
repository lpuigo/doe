package model

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

type Order struct {
	*js.Object

	Ref      string     `js:"Ref"`
	Comment  string     `js:"Comment"`
	Troncons []*Troncon `js:"Troncons"`
}

func NewOrder() *Order {
	o := &Order{Object: tools.O()}
	o.Ref = ""
	o.Comment = ""
	o.Troncons = []*Troncon{}
	return o
}

func (o *Order) Clone() *Order {
	no := &Order{Object: tools.O()}
	no.Copy(o)
	return no
}

func (o *Order) Copy(oo *Order) {
	o.Ref = oo.Ref
	o.Comment = oo.Comment
	o.Troncons = []*Troncon{}
	for _, tr := range oo.Troncons {
		o.Troncons = append(o.Troncons, tr.Clone())
	}
}

func (o *Order) SearchInString() string {
	res := "O_Ref" + o.Ref + "\n"
	res += "O_Comment" + o.Comment + "\n"

	for _, t := range o.Troncons {
		res += t.SearchInString()
	}

	return res
}

func (o *Order) DeleteTroncon(i int) {
	// Not working with GopherJS
	//copy(o.Troncons[i:], o.Troncons[i+1:])
	//o.Troncons[len(o.Troncons)-1] = nil // or the zero value of T
	//o.Troncons = o.Troncons[:len(o.Troncons)-1]
	nts := []*Troncon{}
	for j, t := range o.Troncons {
		if i == j {
			continue
		}
		nts = append(nts, t)
	}
	o.Troncons = nts
}

func (o *Order) AddTroncon() {
	troncon := NewTroncon()
	o.Troncons = append(o.Troncons, troncon)
}

func (o *Order) IsCompleted() bool {
	for _, t := range o.Troncons {
		if !t.IsCompleted() {
			return false
		}
	}
	return true
}
