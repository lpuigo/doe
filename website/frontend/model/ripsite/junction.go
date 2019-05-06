package ripsite

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

type Operation struct {
	*js.Object

	Type        string `js:"Type"`
	TronconName string `js:"TronconName"`
	NbFiber     int    `js:"NbFiber"`
	NbSplice    int    `js:"NbSplice"`
	State       *State `js:"State"`
}

func NewOperation() *Operation {
	o := &Operation{Object: tools.O()}
	o.Type = ""
	o.TronconName = ""
	o.NbFiber = 0
	o.NbSplice = 0
	o.State = NewState()

	return o
}

type Junction struct {
	*js.Object

	NodeName   string       `js:"NodeName"`
	Operations []*Operation `js:"Operations"`
	State      *State       `js:"State"`
}

func NewJunction() *Junction {
	j := &Junction{Object: tools.O()}
	j.NodeName = ""
	j.Operations = nil
	j.State = NewState()

	return j
}
