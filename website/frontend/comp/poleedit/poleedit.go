package poleedit

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/polemap"
	"github.com/lpuig/ewin/doe/website/frontend/model/polesite"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
)

const template string = `<div>
    <h1>
        Poteau: {{editedpolemarker.Pole.Ref}}
    </h1>
    <el-row :gutter="5" type="flex" align="middle" class="spaced">
        <el-col :span="6">Référence:</el-col>
        <el-col :span="18">
            <el-input placeholder="Référence"
                      v-model="editedpolemarker.Pole.Ref" clearable size="mini"
					  @change="UpdateTooltip()"
            ></el-input>
        </el-col>
    </el-row>
    <el-row :gutter="5" type="flex" align="middle" class="spaced">
        <el-col :span="6">Lat / Long:</el-col>
        <el-col :span="9">
            <el-input-number v-model="editedpolemarker.Pole.Lat" size="mini" :precision="8" :controls="false" style="width: 100%"
            ></el-input-number>
        </el-col>
        <el-col :span="9">
            <el-input-number v-model="editedpolemarker.Pole.Long" size="mini" :precision="8" :controls="false" style="width: 100%"
            ></el-input-number>
        </el-col>
    </el-row>
    <el-row :gutter="5" type="flex" align="middle" class="spaced">
        <el-col :span="6">Ville:</el-col>
        <el-col :span="18">
            <el-input placeholder="Ville"
                      v-model="editedpolemarker.Pole.City" clearable size="mini"
            ></el-input>
        </el-col>
    </el-row>
    <el-row :gutter="5" type="flex" align="middle" class="spaced">
        <el-col :span="6">Status:</el-col>
        <el-col :span="18">
            <el-select v-model="editedpolemarker.Pole.State" filterable size="mini" style="width: 100%"
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
		hvue.Props("editedpolemarker"),
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
	EditedPoleMarker *polemap.PoleMarker `js:"editedpolemarker"`

	VM *hvue.VM `js:"VM"`
}

func PoleEditModelFromJS(obj *js.Object) *PoleEditModel {
	return &PoleEditModel{Object: obj}
}

func NewPoleEditModel(vm *hvue.VM) *PoleEditModel {
	pem := &PoleEditModel{Object: tools.O()}
	pem.VM = vm
	pem.EditedPoleMarker = nil
	return pem
}

func (pem *PoleEditModel) GetStates() []*elements.ValueLabel {
	return polesite.GetStatesValueLabel()
}

func (pem *PoleEditModel) UpdateState(vm *hvue.VM) {
	pem = PoleEditModelFromJS(vm.Object)
	pem.EditedPoleMarker.UpdateFromState()
	pem.EditedPoleMarker.Refresh()
}

func (pem *PoleEditModel) UpdateTooltip(vm *hvue.VM) {
	pem = PoleEditModelFromJS(vm.Object)
	pem.EditedPoleMarker.UpdateTitle()
	pem.EditedPoleMarker.Refresh()
}
