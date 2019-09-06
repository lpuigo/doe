package ripsitetable

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/ripprogressbar"
	"github.com/lpuig/ewin/doe/website/frontend/comp/ripsiteinfo"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	fmrip "github.com/lpuig/ewin/doe/website/frontend/model/ripsite"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/date"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
	"sort"
	"strconv"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Registration

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("ripsite-table", ComponentOptions()...)
}

func ComponentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		ripsiteinfo.RegisterComponentRipsiteInfoInfo(),
		ripprogressbar.RegisterComponent(),
		hvue.Template(template),
		hvue.Props("ripsiteinfos", "user"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewRipsiteTableModel(vm)
		}),
		hvue.MethodsOf(&RipsiteTableModel{}),
		hvue.Computed("filteredRipsites", func(vm *hvue.VM) interface{} {
			rtm := &RipsiteTableModel{Object: vm.Object}
			if rtm.Filter == "" {
				return rtm.Ripsiteinfos
			}
			res := []*fm.RipsiteInfo{}
			for _, rsi := range rtm.Ripsiteinfos {
				if rsi.TextFiltered(rtm.Filter) {
					res = append(res, rsi)
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

type RipsiteTableModel struct {
	*js.Object

	Ripsiteinfos []*fm.RipsiteInfo `js:"ripsiteinfos"`
	User         *fm.User          `js:"user"`
	//EnableAddWorksite bool               `js:"enable_add_worksite"`
	Filter string `js:"filter"`

	VM *hvue.VM `js:"VM"`
}

func NewRipsiteTableModel(vm *hvue.VM) *RipsiteTableModel {
	rtm := &RipsiteTableModel{Object: tools.O()}
	rtm.Ripsiteinfos = nil
	rtm.User = fm.NewUser()
	//rtm.EnableAddWorksite = false
	rtm.Filter = ""
	rtm.VM = vm
	return rtm
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Event Related Methods

func (rtm *RipsiteTableModel) SetSelectedRipsite(rsi *fm.RipsiteInfo) {
	rtm.VM.Emit("selected_ripsite", rsi.Id)
}

func (rtm *RipsiteTableModel) AddRipsite(vm *hvue.VM) {
	vm.Emit("new_ripsite")
}

func (rtm *RipsiteTableModel) AttachmentUrl(id int) string {
	return "/api/ripsites/" + strconv.Itoa(id) + "/attach"
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Formatting Related Methods

func (rtm *RipsiteTableModel) TableRowClassName(rowInfo *js.Object) string {
	wsi := &fm.RipsiteInfo{Object: rowInfo.Get("row")}
	return fmrip.RipsiteRowClassName(wsi.Status)
}

func (rtm *RipsiteTableModel) HeaderCellStyle() string {
	return "background: #a1e6e6;"
}

func (rtm *RipsiteTableModel) FormatDate(r, c *js.Object, d string) string {
	return date.DateString(d)
}

func (rtm *RipsiteTableModel) FormatStatus(r, c *js.Object, d string) string {
	return fmrip.RipsiteStatusLabel(d)
}

func (rtm *RipsiteTableModel) SortStatus(a, b *fm.RipsiteInfo) int {
	la := fmrip.RipsiteStatusLabel(a.Status)
	lb := fmrip.RipsiteStatusLabel(b.Status)
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

func (rtm *RipsiteTableModel) FilterHandler(value string, p *js.Object, col *js.Object) bool {
	prop := col.Get("property").String()
	return p.Get(prop).String() == value
}

func (rtm *RipsiteTableModel) FilterList(vm *hvue.VM, prop string) []*elements.ValText {
	rtm = &RipsiteTableModel{Object: vm.Object}
	count := map[string]int{}
	attribs := []string{}

	var translate func(string) string
	switch prop {
	case "Status":
		translate = func(val string) string {
			return fmrip.RipsiteStatusLabel(val)
		}
	default:
		translate = func(val string) string { return val }
	}

	for _, rsi := range rtm.Ripsiteinfos {
		attrib := rsi.Object.Get(prop).String()
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

func (rtm *RipsiteTableModel) FilteredStatusValue() []string {
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
