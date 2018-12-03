package frontmodel

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
	o.Troncons = make([]*Troncon, len(oo.Troncons))
	for i, tr := range oo.Troncons {
		o.Troncons[i] = tr.Clone()
	}
}
