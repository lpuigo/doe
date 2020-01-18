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
	np.Ref = p.Ref
	np.City = p.City
	np.Address = p.Address
	np.Sticker = p.Sticker + suffix
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

func (p *Pole) IsToDo() bool {
	switch p.State {
	case poleconst.StateNotSubmitted:
		return false
	//case poleconst.StateNoGo:
	//case poleconst.StateDictToDo:
	//case poleconst.StateToDo:
	//case poleconst.StateHoleDone:
	//case poleconst.StateIncident:
	//case poleconst.StateDone:
	//case poleconst.StateAttachment:
	case poleconst.StateCancelled:
		return false
	default:
		return true
	}
}

func (p *Pole) IsDone() bool {
	return p.State == poleconst.StateDone
}

func (p *Pole) IsBlocked() bool {
	switch p.State {
	case poleconst.StateNoGo:
		return true
	case poleconst.StateDictToDo:
		return true
	case poleconst.StateIncident:
		return true
	default:
		return false
	}
}

func (p *Pole) IsAttachment() bool {
	return p.State == poleconst.StateAttachment
}

func (p *Pole) SetAttachmentDate(d string) {
	p.AttachmentDate = d
	p.State = poleconst.StateAttachment
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

func GetStatesValueLabel(showAttachment bool) []*elements.ValueLabelDisabled {
	return []*elements.ValueLabelDisabled{
		elements.NewValueLabelDisabled(poleconst.StateNotSubmitted, poleconst.LabelNotSubmitted, false),
		elements.NewValueLabelDisabled(poleconst.StateNoGo, poleconst.LabelNoGo, false),
		elements.NewValueLabelDisabled(poleconst.StateDictToDo, poleconst.LabelDictToDo, false),
		elements.NewValueLabelDisabled(poleconst.StateToDo, poleconst.LabelToDo, false),
		elements.NewValueLabelDisabled(poleconst.StateHoleDone, poleconst.LabelHoleDone, false),
		elements.NewValueLabelDisabled(poleconst.StateIncident, poleconst.LabelIncident, false),
		elements.NewValueLabelDisabled(poleconst.StateDone, poleconst.LabelDone, false),
		elements.NewValueLabelDisabled(poleconst.StateAttachment, poleconst.LabelAttachment, !showAttachment),
		elements.NewValueLabelDisabled(poleconst.StateCancelled, poleconst.LabelCancelled, false),
	}
}

func GetMaterialsValueLabel() []*elements.ValueLabel {
	return []*elements.ValueLabel{
		elements.NewValueLabel(poleconst.MaterialWood, poleconst.MaterialWood),
		elements.NewValueLabel(poleconst.MaterialMetal, poleconst.MaterialMetal),
		elements.NewValueLabel(poleconst.MaterialEnforcedMetal, poleconst.MaterialEnforcedMetal),
		elements.NewValueLabel(poleconst.MaterialComp, poleconst.MaterialComp),
	}
}

func GetProductsValueLabel() []*elements.ValueLabel {
	return []*elements.ValueLabel{
		elements.NewValueLabel(poleconst.ProductCoated, poleconst.ProductCoated),
		elements.NewValueLabel(poleconst.ProductHandDigging, poleconst.ProductHandDigging),
		elements.NewValueLabel(poleconst.ProductMoise, poleconst.ProductMoise),
		elements.NewValueLabel(poleconst.ProductCouple, poleconst.ProductCouple),
		elements.NewValueLabel(poleconst.ProductReplace, poleconst.ProductReplace),
		elements.NewValueLabel(poleconst.ProductTrickyReplace, poleconst.ProductTrickyReplace),
		elements.NewValueLabel(poleconst.ProductRemove, poleconst.ProductRemove),
		elements.NewValueLabel(poleconst.ProductVacuumTruck, poleconst.ProductVacuumTruck),
		elements.NewValueLabel(poleconst.ProductReplenishment, poleconst.ProductReplenishment),
		elements.NewValueLabel(poleconst.ProductFarReplenishment, poleconst.ProductFarReplenishment),
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
