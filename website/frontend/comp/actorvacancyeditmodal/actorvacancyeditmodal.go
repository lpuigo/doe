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
        width="40%"
>
    <!-- 
        Modal Title
    -->
    <span slot="title">
		<el-row :gutter="10" align="middle" type="flex">
			<el-col :span="24">
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
							@change="ApplyChange()">
					></el-date-picker>
				</template>
			</el-table-column>

			<el-table-column label="Commentaire">
				<template slot-scope="scope">
					<el-input
							type="textarea" placeholder="Raison du congé" size="mini"
							v-model="scope.row.Comment" @input="ApplyChange()">
					</el-input>
				</template>
			</el-table-column>
		</el-table>

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
		hvue.Computed("datesComplete", func(vm *hvue.VM) interface{} {
			vemm := ActorVacancyEditModalModelFromJS(vm.Object)
			return vemm.DatesComplete()
		}),
		hvue.Computed("hasChanged", func(vm *hvue.VM) interface{} {
			vemm := ActorVacancyEditModalModelFromJS(vm.Object)
			return vemm.DatesComplete() && vemm.HasChanged()
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
		vemm.VacationDates = append(vemm.VacationDates, NewVacDates(vacPeriod.Begin, vacPeriod.End, vacPeriod.Comment))
	}
	vemm.SortDates()
}

func (vemm *ActorVacancyEditModalModel) UpdateVacation() {
	vemm.CurrentActor.Vacation = []*date.DateRangeComment{}
	for _, vacDates := range vemm.VacationDates {
		vemm.CurrentActor.Vacation = append(vemm.CurrentActor.Vacation, date.NewDateRangeCommentFrom(vacDates.Dates[0], vacDates.Dates[1], vacDates.Comment))
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

func (vemm *ActorVacancyEditModalModel) DatesComplete() bool {
	if len(vemm.VacationDates) == 0 {
		return true
	}
	return vemm.VacationDates[0].IsComplete()
}

func (vemm *ActorVacancyEditModalModel) AddDates(vm *hvue.VM) {
	vemm = ActorVacancyEditModalModelFromJS(vm.Object)
	vemm.VacationDates = append([]*VacDates{NewVacDates("", "", "")}, vemm.VacationDates...)
	vemm.UpdateVacation()
}

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
