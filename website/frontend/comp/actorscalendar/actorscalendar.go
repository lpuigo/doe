package actorscalendar

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/actorstable"
	"github.com/lpuig/ewin/doe/website/frontend/model/actor"
	"github.com/lpuig/ewin/doe/website/frontend/model/actor/actorconst"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/date"
	"strings"
)

const (
	template string = `<el-container style="height: 100%">
    <el-header style="height: auto; padding: 0px">
        <el-row type="flex" align="middle">
            <el-col :offset="13" :span="1">
                <el-button icon="fas fa-chevron-down" size="mini" @click="ResetCurrentDate()"></el-button>
            </el-col>
            <el-col :span="10">
                <el-button icon="fas fa-chevron-left" size="mini" @click="CurrentDateBefore()"></el-button>
                <span style="margin: 0px 10px">{{CurrentDateRange()}}</span>
                <el-button icon="fas fa-chevron-right" size="mini" @click="CurrentDateAfter()"></el-button>
            </el-col>
        </el-row>
    </el-header>

    <div style="height: 100%;overflow-x: hidden;overflow-y: auto;padding: 0px 0px; margin-top: 8px">
        <el-table
                :border=true
                :data="filteredActors"
                :default-sort = "{prop: 'Client', order: 'ascending'}"
                height="100%" size="mini"
				@row-dblclick="HandleDoubleClickedRow"
        >
            <el-table-column
                    :resizable="true" :show-overflow-tooltip=true 
                    prop="Company" label="Société" width="110px"
                    sortable :sort-by="['Company', 'State', 'Role', 'Ref']"
                    :filters="FilterList('Company')" :filter-method="FilterHandler"	filter-placement="bottom-end"
            ></el-table-column>
            
            <el-table-column
                    :resizable="true" :show-overflow-tooltip=true 
                    prop="Client" label="Clients" width="160px"
                    sortable :sort-method="SortClient"
                    :filters="FilterList('Client')" :filter-method="FilterHandler" filter-placement="bottom-end"
            >
                <template slot-scope="scope">
                    <span>{{GetClients(scope.row)}}</span>
                </template>
            </el-table-column>
        
            <el-table-column
                    :resizable="true" :show-overflow-tooltip=true 
                    prop="Role" label="Rôle" width="90px"
                    sortable :sort-by="['Role', 'State', 'Ref']"
                    :filters="FilterList('Role')" :filter-method="FilterHandler" filter-placement="bottom-end"
            ></el-table-column>
            
            <el-table-column
                    :resizable="true" :show-overflow-tooltip=true 
                    prop="Ref" label="Nom Prénom" width="180px"
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
                    <div class="calendar-row">
                        <div v-for="(dayClass, index) in GetHeaderClassState()"
                             :key="index"
                             class="calendar-slot"
                             :class="dayClass"
                        >{{DateOf(index)}}</div>
                    </div>
                </template>        
                <template slot-scope="scope">
                    <div class="calendar-row">
                        <div v-for="(dayClass, index) in GetClassStateFor(scope.row)"
                            :key="index"
                            class="calendar-slot"
                            :class="dayClass.Class"
                        >
							<el-tooltip v-if="dayClass.Comment != ''" placement="top" open-delay=200>
								<div slot="content">{{dayClass.Comment}}</div>
								<i class="fas fa-info icon--small"></i>
							</el-tooltip>
							<span v-else>&nbsp;</span>
                        </div>
                    </div>
                </template>
            </el-table-column>
        </el-table>
    </div>
</el-container>	`
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
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type ActorsCalendarModel struct {
	*actorstable.ActorsTableModel

	CurrentDate string `js:"CurrentDate"`
	DateRange   int    `js:"DateRange"`
}

func NewActorsCalendarModel(vm *hvue.VM) *ActorsCalendarModel {
	acm := &ActorsCalendarModel{ActorsTableModel: actorstable.NewActorsTableModel(vm)}
	acm.ResetCurrentDate()
	acm.DateRange = 5 * 7
	return acm
}

func ActorsCalendarModelFromJS(o *js.Object) *ActorsCalendarModel {
	return &ActorsCalendarModel{ActorsTableModel: actorstable.ActorsTableModelFromJS(o)}
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
		if act.Period.Begin == "" {
			continue // Defected actors
		}
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

func (acm *ActorsCalendarModel) GetHeaderClassState(vm *hvue.VM) []string {
	acm = ActorsCalendarModelFromJS(vm.Object)
	rangeStart := acm.CurrentDate
	rangeLength := acm.DateRange
	today := int(date.NbDaysBetween(rangeStart, date.TodayAfter(0)))

	// calc class array
	res := make([]string, rangeLength)
	for i := 0; i < rangeLength; i++ {
		day := date.After(rangeStart, i)
		res[i] = "header"
		switch {
		case acm.User.IsDayOff(day):
			res[i] += " day-off"
		case i%7 > 4:
			res[i] += " week-end"
		}
		if i == today {
			res[i] += " today"
		}
	}
	return res
}

func (acm *ActorsCalendarModel) GetClassStateFor(vm *hvue.VM, act *actor.Actor) []*CalendarDayInfo {
	acm = ActorsCalendarModelFromJS(vm.Object)
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

	// Vacancy
	isVas := make([]bool, rangeLength)
	vaCmts := make([]string, rangeLength)
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
			vaCmts[i] = vacPeriod.Comment
		}
	}

	// calc class array
	res := make([]*CalendarDayInfo, rangeLength)
	for i := 0; i < rangeLength; i++ {
		res[i] = NewCalendarDayInfo()
		day := date.After(rangeStart, i)
		if i == today {
			res[i].Class = "today "
		}
		switch {
		case !(i >= arrival && i < departure):
			res[i].Class += "inactive "
			continue
		case acm.User.IsDayOff(day):
			res[i].Class += "day-off "
			continue
		case isVas[i]:
			res[i].Class += "holiday "
			res[i].Comment = vaCmts[i]
			continue
		case i%7 > 4:
			res[i].Class += "week-end "
		}
		res[i].Class += "active "
	}
	return res
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Action Methods

func (acm *ActorsCalendarModel) ResetCurrentDate() {
	acm.CurrentDate = date.GetMonday(date.TodayAfter(-7))
}

func (acm *ActorsCalendarModel) CurrentDateBefore() {
	acm.CurrentDate = date.After(acm.CurrentDate, -7)
}

func (acm *ActorsCalendarModel) CurrentDateAfter() {
	acm.CurrentDate = date.After(acm.CurrentDate, 7)
}

func (acm *ActorsCalendarModel) CurrentDateRange() string {
	return date.DateString(acm.CurrentDate) + " à " + date.DateString(acm.CurrentRangeEnd())
}

func (acm *ActorsCalendarModel) DateOf(i int) string {
	return date.Day(date.After(acm.CurrentDate, i))
}

func (acm *ActorsCalendarModel) HandleDoubleClickedRow(vm *hvue.VM, act *actor.Actor) {
	acm.EditActorVacancy(vm, act)
}
