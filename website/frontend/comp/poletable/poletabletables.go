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
	"sort"
	"strconv"
	"strings"
	"time"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Registration

func registerComponentTable(tableType string) hvue.ComponentOption {
	var tableTemplate hvue.ComponentOption
	switch tableType {
	case "creation":
		tableTemplate = hvue.Template(template_creation)
	case "followup":
		tableTemplate = hvue.Template(template_followup)
	case "billing":
		tableTemplate = hvue.Template(template_billing)
	default:
		tableTemplate = hvue.Template("<span>Mode '" + tableType + "' non défini</span>")
	}
	return hvue.Component("pole-table-"+tableType, componentOptionsTable(tableTemplate)...)
}

func componentOptionsTable(tableTemplate hvue.ComponentOption) []hvue.ComponentOption {
	return []hvue.ComponentOption{
		ripprogressbar.RegisterComponent(),
		tableTemplate,
		hvue.Props("user", "polesite", "filter", "filtertype", "context"),
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
	Context    *Context     `js:"context"`
	SizeLimit  int          `js:"SizeLimit"`

	VM *hvue.VM `js:"VM"`
}

func NewPoleTableModel(vm *hvue.VM) *PoleTableModel {
	rtm := &PoleTableModel{Object: tools.O()}
	rtm.Polesite = nil
	rtm.User = fm.NewUser()
	rtm.Filter = ""
	rtm.FilterType = ""
	rtm.Context = NewContext("")
	rtm.SetSizeLimit()
	rtm.VM = vm
	return rtm
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Event Related Methods

func (ptm *PoleTableModel) SetSelectedPole(vm *hvue.VM, p *ps.Pole) {
	if p.Object == nil {
		// implicit callback ... skip
		return
	}
	ptm = &PoleTableModel{Object: vm.Object}
	//if p.Id == ptm.Context.SelectedPole {
	//	// no change ... skip
	//	return
	//}
	ptm.Context.SelectedPole = p.Id
	vm.Emit("update:context", ptm.Context)
}

//func (ptm *PoleTableModel) AddPole(vm *hvue.VM) {
//	message.InfoStr(vm, "AddPole non encore implémenté", false)
//}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Formatting Related Methods

func (ptm *PoleTableModel) TableRowClassName(vm *hvue.VM, rowInfo *js.Object) string {
	ptm = &PoleTableModel{Object: vm.Object}
	p := &ps.Pole{Object: rowInfo.Get("row")}
	selected := ""
	if ptm.Context.SelectedPole == p.Id {
		selected = "pole-selected "
	}
	return selected + ps.PoleRowClassName(p.State)
}

func (ptm *PoleTableModel) HeaderCellStyle() string {
	return "background: #a1e6e6;"
}

func (ptm *PoleTableModel) PoleRefName(p *ps.Pole) string {
	res := p.Ref
	if p.Sticker != "" {
		res += " " + p.Sticker
	}
	return res
}

func (ptm *PoleTableModel) DictEndDate(d string) string {
	if tools.Empty(d) {
		return "-"
	}
	return date.DateString(date.After(d, poleconst.DictValidityDuration))
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

func (ptm *PoleTableModel) FormatType(p *ps.Pole) string {
	return p.Material + " " + strconv.Itoa(p.Height) + "m"
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

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Sorting Methods

func (ptm *PoleTableModel) SortDate(attrib string) func(obj *js.Object) string {
	return func(obj *js.Object) string {
		val := obj.Get(attrib).String()
		if val == "" {
			return "9999-12-31"
		}
		return val
	}
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
// Size Related Methods

const (
	sizeLimitDefault int = 30
	sizeLimitTimer       = 300
)

func (ptm *PoleTableModel) GetSizeLimitedResult(res []*ps.Pole) []*ps.Pole {
	if len(res) == ptm.SizeLimit {
		return res
	}
	if len(res) > sizeLimitDefault {
		ptm.ResetSizeLimit(len(res))
		return res[len(res)-sizeLimitDefault:]
	}
	return res
}

func (ptm *PoleTableModel) SetSizeLimit() {
	ptm.SizeLimit = -1
}

func (ptm *PoleTableModel) ResetSizeLimit(size int) {
	go func() {
		time.Sleep(sizeLimitTimer * time.Millisecond)
		ptm.SizeLimit = size
	}()
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Row Filtering Related Methods

func (ptm *PoleTableModel) GetFilteredPole() []*ps.Pole {
	if ptm.FilterType == poleconst.FilterValueAll && ptm.Filter == "" {
		return ptm.GetSizeLimitedResult(ptm.Polesite.Poles)
	}

	res := []*ps.Pole{}
	expected := strings.ToUpper(strings.Trim(ptm.Filter, " "))
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
	return ptm.GetSizeLimitedResult(res)
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
	case "State":
		translate = func(val string) string {
			return ps.PoleStateLabel(val)
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
