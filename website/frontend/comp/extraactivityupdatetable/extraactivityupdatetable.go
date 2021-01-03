package extraactivityupdatetable

import (
	"strconv"

	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/model/extraactivity"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
)

const (
	template string = `
<el-table
		:border="false"
		:data="value"
>
	<!--	Index   -->
<!--	<el-table-column-->
<!--			label="N°" width="40px" align="right"-->
<!--			type="index"-->
<!--			index=1 -->
<!--	></el-table-column>-->

	<!--	Actions   -->
	<el-table-column v-if="User.Permissions.Invoice" label="" width="58">
		<template slot="header" slot-scope="scope">
			<el-button type="success" plain icon="fas fa-plus fa-fw" size="mini" @click="AddNewEA()"></el-button>
		</template>
		<template slot-scope="scope">
			<el-button type="danger" plain icon="fas fa-ban fa-fw" size="mini" @click="RemoveEA(scope.$index)"></el-button>
		</template>
	</el-table-column>
	
	<!--	Name   -->
	<el-table-column
			:resizable="true"
			prop="Name" label="Activité" width="350px"
	>		
		<template slot-scope="scope">
            <el-input  v-if="User.Permissions.Invoice" placeholder="Intitulé"
                      v-model="scope.row.Name" clearable size="mini"
            ></el-input>
			<span v-else>{{scope.row.Name}}</span>
		</template>
	</el-table-column>

	<!--	Nb Point   -->
	<el-table-column
			:resizable="true"
			prop="NbPoints" label="Points" width="140px"
	>
		<template slot-scope="scope">
			<el-input-number v-if="User.Permissions.Invoice"
				v-model="scope.row.NbPoints" 
				controls-position="right" :precision="1" :step="5" :min="0"
				size="mini"
			></el-input-number>
			<span v-else>{{scope.row.NbPoints}}</span>
		</template>
	</el-table-column>

	<!--	Income   -->
	<el-table-column v-if="User.Permissions.Invoice"
			:resizable="true" width="140px"
			prop="Income" label="Montant €"
	>
		<template slot-scope="scope">
			<el-input-number
				v-model="scope.row.Income" 
				controls-position="right" :precision="2" :step="50" :min="0"
				size="mini"
			></el-input-number>
		</template>
	</el-table-column>

	<!--	Actors   -->
	<el-table-column
			width="270px"
			prop="Actors" label="Acteurs"
	>
		<template slot-scope="scope">
			<el-select v-model="scope.row.Actors" multiple placeholder="Acteurs" size="mini" style="width: 100%"
					   @clear=""
					   @change="UpdateActors(scope.row)"
			>
				<el-option
						v-for="item in GetActors()"
						:key="item.value"
						:label="item.label"
						:value="item.value"
						:disabled="item.disabled"
				>
				</el-option>
			</el-select>
		</template>
	</el-table-column>

	<!--	Date   -->
	<el-table-column
			width="150px"
			prop="Date" label="Date"
	>
		<template slot-scope="scope">
			<el-date-picker format="dd/MM/yyyy" placeholder="Date" size="mini"
							style="width: 100%" type="date"
							v-model="scope.row.Date"
							value-format="yyyy-MM-dd"
							:picker-options="{firstDayOfWeek:1, disabledDate(time) { return time.getTime() > Date.now(); }}"
			></el-date-picker>
		</template>
	</el-table-column>

	<!-- Comment -->
	<el-table-column
			label="Commentaire" prop="Comment"
			min-width="120px" :resizable=true
	>
		<template slot-scope="scope">
			<el-input type="textarea" :autosize="{ minRows: 2, maxRows: 5}" placeholder="Commentaire"
					  v-model="scope.row.Comment" clearable size="mini"
			></el-input>
		</template>
	</el-table-column>
</el-table>
`
)

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("extra-activities-table", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.Props("value", "User", "Client"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewExtraActivityTableModel(vm)
		}),
		hvue.MethodsOf(&ExtraActivityTableModel{}),
		//hvue.Computed("filteredActors", func(vm *hvue.VM) interface{} {
		//	eatm := ExtraActivityTableModelFromJS(vm.Object)
		//	return eatm.User
		//}),
		//hvue.Filter("FormatTronconRef", func(vm *hvue.VM, value *js.Object, args ...*js.Object) interface{} {
		//	eatm := ExtraActivityTableModelFromJS(vm.Object)
		//	t := &fm.Troncon{Object: value}
		//	return rpum.GetFormatTronconRef(t)
		//}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type ExtraActivityTableModel struct {
	*js.Object

	ExtraActivities []*extraactivity.ExtraActivity `js:"value"`
	User            *fm.User                       `js:"User"`
	Client          string                         `js:"Client"`

	VM *hvue.VM `js:"VM"`
}

func NewExtraActivityTableModel(vm *hvue.VM) *ExtraActivityTableModel {
	eatm := &ExtraActivityTableModel{Object: tools.O()}
	eatm.VM = vm
	eatm.ExtraActivities = []*extraactivity.ExtraActivity{}
	eatm.User = fm.NewUser()
	eatm.Client = ""
	return eatm
}

func ExtraActivityTableModelFromJS(o *js.Object) *ExtraActivityTableModel {
	return &ExtraActivityTableModel{Object: o}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Format & Style Functions

func (eatm *ExtraActivityTableModel) AddNewEA(vm *hvue.VM) {
	eatm = ExtraActivityTableModelFromJS(vm.Object)
	eatm.Object.Get("value").Call("push", extraactivity.NewExtraActivity())
}

func (eatm *ExtraActivityTableModel) RemoveEA(vm *hvue.VM, index int) {
	eatm = ExtraActivityTableModelFromJS(vm.Object)
	eatm.Object.Get("value").Call("splice", index, 1)
}

func (eatm *ExtraActivityTableModel) UpdateActors(vm *hvue.VM, ea *extraactivity.ExtraActivity) {
	eatm = ExtraActivityTableModelFromJS(vm.Object)
	client := eatm.User.GetClientByName(eatm.Client)
	if client == nil {
		return
	}
	actors := make(map[string]string)
	for _, actor := range client.Actors {
		actors[strconv.Itoa(actor.Id)] = actor.GetRef()
	}
	ea.Get("Actors").Call("sort", func(a, b string) int {
		// check if actors are not known
		if actors[a] == "" && actors[b] == "" {
			return 0
		}
		if !(actors[a] != "" && actors[b] != "") {
			return 1
		}
		// compare known actors
		if actors[a] < actors[b] {
			return -1
		}
		return 1
	})
}

func (eatm *ExtraActivityTableModel) GetActors(vm *hvue.VM) []*elements.ValueLabelDisabled {
	eatm = ExtraActivityTableModelFromJS(vm.Object)
	client := eatm.User.GetClientByName(eatm.Client)
	if client == nil {
		return nil
	}
	res := []*elements.ValueLabelDisabled{}
	for _, actor := range client.Actors {
		res = append(res, actor.GetElementsValueLabelDisabled())
	}
	return res
}
