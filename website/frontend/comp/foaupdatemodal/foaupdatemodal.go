package foaupdatemodal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/modal"
	"github.com/lpuig/ewin/doe/website/frontend/comp/ripstateupdate"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	fmfoa "github.com/lpuig/ewin/doe/website/frontend/model/foasite"
	"github.com/lpuig/ewin/doe/website/frontend/model/foasite/foaconst"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/date"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
	"strconv"
)

const template string = `
<el-dialog
		:visible.sync="visible" 
		width="80%"
		:before-close="HideWithControl"
>
	<!-- 
		Modal Title
	-->
    <span slot="title">
		<el-row :gutter="10" type="flex" align="middle">
			<el-col :span="24">
				<h2 style="margin: 0 0">
					<i class="far fa-edit icon--left"></i>Mise à jour:
				</h2>
			</el-col>
		</el-row>
    </span>

	<!-- 
		Modal Body
		style="height: 100%;"		
	-->
	<el-container style="padding: 6px 6px;">
		<el-row type="flex" align="middle" :gutter="10" style="width: 100%">
			<!-- Actors -->
			<el-col :span="6">
				<el-select v-model="CurrentState.Actors" filterable multiple placeholder="Acteurs" size="mini" style="width: 100%"
						   @clear="UpdateStatus()"
						   @change="UpdateStatus()"
						   :disabled="DisableTeam()"
				>
					<el-option
							v-for="item in GetActors()"
							:key="item.value"
							:label="item.label"
							:value="item.value"
							:disabled="item.disabled"
					>
					</el-option>
				</el-select>
			</el-col>
		
			<!-- Date -->
			<el-col :span="4">
				<el-date-picker format="dd/MM/yyyy" placeholder="Date" size="mini"
								style="width: 100%" type="date"
								v-model="CurrentState.Date"
								value-format="yyyy-MM-dd"
								:picker-options="{firstDayOfWeek:1, disabledDate(time) { return time.getTime() > Date.now(); }}"
								:disabled="DisableDates()" :clearable="true"
								@change="UpdateStatus()"
				></el-date-picker>
			</el-col>
		
			<!-- Status -->
			<el-col :span="4">
				<el-select v-model="CurrentState.Status" filterable
						   size="mini" style="width: 100%"
						   placeholder="Etat"
						   @clear=""
						   @change="UpdateStatus"
				>
					<el-option
							v-for="item in GetStatuses()"
							:key="item.value"
							:label="item.label"
							:value="item.value">
					</el-option>
				</el-select>
			</el-col>
		
			<!-- Comment -->
			<el-col :span="10">
				<el-input type="textarea" autosize placeholder="Commentaire" size="mini"
						  v-model="CurrentState.Comment"
				></el-input>
			</el-col>
		</el-row>
	</el-container>
	<!-- 
		Body Action Bar
	-->
		
	<span slot="footer">
		<el-row :gutter="15">
			<el-col :span="24" style="text-align: right">
<!--				<el-tooltip effect="light" :open-delay="500">-->
<!--					<div slot="content">Dupliquer ce chantier</div>-->
<!--					<el-button :loading="saving" :disabled="isNewRipsite" type="info" plain size="mini" icon="far fa-clone" @click="Duplicate"></el-button>-->
<!--				</el-tooltip>-->

				<!--				
				<el-tooltip effect="light" :open-delay="500">
					<div slot="content">Annuler les changements</div>
					<el-button :loading="saving" :disabled="!hasChanged" type="info" plain size="mini" icon="fas fa-undo-alt" @click="UndoChange"></el-button>
				</el-tooltip>
				-->
				<el-button size="mini" @click="Hide">Fermer</el-button>
				<el-button type="primary" plain size="mini" @click="ApplyChange">Appliquer</el-button>
			</el-col>
		</el-row>
	</span>
</el-dialog>`

//////////////////////////////////////////////////////////////////////////////////////////////
// Model Methods

type FoaUpdateModalModel struct {
	*modal.ModalModel

	Foas   []*fmfoa.Foa `js:"Foa"`
	User   *fm.User     `js:"user"`
	Client string       `js:"client"`

	CurrentState *fmfoa.State `js:"CurrentState"`
}

func NewFoaUpdateModalModel(vm *hvue.VM) *FoaUpdateModalModel {
	fumm := &FoaUpdateModalModel{ModalModel: modal.NewModalModel(vm)}

	fumm.Foas = []*fmfoa.Foa{}
	fumm.User = fm.NewUser()
	fumm.Client = ""

	fumm.CurrentState = fmfoa.NewState()
	fumm.CurrentState.Date = date.TodayAfter(0)

	return fumm
}

func FoaUpdateModalModelFromJS(o *js.Object) *FoaUpdateModalModel {
	return &FoaUpdateModalModel{ModalModel: &modal.ModalModel{Object: o}}
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Component Methods

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("foa-update-modal", componentOption()...)
}

func componentOption() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		ripstateupdate.RegisterComponent(),
		hvue.Props("user", "client"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewFoaUpdateModalModel(vm)
		}),
		hvue.MethodsOf(&FoaUpdateModalModel{}),
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Modal Methods

func (fumm *FoaUpdateModalModel) Show(foas *fmfoa.FoaSite) {
	fumm.Foas = foas.Foas
	fumm.ModalModel.Show()
}

func (fumm *FoaUpdateModalModel) HideWithControl() {
	fumm.Hide()
}

//////////////////////////////////////////////////////////////////////////////////////////////
// HTML Methods

func (fumm *FoaUpdateModalModel) UpdateStatus(vm *hvue.VM) {
	fumm = FoaUpdateModalModelFromJS(vm.Object)
	if fumm.CurrentState.IsCanceled() {
		fumm.CurrentState.Date = ""
		fumm.CurrentState.Actors = []string{}
		return
	}
	if !tools.Empty(fumm.CurrentState.Date) && len(fumm.CurrentState.Actors) > 0 {
		fumm.CurrentState.Status = foaconst.StateDone
	}
}

func (fumm *FoaUpdateModalModel) ApplyChange(vm *hvue.VM) {
	fumm = FoaUpdateModalModelFromJS(vm.Object)
	for _, foa := range fumm.Foas {
		foa.State.Copy(fumm.CurrentState)
	}
	fumm.Hide()
}

func (fumm *FoaUpdateModalModel) DisableDates(vm *hvue.VM) bool {
	fumm = FoaUpdateModalModelFromJS(vm.Object)
	if !fumm.CurrentState.IsCanceled() && len(fumm.CurrentState.Actors) > 0 {
		return false
	}
	return true
}

func (fumm *FoaUpdateModalModel) DisableTeam(vm *hvue.VM) bool {
	fumm = FoaUpdateModalModelFromJS(vm.Object)
	return fumm.CurrentState.IsCanceled()
}

func (fumm *FoaUpdateModalModel) GetActors(vm *hvue.VM) []*elements.ValueLabelDisabled {
	fumm = FoaUpdateModalModelFromJS(vm.Object)
	client := fumm.User.GetClientByName(fumm.Client)
	if client == nil {
		return nil
	}

	res := []*elements.ValueLabelDisabled{}
	for _, actor := range client.Actors {
		res = append(res, elements.NewValueLabelDisabled(strconv.Itoa(actor.Id), actor.GetRef(), !actor.Active))
	}
	return res
}

func (fumm *FoaUpdateModalModel) GetStatuses(vm *hvue.VM) []*elements.ValueLabel {
	fumm = FoaUpdateModalModelFromJS(vm.Object)
	return fmfoa.GetStatesValueLabel()
}
