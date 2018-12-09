package worksitetable

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	fm "github.com/lpuig/ewin/doe/website/frontend/frontmodel"
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

func ComponentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		//hvue.Component("project-progress-bar", wl_progress_bar.ComponentOptions()...),
		hvue.Props("worksites", "filter"),
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

	SelectedWorksite *fm.Worksite   `js:"selected_worksite"`
	Worksites        []*fm.Worksite `js:"worksites"`
	Filter           string         `js:"filter"`

	VM *hvue.VM `js:"VM"`
}

func NewWorksiteTableModel(vm *hvue.VM) *WorksiteTableModel {
	wtm := &WorksiteTableModel{Object: tools.O()}
	wtm.Worksites = nil
	wtm.SelectedWorksite = nil
	wtm.Filter = ""
	wtm.VM = vm
	return wtm
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Event Related Methods

func (wtm *WorksiteTableModel) SelectRow(vm *hvue.VM, prj *fm.Worksite, event *js.Object) {
	vm.Emit("selected_worksite", prj)
}

func (wtm *WorksiteTableModel) SetSelectedWorksite(nw *fm.Worksite) {
	if nw.Object == nil { // this can happen when Worksites props gets updated
		return
	}
	wtm.SelectedWorksite = nw
	wtm.VM.Emit("update:selected_worksite", nw)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Formatting Related Methods

func (wtm *WorksiteTableModel) TableRowClassName(rowInfo *js.Object) string {
	//ws := &fm.Worksite{Object: rowInfo.Get("row")}
	var res string = ""
	// TODO Add Worksite Status
	//switch ws.Status {
	//case "6 - Done", "0 - Lost":
	//	res = "worksite-row-done"
	//case "5 - Monitoring":
	//	res = "worksite-row-monitoring"
	//case "1 - Candidate", "2 - Outlining":
	//	res = "worksite-row-outline"
	//default:
	//	res = ""
	//}
	return res
}

func (wtm *WorksiteTableModel) HeaderCellStyle() string {
	return "background: #a1e6e6;"
}

func (wtm *WorksiteTableModel) FormatDate(r, c *js.Object, d string) string {
	return date.DateString(d)
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
			fa = "<Vide>"
		}
		res = append(res, elements.NewValText(a, fa+" ("+strconv.Itoa(count[a])+")"))
	}
	return res
}
