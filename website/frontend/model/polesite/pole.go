package polesite

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/model/polesite/poleconst"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

// type Pole reflects backend/model/polesites.pole struct
type Pole struct {
	*js.Object
	Ref      string            `js:"Ref"`
	City     string            `js:"City"`
	Address  string            `js:"Address"`
	Lat      float64           `js:"Lat"`
	Long     float64           `js:"Long"`
	State    string            `js:"State"`
	DtRef    string            `js:"DtRef"`
	DictRef  string            `js:"DictRef"`
	Height   int               `js:"Height"`
	Material string            `js:"Material"`
	Product  map[string]int    `js:"Product"`
	DictInfo map[string]string `js:"DictInfo"`
}

func NewPole() *Pole {
	np := &Pole{
		Object: tools.O(),
	}

	np.Ref = ""
	np.City = ""
	np.Address = ""
	np.Lat = 0.0
	np.Long = 0.0
	np.State = ""
	np.DtRef = ""
	np.DictRef = ""
	np.Height = 0
	np.Product = map[string]int{}
	np.DictInfo = map[string]string{}

	return np
}

func PoleFromJS(o *js.Object) *Pole {
	return &Pole{Object: o}
}

func (p *Pole) SwitchState() {
	switch p.State {
	case poleconst.StateNotSubmitted:
		p.State = poleconst.StateToDo
	case poleconst.StateToDo:
		p.State = poleconst.StateHoleDone
	case poleconst.StateHoleDone:
		p.State = poleconst.StateIncident
	case poleconst.StateIncident:
		p.State = poleconst.StateDone
	case poleconst.StateDone:
		p.State = poleconst.StateCancelled
	case poleconst.StateCancelled:
		p.State = poleconst.StateNotSubmitted
	}
}
