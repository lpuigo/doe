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
	rs.Object = json.Parse(json.Stringify(ors))
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
