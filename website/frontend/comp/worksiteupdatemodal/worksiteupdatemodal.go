package worksiteupdatemodal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/tronconstatustag"
	wem "github.com/lpuig/ewin/doe/website/frontend/comp/worksiteeditmodal"
	"github.com/lpuig/ewin/doe/website/frontend/comp/worksiteinfo"
	"github.com/lpuig/ewin/doe/website/frontend/comp/worksitestatustag"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"strings"
)

type WorksiteUpdateModalModel struct {
	*wem.WorksiteEditModalModel

	Filter string `js:"filter"`
}

func NewWorksiteUpdateModalModel(vm *hvue.VM) *WorksiteUpdateModalModel {
	wumm := &WorksiteUpdateModalModel{WorksiteEditModalModel: wem.NewWorksiteEditModalModel(vm)}
	wumm.Filter = ""
	return wumm
}

func NewWorksiteUpdateModalModelFromJS(o *js.Object) *WorksiteUpdateModalModel {
	wumm := &WorksiteUpdateModalModel{WorksiteEditModalModel: &wem.WorksiteEditModalModel{Object: o}}
	return wumm
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Component Methods

func Register() {
	hvue.NewComponent("worksite-update-modal",
		ComponentOptions()...,
	)
}

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("worksite-update-modal", ComponentOptions()...)
}

func ComponentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		worksiteinfo.RegisterComponent(),
		tronconstatustag.RegisterComponent(),
		worksitestatustag.RegisterComponent(),
		hvue.Template(template),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewWorksiteUpdateModalModel(vm)
		}),
		hvue.MethodsOf(&WorksiteUpdateModalModel{}),

		hvue.Computed("filteredTroncons", func(vm *hvue.VM) interface{} {
			m := NewWorksiteUpdateModalModelFromJS(vm.Object)
			return m.GetTroncons()
		}),
		hvue.Computed("hasChanged", func(vm *hvue.VM) interface{} {
			m := NewWorksiteUpdateModalModelFromJS(vm.Object)
			return m.HasChanged()
		}),
		hvue.Computed("hasWarning", func(vm *hvue.VM) interface{} {
			//m := &WorksiteEditModalModel{Object: vm.Object}
			//if len(m.CurrentProject.Audits) > 0 {
			//	return "warning"
			//}
			return "success"
		}),
		hvue.Filter("FormatTronconRef", func(vm *hvue.VM, value *js.Object, args ...*js.Object) interface{} {
			m := NewWorksiteUpdateModalModelFromJS(vm.Object)
			t := &fm.Troncon{Object: value}
			return m.GetFormatTronconRef(t)
		}),
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Modal Methods

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

func (wumm *WorksiteUpdateModalModel) GetTroncons() []*OrderTroncon {
	res := []*OrderTroncon{}
	for _, o := range wumm.CurrentWorksite.Orders {
		tres := []*fm.Troncon{}
		for _, t := range o.Troncons {
			if wumm.TextFiltered(t) {
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

func (wumm *WorksiteUpdateModalModel) TableRowClassName(rowInfo *js.Object) string {
	//wsi := &fm.WorksiteInfo{Object: rowInfo.Get("row")}
	return ""
}

func (wumm *WorksiteUpdateModalModel) OrderSpanMethod(o *js.Object) interface{} {
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

func (wumm *WorksiteUpdateModalModel) GetFormatTronconRef(t *fm.Troncon) string {
	return t.Pb.Ref + " / " + t.Pb.RefPt
}

func (wumm *WorksiteUpdateModalModel) TextFiltered(t *fm.Troncon) bool {
	filter := wumm.Filter
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
