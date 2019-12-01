package foainfoupdate

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/ripprogressbar"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	fmfoa "github.com/lpuig/ewin/doe/website/frontend/model/foasite"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

const template string = `
<div style="padding: 5px 25px">
	<!-- Client & Ref -->
    <el-row :gutter="10" type="flex" align="middle" class="doublespaced">
        <el-col :span="3" class="align-right">Client :</el-col>
        <el-col :span="8">
            <el-input placeholder="Client"
                      v-model="value.Client" clearable size="mini"
            ></el-input>
        </el-col>

        <el-col :span="3" class="align-right">Référence du chantier :</el-col>
        <el-col :span="8">
            <el-input placeholder="Référence"
                      v-model="value.Ref" clearable size="mini"
            ></el-input>
        </el-col>
    </el-row>

	<!-- Manager & Order Date -->
    <el-row :gutter="10" type="flex" align="middle" class="doublespaced">
        <el-col :span="3" class="align-right">Chargé d'affaire :</el-col>
        <el-col :span="8">
            <el-input placeholder="Caff."
                      v-model="value.Manager" clearable size="mini"
            ></el-input>
        </el-col>

        <el-col :span="3" class="align-right">Date de commande :</el-col>
        <el-col :span="8">
            <el-date-picker format="dd/MM/yyyy" placeholder="Date" size="mini"
                            style="width: 100%" type="date"
                            v-model="value.OrderDate"
                            value-format="yyyy-MM-dd"
                            :picker-options="{firstDayOfWeek:1, disabledDate(time) { return time.getTime() > Date.now(); }}"
                            :clearable="false"
            ></el-date-picker>
        </el-col>
    </el-row>

	<!-- Comment -->
    <el-row :gutter="10" type="flex" align="middle" class="doublespaced">
        <el-col :span="3" class="align-right">Commentaire :</el-col>
        <el-col :span="19">
            <el-input type="textarea" :autosize="{ minRows: 2, maxRows: 5}" placeholder="Commentaire"
                      v-model="value.Comment" clearable size="mini"
            ></el-input>
        </el-col>
    </el-row>

	<!-- Foa Progress -->
    <el-row :gutter="10" type="flex" align="middle" class="doublespaced">
        <el-col :span="3" class="align-right"><h4 style="margin: 20px 0px 10px 0px">Avancement :</h4></el-col>
        <el-col :span="19">
            <ripsiteinfo-progress-bar height="10px" :total="FoaTotal" :done="FoaDone" :blocked="FoaBlocked"></ripsiteinfo-progress-bar>
        </el-col>
    </el-row>
</div>
`

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("foa-info-update", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		ripprogressbar.RegisterComponent(),
		hvue.Template(template),
		hvue.Props("value", "user"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewFoaInfoUpdateModel(vm)
		}),
		hvue.MethodsOf(&FoaInfoUpdateModel{}),
		hvue.Computed("FoaTotal", func(vm *hvue.VM) interface{} {
			fium := FoaInfoUpdateModelFromJS(vm.Object)
			return fium.SetFoaStats()
		}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type FoaInfoUpdateModel struct {
	*js.Object

	Foasite *fmfoa.FoaSite `js:"value"`
	User    *fm.User       `js:"user"`

	//FoaTotal   int `js:"FoaTotal"`
	FoaDone    int `js:"FoaDone"`
	FoaBlocked int `js:"FoaBlocked"`

	VM *hvue.VM `js:"VM"`
}

func NewFoaInfoUpdateModel(vm *hvue.VM) *FoaInfoUpdateModel {
	rmum := &FoaInfoUpdateModel{Object: tools.O()}
	rmum.VM = vm
	rmum.Foasite = fmfoa.NewFoaSite()
	rmum.User = fm.NewUser()

	return rmum
}

func FoaInfoUpdateModelFromJS(o *js.Object) *FoaInfoUpdateModel {
	return &FoaInfoUpdateModel{Object: o}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Formatting Related Methods

// SetFoaStats sets FoaDone and FoaBlocked values, and returns FoaTotal
func (fium *FoaInfoUpdateModel) SetFoaStats() int {
	total, blocked, done := fium.Foasite.GetProgress()
	fium.FoaBlocked = blocked
	fium.FoaDone = done

	return total
}
