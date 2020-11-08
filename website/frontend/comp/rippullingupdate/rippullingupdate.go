package rippullingupdate

import (
	"strings"
	"time"

	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/rippullingdistinfo"
	"github.com/lpuig/ewin/doe/website/frontend/comp/ripstateupdate"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	fmrip "github.com/lpuig/ewin/doe/website/frontend/model/ripsite"
	"github.com/lpuig/ewin/doe/website/frontend/model/ripsite/ripconst"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("rip-pulling-update", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		rippullingdistinfo.RegisterComponent(),
		ripstateupdate.RegisterComponent(),
		hvue.Template(template),
		hvue.Props("value", "user", "filter", "filtertype"),
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
	User       *fm.User `js:"user"`
	Filter     string   `js:"filter"`
	FilterType string   `js:"filtertype"`
	SizeLimit  int      `js:"SizeLimit"`

	VM *hvue.VM `js:"VM"`
}

func NewRipPullingUpdateModel(vm *hvue.VM) *RipPullingUpdateModel {
	rpum := &RipPullingUpdateModel{Object: tools.O()}
	rpum.VM = vm
	rpum.Ripsite = fmrip.NewRisite()
	//rpum.ReferenceWorksite = nil
	rpum.User = nil
	rpum.Filter = ""
	rpum.FilterType = ripconst.FilterValueAll
	rpum.SetSizeLimit()
	return rpum
}

func RipPullingUpdateModelFromJS(o *js.Object) *RipPullingUpdateModel {
	return &RipPullingUpdateModel{Object: o}
}

func (rpum *RipPullingUpdateModel) GetFilteredPullings() []*fmrip.Pulling {
	if rpum.FilterType == ripconst.FilterValueAll && rpum.Filter == "" {
		return rpum.GetSizeLimitedResult(rpum.Ripsite.Pullings)
	}
	res := []*fmrip.Pulling{}
	expected := strings.ToUpper(rpum.Filter)
	filter := func(p *fmrip.Pulling) bool {
		sis := p.SearchString(rpum.FilterType)
		if sis == "" {
			return false
		}
		return strings.Contains(strings.ToUpper(sis), expected)
	}

	for _, pulling := range rpum.Ripsite.Pullings {
		if filter(pulling) {
			res = append(res, pulling)
		}
	}
	return rpum.GetSizeLimitedResult(res)
}

func (rpum *RipPullingUpdateModel) TableRowClassName(rowInfo *js.Object) string {
	pulling := &fmrip.Pulling{Object: rowInfo.Get("row")}
	return pulling.State.GetRowStyle()
}

func (rpum *RipPullingUpdateModel) GetFirstPullingChunk(pulling *fmrip.Pulling) *fmrip.PullingChunk {
	return pulling.Chuncks[0]
}

func (rpum *RipPullingUpdateModel) GetPullingTypeClass(pulling *fmrip.Pulling) string {
	_, _, _, aerial, building := pulling.GetDists()
	if aerial+building > 0 {
		return "pulling-aerial"
	}
	return ""
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

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Size Related Methods

const (
	sizeLimitDefault int = 15
	sizeLimitTimer       = 200
)

func (rpum *RipPullingUpdateModel) GetSizeLimitedResult(res []*fmrip.Pulling) []*fmrip.Pulling {
	if len(res) == rpum.SizeLimit {
		return res
	}
	if len(res) > sizeLimitDefault {
		rpum.ResetSizeLimit(len(res))
		return res[:sizeLimitDefault]
	}
	return res
}

func (rpum *RipPullingUpdateModel) SetSizeLimit() {
	rpum.SizeLimit = -1
}

func (rpum *RipPullingUpdateModel) ResetSizeLimit(size int) {
	go func() {
		time.Sleep(sizeLimitTimer * time.Millisecond)
		rpum.SizeLimit = size
	}()
}
