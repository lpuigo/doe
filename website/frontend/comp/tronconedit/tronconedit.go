package tronconedit

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/ptedit"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

const template string = `
<div>
	<div v-for="(tr, index) in value.Troncons" :key="index" >
		<hr>		
		<el-row :gutter="10">
			<el-col :span="2">
				<el-button type="danger"
						   plain icon="fas fa-share-alt icon--left"
						   size="mini"
						   style="width: 100%"
						   :disabled="value.Troncons.length<=1"
						   @click="DeleteTroncon(index)"
				>Supprimer</el-button>
			</el-col>
			<el-col :span="22">
				<!-- 
					Attributes about TR 
				-->
				<el-row :gutter="10" type="flex" align="middle">
					<el-col :span="12">
						<el-tooltip content="Référence" placement="top" effect="light">
							<el-input v-model="tr.Ref" placeholder="TR-99-9999" 
									  :readonly="readonly" 
									  clearable size="mini"
							></el-input>
						</el-tooltip>
					</el-col>
					<el-col :span="4">
						<el-switch 
								   active-color="#db2828"
								   active-text="Bloquage"
								   inactive-color="#bcbcbc"
						></el-switch>
					</el-col>
					<el-col :span="4">
						<el-switch v-model="tr.NeedSignature"
								   active-color="#db2828"
								   active-text="Signature demandée"
								   inactive-color="#bcbcbc"
						></el-switch>
					</el-col>
					<el-col :span="4">
						<el-switch v-if="tr.NeedSignature"
								   v-model="tr.Signed"
								   active-color="#51a825"
								   active-text="Signature obtenue"
								   inactive-color="#bcbcbc"
						></el-switch>
					</el-col>
				</el-row>
				<!-- 
					Attributes about PB 
				-->
				<el-row :gutter="10" type="flex" align="middle">
					<el-col :span="16">
						<pt-edit title="PB" v-model="tr.Pb" :readonly="readonly"></pt-edit>
					</el-col>
					<el-col :span="4">
						<el-tooltip content="Nb. EL raccordable" placement="top" effect="light">
							<el-input-number
									v-model="tr.NbRacco"
									:min="0" :max="tr.NbFiber"
									:readonly="readonly"
									size="mini"	controls-position="right" style="width: 100%"
							></el-input-number>
						</el-tooltip>
					</el-col>
					<el-col :span="4">
						<el-tooltip content="Nb. Fibre" placement="top" effect="light">
							<el-input-number v-model="tr.NbFiber"
											 :min="6" :step="6"
											 :readonly="readonly"
											 size="mini" controls-position="right" style="width: 100%">
								<template slot="prepend">Nb Fibre</template>
							</el-input-number>
						</el-tooltip>
					</el-col>
				</el-row>

				<!-- 
					Comment Attributes
				-->	
				<el-row :gutter="10">
					<el-col :span="16">
						<el-input :readonly="readonly" clearable placeholder="Commentaire sur tronçon" size="mini" type="textarea" autosize
								  v-model="tr.Comment"
						></el-input>
					</el-col>
					<el-col :span="4">
						<el-tooltip content="Date Pose PB" placement="top" effect="light">
							<el-date-picker :readonly="readonly" format="dd/MM/yyyy" placeholder="Date Installation" size="mini"
											style="width: 100%" type="date"
											v-model="tr.InstallDate"
											value-format="yyyy-MM-dd"
											:picker-options="{firstDayOfWeek:1, disabledDate(time) { return time.getTime() > Date.now(); }}"
											:clearable="false"
							></el-date-picker>
						</el-tooltip>
					</el-col>
					<el-col :span="4">
						<el-tooltip content="Date Mesure" placement="top" effect="light">
							<el-date-picker :readonly="readonly" format="dd/MM/yyyy" placeholder="Date Mesure" size="mini"
											style="width: 100%" type="date"
											v-model="tr.MeasureDate"
											value-format="yyyy-MM-dd"
											:picker-options="{firstDayOfWeek:1, disabledDate(time) { return time.getTime() > Date.now(); }}"
											:clearable="false"
											:disabled="!tr.InstallDate"
							></el-date-picker>
						</el-tooltip>
					</el-col>
				</el-row>
			</el-col>
		</el-row>
	</div>
</div>
`

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Registration

func Register() {
	hvue.NewComponent("troncon-edit",
		ComponentOptions()...,
	)
}

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("troncon-edit", ComponentOptions()...)
}

func ComponentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		ptedit.RegisterComponent(),
		hvue.Template(template),
		hvue.Props("readonly", "value"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewTronconEditModel(vm)
		}),
		//hvue.Computed("refPH", func(vm *hvue.VM) interface{} {
		//	pem := &TronconEditModel{Object: vm.Object}
		//	return pem.Title + "-99999"
		//}),
		hvue.MethodsOf(&TronconEditModel{}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type TronconEditModel struct {
	*js.Object

	Order    *fm.Order `js:"value"`
	Readonly bool      `js:"readonly"`

	VM *hvue.VM `js:"VM"`
}

func NewTronconEditModel(vm *hvue.VM) *TronconEditModel {
	tem := &TronconEditModel{Object: tools.O()}
	tem.VM = vm
	tem.Order = nil
	tem.Readonly = false
	return tem
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Actions

func (tem *TronconEditModel) DeleteTroncon(vm *hvue.VM, i int) {
	tem = &TronconEditModel{Object: vm.Object}
	tem.Order.DeleteTroncon(i)
}
