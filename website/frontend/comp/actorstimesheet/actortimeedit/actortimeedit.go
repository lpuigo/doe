package actortimeedit

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/model/actor/actorconst"
	"github.com/lpuig/ewin/doe/website/frontend/model/timesheet"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

const (
	template string = `
	<div>
		<div v-if="times != null" class="calendar-row">
			<div v-for="(hour, index) in times.Hours"
				:key="index"
				class="calendar-slot"
				:class="GetColumnColor(index)"
				style="padding: 2px 0px;"
			>
				<el-row :gutter="10" type="flex" align="middle">
					<el-col :offset="3" :span="6">
						<el-button 
								type="primary" plain icon="fas fa-calendar-plus" size="mini" 
								@click="SetFullDay(index)"
								:disabled="!IsWorkableDay(index)"
						></el-button>
					</el-col>
					<el-col :span="12">
						<el-input-number v-model="times.Hours[index]"
								size="mini" controls-position="right" style="width: 100%"
								@change="HandleChange"
								:min="0" :max="11"
								:disabled="!IsWorkableDay(index)"
						></el-input-number>
					</el-col>
				</el-row>
			</div>
		</div>
		<div v-else>Sauvegarde nécessaire</div>	
	</div>
`
)

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("actor-time-edit",
		componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.Props("times", "activedays"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewActorTimeEditModel(vm)
		}),
		//hvue.Mounted(func(vm *hvue.VM) {
		//	atem := ActorTimeEditModelFromJS(vm.Object)
		//	print("ActorTimeEditModel : ", atem.Times)
		//}),
		hvue.MethodsOf(&ActorTimeEditModel{}),
		//hvue.Computed("filteredActors", func(vm *hvue.VM) interface{} {
		//	acm := ActorsCalendarModelFromJS(vm.Object)
		//	return acm.GetFilteredActors()
		//}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type ActorTimeEditModel struct {
	*js.Object

	Times      *timesheet.ActorsTime `js:"times"`
	ActiveDays []string              `js:"activedays"`

	VM *hvue.VM `js:"VM"`
}

func NewActorTimeEditModel(vm *hvue.VM) *ActorTimeEditModel {
	atem := &ActorTimeEditModel{Object: tools.O()}
	atem.Times = timesheet.NewActorTime()
	atem.ActiveDays = make([]string, 6)
	atem.VM = vm
	return atem
}

func ActorTimeEditModelFromJS(o *js.Object) *ActorTimeEditModel {
	return &ActorTimeEditModel{Object: o}
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Tools Methods

func (atem *ActorTimeEditModel) HandleChange(vm *hvue.VM) {
	//atem = ActorTimeEditModelFromJS(vm.Object)
}

func (atem *ActorTimeEditModel) SetFullDay(vm *hvue.VM, index int) {
	atem = ActorTimeEditModelFromJS(vm.Object)
	atem.Times.Get("Hours").Call("splice", index, 1, 7)
}

func (atem *ActorTimeEditModel) IsWorkableDay(vm *hvue.VM, index int) bool {
	atem = ActorTimeEditModelFromJS(vm.Object)
	return atem.isWorkableDay(index)
}

func (atem *ActorTimeEditModel) isWorkableDay(index int) bool {
	switch atem.ActiveDays[index] {
	case actorconst.LeaveTypeShortWorking:
		return true
	case actorconst.LeaveTypeShortPublicHoliday:
		return true
	default:
		return false
	}
}

func (atem *ActorTimeEditModel) GetColumnColor(vm *hvue.VM, index int) string {
	atem = ActorTimeEditModelFromJS(vm.Object)
	if atem.Times.Hours[index] == 0 {
		switch atem.ActiveDays[index] {
		case actorconst.LeaveTypeShortInactive:
			return "inactive"
		case actorconst.LeaveTypeShortPublicHoliday:
			return "day-off"
		case actorconst.LeaveTypeShortPaid, actorconst.LeaveTypeShortUnpaid, actorconst.LeaveTypeShortSick, actorconst.LeaveTypeShortInjury:
			return "holiday " + atem.ActiveDays[index]
		}
		if index >= 5 {
			return "inactive"
		}
		return ""
	}
	if !atem.isWorkableDay(index) {
		return "error"
	}
	return "active"
}
