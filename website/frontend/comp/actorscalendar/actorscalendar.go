package actorscalendar

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/actorstable"
	"github.com/lpuig/ewin/doe/website/frontend/model/actor"
	"github.com/lpuig/ewin/doe/website/frontend/model/actor/actorconst"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/date"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
	"sort"
	"strconv"
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
                :default-sort = "{prop: 'Actor.Ref', order: 'ascending'}"
                height="100%" size="mini"
				@row-dblclick="HandleDoubleClickedRow"
        >
            <el-table-column
                    label="N°" width="40px" align="right"
                    type="index"
                    index=1 
            ></el-table-column>
        
			<!--	group   -->
			<el-table-column
					:resizable="true" :show-overflow-tooltip=true 
					prop="GroupName" label="Groupe" width="150px"
					sortable :sort-method="SortGroup"
					:filters="FilterList('GroupName')" :filter-method="FilterHandler"	filter-placement="bottom-end"
			></el-table-column>

			<!--	client
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
        	-->

			<!--	role   -->
            <el-table-column
                    :resizable="true" :show-overflow-tooltip=true 
                    prop="Actor.Role" label="Rôle" width="90px"
                    sortable :sort-by="['Actor.Role', 'Actor.Ref']"
                    :filters="FilterList('Role')" :filter-method="FilterHandler" filter-placement="bottom-end"
            ></el-table-column>
            
			<!--	name   -->
            <el-table-column
                    :resizable="true" :show-overflow-tooltip=true 
                    prop="Actor.Ref" label="Nom Prénom" width="180px"
                    sortable :sort-by="['Actor.Ref']"
            >
                <template slot-scope="scope">
                    <div class="header-menu-container on-hover">
                        <span>{{scope.row.Actor.Ref}}</span>
                        <i class="show link fas fa-edit" @click="EditActorVacancy(scope.row.Actor)"></i>
                    </div>
                </template>
            </el-table-column>
            
			<!--	company   -->
            <el-table-column
                    :resizable="true" :show-overflow-tooltip=true 
                    prop="Actor.Company" label="Société" width="110px"
                    sortable :sort-by="['Actor.Company', 'Actor.Ref']"
                    :filters="FilterList('Company')" :filter-method="FilterHandler"	filter-placement="bottom-end"
            ></el-table-column>
            
			<!--	calendar   -->
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
		hvue.Props("value", "groups", "user", "filter", "filtertype"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewActorsCalendarModel(vm)
		}),
		hvue.MethodsOf(&ActorsCalendarModel{}),
		hvue.Computed("filteredActors", func(vm *hvue.VM) interface{} {
			acm := ActorsCalendarModelFromJS(vm.Object)
			acm.GroupActors = acm.GetFilteredGroupActors()
			return acm.GroupActors
		}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type ActorsCalendarModel struct {
	*actorstable.ActorsTableModel

	CurrentDate string        `js:"CurrentDate"`
	DateRange   int           `js:"DateRange"`
	GroupActors []*GroupActor `js:"GroupActors"`
}

func NewActorsCalendarModel(vm *hvue.VM) *ActorsCalendarModel {
	acm := &ActorsCalendarModel{ActorsTableModel: actorstable.NewActorsTableModel(vm)}
	acm.ResetCurrentDate()
	acm.DateRange = 5 * 7
	acm.GroupActors = []*GroupActor{}
	return acm
}

func ActorsCalendarModelFromJS(o *js.Object) *ActorsCalendarModel {
	return &ActorsCalendarModel{ActorsTableModel: actorstable.ActorsTableModelFromJS(o)}
}

func (acm *ActorsCalendarModel) GetFilteredGroupActors() []*GroupActor {
	calendarRange := date.NewDateRangeFrom(acm.CurrentDate, acm.CurrentRangeEnd())
	res := []*GroupActor{}
	for _, actor := range acm.GetFilteredActors() {
		for _, assign := range actor.Groups.GetGroupAssignsInRange(calendarRange) {
			grpName := acm.GroupStore.GetGroupNameById(assign.Id)
			res = append(res, NewGroupActor(actor, grpName, assign.Period))
		}
	}
	return res
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
func (acm *ActorsCalendarModel) SortGroup(vm *hvue.VM, a, b *GroupActor) int {
	switch {
	case a.GroupName == b.GroupName:
		return acm.SortRoleRef(a.Actor, b.Actor)
	case a.GroupName < b.GroupName:
		return -1
	default:
		return 1
	}
}

func (acm *ActorsCalendarModel) FilterHandler(vm *hvue.VM, value string, p *js.Object, col *js.Object) bool {
	prop := col.Get("property").String()
	prop = strings.TrimPrefix(prop, "Actor.")
	switch prop {
	case "GroupName":
		return p.Get(prop).String() == value
	default:
		return p.Get("Actor").Get(prop).String() == value
	}
}

func (acm *ActorsCalendarModel) FilterList(vm *hvue.VM, prop string) []*elements.ValText {
	acm = ActorsCalendarModelFromJS(vm.Object)
	count := map[string]int{}
	attribs := []string{}

	attrib := ""
	for _, act := range acm.GroupActors {
		if prop == "GroupName" {
			attrib = act.Object.Get(prop).String()
		} else {
			attrib = act.Object.Get("Actor").Get(prop).String()
		}
		var attrs []string
		switch prop {
		case "Client":
			attrs = strings.Split(attrib, ",")
		default:
			attrs = []string{attrib}
		}
		for _, a := range attrs {
			if _, exist := count[a]; !exist {
				attribs = append(attribs, a)
			}
			count[a]++
		}
	}
	sort.Strings(attribs)
	res := []*elements.ValText{}
	for _, a := range attribs {
		fa := a
		if fa == "" {
			fa = "Vide"
		}
		res = append(res, elements.NewValText(a, fa+" ("+strconv.Itoa(count[a])+")"))
	}
	return res
}

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

func (acm *ActorsCalendarModel) GetClassStateFor(vm *hvue.VM, o *js.Object) []*CalendarDayInfo {
	acm = ActorsCalendarModelFromJS(vm.Object)
	gract := GroupActorFromJS(o)
	act := gract.Actor
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

	// assignIn / assignOut
	assignIn := 0
	if gract.Assignment.Begin > rangeStart {
		// Actor assigned afterward
		assignIn = int(date.NbDaysBetween(rangeStart, gract.Assignment.Begin))
	}

	assignOut := rangeLength
	if gract.Assignment.End <= rangeEnd {
		// Actor left before period end
		assignOut -= int(date.NbDaysBetween(gract.Assignment.End, rangeEnd)) + 1
	}

	// Vacancy
	isVas := make([]string, rangeLength)
	vaCmts := make([]string, rangeLength)
	for _, vacPeriod := range act.VacInfo.Vacation {
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
			isVas[i] = leavePeriodTypeClass(vacPeriod.Type)
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
		case !(i >= assignIn && i <= assignOut):
			res[i].Class += "inactive "
			continue
		case acm.User.IsDayOff(day):
			res[i].Class += "day-off "
			res[i].Comment = acm.User.DaysOff[day]
			continue
		case isVas[i] != "":
			res[i].Class += "holiday " + isVas[i]
			res[i].Comment = vaCmts[i]
			continue
		case i%7 > 4:
			res[i].Class += "week-end "
		}
		res[i].Class += "active "
	}
	return res
}

func leavePeriodTypeClass(leaveType string) string {
	switch leaveType {
	case actorconst.LeaveTypePaid:
		return "CP"
	case actorconst.LeaveTypeUnpaid:
		return "CSS"
	case actorconst.LeaveTypeSick:
		return "AM"
	case actorconst.LeaveTypeInjury:
		return "AT"
	case actorconst.LeaveTypePublicHoliday:
		return "JF"
	default:
		return ""
	}
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

func (acm *ActorsCalendarModel) HandleDoubleClickedRow(vm *hvue.VM, act *GroupActor) {
	acm.EditActorVacancy(vm, act.Actor)
}
