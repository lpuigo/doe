package poletable

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	ps "github.com/lpuig/ewin/doe/website/frontend/model/polesite"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/date"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements/message"
	"strconv"
)

const template string = `<el-container  style="height: 100%; padding: 0px">
    <el-header style="height: auto; margin-top: 5px">
        <el-row type="flex" align="middle" :gutter="5">
            <el-col :span="2" style="text-align: right"><span>Mode d'affichage:</span></el-col>
            <el-col :span="5">
                <el-radio-group v-model="context.Mode" @change="ChangeMode" size="mini">
                    <el-tooltip content="Création de poteaux" placement="bottom" effect="light" open-delay="500">
                        <el-radio-button label="creation">Création</el-radio-button>
                    </el-tooltip>
                    <el-tooltip content="Planification d'activité" placement="bottom" effect="light" open-delay="500">
                        <el-radio-button label="followup">Planification</el-radio-button>
                    </el-tooltip>
                    <el-tooltip content="Mise a jour de l'avancement" placement="bottom" effect="light" open-delay="500">
                        <el-radio-button label="billing">Avancement</el-radio-button>
                    </el-tooltip>
                </el-radio-group>
            </el-col> 
            <el-col :offset="0" :span="1">
                <el-popover
                        v-if="context.Mode == 'billing' && user.Permissions.Invoice"
                        placement="bottom" title="Passage en Attachement"
                        trigger="click"
                        width="400"
                        v-model="attachmentVisible"
                >
					<div style="margin: 10px 0 5px">Intervale d'activité : <span v-if="attachmentApplied > 0" style="color: dodgerblue">{{attachmentApplied}} éléments concernés</span></div>
					<el-checkbox v-model="attachmentOverride" size="mini" @change="CountPoleInAttachmentRange">Inclure les éléments déjà attachés</el-checkbox>
					<el-date-picker
							v-model="attachmentRange"
							type="daterange" unlink-panels size="mini" style="width: 100%"
							:picker-options="{firstDayOfWeek:1}" format="dd/MM/yyyy"
							value-format="yyyy-MM-dd"
							range-separator="à"
							start-placeholder="Début"
							end-placeholder="Fin"
							@change="CountPoleInAttachmentRange">
					></el-date-picker>
					<div style="margin: 10px 0 5px">Date de l'attachement :</div>
                    <el-date-picker
                            format="dd/MM/yyyy" size="mini" v-model="attachmentDate"
                            style="width: 100%" type="date"
                            value-format="yyyy-MM-dd"
                            placeholder="Date">
                    </el-date-picker>
					<!-- :picker-options="{firstDayOfWeek:1, disabledDate(time) { return time.getTime() > Date.now(); }}" -->
                    <div style="text-align: right; margin: 15px 0px 0px 0px">
                        <el-button size="mini" type="text" @click="attachmentVisible = false">Annuler</el-button>
                        <el-button size="mini" type="primary"  :disabled="IsAttachmentDisabled" @click="SetAttachments()">Appliquer</el-button>
                    </div>

                    <el-tooltip slot="reference" content="Attachements" placement="bottom" effect="light" open-delay=500>
                        <el-button type="primary" plain class="icon" icon="fas fa-paperclip icon--medium" size="mini" :disabled="attachmentVisible"></el-button>
                    </el-tooltip>
                </el-popover>
            </el-col> 
        </el-row>
    </el-header>
    <div style="height: 100%;overflow-x: hidden;overflow-y: auto;padding: 0px 0px; margin-top: 8px">
		<pole-table-creation v-if="context.Mode == 'creation'"
				:user="user"
				:polesite="polesite"
				:filter="filter"
				:filtertype="filtertype"
				:context.sync="context"
				@update:context="ChangeMode"
		></pole-table-creation>
		<pole-table-followup v-if="context.Mode == 'followup'"
				:user="user"
				:polesite="polesite"
				:filter="filter"
				:filtertype="filtertype"
				:context.sync="context"
				@update:context="ChangeMode"
		></pole-table-followup>
		<pole-table-billing v-if="context.Mode == 'billing'"
				:user="user"
				:polesite="polesite"
				:filter="filter"
				:filtertype="filtertype"
				:context.sync="context"
				@update:context="ChangeMode"
		></pole-table-billing>
    </div>
</el-container>
`

