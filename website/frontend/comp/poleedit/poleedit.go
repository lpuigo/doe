package poleedit

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
	"github.com/lpuig/ewin/mappoc/frontend/mappoc/comp/polemap"
	"github.com/lpuig/ewin/mappoc/frontend/mappoc/model"
)

const template string = `<div>
    <h1>
        Poteau: {{polemarker.Pole.Ref}}
    </h1>
    <el-row :gutter="5" type="flex" align="middle" class="spaced">
        <el-col :span="6">Référence:</el-col>
        <el-col :span="18">
            <el-input placeholder="Référence"
                      v-model="polemarker.Pole.Ref" clearable size="mini"
					  @change="UpdateTooltip()"
            ></el-input>
        </el-col>
    </el-row>
    <el-row :gutter="5" type="flex" align="middle" class="spaced">
        <el-col :span="6">Lat / Long:</el-col>
        <el-col :span="9">
            <el-input-number v-model="polemarker.Pole.Lat" size="mini" :precision="8" :controls="false" style="width: 100%"
            ></el-input-number>
        </el-col>
        <el-col :span="9">
            <el-input-number v-model="polemarker.Pole.Long" size="mini" :precision="8" :controls="false" style="width: 100%"
            ></el-input-number>
        </el-col>
    </el-row>
    <el-row :gutter="5" type="flex" align="middle" class="spaced">
        <el-col :span="6">Ville:</el-col>
        <el-col :span="18">
            <el-input placeholder="Ville"
                      v-model="polemarker.Pole.City" clearable size="mini"
            ></el-input>
        </el-col>
    </el-row>
    <el-row :gutter="5" type="flex" align="middle" class="spaced">
        <el-col :span="6">Status:</el-col>
        <el-col :span="18">
            <el-select v-model="polemarker.Pole.State" filterable size="mini" style="width: 100%"
                       @clear=""
                       @change="UpdateState()"
            >
                <el-option
                        v-for="item in GetStates()"
                        :key="item.value"
                        :label="item.label"
                        :value="item.value"
                ></el-option>
            </el-select>
        </el-col>
    </el-row>
</div>
`

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Registration

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("pole-edit", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.Props("polemarker"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewPoleEditModel(vm)
		}),
		hvue.MethodsOf(&PoleEditModel{}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type PoleEditModel struct {
	*js.Object
	PoleMarker *polemap.PoleMarker `js:"polemarker"`

	VM *hvue.VM `js:"VM"`
}

func PoleEditModelFromJS(obj *js.Object) *PoleEditModel {
	return &PoleEditModel{Object: obj}
}

func NewPoleEditModel(vm *hvue.VM) *PoleEditModel {
	pem := &PoleEditModel{Object: tools.O()}
	pem.VM = vm
	pem.PoleMarker = nil
	return pem
}

func (pem *PoleEditModel) GetStates() []*elements.ValueLabel {
	return model.GetStatesValueLabel()
}

func (pem *PoleEditModel) UpdateState(vm *hvue.VM) {
	pem = PoleEditModelFromJS(vm.Object)
	pem.PoleMarker.UpdateFromState()
	pem.PoleMarker.Refresh()
}

func (pem *PoleEditModel) UpdateTooltip(vm *hvue.VM) {
	pem = PoleEditModelFromJS(vm.Object)
	pem.PoleMarker.UpdateTitle()
	pem.PoleMarker.Refresh()
}
