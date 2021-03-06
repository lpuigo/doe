package polesites

import (
	"fmt"
	"github.com/lpuig/ewin/doe/website/backend/model/bpu"
	"github.com/lpuig/ewin/doe/website/backend/model/items"
	"github.com/lpuig/ewin/doe/website/frontend/model/polesite/poleconst"
	"math"
	"sort"
	"strconv"
	"strings"
)

type Pole struct {
	Id             int
	Ref            string
	City           string
	Address        string
	Sticker        string
	Priority       int
	Lat            float64
	Long           float64
	State          string
	Date           string
	Actors         []string
	DtRef          string
	DictRef        string
	DictDate       string
	DictInfo       string
	DaQueryDate    string
	DaValidation   bool
	DaStartDate    string
	DaEndDate      string
	Height         int
	Material       string
	AspiDate       string
	Kizeo          string
	Comment        string
	Product        []string
	AttachmentDate string
	TimeStamp      string
}

func (p *Pole) SearchString() string {
	var searchBuilder strings.Builder
	fmt.Fprintf(&searchBuilder, "pole%s:%s,", "Ref", strings.ToUpper(p.ExtendedRef()))
	fmt.Fprintf(&searchBuilder, "pole%s:%s,", "City", strings.ToUpper(p.City))
	fmt.Fprintf(&searchBuilder, "pole%s:%s,", "Address", strings.ToUpper(p.Address))
	fmt.Fprintf(&searchBuilder, "pole%s:%s,", "DtRef", strings.ToUpper(p.DtRef))
	fmt.Fprintf(&searchBuilder, "pole%s:%s,", "DictRef", strings.ToUpper(p.DictRef))
	fmt.Fprintf(&searchBuilder, "pole%s:%s,", "Height", strconv.Itoa(p.Height)+"M")
	fmt.Fprintf(&searchBuilder, "pole%s:%s,", "Material", strings.ToUpper(p.Material))
	for _, key := range p.Product {
		fmt.Fprintf(&searchBuilder, "poleProduct:%s,", strings.ToUpper(key))
	}
	for _, actor := range p.Actors {
		fmt.Fprintf(&searchBuilder, "poleActor:%s,", strings.ToUpper(actor))
	}
	fmt.Fprintf(&searchBuilder, "pole%s:%s,", "DictInfo", strings.ToUpper(p.DictInfo))
	fmt.Fprintf(&searchBuilder, "pole%s:%s,", "Date", strings.ToUpper(p.Date))
	fmt.Fprintf(&searchBuilder, "pole%s:%s,", "AspiDate", strings.ToUpper(p.AspiDate))
	fmt.Fprintf(&searchBuilder, "pole%s:%s,", "AttachDate", strings.ToUpper(p.AttachmentDate))
	return searchBuilder.String()
}

