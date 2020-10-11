package poleinfoupdate

import (
	"sort"
	"strconv"

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
	GetColumn func(*polesite.Pole) (string, int)
}

func NewSummarizer() *Summarizer {
	s := &Summarizer{Object: tools.O()}
	s.Colums = []string{}
	s.SummaryDatas = []*SummaryData{}

	s.GetLine = func(pole *polesite.Pole) string {
		return pole.City
	}
	s.GetColumn = func(pole *polesite.Pole) (string, int) {
		state := pole.State
		if state == poleconst.StateDenseNetwork || state == poleconst.StateNoAccess {
			state = poleconst.StateToDo
		}
		return state, 1
	}
	return s
}

const (
	IgnoreColumn string = "Ignore"
)

func GetPoleType(pole *polesite.Pole) (string, int) {
	if !pole.IsToBeDone() {
		return IgnoreColumn, 0
	}
	height := strconv.Itoa(pole.Height)
	nb := 1
	if pole.IsDoublePole() {
		nb = 2
	}
	switch pole.Material {
	case "Bois":
		return "BS " + height, nb
	case "Métal":
		return "MS " + height, nb
	case "Métal Renforcé":
		return "MF " + height, nb
	case "Composite":
		return "FS " + height, nb
	case "Composite Renforcé":
		return "FR " + height, nb
	default:
		mat := pole.Material
		if mat == "" {
			mat = "A définir "
		}
		return mat + height, nb
	}
}

func GetPoleAction(pole *polesite.Pole) (string, int) {
	if !pole.IsToBeDone() {
		return IgnoreColumn, 0
	}
	nb := 1
	create := pole.HasProduct(poleconst.ProductCreation)
	couple := pole.HasProduct(poleconst.ProductCouple)
	moise := pole.HasProduct(poleconst.ProductMoise)
	repl := pole.HasProduct(poleconst.ProductReplace)
	haub := pole.HasProduct(poleconst.ProductHauban)
	switch {
	case repl && create && couple && haub:
		return "Rpl. Couplé Haubané", nb
	case repl && create && couple:
		return "Rempl. Couplé", nb
	case repl && create && moise:
		return "Rempl. Moisé", nb
	case repl && create:
		return "Remplacement", nb
	case create && couple && haub:
		return "Impl. Couplé Haubané", nb
	case create && couple:
		return "Impl. Couplé", nb
	case create && moise:
		return "Impl. Moisé", nb
	case create:
		return "Implantation", nb
	case haub:
		return "Renf. Hauban", nb
	default:
		return "A définir", nb
	}
}

func (s *Summarizer) GetCalcColumns() []string {
	res := []string{}
	for _, column := range s.Colums {
		if column == IgnoreColumn {
			continue
		}
		res = append(res, column)
	}
	sort.Strings(res)
	return res
}

func (s *Summarizer) Calc(poles []*polesite.Pole) {
	//lineSet := make(map[string]bool)
	columnSet := make(map[string]bool)
	summDatas := map[string]*SummaryData{}

	for _, pole := range poles {
		line := s.GetLine(pole)
		//lineSet[line] = true
		column, amount := s.GetColumn(pole)
		columnSet[column] = true

		sd, found := summDatas[line]
		if !found {
			sd = NewSummaryData(line)
			summDatas[line] = sd
		}

		//sd.NbPoles[pole.State]++
		nb := sd.Get("NbPoles").Get(column).Int()
		sd.Get("NbPoles").Set(column, nb+amount)
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
