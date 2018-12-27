package worksiteinfo

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/dates"
)

const template string = `
<div> 
    <i class="fas fa-sitemap icon--left"></i><span>{{NbCommand}}&nbsp;</span>
    <i class="fas fa-share-alt icon--left"></i><span>{{NbTroncon}}&nbsp;</span>
    <i class="fas fa-grip-vertical icon--left"></i><span>{{NbLogement}}</span>
</div>`

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Registration

func Register() {
	hvue.NewComponent("worksite-info",
		ComponentOptions()...,
	)
}

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("worksite-info", ComponentOptions()...)
}

func ComponentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.Props("worksite"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewWorksiteInfoModel(vm)
		}),
		hvue.MethodsOf(&WorksiteInfoModel{}),
		hvue.Computed("NbCommand", func(vm *hvue.VM) interface{} {
			wim := &WorksiteInfoModel{Object: vm.Object}
			nbCommand, nbTroncon, nbLogement := wim.Worksite.GetInfo()
			wim.NbTroncon = nbTroncon
			wim.NbLogement = nbLogement
			return nbCommand
		}),
		hvue.Filter("DateFormat", func(vm *hvue.VM, value *js.Object, args ...*js.Object) interface{} {
			return date.DateString(value.String())
		}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type WorksiteInfoModel struct {
	*js.Object

	Worksite   *fm.Worksite `js:"worksite"`
	NbTroncon  int          `js:"NbTroncon"`
	NbLogement int          `js:"NbLogement"`

	VM *hvue.VM `js:"VM"`
}

func NewWorksiteInfoModel(vm *hvue.VM) *WorksiteInfoModel {
	wim := &WorksiteInfoModel{Object: tools.O()}
	wim.VM = vm
	wim.Worksite = nil
	wim.NbTroncon = 0
	wim.NbLogement = 0
	return wim
}
