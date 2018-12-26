package worksitetable

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/worksiteinfo"
	"github.com/lpuig/ewin/doe/website/frontend/comp/worksitetable/worksitedetail"
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

func ComponentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.Component("worksite-info", worksiteinfo.ComponentOptions()...),
		hvue.Component("worksite-detail", worksitedetail.ComponentOptions()...),
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

//func (wtm *WorksiteTableModel) SelectRow(vm *hvue.VM, ws *fm.Worksite, event *js.Object) {
//	vm.Emit("selected_worksite", ws)
//}
//
//func (wtm *WorksiteTableModel) SetSelectedWorksite(nws *fm.Worksite) {
//	if nws.Object == nil { // this can happen when Worksites props gets updated
//		return
//	}
//	wtm.SelectedWorksite = nws
//	wtm.VM.Emit("update:selected_worksite", nws)
//}

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
	wtm = &WorksiteTableModel{Object: vm.Object}
	print("Add New Worksite number", len(wtm.Worksites))
	ws := fm.NewWorkSite()
	ws.AddOrder()
	wtm.Worksites = append(wtm.Worksites, ws)
	print("Table now has", len(wtm.Worksites), "sites")
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
	//case "InProgress":
	//	res = "worksite-row-inprogress"
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

func (wtm *WorksiteTableModel) FilteredStatusValue() []string {
	res := []string{
		"InProgress",
		"New",
	}
	return res
}
