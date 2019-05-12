package ripmeasurementupdate

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/ripstateupdate"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	fmrip "github.com/lpuig/ewin/doe/website/frontend/model/ripsite"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/json"
	"strings"
)

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("rip-measurement-update", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		ripstateupdate.RegisterComponent(),
		hvue.Template(template),
		hvue.Props("value", "user", "filter"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewRipMeasurementUpdateModel(vm)
		}),
		hvue.MethodsOf(&RipMeasurementUpdateModel{}),
		hvue.Computed("filteredMeasurements", func(vm *hvue.VM) interface{} {
			rpum := RipMeasurementUpdateModelFromJS(vm.Object)
			return rpum.GetFilteredMeasurements()
		}),
		//hvue.Filter("FormatTronconRef", func(vm *hvue.VM, value *js.Object, args ...*js.Object) interface{} {
		//	rpum := RipMeasurementUpdateModelFromJS(vm.Object)
		//	t := &fm.Troncon{Object: value}
		//	return rpum.GetFormatTronconRef(t)
		//}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type RipMeasurementUpdateModel struct {
	*js.Object

	Ripsite *fmrip.Ripsite `js:"value"`
	//ReferenceRipsite *fmrip.Ripsite `js:"refRipsite"`
	User   *fm.User `js:"user"`
	Filter string   `js:"filter"`

	VM *hvue.VM `js:"VM"`
}

func NewRipMeasurementUpdateModel(vm *hvue.VM) *RipMeasurementUpdateModel {
	rmum := &RipMeasurementUpdateModel{Object: tools.O()}
	rmum.VM = vm
	rmum.Ripsite = fmrip.NewRisite()
	//rmum.ReferenceWorksite = nil
	rmum.User = nil
	rmum.Filter = ""
	return rmum
}

func RipMeasurementUpdateModelFromJS(o *js.Object) *RipMeasurementUpdateModel {
	return &RipMeasurementUpdateModel{Object: o}
}

func (rmum *RipMeasurementUpdateModel) GetFilteredMeasurements() []*fmrip.Measurement {
	if rmum.Filter == "" {
		return rmum.Ripsite.Measurements
	}
	res := []*fmrip.Measurement{}
	filter := strings.ToLower(rmum.Filter)
	for _, meas := range rmum.Ripsite.Measurements {
		if strings.Contains(strings.ToLower(json.Stringify(meas)), filter) {
			res = append(res, meas)
		}
	}
	return res
}

func (rmum *RipMeasurementUpdateModel) TableRowClassName(rowInfo *js.Object) string {
	junction := &fmrip.Measurement{Object: rowInfo.Get("row")}
	return junction.State.GetRowStyle()
}

func (rmum *RipMeasurementUpdateModel) GetNode(vm *hvue.VM, name string) *fmrip.Node {
	rmum = RipMeasurementUpdateModelFromJS(vm.Object)
	return rmum.Ripsite.Nodes[name]
}

func (rmum *RipMeasurementUpdateModel) GetDestNodeDist(vm *hvue.VM, meas *fmrip.Measurement) int {
	rmum = RipMeasurementUpdateModelFromJS(vm.Object)
	node := rmum.Ripsite.Nodes[meas.DestNodeName]
	return node.DistFromPm
}
