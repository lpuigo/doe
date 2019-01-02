package worksitetable

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/worksiteedit"
	"github.com/lpuig/ewin/doe/website/frontend/comp/worksiteinfo"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/dates"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
	"sort"
	"strconv"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Registration

func Register() {
	hvue.NewComponent("worksite-table",
		ComponentOptions()...,
	)
}

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("worksite-table", ComponentOptions()...)
}

func ComponentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		worksiteinfo.RegisterComponent(),
		worksiteedit.RegisterComponent(),
		hvue.Template(template),
		hvue.Props("worksites"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewWorksiteTableModel(vm)
		}),
		hvue.MethodsOf(&WorksiteTableModel{}),
		hvue.Computed("filteredWorksites", func(vm *hvue.VM) interface{} {
			wtm := &WorksiteTableModel{Object: vm.Object}
			if wtm.Filter == "" {
				return wtm.Worksites
			}
			res := []*fm.Worksite{}
			for _, ws := range wtm.Worksites {
				if ws.TextFiltered(wtm.Filter) {
					res = append(res, ws)
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

type WorksiteTableModel struct {
	*js.Object

	Worksites []*fm.Worksite `js:"worksites"`
	Filter    string         `js:"filter"`

	VM *hvue.VM `js:"VM"`
}

func NewWorksiteTableModel(vm *hvue.VM) *WorksiteTableModel {
	wtm := &WorksiteTableModel{Object: tools.O()}
	wtm.Worksites = nil
	wtm.Filter = ""
	wtm.VM = vm
	return wtm
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Event Related Methods

//func (wtm *WorksiteTableModel) SelectRow(vm *hvue.VM, ws *fm.Worksite, event *js.Object) {
//	vm.Emit("selected_worksite", ws)
//}
//
func (wtm *WorksiteTableModel) SetSelectedWorksite(nws *fm.Worksite) {
	if nws.Object == nil { // this can happen when Worksites props gets updated
		return
	}
	wtm.VM.Emit("selected_worksite", nws)
}

func (wtm *WorksiteTableModel) SaveWorksite(vm *hvue.VM, uws *fm.Worksite) {
	vm.Emit("save_worksite", uws)
}

func (wtm *WorksiteTableModel) ExpandRow(vm *hvue.VM, ws *fm.Worksite, others *js.Object) {
	if ws.Dirty {
		print("Worksite is Dirty")
	}
	print("Others :", others)
}

func (wtm *WorksiteTableModel) AddWorksite(vm *hvue.VM) {
	vm.Emit("new_worksite")
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Formatting Related Methods

func (wtm *WorksiteTableModel) TableRowClassName(rowInfo *js.Object) string {
	ws := &fm.Worksite{Object: rowInfo.Get("row")}
	var res string = ""
	if ws.Dirty {
		return "worksite-row-need-save"
	}
	switch ws.Status {
	case "Done":
		res = "worksite-row-done"
	case "DOE":
		res = "worksite-row-doe"
	case "Rework":
		res = "worksite-row-rework"
	case "New":
		res = "worksite-row-new"
	default:
		res = ""
	}
	return res
}

func (wtm *WorksiteTableModel) HeaderCellStyle() string {
	return "background: #a1e6e6;"
}

func (wtm *WorksiteTableModel) FormatDate(r, c *js.Object, d string) string {
	return date.DateString(d)
}

func (wtm *WorksiteTableModel) FormatStatus(r, c *js.Object, d string) string {
	return fm.WorksiteStatusLabel(d)
}

func (wtm *WorksiteTableModel) SortStatus(a, b *fm.Worksite) int {
	la := fm.WorksiteStatusLabel(a.Status)
	lb := fm.WorksiteStatusLabel(b.Status)
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

func (wtm *WorksiteTableModel) FilterHandler(value string, p *js.Object, col *js.Object) bool {
	prop := col.Get("property").String()
	return p.Get(prop).String() == value
}

func (wtm *WorksiteTableModel) FilterList(vm *hvue.VM, prop string) []*elements.ValText {
	wtm = &WorksiteTableModel{Object: vm.Object}
	count := map[string]int{}
	attribs := []string{}

	var translate func(string) string
	switch prop {
	case "Status":
		translate = func(val string) string {
			return fm.WorksiteStatusLabel(val)
		}
	default:
		translate = func(val string) string { return val }
	}

	for _, ws := range wtm.Worksites {
		attrib := ws.Object.Get(prop).String()
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

func (wtm *WorksiteTableModel) FilteredStatusValue() []string {
	res := []string{
		"DOE",
		"InProgress",
		"New",
		"Rework",
		//fm.WorksiteStatusLabel("DOE"),
		//fm.WorksiteStatusLabel("InProgress"),
		//fm.WorksiteStatusLabel("New"),
	}
	return res
}
