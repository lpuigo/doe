package actorinfostable

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/actorstable"
	"github.com/lpuig/ewin/doe/website/frontend/model/actor"
	"github.com/lpuig/ewin/doe/website/frontend/model/actor/actorconst"
	"github.com/lpuig/ewin/doe/website/frontend/tools/date"
	"strconv"
)

const (
	template string = `
<el-table
        :border=true
        :data="filteredActors"
		:default-sort = "{prop: 'Ref', order: 'ascending'}"
        :row-class-name="TableRowClassName" height="100%" size="mini"
		@row-dblclick="HandleDoubleClickedRow"
>
	<!--	index -->
	<el-table-column
			label="N°" width="40px" align="right"
			type="index"
			index=1 
	></el-table-column>

	<!--	Company -->
    <el-table-column
            :resizable="true" :show-overflow-tooltip=true 
            prop="Company" label="Société" width="110px"
			sortable :sort-by="['Company', 'State', 'Role', 'Ref']"
			:filters="FilterList('Company')" :filter-method="FilterHandler"	filter-placement="bottom-end"
    ></el-table-column>
    
	<!--	Contract -->
	<el-table-column
            :resizable="true" :show-overflow-tooltip=true
            prop="Contract" label="Contrat" width="110px"
    ></el-table-column>
    
<!--    <el-table-column-->
<!--            :resizable="true" :show-overflow-tooltip=true -->
<!--            prop="Client" label="Clients" width="200px"-->
<!--			sortable :sort-method="SortClient"-->
<!--			:filters="FilterList('Client')" :filter-method="FilterHandler"	filter-placement="bottom-end"-->
<!--    >-->
<!--        <template slot-scope="scope">-->
<!--			<span>{{GetClients(scope.row)}}</span>-->
<!--        </template>-->
<!--	</el-table-column>-->

	<!--	Role -->
    <el-table-column
            :resizable="true" :show-overflow-tooltip=true 
            prop="Role" label="Rôle" width="110px"
			sortable :sort-by="['Role', 'State', 'Ref']"
			:filters="FilterList('Role')" :filter-method="FilterHandler"	filter-placement="bottom-end"
    ></el-table-column>
    
	<!--	Last & First Name -->
    <el-table-column
            :resizable="true" :show-overflow-tooltip=true 
            prop="Ref" label="Nom Prénom" width="170px"
			sortable :sort-by="['Ref']"
    >
        <template slot-scope="scope">
            <div class="header-menu-container on-hover">
            	<span>{{scope.row.Ref}}</span>
				<i v-if="user.Permissions.HR" class="show link fas fa-edit" @click="EditActor(scope.row)"></i>
            </div>
        </template>
	</el-table-column>

	<!--	State -->
    <el-table-column
            :resizable="true" :show-overflow-tooltip=true 
            prop="State" label="Statut" width="80px"
			:formatter="FormatState"
			:filters="FilterList('State')" :filter-method="FilterHandler"	filter-placement="bottom-end" :filtered-value="FilteredStatusValue()"
    ></el-table-column>

	<!--	Salary -->
    <el-table-column
            :resizable="true" :show-overflow-tooltip=true 
            label="Salaire" width="100px" align="center"
    >
        <template slot-scope="scope">
			<span>{{FormatSalary(scope.row)}}</span>
        </template>
	</el-table-column>

	<!--	Travel Subsidy -->
    <el-table-column
            :resizable="true" :show-overflow-tooltip=true 
            label="Frais Dépl." width="100px" align="center"
    >
        <template slot-scope="scope">
			<span>{{FormatTravelSubsidy(scope.row)}}</span>
        </template>
	</el-table-column>

	<!--	Bonus to pay -->
    <el-table-column
            :resizable="true"
            label="Bonus à payer" width="100px"
    >
        <template slot-scope="scope">
			<span>{{FormatToPayBonuses(scope.row)}}</span>
        </template>
	</el-table-column>

	<!--	Bonuses -->
    <el-table-column
            :resizable="true"
            label="Bonus" width="140px"
    >
        <template slot-scope="scope">
			<span>{{FormatBonuses(scope.row)}}</span>
        </template>
	</el-table-column>

	<!--	Trainings -->
    <el-table-column v-for="trainingName in GetTrainingNames()"
            :resizable="true" align="center"
            :label="trainingName" width="90px"
    >
        <template slot-scope="scope">
			<span>{{FormatTrainingDate(scope.row, trainingName)}}</span>
        </template>
	</el-table-column>

	<!--	Debug -->
<!--    <el-table-column-->
<!--            :resizable="true"-->
<!--            label="Debug Info"-->
<!--    >-->
<!--        <template slot-scope="scope">-->
<!--			<pre>{{scope.row.Info}}</pre>-->
<!--        </template>-->
<!--	</el-table-column>-->
</el-table>
`
)

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("actorinfos-table", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.Props("value", "user", "filter", "filtertype"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewActorsInfoTableModel(vm)
		}),
		hvue.MethodsOf(&ActorsInfoTableModel{}),
		hvue.Computed("filteredActors", func(vm *hvue.VM) interface{} {
			acm := ActorsInfoTableModelFromJS(vm.Object)
			return acm.GetFilteredActors()
		}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type ActorsInfoTableModel struct {
	*actorstable.ActorsTableModel
}

func NewActorsInfoTableModel(vm *hvue.VM) *ActorsInfoTableModel {
	aitm := &ActorsInfoTableModel{ActorsTableModel: actorstable.NewActorsTableModel(vm)}
	return aitm
}

func ActorsInfoTableModelFromJS(o *js.Object) *ActorsInfoTableModel {
	return &ActorsInfoTableModel{ActorsTableModel: actorstable.ActorsTableModelFromJS(o)}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Column related Methods

func (aitm *ActorsInfoTableModel) FormatSalary(act *actor.Actor) string {
	currentDac := act.Info.Salary.CurrentDateAmountComment()
	if currentDac == nil {
		return "-"
	}
	suffix := currentDac.Comment
	switch act.Contract {
	case actorconst.ContractTemp:
		suffix = "€/h"
	case actorconst.ContractCDD, actorconst.ContractCDI:
		suffix = "€"
	}
	res := strconv.FormatFloat(currentDac.Amount, 'f', 2, 64) + " " + suffix
	return res
}

func (aitm *ActorsInfoTableModel) FormatTravelSubsidy(act *actor.Actor) string {
	currentDac := act.Info.TravelSubsidy.CurrentDateAmountComment()
	if currentDac == nil {
		return "-"
	}
	suffix := currentDac.Comment
	res := strconv.FormatFloat(currentDac.Amount, 'f', 2, 64) + " " + suffix
	return res
}

func (aitm *ActorsInfoTableModel) FormatBonuses(act *actor.Actor) string {
	if len(act.Info.EarnedBonuses) == 0 {
		return "-"
	}
	res := ""
	thisMonthDate := date.GetFirstOfMonth(date.TodayAfter(0))
	//lastMonthDate := date.GetFirstOfMonth(date.After(thisMonthDate, -2))

	nbBonus := 0
	for _, bonus := range act.Info.EarnedBonuses {
		if nbBonus == 2 {
			break
		}
		if bonus.Date > thisMonthDate {
			continue // skip future bonuses
		}
		if nbBonus > 0 {
			res += "\n"
		}
		res += date.MonthYear(bonus.Date) + " : " + strconv.FormatFloat(bonus.Amount, 'f', 2, 64) + " €"
		nbBonus++
	}
	if res == "" {
		res = "-"
	}
	return res
}

func (aitm *ActorsInfoTableModel) FormatToPayBonuses(act *actor.Actor) string {
	if len(act.Info.EarnedBonuses) == 0 {
		return "-"
	}

	thisMonthDate := date.GetFirstOfMonth(date.TodayAfter(0))

	monthes := map[string]float64{}
	for _, dac := range act.Info.EarnedBonuses {
		if dac.Date > thisMonthDate { // ignore future bonuses
			continue
		}
		monthes[dac.Date] += dac.Amount
	}
	for _, dac := range act.Info.PaidBonuses {
		monthes[dac.Date] -= dac.Amount
		if monthes[dac.Date] < 0.0001 {
			monthes[dac.Date] = 0
		}
	}

	var amount float64
	for _, a := range monthes {
		amount += a
	}

	if amount < 0.001 {
		return "-"
	}

	return strconv.FormatFloat(amount, 'f', 2, 64) + " €"
}

func (aitm *ActorsInfoTableModel) GetTrainingNames() []string {
	return actor.GetDefaultInfoTraining()
}

func (aitm *ActorsInfoTableModel) FormatTrainingDate(vm *hvue.VM, act *actor.Actor, training string) string {
	aitm = ActorsInfoTableModelFromJS(vm.Object)
	if dc, found := act.Info.Trainings[training]; found {
		return date.DateString(dc.Date)
	}
	return "-"
}
