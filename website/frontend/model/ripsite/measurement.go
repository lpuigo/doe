package ripsite

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/model/ripsite/ripconst"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/json"
	"strings"
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

func (m *Measurement) UpdateWith(mr *MeasurementReport, actors []string) {
	// If KO are reported => set measurement to InProgress
	m.State.Actors = actors[:]
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

func (m *Measurement) SearchString(filter string) string {
	searchItem := func(prefix, typ, value string) string {
		if value == "" {
			return ""
		}
		if filter != ripconst.FilterValueAll && filter != typ {
			return ""
		}
		return prefix + typ + value
	}

	res := searchItem("", ripconst.FilterValueComment, m.State.Comment)
	res += searchItem(",", ripconst.FilterValuePtRef, m.DestNodeName)
	for _, nodename := range m.NodeNames {
		res += searchItem(",", ripconst.FilterValuePtRef, nodename)
	}
	return res
}

func (m *Measurement) GetNbFiber() int {
	return m.NbFiber
}

type Warning struct {
	WarnLvl string
	Dist    float64
}

func (m *Measurement) ParseComment() []Warning {
	res := []Warning{}
	// Split multiple msg per line
	msgs := strings.Split(m.State.Comment, "\n")
	for _, msg := range msgs {
		// Split "Fib. #999: " from actual msg
		parts := strings.Split(msg, ": ")
		if len(parts) != 2 {
			continue
		}
		warnmsg := parts[1]
		// split diff warn msg on same fiber
		for _, wmsg := range strings.Split(warnmsg, ", ") {
			// Get distance (msg should be "some text[ à 99.9db]"
			newWarn := Warning{}
			distparts := strings.Split(wmsg, " à ")
			if len(distparts) > 1 {
				newWarn.Dist = js.Global.Call("parseFloat", distparts[1]).Float()
			}
			curWarnMsg := "undef"
			switch {
			case strings.Contains(distparts[0], "KO Max Splice"):
				curWarnMsg = "KO Splice"
			case strings.Contains(distparts[0], "Warn2 Max Splice"):
				curWarnMsg = "Warn2"
			case strings.Contains(distparts[0], "Warn1 Max Splice"):
				curWarnMsg = "Warn1"
			case strings.Contains(distparts[0], "Max Connector"):
				curWarnMsg = "KO Connector"
			}
			newWarn.WarnLvl = curWarnMsg

			res = append(res, newWarn)
		}
	}
	return res
}
