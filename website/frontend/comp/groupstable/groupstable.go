package groupstable

import (
	"sort"
	"strconv"
	"strings"

	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/model/actor"
	"github.com/lpuig/ewin/doe/website/frontend/model/group"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/date"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements/message"
)

const (
	template string = `
<el-table
        :border=false
        :data="filteredGroups"
		:default-sort = "{prop: 'Name', order: 'ascending'}"
        :row-class-name="TableRowClassName" height="100%" size="mini"
>
	<!--	Index   -->
	<el-table-column
			label="N°" width="40px" align="right"
			type="index"
			index=1 
	></el-table-column>

	<!--	Actions   -->
	<el-table-column label="" width="80">
		<template slot="header" slot-scope="scope">
			<el-button type="success" plain icon="fas fa-users fa-fw" size="mini" @click="AddNewGroup()"></el-button>
		</template>
	</el-table-column>
	
	<!--	Nb Active Actors   -->
	<el-table-column label="Nb Acteurs" width="80">
    	<template slot-scope="scope">
			<span>{{activeActors[scope.row.Id]}}</span>
        </template>
	</el-table-column>
	
	<!--	Groupe   -->
    <el-table-column
            :resizable="true" :show-overflow-tooltip=true 
            prop="Name" label="Groupe" width="210px"
			sortable :sort-by="['Name']"
    >        
    	<template slot-scope="scope">
			<el-input 
				v-model="scope.row.Name" 
				size="mini"
			></el-input>
        </template>
	</el-table-column>
	<!-- :filters="FilterList('Name')" :filter-method="FilterHandler"	filter-placement="bottom-end"-->

	<!--	clients -->   
    <el-table-column
            :resizable="true" :show-overflow-tooltip=true 
            prop="Clients" label="Clients" width="600px"
			:filters="FilterList('Client')" :filter-method="FilterHandler"	filter-placement="bottom-end"
    >
        <template slot-scope="scope">
			<el-select multiple placeholder="Client" size="mini"
					   v-model="scope.row.Clients"
					   style="width: 100%"
			>
				<el-option v-for="item in GetClientList()"
						   :key="item.value"
						   :label="item.label"
						   :value="item.value"
				>
				</el-option>
			</el-select>
        </template>
	</el-table-column>

	<!--	ActorDailyWork   -->
    <el-table-column
            :resizable="true" :show-overflow-tooltip=true 
            prop="ActorDailyWork" label="Points / J" width="160px"
    >
        <template slot-scope="scope">
			<el-input-number 
				v-model="scope.row.ActorDailyWork" 
				:precision="1" :step="5" :min="0"
				size="mini"
			></el-input-number>
        </template>
	</el-table-column>

	<!--	ActorDailyIncome   -->
    <el-table-column
            :resizable="true" :show-overflow-tooltip=true 
            prop="ActorDailyIncome" label="€ / J" width="160px"
    >
        <template slot-scope="scope">
			<el-input-number 
				v-model="scope.row.ActorDailyIncome" 
				:precision="1" :step="25" :min="0"
				size="mini"
			></el-input-number>
        </template>
	</el-table-column>
</el-table>
`
)

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("groups-table", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.Props("value", "actors", "user"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewGroupsTableModel(vm)
		}),
		hvue.MethodsOf(&GroupsTableModel{}),
		hvue.Computed("filteredGroups", func(vm *hvue.VM) interface{} {
			gtm := GroupsTableModelFromJS(vm.Object)
			return gtm.GroupStore.Groups
		}),
		hvue.Computed("activeActors", func(vm *hvue.VM) interface{} {
			gtm := GroupsTableModelFromJS(vm.Object)
			return gtm.SetActiveActors()
		}),
		//hvue.Filter("FormatTronconRef", func(vm *hvue.VM, value *js.Object, args ...*js.Object) interface{} {
		//	rpum := GroupsTableModelFromJS(vm.Object)
		//	t := &fm.Troncon{Object: value}
		//	return rpum.GetFormatTronconRef(t)
		//}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type GroupsTableModel struct {
	*js.Object

	Actors       []*actor.Actor    `js:"actors"`
	GroupStore   *group.GroupStore `js:"value"`
	User         *fm.User          `js:"user"`
	ActiveActors map[int]string    `js:"ActiveActors"`

	VM *hvue.VM `js:"VM"`
}

func NewGroupsTableModel(vm *hvue.VM) *GroupsTableModel {
	gtm := &GroupsTableModel{Object: tools.O()}
	gtm.VM = vm
	gtm.Actors = []*actor.Actor{}
	gtm.GroupStore = group.NewGroupStore()
	gtm.User = fm.NewUser()
	gtm.ActiveActors = make(map[int]string)
	return gtm
}

