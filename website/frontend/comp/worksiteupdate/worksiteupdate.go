package worksiteupdate

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/tronconstatustag"
	"github.com/lpuig/ewin/doe/website/frontend/comp/worksitestatustag"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/model/worksite"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	date "github.com/lpuig/ewin/doe/website/frontend/tools/date"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
	"strconv"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Registration

func Register() {
	hvue.NewComponent("worksite-update",
		ComponentOptions()...,
	)
}

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("worksite-update", ComponentOptions()...)
}

func ComponentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		tronconstatustag.RegisterComponent(),
		worksitestatustag.RegisterComponent(),
		hvue.Template(template),
		hvue.Props("worksite", "user"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewWorksiteUpdateModel(vm)
		}),
		hvue.MethodsOf(&WorksiteUpdateModel{}),
		hvue.Computed("filteredTroncons", func(vm *hvue.VM) interface{} {
			wdm := &WorksiteUpdateModel{Object: vm.Object}
			return wdm.GetTroncons()
		}),
		hvue.Filter("FormatTronconRef", func(vm *hvue.VM, value *js.Object, args ...*js.Object) interface{} {
			wdm := &WorksiteUpdateModel{Object: vm.Object}
			t := &worksite.Troncon{Object: value}
			return wdm.GetFormatTronconRef(t)
		}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type WorksiteUpdateModel struct {
	*js.Object

	Worksite          *worksite.Worksite `js:"worksite"`
	ReferenceWorksite *worksite.Worksite `js:"refWorksite"`
	User              *fm.User           `js:"user"`
	Filter            string             `js:"filter"`

	VM *hvue.VM `js:"VM"`
}

func NewWorksiteUpdateModel(vm *hvue.VM) *WorksiteUpdateModel {
	wum := &WorksiteUpdateModel{Object: tools.O()}
	wum.VM = vm
	wum.Worksite = nil
	wum.ReferenceWorksite = nil
	wum.User = nil
	wum.Filter = ""
	return wum
}

func WorksiteUpdateModelFromJS(o *js.Object) *WorksiteUpdateModel {
	return &WorksiteUpdateModel{Object: o}
}

func (wum *WorksiteUpdateModel) DOEArchive() string {
	url := "/api/worksites/" + strconv.Itoa(wum.Worksite.Id) + "/zip"
	return url
}

func (wum *WorksiteUpdateModel) GetTroncons() []*OrderTroncon {
	res := []*OrderTroncon{}
	for _, o := range wum.Worksite.Orders {
		tres := []*worksite.Troncon{}
		for _, t := range o.Troncons {
			if wum.TextFiltered(t) {
				tres = append(tres, t)
			}
		}
		for i, t := range tres {
			span := 0
			if i == 0 {
				span = len(tres)
			}
			res = append(res, NewOrderTroncon(t, o.Ref, span))
		}
	}
	return res
}

func (wum *WorksiteUpdateModel) TableRowClassName(rowInfo *js.Object) string {
	//wsi := &fm.WorksiteInfo{Object: rowInfo.Get("row")}
	return ""
}

func (wum *WorksiteUpdateModel) OrderSpanMethod(o *js.Object) interface{} {
	row := NewOrderTronconFromJS(o.Get("row"))
	col := o.Get("columnIndex").Int()
	if col == 0 {
		if row.Span == 0 {
			return js.M{
				"rowspan": 0,
				"colspan": 0,
			}
		}
		return js.M{
			"rowspan": row.Span,
			"colspan": 1,
		}
	}
	return js.Undefined
}

func (wum *WorksiteUpdateModel) GetFormatTronconRef(t *worksite.Troncon) string {
	return t.Pb.Ref + " / " + t.Pb.RefPt
}

func (wum *WorksiteUpdateModel) TextFiltered(t *worksite.Troncon) bool {
	filter := wum.Filter
	if filter == "" {
		return true
	}
	expected := true
	if strings.HasPrefix(filter, `\`) {
		if len(filter) == 1 {
			return true
		}
		expected = false
		filter = filter[1:]
	}
	return strings.Contains(t.SearchInString(), filter) == expected
}

func (wum *WorksiteUpdateModel) SetInstallDate(t *worksite.Troncon) {
	if tools.Empty(t.InstallDate) {
		t.InstallDate = date.TodayAfter(0)
	}
}

func (wum *WorksiteUpdateModel) SetMeasureDate(t *worksite.Troncon) {
	if tools.Empty(t.MeasureDate) {
		t.MeasureDate = date.TodayAfter(0)
	}
}

func (wum *WorksiteUpdateModel) CheckSignature(t *worksite.Troncon) {
	// t should be a OrderTroncon, bur gopherjs reflection seems to fail
	// also working with
	//func (wumm *WorksiteUpdateModalModel) CheckSignature(o *js.Object) {
	//	NewOrderTronconFromJS(o).CheckSignature()
	//}
	t.CheckSignature()
}

func (wum *WorksiteUpdateModel) CheckInstallDate(t *worksite.Troncon) {
	// t should be a OrderTroncon, bur gopherjs reflection seems to fail
	// also working with
	//func (wumm *WorksiteUpdateModalModel) CheckSignature(o *js.Object) {
	//	NewOrderTronconFromJS(o).CheckSignature()
	//}
	// If blocked without any installation done yet => Set InstallDate as Blockage detection date
	if t.Blockage && tools.Empty(t.InstallActor) && tools.Empty(t.InstallDate) {
		t.InstallDate = date.TodayAfter(0)
	}
	// If unblocked without any installation done yet => Delete InstallDate
	if !t.Blockage && tools.Empty(t.InstallActor) {
		t.InstallDate = ""
	}
}

//func (wum *WorksiteUpdateModel) UserSearch(vm *hvue.VM, query string, callback *js.Object) {
//	wum = WorksiteUpdateModelFromJS(vm.Object)
//
//	q := strings.ToLower(query)
//
//	res := []*autocomplete.Result{}
//	for _, u := range wum.User.Teams {
//		if q == "" || strings.Contains(strings.ToLower(u), q) {
//			res = append(res, autocomplete.NewResult(u))
//		}
//	}
//	callback.Invoke(res)
//}

func (wum *WorksiteUpdateModel) GetTeams(vm *hvue.VM) []*elements.ValueLabel {
	wum = WorksiteUpdateModelFromJS(vm.Object)
	return wum.User.GetTeamValueLabelsFor(wum.Worksite.Client)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp OrderTroncon Model

type OrderTroncon struct {
	*worksite.Troncon
	Order string `js:"Order"`
	Span  int    `js:"Span"`
}

func NewOrderTroncon(t *worksite.Troncon, order string, span int) *OrderTroncon {
	ot := &OrderTroncon{Troncon: t}
	ot.Order = order
	ot.Span = span
	return ot
}

func NewOrderTronconFromJS(o *js.Object) *OrderTroncon {
	return &OrderTroncon{Troncon: &worksite.Troncon{Object: o}}
}
