package actorscalendar

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/model/actor"
	"github.com/lpuig/ewin/doe/website/frontend/model/actor/actorconst"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/date"
	"strings"
)

const (
	//:row-class-name="TableRowClassName"
	template string = `<el-table
			:border=true
			:data="filteredActors"
			:default-sort = "{prop: 'Client', order: 'ascending'}"
			height="100%" size="mini"
	>
		<el-table-column
				:resizable="true" :show-overflow-tooltip=true 
				prop="Company" label="Société" width="110px"
				sortable :sort-by="['Company', 'State', 'Role', 'Ref']"
		></el-table-column>
		
		<el-table-column
				:resizable="true" :show-overflow-tooltip=true 
				prop="Client" label="Clients" width="200px"
				sortable :sort-method="SortClient"
		>
			<template slot-scope="scope">
				<span>{{GetClients(scope.row)}}</span>
			</template>
		</el-table-column>
	
		<el-table-column
				:resizable="true" :show-overflow-tooltip=true 
				prop="Role" label="Rôle" width="110px"
				sortable :sort-by="['Role', 'State', 'Ref']"
		></el-table-column>
		
		<el-table-column
				:resizable="true" :show-overflow-tooltip=true 
				prop="Ref" label="Nom Prénom" width="200px"
				sortable :sort-by="['Ref']"
		>
			<template slot-scope="scope">
				<div class="header-menu-container on-hover">
					<span>{{scope.row.Ref}}</span>
					<i class="show link fas fa-edit" @click="EditActorVacancy(scope.row)"></i>
				</div>
			</template>
		</el-table-column>
		
		<el-table-column
		>
			<template slot="header" slot-scope="scope">
					<el-button icon="fas fa-chevron-left" size="mini" @click="CurrentDateBefore()"></el-button>
					<span style="margin: 0px 10px">{{CurrentDateRange()}}</span>
					<el-button icon="fas fa-chevron-right" size="mini" @click="CurrentDateAfter()"></el-button>
			</template>        
			<template slot-scope="scope">
				<div class="calendar-row">
					<div v-for="(dayClass, index) in GetClassStateFor(scope.row)"
						:key="index"
						class="calendar-slot"
						:class="dayClass"
					>&nbsp;</div>
				</div>
			</template>
		</el-table-column>
		
	</el-table>
	`
)

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("actors-calendar", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.Props("value", "user", "filter", "filtertype"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewActorsCalendarModel(vm)
		}),
		hvue.MethodsOf(&ActorsCalendarModel{}),
		hvue.Computed("filteredActors", func(vm *hvue.VM) interface{} {
			acm := ActorsCalendarModelFromJS(vm.Object)
			return acm.GetFilteredActors()
		}),
		//hvue.Filter("FormatTronconRef", func(vm *hvue.VM, value *js.Object, args ...*js.Object) interface{} {
		//	rpum := ActorsCalendarModelFromJS(vm.Object)
		//	t := &fm.Troncon{Object: value}
		//	return rpum.GetFormatTronconRef(t)
		//}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type ActorsCalendarModel struct {
	*js.Object

	Actors      []*actor.Actor `js:"value"`
	User        *fm.User       `js:"user"`
	Filter      string         `js:"filter"`
	FilterType  string         `js:"filtertype"`
	CurrentDate string         `js:"CurrentDate"`
	DateRange   int            `js:"DateRange"`

	VM *hvue.VM `js:"VM"`
}

func NewActorsCalendarModel(vm *hvue.VM) *ActorsCalendarModel {
	acm := &ActorsCalendarModel{Object: tools.O()}
	acm.VM = vm
	acm.Actors = []*actor.Actor{}
	acm.User = nil
	acm.Filter = ""
	acm.FilterType = ""
	acm.CurrentDate = date.GetMonday(date.TodayAfter(0))
	acm.DateRange = 28
	return acm
}

func ActorsCalendarModelFromJS(o *js.Object) *ActorsCalendarModel {
	return &ActorsCalendarModel{Object: o}
}

func (acm *ActorsCalendarModel) GetFilteredActors() []*actor.Actor {
	acts := acm.GetInRangeActors()
	if acm.FilterType == actorconst.FilterValueAll && acm.Filter == "" {
		return acts
	}
	res := []*actor.Actor{}
	expected := strings.ToUpper(acm.Filter)
	filter := func(a *actor.Actor) bool {
		sis := a.SearchString(acm.FilterType)
		if sis == "" {
			return false
		}
		return strings.Contains(strings.ToUpper(sis), expected)
	}

	for _, actor := range acts {
		if filter(actor) {
			res = append(res, actor)
		}
	}
	return res
}

func (acm *ActorsCalendarModel) CurrentRangeEnd() string {
	return date.After(acm.CurrentDate, acm.DateRange-1)
}

