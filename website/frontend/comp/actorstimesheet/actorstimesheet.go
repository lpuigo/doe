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
	"strconv"
	"strings"
)

const (
	template string = `<el-container style="height: 100%">
    <el-header style="height: auto; padding: 0px">
        <el-row type="flex" align="middle">
            <el-col :offset="8" :span="2">
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
            <el-col :offset="1" :span="5">
                <el-tooltip content="Semaine précédente" placement="bottom" effect="light" open-delay=500>
                	<el-button icon="fas fa-chevron-left" size="mini" @click="HandleCurrentDateBefore()"></el-button>
				</el-tooltip>
                <span style="margin: 0px 10px">{{CurrentDateRange()}}</span>
                <el-tooltip content="Semaine suivante" placement="bottom" effect="light" open-delay=500>
                	<el-button icon="fas fa-chevron-right" size="mini" @click="HandleCurrentDateAfter()"></el-button>
				</el-tooltip>
                <el-tooltip content="Dernière semaine passée" placement="bottom" effect="light" open-delay=500>
                	<el-button icon="fas fa-chevron-down" size="mini" @click="HandleResetCurrentDate()"></el-button>
				</el-tooltip>
            </el-col>

            <el-col :offset="1" :span="6">
				<el-progress :text-inside="true" :stroke-width="20" :percentage="GetInputPct()"></el-progress>
            </el-col>
        </el-row>
    </el-header>

    <div style="height: 100%;overflow-x: hidden;overflow-y: auto;padding: 0px 0px; margin-top: 8px">
        <el-table
                :border=true
                :data="filteredActors"
                :default-sort = "{prop: 'Groups', order: 'ascending'}"
                height="100%" size="mini"
        >
			<!--	company   -->
            <el-table-column
                    :resizable="true" :show-overflow-tooltip=true 
                    prop="Company" label="Société" width="100px"
                    sortable :sort-by="['Company', 'Ref']"
                    :filters="FilterList('Company')" :filter-method="FilterHandler"	filter-placement="bottom-end"
            ></el-table-column>
            
			<!--	group   -->
			<el-table-column
					:resizable="true" :show-overflow-tooltip=true 
					prop="Groups" label="Groupe" width="150px"
					sortable :sort-method="SortGroup"
					:filters="FilterList('Groups')" :filter-method="FilterHandler"	filter-placement="bottom-end"
			>
				<template slot-scope="scope">
					<span>{{GetGroup(scope.row)}}</span>
				</template>
			</el-table-column>

			<!--	Client  
            <el-table-column
                    label="Clients" prop="Client" width="160px"
                    :resizable="true" :show-overflow-tooltip=true 
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
                    label="Rôle" prop="Role" width="110px"
                    :resizable="true" :show-overflow-tooltip=true 
                    sortable :sort-by="['Role', 'Ref']"
                    :filters="FilterList('Role')" :filter-method="FilterHandler" filter-placement="bottom-end"
            ></el-table-column>
            
			<!--	name   -->
            <el-table-column
                    label="Nom Prénom" prop="Ref" width="160px"
                    :resizable="true" :show-overflow-tooltip=true 
                    sortable :sort-by="['Ref']"
            >
                <template slot-scope="scope">
                    <div class="header-menu-container on-hover">
                        <span>{{scope.row.Ref}}</span>
                        <i class="show link fas fa-edit" @click="EditActorVacancy(scope.row)"></i>
                    </div>
                </template>
            </el-table-column>

			<!--	Hours   -->
            <el-table-column
                    label="Heures" width="60px" align="center"
            >
                <template slot-scope="scope">
                    <span>{{GetActiveHours(scope.row.Id)}}</span>
                </template>
            </el-table-column>

			<!--	activity  -->
			<el-table-column
					label="Activité" width="80px" align="center"
			>
                <template slot-scope="scope">
                    <div style="padding: 3px 0px">
						<el-tooltip content="Remplir la semaine" placement="left" effect="light" open-delay=500>
							<el-button 
									type="primary" size="mini" icon="fas fa-calendar-check" 
									@click="SetActorWeek(scope.row)"
									:disabled="!TimesheetLoaded"
							></el-button>
						</el-tooltip>
                    </div>                        
                </template>
			</el-table-column>

			<!--	timesheet  -->
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
		hvue.Props("value", "groups", "user", "filter", "filtertype"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewActorsTimeSheetModel(vm)
		}),
		hvue.Mounted(func(vm *hvue.VM) {
			acm := ActorsTimeSheetModelFromJS(vm.Object)
			acm.GetTimeSheet()
		}),
		hvue.BeforeDestroy(func(vm *hvue.VM) {
			print("exiting TimeSheetTable")
		}),
		hvue.MethodsOf(&ActorsTimeSheetModel{}),
		hvue.Computed("filteredActors", func(vm *hvue.VM) interface{} {
			acm := ActorsTimeSheetModelFromJS(vm.Object)
			return acm.GetFilteredActors()
		}),
		hvue.Computed("IsDirty", func(vm *hvue.VM) interface{} {
			acm := ActorsTimeSheetModelFromJS(vm.Object)
			acm.Dirty = acm.CheckReference()
			return acm.Dirty
		}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type ActorsTimeSheetModel struct {
	*actorstable.ActorsTableModel

	CurrentDate string `js:"CurrentDate"`
	DateRange   int    `js:"DateRange"`

	TimeSheet       *timesheet.TimeSheet `js:"TimeSheet"`
	ProgressWorking map[int]int          `js:"ProgressWorking"`
	TimesheetLoaded bool                 `js:"TimesheetLoaded"`
	Reference       string               `js:"Reference"`
	Dirty           bool                 `js:"Dirty"`
}

func NewActorsTimeSheetModel(vm *hvue.VM) *ActorsTimeSheetModel {
	acm := &ActorsTimeSheetModel{ActorsTableModel: actorstable.NewActorsTableModel(vm)}
	acm.ResetCurrentDate()
	acm.DateRange = 6
	acm.TimeSheet = timesheet.NewTimeSheet()
	acm.ProgressWorking = make(map[int]int)
	acm.TimesheetLoaded = false
	acm.Reference = ""
	acm.Dirty = true

	return acm
}

func ActorsTimeSheetModelFromJS(o *js.Object) *ActorsTimeSheetModel {
	return &ActorsTimeSheetModel{ActorsTableModel: actorstable.ActorsTableModelFromJS(o)}
}

func (atsm *ActorsTimeSheetModel) GetFilteredActors() []*actor.Actor {
	acts := atsm.GetInRangeActors()
	if atsm.FilterType == actorconst.FilterValueAll && atsm.Filter == "" {
		return acts
	}
	res := []*actor.Actor{}
	expected := strings.ToUpper(atsm.Filter)
	filter := func(a *actor.Actor) bool {
		sis := a.SearchString(atsm.FilterType)
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

func (atsm *ActorsTimeSheetModel) CurrentRangeEnd() string {
	return date.After(atsm.CurrentDate, atsm.DateRange-1)
}

func (atsm *ActorsTimeSheetModel) GetInRangeActors() []*actor.Actor {
	rangeDeb := atsm.CurrentDate
	rangeEnd := atsm.CurrentRangeEnd()
	res := []*actor.Actor{}
	for _, act := range atsm.Actors {
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

func (atsm *ActorsTimeSheetModel) GetInputPct() float64 {
	twd := 0
	tad := 0
	for actId, wd := range atsm.ProgressWorking {
		actTimes, exist := atsm.TimeSheet.ActorsTimes[actId]
		if !exist {
			continue
		}
		twd += wd
		tad += actTimes.NbActiveDays()
	}
	if twd == 0 {
		return 100.0

	}
	return float64(tad*1000/twd) / 10.0
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Format & Style Functions

func (atsm *ActorsTimeSheetModel) GetHeaderClassState(vm *hvue.VM) []string {
	atsm = ActorsTimeSheetModelFromJS(vm.Object)
	rangeStart := atsm.CurrentDate
	rangeLength := atsm.DateRange

	today := int(date.NbDaysBetween(rangeStart, date.TodayAfter(0)))

	// calc class array
	res := make([]string, rangeLength)
	for i := 0; i < rangeLength; i++ {
		day := date.After(rangeStart, i)
		res[i] = "header"
		switch {
		case atsm.User.IsDayOff(day):
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

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Action Methods

func (atsm *ActorsTimeSheetModel) HandleSaveTimeSheet() {
	callback := func() {
		atsm.GetTimeSheet()
	}
	go atsm.callUpdateActors(callback)
}

func (atsm *ActorsTimeSheetModel) HandleReloadTimeSheet() {
	atsm.ResetToReference()
}

func (atsm *ActorsTimeSheetModel) HandleResetCurrentDate() {
	currentMonday := GetCurrentDate()
	if atsm.CurrentDate == currentMonday {
		return
	}

	atsm.CheckAskSaveDialogBefore(func() {
		atsm.CurrentDate = currentMonday
		atsm.GetTimeSheet()
	})
}

func (atsm *ActorsTimeSheetModel) HandleCurrentDateBefore() {
	atsm.CheckAskSaveDialogBefore(func() {
		atsm.CurrentDateBefore()
		atsm.GetTimeSheet()
	})
}

func (atsm *ActorsTimeSheetModel) HandleCurrentDateAfter() {
	atsm.CheckAskSaveDialogBefore(func() {
		atsm.CurrentDateAfter()
		atsm.GetTimeSheet()
	})
}

func (atsm *ActorsTimeSheetModel) HandleDoubleClickedRow(vm *hvue.VM, act *actor.Actor) {
	atsm.EditActorVacancy(vm, act)
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Tools Methods

func GetCurrentDate() string {
	return date.GetMonday(date.TodayAfter(-7))
}

func (atsm *ActorsTimeSheetModel) ResetCurrentDate() {
	atsm.CurrentDate = GetCurrentDate()
}

func (atsm *ActorsTimeSheetModel) CurrentDateBefore() {
	atsm.CurrentDate = date.After(atsm.CurrentDate, -7)
}

func (atsm *ActorsTimeSheetModel) CurrentDateAfter() {
	atsm.CurrentDate = date.After(atsm.CurrentDate, 7)
}

func (atsm *ActorsTimeSheetModel) CurrentDateRange() string {
	return date.DateString(atsm.CurrentDate) + " à " + date.DateString(atsm.CurrentRangeEnd())
}

func (atsm *ActorsTimeSheetModel) DateOf(i int) string {
	return date.DayMonth(date.After(atsm.CurrentDate, i))
}

func (atsm *ActorsTimeSheetModel) GetReference() string {
	return json.Stringify(atsm.TimeSheet)
}

func (atsm *ActorsTimeSheetModel) SetReference() {
	atsm.Reference = atsm.GetReference()
}

func (atsm *ActorsTimeSheetModel) ResetToReference() {
	if atsm.Reference == "" {
		return
	}
	atsm.TimeSheet = timesheet.TimeSheetFromJS(json.Parse(atsm.Reference))
}

// CheckReference returns true when some data has change
func (atsm *ActorsTimeSheetModel) CheckReference() bool {
	return atsm.Reference != atsm.GetReference()
}

func (atsm *ActorsTimeSheetModel) GetTimeSheet() {
	atsm.TimesheetLoaded = false
	callback := func() {
		atsm.SetReference()
		atsm.ProgressWorking = make(map[int]int)
		atsm.TimesheetLoaded = true
	}
	go atsm.callGetTimeSheet(callback)
}

func (atsm *ActorsTimeSheetModel) GetActorsTime(id int) *timesheet.ActorsTime {
	return atsm.TimeSheet.ActorsTimes[id]
	//at, found := atsm.TimeSheet.ActorsTimes[id]
	//if !found {
	//	at = timesheet.NewActorTime()
	//	atsm.TimeSheet.AddActorTime(id, at)
	//}
	//return at
}

func (atsm *ActorsTimeSheetModel) GetActiveDays(vm *hvue.VM, act *actor.Actor) []int {
	atsm = ActorsTimeSheetModelFromJS(vm.Object)
	activeDays := act.GetActiveDays(atsm.CurrentDate, atsm.User.DaysOff)
	pw := 0
	for i, val := range activeDays {
		if val == 1 && i < 5 {
			pw++
		}
	}
	atsm.Get("ProgressWorking").SetIndex(act.Id, pw)
	return activeDays
}

func (atsm *ActorsTimeSheetModel) GetActiveHours(id int) string {
	at, found := atsm.TimeSheet.ActorsTimes[id]
	if !found {
		return ""
	}
	hours, supHours := 0, 0
	for i := 0; i < 6; i++ {
		todayHours := at.Hours[i]
		todaySupHours := 0
		if todayHours > 7 {
			todaySupHours = todayHours - 7
			todayHours = 7
		}
		if i > 4 {
			todaySupHours += todayHours
			todayHours = 0
		}
		hours += todayHours
		supHours += todaySupHours
	}
	return strconv.Itoa(hours) + " + " + strconv.Itoa(supHours)
}

func (atsm *ActorsTimeSheetModel) SetActorWeek(vm *hvue.VM, act *actor.Actor) {
	atsm = ActorsTimeSheetModelFromJS(vm.Object)
	atsm.TimeSheet.ActorsTimes[act.Id].SetActiveWeek(act.GetActiveDays(atsm.CurrentDate, atsm.User.DaysOff))
}

func (atsm *ActorsTimeSheetModel) getUpdatedTimeSheet() *timesheet.TimeSheet {
	updatedTS := atsm.TimeSheet.Clone()
	refTS := timesheet.TimeSheetFromJS(json.Parse(atsm.Reference))
	updatedTS.AddUpdatedActorsTimes(refTS, atsm.TimeSheet)
	return updatedTS
}

func (atsm *ActorsTimeSheetModel) CheckAskSaveDialogBefore(callback func()) {
	if !atsm.Dirty {
		callback()
		return
	}
	message.ConfirmCancelWarning(atsm.VM, "Sauvegarder les modifications apportées ?",
		func() { // confirm
			go atsm.callUpdateActors(func() { callback() })
		},
		func() { // cancel
			callback()
		},
	)
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// WS call Methods

func (atsm *ActorsTimeSheetModel) callGetTimeSheet(callback func()) {
	req := xhr.NewRequest("GET", "/api/timesheet/"+atsm.CurrentDate)
	req.Timeout = tools.LongTimeOut
	req.ResponseType = xhr.JSON
	err := req.Send(nil)
	if err != nil {
		message.ErrorStr(atsm.VM, "Oups! "+err.Error(), true)
		return
	}
	if req.Status != tools.HttpOK {
		message.ErrorRequestMessage(atsm.VM, req)
		return
	}
	atsm.TimeSheet = timesheet.TimeSheetFromJS(req.Response)
	callback()
}

func (atsm *ActorsTimeSheetModel) callUpdateActors(callback func()) {
	updatedTs := atsm.getUpdatedTimeSheet()
	defer callback()
	if len(updatedTs.ActorsTimes) == 0 {
		message.ErrorStr(atsm.VM, "Could not find any updated Actors Time", false)
		return
	}

	req := xhr.NewRequest("PUT", "/api/timesheet/"+updatedTs.WeekDate)
	req.Timeout = tools.TimeOut
	req.ResponseType = xhr.JSON
	err := req.Send(json.Stringify(updatedTs))
	if err != nil {
		message.ErrorStr(atsm.VM, "Oups! "+err.Error(), true)
		return
	}
	if req.Status != tools.HttpOK {
		message.ErrorRequestMessage(atsm.VM, req)
		return
	}
	message.NotifySuccess(atsm.VM, "Pointage Horaire", "Modifications sauvegardées")
}
