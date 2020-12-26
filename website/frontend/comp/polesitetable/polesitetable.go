package polesitetable

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/ripprogressbar"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	ps "github.com/lpuig/ewin/doe/website/frontend/model/polesite"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/date"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
	"sort"
	"strconv"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Registration

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("polesite-table", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		ripprogressbar.RegisterComponent(),
		hvue.Template(template),
		hvue.Props("polesiteinfos", "user"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewPolesiteTableModel(vm)
		}),
		hvue.MethodsOf(&PolesiteTableModel{}),
		hvue.Computed("filteredPolesites", func(vm *hvue.VM) interface{} {
			rtm := &PolesiteTableModel{Object: vm.Object}
			if rtm.Filter == "" {
				return rtm.Polesiteinfos
			}
			res := []*fm.PolesiteInfo{}
			for _, psi := range rtm.Polesiteinfos {
				if psi.TextFiltered(rtm.Filter) {
					res = append(res, psi)
				}
			}
			return res
		}),
		hvue.Filter("DateFormat", func(vm *hvue.VM, value *js.Object, args ...*js.Object) interface{} {
			return date.DateString(value.String())
		}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type PolesiteTableModel struct {
	*js.Object

	Polesiteinfos []*fm.PolesiteInfo `js:"polesiteinfos"`
	User          *fm.User           `js:"user"`
	Filter        string             `js:"filter"`

	VM *hvue.VM `js:"VM"`
}

func NewPolesiteTableModel(vm *hvue.VM) *PolesiteTableModel {
	rtm := &PolesiteTableModel{Object: tools.O()}
	rtm.Polesiteinfos = nil
	rtm.User = fm.NewUser()
	//rtm.EnableAddWorksite = false
	rtm.Filter = ""
	rtm.VM = vm
	return rtm
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Event Related Methods

func (ptm *PolesiteTableModel) SetSelectedPolesite(psi *fm.PolesiteInfo) {
	ptm.OpenPolesite(psi.Id)
}

func (ptm *PolesiteTableModel) AddPolesite(vm *hvue.VM) {
	vm.Emit("new_polesite")
}

func (ptm *PolesiteTableModel) AttachmentUrl(id int) string {
	return "/api/polesites/" + strconv.Itoa(id) + "/attach"
}

func (ptm *PolesiteTableModel) ExportUrl(id int) string {
	return "/api/polesites/" + strconv.Itoa(id) + "/export"
}

func (ptm *PolesiteTableModel) OpenPolesite(id int) {
	js.Global.Get("window").Call("open", "polesite.html?v=1.0&psid="+strconv.Itoa(id))
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Formatting Related Methods

func (ptm *PolesiteTableModel) TableRowClassName(rowInfo *js.Object) string {
	psi := &fm.PolesiteInfo{Object: rowInfo.Get("row")}
	return ps.PolesiteRowClassName(psi.Status)
}

func (ptm *PolesiteTableModel) HeaderCellStyle() string {
	return "background: #a1e6e6;"
}

func (ptm *PolesiteTableModel) FormatDate(r, c *js.Object, d string) string {
	return date.DateString(d)
}

func (ptm *PolesiteTableModel) FormatStatus(r, c *js.Object, d string) string {
	return ps.PolesiteStatusLabel(d)
}

func (ptm *PolesiteTableModel) SortStatus(a, b *fm.RipsiteInfo) int {
	la := ps.PolesiteStatusLabel(a.Status)
	lb := ps.PolesiteStatusLabel(b.Status)
	if la < lb {
		return -1
	}
	if la == lb {
		return 0
	}
	return 1
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Column Filtering Related Methods

func (ptm *PolesiteTableModel) FilterHandler(value string, p *js.Object, col *js.Object) bool {
	prop := col.Get("property").String()
	return p.Get(prop).String() == value
}

func (ptm *PolesiteTableModel) FilterList(vm *hvue.VM, prop string) []*elements.ValText {
	ptm = &PolesiteTableModel{Object: vm.Object}
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

	for _, psi := range ptm.Polesiteinfos {
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

func (ptm *PolesiteTableModel) FilteredStatusValue() []string {
	res := []string{
		//poleconst.PsStatusNew,
		//poleconst.PsStatusInProgress,
		//poleconst.PsStatusBlocked,
		//poleconst.PsStatusCancelled,
		//poleconst.PsStatusDone,
	}
	return res
}
