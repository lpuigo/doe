package polesite

import (
	"github.com/lpuig/ewin/doe/website/frontend/model/polesite/poleconst"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
)

func GetStatesValueLabel() []*elements.ValueLabel {
	return []*elements.ValueLabel{
		elements.NewValueLabel(poleconst.StateNotSubmitted, poleconst.LabelNotSubmitted),
		elements.NewValueLabel(poleconst.StateToDo, poleconst.LabelToDo),
		elements.NewValueLabel(poleconst.StateHoleDone, poleconst.LabelHoleDone),
		elements.NewValueLabel(poleconst.StateIncident, poleconst.LabelIncident),
		elements.NewValueLabel(poleconst.StateDone, poleconst.LabelDone),
		elements.NewValueLabel(poleconst.StateCancelled, poleconst.LabelCancelled),
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
