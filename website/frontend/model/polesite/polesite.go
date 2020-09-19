package polesite

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/model/polesite/poleconst"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
	"github.com/lpuig/ewin/doe/website/frontend/tools/json"
)

// Polesite reflects backend/model/polesites.polesite struct
type Polesite struct {
	*js.Object

	Id         int     `js:"Id"`
	Client     string  `js:"Client"`
	Ref        string  `js:"Ref"`
	Manager    string  `js:"Manager"`
	OrderDate  string  `js:"OrderDate"`
	UpdateDate string  `js:"UpdateDate"`
	Status     string  `js:"Status"`
	Comment    string  `js:"Comment"`
	Poles      []*Pole `js:"Poles"`
}

func (ps *Polesite) getNextId() int {
	// ps.Poles must be sorted by ascending Pole Ids
	// so last Pole has the highest id
	nbPoles := len(ps.Poles)
	if nbPoles == 0 {
		return 0
	}
	// returns highest + 1
	return ps.Poles[nbPoles-1].Id + 1

}

// AddPole adds the given pole to polesite, and sets pole's new Id to ensure Id unicity
func (ps *Polesite) AddPole(pole *Pole) {
	pole.Id = ps.getNextId()
	ps.Poles = append(ps.Poles, pole)
}

// DeletePole deletes the given pole and returns true if it was found and deleted, false otherwise (no-op)
func (ps *Polesite) DeletePole(pole *Pole) bool {
	// Set the state of given Pole to StateDeleted
	pole.State = poleconst.StateDeleted
	return true
	//for _, p := range ps.Poles {
	//	if p.Id == pole.Id {
	//		// remove the item the JS way, to trigger vueJS observers
	//		ps.Object.Get("Poles").Call("splice", i, 1)
	//		return true
	//	}
	//}
	//return false
}

// DuplicatePole duplicates the given pole and returns true if it was found and deleted, false otherwise (no-op)
func (ps *Polesite) DuplicatePole(pole *Pole) bool {
	print("DuplicatePole", pole.Object)
	return true
}

func NewPolesite() *Polesite {
	return &Polesite{Object: tools.O()}
}

func PolesiteFromJS(o *js.Object) *Polesite {
	return &Polesite{Object: o}
}

func (ps *Polesite) SearchInString() string {
	return json.Stringify(ps)
}

func (ps *Polesite) Copy(ops *Polesite) {
	ps.Id = ops.Id
	ps.Client = ops.Client
	ps.Manager = ops.Manager
	ps.OrderDate = ops.OrderDate
	ps.UpdateDate = ops.UpdateDate
	ps.Status = ops.Status
	ps.Comment = ops.Comment
	poles := make([]*Pole, len(ops.Poles))
	for ip, pole := range ops.Poles {
		ps.Poles[ip] = pole.Clone()
	}
	ps.Poles = poles
}

func (ps *Polesite) Clone() *Polesite {
	return &Polesite{Object: json.Parse(json.Stringify(ps))}
}

func (ps *Polesite) CheckPolesStatus() {
	for _, pole := range ps.Poles {
		pole.CheckInfoConsistency()
		pole.UpdateState()
	}
}

func PolesiteStatusLabel(status string) string {
	switch status {
	case poleconst.PsStatusNew:
		return poleconst.PsStatusLabelNew
	case poleconst.PsStatusInProgress:
		return poleconst.PsStatusLabelInProgress
	case poleconst.PsStatusBlocked:
		return poleconst.PsStatusLabelBlocked
	case poleconst.PsStatusCancelled:
		return poleconst.PsStatusLabelCancelled
	case poleconst.PsStatusDone:
		return poleconst.PsStatusLabelDone
	default:
		return "<" + status + ">"
	}
}

func PolesiteRowClassName(status string) string {
	var res string = ""
	switch status {
	case poleconst.PsStatusNew:
		return "worksite-row-new"
	case poleconst.PsStatusInProgress:
		return "worksite-row-inprogress"
	case poleconst.PsStatusBlocked:
		return "worksite-row-blocked"
	case poleconst.PsStatusCancelled:
		return "worksite-row-canceled"
	case poleconst.PsStatusDone:
		return "worksite-row-done"
	default:
		res = "worksite-row-error"
	}
	return res
}

func GetPoleSiteStatesValueLabel() []*elements.ValueLabel {
	return []*elements.ValueLabel{
		elements.NewValueLabel(poleconst.PsStatusNew, PolesiteStatusLabel(poleconst.PsStatusNew)),
		elements.NewValueLabel(poleconst.PsStatusInProgress, PolesiteStatusLabel(poleconst.PsStatusInProgress)),
		elements.NewValueLabel(poleconst.PsStatusBlocked, PolesiteStatusLabel(poleconst.PsStatusBlocked)),
		elements.NewValueLabel(poleconst.PsStatusCancelled, PolesiteStatusLabel(poleconst.PsStatusCancelled)),
		elements.NewValueLabel(poleconst.PsStatusDone, PolesiteStatusLabel(poleconst.PsStatusDone)),
	}
}

type pos struct {
	lat, long int
}

func (ps *Polesite) DetectDuplicate() {
	const prec float64 = 100000.0
	duplicatePolesByPos := map[pos][]*Pole{}
	duplicatePolesByTitle := map[string][]*Pole{}
	foundDuplicateByPos := false
	foundDuplicateByTitle := false
	for _, pole := range ps.Poles {
		polePos := pos{lat: int(pole.Lat*prec) / 2, long: int(pole.Long*prec) / 2}
		polesWithPos, found := duplicatePolesByPos[polePos]
		if !found {
			duplicatePolesByPos[polePos] = []*Pole{pole}
		} else {
			duplicatePolesByPos[polePos] = append(polesWithPos, pole)
			foundDuplicateByPos = true
		}

		title := pole.GetTitle()
		polesWithTitle, found := duplicatePolesByTitle[title]
		if !found {
			duplicatePolesByTitle[title] = []*Pole{pole}
		} else {
			duplicatePolesByTitle[title] = append(polesWithTitle, pole)
			foundDuplicateByTitle = true
		}
	}

	if foundDuplicateByPos {
		print("====================== Duplicate by Position ====================")
		for _, poles := range duplicatePolesByPos {
			if len(poles) > 1 {
				print("Duplicate from", poles[0].GetTitle())
				for _, pole := range poles[1:] {
					print("==> ", pole.GetTitle())
				}
			}
		}
	}

	if foundDuplicateByTitle {
		print("====================== Duplicate by Title =======================")
		for _, poles := range duplicatePolesByTitle {
			if len(poles) > 1 {
				print("Duplicate from", poles[0].GetTitle())
				for _, pole := range poles[1:] {
					print("==> ", pole.GetTitle())
				}
			}
		}
	}
}

func (ps *Polesite) DetectProductInconsistency() {
	for _, pole := range ps.Poles {
		if !(pole.State != poleconst.StateCancelled && pole.State != poleconst.StateNotSubmitted) {
			continue
		}
		pole.AddProduct(poleconst.ProductCreation)
	}
}

func (ps *Polesite) DetectMissingDAValidation() bool {
	updateMap := false
	for _, pole := range ps.Poles {
		if !(!tools.Empty(pole.DaQueryDate) && !tools.Empty(pole.DaStartDate) && !tools.Empty(pole.DaEndDate) && !pole.DaValidation) {
			continue
		}
		pole.DaValidation = true
		currentState := pole.State
		pole.UpdateState()
		if pole.State != currentState {
			updateMap = true
		}
	}
	return updateMap
}
