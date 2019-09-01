package reworkedit

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/worksitestatustag"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/model/worksite"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Registration

func Register() {
	hvue.NewComponent("rework-edit",
		componentOptions()...,
	)
}

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("rework-edit", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		worksitestatustag.RegisterComponent(),
		hvue.Template(template),
		hvue.Props("worksite", "user"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewReworkEditModel(vm)
		}),
		hvue.Computed("filteredReworks", func(vm *hvue.VM) interface{} {
			m := ReworkEditModelFromJS(vm.Object)
			return m.GetReworks()
		}),
		hvue.MethodsOf(&ReworkEditModel{}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type ReworkEditModel struct {
	*js.Object

	Worksite *worksite.Worksite `js:"worksite"`
	User     *fm.User           `js:"user"`

	VM *hvue.VM `js:"VM"`
}

func NewReworkEditModel(vm *hvue.VM) *ReworkEditModel {
	wum := &ReworkEditModel{Object: tools.O()}
	wum.VM = vm
	wum.Worksite = worksite.NewWorkSite()
	wum.User = fm.NewUser()
	return wum
}

func ReworkEditModelFromJS(o *js.Object) *ReworkEditModel {
	return &ReworkEditModel{Object: o}
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Modal Methods

//func (rem *ReworkEditModel) TableRowClassName(rowInfo *js.Object) string {
//	//wsi := &fm.WorksiteInfo{Object: rowInfo.Get("row")}
//	return ""
//}

func (rem *ReworkEditModel) GetReworks() []*worksite.Defect {
	return rem.Worksite.Rework.Defects
}

func (rem *ReworkEditModel) GetPTs(vm *hvue.VM) []*elements.ValueLabel {
	rem = ReworkEditModelFromJS(vm.Object)
	res := []*elements.ValueLabel{}
	for _, o := range rem.Worksite.Orders {
		for _, t := range o.Troncons {
			label := t.Pb.Ref + " / " + t.Pb.RefPt + " (" + t.Ref + ")"
			res = append(res, elements.NewValueLabel(t.Pb.RefPt, label))
		}
	}
	return res
}

func (rem *ReworkEditModel) AddAllDefect(vm *hvue.VM) {
	rem = ReworkEditModelFromJS(vm.Object)
	for _, vl := range rem.GetPTs(vm) {
		rem.addDefect(vl.Value)
	}
}

func (rem *ReworkEditModel) addDefect(ptName string) {
	r := rem.Worksite.Rework
	d := worksite.NewDefect()
	d.PT = ptName
	d.SubmissionDate = r.ControlDate
	r.Defects = append(r.Defects, d)
}

func (rem *ReworkEditModel) RemoveDefect(vm *hvue.VM, i int) {
	rem = ReworkEditModelFromJS(vm.Object)
	rem.Worksite.Rework.Object.Get("Defects").Call("splice", i, 1)
}
