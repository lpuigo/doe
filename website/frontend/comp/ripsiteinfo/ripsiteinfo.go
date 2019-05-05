package ripsiteinfo

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

const (
	template1 string = `
<div> 
    <i class="fas fa-code-branch icon--left"></i><span>{{NbCommand}}&nbsp;</span>
    <i class="fas fa-share-alt icon--left"></i><span>{{NbTroncon}}&nbsp;</span>
    <i class="fas fa-home icon--left"></i>
	<span v-if="NbAvailLogement != NbLogement">{{NbAvailLogement}} / {{NbLogement}}</span>
	<span v-else>{{NbLogement}}</span>
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

/* TODO implement Ripsite & RipsiteInfoModel

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Registration Ripsite version

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("ripsite-info", ComponentRipsiteOptions()...)
}

func ComponentRipsiteOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template1),
		hvue.Props("ripsite"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewRipsiteInfoModel(vm)
		}),
		hvue.MethodsOf(&RipsiteInfoModel{}),
		hvue.Computed("NbCommand", func(vm *hvue.VM) interface{} {
			wim := &RipsiteInfoModel{Object: vm.Object}
			nbCommand, nbTroncon, nbAvailLogement, nbLogement := wim.Ripsite.GetInfo()
			wim.NbTroncon = nbTroncon
			wim.NbLogement = nbLogement
			wim.NbAvailLogement = nbAvailLogement
			return nbCommand
		}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type RipsiteInfoModel struct {
	*js.Object

	Ripsite         *fm.Worksite `js:"ripsite"`
	NbTroncon       int          `js:"NbTroncon"`
	NbLogement      int          `js:"NbLogement"`
	NbAvailLogement int          `js:"NbAvailLogement"`

	VM *hvue.VM `js:"VM"`
}

func NewRipsiteInfoModel(vm *hvue.VM) *RipsiteInfoModel {
	wim := &RipsiteInfoModel{Object: tools.O()}
	wim.VM = vm
	wim.Ripsite = nil
	wim.NbTroncon = 0
	wim.NbLogement = 0
	wim.NbAvailLogement = 0
	return wim
}
*/

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Registration RipsiteInfo version

func RegisterComponentRipsiteInfoInfo() hvue.ComponentOption {
	return hvue.Component("ripsiteinfo-info", ComponentRipsiteInfoOptions()...)
}

func ComponentRipsiteInfoOptions() []hvue.ComponentOption {
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
