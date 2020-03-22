package actorupdatemodal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/model/actor/actorconst"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
	"strings"
)

const template string = `<el-dialog
        :before-close="HideWithControl"
        :visible.sync="visible"
        width="70%"
>
    <!-- 
        Modal Title
    -->
    <span slot="title">
		<el-row :gutter="10" align="middle" type="flex">
			<el-col :span="12">
				<h2 style="margin: 0 0" v-if="current_actor">
					<i class="far fa-edit icon--left"></i>Edition de l'acteur : <span style="color: #ccebff">{{current_actor.Ref}}</span>
				</h2>
			</el-col>
		</el-row>
    </span>

    <!-- 
        Modal Body
        style="height: 100%;"
        
    -->
    <div style="height: 45vh; padding: 5px 25px; overflow-x: hidden;overflow-y: auto;">
        <!-- Last & First Name -->
        <el-row :gutter="10" align="middle" class="spaced" type="flex">
            <el-col :span="4" class="align-right">Nom :</el-col>
            <el-col :span="8">
                <el-input @change="CheckName"
                          clearable placeholder="Nom" size="mini"
                          v-model="current_actor.LastName"
                ></el-input>
            </el-col>

            <el-col :span="4" class="align-right">Prénom :</el-col>
            <el-col :span="8">
                <el-input @change="CheckName"
                          clearable placeholder="Prénom" size="mini"
                          v-model="current_actor.FirstName"
                ></el-input>
            </el-col>
        </el-row>

        <!-- Company & Contract Name -->
        <el-row :gutter="10" align="middle" class="doublespaced" type="flex">
            <el-col :span="4" class="align-right">Contrat :</el-col>
            <el-col :span="8">
                <el-select filterable allow-create clearable placeholder="Contrat" size="mini"
                           v-model="current_actor.Contract"
                           @change="CheckName" style="width: 100%"
                >
                    <el-option v-for="item in GetContractList()"
                               :key="item.value"
                               :label="item.label"
                               :value="item.value"
                    >
                    </el-option>
                </el-select>
            </el-col>

            <el-col :span="4" class="align-right">Société :</el-col>
            <el-col :span="8">
                <el-input @change="CheckName"
                          clearable placeholder="Nom" size="mini"
                          v-model="current_actor.Company"
                ></el-input>
            </el-col>
        </el-row>

        <!-- Start / Leave Dates -->
        <el-row :gutter="10" align="middle" class="doublespaced" type="flex">
            <el-col :span="4" class="align-right">Premier jour d'activité :</el-col>
            <el-col :span="8">
                <el-date-picker :picker-options="{firstDayOfWeek:1}" format="dd/MM/yyyy"
                                placeholder="Date" size="mini"
                                style="width: 100%"
                                type="date"
                                v-model="current_actor.Period.Begin"
                                value-format="yyyy-MM-dd"
                ></el-date-picker>
            </el-col>

            <el-col :span="4" class="align-right">Dernier jour d'activité :</el-col>
            <el-col :span="8">
                <el-date-picker :picker-options="{firstDayOfWeek:1}" format="dd/MM/yyyy"
                                placeholder="Date" size="mini"
                                style="width: 100%"
                                type="date"
                                v-model="current_actor.Period.End"
                                value-format="yyyy-MM-dd"
                ></el-date-picker>
            </el-col>
        </el-row>

        <!-- Client & Role -->
        <el-row :gutter="10" align="middle" class="doublespaced" type="flex">
            <el-col :span="4" class="align-right">Client(s) :</el-col>
            <el-col :span="8">
                <el-select multiple placeholder="Client" size="mini"
                           v-model="current_actor.Client"
                           style="width: 100%"
                >
                    <el-option v-for="item in GetClientList()"
                               :key="item.value"
                               :label="item.label"
                               :value="item.value"
                    >
                    </el-option>
                </el-select>
            </el-col>

            <el-col :span="4" class="align-right">Rôle :</el-col>
            <el-col :span="8">
                <el-select filterable allow-create clearable placeholder="Rôle" size="mini"
                           v-model="current_actor.Role"
                           @change="CheckName" style="width: 100%"
                >
                    <el-option v-for="item in GetRoleList()"
                               :key="item.value"
                               :label="item.label"
                               :value="item.value"
                    >
                    </el-option>
                </el-select>
            </el-col>
        </el-row>

        <!-- Comment -->
        <el-row :gutter="10" align="middle" class="doublespaced" type="flex">
            <el-col :span="4" class="align-right">Commentaire :</el-col>
            <el-col :span="20">
                <el-input type="textarea" :autosize="{ minRows: 2, maxRows: 5}" placeholder="Commentaire"
                          v-model="current_actor.Comment" clearable size="mini"
                ></el-input>
            </el-col>
        </el-row>
    </div>

    <!-- 
        Modal Footer Action Bar
    -->
    <span slot="footer">
		<el-row :gutter="15">
			<el-col :span="24" style="text-align: right">
				<el-tooltip :open-delay="500" effect="light">
					<div slot="content">Annuler les changements</div>
					<el-button :disabled="!hasChanged" @click="UndoChange" icon="fas fa-undo-alt" plain size="mini"
                               type="info"></el-button>
				</el-tooltip>
				
				<el-button @click="Hide" size="mini">Fermer</el-button>
				
				<el-button :disabled="!hasChanged" type="success" @click="ConfirmChange" plain size="mini">Enregistrer</el-button>
			</el-col>
		</el-row>
	</span>
</el-dialog>`

