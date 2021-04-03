package actorupdatemodal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/model/actor"
	"github.com/lpuig/ewin/doe/website/frontend/model/actor/actorconst"
	"github.com/lpuig/ewin/doe/website/frontend/model/group"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/date"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
	"strconv"
	"strings"
)

const template string = `<el-dialog
        :before-close="HideWithControl"
        :visible.sync="visible" :close-on-click-modal="false"
        width="70%" top="5vh"
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
    <el-tabs type="border-card" tab-position="left" style="height: 75vh">
		<!-- ===================================== Acteur Tab ======================================================= -->
		<el-tab-pane v-if="user.Permissions.HR" label="Acteur" lazy=true style="height: 75vh; padding: 5px 25px; overflow-x: hidden;overflow-y: auto;">
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
	
			<!-- Group -->
			<el-row :gutter="10" align="middle" class="doublespaced" type="flex">
				<el-col :span="4" class="align-right">Affectation :</el-col>
				<el-col :span="20">
					<el-table 
							:data="GroupsControl.Assignments"
							max-height="160" size="mini" border
					>
						<el-table-column label="" width="80">
							<template slot="header" slot-scope="scope">
								<el-button type="success" plain icon="fas fa-users fa-fw" size="mini" @click="AddAssignment()"></el-button>
							</template>
							<template slot-scope="scope">
								<el-button type="danger" plain icon="far fa-trash-alt fa-fw" size="mini" @click="RemoveAssignment(scope.$index)"></el-button>
							</template>
						</el-table-column>
	
						<el-table-column label="Depuis le" width="180">
							<template slot-scope="scope">
								<el-date-picker :picker-options="{firstDayOfWeek:1}" format="dd/MM/yyyy"
											placeholder="Date" size="mini"
											style="width: 100%"
											type="date"
											v-model="scope.row.Date"
											value-format="yyyy-MM-dd"
											@change="UpdateAssignments"
								></el-date-picker>
							</template>
						</el-table-column>
	
						<el-table-column label="Groupe">
							<template slot-scope="scope">
								<el-select placeholder="Groupe" size="mini"
										   v-model="scope.row.GroupId"
										   @change="UpdateAssignments" style="width: 100%"
								>
									<el-option v-for="item in GetAssignmentGroups()"
											   :key="item.value"
											   :label="item.label"
											   :value="item.value"
									>
									</el-option>
								</el-select>
							</template>
						</el-table-column>
					</el-table>
				</el-col>
			</el-row>
	
	        <!-- Salary -->
			<el-row :gutter="10" align="top" class="doublespaced" type="flex">
				<el-col :span="4" class="align-right">Salaire Net :</el-col>
				<el-col :span="20">
					<el-table 
							:data="current_actor.Info.Salary"
							max-height="160" size="mini" border
					>
						<el-table-column label="" width="80">
							<template slot="header" slot-scope="scope">
								<el-button type="success" plain icon="fas fa-euro-sign fa-fw" size="mini" @click="AddSalary()"></el-button>
							</template>
							<template slot-scope="scope">
								<el-button type="danger" plain icon="far fa-trash-alt fa-fw" size="mini" @click="RemoveSalary(scope.$index)"></el-button>
							</template>
						</el-table-column>
	
						<el-table-column label="Date" width="180">
							<template slot-scope="scope">
								<el-date-picker v-model="scope.row.Date"
												type="month" :clearable="false"
												:picker-options="{firstDayOfWeek:1}" value-format="yyyy-MM-dd" format="dd/MM/yyyy" 
												placeholder="Date" size="mini" style="width: 100%"
												@change="CheckSalary()"
								></el-date-picker>
							</template>
						</el-table-column>
	
						<el-table-column label="Montant" width="180">
							<template slot-scope="scope">
								<el-input-number 
										v-model="scope.row.Amount" 
										:min="0" :step="100" :precision="2" size="mini" 
										controls-position="right" style="width: 100%"
								></el-input-number>
							</template>
						</el-table-column>
	
						<el-table-column label="Commentaire">
							<template slot-scope="scope">
								<el-input 
										v-model="scope.row.Comment" placeholder="Commentaire"
										clearable size="mini" style="width: 100%"
								></el-input>
							</template>
						</el-table-column>
					</el-table>
				</el-col>
			</el-row>

			<!-- TravelSubsidy -->
			<el-row :gutter="10" align="top" class="doublespaced" type="flex">
				<el-col :span="4" class="align-right">Déplacements :</el-col>
				<el-col :span="20">
					<el-table 
							:data="current_actor.Info.TravelSubsidy"
							max-height="200" size="mini" border
					>
						<el-table-column label="" width="80">
							<template slot="header" slot-scope="scope">
								<el-button type="success" plain icon="fas fa-euro-sign fa-fw" size="mini" @click="AddTravelSubsidy()"></el-button>
							</template>
							<template slot-scope="scope">
								<el-button type="danger" plain icon="far fa-trash-alt fa-fw" size="mini" @click="RemoveTravelSubsidy(scope.$index)"></el-button>
							</template>
						</el-table-column>
	
						<el-table-column label="Date" width="180">
							<template slot-scope="scope">
								<el-date-picker v-model="scope.row.Date"
												type="month" :clearable="false"
												:picker-options="{firstDayOfWeek:1}" value-format="yyyy-MM-dd" format="dd/MM/yyyy" 
												placeholder="Date" size="mini" style="width: 100%"
												@change="CheckTravelSubsidy()"
								></el-date-picker>
							</template>
						</el-table-column>
	
						<el-table-column label="Montant" width="180">
							<template slot-scope="scope">
								<el-input-number 
										v-model="scope.row.Amount" 
										:min="0" :step="50" :precision="2" size="mini" 
										controls-position="right" style="width: 100%"
								></el-input-number>
							</template>
						</el-table-column>
	
						<el-table-column label="Commentaire">
							<template slot-scope="scope">
								<el-input 
										v-model="scope.row.Comment" placeholder="Commentaire"
										clearable size="mini" style="width: 100%"
								></el-input>
							</template>
						</el-table-column>
					</el-table>
				</el-col>
			</el-row>

			<!-- Comment -->
			<el-row :gutter="10" align="top" class="doublespaced" type="flex">
				<el-col :span="4" class="align-right">Commentaire :</el-col>
				<el-col :span="20">
					<el-input type="textarea" :autosize="{ minRows: 2, maxRows: 5}" placeholder="Commentaire"
							  v-model="current_actor.Comment" clearable size="mini"
					></el-input>
				</el-col>
			</el-row>

		</el-tab-pane>

		<!-- ===================================== Bonus Tab ======================================================= -->
		<el-tab-pane v-if="user.Permissions.HR" label="Bonus" lazy=true style="height: 75vh; padding: 5px 25px; overflow-x: hidden;overflow-y: auto;">
			<!-- Bonuses -->
			<el-row :gutter="10" align="top" class="spaced" type="flex">
				<el-col :span="3" class="align-right">Bonus :</el-col>
				<el-col :span="21">
					<el-table 
							:data="current_actor.Info.EarnedBonuses"
							max-height="260" size="mini" border
					>
						<el-table-column label="" width="70" :resizable="false">
							<template slot="header" slot-scope="scope">
								<el-button type="success" plain icon="fas fa-euro-sign" size="mini" @click="AddBonus()"></el-button>
							</template>
							<template slot-scope="scope">
								<el-button type="danger" plain icon="far fa-trash-alt" size="mini" @click="RemoveBonus(scope.$index)"></el-button>
							</template>
						</el-table-column>
	
						<el-table-column label="Date Paiement" width="160" :resizable="false">
							<template slot-scope="scope">
								<el-date-picker v-model="scope.row.Date"
												type="month" :clearable="false"
												:picker-options="{firstDayOfWeek:1}" value-format="yyyy-MM-dd" format="dd/MM/yyyy" 
												placeholder="Date" size="mini" style="width: 100%"
												@change="CheckBonus(scope.row)"
								></el-date-picker>
							</template>
						</el-table-column>
	
						<el-table-column label="Montant" width="200" :resizable="false">
							<template slot-scope="scope">
								<div class="header-menu-container">
									<el-input-number 
											v-model="scope.row.Amount" 
											:min="0" :step="100" :precision="2" size="mini" 
											controls-position="right"
									></el-input-number>
									<el-button v-if="IsBonusPayable(scope.row.Date)" 
											type="success" plain icon="fas fa-caret-down" size="mini" 
											:disabled="IsBonusPaid(scope.$index)" 
											@click="PayBonus(scope.$index)"
									></el-button>
								</div>
							</template>
						</el-table-column>
	
						<el-table-column label="Commentaire">
							<template slot-scope="scope">
								<el-input 
										v-model="scope.row.Comment" placeholder="Commentaire"
										clearable size="mini" style="width: 100%"
								></el-input>
							</template>
						</el-table-column>
					</el-table>
				</el-col>
			</el-row>
	
			<!-- Paid Bonuses -->
			<el-row :gutter="10" align="top" class="doublespaced" type="flex">
				<el-col :span="3" class="align-right">Bonus payés :</el-col>
				<el-col :span="21">
					<el-table 
							:data="current_actor.Info.PaidBonuses"
							max-height="260" size="mini" border
					>
						<el-table-column label="" width="70">
							<template slot="header" slot-scope="scope">
								<el-button type="success" plain icon="fas fa-euro-sign" size="mini" @click="AddPaidBonus()"></el-button>
							</template>
							<template slot-scope="scope">
								<el-button type="danger" plain icon="far fa-trash-alt" size="mini" @click="RemovePaidBonus(scope.$index)"></el-button>
							</template>
						</el-table-column>
	
						<el-table-column label="Date" width="160">
							<template slot-scope="scope">
								<el-date-picker v-model="scope.row.Date"
												type="month" :clearable="false"
												:picker-options="{firstDayOfWeek:1}" value-format="yyyy-MM-dd" format="dd/MM/yyyy" 
												placeholder="Date" size="mini" style="width: 100%"
												@change="CheckPaidBonus(scope.row)"
								></el-date-picker>
							</template>
						</el-table-column>
	
						<el-table-column label="Montant" width="200">
							<template slot-scope="scope">
								<el-input-number 
										v-model="scope.row.Amount" 
										:min="0" :step="100" :precision="2" size="mini" 
										controls-position="right" style="width: 100%"
								></el-input-number>
							</template>
						</el-table-column>
	
						<el-table-column label="Commentaire">
							<template slot-scope="scope">
								<el-input 
										v-model="scope.row.Comment" placeholder="Commentaire"
										clearable size="mini" style="width: 100%"
								></el-input>
							</template>
						</el-table-column>
					</el-table>
				</el-col>
			</el-row>
		</el-tab-pane>

		<!-- ===================================== Vacancy Tab ======================================================= -->
		<el-tab-pane label="Absences" lazy=true style="height: 75vh; padding: 5px 25px; overflow-x: hidden;overflow-y: auto;">
			<el-table 
					:border="false"
					:data="VacationDates"
					height="100%" size="mini"
	
			>
				<el-table-column
						label="" width="70px" align="center"
				>
					<template slot="header" slot-scope="scope">
						<el-button type="success" plain icon="fas fa-user-plus" circle size="mini" :disabled="!datesComplete" @click="AddDates()"></el-button>
					</template>
					<template slot-scope="scope">
						<el-button type="danger" icon="fas fa-user-minus" circle size="mini" @click="DeleteDates(scope.$index)"></el-button>
					</template>
				</el-table-column>
				
				<el-table-column
						label="Congés" 
				>
					<template slot-scope="scope">
						<el-date-picker
								v-model="scope.row.Dates"
								type="daterange" unlink-panels size="mini" style="width: 100%"
								:picker-options="{firstDayOfWeek:1}" format="dd/MM/yyyy"
								value-format="yyyy-MM-dd"
								range-separator="au"
								start-placeholder="Début"
								end-placeholder="Fin"
								@change="ApplyVacationChange()">
						></el-date-picker>
					</template>
				</el-table-column>
	
				<el-table-column label="Commentaire">
					<template slot-scope="scope">
						<el-input
								type="textarea" placeholder="Nature du congé" size="mini"
								v-model="scope.row.Comment" @input="ApplyVacationChange()">
						</el-input>
					</template>
				</el-table-column>
			</el-table>
		</el-tab-pane>

		<!-- ===================================== Training Tab ======================================================= -->
		<el-tab-pane label="Formations" lazy=true style="height: 75vh; padding: 5px 25px; overflow-x: hidden;overflow-y: auto;">
			<el-row :gutter="10" align="top" class="spaced" type="flex">
				<el-col :span="3" class="align-right">Visites et Formations :</el-col>
				<el-col :span="21">
					<el-row v-for="trainingName in GetTrainingNames()"
							:gutter="10" align="middle"class="spaced" type="flex"
					>
						<el-col :span="3" class="align-right">{{trainingName}} :</el-col>
						<el-col v-if="current_actor.Info.Trainings[trainingName] == null" :span="8">
							<el-button type="success" plain icon="fas fa-user-graduate" size="mini" @click="AddTraining(trainingName)"></el-button>
						</el-col>
						<el-col v-else :span="21">
							<el-row :gutter="10" align="middle" type="flex">
								<el-col :span="2">
									<el-button type="danger" plain icon="far fa-trash-alt" size="mini" @click="RemoveTraining(trainingName)"></el-button>
								</el-col>
								<el-col :span="8">
									<el-date-picker v-model="current_actor.Info.Trainings[trainingName].Date"
													type="date"
													:picker-options="{firstDayOfWeek:1}" value-format="yyyy-MM-dd" format="dd/MM/yyyy" 
													placeholder="Date" size="mini" style="width: 100%"
													@change="CheckTrainingDate(trainingName)"
									></el-date-picker>
								</el-col>
								<el-col :span="14">
										<el-input 
												v-model="current_actor.Info.Trainings[trainingName].Comment" placeholder="Commentaire"
												clearable size="mini" style="width: 100%"
										></el-input>
								</el-col>
							</el-row>
						</el-col>
					</el-row>
				</el-col>
			</el-row>
		</el-tab-pane>
    </el-tabs>

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
				
				<el-button @click="Hide" size="mini">Ignorer</el-button>
				
				<el-button :disabled="!hasChanged" type="success" @click="ConfirmChange" plain size="mini">Valider</el-button>
			</el-col>
		</el-row>
	</span>
</el-dialog>`

