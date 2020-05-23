package polesite

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/model/polesite/poleconst"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/date"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
	"github.com/lpuig/ewin/doe/website/frontend/tools/json"
	"strconv"
	"strings"
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
	DaQueryDate    string   `js:"DaQueryDate"`
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
	np.DaQueryDate = ""
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

func (p *Pole) Duplicate(newname string, offset float64) *Pole {
	np := NewPole()
	np.Ref = p.Ref
	np.City = p.City
	np.Address = p.Address
	np.Sticker = newname
	np.Lat = p.Lat + offset
	np.Long = p.Long + offset
	np.State = p.State
	//np.Date = ""
	np.Actors = []string{}
	np.DtRef = p.DtRef
	np.DictRef = p.DictRef
	np.DictDate = p.DictDate
	np.DictInfo = p.DictInfo
	np.DaQueryDate = p.DaQueryDate
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

//func (p *Pole) SwitchState() {
//	switch p.State {
//	case poleconst.StateNotSubmitted:
//		p.State = poleconst.StateToDo
//	case poleconst.StateToDo:
//		p.State = poleconst.StateHoleDone
//	case poleconst.StateHoleDone:
//		p.State = poleconst.StateIncident
//	case poleconst.StateIncident:
//		p.State = poleconst.StateDone
//	case poleconst.StateDone:
//		p.State = poleconst.StateCancelled
//	case poleconst.StateCancelled:
//		p.State = poleconst.StateNotSubmitted
//	}
//}

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
	case poleconst.StateNoGo:
		return false
	//case poleconst.StateDictToDo:
	//case poleconst.StateDaToDo:
	//case poleconst.StateDaExpected:
	//case poleconst.StatePermissionPending:
	//case poleconst.StateToDo:
	//case poleconst.StateNoAccess:
	//case poleconst.StateDenseNetwork:
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

//func (p *Pole) IsDoable() bool {
//	switch p.State {
//	case poleconst.StateNotSubmitted:
//		return false
//	case poleconst.StateNoGo:
//		return false
//	case poleconst.StateDictToDo:
//		return false
//	case poleconst.StateDaToDo:
//		return false
//	case poleconst.StateDaExpected:
//		return false
//	//case poleconst.StatePermissionPending:
//	//case poleconst.StateToDo:
//	//case poleconst.StateNoAccess:
//	//case poleconst.StateDenseNetwork:
//	//case poleconst.StateHoleDone:
//	//case poleconst.StateIncident:
//	case poleconst.StateDone:
//		return false
//	case poleconst.StateAttachment:
//		return false
//	case poleconst.StateCancelled:
//		return false
//	default:
//		return true
//	}
//}
//
//func (p *Pole) IsSettingUp() bool {
//	switch p.State {
//	//case poleconst.StateNotSubmitted:
//	//case poleconst.StateNoGo:
//	//case poleconst.StateDictToDo:
//	case poleconst.StateDaToDo:
//		return true
//	case poleconst.StateDaExpected:
//		return true
//	//case poleconst.StatePermissionPending:
//	//case poleconst.StateToDo:
//	//case poleconst.StateNoAccess:
//	//case poleconst.StateDenseNetwork:
//	//case poleconst.StateHoleDone:
//	//case poleconst.StateIncident:
//	//case poleconst.StateDone:
//	//case poleconst.StateAttachment:
//	//case poleconst.StateCancelled:
//	default:
//		return false
//	}
//}

func (p *Pole) IsInStateToBeChecked() bool {
	switch p.State {
	case poleconst.StateDictToDo:
		return true
	case poleconst.StateDaToDo:
		return true
	case poleconst.StateDaExpected:
		return true
	case poleconst.StatePermissionPending:
		return true
	case poleconst.StateToDo:
		return true
	case poleconst.StateNoAccess:
		return true
	case poleconst.StateDenseNetwork:
		return true
	case poleconst.StateHoleDone:
		return true
	case poleconst.StateIncident:
		return true
	default:
		return false
	}
}

func (p *Pole) IsBlocked() bool {
	switch p.State {
	//case poleconst.StateNoGo:
	//	return true
	case poleconst.StateDictToDo:
		return true
	case poleconst.StateDaToDo:
		return true
	case poleconst.StateDaExpected:
		return true
	case poleconst.StatePermissionPending:
		return true
	case poleconst.StateIncident:
		return true
	default:
		return false
	}
}

func (p *Pole) IsDone() bool {
	return p.State == poleconst.StateDone
}

func (p *Pole) IsAttachment() bool {
	return p.State == poleconst.StateAttachment
}

func (p *Pole) SetAttachmentDate(d string) {
	p.AttachmentDate = d
	p.State = poleconst.StateAttachment
}

// CheckState updates pole state according to DICT and DA value (depending on current date) if pole is Actually To Do
func (p *Pole) CheckState() {
	today := date.TodayAfter(0)

	testDICTandDAdates := func() {
		// check if DICT if to be done => StateDictToDo
		if tools.Empty(p.DictRef) {
			p.SetState(poleconst.StateDictToDo)
			return
		}
		if !(!tools.Empty(p.DictDate) && today <= date.After(p.DictDate, poleconst.DictValidityDuration)) {
			p.SetState(poleconst.StateDictToDo)
			return
		}
		// check if DAC is to be done => StateDaToDo
		if tools.Empty(p.DaQueryDate) {
			p.SetState(poleconst.StateDaToDo)
			return
		}
		// DAC query is Set. Check DaStart and End Dates
		if !tools.Empty(p.DaStartDate) && !tools.Empty(p.DaEndDate) && today > p.DaEndDate {
			p.SetState(poleconst.StateDaToDo)
			return
		}
		if tools.Empty(p.DaStartDate) {
			p.SetState(poleconst.StateDaExpected)
			return
		}
		// DA dates are not applicable => StatePermissionPending
		if !(p.DictDate <= today && p.DaStartDate <= today) {
			p.SetState(poleconst.StatePermissionPending)
			return
		}
		// All checked => StateToDo
		p.SetState(poleconst.StateToDo)
	}

	if p.IsInStateToBeChecked() {
		testDICTandDAdates()
	}
}

func (p *Pole) SetState(state string) {
	switch state {
	case poleconst.StateToDo:
		switch {
		case p.HasProduct(poleconst.ProductDenseNetwork):
			p.State = poleconst.StateDenseNetwork
		case p.HasProduct(poleconst.ProductNoAccess):
			p.State = poleconst.StateNoAccess
		default:
			p.State = state
		}
	case poleconst.StateDenseNetwork:
		p.AddProduct(poleconst.ProductDenseNetwork)
		p.State = state
	case poleconst.StateNoAccess:
		p.AddProduct(poleconst.ProductNoAccess)
		p.State = state
	default:
		p.State = state
	}
}

func (p *Pole) HasProduct(prd string) bool {
	for _, product := range p.Product {
		if product == prd {
			return true
		}
	}
	return false
}

func (p *Pole) AddProduct(prd string) {
	if p.HasProduct(prd) {
		return
	}
	p.Product = append(p.Product, prd)
	p.Get("Product").Call("sort")
}

func (p *Pole) CheckProductConsistency() {
	if p.HasProduct(poleconst.ProductTrickyReplace) {
		p.AddProduct(poleconst.ProductReplace)
	}
}

func (p *Pole) CheckInfoConsistency() {
	// City
	p.City = strings.Trim(strings.Title(strings.ToLower(p.City)), " ")
	// DICT & DT
	p.DictRef = strings.Trim(strings.ToUpper(p.DictRef), " ")
	p.DtRef = strings.Trim(strings.ToUpper(p.DtRef), " ")
	// DA dates
	if tools.Empty(p.DaQueryDate) && !tools.Empty(p.DaEndDate) {
		if tools.Empty(p.DaStartDate) {
			p.DaQueryDate = p.DaEndDate
			p.DaEndDate = ""
		} else {
			p.DaQueryDate = p.DaStartDate
		}
	}
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
	case poleconst.StateDaToDo:
		return poleconst.LabelDaToDo
	case poleconst.StateDaExpected:
		return poleconst.LabelDaExpected
	case poleconst.StatePermissionPending:
		return poleconst.LabelPermissionPending
	case poleconst.StateToDo:
		return poleconst.LabelToDo
	case poleconst.StateNoAccess:
		return poleconst.LabelNoAccess
	case poleconst.StateDenseNetwork:
		return poleconst.LabelDenseNetwork
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
		elements.NewValueLabelDisabled(poleconst.StateToDo, poleconst.LabelToDo, false),
		elements.NewValueLabelDisabled(poleconst.StateNoAccess, poleconst.LabelNoAccess, false),
		elements.NewValueLabelDisabled(poleconst.StateDenseNetwork, poleconst.LabelDenseNetwork, false),
		elements.NewValueLabelDisabled(poleconst.StateHoleDone, poleconst.LabelHoleDone, false),
		elements.NewValueLabelDisabled(poleconst.StateIncident, poleconst.LabelIncident, false),
		elements.NewValueLabelDisabled(poleconst.StateDone, poleconst.LabelDone, false),
		elements.NewValueLabelDisabled(poleconst.StateCancelled, poleconst.LabelCancelled, false),
		elements.NewValueLabelDisabled(poleconst.StateAttachment, poleconst.LabelAttachment, !showAttachment),
		elements.NewValueLabelDisabled(poleconst.StateDictToDo, poleconst.LabelDictToDo, true),
		elements.NewValueLabelDisabled(poleconst.StateDaToDo, poleconst.LabelDaToDo, true),
		elements.NewValueLabelDisabled(poleconst.StateDaExpected, poleconst.LabelDaExpected, true),
		elements.NewValueLabelDisabled(poleconst.StatePermissionPending, poleconst.LabelPermissionPending, true),
	}
}

