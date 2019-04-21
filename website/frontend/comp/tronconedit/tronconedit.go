package tronconedit

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/ptedit"
	"github.com/lpuig/ewin/doe/website/frontend/comp/tronconstatustag"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/autocomplete"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
	"strings"
)

const template string = `
<div>
    <!-- 
        Status, Troncon and PB info, Size, Nb El
    -->
    <el-row :gutter="10" type="flex" align="middle">
        <el-col :span="3">
            <troncon-status-tag v-model="value"></troncon-status-tag>
        </el-col>

        <el-col :span="3">
            <el-autocomplete v-model="value.Ref"
                             :fetch-suggestions="RefSearch"
                             placeholder="TR-99-9999"
                             clearable size="mini" style="width: 100%"
                             @input="CheckRef(value)"
            >
                <template slot="prepend">Tronçon:</template>
            </el-autocomplete>
        </el-col>

        <el-col :span="12">
            <pt-edit title="PB" v-model="value.Pb" :readonly="readonly" :info="LastPBinfo()"></pt-edit>
        </el-col>

        <el-col :span="2">
            <el-tooltip content="Nb. Fibre" placement="bottom" effect="light" :open-delay="500">
                <el-input-number v-model="value.NbFiber"
                                 :min="6" :step="6"
                                 :readonly="readonly"
                                 size="mini" style="width: 100%"
                                 @input="CheckFiber(value)"
                ></el-input-number>
<!--                    <template slot="prepend">Nb Fibre</template>-->
            </el-tooltip>
        </el-col>

        <el-col :span="4">
            <el-tooltip style="padding: 0px 12px;" content="Nb. EL raccordable" placement="bottom" effect="light" :open-delay="500">
                <el-slider
                        v-model="value.NbRacco"
                        :min="0" :max="value.NbFiber"
                        :step="1"
                        :readonly="readonly"
                        show-stops
                ></el-slider>
<!--					show-input input-size="mini" :show-input-controls="false"-->
            </el-tooltip>
        </el-col>  
        
    </el-row>

    <!-- 
        Blockage & Signature Request, Comment
    -->
    <el-row :gutter="10" type="flex" align="middle">
        <el-col :span="3">
            <el-switch v-model="value.NeedSignature"
                       active-color="#db2828"
                       active-text="Signature demandée"
                       inactive-color="#bcbcbc"
                       @input="CheckSignature(value)"
            ></el-switch>
        </el-col>

        <el-col :span="3">
            <el-switch v-model="value.Blockage"
                       active-color="#db2828"
                       active-text="Blocage"
                       inactive-color="#bcbcbc"
                       :disabled="value.NeedSignature && !value.Signed"
            ></el-switch>
        </el-col>

        <el-col :span="14">
            <el-input :readonly="readonly" clearable placeholder="Commentaire sur tronçon" size="mini" type="textarea" autosize
                      v-model="value.Comment" :disabled="!value.Blockage"
            ></el-input>
        </el-col>

		<el-col :span="4">
			<el-select v-model="value.Article" filterable placeholder="CEM" size="mini" style="width: 100%">
				<el-option
						v-for="item in articles"
						:key="item.value"
						:label="item.label"
						:value="item.value">
				</el-option>
			</el-select>                            
        </el-col>
    </el-row>
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
		tronconstatustag.RegisterComponent(),
		hvue.Template(template),
		hvue.Props("value", "readonly", "previous", "articles"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewTronconEditModel(vm)
		}),
		hvue.MethodsOf(&TronconEditModel{}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type TronconEditModel struct {
	*js.Object

	Articles    []*elements.ValueLabel `js:"articles"`
	Troncon     *fm.Troncon            `js:"value"`
	PrevTroncon *fm.Troncon            `js:"previous"`
	Readonly    bool                   `js:"readonly"`

	VM *hvue.VM `js:"VM"`
}

func NewTronconEditModel(vm *hvue.VM) *TronconEditModel {
	tem := &TronconEditModel{Object: tools.O()}
	tem.VM = vm
	tem.Articles = nil
	tem.Troncon = nil
	tem.Readonly = false
	return tem
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Actions

func (tem *TronconEditModel) CheckRef(tr *fm.Troncon) {
	if strings.HasPrefix(tr.Ref, "TR") && len(tr.Ref) > 3 {
		tr.Ref = strings.Replace(tr.Ref, " ", "-", -1)
		return
	}
	if !strings.HasPrefix(tr.Ref, "TR-") {
		tr.Ref = "TR-" + tr.Ref
	}
}

func (tem *TronconEditModel) CheckFiber(tr *fm.Troncon) {
	if tr.NbFiber < tr.NbRacco {
		tr.NbRacco = tr.NbFiber
	}
}

func (tem *TronconEditModel) CheckSignature(tr *fm.Troncon) {
	tr.CheckSignature()
}

func (tem *TronconEditModel) RefSearch(vm *hvue.VM, query string, callback *js.Object) {
	tem = &TronconEditModel{Object: vm.Object}
	// if no previous troncon.ref return default choice list
	res := []*autocomplete.Result{}
	if tem.PrevTroncon == nil || tem.PrevTroncon.Object == nil || tem.PrevTroncon.Object == js.Undefined {
		callback.Invoke(res)
		return
	}
	// retrieve previous troncon.Ref
	lastref := tem.PrevTroncon.Ref
	if lastref == "" || !strings.HasPrefix(lastref, "TR-") {
		callback.Invoke(res)
		return
	}
	refchunck := strings.Split(lastref, "-")
	res = autocomplete.GenResults(strings.Join(refchunck[:2], "-")+"-", refchunck[2], 4)
	callback.Invoke(res)
}

func (tem *TronconEditModel) LastPBinfo(vm *hvue.VM) js.M {
	tem = &TronconEditModel{Object: vm.Object}
	pbRef := ""
	ptRef := ""
	// if no previous troncon return default choice list
	if tem.PrevTroncon == nil || tem.PrevTroncon.Object == nil || tem.PrevTroncon.Object == js.Undefined {
		return js.M{"PB": pbRef, "PT": ptRef}
	}
	// retrieve last troncon.PB
	lastPb := tem.PrevTroncon.Pb
	if lastPb.Ref != "" && strings.HasPrefix(lastPb.Ref, "PB-") {
		pbRef = strings.TrimPrefix(lastPb.Ref, "PB-")
	}
	if lastPb.RefPt != "" && strings.HasPrefix(lastPb.RefPt, "PT-") {
		ptRef = strings.TrimPrefix(lastPb.RefPt, "PT-")
	}
	return js.M{"PB": pbRef, "PT": ptRef}
}

func (tem *TronconEditModel) CheckMeasureDate(vm *hvue.VM, date string) {
	tem = &TronconEditModel{Object: vm.Object}
	if date < tem.Troncon.InstallDate {
		tem.Troncon.MeasureDate = tem.Troncon.InstallDate
	}
}
