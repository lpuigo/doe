package polesite

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/model/polesite/poleconst"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
	"strconv"
)

// type Pole reflects backend/model/polesites.pole struct
type Pole struct {
	*js.Object
	Ref            string   `js:"Ref"`
	City           string   `js:"City"`
	Address        string   `js:"Address"`
	Sticker        string   `js:"Sticker"`
	Lat            float64  `js:"Lat"`
	Long           float64  `js:"Long"`
	State          string   `js:"State"`
	Date           string   `js:"Date"`
	Actors         []string `js:"Actors"`
	DtRef          string   `js:"DtRef"`
	DictRef        string   `js:"DictRef"`
	DictInfo       string   `js:"DictInfo"`
	Height         int      `js:"Height"`
	Material       string   `js:"Material"`
	AspiDate       string   `js:"AspiDate"`
	Kizeo          string   `js:"Kizeo"`
	Comment        string   `js:"Comment"`
	Product        []string `js:"Product"`
	AttachmentDate string   `js:"AttachmentDate"`
}

func NewPole() *Pole {
	np := &Pole{
		Object: tools.O(),
	}

	np.Ref = ""
	np.City = ""
	np.Address = ""
	np.Sticker = ""
	np.Lat = 0.0
	np.Long = 0.0
	np.State = ""
	np.Date = ""
	np.Actors = []string{}
	np.DtRef = ""
	np.DictRef = ""
	np.DictInfo = ""
	np.Height = 0
	np.Material = ""
	np.AspiDate = ""
	np.Kizeo = ""
	np.Comment = ""
	np.Product = []string{}
	np.AttachmentDate = ""

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

func (p *Pole) SearchString(filter string) string {
	searchItem := func(prefix, typ, value string) string {
		if value == "" {
			return ""
		}
		if filter != poleconst.FilterValueAll && filter != typ {
			return ""
		}
		return prefix + typ + value
	}

	res := searchItem("", poleconst.FilterValueRef, p.Ref)
	res += searchItem(",", poleconst.FilterValueCity, p.City)
	res += searchItem(",", poleconst.FilterValueComment, p.Comment)
	res += searchItem(",", poleconst.FilterValueHeigth, strconv.Itoa(p.Height))
	if len(p.Product) > 0 {
		for _, prd := range p.Product {
			res += searchItem(",", poleconst.FilterValueProduct, prd)
		}
	}
	res += searchItem(",", poleconst.FilterValueDt, p.DtRef)
	res += searchItem(",", poleconst.FilterValueDict, p.DictRef)
	res += searchItem(",", poleconst.FilterValueDictInfo, p.DictInfo)
	return res
}

func GetFilterTypeValueLabel() []*elements.ValueLabel {
	return []*elements.ValueLabel{
		elements.NewValueLabel(poleconst.FilterValueAll, poleconst.FilterLabelAll),
		elements.NewValueLabel(poleconst.FilterValueRef, poleconst.FilterLabelRef),
		elements.NewValueLabel(poleconst.FilterValueCity, poleconst.FilterLabelCity),
		elements.NewValueLabel(poleconst.FilterValueComment, poleconst.FilterLabelComment),
		elements.NewValueLabel(poleconst.FilterValueHeigth, poleconst.FilterLabelHeigth),
		elements.NewValueLabel(poleconst.FilterValueProduct, poleconst.FilterLabelProduct),
		elements.NewValueLabel(poleconst.FilterValueDt, poleconst.FilterLabelDt),
		elements.NewValueLabel(poleconst.FilterValueDict, poleconst.FilterLabelDict),
		elements.NewValueLabel(poleconst.FilterValueDictInfo, poleconst.FilterLabelDictInfo),
	}
}
