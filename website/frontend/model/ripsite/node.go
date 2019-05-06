package ripsite

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

type Node struct {
	*js.Object

	Name          string `js:"Name"`
	Address       string `js:"Address"`
	Type          string `js:"Type"`
	BoxType       string `js:"BoxType"`
	Ref           string `js:"Ref"`
	TronconInName string `js:"TronconInName"`
	DistFromPm    int    `js:"DistFromPm"`
}

func NewNode() *Node {
	n := &Node{Object: tools.O()}
	n.Name = ""
	n.Address = ""
	n.Type = ""
	n.BoxType = ""
	n.Ref = ""
	n.TronconInName = ""
	n.DistFromPm = 0

	return n
}
