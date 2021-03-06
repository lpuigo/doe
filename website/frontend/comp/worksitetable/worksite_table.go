package worksitetable

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/progressbar"
	"github.com/lpuig/ewin/doe/website/frontend/comp/worksiteinfo"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/model/worksite"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/date"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
	"sort"
	"strconv"
	"time"
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
				return wtm.GetSizeLimitedResult(wtm.Worksiteinfos)
			}
			res := []*fm.WorksiteInfo{}
			for _, ws := range wtm.Worksiteinfos {
				if ws.TextFiltered(wtm.Filter) {
					res = append(res, ws)
				}
			}
			return wtm.GetSizeLimitedResult(res)
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
	SizeLimit         int                `js:"SizeLimit"`

	VM *hvue.VM `js:"VM"`
}

func NewWorksiteTableModel(vm *hvue.VM) *WorksiteTableModel {
	wtm := &WorksiteTableModel{Object: tools.O()}
	wtm.Worksiteinfos = nil
	wtm.EnableAddWorksite = false
	wtm.Filter = ""
	wtm.SetSizeLimit()
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
	return worksite.WorksiteIsReworkable(status)
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
	return worksite.WorksiteRowClassName(wsi.Status)
}

func (wtm *WorksiteTableModel) HeaderCellStyle() string {
	return "background: #a1e6e6;"
}

func (wtm *WorksiteTableModel) FormatDate(r, c *js.Object, d string) string {
	return date.DateString(d)
}

func (wtm *WorksiteTableModel) FormatStatus(r, c *js.Object, d string) string {
	return worksite.WorksiteStatusLabel(d)
}

func (wtm *WorksiteTableModel) SortStatus(a, b *worksite.Worksite) int {
	la := worksite.WorksiteStatusLabel(a.Status)
	lb := worksite.WorksiteStatusLabel(b.Status)
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

func (wtm *WorksiteTableModel) GetSizeLimitedResult(res []*fm.WorksiteInfo) []*fm.WorksiteInfo {
	if len(res) == wtm.SizeLimit {
		return res
	}
	if len(res) > sizeLimitDefault {
		wtm.ResetSizeLimit(len(res))
		return res[len(res)-sizeLimitDefault:]
	}
	return res
}

func (wtm *WorksiteTableModel) SetSizeLimit() {
	wtm.SizeLimit = -1
}

func (wtm *WorksiteTableModel) ResetSizeLimit(size int) {
	go func() {
		time.Sleep(sizeLimitTimer * time.Millisecond)
		wtm.SizeLimit = size
	}()
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
			return worksite.WorksiteStatusLabel(val)
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
		//fm.WsStatusNew,
		//fm.WsStatusFormInProgress,
		//fm.WsStatusInProgress,
		//fm.WsStatusDOE,
		//fm.WsStatusAttachment,
		//fm.WsStatusPayment,
		//fm.WsStatusRework,
	}
	return res
}