//@pole-selected="SetSelectedPole"

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Registration

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("pole-table", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		registerComponentTable("creation"),
		registerComponentTable("followup"),
		registerComponentTable("billing"),
		hvue.Template(template),
		hvue.Props("user", "polesite", "filter", "filtertype", "context"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewPoleTablesModel(vm)
		}),
		hvue.Computed("attachmentApplied", func(vm *hvue.VM) interface{} {
			ptm := &PoleTablesModel{Object: vm.Object}
			return ptm.CountPoleInAttachmentRange()
		}),
		hvue.Computed("IsAttachmentDisabled", func(vm *hvue.VM) interface{} {
			ptm := &PoleTablesModel{Object: vm.Object}
			return ptm.CheckAttachmentDisabled()
		}),
		hvue.MethodsOf(&PoleTablesModel{}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type PoleTablesModel struct {
	*js.Object

	Polesite   *ps.Polesite `js:"polesite"`
	User       *fm.User     `js:"user"`
	Filter     string       `js:"filter"`
	FilterType string       `js:"filtertype"`
	Context    *Context     `js:"context"`

	AttachmentVisible  bool     `js:"attachmentVisible"`
	AttachmentRange    []string `js:"attachmentRange"`
	AttachmentDate     string   `js:"attachmentDate"`
	AttachmentOverride bool     `js:"attachmentOverride"`

	VM *hvue.VM `js:"VM"`
}

func NewPoleTablesModel(vm *hvue.VM) *PoleTablesModel {
	rtm := &PoleTablesModel{Object: tools.O()}
	rtm.Polesite = ps.NewPolesite()
	rtm.User = fm.NewUser()
	rtm.Filter = ""
	rtm.FilterType = ""
	rtm.Context = NewContext("")

	rtm.AttachmentVisible = false
	rtm.AttachmentDate = date.TodayAfter(0)
	rtm.AttachmentRange = []string{"", ""}
	rtm.AttachmentOverride = false

	rtm.VM = vm
	return rtm
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Actions related Methods

func (ptm *PoleTablesModel) ChangeMode(vm *hvue.VM) {
	ptm = &PoleTablesModel{Object: vm.Object}
	vm.Emit("update:context", ptm.Context)
}

func (ptm *PoleTablesModel) CheckAttachmentDisabled() bool {
	if tools.Empty(ptm.AttachmentDate) {
		return true
	}
	if tools.Empty(ptm.AttachmentRange[0]) && tools.Empty(ptm.AttachmentRange[0]) {
		return true
	}
	return false
}

func (ptm *PoleTablesModel) IsPoleInAttachmentRange(pole *ps.Pole) bool {
	if !(!tools.Empty(pole.Date) && pole.Date >= ptm.AttachmentRange[0] && pole.Date <= ptm.AttachmentRange[1]) {
		return false
	}
	return !(!pole.IsDone() && !(ptm.AttachmentOverride && pole.IsAttachment()))
}

func (ptm *PoleTablesModel) CountPoleInAttachmentRange() int {
	if ptm.CheckAttachmentDisabled() {
		return 0
	}
	nbApplied := 0
	for _, pole := range ptm.Polesite.Poles {
		if !ptm.IsPoleInAttachmentRange(pole) {
			continue
		}
		nbApplied++
	}
	return nbApplied
}

func (ptm *PoleTablesModel) SetAttachments(vm *hvue.VM) {
	ptm = &PoleTablesModel{Object: vm.Object}
	nbApplied := 0
	for _, pole := range ptm.Polesite.Poles {
		if !ptm.IsPoleInAttachmentRange(pole) {
			continue
		}
		pole.SetAttachmentDate(ptm.AttachmentDate)
		nbApplied++
	}
	msg := "Attachement appliqué sur " + strconv.Itoa(nbApplied) + " élément"
	if nbApplied > 1 {
		msg += "s"
	}
	message.NotifySuccess(ptm.VM, "Passage en attachement", msg)
	ptm.AttachmentVisible = false
}
