package ripinfoupdate

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/extraactivityupdatetable"
	"github.com/lpuig/ewin/doe/website/frontend/comp/ripprogressbar"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	fmrip "github.com/lpuig/ewin/doe/website/frontend/model/ripsite"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
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

	<!-- Status -->
    <el-row v-if="user.Permissions.Invoice" :gutter="10" type="flex" align="middle" class="doublespaced">
        <el-col :span="3" class="align-right">Statut Chantier :</el-col>
        <el-col :span="8">
			<el-select v-model="value.Status" placeholder="Statut" size="mini" style="width: 100%"
					   @clear=""
					   @change=""
			>
				<el-option
						v-for="item in GetStates()"
						:key="item.value"
						:label="item.label"
						:value="item.value"
				>
				</el-option>
			</el-select>
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

	<!-- Pulling Progress -->
    <el-row :gutter="10" type="flex" align="middle" class="doublespaced">
        <el-col :span="3" class="align-right"><h4 style="margin: 20px 0px 10px 0px">Total Tirage :</h4></el-col>
        <el-col :span="19">
            <ripsiteinfo-progress-bar height="10px" :total="PullingTotal" :done="PullingDone" :blocked="PullingBlocked"
            ></ripsiteinfo-progress-bar>
        </el-col>
    </el-row>

    <el-row :gutter="10" type="flex" align="middle" class="spaced">
        <el-col :span="3" class="align-right">Câbles :</el-col>
        <el-col :span="19">
            <ripsiteinfo-progress-bar height="7px" :total="PullCableTotal" :done="PullCableDone" :blocked="PullCableBlocked"
            ></ripsiteinfo-progress-bar>
        </el-col>
    </el-row>
    <el-row :gutter="10" type="flex" align="middle" class="spaced">
        <el-col :span="3" class="align-right">Souterrain :</el-col>
        <el-col :span="19">
            <ripsiteinfo-progress-bar height="7px" :total="PullUndTotal" :done="PullUndDone" :blocked="PullUndBlocked"
            ></ripsiteinfo-progress-bar>
        </el-col>
    </el-row>
    <el-row :gutter="10" type="flex" align="middle" class="spaced">
        <el-col :span="3" class="align-right">Aérien :</el-col>
        <el-col :span="19">
            <ripsiteinfo-progress-bar height="7px" :total="PullAerTotal" :done="PullAerDone" :blocked="PullAerBlocked"
            ></ripsiteinfo-progress-bar>
        </el-col>
    </el-row>
    <el-row :gutter="10" type="flex" align="middle" class="spaced">
        <el-col :span="3" class="align-right">Façade :</el-col>
        <el-col :span="19">
            <ripsiteinfo-progress-bar height="7px" :total="PullBuildTotal" :done="PullBuildDone" :blocked="PullBuildBlocked"
            ></ripsiteinfo-progress-bar>
        </el-col>
    </el-row>

	<!-- Junction Progress -->
    <el-row :gutter="10" type="flex" align="middle" class="doublespaced">
        <el-col :span="3" class="align-right"><h4 style="margin: 20px 0px 10px 0px">Total Raccordement :</h4></el-col>
        <el-col :span="19">
            <ripsiteinfo-progress-bar height="10px" :total="JunctionTotal" :done="JunctionDone" :blocked="JunctionBlocked"
            ></ripsiteinfo-progress-bar>
        </el-col>
    </el-row>

    <el-row :gutter="10" type="flex" align="middle" class="spaced">
        <el-col :span="3" class="align-right">Boitiers :</el-col>
        <el-col :span="19">
            <ripsiteinfo-progress-bar height="7px" :total="JunctionNodeTotal" :done="JunctionNodeDone" :blocked="JunctionNodeBlocked"
            ></ripsiteinfo-progress-bar>
        </el-col>
    </el-row>
    <el-row :gutter="10" type="flex" align="middle" class="spaced">
        <el-col :span="3" class="align-right">Epissures :</el-col>
        <el-col :span="19">
            <ripsiteinfo-progress-bar height="7px" :total="JunctionSpliceTotal" :done="JunctionSpliceDone" :blocked="JunctionSpliceBlocked"
            ></ripsiteinfo-progress-bar>
        </el-col>
    </el-row>

	<!-- Measurement Progress -->
    <el-row :gutter="10" type="flex" align="middle" class="doublespaced">
        <el-col :span="3" class="align-right"><h4 style="margin: 20px 0px 10px 0px">Total Mesures :</h4></el-col>
        <el-col :span="19">
            <ripsiteinfo-progress-bar height="10px" :total="MeasurementFiberTotal" :done="MeasurementFiberDone" :blocked="MeasurementFiberBlocked"
            ></ripsiteinfo-progress-bar>
        </el-col>
    </el-row>
    <el-row :gutter="10" type="flex" align="middle" class="spaced">
        <el-col :span="3" class="align-right">Boitiers :</el-col>
        <el-col :span="19">
            <ripsiteinfo-progress-bar height="7px" :total="MeasurementNodeTotal" :done="MeasurementNodeDone" :blocked="MeasurementNodeBlocked"
            ></ripsiteinfo-progress-bar>
        </el-col>
    </el-row>

	<!--	Extra Activities-->
    <el-row :gutter="10" type="flex" align="top" class="doublespaced">
        <el-col :span="3" class="align-right"><h4 style="margin: 20px 0px 10px 0px">Activité complémentaire :</h4></el-col>
        <el-col :span="19">
			<extra-activities-table
					v-model="value.ExtraActivities"
					:User="user"
					:Client="value.Client"
			></extra-activities-table>
		</el-col>
    </el-row>

