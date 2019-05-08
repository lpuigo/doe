package rippullingupdate

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/rippullingdistinfo"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	fmrip "github.com/lpuig/ewin/doe/website/frontend/model/ripsite"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/json"
	"strings"
)

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("rip-pulling-update", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		rippullingdistinfo.RegisterComponent(),
		hvue.Template(template),
		hvue.Props("value", "user"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewRipPullingUpdateModel(vm)
		}),
		hvue.MethodsOf(&RipPullingUpdateModel{}),
		hvue.Computed("filteredPullings", func(vm *hvue.VM) interface{} {
			rpum := RipPullingUpdateModelFromJS(vm.Object)
			return rpum.GetFilteredPullings()
		}),
		//hvue.Filter("FormatTronconRef", func(vm *hvue.VM, value *js.Object, args ...*js.Object) interface{} {
		//	rpum := RipPullingUpdateModelFromJS(vm.Object)
		//	t := &fm.Troncon{Object: value}
		//	return rpum.GetFormatTronconRef(t)
		//}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type RipPullingUpdateModel struct {
	*js.Object

	Ripsite *fmrip.Ripsite `js:"value"`
	//ReferenceRipsite *fmrip.Ripsite `js:"refRipsite"`
	User   *fm.User `js:"user"`
	Filter string   `js:"filter"`

	VM *hvue.VM `js:"VM"`
}

func NewRipPullingUpdateModel(vm *hvue.VM) *RipPullingUpdateModel {
	rpum := &RipPullingUpdateModel{Object: tools.O()}
	rpum.VM = vm
	rpum.Ripsite = fmrip.NewRisite()
	//rpum.ReferenceWorksite = nil
	rpum.User = nil
	rpum.Filter = ""
	return rpum
}

func RipPullingUpdateModelFromJS(o *js.Object) *RipPullingUpdateModel {
	return &RipPullingUpdateModel{Object: o}
}

func (rpum *RipPullingUpdateModel) GetFilteredPullings() []*fmrip.Pulling {
	if rpum.Filter == "" {
		return rpum.Ripsite.Pullings
	}
	res := []*fmrip.Pulling{}
	filter := strings.ToLower(rpum.Filter)
	for _, pulling := range rpum.Ripsite.Pullings {
		if strings.Contains(strings.ToLower(json.Stringify(pulling)), filter) {
			res = append(res, pulling)
		}
	}
	return res
}

func (rpum *RipPullingUpdateModel) TableRowClassName(rowInfo *js.Object) string {
	return ""
}

func (rpum *RipPullingUpdateModel) GetFirstPullingChunk(pulling *fmrip.Pulling) *fmrip.PullingChunk {
	return pulling.Chuncks[0]
}

func (rpum *RipPullingUpdateModel) GetLastPullingChunk(pulling *fmrip.Pulling) *fmrip.PullingChunk {
	return pulling.Chuncks[len(pulling.Chuncks)-1]
}

func (rpum *RipPullingUpdateModel) GetTroncon(vm *hvue.VM, trName string) *fmrip.Troncon {
	rpum = RipPullingUpdateModelFromJS(vm.Object)
	return rpum.Ripsite.Troncons[trName]
}

func (rpum *RipPullingUpdateModel) GetNode(vm *hvue.VM, nodeName string) *fmrip.Node {
	rpum = RipPullingUpdateModelFromJS(vm.Object)
	node, exist := rpum.Ripsite.Nodes[nodeName]
	if !exist {
		print(nodeName, "not found in", rpum.Ripsite.Pullings)
		return fmrip.NewNode()
	}
	return node
}
