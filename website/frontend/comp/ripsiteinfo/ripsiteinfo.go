package ripsiteinfo

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	fmrip "github.com/lpuig/ewin/doe/website/frontend/model/ripsite"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

const (
	template1 string = `
<div> 
    <i class="fas fa-arrows-alt-h icon--left"></i>
    <span v-if="NbAvailPulling != NbPulling">{{NbAvailPulling}} / {{NbPulling}}&nbsp;</span>
    <span v-else>{{NbPulling}}&nbsp;</span>
    <i class="fas fa-project-diagram icon--left"></i>
    <span v-if="NbAvailJunction != NbJunction">{{NbAvailJunction}} / {{NbJunction}}&nbsp;</span>
    <span v-else>{{NbJunction}}&nbsp;</span>
    <i class="fas fa-weight icon--left"></i>
	<span v-if="NbAvailMeasurement != NbMeasurement">{{NbAvailMeasurement}} / {{NbMeasurement}}</span>
	<span v-else>{{NbMeasurement}}</span>
</div>`

	template2 string = `
<div> 
    <i class="fas fa-arrows-alt-h icon--left"></i>
    <span v-if="value.NbPullingBlocked > 0">{{value.NbPulling - value.NbPullingBlocked}} / {{value.NbPulling}}&nbsp;</span>
    <span v-else>{{value.NbPulling}}&nbsp;</span>
    <i class="fas fa-project-diagram icon--left"></i>
    <span v-if="value.NbJunctionBlocked > 0">{{value.NbJunction - value.NbJunctionBlocked}} / {{value.NbJunction}}&nbsp;</span>
    <span v-else>{{value.NbJunction}}&nbsp;</span>
    <i class="fas fa-weight icon--left"></i>
    <span v-if="value.NbMeasurementBlocked > 0">{{value.NbMeasurement - value.NbMeasurementBlocked}} / {{value.NbMeasurement}}&nbsp;</span>
    <span v-else>{{value.NbMeasurement}}&nbsp;</span>
</div>`
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Registration Ripsite version

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("ripsite-info", componentRipsiteOptions()...)
}

func componentRipsiteOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template1),
		hvue.Props("ripsite"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewRipsiteInfoModel(vm)
		}),
		hvue.MethodsOf(&RipsiteInfoModel{}),
		hvue.Computed("NbAvailPulling", func(vm *hvue.VM) interface{} {
			wim := &RipsiteInfoModel{Object: vm.Object}
			NbAvailPulling, NbPulling, NbAvailJunction, NbJunction, NbAvailMeasurement, NbMeasurement := wim.Ripsite.GetInfo()
			wim.NbPulling = NbPulling
			wim.NbAvailJunction = NbAvailJunction
			wim.NbJunction = NbJunction
			wim.NbAvailMeasurement = NbAvailMeasurement
			wim.NbMeasurement = NbMeasurement
			return NbAvailPulling
		}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type RipsiteInfoModel struct {
	*js.Object

	Ripsite            *fmrip.Ripsite `js:"ripsite"`
	NbPulling          int            `js:"NbPulling"`
	NbJunction         int            `js:"NbJunction"`
	NbAvailJunction    int            `js:"NbAvailJunction"`
	NbMeasurement      int            `js:"NbMeasurement"`
	NbAvailMeasurement int            `js:"NbAvailMeasurement"`

	VM *hvue.VM `js:"VM"`
}

func NewRipsiteInfoModel(vm *hvue.VM) *RipsiteInfoModel {
	rim := &RipsiteInfoModel{Object: tools.O()}
	rim.VM = vm
	rim.Ripsite = nil
	rim.NbPulling = 0
	rim.NbJunction = 0
	rim.NbAvailJunction = 0
	rim.NbMeasurement = 0
	rim.NbAvailMeasurement = 0

	return rim
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Registration RipsiteInfo version

func RegisterComponentRipsiteInfoInfo() hvue.ComponentOption {
	return hvue.Component("ripsiteinfo-info", componentRipsiteInfoOptions()...)
}

func componentRipsiteInfoOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template2),
		hvue.Props("value"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewRipsiteInfoInfoModel(vm)
		}),
		hvue.MethodsOf(&RipsiteInfoInfoModel{}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type RipsiteInfoInfoModel struct {
	*js.Object

	RipsiteInfo *fm.RipsiteInfo `js:"value"`
	VM          *hvue.VM        `js:"VM"`
}

func NewRipsiteInfoInfoModel(vm *hvue.VM) *RipsiteInfoInfoModel {
	wim := &RipsiteInfoInfoModel{Object: tools.O()}
	wim.VM = vm
	wim.RipsiteInfo = nil
	return wim
}
