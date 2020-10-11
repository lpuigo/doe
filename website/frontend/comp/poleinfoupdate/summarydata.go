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
	GetColumn func(*polesite.Pole) ([]string, []int)
}

func NewSummarizer() *Summarizer {
	s := &Summarizer{Object: tools.O()}
	s.Colums = []string{}
	s.SummaryDatas = []*SummaryData{}

	s.GetLine = func(pole *polesite.Pole) string {
		return pole.City
	}
	s.GetColumn = func(pole *polesite.Pole) ([]string, []int) {
		state := pole.State
		if state == poleconst.StateDenseNetwork || state == poleconst.StateNoAccess {
			state = poleconst.StateToDo
		}
		return []string{state}, []int{1}
	}
	return s
}

const (
	IgnoreColumn string = "Ignore"
)

func GetPoleType(pole *polesite.Pole) ([]string, []int) {
	if !pole.IsToBeDone() {
		return []string{IgnoreColumn}, []int{0}
	}
	height := strconv.Itoa(pole.Height)
	crea := pole.HasProduct(poleconst.ProductCreation)
	coupl := pole.HasProduct(poleconst.ProductCouple)
	moise := pole.HasProduct(poleconst.ProductMoise)
	nb := 0
	if crea || coupl || moise {
		nb = 1
	}
	if crea && (coupl || moise) {
		nb = 2
	}
	if nb == 0 {
		return []string{IgnoreColumn}, []int{0}
	}
	poles := []string{}
	nbs := []int{}
	switch pole.Material {
	case "Bois":
		poles = append(poles, "BS"+height)
		nbs = append(nbs, nb)
	case "Métal":
		poles = append(poles, "MS"+height)
		nbs = append(nbs, nb)
	case "Métal Renforcé":
		poles = append(poles, "MF"+height)
		nbs = append(nbs, nb)
	case "Composite":
		poles = append(poles, "FS"+height)
		nbs = append(nbs, nb)
	case "Composite Renforcé":
		poles = append(poles, "FR"+height)
		nbs = append(nbs, nb)
	default:
		materl := pole.Material
		if materl == "" {
			materl = "A définir "
		}
		poles = append(poles, materl+height)
		nbs = append(nbs, nb)
	}
	return poles, nbs
}

func GetPoleItem(pole *polesite.Pole) ([]string, []int) {
	if !pole.IsToBeDone() {
		return []string{IgnoreColumn}, []int{0}
	}
	coupl := pole.HasProduct(poleconst.ProductCouple)
	moise := pole.HasProduct(poleconst.ProductMoise)
	mat := ""
	switch pole.Material {
	case "Bois":
		mat = "Bois"
	case "Métal":
		mat = "Mét."
	case "Métal Renforcé":
		mat = "Mét."
	case "Composite":
		mat = "Comp."
	case "Composite Renforcé":
		mat = "Comp."
	default:
		mat = "à définir"
	}
	poles := []string{}
	nbs := []int{}
	if pole.HasProduct(poleconst.ProductHauban) {
		poles = append(poles, "Hauban "+mat)
		nbs = append(nbs, 1)
	}
	if moise {
		poles = append(poles, "Mois. "+mat)
		nbs = append(nbs, 1)
	}
	if coupl {
		poles = append(poles, "Entr. "+mat)
		nbs = append(nbs, 1)
	}
	return poles, nbs
}

func GetPoleAction(pole *polesite.Pole) ([]string, []int) {
	if !pole.IsToBeDone() {
		return []string{IgnoreColumn}, []int{0}
	}
	nb := 1
	create := pole.HasProduct(poleconst.ProductCreation)
	couple := pole.HasProduct(poleconst.ProductCouple)
	moise := pole.HasProduct(poleconst.ProductMoise)
	repl := pole.HasProduct(poleconst.ProductReplace)
	haub := pole.HasProduct(poleconst.ProductHauban)
	switch {
	case repl && create && couple && haub:
		return []string{"Rempl. Couplé Haub."}, []int{nb}
	case repl && create && couple:
		return []string{"Rempl. Couplé"}, []int{nb}
	case repl && create && moise:
		return []string{"Rempl. Moisé"}, []int{nb}
	case repl && create:
		return []string{"Remplacement"}, []int{nb}
	case create && couple && haub:
		return []string{"Impl. Couplé Haub."}, []int{nb}
	case create && couple:
		return []string{"Impl. Couplé"}, []int{nb}
	case create && moise:
		return []string{"Impl. Moisé"}, []int{nb}
	case create:
		return []string{"Implantation"}, []int{nb}
	case couple:
		return []string{"Renf. Couplé"}, []int{nb}
	case haub:
		return []string{"Renf. Hauban"}, []int{nb}
	default:
		return []string{"A définir"}, []int{nb}
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
		sd, found := summDatas[line]
		if !found {
			sd = NewSummaryData(line)
			summDatas[line] = sd
		}
		columns, amounts := s.GetColumn(pole)
		for i, column := range columns {
			columnSet[column] = true
			//sd.NbPoles[pole.State]++
			nb := sd.Get("NbPoles").Get(column).Int()
			sd.Get("NbPoles").Set(column, nb+amounts[i])
		}
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
