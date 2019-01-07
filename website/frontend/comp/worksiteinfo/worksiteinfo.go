package worksiteinfo

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

const (
	template1 string = `
<div> 
    <i class="fas fa-sitemap icon--left"></i><span>{{NbCommand}}&nbsp;</span>
    <i class="fas fa-share-alt icon--left"></i><span>{{NbTroncon}}&nbsp;</span>
    <i class="fas fa-grip-vertical icon--left"></i>
	<span v-if="NbAvailLogement != NbLogement">{{NbAvailLogement}} / {{NbLogement}}</span>
	<span v-else>{{NbLogement}}</span>
</div>`

	template2 string = `
<div> 
    <i class="fas fa-sitemap icon--left"></i><span>{{value.NbOrder}}&nbsp;</span>
    <i class="fas fa-share-alt icon--left"></i><span>{{value.NbTroncon}}&nbsp;</span>
    <i class="fas fa-grip-vertical icon--left"></i>
	<span v-if="value.NbElBlocked > 0">{{value.NbElTotal - value.NbElBlocked}} / {{value.NbElTotal}}</span>
	<span v-else>{{value.NbElTotal}}</span>
</div>`
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Registration Worksite version

func RegisterComponentWorksite() hvue.ComponentOption {
	return hvue.Component("worksite-info", ComponentWorksiteOptions()...)
}

func ComponentWorksiteOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template1),
		hvue.Props("worksite"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewWorksiteInfoModel(vm)
		}),
		hvue.MethodsOf(&WorksiteInfoModel{}),
		hvue.Computed("NbCommand", func(vm *hvue.VM) interface{} {
			wim := &WorksiteInfoModel{Object: vm.Object}
			nbCommand, nbTroncon, nbAvailLogement, nbLogement := wim.Worksite.GetInfo()
			wim.NbTroncon = nbTroncon
			wim.NbLogement = nbLogement
			wim.NbAvailLogement = nbAvailLogement
			return nbCommand
		}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type WorksiteInfoModel struct {
	*js.Object

	Worksite        *fm.Worksite `js:"worksite"`
	NbTroncon       int          `js:"NbTroncon"`
	NbLogement      int          `js:"NbLogement"`
	NbAvailLogement int          `js:"NbAvailLogement"`

	VM *hvue.VM `js:"VM"`
}

func NewWorksiteInfoModel(vm *hvue.VM) *WorksiteInfoModel {
	wim := &WorksiteInfoModel{Object: tools.O()}
	wim.VM = vm
	wim.Worksite = nil
	wim.NbTroncon = 0
	wim.NbLogement = 0
	wim.NbAvailLogement = 0
	return wim
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Registration Worksite version

func RegisterComponentWorksiteInfo() hvue.ComponentOption {
	return hvue.Component("worksiteinfo-info", ComponentWorksiteInfoOptions()...)
}

func ComponentWorksiteInfoOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template2),
		hvue.Props("value"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewWorksiteInfoModel(vm)
		}),
		hvue.MethodsOf(&WorksiteInfoModel{}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type WorksiteInfoInfoModel struct {
	*js.Object

	WorksiteInfo *fm.WorksiteInfo `js:"value"`
	VM           *hvue.VM         `js:"VM"`
}

func NewWorksiteInfoInfoModel(vm *hvue.VM) *WorksiteInfoInfoModel {
	wim := &WorksiteInfoInfoModel{Object: tools.O()}
	wim.VM = vm
	wim.WorksiteInfo = nil
	return wim
}
