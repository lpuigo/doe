package ripsite

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/json"
)

type Ripsite struct {
	*js.Object

	Id        int                 `js:"Id"`
	Client    string              `js:"Client"`
	Ref       string              `js:"Ref"`
	Manager   string              `js:"Manager"`
	OrderDate string              `js:"OrderDate"`
	Status    string              `js:"Status"`
	Comment   string              `js:"Comment"`
	Nodes     map[string]*Node    `js:"Nodes"`
	Troncons  map[string]*Troncon `js:"Troncons"`

	Pullings     []*Pulling     `js:"Pullings"`
	Junctions    []*Junction    `js:"Junctions"`
	Measurements []*Measurement `js:"Measurements"`
}

func RipsiteFromJS(o *js.Object) *Ripsite {
	return &Ripsite{Object: o}
}

func NewRisite() *Ripsite {
	rs := &Ripsite{Object: tools.O()}
	rs.Id = -1
	rs.Client = ""
	rs.Ref = ""
	rs.Manager = ""
	rs.OrderDate = ""
	rs.Status = ""
	rs.Comment = ""
	rs.Nodes = nil
	rs.Troncons = nil
	rs.Pullings = nil
	rs.Junctions = nil
	rs.Measurements = nil

	return rs
}

func (rs *Ripsite) SearchInString() string {
	return json.Stringify(rs)
}

func (rs *Ripsite) Copy(ors *Ripsite) {
	nodes := make(map[string]*Node)
	for nodeName, node := range ors.Nodes {
		nodes[nodeName] = node.Clone()
	}
	troncons := make(map[string]*Troncon)
	for trName, tr := range ors.Troncons {
		troncons[trName] = tr.Clone()
	}
	pullings := make([]*Pulling, len(ors.Pullings))
	for ip, pull := range ors.Pullings {
		pullings[ip] = pull.Clone()
	}
	junctions := make([]*Junction, len(ors.Junctions))
	for ij, junc := range ors.Junctions {
		junctions[ij] = junc.Clone()
	}
	measurements := make([]*Measurement, len(ors.Measurements))
	for im, meas := range ors.Measurements {
		measurements[im] = meas.Clone()
	}
	rs.Id = ors.Id
	rs.Client = ors.Client
	rs.Ref = ors.Ref
	rs.Manager = ors.Manager
	rs.OrderDate = ors.OrderDate
	rs.Status = ors.Status
	rs.Comment = ors.Comment
	rs.Nodes = nodes
	rs.Troncons = troncons
	rs.Pullings = pullings
	rs.Junctions = junctions
	rs.Measurements = measurements
}

func (rs *Ripsite) Clone() *Ripsite {
	return &Ripsite{Object: json.Parse(json.Stringify(rs))}
}

func (rs *Ripsite) GetInfo() (nbAvailPulling, nbPulling, nbAvailJunction, nbJunction, nbAvailMeas, nbMeas int) {
	nbPulling = len(rs.Pullings)
	nbAvailPulling = nbPulling
	for _, pulling := range rs.Pullings {
		if pulling.State.IsBlocked() {
			nbAvailPulling--
		}
	}
	nbJunction = len(rs.Junctions)
	nbAvailJunction = nbJunction
	for _, junction := range rs.Junctions {
		if junction.State.IsBlocked() {
			nbAvailJunction--
		}
	}
	nbMeas = len(rs.Measurements)
	nbAvailMeas = nbMeas
	for _, meas := range rs.Measurements {
		if meas.State.IsBlocked() {
			nbAvailMeas--
		}
	}
	return
}
