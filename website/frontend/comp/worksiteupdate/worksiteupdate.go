package worksiteupdate

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/orderedit"
	"github.com/lpuig/ewin/doe/website/frontend/comp/ptedit"
	"github.com/lpuig/ewin/doe/website/frontend/comp/tronconstatustag"
	"github.com/lpuig/ewin/doe/website/frontend/comp/worksitestatustag"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/autocomplete"
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
		ptedit.RegisterComponent(),
		orderedit.RegisterComponent(),
		tronconstatustag.RegisterComponent(),
		worksitestatustag.RegisterComponent(),
		hvue.Template(template),
		hvue.Props("worksite"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewWorksiteUpdateModel(vm)
		}),
		hvue.MethodsOf(&WorksiteUpdateModel{}),
		//hvue.Computed("HasChanged", func(vm *hvue.VM) interface{} {
		//	wdm := &WorksiteUpdateModel{Object: vm.Object}
		//	if wdm.ReferenceWorksite.Object == nil {
		//		wdm.ReferenceWorksite = wdm.Worksite.Clone()
		//		return wdm.Worksite.Dirty
		//	}
		//	s1 := wdm.Worksite.SearchInString()
		//	s2 := wdm.ReferenceWorksite.SearchInString()
		//	wdm.Worksite.Dirty = s1 != s2
		//	return wdm.Worksite.Dirty
		//}),
		hvue.Computed("filteredTroncons", func(vm *hvue.VM) interface{} {
			wdm := &WorksiteUpdateModel{Object: vm.Object}
			return wdm.GetTroncons()
		}),
		hvue.Filter("FormatTronconRef", func(vm *hvue.VM, value *js.Object, args ...*js.Object) interface{} {
			wdm := &WorksiteUpdateModel{Object: vm.Object}
			t := &fm.Troncon{Object: value}
			return wdm.GetFormatTronconRef(t)
		}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type WorksiteUpdateModel struct {
	*js.Object

	Worksite          *fm.Worksite `js:"worksite"`
	ReferenceWorksite *fm.Worksite `js:"refWorksite"`
	Filter            string       `js:"filter"`

	VM *hvue.VM `js:"VM"`
}

func NewWorksiteUpdateModel(vm *hvue.VM) *WorksiteUpdateModel {
	wum := &WorksiteUpdateModel{Object: tools.O()}
	wum.VM = vm
	wum.Worksite = nil
	wum.ReferenceWorksite = nil
	wum.Filter = ""
	return wum
}

func (wum *WorksiteUpdateModel) DOEArchive() string {
	url := "/api/worksites/" + strconv.Itoa(wum.Worksite.Id) + "/zip"
	return url
}

func (wum *WorksiteUpdateModel) GetTroncons() []*OrderTroncon {
	res := []*OrderTroncon{}
	for _, o := range wum.Worksite.Orders {
		tres := []*fm.Troncon{}
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

func (wum *WorksiteUpdateModel) GetFormatTronconRef(t *fm.Troncon) string {
	return t.Pb.Ref + " / " + t.Pb.RefPt
}

func (wum *WorksiteUpdateModel) TextFiltered(t *fm.Troncon) bool {
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

func (wum *WorksiteUpdateModel) CheckSignature(t *fm.Troncon) {
	// t should be a OrderTroncon, bur gopherjs reflection seems to fail
	// also working with
	//func (wumm *WorksiteUpdateModalModel) CheckSignature(o *js.Object) {
	//	NewOrderTronconFromJS(o).CheckSignature()
	//}
	t.CheckSignature()
}

func (wum *WorksiteUpdateModel) UserSearch(vm *hvue.VM, query string, callback *js.Object) {
	users := fm.GetTeamUsers()

	q := strings.ToLower(query)

	res := []*autocomplete.Result{}
	for _, u := range users {
		if q == "" || strings.Contains(strings.ToLower(u), q) {
			res = append(res, autocomplete.NewResult(u))
		}
	}
	callback.Invoke(res)
}

//func (wum *WorksiteUpdateModel) DeleteOrder(vm *hvue.VM, i int) {
//	wum = &WorksiteUpdateModel{Object: vm.Object}
//	wum.Worksite.DeleteOrder(i)
//}
//
//func (wum *WorksiteUpdateModel) AddOrder(vm *hvue.VM) {
//	wum = &WorksiteUpdateModel{Object: vm.Object}
//	wum.Worksite.AddOrder()
//}
//
//func (wum *WorksiteUpdateModel) Save(vm *hvue.VM) {
//	wum = &WorksiteUpdateModel{Object: vm.Object}
//	vm.Emit("save_worksite", wum.Worksite)
//}
//
//func (wum *WorksiteUpdateModel) Undo(vm *hvue.VM) {
//	wum = &WorksiteUpdateModel{Object: vm.Object}
//	wum.Worksite.Copy(wum.ReferenceWorksite)
//}
//
//func (wdm *WorksiteUpdateModel) CheckDoeDate(vm *hvue.VM) {
//	wdm = &WorksiteUpdateModel{Object: vm.Object}
//	if tools.Empty(wdm.Worksite.DoeDate) {
//		wdm.Worksite.Status = fm.WsStatusDOE
//		return
//	}
//	wdm.Worksite.Status = fm.WsStatusDone
//}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp OrderTroncon Model

type OrderTroncon struct {
	*fm.Troncon
	Order string `js:"Order"`
	Span  int    `js:"Span"`
}

func NewOrderTroncon(t *fm.Troncon, order string, span int) *OrderTroncon {
	ot := &OrderTroncon{Troncon: t}
	ot.Order = order
	ot.Span = span
	return ot
}

func NewOrderTronconFromJS(o *js.Object) *OrderTroncon {
	return &OrderTroncon{Troncon: &fm.Troncon{Object: o}}
}
