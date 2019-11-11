package actorstimesheet

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/actorstable"
	"github.com/lpuig/ewin/doe/website/frontend/comp/actorstimesheet/actortimeedit"
	"github.com/lpuig/ewin/doe/website/frontend/model/actor"
	"github.com/lpuig/ewin/doe/website/frontend/model/actor/actorconst"
	"github.com/lpuig/ewin/doe/website/frontend/model/timesheet"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/date"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements/message"
	"github.com/lpuig/ewin/doe/website/frontend/tools/json"
	"honnef.co/go/js/xhr"
	"strings"
)

const (
	template string = `<el-container style="height: 100%">
    <el-header style="height: auto; padding: 0px">
        <el-row type="flex" align="middle">
            <el-col :offset="10" :span="2">
                <el-button-group>
                    <el-tooltip content="Enregistrer les temps saisis" placement="bottom" effect="light" open-delay=500>
                        <el-button class="icon" icon="fas fa-cloud-upload-alt icon--medium" @click="HandleSaveTimeSheet"
                                   :disabled="!IsDirty" size="mini"></el-button>
                    </el-tooltip>
                    <el-tooltip content="Annuler les modifications" placement="bottom" effect="light" open-delay=500>
                        <el-button class="icon" icon="fas fa-undo-alt icon--medium" @click="HandleReloadTimeSheet"
                                   :disabled="!IsDirty" size="mini"></el-button>
                    </el-tooltip>
                </el-button-group>
            </el-col>
            <el-col :offset="2" :span="1">
                <el-tooltip content="Dernière semaine passée" placement="bottom" effect="light" open-delay=500>
                	<el-button icon="fas fa-chevron-down" size="mini" @click="HandleResetCurrentDate()"></el-button>
				</el-tooltip>
            </el-col>
            <el-col :span="10">
                <el-tooltip content="Semaine précédente" placement="bottom" effect="light" open-delay=500>
                	<el-button icon="fas fa-chevron-left" size="mini" @click="HandleCurrentDateBefore()"></el-button>
				</el-tooltip>
                <span style="margin: 0px 10px">{{CurrentDateRange()}}</span>
                <el-tooltip content="Semaine suivante" placement="bottom" effect="light" open-delay=500>
                	<el-button icon="fas fa-chevron-right" size="mini" @click="HandleCurrentDateAfter()"></el-button>
				</el-tooltip>
            </el-col>
        </el-row>
    </el-header>

    <div style="height: 100%;overflow-x: hidden;overflow-y: auto;padding: 0px 0px; margin-top: 8px">
        <el-table
                :border=true
                :data="filteredActors"
                :default-sort = "{prop: 'Client', order: 'ascending'}"
                height="100%" size="mini"
        >
            <el-table-column
                    :resizable="true" :show-overflow-tooltip=true 
                    prop="Company" label="Société" width="110px"
                    sortable :sort-by="['Company', 'State', 'Role', 'Ref']"
                    :filters="FilterList('Company')" :filter-method="FilterHandler"	filter-placement="bottom-end"
            ></el-table-column>
            
            <el-table-column
                    label="Clients" prop="Client" width="200px"
                    :resizable="true" :show-overflow-tooltip=true 
                    sortable :sort-method="SortClient"
                    :filters="FilterList('Client')" :filter-method="FilterHandler" filter-placement="bottom-end"
            >
                <template slot-scope="scope">
                    <span>{{GetClients(scope.row)}}</span>
                </template>
            </el-table-column>
        
            <el-table-column
                    label="Rôle" prop="Role" width="110px"
                    :resizable="true" :show-overflow-tooltip=true 
                    sortable :sort-by="['Role', 'State', 'Ref']"
                    :filters="FilterList('Role')" :filter-method="FilterHandler" filter-placement="bottom-end"
            ></el-table-column>
            
            <el-table-column
                    label="Nom Prénom" prop="Ref" width="200px"
                    :resizable="true" :show-overflow-tooltip=true 
                    sortable :sort-by="['Ref']"
            ></el-table-column>

			<el-table-column
					label="Activité" width="80px" align="center"
			>
                <template slot-scope="scope">
                    <div style="padding: 3px 0px">
						<el-tooltip content="Remplir la semaine" placement="bottom" effect="light" open-delay=500>
							<el-button 
									type="primary" size="mini" icon="fas fa-calendar-check" 
									@click="SetActorWeek(scope.row)"
									:disabled="!TimesheetLoaded"
							></el-button>
						</el-tooltip>
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
                    <actor-time-edit v-if="TimesheetLoaded" 
                    		:times="GetActorsTime(scope.row.Id)" 
                    		:activedays="GetActiveDays(scope.row)"
                    ></actor-time-edit>
                </template>
            </el-table-column>
        </el-table>
    </div>
</el-container>	`
)

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("actors-time-sheet", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		actortimeedit.RegisterComponent(),
		hvue.Template(template),
		hvue.Props("value", "user", "filter", "filtertype"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewActorsCalendarModel(vm)
		}),
		hvue.Mounted(func(vm *hvue.VM) {
			acm := ActorsCalendarModelFromJS(vm.Object)
			acm.GetTimeSheet()
		}),
		hvue.BeforeDestroy(func(vm *hvue.VM) {
			print("exiting TimeSheetTable")
		}),
		hvue.MethodsOf(&ActorsCalendarModel{}),
		hvue.Computed("filteredActors", func(vm *hvue.VM) interface{} {
			acm := ActorsCalendarModelFromJS(vm.Object)
			return acm.GetFilteredActors()
		}),
		hvue.Computed("IsDirty", func(vm *hvue.VM) interface{} {
			acm := ActorsCalendarModelFromJS(vm.Object)
			acm.Dirty = acm.CheckReference()
			return acm.Dirty
		}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type ActorsCalendarModel struct {
	*actorstable.ActorsTableModel

	CurrentDate string `js:"CurrentDate"`
	DateRange   int    `js:"DateRange"`

	TimeSheet       *timesheet.TimeSheet `js:"TimeSheet"`
	TimesheetLoaded bool                 `js:"TimesheetLoaded"`
	Reference       string               `js:"Reference"`
	Dirty           bool                 `js:"Dirty"`
}

func NewActorsCalendarModel(vm *hvue.VM) *ActorsCalendarModel {
	acm := &ActorsCalendarModel{ActorsTableModel: actorstable.NewActorsTableModel(vm)}
	acm.ResetCurrentDate()
	acm.DateRange = 6
	acm.TimeSheet = timesheet.NewTimeSheet()
	acm.TimesheetLoaded = false
	acm.Reference = ""
	acm.Dirty = true

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
	rangeStart := acm.CurrentDate
	rangeLength := acm.DateRange

	today := int(date.NbDaysBetween(rangeStart, date.TodayAfter(0)))

	// calc class array
	res := make([]string, rangeLength)
	for i := 0; i < rangeLength; i++ {
		res[i] = "header"
		if i%7 > 4 {
			res[i] += " week-end"
		}
		if i == today {
			res[i] += " today"
		}
	}
	return res
}

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
			res[i] += "inactive "
			continue
		}
		if isVas[i] {
			res[i] += "holiday "
			continue
		}
		if i%7 > 4 {
			res[i] += "week-end "
		}
		res[i] += "active "
	}
	return res
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Action Methods

