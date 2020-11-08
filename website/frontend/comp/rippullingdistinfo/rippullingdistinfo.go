package rippullingdistinfo

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	fmrip "github.com/lpuig/ewin/doe/website/frontend/model/ripsite"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

const (
	template1 string = `
<div>
	<b v-if="total > 0">Total: {{total}}m </b>
	<span v-if="Love > 0">Lov.: {{Love}}m </span>
	<span v-if="Underground > 0">Sout.: {{Underground}}m </span>
	<span v-if="Aerial > 0" class="pulling-aerial">Aér.: {{Aerial}}m </span>
	<span v-if="Building > 0" class="pulling-aerial">Faç.: {{Building}}m</span>
</div>
`
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Registration Pulling version

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("pulling-distances-info", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template1),
		hvue.Props("value"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewRipPullingDistanceInfoModel(vm)
		}),
		hvue.MethodsOf(&RipPullingDistanceInfoModel{}),
		hvue.Computed("total", func(vm *hvue.VM) interface{} {
			rim := &RipPullingDistanceInfoModel{Object: vm.Object}
			var total int
			total, rim.Love, rim.Underground, rim.Aerial, rim.Building = rim.Pulling.GetDists()
			nbval := 0
			if rim.Love > 0 {
				nbval++
			}
			if rim.Underground > 0 {
				nbval++
			}
			if rim.Aerial > 0 {
				nbval++
			}
			if rim.Building > 0 {
				nbval++
			}
			if nbval == 1 {
				total = 0
			}
			return total
		}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type RipPullingDistanceInfoModel struct {
	*js.Object

	Pulling     *fmrip.Pulling `js:"value"`
	Love        int            `js:"Love"`
	Underground int            `js:"Underground"`
	Aerial      int            `js:"Aerial"`
	Building    int            `js:"Building"`

	VM *hvue.VM `js:"VM"`
}

func NewRipPullingDistanceInfoModel(vm *hvue.VM) *RipPullingDistanceInfoModel {
	rpdim := &RipPullingDistanceInfoModel{Object: tools.O()}
	rpdim.VM = vm
	rpdim.Pulling = nil
	rpdim.Love = 0
	rpdim.Underground = 0
	rpdim.Aerial = 0
	rpdim.Building = 0
	return rpdim
}
