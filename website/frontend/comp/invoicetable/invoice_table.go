package invoicetable

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

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("invoice-table", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		worksiteinfo.RegisterComponentWorksiteInfo(),
		progressbar.RegisterComponent(),
		hvue.Template(template),
		hvue.Props("worksiteinfos"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewInvoiceTableModel(vm)
		}),
		hvue.MethodsOf(&InvoiceTableModel{}),
		hvue.Computed("filteredWorksites", func(vm *hvue.VM) interface{} {
			wtm := &InvoiceTableModel{Object: vm.Object}
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

type InvoiceTableModel struct {
	*js.Object

	Worksiteinfos []*fm.WorksiteInfo `js:"worksiteinfos"`
	Filter        string             `js:"filter"`

	VM *hvue.VM `js:"VM"`
}

func NewInvoiceTableModel(vm *hvue.VM) *InvoiceTableModel {
	wtm := &InvoiceTableModel{Object: tools.O()}
	wtm.Worksiteinfos = nil
	wtm.Filter = ""
	wtm.VM = vm
	return wtm
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Event Related Methods

func (itm *InvoiceTableModel) SetSelectedWorksite(wsi *fm.WorksiteInfo) {
	itm.VM.Emit("selected_worksite", wsi.Id)
}

func (itm *InvoiceTableModel) ExpandRow(vm *hvue.VM, ws *fm.WorksiteInfo, others *js.Object) {
	print("Others :", others)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Formatting Related Methods

func (itm *InvoiceTableModel) TableRowClassName(rowInfo *js.Object) string {
	wsi := &fm.WorksiteInfo{Object: rowInfo.Get("row")}
	return fm.WorksiteRowClassName(wsi.Status)
}

func (itm *InvoiceTableModel) HeaderCellStyle() string {
	return "background: #a1e6e6;"
}

func (itm *InvoiceTableModel) FormatDate(r, c *js.Object, d string) string {
	return date.DateString(d)
}

func (itm *InvoiceTableModel) FormatStatus(r, c *js.Object, d string) string {
	return fm.WorksiteStatusLabel(d)
}

func (itm *InvoiceTableModel) SortStatus(a, b *fm.Worksite) int {
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

func (itm *InvoiceTableModel) FilterHandler(value string, p *js.Object, col *js.Object) bool {
	prop := col.Get("property").String()
	return p.Get(prop).String() == value
}

func (itm *InvoiceTableModel) FilterList(vm *hvue.VM, prop string) []*elements.ValText {
	itm = &InvoiceTableModel{Object: vm.Object}
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

	for _, ws := range itm.Worksiteinfos {
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

func (itm *InvoiceTableModel) FilteredStatusValue() []string {
	res := []string{
		fm.WsStatusAttachment,
		fm.WsStatusInvoice,
		fm.WsStatusPayment,
		//		fm.WsStatusRework,
	}
	return res
}
