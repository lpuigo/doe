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
	Priority       int      `js:"Priority"`
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
	DaValidation   bool     `js:"DaValidation"`
	DaStartDate    string   `js:"DaStartDate"`
	DaEndDate      string   `js:"DaEndDate"`
	Height         int      `js:"Height"`
	Material       string   `js:"Material"`
	AspiDate       string   `js:"AspiDate"`
	Kizeo          string   `js:"Kizeo"`
	Comment        string   `js:"Comment"`
	Product        []string `js:"Product"`
	AttachmentDate string   `js:"AttachmentDate"`
	TimeStamp      string   `js:"TimeStamp"`
}

func NewPole() *Pole {
	np := &Pole{
		Object: tools.O(),
	}

	np.Id = -100000
	np.Ref = ""
	np.City = ""
	np.Address = ""
	np.Sticker = ""
	np.Priority = 1
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
	np.DaValidation = false
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
	np.Priority = p.Priority
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
	np.DaValidation = p.DaValidation
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
	res += searchItem("", poleconst.FilterValuePrio, strconv.Itoa(p.Priority))
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
	//case poleconst.StateMarked:
	//case poleconst.StateNoAccess:
	//case poleconst.StateDenseNetwork:
	//case poleconst.StateHoleDone:
	//case poleconst.StateIncident:
	//case poleconst.StateDone:
	//case poleconst.StateAttachment:
	case poleconst.StateCancelled:
		return false
	case poleconst.StateDeleted:
		return false
	default:
		return true
	}
}

// IsToBeDone returns true if pole receiver is yet to be done (ie state not canceled, or already done)
func (p *Pole) IsToBeDone() bool {
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
	//case poleconst.StateMarked:
	//case poleconst.StateNoAccess:
	//case poleconst.StateDenseNetwork:
	//case poleconst.StateHoleDone:
	//case poleconst.StateIncident:
	case poleconst.StateDone:
		return false
	case poleconst.StateAttachment:
		return false
	case poleconst.StateCancelled:
		return false
	case poleconst.StateDeleted:
		return false
	default:
		return true
	}
}

// IsAlreadyDone returns true if pole receiver is done (ie has strate done or attachement)
func (p *Pole) IsAlreadyDone() bool {
	switch p.State {
	//case poleconst.StateNotSubmitted:
	//case poleconst.StateNoGo:
	//case poleconst.StateDictToDo:
	//case poleconst.StateDaToDo:
	//case poleconst.StateDaExpected:
	//case poleconst.StatePermissionPending:
	//case poleconst.StateToDo:
	//case poleconst.StateMarked:
	//case poleconst.StateNoAccess:
	//case poleconst.StateDenseNetwork:
	//case poleconst.StateHoleDone:
	//case poleconst.StateIncident:
	case poleconst.StateDone:
		return true
	case poleconst.StateAttachment:
		return true
	//case poleconst.StateCancelled:
	//case poleconst.StateDeleted:
	default:
		return false
	}
}

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
	case poleconst.StateMarked:
		return true
	case poleconst.StateNoAccess:
		return true
	case poleconst.StateDenseNetwork:
		return true
	//case poleconst.StateHoleDone:
	//	return true
	//case poleconst.StateIncident:
	//	return true
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

func (p *Pole) Deleted() bool {
	return p.State == poleconst.StateDeleted
}

func (p *Pole) IsAttachment() bool {
	return p.State == poleconst.StateAttachment
}

func (p *Pole) SetAttachmentDate(d string) {
	p.AttachmentDate = d
	p.State = poleconst.StateAttachment
}

// UpdateState updates pole receiver state.
//
// According to DICT and DA info (depending on current date) if pole is actually To be processed
//
// If State is StateDone, Date is checked and filled with current date if empty
func (p *Pole) UpdateState() {
	today := date.TodayAfter(0)

	if p.State == poleconst.StateIncident {
		p.AddProduct(poleconst.ProductIncident)
	}
	if p.State == poleconst.StateNoAccess {
		p.AddProduct(poleconst.ProductNoAccess)
	}

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
		if !(p.DaValidation && !tools.Empty(p.DaStartDate)) {
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
		return
	}
	if p.State == poleconst.StateDone {
		if tools.Empty(p.Date) {
			p.Date = date.TodayAfter(0)
		}
		if !tools.Empty(p.AttachmentDate) {
			p.State = poleconst.StateAttachment
		}
	}
}

