package actorstable

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
	template string = `<el-table
        :border=true
        :data="filteredActors"
        :row-class-name="TableRowClassName" height="100%" size="mini"
>
    <el-table-column
            :resizable="true" :show-overflow-tooltip=true 
            prop="Company" label="Compagnie" width="100px"
    ></el-table-column>
    
    <el-table-column
            :resizable="true" :show-overflow-tooltip=true
            prop="Contract" label="Contrat" width="100px"
    ></el-table-column>
    
    <el-table-column
            :resizable="true" :show-overflow-tooltip=true 
            prop="Ref" label="Nom Prénom" width="100px"
    ></el-table-column>
    
    <el-table-column
            :resizable="true" :show-overflow-tooltip=true 
            prop="State" label="Statut" width="100px"
    ></el-table-column>
    
    <el-table-column
            :resizable="true" :show-overflow-tooltip=true 
            prop="Role" label="Rôle" width="100px"
    ></el-table-column>
    
    <el-table-column
            :resizable="true" :show-overflow-tooltip=true 
            label="Congés" width="200px"
    >
        <template slot-scope="scope">
            <span>{{GetHoliday(scope.row)}}</span>
        </template>
    </el-table-column>
    
    <el-table-column
            label="Etat"
    >
        <template slot-scope="scope">
            <rip-state-update :client="value.Client" :user="user" v-model="scope.row.State"></rip-state-update>
        </template>
    </el-table-column>
    
    <el-table-column
            :resizable="true" :show-overflow-tooltip=true
            prop="Comment" label="Commentaire"
    ></el-table-column>
</el-table>
`
)

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("actors-table", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.Props("value", "user", "filter", "filtertype"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewActorsTableModel(vm)
		}),
		hvue.MethodsOf(&ActorsTableModel{}),
		hvue.Computed("filteredActors", func(vm *hvue.VM) interface{} {
			atm := ActorsTableModelFromJS(vm.Object)
			return atm.GetFilteredActors()
		}),
		//hvue.Filter("FormatTronconRef", func(vm *hvue.VM, value *js.Object, args ...*js.Object) interface{} {
		//	rpum := ActorsTableModelFromJS(vm.Object)
		//	t := &fm.Troncon{Object: value}
		//	return rpum.GetFormatTronconRef(t)
		//}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type ActorsTableModel struct {
	*js.Object

	Actors     []*actor.Actor `js:"value"`
	User       *fm.User       `js:"user"`
	Filter     string         `js:"filter"`
	FilterType string         `js:"filtertype"`

	VM *hvue.VM `js:"VM"`
}

func NewActorsTableModel(vm *hvue.VM) *ActorsTableModel {
	atm := &ActorsTableModel{Object: tools.O()}
	atm.VM = vm
	atm.Actors = []*actor.Actor{}
	atm.User = nil
	atm.Filter = ""
	atm.FilterType = ""
	return atm
}

func ActorsTableModelFromJS(o *js.Object) *ActorsTableModel {
	return &ActorsTableModel{Object: o}
}

func (atm *ActorsTableModel) GetFilteredActors() []*actor.Actor {
	if atm.FilterType == actorconst.FilterValueAll && atm.Filter == "" {
		return atm.Actors
	}
	res := []*actor.Actor{}
	expected := strings.ToUpper(atm.Filter)
	filter := func(a *actor.Actor) bool {
		sis := a.SearchString(atm.FilterType)
		if sis == "" {
			return false
		}
		return strings.Contains(strings.ToUpper(sis), expected)
	}

	for _, actor := range atm.Actors {
		if filter(actor) {
			res = append(res, actor)
		}
	}
	return res
}

func (atm *ActorsTableModel) TableRowClassName(rowInfo *js.Object) string {
	actor := actor.NewActorFromJS(rowInfo.Get("row"))
	return GetRowStyle(actor)
}

func (atm *ActorsTableModel) GetHoliday(act *actor.Actor) string {
	if len(act.Vacation) == 0 {
		return ""
	}
	today := date.TodayAfter(0)
	vacBegin := ""
	vacEnd := ""
	for _, vacPeriod := range act.Vacation {
		if vacPeriod.End < today {
			continue
		}
		if vacBegin == "" && vacPeriod.End >= today {
			vacBegin = vacPeriod.Begin
			vacEnd = vacPeriod.End
			continue
		}
		// vacBegin != ""
		if vacPeriod.Begin < vacBegin {
			vacBegin = vacPeriod.Begin
			vacEnd = vacPeriod.End
		}
	}

	if vacBegin == "" {
		return ""
	}
	return date.DateString(vacBegin) + " au " + date.DateString(vacEnd)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Style Functions

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
