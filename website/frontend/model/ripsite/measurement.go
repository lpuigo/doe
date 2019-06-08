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
	NbOK         int    `js:"NbOK"`
	NbWarn1      int    `js:"NbWarn1"`
	NbWarn2      int    `js:"NbWarn2"`
	NbKO         int    `js:"NbKO"`
	//NbEvent      int // == len(NodeNames)
	Dist      int      `js:"Dist"` // DestNode.DistFromPM
	NodeNames []string `js:"NodeNames"`
	State     *State   `js:"State"`
}

func NewMeasurement() *Measurement {
	m := &Measurement{Object: tools.O()}
	m.DestNodeName = ""
	m.NbFiber = 0
	m.NbOK = 0
	m.NbWarn1 = 0
	m.NbWarn2 = 0
	m.NbKO = 0
	m.Dist = 0
	m.NodeNames = []string{}
	m.State = NewState()

	return m
}

func (m *Measurement) Clone() *Measurement {
	return &Measurement{Object: json.Parse(json.Stringify(m))}
}

func (m *Measurement) UpdateWith(mr *MeasurementReport, team string) {
	// If KO are reported => set measurement to InProgress
	m.State.Team = team
	m.NbOK = mr.FiberOK
	m.NbKO = mr.FiberKO
	m.NbWarn1 = mr.FiberWarning1
	m.NbWarn2 = mr.FiberWarning2

	if !(mr.FiberKO == 0 && mr.ConnectorKO == 0) {
		m.State.SetBlocked()
		if !(!tools.Empty(m.State.DateStart) && m.State.DateStart <= mr.Date) {
			m.State.DateStart = mr.Date
		}
		m.State.DateEnd = ""
	} else { // OK or just warning => set measurement to Done
		switch {
		case mr.FiberWarning2 > 0:
			m.State.SetWarning2()
		case mr.FiberWarning1 > 0:
			m.State.SetWarning1()
		default:
			m.State.SetDone()
		}

		if !(!tools.Empty(m.State.DateStart) && m.State.DateStart <= mr.Date) {
			m.State.DateStart = mr.Date
		}

		if !(!tools.Empty(m.State.DateEnd) && m.State.DateEnd <= mr.Date) {
			m.State.DateEnd = mr.Date
		}
	}
	if mr.Results != nil && len(mr.Results) > 0 {
		m.State.Comment = mr.Comments()
	}
}