func GetMaterialsValueLabel() []*elements.ValueLabel {
	return []*elements.ValueLabel{
		elements.NewValueLabel(poleconst.MaterialWood, poleconst.MaterialWood),
		elements.NewValueLabel(poleconst.MaterialMetal, poleconst.MaterialMetal),
		elements.NewValueLabel(poleconst.MaterialEnforcedMetal, poleconst.MaterialEnforcedMetal),
		elements.NewValueLabel(poleconst.MaterialComp, poleconst.MaterialComp),
		elements.NewValueLabel(poleconst.MaterialEnforcedComp, poleconst.MaterialEnforcedComp),
	}
}

func GetProductsValueLabel() []*elements.ValueLabel {
	return []*elements.ValueLabel{
		elements.NewValueLabel(poleconst.ProductCreation, poleconst.ProductCreation),
		elements.NewValueLabel(poleconst.ProductCoated, poleconst.ProductCoated),
		elements.NewValueLabel(poleconst.ProductHandDigging, poleconst.ProductHandDigging),
		elements.NewValueLabel(poleconst.ProductMoise, poleconst.ProductMoise),
		elements.NewValueLabel(poleconst.ProductCouple, poleconst.ProductCouple),
		elements.NewValueLabel(poleconst.ProductHauban, poleconst.ProductHauban),
		elements.NewValueLabel(poleconst.ProductReplace, poleconst.ProductReplace),
		elements.NewValueLabel(poleconst.ProductTrickyReplace, poleconst.ProductTrickyReplace),
		elements.NewValueLabel(poleconst.ProductRemove, poleconst.ProductRemove),
		elements.NewValueLabel(poleconst.ProductNoAccess, poleconst.ProductNoAccess),
		elements.NewValueLabel(poleconst.ProductDenseNetwork, poleconst.ProductDenseNetwork),
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
		return "pole-row-not-submitted"
	case poleconst.StateDictToDo:
		return "pole-row-nogo"
	case poleconst.StateDaToDo:
		return "pole-row-nogo"
	case poleconst.StateDaExpected:
		return "pole-row-nogo"
	case poleconst.StatePermissionPending:
		return "pole-row-nogo"
	case poleconst.StateToDo:
		return "pole-row-todo"
	case poleconst.StateNoAccess:
		return "pole-row-todo-tricky"
	case poleconst.StateDenseNetwork:
		return "pole-row-todo-tricky"
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