func GroupsTableModelFromJS(o *js.Object) *GroupsTableModel {
	return &GroupsTableModel{Object: o}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Format & Style Functions

func (gtm *GroupsTableModel) TableRowClassName(rowInfo *js.Object) string {
	//group := group.GroupFromJS(rowInfo.Get("row"))
	//return GetRowStyle(group)
	return ""
}

func (gtm *GroupsTableModel) GetClientList(vm *hvue.VM) []*elements.ValueLabel {
	gtm = GroupsTableModelFromJS(vm.Object)
	res := []*elements.ValueLabel{}
	for _, clientName := range gtm.User.GetSortedClientNames() {
		res = append(res, elements.NewValueLabel(clientName, clientName))
	}
	return res
}

func (gtm *GroupsTableModel) SortGroup(vm *hvue.VM, a, b *group.Group) int {
	switch {
	case a.Name == b.Name:
		return 0
	case a.Name < b.Name:
		return -1
	default:
		return 1
	}
}

func (gtm *GroupsTableModel) FormatDate(d string) string {
	return date.DateString(d)
}

func (gtm *GroupsTableModel) SortDate(attrib1, attrib2 string) func(obj *js.Object) string {
	return func(obj *js.Object) string {
		val := obj.Get(attrib1).Get(attrib2).String()
		if val == "" {
			return "9999-12-31"
		}
		return val
	}
}

func (gtm *GroupsTableModel) SetActiveActors() map[int]string {
	aa := make(map[int]int)
	for _, act := range gtm.Actors {
		if !act.IsActive() {
			continue
		}
		grpId, _ := act.Groups.GetCurrentInfo()
		aa[grpId]++
	}
	activAct := make(map[int]string)
	for grbId, nb := range aa {
		val := strconv.Itoa(nb)
		if nb > 1 {
			val += " acteurs"
		} else {
			val += " acteur"
		}
		activAct[grbId] = val
	}
	gtm.ActiveActors = activAct
	return gtm.ActiveActors
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Column Filtering Related Methods

func (gtm *GroupsTableModel) FilterHandler(vm *hvue.VM, value string, p *js.Object, col *js.Object) bool {
	prop := col.Get("property").String()
	grp := group.GroupFromJS(p)
	switch prop {
	case "Client":
		for _, c := range grp.Clients {
			if c == value {
				return true
			}
		}
		return false
	}
	return p.Get(prop).String() == value
}

func (gtm *GroupsTableModel) FilterList(vm *hvue.VM, prop string) []*elements.ValText {
	gtm = GroupsTableModelFromJS(vm.Object)
	count := map[string]int{}
	attribs := []string{}

	var translate func(string) string
	//switch prop {
	//case "State":
	//	translate = func(state string) string {
	//		return GetStateLabel(state)
	//	}
	//default:
	translate = func(val string) string { return val }
	//}

	for _, grp := range gtm.GroupStore.Groups {
		var attrs []string
		switch prop {
		case "Client":
			attrs = grp.Clients
		default:
			attrs = []string{grp.Object.Get(prop).String()}
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
		res = append(res, elements.NewValText(a, translate(fa)+" ("+strconv.Itoa(count[a])+")"))
	}
	return res
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Tools Functions

//func GetRowStyle(group *group.Group) string {
//	return ""
//}
//
//func GetStateLabel(state string) string {
//	switch state {
//	case actorconst.StateCandidate:
//		return actorconst.StateLabelCandidate
//	case actorconst.StateActive:
//		return actorconst.StateLabelActive
//	case actorconst.StateOnHoliday:
//		return actorconst.StateLabelOnHoliday
//	case actorconst.StateGone:
//		return actorconst.StateLabelGone
//	case actorconst.StateDefection:
//		return actorconst.StateLabelDefection
//	default:
//		return "Erreur"
//	}
//}

func (gtm *GroupsTableModel) GetClients(grp *group.Group) string {
	return strings.Join(grp.Clients, ", ")
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Action Methods

func (gtm *GroupsTableModel) AddNewGroup(vm *hvue.VM) {
	gtm = GroupsTableModelFromJS(vm.Object)
	ngrp := group.NewGroup()
	ngrp.Name = "Nouveau Groupe"
	ngrp.ActorDailyIncome = 600.0
	ngrp.ActorDailyWork = 30.0
	gtm.GroupStore.AddNewGroup(ngrp)
}

func (gtm *GroupsTableModel) RemoveGroup(vm *hvue.VM, groupId int) {
	dgrp := gtm.GroupStore.GetGroupById(groupId)
	if dgrp == nil {
		message.NotifyWarning(vm, "Suppression de groupe", "Impossible de supprimer le groupe avec l'Id "+strconv.Itoa(groupId))
		return
	}
	print("Suppression du groupe", groupId, dgrp.Name)
	gtm.GroupStore.RemoveGroupById(groupId)
}

func (gtm *GroupsTableModel) EditGroup(vm *hvue.VM, group *group.Group) {
	vm.Emit("edit-group", group)
}
