package ripstateupdate

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	fmrip "github.com/lpuig/ewin/doe/website/frontend/model/ripsite"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
)

const template string = `
	<el-row type="flex" align="middle" :gutter="10">
		<el-col :span="5">
			<el-select v-model="value.Team" clearable filterable
					   size="mini" style="width: 100%"
					   placeholder="Equipe"
					   @clear="UpdateStatus()"
					   @change="UpdateStatus()"
                       :disabled="DisableTeam()"
			>
				<el-option
						v-for="item in GetTeams()"
						:key="item.value"
						:label="item.label"
						:value="item.value">
				</el-option>
			</el-select>
		</el-col>
		<el-col :span="3">
			<el-date-picker format="dd/MM/yyyy" placeholder="Début" size="mini"
							style="width: 100%" type="date"
							v-model="value.DateStart"
							value-format="yyyy-MM-dd"
							:picker-options="{firstDayOfWeek:1, disabledDate(time) { return time.getTime() > Date.now(); }}"
							:disabled="DisableDates()" :clearable="true"
                            @change="UpdateStatus()"
			></el-date-picker>
		</el-col>
		<el-col :span="3">
			<el-date-picker format="dd/MM/yyyy" placeholder="Fin" size="mini"
							style="width: 100%" type="date"
							v-model="value.DateEnd"
							value-format="yyyy-MM-dd"
							:picker-options="{firstDayOfWeek:1, disabledDate(time) { return time.getTime() > Date.now(); }}"
							:disabled="DisableDates()" :clearable="true"
                            @change="UpdateStatus()"
			></el-date-picker>
		</el-col>
        <el-col :span="3">
            <el-select v-model="value.Status" filterable
                       size="mini" style="width: 100%"
                       placeholder="Etat"
                       @clear=""
                       @change="UpdateStatus()"
            >
                <el-option
                        v-for="item in GetStatuses()"
                        :key="item.value"
                        :label="item.label"
                        :value="item.value">
                </el-option>
            </el-select>
        </el-col>
		<el-col :span="10">
			<el-input type="textarea" autosize placeholder="Commentaire" size="mini"
					  v-model="value.Comment"
			></el-input>
		</el-col>
	</el-row>
`

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Registration Ripsite version

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("rip-state-update", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.Props("value", "user", "client"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewStateUpdateModel(vm)
		}),
		hvue.MethodsOf(&StateUpdateModel{}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type StateUpdateModel struct {
	*js.Object

	State  *fmrip.State `js:"value"`
	User   *fm.User     `js:"user"`
	Client string       `js:"client"`

	VM *hvue.VM `js:"VM"`
}

func NewStateUpdateModel(vm *hvue.VM) *StateUpdateModel {
	sum := &StateUpdateModel{Object: tools.O()}
	sum.VM = vm
	sum.State = fmrip.NewState()
	sum.User = fm.NewUser()
	sum.Client = ""

	return sum
}

func StateUpdateModelFromJS(o *js.Object) *StateUpdateModel {
	return &StateUpdateModel{Object: o}
}

func (sum *StateUpdateModel) GetTeams(vm *hvue.VM) []*elements.ValueLabel {
	sum = StateUpdateModelFromJS(vm.Object)
	return sum.User.GetTeamValueLabelsFor(sum.Client)
}

func (sum *StateUpdateModel) GetStatuses() []*elements.ValueLabel {
	return fmrip.GetStateStatusesValueLabel()
}

func (sum *StateUpdateModel) GetAllStatuses() []*elements.ValueLabel {
	return fmrip.GetStateStatusesWithWarningValueLabel()
}

func (sum *StateUpdateModel) UpdateStatus(vm *hvue.VM) {
	sum = StateUpdateModelFromJS(vm.Object)
	sum.State.UpdateStatus()
}

func (sum *StateUpdateModel) DisableDates(vm *hvue.VM) bool {
	sum = StateUpdateModelFromJS(vm.Object)
	if !sum.State.IsCanceled() && !tools.Empty(sum.State.Team) {
		return false
	}
	return true
}

func (sum *StateUpdateModel) DisableTeam(vm *hvue.VM) bool {
	sum = StateUpdateModelFromJS(vm.Object)
	return sum.State.IsCanceled()
}
