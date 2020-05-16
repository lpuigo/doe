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
	// naive algorithm ... something smarter must be possible
	maxid := -1
	for _, pole := range ps.Poles {
		if pole.Id >= maxid {
			maxid = pole.Id + 1
		}
	}
	return maxid
}

// AddPole adds the given pole to polesite, and sets pole's new Id to ensure Id unicity
func (ps *Polesite) AddPole(pole *Pole) {
	pole.Id = ps.getNextId()
	ps.Poles = append(ps.Poles, pole)
}

// DeletePole deletes the given pole and returns true if it was found and deleted, false otherwise (no-op)
func (ps *Polesite) DeletePole(pole *Pole) bool {
	for i, p := range ps.Poles {
		if p.Id == pole.Id {
			// remove the item the JS way, to triggger vueJS observers
			ps.Object.Get("Poles").Call("splice", i, 1)
			return true
		}
	}
	return false
}

// DeletePole deletes the given pole and returns true if it was found and deleted, false otherwise (no-op)
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
		pole.CheckState()
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
