package actorupdatemodal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/model/actor"
	"github.com/lpuig/ewin/doe/website/frontend/model/actor/actorconst"
	"github.com/lpuig/ewin/doe/website/frontend/tools/date"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
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
        <!-- Last & First Name -->
		<el-tab-pane label="Acteur" lazy=true style="height: 75vh; padding: 5px 25px; overflow-x: hidden;overflow-y: auto;">
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
								<el-button type="success" plain icon="fas fa-euro-sign" size="mini" @click="AddSalary()"></el-button>
							</template>
							<template slot-scope="scope">
								<el-button type="danger" plain icon="far fa-trash-alt" size="mini" @click="RemoveSalary(scope.$index)"></el-button>
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

		<el-tab-pane label="Bonus" lazy=true style="height: 75vh; padding: 5px 25px; overflow-x: hidden;overflow-y: auto;">
			<!-- Bonuses -->
			<el-row :gutter="10" align="top" class="spaced" type="flex">
				<el-col :span="3" class="align-right">Bonus :</el-col>
				<el-col :span="21">
					<el-table 
							:data="current_actor.Info.EarnedBonuses"
							max-height="160" size="mini" border
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
							max-height="160" size="mini" border
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

		<el-tab-pane label="Déplacements" lazy=true style="height: 75vh; padding: 5px 25px; overflow-x: hidden;overflow-y: auto;">
			<!-- TravelSubsidy -->
			<el-row :gutter="10" align="top" class="spaced" type="flex">
				<el-col :span="3" class="align-right">Déplacements :</el-col>
				<el-col :span="21">
					<el-table 
							:data="current_actor.Info.TravelSubsidy"
							max-height="200" size="mini" border
					>
						<el-table-column label="" width="80">
							<template slot="header" slot-scope="scope">
								<el-button type="success" plain icon="fas fa-euro-sign" size="mini" @click="AddTravelSubsidy()"></el-button>
							</template>
							<template slot-scope="scope">
								<el-button type="danger" plain icon="far fa-trash-alt" size="mini" @click="RemoveTravelSubsidy(scope.$index)"></el-button>
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
		</el-tab-pane>

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
				
				<el-button @click="Hide" size="mini">Fermer</el-button>
				
				<el-button :disabled="!hasChanged" type="success" @click="ConfirmChange" plain size="mini">Enregistrer</el-button>
			</el-col>
		</el-row>
	</span>
</el-dialog>`

type ActorUpdateModalModel struct {
	*ActorModalModel

	ThisMonth string `js:"ThisMonth"`
}

func NewActorUpdateModalModel(vm *hvue.VM) *ActorUpdateModalModel {
	aumm := &ActorUpdateModalModel{ActorModalModel: NewActorModalModel(vm)}
	aumm.ThisMonth = date.GetFirstOfMonth(date.TodayAfter(0))
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

func (aumm *ActorUpdateModalModel) GetClientList(vm *hvue.VM) []*elements.ValueLabel {
	aumm = ActorUpdateModalModelFromJS(vm.Object)
	res := []*elements.ValueLabel{}
	for _, client := range aumm.User.Clients {
		res = append(res, elements.NewValueLabel(client.Name, client.Name))
	}
	return res
}

func (aumm *ActorUpdateModalModel) ConfirmChange(vm *hvue.VM) {
	aumm = ActorUpdateModalModelFromJS(vm.Object)
	aumm.ActorModalModel.ConfirmChange()
	vm.Emit("edited-actor", aumm.EditedActor)
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

// TravelSubsidy

func (aumm *ActorUpdateModalModel) AddTravelSubsidy(vm *hvue.VM) {
	aumm = ActorUpdateModalModelFromJS(vm.Object)
	nd := actor.NewDateAmountComment()
	nd.Date = date.GetFirstOfMonth(date.TodayAfter(0))
	aumm.CurrentActor.Info.TravelSubsidy = append(actor.Earnings{*nd}, aumm.CurrentActor.Info.TravelSubsidy...)
}

func (aumm *ActorUpdateModalModel) RemoveTravelSubsidy(vm *hvue.VM, pos int) {
	aumm = ActorUpdateModalModelFromJS(vm.Object)
	aumm.CurrentActor.Info.Get("TravelSubsidy").Call("splice", pos, 1)
}

func (aumm *ActorUpdateModalModel) CheckTravelSubsidy(vm *hvue.VM) {
	aumm = ActorUpdateModalModelFromJS(vm.Object)
	aumm.CurrentActor.Info.Get("TravelSubsidy").Call("sort", actor.CompareDateAmountComment)
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