// SetState sets pole receiver state.
//
// StateToDo can be modified according to already declared specific product (ProductDenseNetwork or ProductNoAccess)
func (p *Pole) SetState(state string) {
	switch state {
	case poleconst.StateToDo:
		switch {
		case p.State == poleconst.StateMarked:
			// do not change
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

func (p *Pole) IsDoublePole() bool {
	if !p.HasProduct(poleconst.ProductCreation) {
		return false
	}
	for _, product := range p.Product {
		if !(product != poleconst.ProductCouple && product != poleconst.ProductMoise) {
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

func (p *Pole) RemoveProduct(prd string) {
	for i, product := range p.Product {
		if product == prd {
			p.Product = append(p.Product[:i], p.Product[i+1:]...)
			return
		}
	}
}

// CheckProductConsistency checks for reciever pole's products consistency. Product can be added or discarded on pole receiver
func (p *Pole) CheckProductConsistency() {
	if p.HasProduct(poleconst.ProductTrickyReplace) {
		p.AddProduct(poleconst.ProductReplace)
	}
	if p.HasProduct(poleconst.ProductReplace) {
		p.AddProduct(poleconst.ProductCreation)
		p.RemoveProduct(poleconst.ProductStraighten)
	}
	if p.HasProduct(poleconst.ProductStraighten) {
		p.RemoveProduct(poleconst.ProductCreation)
		p.RemoveProduct(poleconst.ProductReplace)
	}
	if p.HasProduct(poleconst.ProductTightenHauban) {
		p.RemoveProduct(poleconst.ProductCreation)
		p.RemoveProduct(poleconst.ProductReplace)
		p.RemoveProduct(poleconst.ProductHauban)
	}
	ucComment := strings.ToUpper(p.Comment)
	if len(p.Product) == 0 && (strings.Contains(ucComment, "REDRESSEM") || strings.Contains(ucComment, "RECALAGE")) {
		p.AddProduct(poleconst.ProductStraighten)
	}
}

func (p *Pole) CheckInfoConsistency() {
	// Products
	p.CheckProductConsistency()
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
		elements.NewValueLabel(poleconst.FilterValuePrio, poleconst.FilterLabelPrio),
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
	case poleconst.StateMarked:
		return poleconst.LabelMarked
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
	case poleconst.StateDeleted:
		return poleconst.LabelDeleted
	default:
		return "<" + state + ">"
	}
}

func GetStatesValueLabel(showAttachment bool) []*elements.ValueLabelDisabled {
	return []*elements.ValueLabelDisabled{
		elements.NewValueLabelDisabled(poleconst.StateNotSubmitted, poleconst.LabelNotSubmitted, false),
		elements.NewValueLabelDisabled(poleconst.StateNoGo, poleconst.LabelNoGo, false),
		elements.NewValueLabelDisabled(poleconst.StateToDo, poleconst.LabelToDo, false),
		elements.NewValueLabelDisabled(poleconst.StateMarked, poleconst.LabelMarked, false),
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
		elements.NewValueLabel(poleconst.ProductCoated, poleconst.ProductCoated),
		elements.NewValueLabel(poleconst.ProductPruning, poleconst.ProductPruning),
		elements.NewValueLabel(poleconst.ProductHandDigging, poleconst.ProductHandDigging),
		elements.NewValueLabel(poleconst.ProductMechDigging, poleconst.ProductMechDigging),
		elements.NewValueLabel(poleconst.ProductTrickyReplace, poleconst.ProductTrickyReplace),
		elements.NewValueLabel(poleconst.ProductDenseNetwork, poleconst.ProductDenseNetwork),
		elements.NewValueLabel(poleconst.ProductNoAccess, poleconst.ProductNoAccess),
		elements.NewValueLabel(poleconst.ProductIncident, poleconst.ProductIncident),
		elements.NewValueLabel(poleconst.ProductCreation, poleconst.ProductCreation),
		elements.NewValueLabel(poleconst.ProductReplace, poleconst.ProductReplace),
		elements.NewValueLabel(poleconst.ProductCouple, poleconst.ProductCouple),
		elements.NewValueLabel(poleconst.ProductMoise, poleconst.ProductMoise),
		elements.NewValueLabel(poleconst.ProductHauban, poleconst.ProductHauban),
		elements.NewValueLabel(poleconst.ProductTightenHauban, poleconst.ProductTightenHauban),
		elements.NewValueLabel(poleconst.ProductStraighten, poleconst.ProductStraighten),
		elements.NewValueLabel(poleconst.ProductRemove, poleconst.ProductRemove),
		elements.NewValueLabel(poleconst.ProductInRow, poleconst.ProductInRow),
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
		return "pole-row-expectda"
	case poleconst.StatePermissionPending:
		return "pole-row-permission"
	case poleconst.StateToDo:
		return "pole-row-todo"
	case poleconst.StateMarked:
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
	case poleconst.StateDeleted:
		return "pole-row-deleted"

	default:
		res = "worksite-row-error"
	}
	return res
}

// GetPermissionDateRange returns a number (nb week) depending on receiver pole permissions dates
//
// - positive number : the pole will be ok in N week(s)
//
// - negative number : the pole is OK for N weeks
//
// - "NS" if result is not significant (ie pole aldready done, or not submited)
//
// - "BO" if Back Office work has yet to be done
func (p *Pole) GetPermissionDateRange(checkDA bool) string {
	waitFor := func() string {
		if tools.Empty(p.DictDate) {
			return "BO"
		}
		today := date.TodayAfter(0)
		nbDays := int(date.NbDaysBetween(today, p.DictDate))
		if checkDA {
			if !(p.DaValidation && !tools.Empty(p.DaStartDate)) {
				return "BO"
			}
			nbDaDays := int(date.NbDaysBetween(today, p.DaStartDate))
			if nbDaDays > nbDays {
				nbDays = nbDaDays
			}
		}
		if nbDays < 7 {
			nbDays = 7
		}
		return strconv.Itoa(nbDays / 7)
	}

	stillOkFor := func() string {
		if tools.Empty(p.DictDate) {
			return "BO"
		}
		today := date.TodayAfter(0)
		nbDays := int(date.NbDaysBetween(today, p.DictDate)) + poleconst.DictValidityDuration
		if checkDA {
			if !(p.DaValidation && !tools.Empty(p.DaEndDate)) {
				return "BO"
			}
			nbDaDays := int(date.NbDaysBetween(today, p.DaEndDate))
			if nbDaDays < nbDays {
				nbDays = nbDaDays
			}
		}
		return strconv.Itoa(-nbDays / 7)
	}

	// check for early decision
	if checkDA { // early decision depending on DA
		switch p.State {
		case poleconst.StateDaToDo:
			return "BO"
		case poleconst.StateDaExpected:
			return "BO"
		}
	}
	switch p.State {
	case poleconst.StateNotSubmitted:
		return "NS"
	case poleconst.StateNoGo:
		return "NS"
	case poleconst.StateDictToDo:
		return "BO"
	case poleconst.StateDaToDo:
		return waitFor()
	case poleconst.StateDaExpected:
		return waitFor()
	case poleconst.StatePermissionPending:
		return waitFor()
	case poleconst.StateToDo:
		return stillOkFor()
	case poleconst.StateMarked:
		return stillOkFor()
	case poleconst.StateNoAccess:
		return stillOkFor()
	case poleconst.StateDenseNetwork:
		return stillOkFor()
	case poleconst.StateHoleDone:
		return stillOkFor()
	case poleconst.StateIncident:
		return stillOkFor()
	case poleconst.StateDone:
		return "NS"
	case poleconst.StateAttachment:
		return "NS"
	case poleconst.StateCancelled:
		return "NS"
	case poleconst.StateDeleted:
		return "NS"
	}
	return "ERR"
}

func GetGroupNameByAge(name string) string {
	switch name {
	case "NS":
		return "Non Significatif"
	case "BO":
		return "A traiter par BO"
	case "ERR":
		return "Erreur"
	}
	nbWeeks, err := strconv.Atoi(name)
	if err != nil {
		return "Erreur " + name
	}
	if nbWeeks >= 0 {
		return "Dans " + strconv.Itoa(nbWeeks) + " sem"
	}
	return "OK pour " + strconv.Itoa(-nbWeeks) + " sem"
}

func (p *Pole) GetDateInfo() string {
	dictStart := func() string {
		if tools.Empty(p.DictDate) {
			return ""
		}
		return p.DictDate
	}

	daStart := func() string {
		if !(p.DaValidation && !tools.Empty(p.DaStartDate)) {
			return ""
		}
		return p.DaStartDate
	}

	dictEnd := func() string {
		if tools.Empty(p.DictDate) {
			return ""
		}
		return date.After(p.DictDate, poleconst.DictValidityDuration)
	}

	daEnd := func() string {
		if !(p.DaValidation && !tools.Empty(p.DaEndDate)) {
			return ""
		}
		return p.DaEndDate
	}

	startDate := func() string {
		begDict := dictStart()
		begDA := daStart()
		if begDict == "" {
			return ""
		}
		if begDA != "" && begDA > begDict {
			return begDA
		}
		return begDict
	}

	endDate := func() string {
		endDict := dictEnd()
		endDA := daEnd()
		if endDict == "" {
			return ""
		}
		if endDA != "" && endDA < endDict {
			return endDA
		}
		return endDict
	}

	switch p.State {
	case poleconst.StateNotSubmitted:
		return ""
	case poleconst.StateNoGo:
		return ""
	case poleconst.StateDictToDo:
		return ""
	case poleconst.StateDaToDo:
		return startDate()
	case poleconst.StateDaExpected:
		return startDate()
	case poleconst.StatePermissionPending:
		return startDate()
	case poleconst.StateToDo:
		return endDate()
	case poleconst.StateMarked:
		return endDate()
	case poleconst.StateNoAccess:
		return endDate()
	case poleconst.StateDenseNetwork:
		return endDate()
	case poleconst.StateHoleDone:
		return endDate()
	case poleconst.StateIncident:
		return endDate()
	case poleconst.StateDone:
		return ""
	case poleconst.StateAttachment:
		return ""
	case poleconst.StateCancelled:
		return ""
	case poleconst.StateDeleted:
		return ""
	}
	return "ERR"
}
