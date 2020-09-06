package poleinfoupdate

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/ripprogressbar"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/model/polesite"
	"github.com/lpuig/ewin/doe/website/frontend/model/polesite/poleconst"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
	"strconv"
)

const template string = `
<div style="padding: 20px 25px; height: calc(100% - 40px);">
	<!-- Client & Ref -->
    <el-row :gutter="10" type="flex" align="middle">
        <el-col :span="3" class="align-right">Client :</el-col>
        <el-col :span="8">
            <el-input placeholder="Client"
                      v-model="value.Client" clearable size="mini"
					  @change="UpdateTitle"
            ></el-input>
        </el-col>

        <el-col :span="3" class="align-right">Référence du chantier :</el-col>
        <el-col :span="8">
            <el-input placeholder="Référence"
                      v-model="value.Ref" clearable size="mini"
					  @change="UpdateTitle"
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

	<!-- Progression -->
    <el-row :gutter="10" type="flex" align="middle" class="doublespaced">
        <el-col :span="3" class="align-right"><h4 style="margin: 20px 0px 10px 0px">Avancement :</h4></el-col>
        <el-col :span="19">
				<ripsiteinfo-progress-bar height="10px" :total="statNbPole" :blocked="statNbPoleBlocked" :billed="statNbPoleBilled" :done="statNbPoleDone"></ripsiteinfo-progress-bar>
        </el-col>
    </el-row>


	<!-- Summary -->
    <el-row :gutter="10" type="flex" class="doublespaced">
        <el-col :span="3" class="align-right"><h4 style="margin: 20px 0px 10px 0px">Synthèse :</h4></el-col>
        <el-col :span="19" >
			<el-tabs tab-position="top">
				<!-- ===================================== Bonus Tab ======================================================= -->
				<el-tab-pane label="Statuts" lazy=true style="height: 44vh; padding: 5px 25px; overflow-x: hidden;overflow-y: auto;">
					<el-table
							:data="summaryInfos"
							stripe size="mini" show-summary :summary-method="SummaryTotal"
							:default-sort = "{prop: 'City', order: 'ascending'}"
					>
						<el-table-column
								label="Ville" prop="Line" sortable :sort-by="['Line']"
								width="200px" :resizable=true :show-overflow-tooltip=true
						></el-table-column>
		
						<el-table-column v-for="status in GetSummaryStatuses()"
								:resizable="true" align="center"
								:label="StateName(status)" width="150px"
								sortable :sort-by="SortBy(status)"
						>
							<template slot-scope="scope">
								<span>{{scope.row.NbPoles[status]}}</span>
							</template>
						</el-table-column>
					</el-table>
				</el-tab-pane>
			</el-tabs>
        </el-col>
    </el-row>
</div>
`

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("polesite-info-update", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		ripprogressbar.RegisterComponent(),
		hvue.Template(template),
		hvue.Props("value", "user"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewPoleInfoUpdateModel(vm)
		}),
		hvue.MethodsOf(&PoleInfoUpdateModel{}),
		hvue.Computed("statNbPole", func(vm *hvue.VM) interface{} {
			pium := PoleInfoUpdateModelFromJS(vm.Object)
			return pium.CalcStat()
		}),
		hvue.Computed("summaryInfos", func(vm *hvue.VM) interface{} {
			pium := PoleInfoUpdateModelFromJS(vm.Object)
			summer := NewSummarizer()
			summer.Colums = pium.GetSummaryStatuses()
			summer.Calc(pium.Polesite.Poles)
			return summer.SummaryDatas
			//return CalcSummaryDatas(pium.Polesite.Poles, pium.GetSummaryStatuses())
		}),
		//hvue.Computed("PullingTotal", func(vm *hvue.VM) interface{} {
		//	pium := PoleInfoUpdateModelFromJS(vm.Object)
		//	return rium.SetPullingStats()
		//}),
		//hvue.Filter("FormatTronconRef", func(vm *hvue.VM, value *js.Object, args ...*js.Object) interface{} {
		//	rpum := RipPullingUpdateModelFromJS(vm.Object)
		//	t := &fm.Troncon{Object: value}
		//	return rpum.GetFormatTronconRef(t)
		//}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type PoleInfoUpdateModel struct {
	*js.Object

	Polesite *polesite.Polesite `js:"value"`
	User     *fm.User           `js:"user"`

	StatNbPoleBlocked  int `js:"statNbPoleBlocked"`
	StatNbPoleDone     int `js:"statNbPoleDone"`
	StatNbPoleBilled   int `js:"statNbPoleBilled"`
	StatNbPoleDICTToDo int `js:"statNbPoleDICTToDo"`
	StatNbPolePending  int `js:"statNbPolePending"`

	VM *hvue.VM `js:"VM"`
}

