package ptedit

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

const template string = `
<el-row :gutter="10" type="flex" align="middle">
    <el-col :span="1">
        <span><strong>{{title}}:</strong></span>
    </el-col>
    <el-col :span="4">
        <el-input :placeholder="refPH" :readonly="readonly" clearable size="mini"
                  v-model="value.Ref"
        ></el-input>
    </el-col>
    <el-col :span="4">
        <el-input placeholder="PT-009999" :readonly="readonly" clearable size="mini"
                  v-model="value.RefPt"
        ></el-input>
    </el-col>
    <el-col :span="15">
        <el-input :placeholder="addressPH" :readonly="readonly" clearable size="mini"
                  v-model="value.Address"
        >
        </el-input>
    </el-col>
</el-row>
`

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Registration

func Register() {
	hvue.NewComponent("pt-edit",
		ComponentOptions()...,
	)
}

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("pt-edit", ComponentOptions()...)
}

func ComponentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.Props("title", "readonly", "value"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewPtEditModel(vm)
		}),
		hvue.Computed("refPH", func(vm *hvue.VM) interface{} {
			pem := &PtEditModel{Object: vm.Object}
			return pem.Title + "-99999"
		}),
		hvue.Computed("refptPH", func(vm *hvue.VM) interface{} {
			pem := &PtEditModel{Object: vm.Object}
			return pem.Title + "-0099999"
		}),
		hvue.Computed("addressPH", func(vm *hvue.VM) interface{} {
			pem := &PtEditModel{Object: vm.Object}
			return "Adresse " + pem.Title
		}),
		hvue.MethodsOf(&PtEditModel{}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type PtEditModel struct {
	*js.Object

	Pt       *fm.PT `js:"value"`
	Readonly bool   `js:"readonly"`
	Title    string `js:"title"`

	VM *hvue.VM `js:"VM"`
}

func NewPtEditModel(vm *hvue.VM) *PtEditModel {
	pem := &PtEditModel{Object: tools.O()}
	pem.VM = vm
	pem.Pt = nil
	pem.Readonly = false
	pem.Title = "PT"
	return pem
}
