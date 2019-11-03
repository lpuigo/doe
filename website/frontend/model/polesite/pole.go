package polesite

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/model/polesite/poleconst"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
	"github.com/lpuig/ewin/doe/website/frontend/tools/json"
	"strconv"
)

// type Pole reflects backend/model/polesites.pole struct
type Pole struct {
	*js.Object
	Id             int      `js:"Id"`
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
	DictDate       string   `js:"DictDate"`
	DictInfo       string   `js:"DictInfo"`
	DaStartDate    string   `js:"DaStartDate"`
	DaEndDate      string   `js:"DaEndDate"`
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

	np.Id = -100
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
	np.DictDate = ""
	np.DictInfo = ""
	np.DaStartDate = ""
	np.DaEndDate = ""
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

func (p *Pole) Duplicate(suffix string, offset float64) *Pole {
	np := NewPole()
	np.Ref = p.Ref + suffix
	np.City = p.City
	np.Address = p.Address
	np.Sticker = p.Sticker
	np.Lat = p.Lat + offset
	np.Long = p.Long + offset
	np.State = p.State
	//np.Date = ""
	np.Actors = []string{}
	np.DtRef = p.DtRef
	np.DictRef = p.DictRef
	np.DictDate = p.DictDate
	np.DictInfo = p.DictInfo
	np.DaStartDate = p.DaStartDate
	np.DaEndDate = p.DaEndDate
	np.Height = p.Height
	np.Material = p.Material
	//np.AspiDate = ""
	//np.Kizeo = ""
	np.Comment = p.Comment
	np.Product = p.Product[:]
	//np.AttachmentDate = ""
	return np
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
	res := searchItem("", poleconst.FilterValueRef, p.GetTitle())
	res += searchItem(",", poleconst.FilterValueCity, p.City)
	res += searchItem(",", poleconst.FilterValueAddr, p.Address)
	res += searchItem(",", poleconst.FilterValueComment, p.Comment)
	res += searchItem(",", poleconst.FilterValueMaterial, p.Material)
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

func (p *Pole) Clone() *Pole {
	return &Pole{Object: json.Parse(json.Stringify(p))}
}

func (p *Pole) GetTitle() string {
	title := p.Ref
	if p.Sticker != "" {
		title += " " + p.Sticker
	}
	return title
}

func GetFilterTypeValueLabel() []*elements.ValueLabel {
	return []*elements.ValueLabel{
		elements.NewValueLabel(poleconst.FilterValueAll, poleconst.FilterLabelAll),
		elements.NewValueLabel(poleconst.FilterValueRef, poleconst.FilterLabelRef),
		elements.NewValueLabel(poleconst.FilterValueCity, poleconst.FilterLabelCity),
		elements.NewValueLabel(poleconst.FilterValueAddr, poleconst.FilterLabelAddr),
		elements.NewValueLabel(poleconst.FilterValueComment, poleconst.FilterLabelComment),
		elements.NewValueLabel(poleconst.FilterValueMaterial, poleconst.FilterLabelMaterial),
		elements.NewValueLabel(poleconst.FilterValueHeigth, poleconst.FilterLabelHeigth),
		elements.NewValueLabel(poleconst.FilterValueProduct, poleconst.FilterLabelProduct),
		elements.NewValueLabel(poleconst.FilterValueDt, poleconst.FilterLabelDt),
		elements.NewValueLabel(poleconst.FilterValueDict, poleconst.FilterLabelDict),
		elements.NewValueLabel(poleconst.FilterValueDictInfo, poleconst.FilterLabelDictInfo),
	}
}

func PoleStateLabel(state string) string {
	switch state {
	case poleconst.StateNotSubmitted:
		return poleconst.LabelNotSubmitted
	case poleconst.StateNoGo:
		return poleconst.LabelNoGo
	case poleconst.StateDictToDo:
		return poleconst.LabelDictToDo
	case poleconst.StateToDo:
		return poleconst.LabelToDo
	case poleconst.StateHoleDone:
		return poleconst.LabelHoleDone
	case poleconst.StateIncident:
		return poleconst.LabelIncident
	case poleconst.StateDone:
		return poleconst.LabelDone
	case poleconst.StateAttachment:
		return poleconst.LabelAttachment
	case poleconst.StateCancelled:
		return poleconst.LabelCancelled
	default:
		return "<" + state + ">"
	}
}

func GetStatesValueLabel() []*elements.ValueLabel {
	return []*elements.ValueLabel{
		elements.NewValueLabel(poleconst.StateNotSubmitted, poleconst.LabelNotSubmitted),
		elements.NewValueLabel(poleconst.StateNoGo, poleconst.LabelNoGo),
		elements.NewValueLabel(poleconst.StateDictToDo, poleconst.LabelDictToDo),
		elements.NewValueLabel(poleconst.StateToDo, poleconst.LabelToDo),
		elements.NewValueLabel(poleconst.StateHoleDone, poleconst.LabelHoleDone),
		elements.NewValueLabel(poleconst.StateIncident, poleconst.LabelIncident),
		elements.NewValueLabel(poleconst.StateDone, poleconst.LabelDone),
		elements.NewValueLabel(poleconst.StateAttachment, poleconst.LabelAttachment),
		elements.NewValueLabel(poleconst.StateCancelled, poleconst.LabelCancelled),
	}
}

func GetMaterialsValueLabel() []*elements.ValueLabel {
	return []*elements.ValueLabel{
		elements.NewValueLabel(poleconst.MaterialWood, poleconst.MaterialWood),
		elements.NewValueLabel(poleconst.MaterialMetal, poleconst.MaterialMetal),
		elements.NewValueLabel(poleconst.MaterialComp, poleconst.MaterialComp),
	}
}

func GetProductsValueLabel() []*elements.ValueLabel {
	return []*elements.ValueLabel{
		elements.NewValueLabel(poleconst.ProductCoated, poleconst.ProductCoated),
		elements.NewValueLabel(poleconst.ProductMoise, poleconst.ProductMoise),
		elements.NewValueLabel(poleconst.ProductCouple, poleconst.ProductCouple),
		elements.NewValueLabel(poleconst.ProductReplace, poleconst.ProductReplace),
		elements.NewValueLabel(poleconst.ProductRemove, poleconst.ProductRemove),
	}
}

func PoleRowClassName(status string) string {
	var res string = ""
	switch status {
	case poleconst.StateNotSubmitted:
		return "pole-row-not-submitted"
	case poleconst.StateNoGo:
		return "pole-row-nogo"
	case poleconst.StateDictToDo:
		return "pole-row-nogo"
	case poleconst.StateToDo:
		return "pole-row-todo"
	case poleconst.StateHoleDone:
		return "pole-row-hole-done"
	case poleconst.StateIncident:
		return "pole-row-incident"
	case poleconst.StateDone:
		return "pole-row-done"
	case poleconst.StateAttachment:
		return "pole-row-attachment"
	case poleconst.StateCancelled:
		return "pole-row-cancelled"

	default:
		res = "worksite-row-error"
	}
	return res
}