</div>
`

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("rip-info-update", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		ripprogressbar.RegisterComponent(),
		extraactivityupdatetable.RegisterComponent(),
		hvue.Template(template),
		hvue.Props("value", "user"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewRipInfoUpdateModel(vm)
		}),
		hvue.MethodsOf(&RipInfoUpdateModel{}),
		hvue.Computed("PullingTotal", func(vm *hvue.VM) interface{} {
			rium := RipInfoUpdateModelFromJS(vm.Object)
			return rium.SetPullingStats()
		}),
		hvue.Computed("JunctionTotal", func(vm *hvue.VM) interface{} {
			rium := RipInfoUpdateModelFromJS(vm.Object)
			return rium.SetJunctionStats()
		}),
		hvue.Computed("MeasurementFiberTotal", func(vm *hvue.VM) interface{} {
			rium := RipInfoUpdateModelFromJS(vm.Object)
			return rium.SetMeasurementStats()
		}),
		//hvue.Filter("FormatTronconRef", func(vm *hvue.VM, value *js.Object, args ...*js.Object) interface{} {
		//	rpum := RipPullingUpdateModelFromJS(vm.Object)
		//	t := &fm.Troncon{Object: value}
		//	return rpum.GetFormatTronconRef(t)
		//}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type RipInfoUpdateModel struct {
	*js.Object

	Ripsite *fmrip.Ripsite `js:"value"`
	User    *fm.User       `js:"user"`

	//PullingTotal int `js:"PullingTotal"`
	PullingDone      int `js:"PullingDone"`
	PullingBlocked   int `js:"PullingBlocked"`
	PullCableTotal   int `js:"PullCableTotal"`
	PullCableDone    int `js:"PullCableDone"`
	PullCableBlocked int `js:"PullCableBlocked"`
	PullUndTotal     int `js:"PullUndTotal"`
	PullUndDone      int `js:"PullUndDone"`
	PullUndBlocked   int `js:"PullUndBlocked"`
	PullAerTotal     int `js:"PullAerTotal"`
	PullAerDone      int `js:"PullAerDone"`
	PullAerBlocked   int `js:"PullAerBlocked"`
	PullBuildTotal   int `js:"PullBuildTotal"`
	PullBuildDone    int `js:"PullBuildDone"`
	PullBuildBlocked int `js:"PullBuildBlocked"`

	//JunctionTotal int `js:"JunctionTotal"`
	JunctionDone          int `js:"JunctionDone"`
	JunctionBlocked       int `js:"JunctionBlocked"`
	JunctionNodeTotal     int `js:"JunctionNodeTotal"`
	JunctionNodeDone      int `js:"JunctionNodeDone"`
	JunctionNodeBlocked   int `js:"JunctionNodeBlocked"`
	JunctionSpliceTotal   int `js:"JunctionSpliceTotal"`
	JunctionSpliceDone    int `js:"JunctionSpliceDone"`
	JunctionSpliceBlocked int `js:"JunctionSpliceBlocked"`

	//MeasurementFiberTotal   int `js:"MeasurementFiberTotal"`
	MeasurementFiberDone    int `js:"MeasurementFiberDone"`
	MeasurementFiberBlocked int `js:"MeasurementFiberBlocked"`
	MeasurementNodeTotal    int `js:"MeasurementNodeTotal"`
	MeasurementNodeDone     int `js:"MeasurementNodeDone"`
	MeasurementNodeBlocked  int `js:"MeasurementNodeBlocked"`

	VM *hvue.VM `js:"VM"`
}

func NewRipInfoUpdateModel(vm *hvue.VM) *RipInfoUpdateModel {
	rmum := &RipInfoUpdateModel{Object: tools.O()}
	rmum.VM = vm
	rmum.Ripsite = fmrip.NewRisite()
	rmum.User = nil

	return rmum
}

func RipInfoUpdateModelFromJS(o *js.Object) *RipInfoUpdateModel {
	return &RipInfoUpdateModel{Object: o}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Formatting Related Methods

// SetPullingStats sets PullingDone and PullingBlocked values, and returns PullingTotal
func (rium *RipInfoUpdateModel) SetPullingStats() int {
	tot, cable, und, aer, build := rium.Ripsite.GetPullingProgresses()
	rium.PullingBlocked = tot.Blocked
	rium.PullingDone = tot.Done

	rium.PullCableTotal = cable.Total
	rium.PullCableBlocked = cable.Blocked
	rium.PullCableDone = cable.Done

	rium.PullUndTotal = und.Total
	rium.PullUndBlocked = und.Blocked
	rium.PullUndDone = und.Done

	rium.PullAerTotal = aer.Total
	rium.PullAerBlocked = aer.Blocked
	rium.PullAerDone = aer.Done

	rium.PullBuildTotal = build.Total
	rium.PullBuildBlocked = build.Blocked
	rium.PullBuildDone = build.Done

	return tot.Total
}

// SetJunctionStats sets JunctionDone and JunctionBlocked values, and returns JunctionTotal
func (rium *RipInfoUpdateModel) SetJunctionStats() int {
	fibers, nodes, splices := rium.Ripsite.GetJunctionProgresses()
	rium.JunctionBlocked = fibers.Blocked
	rium.JunctionDone = fibers.Done

	rium.JunctionNodeTotal = nodes.Total
	rium.JunctionNodeBlocked = nodes.Blocked
	rium.JunctionNodeDone = nodes.Done

	rium.JunctionSpliceTotal = splices.Total
	rium.JunctionSpliceBlocked = splices.Blocked
	rium.JunctionSpliceDone = splices.Done

	return fibers.Total
}

// SetMeasurementStats sets MeasurementDone and MeasurementBlocked values, and returns MeasurementTotal
func (rium *RipInfoUpdateModel) SetMeasurementStats() int {
	nodes, fibers := rium.Ripsite.GetMeasurementProgresses()

	rium.MeasurementNodeTotal = nodes.Total
	rium.MeasurementNodeBlocked = nodes.Blocked
	rium.MeasurementNodeDone = nodes.Done

	//rium.MeasurementFiberTotal = fibers.Total
	rium.MeasurementFiberBlocked = fibers.Blocked
	rium.MeasurementFiberDone = fibers.Done

	return fibers.Total
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Data Items Related Methods

func (rium *RipInfoUpdateModel) GetStates(vm *hvue.VM) []*elements.ValueLabel {
	//rium := RipInfoUpdateModelFromJS(vm.Object)
	return fmrip.GetStatesValueLabel()
}