type ActorUpdateModalModel struct {
	*ActorModalModel
}

func NewActorUpdateModalModel(vm *hvue.VM) *ActorUpdateModalModel {
	aumm := &ActorUpdateModalModel{ActorModalModel: NewActorModalModel(vm)}
	return aumm
}

func ActorUpdateModalModelFromJS(o *js.Object) *ActorUpdateModalModel {
	return &ActorUpdateModalModel{ActorModalModel: ActorModalModelFromJS(o)}
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Component Methods

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("actor-update-modal", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewActorUpdateModalModel(vm)
		}),
		hvue.MethodsOf(&ActorUpdateModalModel{}),
		hvue.Computed("isNewActor", func(vm *hvue.VM) interface{} {
			aumm := ActorUpdateModalModelFromJS(vm.Object)
			return aumm.CurrentActor.Id == -1
		}),
		hvue.Computed("hasChanged", func(vm *hvue.VM) interface{} {
			aumm := ActorUpdateModalModelFromJS(vm.Object)
			return aumm.HasChanged()
		}),
		//hvue.Computed("currentActorRef", func(vm *hvue.VM) interface{} {
		//	aumm := ActorUpdateModalModelFromJS(vm.Object)
		//	aumm.CurrentActor.Ref = aumm.CurrentActor.LastName + " " + aumm.CurrentActor.FirstName
		//	return aumm.CurrentActor.Ref
		//}),
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Tools Button Methods

func (aumm *ActorUpdateModalModel) CheckName(vm *hvue.VM) {
	aumm = ActorUpdateModalModelFromJS(vm.Object)

	aumm.CurrentActor.LastName = strings.Trim(strings.ToUpper(aumm.CurrentActor.LastName), " \t")
	aumm.CurrentActor.FirstName = strings.Trim(aumm.CurrentActor.FirstName, " \t")
	aumm.CurrentActor.Ref = aumm.CurrentActor.LastName + " " + aumm.CurrentActor.FirstName

	aumm.CurrentActor.Company = strings.Trim(strings.ToUpper(aumm.CurrentActor.Company), " \t")
}

func (aumm *ActorUpdateModalModel) GetContractList() []*elements.ValueLabel {
	return []*elements.ValueLabel{
		elements.NewValueLabel(actorconst.ContractTemp, actorconst.ContractTemp),
		elements.NewValueLabel(actorconst.ContractCDI, actorconst.ContractCDI),
		elements.NewValueLabel(actorconst.ContractCDD, actorconst.ContractCDD),
	}
}

func (aumm *ActorUpdateModalModel) GetRoleList() []*elements.ValueLabel {
	return []*elements.ValueLabel{
		elements.NewValueLabel(actorconst.RolePuller, actorconst.RolePuller),
		elements.NewValueLabel(actorconst.RoleJuncter, actorconst.RoleJuncter),
		elements.NewValueLabel(actorconst.RoleDriver, actorconst.RoleDriver),
		elements.NewValueLabel(actorconst.RoleTeamleader, actorconst.RoleTeamleader),
	}
}

func (aumm *ActorUpdateModalModel) GetClientList(vm *hvue.VM) []*elements.ValueLabel {
	aumm = ActorUpdateModalModelFromJS(vm.Object)
	res := []*elements.ValueLabel{}
	for _, client := range aumm.User.Clients {
		res = append(res, elements.NewValueLabel(client.Name, client.Name))
	}
	return res
}
