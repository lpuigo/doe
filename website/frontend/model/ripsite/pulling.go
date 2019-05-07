package ripsite

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/json"
)

type PullingChunk struct {
	*js.Object

	TronconName      string `js:"TronconName"`
	StartingNodeName string `js:"StartingNodeName"`
	EndingNodeName   string `js:"EndingNodeName"`
	LoveDist         int    `js:"LoveDist"`
	UndergroundDist  int    `js:"UndergroundDist"`
	AerialDist       int    `js:"AerialDist"`
	BuildingDist     int    `js:"BuildingDist"`
	State            *State `js:"State"`
}

func NewPullingChunk() *PullingChunk {
	pc := &PullingChunk{Object: tools.O()}
	pc.TronconName = ""
	pc.StartingNodeName = ""
	pc.EndingNodeName = ""
	pc.LoveDist = 0
	pc.UndergroundDist = 0
	pc.AerialDist = 0
	pc.BuildingDist = 0
	pc.State = NewState()
	return pc
}

type Pulling struct {
	*js.Object

	CableName string          `js:"CableName"`
	Chuncks   []*PullingChunk `js:"Chuncks"`
	State     *State          `js:"State"`
}

func NewPulling() *Pulling {
	p := &Pulling{Object: tools.O()}
	p.CableName = ""
	p.Chuncks = nil
	p.State = NewState()

	return p
}

func (p *Pulling) Clone() *Pulling {
	return &Pulling{Object: json.Parse(json.Stringify(p))}
}
