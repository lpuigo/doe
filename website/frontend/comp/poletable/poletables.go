package poletable

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	ps "github.com/lpuig/ewin/doe/website/frontend/model/polesite"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

const template string = `<el-container  style="height: 100%; padding: 0px">
    <el-header style="height: auto; margin-top: 5px">
        <el-row type="flex" align="middle" :gutter="5">
            <el-col :span="2" style="text-align: right"><span>Mode d'affichage:</span></el-col>
            <el-col :span="10">
                <el-radio-group v-model="context.Mode" @change="ChangeMode" size="mini">
                    <el-tooltip content="Création de poteaux" placement="bottom" effect="light" open-delay="500">
                        <el-radio-button label="creation">Création</el-radio-button>
                    </el-tooltip>
                    <el-tooltip content="Planification d'activité" placement="bottom" effect="light" open-delay="500">
                        <el-radio-button label="followup">Planification</el-radio-button>
                    </el-tooltip>
                    <el-tooltip content="Mise a jour de l'avancement" placement="bottom" effect="light" open-delay="500">
                        <el-radio-button label="billing">Avancement</el-radio-button>
                    </el-tooltip>
                </el-radio-group>
            </el-col>  
        </el-row>
    </el-header>
    <div style="height: 100%;overflow-x: hidden;overflow-y: auto;padding: 0px 0px; margin-top: 8px">
		<pole-table-creation v-if="context.Mode == 'creation'"
				:user="user"
				:polesite="polesite"
				:filter="filter"
				:filtertype="filtertype"
				:context.sync="context"
				@update:context="ChangeMode"
		></pole-table-creation>
		<pole-table-followup v-if="context.Mode == 'followup'"
				:user="user"
				:polesite="polesite"
				:filter="filter"
				:filtertype="filtertype"
				:context.sync="context"
				@update:context="ChangeMode"
		></pole-table-followup>
		<pole-table-billing v-if="context.Mode == 'billing'"
				:user="user"
				:polesite="polesite"
				:filter="filter"
				:filtertype="filtertype"
				:context.sync="context"
				@update:context="ChangeMode"
		></pole-table-billing>
    </div>
</el-container>
`

//@pole-selected="SetSelectedPole"

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Registration

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("pole-table", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		registerComponentTable("creation"),
		registerComponentTable("followup"),
		registerComponentTable("billing"),
		hvue.Template(template),
		hvue.Props("user", "polesite", "filter", "filtertype", "context"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewPoleTablesModel(vm)
		}),
		hvue.MethodsOf(&PoleTablesModel{}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type PoleTablesModel struct {
	*js.Object

	Polesite   *ps.Polesite `js:"polesite"`
	User       *fm.User     `js:"user"`
	Filter     string       `js:"filter"`
	FilterType string       `js:"filtertype"`
	Context    *Context     `js:"context"`

	VM *hvue.VM `js:"VM"`
}

func NewPoleTablesModel(vm *hvue.VM) *PoleTablesModel {
	rtm := &PoleTablesModel{Object: tools.O()}
	rtm.Polesite = nil
	rtm.User = fm.NewUser()
	rtm.Filter = ""
	rtm.FilterType = ""
	rtm.Context = NewContext("")
	rtm.VM = vm
	return rtm
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Actions related Methods

func (ptm *PoleTablesModel) ChangeMode(vm *hvue.VM) {
	ptm = &PoleTablesModel{Object: vm.Object}
	vm.Emit("update:context", ptm.Context)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Row Filtering Related Methods

//func (ptm *PoleTablesModel) GetFilteredPole() []*ps.Pole {
//	if ptm.FilterType == poleconst.FilterValueAll && ptm.Filter == "" {
//		return ptm.Polesite.Poles
//	}
//
//	res := []*ps.Pole{}
//	expected := strings.ToUpper(ptm.Filter)
//	filter := func(p *ps.Pole) bool {
//		sis := p.SearchString(ptm.FilterType)
//		if sis == "" {
//			return false
//		}
//		return strings.Contains(strings.ToUpper(sis), expected)
//	}
//	for _, pole := range ptm.Polesite.Poles {
//		if filter(pole) {
//			res = append(res, pole)
//		}
//	}
//	return res
//}