func (p *Pole) IsTodo() bool {
	//switch p.State {
	////case poleconst.StateNotSubmitted:
	////case poleconst.StateNoGo:
	//case poleconst.StateDictToDo:
	//	return true
	//case poleconst.StateDaToDo:
	//	return true
	//case poleconst.StateDaExpected:
	//	return true
	//case poleconst.StatePermissionPending:
	//	return true
	//case poleconst.StateToDo:
	//	return true
	//case poleconst.StateMarked:
	//	return true
	//case poleconst.StateHoleDone:
	//	return true
	//case poleconst.StateIncident:
	//	return true
	//case poleconst.StateDone:
	//	return true
	//case poleconst.StateAttachment:
	//	return true
	////case poleconst.StateCancelled:
	//default:
	//	return false
	//}
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

func (p *Pole) IsDone() bool {
	switch p.State {
	//case poleconst.StateNotSubmitted:
	//case poleconst.StateNoGo:
	//case poleconst.StateDictToDo:
	//case poleconst.StateDaToDo:
	//case poleconst.StateDaExpected:
	//case poleconst.StatePermissionPending:
	//case poleconst.StateToDo:
	//case poleconst.StateMarked:
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

func (p *Pole) IsBlocked() bool {
	switch p.State {
	//case poleconst.StateNotSubmitted:
	case poleconst.StateNoGo:
		return true
	//case poleconst.StateDictToDo:
	//case poleconst.StateToDo:
	//case poleconst.StateHoleDone:
	case poleconst.StateIncident:
		return true
	//case poleconst.StateDone:
	//case poleconst.StateCancelled:
	//case poleconst.StateDeleted:
	default:
		return false
	}
}

func (p *Pole) IsBilled() bool {
	switch p.State {
	//case poleconst.StateNotSubmitted:
	//case poleconst.StateNoGo:
	//case poleconst.StateDictToDo:
	//case poleconst.StateDaToDo:
	//case poleconst.StateDaExpected:
	//case poleconst.StatePermissionPending:
	//case poleconst.StateToDo:
	//case poleconst.StateHoleDone:
	//case poleconst.StateIncident:
	//case poleconst.StateDone:
	case poleconst.StateAttachment:
		return true
	//case poleconst.StateCancelled:
	//case poleconst.StateDeleted:
	default:
		return false
	}
}

func (p *Pole) IsArchivable() bool {
	switch p.State {
	//case poleconst.StateNotSubmitted:
	//case poleconst.StateNoGo:
	//case poleconst.StateDictToDo:
	//case poleconst.StateDaToDo:
	//case poleconst.StateDaExpected:
	//case poleconst.StatePermissionPending:
	//case poleconst.StateToDo:
	//case poleconst.StateHoleDone:
	//case poleconst.StateIncident:
	//case poleconst.StateDone:
	case poleconst.StateAttachment:
		return true
	case poleconst.StateCancelled:
		return true
	case poleconst.StateDeleted:
		return true
	default:
		return false
	}
}

const (
	activityPole    string = "Poteaux"
	catPoleCreation string = "Création"
)

func (p *Pole) ExtendedRef() string {
	ref := strings.Trim(p.Ref, " ")
	if p.Sticker != "" {
		ref += " " + strings.Trim(p.Sticker, " ")
	}
	return ref
}

func (p *Pole) Itemize(client, site string, currentBpu *bpu.Bpu) ([]*items.Item, error) {
	res := []*items.Item{}

	poleArticles := currentBpu.GetCategoryArticles(activityPole)

	todo, done, blocked, billed := p.IsTodo(), p.IsDone(), p.IsBlocked(), p.IsBilled()

	//article, err := poleArticles.GetArticleFor(catPoleCreation, p.Height)
	//if err != nil {
	//	return nil, fmt.Errorf("can not define pole creation Item: %s", err.Error())
	//}
	//
	//info := fmt.Sprintf("Création poteau %s %dm", p.Material, p.Height)
	//
	//it := items.NewItem(
	//	client,
	//	site,
	//	activityPole,
	//	ref,
	//	info,
	//	p.Date,
	//	"",
	//	article,
	//	1,
	//	1,
	//	todo,
	//	done,
	//	blocked,
	//	billed,
	//)
	//it.Comment = p.Comment
	//it.Actors = p.Actors
	//if billed {
	//	it.AttachDate = p.AttachmentDate
	//}
	//res = append(res, it)
	ref := p.ExtendedRef()

	for _, product := range p.Product {
		article, err := poleArticles.GetArticleFor(product, p.Height)
		if err != nil {
			return nil, fmt.Errorf("can not define pole product item: %s", err.Error())
		}
		info := ""
		comment := ""
		if product == poleconst.ProductCreation {
			info = fmt.Sprintf("Création poteau %s %dm", p.Material, p.Height)
			comment = p.Comment
		} else {
			info = fmt.Sprintf("prestation complémentaire %s", product)
		}

		it := items.NewItem(
			client,
			site,
			activityPole,
			ref,
			info,
			p.Date,
			"",
			article,
			1,
			1,
			todo,
			done,
			blocked,
			billed,
		)
		it.Comment = comment
		it.Actors = p.Actors
		if billed {
			it.AttachDate = p.AttachmentDate
		}
		res = append(res, it)
	}

	return res, nil
}

const polePosPrecision float64 = 0.000000001

// IsEqual returns true if both Pole are identical (long and lat must be 1e-9 near, TimeStamp is not checked)
func (p *Pole) IsEqual(pole *Pole) bool {
	//Id             int
	if p.Id != pole.Id {
		return false
	}
	return p.IsEquivalent(pole)
}

// IsEquivalent returns true if both Pole are identical (long and lat must be 1e-9 near, Id and TimeStamp is not checked)
func (p *Pole) IsEquivalent(pole *Pole) bool {
	//Ref            string
	if p.Ref != pole.Ref {
		return false
	}
	//City           string
	if p.City != pole.City {
		return false
	}
	//Address        string
	if p.Address != pole.Address {
		return false
	}
	//Sticker        string
	if p.Sticker != pole.Sticker {
		return false
	}
	//Priority        int
	if p.Priority != pole.Priority {
		return false
	}
	//Lat            float64
	if math.Abs(p.Lat-pole.Lat) > polePosPrecision {
		return false
	}
	//Long           float64
	if math.Abs(p.Long-pole.Long) > polePosPrecision {
		return false
	}
	//State          string
	if p.State != pole.State {
		return false
	}
	//Date           string
	if p.Date != pole.Date {
		return false
	}
	//Actors         []string
	if len(p.Actors) != len(pole.Actors) {
		return false
	}
	for i, act := range pole.Actors {
		if p.Actors[i] != act {
			return false
		}
	}
	//DtRef          string
	if p.DtRef != pole.DtRef {
		return false
	}
	//DictRef        string
	if p.DictRef != pole.DictRef {
		return false
	}
	//DictDate       string
	if p.DictDate != pole.DictDate {
		return false
	}
	//DictInfo       string
	if p.DictInfo != pole.DictInfo {
		return false
	}
	//DaQueryDate    string
	if p.DaQueryDate != pole.DaQueryDate {
		return false
	}
	//DaValidation   bool
	if p.DaValidation != pole.DaValidation {
		return false
	}
	//DaStartDate    string
	if p.DaStartDate != pole.DaStartDate {
		return false
	}
	//DaEndDate      string
	if p.DaEndDate != pole.DaEndDate {
		return false
	}
	//Height         int
	if p.Height != pole.Height {
		return false
	}
	//Material       string
	if p.Material != pole.Material {
		return false
	}
	//AspiDate       string
	if p.AspiDate != pole.AspiDate {
		return false
	}
	//Kizeo          string
	if p.Kizeo != pole.Kizeo {
		return false
	}
	//Comment        string
	if p.Comment != pole.Comment {
		return false
	}
	//Product        []string
	if len(p.Product) != len(pole.Product) {
		return false
	}
	for i, prd := range pole.Product {
		if p.Product[i] != prd {
			return false
		}
	}
	//AttachmentDate string
	if p.AttachmentDate != pole.AttachmentDate {
		return false
	}
	//TimeStamp      string
	// not compared
	return true
}

func (p *Pole) UpdateProductFrom(npole *Pole) string {
	productDict := func(prds []string) map[string]bool {
		prdDict := make(map[string]bool)
		for _, prd := range prds {
			prdDict[prd] = true
		}
		return prdDict
	}

	isMainProduct := func(prd string) bool {
		switch prd {
		case
			poleconst.ProductCreation,
			poleconst.ProductInRow,
			poleconst.ProductReplace,
			poleconst.ProductRemove,
			poleconst.ProductStraighten,
			poleconst.ProductMoise,
			poleconst.ProductCouple,
			poleconst.ProductHauban,
			poleconst.ProductFarReplenishment,
			poleconst.ProductReplenishment:
			return true
		default:
			return false
		}
	}

	existingPrdDict := productDict(p.Product)
	newPrdDict := productDict(npole.Product)

	changes := ""

	// add product from npole
	for nprd, _ := range newPrdDict {
		if existingPrdDict[nprd] {
			continue
		}
		existingPrdDict[nprd] = true
		changes += " +" + nprd
	}

	// remove main product not in npole
	for prd, _ := range existingPrdDict {
		if !isMainProduct(prd) { // not main product must not be remove from existing list
			continue
		}
		if !newPrdDict[prd] {
			delete(existingPrdDict, prd)
			changes += " -" + prd
		}
	}

	nprdList := make([]string, len(existingPrdDict))
	i := 0
	for prd, _ := range existingPrdDict {
		nprdList[i] = prd
		i++
	}
	sort.Strings(nprdList)
	p.Product = nprdList
	return changes
}

func CheckCAPFTPoleInfo(info string) bool {
	if len(info) < 2 {
		return false
	}
	info = strings.ToUpper(info)
	switch info[0:2] {
	case "BS", "BH", "BM", "BC":
		return true
	case "MI", "MS", "MH", "MC", "XC", "MT", "M4", "MF", "MR":
		return true
	case "FR":
		return true
	default:
		return false
	}
}

func DecodeCAPFTPoleInfo(info string, products *[]string) (mat string, height int) {
	targetWords := strings.Split(strings.Trim(strings.ToUpper(info), " "), " ")
	if len(targetWords) > 0 && len(targetWords[0]) > 2 {
		switch targetWords[0][0:2] {
		case "BS":
			mat = poleconst.MaterialWood

		case "BH":
			mat = poleconst.MaterialWood
			*products = append(*products, poleconst.ProductHauban)

		case "BM":
			mat = poleconst.MaterialWood
			*products = append(*products, poleconst.ProductMoise)

		case "BC":
			mat = poleconst.MaterialWood
			*products = append(*products, poleconst.ProductCouple)

		case "MI", "MS":
			mat = poleconst.MaterialMetal

		case "MH":
			mat = poleconst.MaterialMetal
			*products = append(*products, poleconst.ProductHauban)

		case "MC", "XC":
			mat = poleconst.MaterialMetal
			*products = append(*products, poleconst.ProductCouple)

		case "M4", "MF", "MR":
			mat = poleconst.MaterialEnforcedMetal

		case "FR":
			mat = poleconst.MaterialEnforcedComp
		}
		if mat != "" {
			height = int(targetWords[0][2] - '0')
			if height <= 1 { // 0 or 1
				height = 10
			}
		} else {
			height = 6
		}
	}
	return
}
