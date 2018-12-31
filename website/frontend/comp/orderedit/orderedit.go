package orderedit

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/tronconedit"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

const template string = `
<div>
	<!-- 
		 Order attributes
	 -->
	<el-row :gutter="10">
		<el-col :span="5">
			<el-input placeholder="F99999jjmmaa"
                      v-model="value.Ref"
                      :readonly="readonly" clearable size="mini"
			>
                <template slot="prepend">Commande:</template>
            </el-input>
		</el-col>
		<el-col :span="19">
			<el-input placeholder="Commentaire sur la commande" :readonly="readonly" clearable size="mini" type="textarea" autosize
					  v-model="value.Comment"
			></el-input>
		</el-col>
	</el-row>
	<!-- 
		 Attributes about Order.Troncons 
	 -->
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
                <troncon-edit v-model="tr" :previous="PreviousTroncon(index)" :readonly="readonly"></troncon-edit>
            </el-col>
        </el-row>
    </div>
	<hr>
	<el-row :gutter="10" type="flex" align="middle">
		<el-col :span="2">
			<el-button type="primary" 
					   plain icon="fas fa-share-alt icon--left" 
					   size="mini" style="width: 100%"
					   @click="AddTroncon()"
			>Ajouter</el-button>
		</el-col>
	</el-row>
</div>`

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Registration

func Register() {
	hvue.NewComponent("order-edit",
		ComponentOptions()...,
	)
}

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("order-edit", ComponentOptions()...)
}

func ComponentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		tronconedit.RegisterComponent(),
		hvue.Template(template),
		hvue.Props("readonly", "value"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewOrderEditModel(vm)
		}),
		//hvue.Computed("refPH", func(vm *hvue.VM) interface{} {
		//	pem := &OrderEditModel{Object: vm.Object}
		//	return pem.Title + "-99999"
		//}),
		hvue.MethodsOf(&OrderEditModel{}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type OrderEditModel struct {
	*js.Object

	Order    *fm.Order `js:"value"`
	Readonly bool      `js:"readonly"`

	VM *hvue.VM `js:"VM"`
}

func NewOrderEditModel(vm *hvue.VM) *OrderEditModel {
	oem := &OrderEditModel{Object: tools.O()}
	oem.VM = vm
	oem.Order = nil
	oem.Readonly = false
	return oem
}

func (oem *OrderEditModel) AddTroncon(vm *hvue.VM) {
	oem = &OrderEditModel{Object: vm.Object}
	oem.Order.AddTroncon()
}

func (oem *OrderEditModel) DeleteTroncon(vm *hvue.VM, i int) {
	oem = &OrderEditModel{Object: vm.Object}
	oem.Order.DeleteTroncon(i)
}

func (oem *OrderEditModel) PreviousTroncon(vm *hvue.VM, i int) *fm.Troncon {
	oem = &OrderEditModel{Object: vm.Object}
	if len(oem.Order.Troncons) <= 1 || i == 0 {
		return nil
	}
	return oem.Order.Troncons[i-1]
}
