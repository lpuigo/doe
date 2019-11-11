package actortimeedit

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/model/timesheet"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

const (
	template string = `
	<div class="calendar-row">
		<div v-for="(hour, index) in times.Hours"
			:key="index"
			class="calendar-slot"
 			:class="GetColumnColor(index)"
			style="padding: 2px 0px"
		>
        <el-row :gutter="10" type="flex" align="middle">
			<el-col :offset="3" :span="6">
				<el-button type="primary" plain icon="fas fa-calendar-plus" size="mini" @click="SetFullDay(index)"></el-button>
			</el-col>
            <el-col :span="12">
				<el-input-number v-model="times.Hours[index]"
								size="mini" controls-position="right" style="width: 100%"
								@change="HandleChange"
								:min="0" :max="11"
				></el-input-number>
            </el-col>
        </el-row>
		</div>
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
		hvue.Props("times"),
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

	Times *timesheet.ActorsTime `js:"times"`

	VM *hvue.VM `js:"VM"`
}

func NewActorTimeEditModel(vm *hvue.VM) *ActorTimeEditModel {
	atem := &ActorTimeEditModel{Object: tools.O()}
	atem.Times = timesheet.NewActorTime()
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

func (atem *ActorTimeEditModel) GetColumnColor(vm *hvue.VM, index int) string {
	atem = ActorTimeEditModelFromJS(vm.Object)
	if atem.Times.Hours[index] == 0 {
		if index >= 5 {
			return "inactive"
		}
		return ""
	}
	return "active"
}
