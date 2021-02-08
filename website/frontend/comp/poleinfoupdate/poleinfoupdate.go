package poleinfoupdate

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/extraactivityupdatetable"
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

	<!-- Summary -->
    <el-row :gutter="10" type="flex" class="doublespaced">
        <el-col :span="3" class="align-right"><h4 style="margin: 20px 0px 10px 0px">Synthèse :</h4></el-col>
        <el-col :span="19" >
			<el-switch
				  style="display: block"
				  v-model="LineByCity"
				  active-color="#13ce66"
				  inactive-color="#ff4949"
				  active-text="Par Ville"
				  inactive-text="Par Référence"
			></el-switch>
			<el-tabs tab-position="top">
				<!-- ================================== Per Status Tab ============================================= -->
				<el-tab-pane label="Statuts" lazy=true style="padding: 0px 0px;">
					<el-table
							:data="summaryStatusInfos"
							stripe size="mini" show-summary :summary-method="SummaryStatusTotal"
							:default-sort = "{prop: 'Line', order: 'ascending'}"
					>
						<el-table-column
								:label='LineByCity?"Ville":"Référence"' prop="Line" sortable :sort-by="['Line']"
								width="160px" :resizable=true :show-overflow-tooltip=true
						>
							<template slot-scope="scope">
								<span>{{scope.row.Line}} : {{scope.row.Total}}</span>
							</template>
						</el-table-column>
		
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
				<!-- ================================== Per Pole Action To Do Tab ========================================== -->
				<el-tab-pane label="Trx à faire" lazy=true style="padding: 0px 0px;">
					<el-table
							:data="summaryPoleActionInfos"
							stripe size="mini" show-summary :summary-method="SummaryPoleActionTotal"
							:default-sort = "{prop: 'Line', order: 'ascending'}"
					>
						<el-table-column
								:label='LineByCity?"Ville":"Référence"' prop="Line" sortable :sort-by="['Line']"
								width="160px" :resizable=true :show-overflow-tooltip=true
						>
							<template slot-scope="scope">
								<span>{{scope.row.Line}} : {{scope.row.Total}}</span>
							</template>
						</el-table-column>
		
						<el-table-column v-for="poleAction in PoleActionCols"
								:resizable="true" align="center"
								:label="poleAction" width="150px"
								sortable :sort-by="SortBy(poleAction)"
						>
							<template slot-scope="scope">
								<span>{{scope.row.NbPoles[poleAction]}}</span>
							</template>
						</el-table-column>
					</el-table>
				</el-tab-pane>
				<!-- ================================== Per Pole Action Done Tab ========================================== -->
				<el-tab-pane label="Trx faits" lazy=true style="padding: 0px 0px;">
					<el-table
							:data="summaryPoleActionDoneInfos"
							stripe size="mini" show-summary :summary-method="SummaryPoleActionDoneTotal"
							:default-sort = "{prop: 'Line', order: 'ascending'}"
					>
						<el-table-column
								:label='LineByCity?"Ville":"Référence"' prop="Line" sortable :sort-by="['Line']"
								width="160px" :resizable=true :show-overflow-tooltip=true
						>
							<template slot-scope="scope">
								<span>{{scope.row.Line}} : {{scope.row.Total}}</span>
							</template>
						</el-table-column>
		
						<el-table-column v-for="poleAction in PoleActionDoneCols"
								:resizable="true" align="center"
								:label="poleAction" width="150px"
								sortable :sort-by="SortBy(poleAction)"
						>
							<template slot-scope="scope">
								<span>{{scope.row.NbPoles[poleAction]}}</span>
							</template>
						</el-table-column>
					</el-table>
				</el-tab-pane>
				<!-- ================================== Per Pole Type Tab ========================================== -->
				<el-tab-pane label="Besoin Appuis" lazy=true style="padding: 0px 0px;">
					<el-table
							:data="summaryPoleTypeInfos"
							stripe size="mini" show-summary :summary-method="SummaryPoleTypeTotal"
							:default-sort = "{prop: 'Line', order: 'ascending'}"
					>
						<el-table-column
								:label='LineByCity?"Ville":"Référence"' prop="Line" sortable :sort-by="['Line']"
								width="160px" :resizable=true :show-overflow-tooltip=true
						>
							<template slot-scope="scope">
								<span>{{scope.row.Line}} : {{scope.row.Total}}</span>
							</template>
						</el-table-column>
		
						<el-table-column v-for="poleType in PoleTypeCols"
								:resizable="true" align="center"
								:label="poleType" width="97px"
								sortable :sort-by="SortBy(poleType)"
						>
							<template slot-scope="scope">
								<span>{{scope.row.NbPoles[poleType]}}</span>
							</template>
						</el-table-column>
					</el-table>
				</el-tab-pane>
				<!-- ================================== Per Pole Item Tab ========================================== -->
				<el-tab-pane label="Besoin Matériel" lazy=true style="padding: 0px 0px;">
					<el-table
							:data="summaryPoleItemInfos"
							stripe size="mini" show-summary :summary-method="SummaryPoleItemTotal"
							:default-sort = "{prop: 'Line', order: 'ascending'}"
					>
						<el-table-column
								:label='LineByCity?"Ville":"Référence"' prop="Line" sortable :sort-by="['Line']"
								width="160px" :resizable=true :show-overflow-tooltip=true
						>
							<template slot-scope="scope">
								<span>{{scope.row.Line}} : {{scope.row.Total}}</span>
							</template>
						</el-table-column>
		
						<el-table-column v-for="poleItem in PoleItemCols"
								:resizable="true" align="center"
								:label="poleItem" width="120px"
								sortable :sort-by="SortBy(poleItem)"
						>
							<template slot-scope="scope">
								<span>{{scope.row.NbPoles[poleItem]}}</span>
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
		extraactivityupdatetable.RegisterComponent(),
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
		hvue.Computed("summaryStatusInfos", func(vm *hvue.VM) interface{} {
			pium := PoleInfoUpdateModelFromJS(vm.Object)
			summer := NewSummarizer(pium.LineByCity)
			summer.Calc(pium.Polesite.Poles)
			pium.StateCols = pium.GetSummaryStatuses()
			return summer.SummaryDatas
		}),
		hvue.Computed("summaryPoleActionInfos", func(vm *hvue.VM) interface{} {
			pium := PoleInfoUpdateModelFromJS(vm.Object)
			summer := NewSummarizer(pium.LineByCity)
			summer.GetColumn = GetPoleAction
			summer.Calc(pium.Polesite.Poles)
			pium.PoleActionCols = summer.GetCalcColumns()
			return summer.SummaryDatas
		}),
		hvue.Computed("summaryPoleActionDoneInfos", func(vm *hvue.VM) interface{} {
			pium := PoleInfoUpdateModelFromJS(vm.Object)
			summer := NewSummarizer(pium.LineByCity)
			summer.GetColumn = GetPoleActionDone
			summer.Calc(pium.Polesite.Poles)
			pium.PoleActionDoneCols = summer.GetCalcColumns()
			return summer.SummaryDatas
		}),
		hvue.Computed("summaryPoleTypeInfos", func(vm *hvue.VM) interface{} {
			pium := PoleInfoUpdateModelFromJS(vm.Object)
			summer := NewSummarizer(pium.LineByCity)
			summer.GetColumn = GetPoleType
			summer.Calc(pium.Polesite.Poles)
			pium.PoleTypeCols = summer.GetCalcColumns()
			return summer.SummaryDatas
		}),
		hvue.Computed("summaryPoleItemInfos", func(vm *hvue.VM) interface{} {
			pium := PoleInfoUpdateModelFromJS(vm.Object)
			summer := NewSummarizer(pium.LineByCity)
			summer.GetColumn = GetPoleItem
			summer.Calc(pium.Polesite.Poles)
			pium.PoleItemCols = summer.GetCalcColumns()
			return summer.SummaryDatas
		}),
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

	LineByCity         bool     `js:"LineByCity"`
	StateCols          []string `js:"StateCols"`
	PoleTypeCols       []string `js:"PoleTypeCols"`
	PoleActionCols     []string `js:"PoleActionCols"`
	PoleActionDoneCols []string `js:"PoleActionDoneCols"`
	PoleItemCols       []string `js:"PoleItemCols"`

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

	pium.LineByCity = true
	pium.StateCols = []string{}
	pium.PoleTypeCols = []string{}
	pium.PoleActionCols = []string{}
	pium.PoleActionDoneCols = []string{}
	pium.PoleItemCols = []string{}

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
		//poleconst.StateNotSubmitted,
		poleconst.StateDictToDo,
		poleconst.StateDaToDo,
		poleconst.StateDaExpected,
		poleconst.StatePermissionPending,
		poleconst.StateToDo,
		poleconst.StateHoleDone,
		poleconst.StateIncident,
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

