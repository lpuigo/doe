package poleinfoupdate

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/model/polesite"
	"github.com/lpuig/ewin/doe/website/frontend/model/polesite/poleconst"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

type SummaryData struct {
	*js.Object

	City    string         `js:"City"`
	NbPoles map[string]int `js:"NbPoles"`
}

func NewSummaryData(city string) *SummaryData {
	sd := &SummaryData{Object: tools.O()}
	sd.City = city
	sd.NbPoles = make(map[string]int)
	return sd
}

func CalcSummaryDatas(poles []*polesite.Pole, interestingStatuses []string) []*SummaryData {
	statuses := map[string]bool{}
	for _, status := range interestingStatuses {
		statuses[status] = true
	}

	summDataByCity := map[string]*SummaryData{}
	for _, pole := range poles {
		state := pole.State
		if state == poleconst.StateDenseNetwork || state == poleconst.StateNoAccess {
			state = poleconst.StateToDo
		}
		if !statuses[state] {
			continue
		}
		sd, found := summDataByCity[pole.City]
		if !found {
			sd = NewSummaryData(pole.City)
			summDataByCity[pole.City] = sd
		}
		//sd.NbPoles[pole.State]++
		nb := sd.Get("NbPoles").Get(state).Int()
		sd.Get("NbPoles").Set(state, nb+1)
	}

	res := []*SummaryData{}
	for _, sd := range summDataByCity {
		res = append(res, sd)
	}
	return res
}
