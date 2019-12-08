package foasite

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/model/foasite/foaconst"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
	"github.com/lpuig/ewin/doe/website/frontend/tools/json"
)

// type Foa reflects backend/model/foasites.foa struct
type Foa struct {
	*js.Object
	Id    int    `js:"Id"`
	Ref   string `js:"Ref"`
	Insee string `js:"Insee"`
	Type  string `js:"Type"`

	State *State `js:"State"`
}

func NewFoa() *Foa {
	nf := &Foa{Object: tools.O()}
	nf.Id = -1
	nf.Ref = ""
	nf.Insee = ""
	nf.Type = ""
	nf.State = NewState()
	return nf
}

func FoaFromJs(o *js.Object) *Foa {
	return &Foa{Object: o}

}

func (f *Foa) SearchString(filter string) string {
	searchItem := func(prefix, typ, value string) string {
		if value == "" {
			return ""
		}
		if filter != foaconst.FilterValueAll && filter != typ {
			return ""
		}
		return prefix + typ + value
	}
	res := searchItem("", foaconst.FilterValueRef, f.Ref)
	res += searchItem(",", foaconst.FilterValueInsee, f.Insee)
	res += searchItem(",", foaconst.FilterValueComment, f.State.Comment)
	res += searchItem(",", foaconst.FilterValueType, f.Type)
	return res
}

func (f *Foa) Clone() *Foa {
	return FoaFromJs(json.Clone(f.Object))
}

func (f *Foa) Copy(of *Foa) {
	f.Insee = of.Insee
	f.Ref = of.Ref
	f.Type = of.Type
	f.State.Copy(of.State)
}

func (f *Foa) TotalBlockedDone() (total, blocked, done int) {
	switch f.State.Status {
	case foaconst.StateToDo:
		total = 1
	case foaconst.StateIncident:
		total = 1
		blocked = 1
	case foaconst.StateDone:
		total = 1
		done = 1
	case foaconst.StateAttachment:
		total = 1
		done = 1
	case foaconst.StateCancelled:
	default:
	}
	return
}

func (f *Foa) GetRowStyle() string {
	return FoaRowClassName(f.State.Status)
}

func FoaStateLabel(state string) string {
	switch state {
	case foaconst.StateToDo:
		return foaconst.LabelToDo
	case foaconst.StateIncident:
		return foaconst.LabelIncident
	case foaconst.StateDone:
		return foaconst.LabelDone
	case foaconst.StateAttachment:
		return foaconst.LabelAttachment
	case foaconst.StateCancelled:
		return foaconst.LabelCancelled
	default:
		return "<" + state + ">"
	}
}

func GetStatesValueLabel() []*elements.ValueLabel {
	return []*elements.ValueLabel{
		elements.NewValueLabel(foaconst.StateToDo, foaconst.LabelToDo),
		elements.NewValueLabel(foaconst.StateIncident, foaconst.LabelIncident),
		elements.NewValueLabel(foaconst.StateDone, foaconst.LabelDone),
		elements.NewValueLabel(foaconst.StateAttachment, foaconst.LabelAttachment),
		elements.NewValueLabel(foaconst.StateCancelled, foaconst.LabelCancelled),
	}
}

func GetFilterType() []*elements.ValueLabel {
	return []*elements.ValueLabel{
		elements.NewValueLabel(foaconst.FilterValueAll, foaconst.FilterLabelAll),
		elements.NewValueLabel(foaconst.FilterValueRef, foaconst.FilterLabelRef),
		elements.NewValueLabel(foaconst.FilterValueInsee, foaconst.FilterLabelInsee),
		elements.NewValueLabel(foaconst.FilterValueType, foaconst.FilterLabelType),
		elements.NewValueLabel(foaconst.FilterValueComment, foaconst.FilterLabelComment),
	}
}

func FoaRowClassName(status string) string {
	switch status {
	case foaconst.StateToDo:
		return "foa-row-todo"
	case foaconst.StateIncident:
		return "foa-row-incident"
	case foaconst.StateDone:
		return "foa-row-done"
	case foaconst.StateAttachment:
		return "foa-row-attachment"
	case foaconst.StateCancelled:
		return "foa-row-cancelled"
	default:
		return "foa-row-error"
	}
}
