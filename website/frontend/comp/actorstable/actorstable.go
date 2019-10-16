package actorstable

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
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
	template string = `
<el-table
        :border=true
        :data="filteredActors"
		:default-sort = "{prop: 'Client', order: 'ascending'}"
        :row-class-name="TableRowClassName" height="100%" size="mini"
>
    <el-table-column
            :resizable="true" :show-overflow-tooltip=true 
            prop="Company" label="Société" width="110px"
			sortable :sort-by="['Company', 'State', 'Role', 'Ref']"
			:filters="FilterList('Company')" :filter-method="FilterHandler"	filter-placement="bottom-end" :filtered-value="FilteredStatusValue()"
    ></el-table-column>
    
    <el-table-column
            :resizable="true" :show-overflow-tooltip=true
            prop="Contract" label="Contrat" width="110px"
    ></el-table-column>
    
    <el-table-column
            :resizable="true" :show-overflow-tooltip=true 
            prop="Client" label="Clients" width="200px"
			sortable :sort-method="SortClient"
			:filters="FilterList('Client')" :filter-method="FilterHandler"	filter-placement="bottom-end" :filtered-value="FilteredStatusValue()"
    >
        <template slot-scope="scope">
			<span>{{GetClients(scope.row)}}</span>
        </template>
	</el-table-column>

    <el-table-column
            :resizable="true" :show-overflow-tooltip=true 
            prop="Role" label="Rôle" width="110px"
			sortable :sort-by="['Role', 'State', 'Ref']"
			:filters="FilterList('Role')" :filter-method="FilterHandler"	filter-placement="bottom-end" :filtered-value="FilteredStatusValue()"
    ></el-table-column>
    
    <el-table-column
            :resizable="true" :show-overflow-tooltip=true 
            prop="Ref" label="Nom Prénom" width="200px"
			sortable :sort-by="['Ref']"
    >
        <template slot-scope="scope">
            <div class="header-menu-container on-hover">
            	<span>{{scope.row.Ref}}</span>
				<i v-if="user.Permissions.Invoice" class="show link fas fa-edit" @click="EditActor(scope.row)"></i>
            </div>
        </template>
	</el-table-column>
    
    <el-table-column
            :resizable="true" :show-overflow-tooltip=true 
            prop="State" label="Statut" width="100px"
			:formatter="FormatState"
			:filters="FilterList('State')" :filter-method="FilterHandler"	filter-placement="bottom-end" :filtered-value="FilteredStatusValue()"
    ></el-table-column>
    
    <el-table-column
            :resizable="true" :show-overflow-tooltip=true 
            label="Congés" width="200px"
    >
        <template slot-scope="scope">
            <div class="header-menu-container on-hover">
            	<span>{{GetHoliday(scope.row)}}</span>
				<i class="show link fas fa-edit" @click="EditActorVacancy(scope.row)"></i>
            </div>
        </template>
    </el-table-column>
    
    <el-table-column
            :resizable="true" prop="Comment" label="Commentaire"
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

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Format & Style Functions

func (atm *ActorsTableModel) TableRowClassName(rowInfo *js.Object) string {
	actor := actor.NewActorFromJS(rowInfo.Get("row"))
	return GetRowStyle(actor)
}

func (atm *ActorsTableModel) GetHoliday(act *actor.Actor) string {
	if act.State == actorconst.StateGone {
		return "parti le " + date.DateString(act.Period.End)
	}
	if act.State == actorconst.StateCandidate {
		return "débute le " + date.DateString(act.Period.Begin)
	}
	vacPeriod := act.GetNextVacation()
	if vacPeriod == nil {
		return ""
	}
	return "du " + date.DateString(vacPeriod.Begin) + " au " + date.DateString(vacPeriod.End)
}

func (atm *ActorsTableModel) GetClients(act *actor.Actor) string {
	return strings.Join(act.Client, ", ")
}

func (atm *ActorsTableModel) SortClient(a, b *actor.Actor) int {
	ca, cb := atm.GetClients(a), atm.GetClients(b)
	switch {
	case ca == cb:
		switch {
		case a.State == b.State:
			return atm.SortRoleRef(a, b)
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

func (atm *ActorsTableModel) SortRoleRef(a, b *actor.Actor) int {
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

func (atm *ActorsTableModel) FormatState(row, column, cellValue, index *js.Object) string {
	return GetStateLabel(cellValue.String())
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Column Filtering Related Methods

func (atm *ActorsTableModel) FilterHandler(value string, p *js.Object, col *js.Object) bool {
	prop := col.Get("property").String()
	return p.Get(prop).String() == value
}

func (atm *ActorsTableModel) FilterList(vm *hvue.VM, prop string) []*elements.ValText {
	atm = ActorsTableModelFromJS(vm.Object)
	count := map[string]int{}
	attribs := []string{}

	var translate func(string) string
	switch prop {
	case "State":
		translate = func(state string) string {
			return GetStateLabel(state)
		}
	default:
		translate = func(val string) string { return val }
	}

	for _, act := range atm.Actors {
		attrib := act.Object.Get(prop).String()
		if _, exist := count[attrib]; !exist {
			attribs = append(attribs, attrib)
		}
		count[attrib]++
	}
	sort.Strings(attribs)
	res := []*elements.ValText{}
	for _, a := range attribs {
		fa := a
		if fa == "" {
			fa = "Vide"
		}
		res = append(res, elements.NewValText(a, translate(fa)+" ("+strconv.Itoa(count[a])+")"))
	}
	return res
}

func (atm *ActorsTableModel) FilteredStatusValue() []string {
	res := []string{
		//poleconst.PsStatusNew,
		//poleconst.PsStatusInProgress,
		//poleconst.PsStatusBlocked,
		//poleconst.PsStatusCancelled,
		//poleconst.PsStatusDone,
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

func (atm *ActorsTableModel) EditActor(vm *hvue.VM, act *actor.Actor) {
	vm.Emit("edit-actor", act)
}

func (atm *ActorsTableModel) EditActorVacancy(vm *hvue.VM, act *actor.Actor) {
	vm.Emit("edit-actor-vacancy", act)
}