func (pium *PoleInfoUpdateModel) summaryTotal(vm *hvue.VM, param *js.Object, statuses []string) []string {
	nbRes := make([]int, len(statuses))
	param.Get("data").Call("forEach", func(sd *SummaryData) {
		for i, status := range statuses {
			nbRes[i] += sd.NbPoles[status]
		}
	})
	res := make([]string, len(statuses)+1)
	total := 0
	for i, nb := range nbRes {
		if nb == 0 {
			continue
		}
		res[i+1] = strconv.Itoa(nb)
		total += nb
	}
	res[0] = "Total : " + strconv.Itoa(total)
	return res
}

func (pium *PoleInfoUpdateModel) SummaryStatusTotal(vm *hvue.VM, param *js.Object) []string {
	pium = PoleInfoUpdateModelFromJS(vm.Object)
	return pium.summaryTotal(vm, param, pium.StateCols)
}

func (pium *PoleInfoUpdateModel) SummaryPoleTypeTotal(vm *hvue.VM, param *js.Object) []string {
	pium = PoleInfoUpdateModelFromJS(vm.Object)
	return pium.summaryTotal(vm, param, pium.PoleTypeCols)
}

func (pium *PoleInfoUpdateModel) SummaryPoleActionTotal(vm *hvue.VM, param *js.Object) []string {
	pium = PoleInfoUpdateModelFromJS(vm.Object)
	return pium.summaryTotal(vm, param, pium.PoleActionCols)
}

func (pium *PoleInfoUpdateModel) SummaryPoleActionDoneTotal(vm *hvue.VM, param *js.Object) []string {
	pium = PoleInfoUpdateModelFromJS(vm.Object)
	return pium.summaryTotal(vm, param, pium.PoleActionDoneCols)
}

func (pium *PoleInfoUpdateModel) SummaryPoleItemTotal(vm *hvue.VM, param *js.Object) []string {
	pium = PoleInfoUpdateModelFromJS(vm.Object)
	return pium.summaryTotal(vm, param, pium.PoleItemCols)
}