func NewPoleInfoUpdateModel(vm *hvue.VM) *PoleInfoUpdateModel {
	pium := &PoleInfoUpdateModel{Object: tools.O()}
	pium.VM = vm
	pium.Polesite = polesite.NewPolesite()
	pium.User = nil

	pium.StatNbPoleBlocked = 0
	pium.StatNbPoleDone = 0
	pium.StatNbPoleBilled = 0
	pium.StatNbPoleDICTToDo = 0
	pium.StatNbPolePending = 0

	return pium
}

func PoleInfoUpdateModelFromJS(o *js.Object) *PoleInfoUpdateModel {
	return &PoleInfoUpdateModel{Object: o}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Data Items Related Methods

func (pium *PoleInfoUpdateModel) UpdateTitle(vm *hvue.VM) {
	pium = PoleInfoUpdateModelFromJS(vm.Object)
	newTitle := pium.Polesite.Ref + " - " + pium.Polesite.Client
	js.Global.Get("document").Set("title", newTitle)
}

func (pium *PoleInfoUpdateModel) GetStates(vm *hvue.VM) []*elements.ValueLabel {
	//pium := PoleInfoUpdateModelFromJS(vm.Object)
	return polesite.GetPoleSiteStatesValueLabel()
}

func (pium *PoleInfoUpdateModel) CalcStat() int {
	var tot, done, billed, blocked, dict, pending int
	for _, pole := range pium.Polesite.Poles {
		if !pole.IsToDo() {
			continue
		}
		tot++
		switch {
		case pole.IsAttachment():
			billed++
		case pole.IsDone():
			done++
		case pole.IsBlocked():
			blocked++
			switch pole.State {
			case poleconst.StateDictToDo:
				dict++
			case poleconst.StatePermissionPending:
				pending++
			}
		}

	}
	pium.StatNbPoleBilled = billed
	pium.StatNbPoleBlocked = blocked
	pium.StatNbPoleDone = done
	pium.StatNbPoleDICTToDo = dict
	pium.StatNbPolePending = pending
	return tot
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Summary Data Table Related Methods

func (pium *PoleInfoUpdateModel) GetSummaryStatuses() []string {
	//pium := PoleInfoUpdateModelFromJS(vm.Object)
	return []string{
		poleconst.StateNotSubmitted,
		poleconst.StateDictToDo,
		poleconst.StateDaToDo,
		poleconst.StateDaExpected,
		poleconst.StatePermissionPending,
		poleconst.StateToDo,
		poleconst.StateDone,
		poleconst.StateAttachment,
	}
}

func (pium *PoleInfoUpdateModel) StateName(status string) string {
	return polesite.PoleStateLabel(status)
}

func (pium *PoleInfoUpdateModel) SortBy(attrib string) func(obj *js.Object) int {
	return func(obj *js.Object) int {
		return obj.Get("NbPoles").Get(attrib).Int()
	}
}

func (pium *PoleInfoUpdateModel) SummaryTotal(vm *hvue.VM, param *js.Object) []string {
	pium = PoleInfoUpdateModelFromJS(vm.Object)
	statuses := pium.GetSummaryStatuses()
	nbRes := make([]int, len(statuses))
	param.Get("data").Call("forEach", func(sd *SummaryData) {
		for i, status := range statuses {
			nbRes[i] += sd.NbPoles[status]
		}
	})
	res := make([]string, len(statuses)+1)
	res[0] = "Total :"
	for i, nb := range nbRes {
		if nb == 0 {
			continue
		}
		res[i+1] = strconv.Itoa(nb)
	}
	return res
}
