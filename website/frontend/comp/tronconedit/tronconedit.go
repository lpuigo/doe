package tronconedit

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/ptedit"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

const template string = `
<el-row :gutter="10">
    <el-col :span="2">
        <el-button type="danger" plain icon="fas fa-share-alt icon--left" size="mini" style="width: 100%">Supprimer</el-button>
    </el-col>
    <el-col :span="22">
        <!-- 
               Attributes about each Troncon 
        -->
        <!-- 
            Attributes about PB 
        -->
        <el-row :gutter="10" type="flex" align="middle">
            <el-col :span="20">
				<pt-edit title="PB" v-model="value.Pb" :readonly="readonly"></pt-edit>
            </el-col>
			<el-col :span="4">
				<el-switch v-model="value.NeedSignature"
						   active-color="#ff3200"
						   active-text="Signature demandée"
						   inactive-color="#bcbcbc">
				</el-switch>
			</el-col>
        </el-row>
        <!-- 
            Attributes about TR 
        -->
        <el-row :gutter="10" type="flex" align="middle">
            <el-col :span="6">
                <el-input placeholder="TR-99-9999" :readonly="readonly" clearable size="mini"
                          v-model="value.Ref"
                ></el-input>
            </el-col>
            <el-col :span="3">
                <el-input-number v-model="value.NbRacco" :min="0" :max="value.NbFiber" :readonly="readonly" size="mini" label="Nb Racco" controls-position="right" style="width: 100%">
                    <template slot="prepend">Nb El</template>
                </el-input-number>
            </el-col>
            <el-col :span="3">
                <el-input-number v-model="value.NbFiber" :min="6" :step="6" :readonly="readonly" size="mini" label="Nb Fibre" controls-position="right" style="width: 100%">
                    <template slot="prepend">Nb Fibre</template>
                </el-input-number>
            </el-col>
            <el-col :span="4">
                <el-date-picker :readonly="readonly" format="dd/MM/yyyy" placeholder="Date Installation" size="mini"
                                style="width: 100%" type="date"
                                v-model="value.InstallDate"
                                value-format="yyyy-MM-dd"
                                :picker-options="{firstDayOfWeek:1, disabledDate(time) { return time.getTime() > Date.now(); }}"
                                :clearable="false"
                ></el-date-picker>
            </el-col>
            <el-col :span="4">
                <el-date-picker :readonly="readonly" format="dd/MM/yyyy" placeholder="Date Mesure" size="mini"
                                style="width: 100%" type="date"
                                v-model="value.MeasureDate"
                                value-format="yyyy-MM-dd"
                                :picker-options="{firstDayOfWeek:1, disabledDate(time) { return time.getTime() > Date.now(); }}"
                                :clearable="false"
								:disabled="!value.InstallDate"
                ></el-date-picker>
            </el-col>
			<el-col v-if="value.NeedSignature" :span="4">
				<el-switch v-model="value.Signed"
						   active-color="#51a825"
						   active-text="Signature obtenue"
						   inactive-color="#bcbcbc"
				></el-switch>
			</el-col>
        </el-row>
		<el-row :gutter="10">
			<el-col :span="24">
				<el-input :readonly="readonly" clearable placeholder="Commentaire sur tronçon" size="mini" type="textarea" autosize
						  v-model="value.Comment"
				></el-input>
			</el-col>
		</el-row>

	</el-col>  
</el-row>
`

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Registration

func Register() {
	hvue.NewComponent("troncon-edit",
		ComponentOptions()...,
	)
}

func ComponentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Component("pt-edit", ptedit.ComponentOptions()...),
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

	Troncon  *fm.Troncon `js:"value"`
	Readonly bool        `js:"readonly"`

	VM *hvue.VM `js:"VM"`
}

func NewTronconEditModel(vm *hvue.VM) *TronconEditModel {
	tem := &TronconEditModel{Object: tools.O()}
	tem.VM = vm
	tem.Troncon = nil
	tem.Readonly = false
	return tem
}
