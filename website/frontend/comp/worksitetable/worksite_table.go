package worksitetable

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/progressbar"
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
		worksiteinfo.RegisterComponentWorksiteInfo(),
		progressbar.RegisterComponent(),
		hvue.Template(template),
		hvue.Props("worksiteinfos", "enable_add_worksite"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewWorksiteTableModel(vm)
		}),
		hvue.MethodsOf(&WorksiteTableModel{}),
		hvue.Computed("filteredWorksites", func(vm *hvue.VM) interface{} {
			wtm := &WorksiteTableModel{Object: vm.Object}
			if wtm.Filter == "" {
				return wtm.Worksiteinfos
			}
			res := []*fm.WorksiteInfo{}
			for _, ws := range wtm.Worksiteinfos {
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

	Worksiteinfos     []*fm.WorksiteInfo `js:"worksiteinfos"`
	EnableAddWorksite bool               `js:"enable_add_worksite"`
	Filter            string             `js:"filter"`

	VM *hvue.VM `js:"VM"`
}

func NewWorksiteTableModel(vm *hvue.VM) *WorksiteTableModel {
	wtm := &WorksiteTableModel{Object: tools.O()}
	wtm.Worksiteinfos = nil
	wtm.EnableAddWorksite = false
	wtm.Filter = ""
	wtm.VM = vm
	return wtm
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Event Related Methods

func (wtm *WorksiteTableModel) SetSelectedWorksite(wsi *fm.WorksiteInfo) {
	wtm.VM.Emit("selected_worksite", wsi.Id)
}

//func (wtm *WorksiteTableModel) SaveWorksite(vm *hvue.VM, uws *fm.Worksite) {
//	vm.Emit("save_worksite", uws)
//}

func (wtm *WorksiteTableModel) ExpandRow(vm *hvue.VM, ws *fm.WorksiteInfo, others *js.Object) {
	print("Others :", others)
}

func (wtm *WorksiteTableModel) AddWorksite(vm *hvue.VM) {
	vm.Emit("new_worksite")
}

func (wtm *WorksiteTableModel) IsReworkable(status string) bool {
	return fm.WorksiteIsReworkable(status)
}

func (wtm *WorksiteTableModel) ReworkIconColor(wsi *fm.WorksiteInfo) string {
	if !wsi.Inspected {
		return ""
	}
	if wsi.NbRework > 0 {
		if wsi.NbReworkDone == wsi.NbRework {
			return "rework-orange"
		}
		return "rework-red"
	}
	return "rework-green"
}

func (wtm *WorksiteTableModel) CreateRework(vm *hvue.VM, wsi *fm.WorksiteInfo) {
	vm.Emit("edit_rework", wsi.Id)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Formatting Related Methods

func (wtm *WorksiteTableModel) TableRowClassName(rowInfo *js.Object) string {
	wsi := &fm.WorksiteInfo{Object: rowInfo.Get("row")}
	var res string = ""
	switch wsi.Status {
	case fm.WsStatusDone:
		res = "worksite-row-done"
	case fm.WsStatusRework:
		res = "worksite-row-rework"
	case fm.WsStatusPayment:
		res = "worksite-row-payment"
	case fm.WsStatusAttachment:
		res = "worksite-row-attachment"
	case fm.WsStatusDOE:
		res = "worksite-row-doe"
	case fm.WsStatusInProgress:
		res = "worksite-row-inprogress"
	case fm.WsStatusFormInProgress:
		res = "worksite-row-forminprogress"
	case fm.WsStatusNew:
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

	for _, ws := range wtm.Worksiteinfos {
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
		fm.WsStatusNew,
		fm.WsStatusFormInProgress,
		fm.WsStatusInProgress,
		fm.WsStatusDOE,
		fm.WsStatusAttachment,
		fm.WsStatusPayment,
		fm.WsStatusRework,
	}
	return res
}
