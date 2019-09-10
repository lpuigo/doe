package poletable

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/ripprogressbar"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	ps "github.com/lpuig/ewin/doe/website/frontend/model/polesite"
	"github.com/lpuig/ewin/doe/website/frontend/model/polesite/poleconst"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/date"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements/message"
	"sort"
	"strconv"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Registration

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("pole-table", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		ripprogressbar.RegisterComponent(),
		hvue.Template(template),
		hvue.Props("user", "polesite", "filter", "filtertype"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewPoleTableModel(vm)
		}),
		hvue.MethodsOf(&PoleTableModel{}),
		hvue.Computed("filteredPoles", func(vm *hvue.VM) interface{} {
			rtm := &PoleTableModel{Object: vm.Object}
			return rtm.GetFilteredPole()
		}),
		hvue.Filter("DateFormat", func(vm *hvue.VM, value *js.Object, args ...*js.Object) interface{} {
			return date.DateString(value.String())
		}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type PoleTableModel struct {
	*js.Object

	Polesite   *ps.Polesite `js:"polesite"`
	User       *fm.User     `js:"user"`
	Filter     string       `js:"filter"`
	FilterType string       `js:"filtertype"`

	VM *hvue.VM `js:"VM"`
}

func NewPoleTableModel(vm *hvue.VM) *PoleTableModel {
	rtm := &PoleTableModel{Object: tools.O()}
	rtm.Polesite = nil
	rtm.User = fm.NewUser()
	rtm.Filter = ""
	rtm.FilterType = ""
	rtm.VM = vm
	return rtm
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Event Related Methods

func (ptm *PoleTableModel) SetSelectedPole(vm *hvue.VM, p *ps.Pole) {
	vm.Emit("pole-selected", p)
}

func (ptm *PoleTableModel) AddPole(vm *hvue.VM) {
	message.InfoStr(vm, "AddPole non encore implémenté", false)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Formatting Related Methods

func (ptm *PoleTableModel) TableRowClassName(rowInfo *js.Object) string {
	p := &ps.Pole{Object: rowInfo.Get("row")}
	return ps.PoleRowClassName(p.State)
}

func (ptm *PoleTableModel) HeaderCellStyle() string {
	return "background: #a1e6e6;"
}

func (ptm *PoleTableModel) FormatDate(r, c *js.Object, d string) string {
	return date.DateString(d)
}

func (ptm *PoleTableModel) FormatState(r, c *js.Object, d string) string {
	return ps.PoleStateLabel(d)
}

func (ptm *PoleTableModel) FormatProduct(p *ps.Pole) string {
	return strings.Join(p.Product, "\n")
}

func (ptm *PoleTableModel) FormatActors(vm *hvue.VM, p *ps.Pole) string {
	ptm = &PoleTableModel{Object: vm.Object}
	client := ptm.User.GetClientByName(ptm.Polesite.Client)
	actors := []string{}
	for _, actId := range p.Actors {
		actor := client.GetActorBy(actId)
		if actor != nil {
			actors = append(actors, actor.LastName)
		}
	}
	return strings.Join(actors, "\n")
}

func (ptm *PoleTableModel) SortState(a, b *ps.Pole) int {
	la := ps.PoleStateLabel(a.State)
	lb := ps.PoleStateLabel(b.State)
	if la < lb {
		return -1
	}
	if la == lb {
		return 0
	}
	return 1
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Row Filtering Related Methods

func (ptm *PoleTableModel) GetFilteredPole() []*ps.Pole {
	if ptm.FilterType == poleconst.FilterValueAll && ptm.Filter == "" {
		return ptm.Polesite.Poles
	}

	res := []*ps.Pole{}
	expected := strings.ToUpper(ptm.Filter)
	filter := func(p *ps.Pole) bool {
		sis := p.SearchString(ptm.FilterType)
		if sis == "" {
			return false
		}
		return strings.Contains(strings.ToUpper(sis), expected)
	}
	for _, pole := range ptm.Polesite.Poles {
		if filter(pole) {
			res = append(res, pole)
		}
	}
	return res
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Column Filtering Related Methods

func (ptm *PoleTableModel) FilterHandler(value string, p *js.Object, col *js.Object) bool {
	prop := col.Get("property").String()
	return p.Get(prop).String() == value
}

func (ptm *PoleTableModel) FilterList(vm *hvue.VM, prop string) []*elements.ValText {
	ptm = &PoleTableModel{Object: vm.Object}
	count := map[string]int{}
	attribs := []string{}

	var translate func(string) string
	switch prop {
	case "Status":
		translate = func(val string) string {
			return ps.PolesiteStatusLabel(val)
		}
	default:
		translate = func(val string) string { return val }
	}

	for _, psi := range ptm.Polesite.Poles {
		attrib := psi.Object.Get(prop).String()
		if _, exist := count[attrib]; !exist {
			attribs = append(attribs, attrib)
		}
		count[attrib]++
	}
	sort.Strings(attribs)
	res := []*elements.ValText{}
	for _, a := range attribs {
		fa := a
		if fa == "" {
			fa = "Vide"
		}
		res = append(res, elements.NewValText(a, translate(fa)+" ("+strconv.Itoa(count[a])+")"))
	}
	return res
}

func (ptm *PoleTableModel) FilteredStatusValue() []string {
	res := []string{
		//poleconst.PsStatusNew,
		//poleconst.PsStatusInProgress,
		//poleconst.PsStatusBlocked,
		//poleconst.PsStatusCancelled,
		//poleconst.PsStatusDone,
	}
	return res
}