type ActorUpdateModalModel struct {
	*ActorModalModel

	ThisMonth     string            `js:"ThisMonth"`
	GroupStore    *group.GroupStore `js:"groups"`
	GroupsControl *GroupsControl    `js:"GroupsControl"`
	VacationDates []*VacDates       `js:"VacationDates"`
}

func NewActorUpdateModalModel(vm *hvue.VM) *ActorUpdateModalModel {
	aumm := &ActorUpdateModalModel{ActorModalModel: NewActorModalModel(vm)}
	aumm.ThisMonth = date.GetFirstOfMonth(date.TodayAfter(0))
	aumm.GroupStore = nil
	aumm.GroupsControl = NewGroupsControl(nil)
	aumm.VacationDates = []*VacDates{}
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
		hvue.Props("groups"),
		hvue.Computed("isNewActor", func(vm *hvue.VM) interface{} {
			aumm := ActorUpdateModalModelFromJS(vm.Object)
			return aumm.CurrentActor.Id == -1
		}),
		hvue.Computed("hasChanged", func(vm *hvue.VM) interface{} {
			aumm := ActorUpdateModalModelFromJS(vm.Object)
			return aumm.HasChanged()
		}),
		hvue.Computed("datesComplete", func(vm *hvue.VM) interface{} {
			aumm := ActorUpdateModalModelFromJS(vm.Object)
			return aumm.DatesComplete()
		}),
		//hvue.Computed("currentActorRef", func(vm *hvue.VM) interface{} {
		//	aumm := ActorUpdateModalModelFromJS(vm.Object)
		//	aumm.CurrentActor.Ref = aumm.CurrentActor.LastName + " " + aumm.CurrentActor.FirstName
		//	return aumm.CurrentActor.Ref
		//}),
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Modal Methods

func (aumm *ActorUpdateModalModel) Show(act *actor.Actor, user *fm.User) {
	aumm.GroupsControl = NewGroupsControl(aumm.GroupStore)
	aumm.EditedActor = act
	aumm.CurrentActor = act.Copy()
	aumm.GroupsControl.SetAssignments(aumm.CurrentActor)
	aumm.User = user
	aumm.UpdateDates()
	aumm.ShowConfirmDelete = false
	aumm.Visible = true
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Tools Button Methods

func (aumm *ActorUpdateModalModel) CheckName(vm *hvue.VM) {
	aumm = ActorUpdateModalModelFromJS(vm.Object)

	aumm.CurrentActor.LastName = strings.Trim(strings.ToUpper(aumm.CurrentActor.LastName), " \t")
	aumm.CurrentActor.FirstName = strings.Trim(strings.Title(aumm.CurrentActor.FirstName), " \t")
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

func (aumm *ActorUpdateModalModel) ConfirmChange(vm *hvue.VM) {
	aumm = ActorUpdateModalModelFromJS(vm.Object)
	aumm.UpdateVacation()
	aumm.ActorModalModel.ConfirmChange()
	vm.Emit("edited-actor", aumm.EditedActor)
}

func (aumm *ActorUpdateModalModel) UndoChange() {
	aumm.ActorModalModel.UndoChange()
	aumm.UpdateDates()
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Actor Group Tools Methods

func (aumm *ActorUpdateModalModel) GetCurrentGroupInfo(vm *hvue.VM) string {
	aumm = ActorUpdateModalModelFromJS(vm.Object)
	id, assignDate := aumm.CurrentActor.Groups.GetCurrentInfo()
	if id < 0 {
		return "Pas de groupe affecté"
	}
	group := aumm.GroupStore.GetGroupById(id)
	if group == nil {
		return "Pas de groupe connu affecté"
	}
	return group.Name + " depuis le " + date.DateString(assignDate)
}

func (aumm *ActorUpdateModalModel) AddAssignment(vm *hvue.VM) {
	aumm = ActorUpdateModalModelFromJS(vm.Object)
	aumm.GroupsControl.Add()
}

func (aumm *ActorUpdateModalModel) RemoveAssignment(vm *hvue.VM, index int) {
	aumm = ActorUpdateModalModelFromJS(vm.Object)
	aumm.GroupsControl.Remove(index)
}

func (aumm *ActorUpdateModalModel) UpdateAssignments(vm *hvue.VM) {
	aumm = ActorUpdateModalModelFromJS(vm.Object)
	aumm.GroupsControl.SortAssignments()
	aumm.GroupsControl.UpdateActor()
}

func (aumm *ActorUpdateModalModel) GetAssignmentGroups(vm *hvue.VM) []*elements.ValueLabel {
	aumm = ActorUpdateModalModelFromJS(vm.Object)
	res := make([]*elements.ValueLabel, len(aumm.GroupStore.Groups))
	for i, grp := range aumm.GroupStore.GetGroupsSortedByName() {
		res[i] = elements.NewValueLabel(strconv.Itoa(grp.Id), grp.Name)
	}
	return res
}

//////////////////////////////////////////////////////////////////////////////////////////////
// ActorInfos Tools Button Methods

// Salary

func (aumm *ActorUpdateModalModel) AddSalary(vm *hvue.VM) {
	aumm = ActorUpdateModalModelFromJS(vm.Object)
	nd := actor.NewDateAmountComment()
	if len(aumm.CurrentActor.Info.Salary) == 0 {
		nd.Date = aumm.CurrentActor.Period.Begin
		nd.Amount = 1500
	} else {
		nd.Date = date.GetFirstOfMonth(date.TodayAfter(0))
		nd.Amount = aumm.CurrentActor.Info.Salary[0].Amount + 200
	}
	aumm.CurrentActor.Info.Salary = append(actor.EarningHistory{*nd}, aumm.CurrentActor.Info.Salary...)
	//aumm.CurrentActor.Info.Get("Salary").Call("unshift", nd)
}

func (aumm *ActorUpdateModalModel) RemoveSalary(vm *hvue.VM, pos int) {
	aumm = ActorUpdateModalModelFromJS(vm.Object)
	aumm.CurrentActor.Info.Get("Salary").Call("splice", pos, 1)
}

func (aumm *ActorUpdateModalModel) CheckSalary(vm *hvue.VM) {
	aumm = ActorUpdateModalModelFromJS(vm.Object)
	aumm.CurrentActor.Info.Get("Salary").Call("sort", actor.CompareDateAmountComment)
}

// TravelSubsidy

func (aumm *ActorUpdateModalModel) AddTravelSubsidy(vm *hvue.VM) {
	aumm = ActorUpdateModalModelFromJS(vm.Object)
	nd := actor.NewDateAmountComment()
	nd.Date = date.GetFirstOfMonth(date.TodayAfter(0))
	aumm.CurrentActor.Info.TravelSubsidy = append(actor.EarningHistory{*nd}, aumm.CurrentActor.Info.TravelSubsidy...)
}

func (aumm *ActorUpdateModalModel) RemoveTravelSubsidy(vm *hvue.VM, pos int) {
	aumm = ActorUpdateModalModelFromJS(vm.Object)
	aumm.CurrentActor.Info.Get("TravelSubsidy").Call("splice", pos, 1)
}

func (aumm *ActorUpdateModalModel) CheckTravelSubsidy(vm *hvue.VM) {
	aumm = ActorUpdateModalModelFromJS(vm.Object)
	aumm.CurrentActor.Info.Get("TravelSubsidy").Call("sort", actor.CompareDateAmountComment)
}

// EarnedBonuses

func (aumm *ActorUpdateModalModel) AddBonus(vm *hvue.VM) {
	aumm = ActorUpdateModalModelFromJS(vm.Object)
	nd := actor.NewDateAmountComment()
	nd.Date = date.GetFirstOfMonth(date.TodayAfter(0))
	aumm.CurrentActor.Info.EarnedBonuses = append(actor.Earnings{*nd}, aumm.CurrentActor.Info.EarnedBonuses...)
}

func (aumm *ActorUpdateModalModel) RemoveBonus(vm *hvue.VM, pos int) {
	aumm = ActorUpdateModalModelFromJS(vm.Object)
	aumm.CurrentActor.Info.Get("EarnedBonuses").Call("splice", pos, 1)
}

func (aumm *ActorUpdateModalModel) CheckBonus(vm *hvue.VM, value *js.Object) {
	aumm = ActorUpdateModalModelFromJS(vm.Object)
	currentDate := value.Get("Date").String()
	if !strings.HasSuffix(currentDate, "01") {
		value.Set("Date", date.GetFirstOfMonth(currentDate))
	}
	aumm.CurrentActor.Info.Get("EarnedBonuses").Call("sort", actor.CompareDateAmountComment)
}

func (aumm *ActorUpdateModalModel) IsBonusPaid(vm *hvue.VM, index int) bool {
	aumm = ActorUpdateModalModelFromJS(vm.Object)
	bonus := aumm.CurrentActor.Info.EarnedBonuses[index]
	nbMatch := 0
	if index < len(aumm.CurrentActor.Info.EarnedBonuses)-1 {
		for _, earnedBonus := range aumm.CurrentActor.Info.EarnedBonuses[index+1:] {
			if earnedBonus.Date != bonus.Date {
				break
			}
			diffAmount := earnedBonus.Amount - bonus.Amount
			if diffAmount < 0 {
				diffAmount = -diffAmount
			}
			if diffAmount < 0.001 {
				nbMatch++
			}
		}
	}
	nbPaidMatch := 0
	for _, paidBonus := range aumm.CurrentActor.Info.PaidBonuses {
		if paidBonus.Date < bonus.Date {
			break
		}
		if paidBonus.Date != bonus.Date {
			continue
		}
		diffAmount := paidBonus.Amount - bonus.Amount
		if diffAmount < 0 {
			diffAmount = -diffAmount
		}
		if diffAmount < 0.001 {
			nbPaidMatch++
		}
	}
	return nbPaidMatch > nbMatch
}

func (aumm *ActorUpdateModalModel) IsBonusPayable(vm *hvue.VM, month string) bool {
	aumm = ActorUpdateModalModelFromJS(vm.Object)
	if date.GetFirstOfMonth(month) >= aumm.ThisMonth {
		return false
	}
	return true
}

func (aumm *ActorUpdateModalModel) PayBonus(vm *hvue.VM, index int) {
	aumm = ActorUpdateModalModelFromJS(vm.Object)
	bonusToPay := aumm.CurrentActor.Info.EarnedBonuses[index].Copy()
	aumm.CurrentActor.Info.PaidBonuses = append(actor.Earnings{*bonusToPay}, aumm.CurrentActor.Info.PaidBonuses...)
	aumm.CurrentActor.Info.Get("PaidBonuses").Call("sort", actor.CompareDateAmountComment)
}

// PaidBonuses

func (aumm *ActorUpdateModalModel) AddPaidBonus(vm *hvue.VM) {
	aumm = ActorUpdateModalModelFromJS(vm.Object)
	nd := actor.NewDateAmountComment()
	nd.Date = date.GetFirstOfMonth(date.TodayAfter(0))
	aumm.CurrentActor.Info.PaidBonuses = append(actor.Earnings{*nd}, aumm.CurrentActor.Info.PaidBonuses...)
}

func (aumm *ActorUpdateModalModel) RemovePaidBonus(vm *hvue.VM, pos int) {
	aumm = ActorUpdateModalModelFromJS(vm.Object)
	aumm.CurrentActor.Info.Get("PaidBonuses").Call("splice", pos, 1)
}

func (aumm *ActorUpdateModalModel) CheckPaidBonus(vm *hvue.VM, value *js.Object) {
	aumm = ActorUpdateModalModelFromJS(vm.Object)
	currentDate := value.Get("Date").String()
	if !strings.HasSuffix(currentDate, "01") {
		value.Set("Date", date.GetFirstOfMonth(currentDate))
	}
	aumm.CurrentActor.Info.Get("PaidBonuses").Call("sort", actor.CompareDateAmountComment)
}

// Trainings

func (aumm *ActorUpdateModalModel) GetTrainingNames() []string {
	return actor.GetDefaultInfoTraining()
}

func (aumm *ActorUpdateModalModel) AddTraining(vm *hvue.VM, name string) {
	aumm = ActorUpdateModalModelFromJS(vm.Object)
	dc := *actor.NewDateComment()
	dc.Date = date.TodayAfter(0)
	trainings := aumm.CurrentActor.Info.Trainings
	trainings[name] = dc
	aumm.CurrentActor.Info.Trainings = trainings
}

func (aumm *ActorUpdateModalModel) RemoveTraining(vm *hvue.VM, name string) {
	aumm = ActorUpdateModalModelFromJS(vm.Object)
	if _, found := aumm.CurrentActor.Info.Trainings[name]; !found {
		return
	}
	trainings := aumm.CurrentActor.Info.Trainings
	delete(trainings, name)
	aumm.CurrentActor.Info.Trainings = trainings
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Vacancy Tab Methods

func (aumm *ActorUpdateModalModel) DatesComplete() bool {
	if len(aumm.VacationDates) == 0 {
		return true
	}
	return aumm.VacationDates[0].IsComplete()
}

func (aumm *ActorUpdateModalModel) AddDates(vm *hvue.VM) {
	aumm = ActorUpdateModalModelFromJS(vm.Object)
	aumm.VacationDates = append([]*VacDates{NewVacDates("", "", "")}, aumm.VacationDates...)
	aumm.UpdateVacation()
}

func (aumm *ActorUpdateModalModel) DeleteDates(vm *hvue.VM, index int) {
	aumm = ActorUpdateModalModelFromJS(vm.Object)
	aumm.Get("VacationDates").Call("splice", index, 1)
	aumm.UpdateVacation()
}

func (aumm *ActorUpdateModalModel) ApplyVacationChange(vm *hvue.VM) {
	aumm = ActorUpdateModalModelFromJS(vm.Object)
	aumm.SortDates()
	aumm.UpdateVacation()
}

func (aumm *ActorUpdateModalModel) SortDates() {
	aumm.Get("VacationDates").Call("sort", func(a, b *VacDates) int {
		switch {
		case a.Dates[0] == b.Dates[0]:
			return 0
		case a.Dates[0] < b.Dates[0]:
			return 1
		default:
			return -1
		}
	})
}

// UpdateDates copies currentactors Vacation' Dates to modal VacationDates
func (aumm *ActorUpdateModalModel) UpdateDates() {
	aumm.VacationDates = []*VacDates{}
	for _, vacPeriod := range aumm.CurrentActor.Vacation {
		aumm.VacationDates = append(aumm.VacationDates, NewVacDates(vacPeriod.Begin, vacPeriod.End, vacPeriod.Comment))
	}
	aumm.SortDates()
}

// UpdateVacation copies modal VacationDates to currentactors Vacation' Dates
func (aumm *ActorUpdateModalModel) UpdateVacation() {
	aumm.CurrentActor.Vacation = []*date.DateRangeComment{}
	for _, vacDates := range aumm.VacationDates {
		aumm.CurrentActor.Vacation = append(aumm.CurrentActor.Vacation, date.NewDateRangeCommentFrom(vacDates.Dates[0], vacDates.Dates[1], vacDates.Comment))
	}
}

// VacDates type
type VacDates struct {
	*js.Object
	Dates   []string `js:"Dates"`
	Comment string   `js:"Comment"`
}

func NewVacDates(beg, end, cmt string) *VacDates {
	vd := &VacDates{Object: tools.O()}
	vd.Dates = []string{beg, end}
	vd.Comment = cmt
	return vd
}

func (vd *VacDates) IsComplete() bool {
	return !tools.Empty(vd.Dates[0]) && !tools.Empty(vd.Dates[1])
}
