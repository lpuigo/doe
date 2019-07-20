package polesite

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
)

const (
	PoleStateNotSubmitted string = "00 Not Submitted"
	PoleStateToDo         string = "10 To Do"
	PoleStateHoleDone     string = "20 Hole Done"
	PoleStateIncident     string = "25 Incident"
	PoleStateDone         string = "90 Done"
	PoleStateCancelled    string = "99 Cancelled"
)

func GetStatesValueLabel() []*elements.ValueLabel {
	return []*elements.ValueLabel{
		elements.NewValueLabel(PoleStateNotSubmitted, "Non soumis"),
		elements.NewValueLabel(PoleStateToDo, "A faire"),
		elements.NewValueLabel(PoleStateHoleDone, "Trou fait"),
		elements.NewValueLabel(PoleStateIncident, "Incident"),
		elements.NewValueLabel(PoleStateDone, "Fait"),
		elements.NewValueLabel(PoleStateCancelled, "Annulé"),
	}
}

type Pole struct {
	*js.Object
	Ref   string  `js:"Ref"`
	City  string  `js:"City"`
	Lat   float64 `js:"Lat"`
	Long  float64 `js:"Long"`
	State string  `js:"State"`
}

func NewPole(pole BePole) *Pole {
	np := &Pole{
		Object: tools.O(),
	}

	np.Ref = pole.Ref
	np.City = pole.City
	np.Lat = pole.Lat
	np.Long = pole.Long
	np.State = pole.State

	return np
}

func (p *Pole) SwitchState() {
	switch p.State {
	case PoleStateNotSubmitted:
		p.State = PoleStateToDo
	case PoleStateToDo:
		p.State = PoleStateHoleDone
	case PoleStateHoleDone:
		p.State = PoleStateIncident
	case PoleStateIncident:
		p.State = PoleStateDone
	case PoleStateDone:
		p.State = PoleStateCancelled
	case PoleStateCancelled:
		p.State = PoleStateNotSubmitted
	}
}

func GetCenterAndBounds(poles []*Pole) (clat, clong, blat1, blong1, blat2, blong2 float64) {
	if len(poles) == 0 {
		return 47, 5, 46.5, 4.5, 47.5, 5.5
	}

	min := func(pole *Pole) {
		if pole.Lat < blat1 {
			blat1 = pole.Lat
		}
		if pole.Long < blong1 {
			blong1 = pole.Long
		}
	}

	max := func(pole *Pole) {
		if pole.Lat > blat2 {
			blat2 = pole.Lat
		}
		if pole.Long > blong2 {
			blong2 = pole.Long
		}
	}

	blat1, blong1 = 500, 500
	for _, pole := range poles {
		clat += pole.Lat
		clong += pole.Long
		min(pole)
		max(pole)
	}

	nb := float64(len(poles))
	clat /= nb
	clong /= nb
	return
}

func GenPoles(poles []BePole) []*Pole {
	res := make([]*Pole, len(poles))

	for i, pole := range poles {
		res[i] = NewPole(pole)
	}
	return res
}