func (acm *ActorsCalendarModel) GetInRangeActors() []*actor.Actor {
	rangeDeb := acm.CurrentDate
	rangeEnd := acm.CurrentRangeEnd()
	res := []*actor.Actor{}
	for _, act := range acm.Actors {
		if act.Period.Begin > rangeEnd {
			continue // actors came in after current period
		}
		if !tools.Empty(act.Period.End) && act.Period.End < rangeDeb {
			continue // actors left out before current period
		}
		res = append(res, act)
	}
	return res
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Format & Style Functions

func (acm *ActorsCalendarModel) TableRowClassName(rowInfo *js.Object) string {
	actor := actor.NewActorFromJS(rowInfo.Get("row"))
	return GetRowStyle(actor)
}

//func (acm *ActorsCalendarModel) GetHoliday(act *actor.Actor) string {
//	if act.State == actorconst.StateCandidate {
//		return "débute " + date.DateString(act.Period.Begin)
//	}
//	vacPeriod := act.GetNextVacation()
//	if vacPeriod == nil {
//		return ""
//	}
//	return "du " + date.DateString(vacPeriod.Begin) + " au " + date.DateString(vacPeriod.End)
//}

func (acm *ActorsCalendarModel) GetClients(act *actor.Actor) string {
	return strings.Join(act.Client, ", ")
}

func (acm *ActorsCalendarModel) SortClient(a, b *actor.Actor) int {
	ca, cb := acm.GetClients(a), acm.GetClients(b)
	switch {
	case ca == cb:
		switch {
		case a.State == b.State:
			return acm.SortRoleRef(a, b)
		case a.State < b.State:
			return -1
		default:
			return 1
		}
	case ca < cb:
		return -1
	default:
		return 1
	}
}

func (acm *ActorsCalendarModel) SortRoleRef(a, b *actor.Actor) int {
	switch {
	case a.Role == b.Role:
		switch {
		case a.Ref == b.Ref:
			return 0
		case a.Ref < b.Ref:
			return -1
		default:
			return 1
		}
	case a.Role < b.Role:
		return -1
	default:
		return 1
	}
}

//func (acm *ActorsCalendarModel) FormatState(row, column, cellValue, index *js.Object) string {
//	return GetStateLabel(cellValue.String())
//}

func (acm *ActorsCalendarModel) GetClassStateFor(vm *hvue.VM, act *actor.Actor) []string {
	rangeStart := acm.CurrentDate
	rangeEnd := acm.CurrentRangeEnd()
	rangeLength := acm.DateRange

	today := int(date.NbDaysBetween(rangeStart, date.TodayAfter(0)))

	// arrival / departure
	arrival := 0
	if act.Period.Begin > rangeStart {
		// Actor arrived afterward
		arrival = int(date.NbDaysBetween(rangeStart, act.Period.Begin))
	}

	departure := rangeLength
	if !tools.Empty(act.Period.End) && act.Period.End <= rangeEnd {
		// Actor left before period end
		departure -= int(date.NbDaysBetween(act.Period.End, rangeEnd))
	}

	// Vancancy
	isVas := make([]bool, rangeLength)
	for _, vacPeriod := range act.Vacation {
		if vacPeriod.End < rangeStart || vacPeriod.Begin > rangeEnd {
			continue
		}
		vacPeriodBeg := 0
		vacPeriodEnd := rangeLength
		if vacPeriod.Begin > rangeStart {
			vacPeriodBeg = int(date.NbDaysBetween(rangeStart, vacPeriod.Begin))
		}
		if vacPeriod.End < rangeEnd {
			vacPeriodEnd -= int(date.NbDaysBetween(vacPeriod.End, rangeEnd))
		}
		for i := vacPeriodBeg; i < vacPeriodEnd; i++ {
			isVas[i] = true
		}
	}

	// calc class array
	res := make([]string, rangeLength)
	for i := 0; i < rangeLength; i++ {
		if i == today {
			res[i] = "today "
		}
		if !(i >= arrival && i < departure) {
			res[i] += "inactive"
			continue
		}
		if isVas[i] {
			res[i] += "holiday"
			continue
		}
		if i%7 > 4 {
			res[i] += "week-end"
			continue
		}
		res[i] += "active"
	}
	return res
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Tools Functions

func GetRowStyle(actor *actor.Actor) string {
	switch actor.State {
	case actorconst.StateCandidate:
		return "actor-row-candidate"
	case actorconst.StateActive:
		return "actor-row-active"
	case actorconst.StateOnHoliday:
		return "actor-row-holiday"
	case actorconst.StateGone:
		return "actor-row-gone"
	default:
		return "actor-row-error"
	}
}

func GetStateLabel(state string) string {
	switch state {
	case actorconst.StateCandidate:
		return actorconst.StateLabelCandidate
	case actorconst.StateActive:
		return actorconst.StateLabelActive
	case actorconst.StateOnHoliday:
		return actorconst.StateLabelOnHoliday
	case actorconst.StateGone:
		return actorconst.StateLabelGone
	default:
		return "Erreur"
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Action Methods

func (acm *ActorsCalendarModel) CurrentDateBefore() {
	acm.CurrentDate = date.After(acm.CurrentDate, -7)
}

func (acm *ActorsCalendarModel) CurrentDateAfter() {
	acm.CurrentDate = date.After(acm.CurrentDate, 7)
}

func (acm *ActorsCalendarModel) CurrentDateRange() string {
	return date.DateString(acm.CurrentDate) + " à " + date.DateString(acm.CurrentRangeEnd())
}

func (acm *ActorsCalendarModel) EditActorVacancy(vm *hvue.VM, act *actor.Actor) {
	vm.Emit("edit-actor-vacancy", act)
}