func (acm *ActorsCalendarModel) HandleSaveTimeSheet() {
	callback := func() {
		acm.GetTimeSheet()
	}
	go acm.callUpdateActors(callback)
}

func (acm *ActorsCalendarModel) HandleReloadTimeSheet() {
	acm.ResetToReference()
}

func (acm *ActorsCalendarModel) HandleResetCurrentDate() {
	if acm.ResetCurrentDate() {
		acm.GetTimeSheet()
	}
}

func (acm *ActorsCalendarModel) HandleCurrentDateBefore() {
	acm.CheckAskSaveDialogBefore(func() {
		acm.CurrentDateBefore()
		acm.GetTimeSheet()
	})
}

func (acm *ActorsCalendarModel) HandleCurrentDateAfter() {
	acm.CheckAskSaveDialogBefore(func() {
		acm.CurrentDateAfter()
		acm.GetTimeSheet()
	})
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Tools Methods

func (acm *ActorsCalendarModel) ResetCurrentDate() bool {
	currentMonday := date.GetMonday(date.TodayAfter(-7))
	if acm.CurrentDate == currentMonday {
		return false
	}
	acm.CurrentDate = currentMonday
	return true
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

func (acm *ActorsCalendarModel) GetReference() string {
	return json.Stringify(acm.TimeSheet)
}

func (acm *ActorsCalendarModel) SetReference() {
	acm.Reference = acm.GetReference()
}

func (acm *ActorsCalendarModel) ResetToReference() {
	if acm.Reference == "" {
		return
	}
	acm.TimeSheet = timesheet.TimeSheetFromJS(json.Parse(acm.Reference))
}

// CheckReference returns true when some data has change
func (acm *ActorsCalendarModel) CheckReference() bool {
	return acm.Reference != acm.GetReference()
}

func (acm *ActorsCalendarModel) GetTimeSheet() {
	acm.TimesheetLoaded = false
	callback := func() {
		acm.SetReference()
		acm.TimesheetLoaded = true
	}
	go acm.callGetTimeSheet(callback)
}

func (acm *ActorsCalendarModel) GetActorsTime(id int) *timesheet.ActorsTime {
	at, found := acm.TimeSheet.ActorsTimes[id]
	if !found {
		at = timesheet.NewActorTime()
		at.Get("Hours").SetIndex(0, -100)
	}
	return at
}

func (acm *ActorsCalendarModel) GetActiveDays(act *actor.Actor) []int {
	return act.GetActiveDays(acm.CurrentDate)
}

func (acm *ActorsCalendarModel) SetActorWeek(vm *hvue.VM, act *actor.Actor) {
	acm = ActorsCalendarModelFromJS(vm.Object)
	acm.TimeSheet.ActorsTimes[act.Id].SetActiveWeek(act.GetActiveDays(acm.CurrentDate))
}

func (acm *ActorsCalendarModel) getUpdatedTimeSheet() *timesheet.TimeSheet {
	updatedTS := acm.TimeSheet.Clone()
	refTS := timesheet.TimeSheetFromJS(json.Parse(acm.Reference))
	updatedTS.AddUpdatedActorsTimes(refTS, acm.TimeSheet)
	return updatedTS
}

func (acm *ActorsCalendarModel) CheckAskSaveDialogBefore(callback func()) {
	if !acm.Dirty {
		callback()
		return
	}
	message.ConfirmCancelWarning(acm.VM, "Sauvegarder les modifications apportées ?",
		func() { // confirm
			go acm.callUpdateActors(func() { callback() })
		},
		func() { // cancel
			callback()
		},
	)
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// WS call Methods

func (acm *ActorsCalendarModel) callGetTimeSheet(callback func()) {
	req := xhr.NewRequest("GET", "/api/timesheet/"+acm.CurrentDate)
	req.Timeout = tools.LongTimeOut
	req.ResponseType = xhr.JSON
	err := req.Send(nil)
	if err != nil {
		message.ErrorStr(acm.VM, "Oups! "+err.Error(), true)
		return
	}
	if req.Status != tools.HttpOK {
		message.ErrorRequestMessage(acm.VM, req)
		return
	}
	acm.TimeSheet = timesheet.TimeSheetFromJS(req.Response)
	callback()
}

func (acm *ActorsCalendarModel) callUpdateActors(callback func()) {
	updatedTs := acm.getUpdatedTimeSheet()
	defer callback()
	if len(updatedTs.ActorsTimes) == 0 {
		message.ErrorStr(acm.VM, "Could not find any updated Actors Time", false)
		return
	}

	req := xhr.NewRequest("PUT", "/api/timesheet/"+updatedTs.WeekDate)
	req.Timeout = tools.TimeOut
	req.ResponseType = xhr.JSON
	err := req.Send(json.Stringify(updatedTs))
	if err != nil {
		message.ErrorStr(acm.VM, "Oups! "+err.Error(), true)
		return
	}
	if req.Status != tools.HttpOK {
		message.ErrorRequestMessage(acm.VM, req)
		return
	}
	message.NotifySuccess(acm.VM, "Pointage Horaire", "Modifications sauvegardées")
}
