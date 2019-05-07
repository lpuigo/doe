package ripsite

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/json"
)

type Measurement struct {
	*js.Object

	DestNodeName string `js:"DestNodeName"`
	NbFiber      int    `js:"NbFiber"` // DestNode.Operation.Attente.NbFiber
	//NbEvent      int // == len(NodeNames)
	Dist      int      `js:"Dist"` // DestNode.DistFromPM
	NodeNames []string `js:"NodeNames"`
	State     *State   `js:"State"`
}

func NewMeasurement() *Measurement {
	m := &Measurement{Object: tools.O()}
	m.DestNodeName = ""
	m.NbFiber = 0
	m.Dist = 0
	m.NodeNames = nil
	m.State = NewState()

	return m
}

func (m *Measurement) Clone() *Measurement {
	return &Measurement{Object: json.Parse(json.Stringify(m))}
}
