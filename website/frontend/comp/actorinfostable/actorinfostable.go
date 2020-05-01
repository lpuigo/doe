package actorinfostable

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/actorstable"
)

const (
	template string = `
<el-table
        :border=true
        :data="filteredActors"
		:default-sort = "{prop: 'Client', order: 'ascending'}"
        :row-class-name="TableRowClassName" height="100%" size="mini"
		@row-dblclick="HandleDoubleClickedRow"
>
	<el-table-column
			label="N°" width="40px" align="right"
			type="index"
			index=1 
	></el-table-column>

    <el-table-column
            :resizable="true" :show-overflow-tooltip=true 
            prop="Company" label="Société" width="110px"
			sortable :sort-by="['Company', 'State', 'Role', 'Ref']"
			:filters="FilterList('Company')" :filter-method="FilterHandler"	filter-placement="bottom-end"
    ></el-table-column>
    
    <el-table-column
            :resizable="true" :show-overflow-tooltip=true
            prop="Contract" label="Contrat" width="110px"
    ></el-table-column>
    
<!--    <el-table-column-->
<!--            :resizable="true" :show-overflow-tooltip=true -->
<!--            prop="Client" label="Clients" width="200px"-->
<!--			sortable :sort-method="SortClient"-->
<!--			:filters="FilterList('Client')" :filter-method="FilterHandler"	filter-placement="bottom-end"-->
<!--    >-->
<!--        <template slot-scope="scope">-->
<!--			<span>{{GetClients(scope.row)}}</span>-->
<!--        </template>-->
<!--	</el-table-column>-->

    <el-table-column
            :resizable="true" :show-overflow-tooltip=true 
            prop="Role" label="Rôle" width="110px"
			sortable :sort-by="['Role', 'State', 'Ref']"
			:filters="FilterList('Role')" :filter-method="FilterHandler"	filter-placement="bottom-end"
    ></el-table-column>
    
    <el-table-column
            :resizable="true" :show-overflow-tooltip=true 
            prop="Ref" label="Nom Prénom" width="200px"
			sortable :sort-by="['Ref']"
    >
        <template slot-scope="scope">
            <div class="header-menu-container on-hover">
            	<span>{{scope.row.Ref}}</span>
				<i v-if="user.Permissions.HR" class="show link fas fa-edit" @click="EditActor(scope.row)"></i>
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
            label="Info" width="200px"
    >
        <template slot-scope="scope">
			<pre>{{scope.row.Info}}</pre>
        </template>
	</el-table-column>

</el-table>
`
)

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("actorinfos-table", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.Props("value", "user", "filter", "filtertype"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewActorsInfoTableModel(vm)
		}),
		hvue.MethodsOf(&ActorsInfoTableModel{}),
		hvue.Computed("filteredActors", func(vm *hvue.VM) interface{} {
			acm := ActorsInfoTableModelFromJS(vm.Object)
			return acm.GetFilteredActors()
		}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type ActorsInfoTableModel struct {
	*actorstable.ActorsTableModel
}

func NewActorsInfoTableModel(vm *hvue.VM) *ActorsInfoTableModel {
	aitm := &ActorsInfoTableModel{ActorsTableModel: actorstable.NewActorsTableModel(vm)}
	return aitm
}

func ActorsInfoTableModelFromJS(o *js.Object) *ActorsInfoTableModel {
	return &ActorsInfoTableModel{ActorsTableModel: actorstable.ActorsTableModelFromJS(o)}
}
