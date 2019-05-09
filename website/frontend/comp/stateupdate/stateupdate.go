package stateupdate

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
		<el-col :span="3">
			<el-select v-model="value.Status" filterable
					   size="mini" style="width: 100%"
					   placeholder="Etat"
					   @clear=""
					   @change=""
			>
				<el-option
						v-for="item in GetStatuses()"
						:key="item.value"
						:label="item.label"
						:value="item.value">
				</el-option>
			</el-select>
		</el-col>
		<el-col :span="5">
			<el-select v-model="value.Team" clearable filterable
					   size="mini" style="width: 100%"
					   placeholder="Equipe"
					   @clear=""
					   @change="SetDates(scope.row)"
			>
				<el-option
						v-for="item in GetTeams()"
						:key="item.value"
						:label="item.label"
						:value="item.value">
				</el-option>
			</el-select>
		</el-col>
		<el-col :span="5">
			<el-date-picker format="dd/MM/yyyy" placeholder="Début" size="mini"
							style="width: 100%" type="date"
							v-model="value.DateStart"
							value-format="yyyy-MM-dd"
							:picker-options="{firstDayOfWeek:1, disabledDate(time) { return time.getTime() > Date.now(); }}"
							:disabled="false" :clearable="false"
			></el-date-picker>
		</el-col>
		<el-col :span="5">
			<el-date-picker format="dd/MM/yyyy" placeholder="Fin" size="mini"
							style="width: 100%" type="date"
							v-model="value.DateEnd"
							value-format="yyyy-MM-dd"
							:picker-options="{firstDayOfWeek:1, disabledDate(time) { return time.getTime() > Date.now(); }}"
							:disabled="false" :clearable="false"
			></el-date-picker>
		</el-col>
		<el-col :span="6">
			<el-input type="textarea" autosize placeholder="Commentaire" size="mini"
					  v-model="value.Comment"
			></el-input>
		</el-col>
	</el-row>
`

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Registration Ripsite version

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("state-update", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.Props("value", "user", "client"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewStateUpdateModel(vm)
		}),
		hvue.MethodsOf(&StateUpdateModel{}),
		//hvue.Computed("NbAvailPulling", func(vm *hvue.VM) interface{} {
		//	rim := &StateUpdateModel{Object: vm.Object}
		//	return NbAvailPulling
		//}),
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
