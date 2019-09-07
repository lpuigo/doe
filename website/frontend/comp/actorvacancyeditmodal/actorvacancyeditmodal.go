package actorvacancyeditmodal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/actorupdatemodal"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/model/actor"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/date"
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
					<i class="fas fa-user-clock icon--left"></i>Edition des congés de : <span style="color: #ccebff">{{current_actor.Ref}}</span>
				</h2>
			</el-col>
		</el-row>
    </span>

    <!-- 
        Modal Body
        style="height: 100%;"
        
    -->
    <div style="height: 45vh; padding: 5px 25px; ">
		<el-table 
				:border=true
        		:data="VacationDates"
		        max-height="100%" size="mini"

		>
			<el-table-column
					label="Action" width="80px"
			>
				<template slot-scope="scope">
					<el-button type="danger" icon="el-icon-delete" circle size="small" @click="DeleteDates(scope.$index)"></el-button>
				</template>
			</el-table-column>
			
			<el-table-column
					prop="Begin" label="Début - Fin" width="700px" 
			>
				<template slot-scope="scope">
					<el-date-picker
							v-model="scope.row.Dates"
							type="daterange" unlink-panels
							:picker-options="{firstDayOfWeek:1}" format="dd/MM/yyyy"
							value-format="yyyy-MM-dd"
							range-separator="au"
							start-placeholder="Début"
							end-placeholder="Fin"
							@change="ApplyChange()">
					></el-date-picker>

				</template>
			</el-table-column>
		</el-table>

		<pre>{{VacationDates}}</pre>

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

type VacDates struct {
	*js.Object
	Dates []string `js:"Dates"`
}

func NewVacDates(beg, end string) *VacDates {
	vd := &VacDates{Object: tools.O()}
	vd.Dates = []string{beg, end}
	return vd
}

type ActorVacancyEditModalModel struct {
	*actorupdatemodal.ActorModalModel

	VacationDates []*VacDates `js:"VacationDates"`
}

func NewActorVacancyEditModalModel(vm *hvue.VM) *ActorVacancyEditModalModel {
	vemm := &ActorVacancyEditModalModel{ActorModalModel: actorupdatemodal.NewActorModalModel(vm)}
	vemm.VacationDates = []*VacDates{}
	return vemm
}

func ActorVacancyEditModalModelFromJS(o *js.Object) *ActorVacancyEditModalModel {
	return &ActorVacancyEditModalModel{ActorModalModel: actorupdatemodal.ActorModalModelFromJS(o)}
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Component Methods

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("actor-vacancy-edit-modal", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewActorVacancyEditModalModel(vm)
		}),
		hvue.MethodsOf(&ActorVacancyEditModalModel{}),
		hvue.Computed("isNewActor", func(vm *hvue.VM) interface{} {
			vemm := ActorVacancyEditModalModelFromJS(vm.Object)
			return vemm.CurrentActor.Id == -1
		}),
		hvue.Computed("hasChanged", func(vm *hvue.VM) interface{} {
			vemm := ActorVacancyEditModalModelFromJS(vm.Object)
			return vemm.HasChanged()
		}),
		//hvue.Computed("currentActorRef", func(vm *hvue.VM) interface{} {
		//	aumm := ActorVacancyEditModalModelFromJS(vm.Object)
		//	aumm.CurrentActor.Ref = aumm.CurrentActor.LastName + " " + aumm.CurrentActor.FirstName
		//	return aumm.CurrentActor.Ref
		//}),
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Modal Methods

func (vemm *ActorVacancyEditModalModel) SortDates() {
	vemm.Get("VacationDates").Call("sort", func(a, b *VacDates) int {
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

func (vemm *ActorVacancyEditModalModel) UpdateDates() {
	vemm.VacationDates = []*VacDates{}
	for _, vacPeriod := range vemm.CurrentActor.Vacation {
		vemm.VacationDates = append(vemm.VacationDates, NewVacDates(vacPeriod.Begin, vacPeriod.End))
	}
	vemm.SortDates()
}

func (vemm *ActorVacancyEditModalModel) UpdateVacation() {
	vemm.CurrentActor.Vacation = []*date.DateRange{}
	for _, vacDates := range vemm.VacationDates {
		vemm.CurrentActor.Vacation = append(vemm.CurrentActor.Vacation, date.NewDateRangeFrom(vacDates.Dates[0], vacDates.Dates[1]))
	}
}

func (vemm *ActorVacancyEditModalModel) Show(act *actor.Actor, user *fm.User) {
	vemm.ActorModalModel.Show(act, user)
	vemm.UpdateDates()
}

func (vemm *ActorVacancyEditModalModel) Hide() {
	vemm.ActorModalModel.Hide()
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Action Button Methods

func (vemm *ActorVacancyEditModalModel) ConfirmChange() {
	vemm.UpdateVacation()
	vemm.ActorModalModel.ConfirmChange()
}

func (vemm *ActorVacancyEditModalModel) UndoChange() {
	vemm.ActorModalModel.UndoChange()
	vemm.UpdateDates()
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Tools Button Methods

func (vemm *ActorVacancyEditModalModel) DeleteDates(vm *hvue.VM, index int) {
	vemm = ActorVacancyEditModalModelFromJS(vm.Object)
	vemm.Get("VacationDates").Call("splice", index, 1)
	vemm.UpdateVacation()
}

func (vemm *ActorVacancyEditModalModel) ApplyChange(vm *hvue.VM) {
	vemm = ActorVacancyEditModalModelFromJS(vm.Object)
	vemm.SortDates()
	vemm.UpdateVacation()
}
