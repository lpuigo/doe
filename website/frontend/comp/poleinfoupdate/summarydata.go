package poleinfoupdate

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/model/polesite"
	"github.com/lpuig/ewin/doe/website/frontend/model/polesite/poleconst"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

type SummaryData struct {
	*js.Object

	Line    string         `js:"Line"`
	NbPoles map[string]int `js:"NbPoles"`
}

func NewSummaryData(line string) *SummaryData {
	sd := &SummaryData{Object: tools.O()}
	sd.Line = line
	sd.NbPoles = make(map[string]int)
	return sd
}

func SummaryDataFromJS(obj *js.Object) *SummaryData {
	return &SummaryData{Object: obj}
}

type Summarizer struct {
	*js.Object

	Colums       []string       `js:"Colums"`
	SummaryDatas []*SummaryData `js:"SummaryDatas"`

	GetLine   func(*polesite.Pole) string
	GetColumn func(*polesite.Pole) string
}

func NewSummarizer() *Summarizer {
	s := &Summarizer{Object: tools.O()}
	s.Colums = []string{}
	s.SummaryDatas = []*SummaryData{}

	s.GetLine = func(pole *polesite.Pole) string {
		return pole.City
	}
	s.GetColumn = func(pole *polesite.Pole) string {
		state := pole.State
		if state == poleconst.StateDenseNetwork || state == poleconst.StateNoAccess {
			state = poleconst.StateToDo
		}
		return state
	}
	return s
}

func (s *Summarizer) Calc(poles []*polesite.Pole) {
	//lineSet := make(map[string]bool)
	columnSet := make(map[string]bool)
	summDatas := map[string]*SummaryData{}

	for _, pole := range poles {
		line := s.GetLine(pole)
		//lineSet[line] = true
		column := s.GetColumn(pole)
		columnSet[column] = true

		sd, found := summDatas[line]
		if !found {
			sd = NewSummaryData(line)
			summDatas[line] = sd
		}

		//sd.NbPoles[pole.State]++
		nb := sd.Get("NbPoles").Get(column).Int()
		sd.Get("NbPoles").Set(column, nb+1)
	}

	for column, _ := range columnSet {
		s.Colums = append(s.Colums, column)
	}
	for _, sd := range summDatas {
		s.SummaryDatas = append(s.SummaryDatas, sd)
	}
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
