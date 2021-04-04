package vehiculestable

import (
	"github.com/lpuig/ewin/doe/website/frontend/model/actor"
	"sort"
	"strconv"
	"strings"

	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/model/vehicule"
	"github.com/lpuig/ewin/doe/website/frontend/model/vehicule/vehiculeconst"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/date"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
)

const (
	template string = `
<el-table
        :border=true
        :data="filteredVehicules"
		:default-sort = "{prop: 'Immat', order: 'ascending'}"
        :row-class-name="TableRowClassName" height="100%" size="mini"
		@row-dblclick="HandleDoubleClickedRow"
>
	<!--	Index   -->
	<el-table-column
			label="N°" width="40px" align="right"
			type="index"
			index=1 
	></el-table-column>

	<!--	Status   -->
    <el-table-column
            :resizable="true" :show-overflow-tooltip=true 
            label="Statut" width="90px"
			sortable :sort-method="CompStatus"
			:filters="FilterList('Status')" :filter-method="FilterHandler"	filter-placement="bottom-end"
    >
		<template slot-scope="scope">{{GetStatus(scope.row)}}</template>
	</el-table-column>
    
	<!--	Compagny   -->
    <el-table-column
            :resizable="true" :show-overflow-tooltip=true 
            prop="Company" label="Société" width="110px"
			sortable :sort-by="['Company', 'Type', 'Immat']"
			:filters="FilterList('Company')" :filter-method="FilterHandler"	filter-placement="bottom-end"
    ></el-table-column>
    
	<!--	Type   -->
    <el-table-column
            :resizable="true" :show-overflow-tooltip=true 
            prop="Type" label="Type" width="110px"
			sortable :sort-by="['Type', 'Immat']"
			:filters="FilterList('Type')" :filter-method="FilterHandler"	filter-placement="bottom-end"
    ></el-table-column>
        
	<!--	Immat   -->
    <el-table-column
            :resizable="true" :show-overflow-tooltip=true
            prop="Immat" label="Immatriculation" width="130px"
			sortable
    >
		<template slot-scope="scope">
            <div class="header-menu-container on-hover">
            	<span>{{scope.row.Immat}}</span>
				<i class="show link fas fa-edit" @click="EditVehicule(scope.row)"></i>
            </div>
        </template>
	</el-table-column>
    
	<!--	group   -->
<!--    <el-table-column-->
<!--            :resizable="true" :show-overflow-tooltip=true -->
<!--            prop="Groups" label="Groupe" width="150px"-->
<!--			sortable :sort-method="SortGroup"-->
<!--			:filters="FilterList('Groups')" :filter-method="FilterHandler"	filter-placement="bottom-end"-->
<!--    >-->
<!--        <template slot-scope="scope">-->
<!--			<span>{{GetGroup(scope.row)}}</span>-->
<!--        </template>-->
<!--	</el-table-column>-->

<!--	&lt;!&ndash;	Last & First Name   &ndash;&gt;-->
<!--    <el-table-column-->
<!--            :resizable="true" :show-overflow-tooltip=true -->
<!--            prop="Ref" label="Nom Prénom" width="200px"-->
<!--			sortable :sort-by="['Ref']"-->
<!--    >-->
<!--        <template slot-scope="scope">-->
<!--            <div class="header-menu-container on-hover">-->
<!--            	<span>{{scope.row.Ref}}</span>-->
<!--				<i v-if="user.Permissions.HR" class="show link fas fa-edit" @click="EditVehicule(scope.row)"></i>-->
<!--            </div>-->
<!--        </template>-->
<!--	</el-table-column>-->
    
	<!--	InCharge   -->
    <el-table-column
            label="Responsable"
            width="140px" :resizable="true"
    >
		<template slot-scope="scope">
			<span>{{InChargeName(scope.row)}}</span>
		</template>
    </el-table-column>
        
	<!--	Start Day   -->
    <el-table-column
            label="Mise en Service" sortable :sort-by="SortDate('ServiceDate')"
            width="130px" :resizable="true" 
			align="center"
    >
		<template slot-scope="scope">
			<span>{{FormatServiceDate(scope.row)}}</span>
		</template>
    </el-table-column>
        
	<!--	Carte Carb.   -->
    <el-table-column
            label="Carte Carb." prop="FuelCard" 
            width="130px" :resizable="true" 
			align="center"
    ></el-table-column>
        
	<!--	Kilometrage   -->
    <el-table-column
            label="Kilométrage"
            width="145px" :resizable="true"
    >
		<template slot-scope="scope">
			<span>{{TravelledKms(scope.row)}}</span>
		</template>
    </el-table-column>
        
	<!--	Inventaire   -->
    <el-table-column
            label="Inventaire"
            width="90px" :resizable="true"
    >
		<template slot-scope="scope">
			<span>{{LastInventory(scope.row)}}</span>
		</template>
    </el-table-column>
        
	<!--	Event   -->
    <el-table-column
            label="Évènements"
            width="140px" :resizable="true"
    >
		<template slot-scope="scope">
			<span v-html="LastEvent(scope.row)"></span>
		</template>
    </el-table-column>
        
	<!--	Comment   -->
    <el-table-column
            :resizable="true" prop="Comment" label="Commentaire"
    ></el-table-column>
</el-table>
`
)

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("vehicules-table", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.Props("value", "user", "actorstore", "filter", "filtertype"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewVehiculesTableModel(vm)
		}),
		hvue.MethodsOf(&VehiculesTableModel{}),
		hvue.Computed("filteredVehicules", func(vm *hvue.VM) interface{} {
			vtm := VehiculesTableModelFromJS(vm.Object)
			return vtm.GetFilteredVehicules()
		}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type VehiculesTableModel struct {
	*js.Object

	Vehicules  []*vehicule.Vehicule `js:"value"`
	User       *fm.User             `js:"user"`
	ActorStr   *actor.ActorStore    `js:"actorstore"`
	Filter     string               `js:"filter"`
	FilterType string               `js:"filtertype"`

	VM *hvue.VM `js:"VM"`
}

func NewVehiculesTableModel(vm *hvue.VM) *VehiculesTableModel {
	atm := &VehiculesTableModel{Object: tools.O()}
	atm.VM = vm
	atm.Vehicules = []*vehicule.Vehicule{}
	atm.User = fm.NewUser()
	atm.ActorStr = actor.NewActorStore()
	atm.Filter = ""
	atm.FilterType = ""
	return atm
}

func VehiculesTableModelFromJS(o *js.Object) *VehiculesTableModel {
	return &VehiculesTableModel{Object: o}
}

func (vtm *VehiculesTableModel) GetFilteredVehicules() []*vehicule.Vehicule {
	if vtm.FilterType == vehiculeconst.FilterValueAll && vtm.Filter == "" {
		return vtm.Vehicules
	}
	res := []*vehicule.Vehicule{}
	expected := strings.ToUpper(vtm.Filter)
	filter := func(v *vehicule.Vehicule) bool {
		sis := v.SearchString(vtm.FilterType)
		if sis == "" {
			return false
		}
		return strings.Contains(strings.ToUpper(sis), expected)
	}

	for _, vehic := range vtm.Vehicules {
		if filter(vehic) {
			res = append(res, vehic)
		}
	}
	return res
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Format & Style Functions

func (vtm *VehiculesTableModel) TableRowClassName(rowInfo *js.Object) string {
	vehic := vehicule.VehiculeFromJS(rowInfo.Get("row"))
	return GetRowStyle(vehic)
}

//func (vtm *VehiculesTableModel) GetClients(act *actor.Actor) string {
//	return strings.Join(act.Client, ", ")
//}
//
//func (vtm *VehiculesTableModel) SortClient(a, b *actor.Actor) int {
//	ca, cb := vtm.GetClients(a), vtm.GetClients(b)
//	switch {
//	case ca == cb:
//		return vtm.SortRoleRef(a, b)
//	case ca < cb:
//		return -1
//	default:
//		return 1
//	}
//}

func (vtm *VehiculesTableModel) FormatDate(d string) string {
	return date.DateString(d)
}

func (vtm *VehiculesTableModel) FormatServiceDate(vehic *vehicule.Vehicule) string {
	if tools.Empty(vehic.EndServiceDate) {
		return date.DateString(vehic.ServiceDate)
	}
	return date.DateString(vehic.ServiceDate) + " à " + date.DateString(vehic.EndServiceDate)
}

func (vtm *VehiculesTableModel) InChargeName(vm *hvue.VM, vehic *vehicule.Vehicule) string {
	vtm = VehiculesTableModelFromJS(vm.Object)
	actId := vehic.GetInChargeActorId(date.TodayAfter(0))
	act := vtm.ActorStr.GetActorById(actId)
	if act == nil {
		return vehiculeconst.InChargeNotAffected
	}
	return act.GetRefStatus()
}

func (vtm *VehiculesTableModel) TravelledKms(vm *hvue.VM, vehic *vehicule.Vehicule) string {
	vtm = VehiculesTableModelFromJS(vm.Object)
	cth := vehic.GetCurrentTravelledKms()
	if cth == nil {
		return "Non connu"
	}
	return date.DateString(cth.Date) + " : " + strconv.Itoa(cth.Kms) + " kms"
}

func (vtm *VehiculesTableModel) SortDate(attrib1 string) func(obj *js.Object) string {
	return func(obj *js.Object) string {
		val := obj.Get(attrib1).String()
		if val == "" {
			return "9999-12-31"
		}
		return val
	}
}

func (vtm *VehiculesTableModel) GetStatus(vehic *vehicule.Vehicule) string {
	return vehic.Status()
}

func (vtm *VehiculesTableModel) CompStatus(a, b *vehicule.Vehicule) int {
	sa, sb := a.Status(), b.Status()
	if sa == sb {
		if a.Immat < b.Immat {
			return -1
		}
		return 1
	}
	if sa < sb {
		return -1
	}
	return 1
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Column Filtering Related Methods

func (vtm *VehiculesTableModel) FilterHandler(vm *hvue.VM, value string, p *js.Object, col *js.Object) bool {
	prop := col.Get("property").String()
	if prop == "Status" {
		return vehicule.VehiculeFromJS(p).Status() == value
	}
	return p.Get(prop).String() == value
}

func (vtm *VehiculesTableModel) FilterList(vm *hvue.VM, prop string) []*elements.ValText {
	vtm = VehiculesTableModelFromJS(vm.Object)
	count := map[string]int{}
	attribs := []string{}

	var getValue func(vehic *vehicule.Vehicule) string
	switch prop {
	case "Status":
		getValue = func(vehic *vehicule.Vehicule) string {
			return vehic.Status()
		}
	default:
		getValue = func(vehic *vehicule.Vehicule) string {
			return vehic.Object.Get(prop).String()
		}
	}

	attrib := ""
	for _, vehic := range vtm.Vehicules {
		attrib = getValue(vehic)
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
		res = append(res, elements.NewValText(a, fa+" ("+strconv.Itoa(count[a])+")"))
	}
	return res
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Tools Functions

func GetRowStyle(vehic *vehicule.Vehicule) string {
	class := ""
	if vehic.Status() == vehiculeconst.StatusReturned {
		class = "vehicule-returned "
	}
	switch vehic.Type {
	case vehiculeconst.TypeTariere:
		class += "vehicule-row-tariere"
	case vehiculeconst.TypeNacelle:
		class += "vehicule-row-nacelle"
	case vehiculeconst.TypeFourgon:
		class += "vehicule-row-fourgon"
	case vehiculeconst.TypeCar:
		class += "vehicule-row-car"
	case vehiculeconst.TypePorteTouret:
		class += "vehicule-row-portetouret"
	default:
		class += "vehicule-row-error"
	}
	return class
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Action Methods

func (vtm *VehiculesTableModel) EditVehicule(vm *hvue.VM, vehic *vehicule.Vehicule) {
	vm.Emit("edit-vehicule", vehic)
}

func (vtm *VehiculesTableModel) HandleDoubleClickedRow(vm *hvue.VM, vehic *vehicule.Vehicule) {
	vtm = VehiculesTableModelFromJS(vm.Object)
	vtm.EditVehicule(vm, vehic)
}

func (vtm *VehiculesTableModel) LastInventory(vm *hvue.VM, vehic *vehicule.Vehicule) string {
	vtm = VehiculesTableModelFromJS(vm.Object)
	if len(vehic.Inventories) == 0 {
		return "A Faire"
	}
	lastInvent := vehic.Inventories[0]
	return date.DateString(lastInvent.ReferenceDate)
}

func (vtm *VehiculesTableModel) LastEvent(vm *hvue.VM, vehic *vehicule.Vehicule) string {
	vtm = VehiculesTableModelFromJS(vm.Object)
	event, future := vehic.GetInterestEvent()
	if event == nil {
		return ""
	}
	res := date.DateString(event.StartDate) + " " + event.Type
	if future {
		res = `<strong style="color: blue;">` + res + "</strong>"
	}
	return res
}
