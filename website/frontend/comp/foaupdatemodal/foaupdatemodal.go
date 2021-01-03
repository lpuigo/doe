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
)

const template string = `
<el-dialog
		:visible.sync="visible" 
		width="80%"
		:before-close="Hide"
>
	<!-- 
		Modal Title
	-->
    <span slot="title">
		<el-row :gutter="10" type="flex" align="middle">
			<el-col :span="24">
				<h2 style="margin: 0 0">
					<i class="far fa-edit icon--left"></i>
					<span v-if="Foas.length > 1">Mise à jour de {{Foas.length}} FOAs</span>
					<span v-else>Mise à jour de 1 FOA</span>
				</h2>
			</el-col>
		</el-row>
    </span>

	<!-- 
		Modal Body
		style="height: 100%;"		
	-->
	<div style="padding: 6px 6px;">
		<div v-if="EditMode">
			<!-- Insee & Ref & Type-->
			<el-row :gutter="10" type="flex" align="middle" class="doublespaced">
				<el-col :span="2" class="align-right">Code Insee :</el-col>
				<el-col :span="5">
					<el-input placeholder="Insee"
							  v-model="EditedFoa.Insee" clearable size="mini"
					></el-input>
				</el-col>
		
				<el-col :span="3" class="align-right">Référence de la chambre :</el-col>
				<el-col :span="5">
					<el-input placeholder="Référence"
							  v-model="EditedFoa.Ref" clearable size="mini"
					></el-input>
				</el-col>
		
				<el-col :span="3" class="align-right">Type de chambre :</el-col>
				<el-col :span="5">
					<el-input placeholder="Type"
							  v-model="EditedFoa.Type" clearable size="mini"
					></el-input>
				</el-col>
			</el-row>
		</div>
		<div v-else>
			<el-row :gutter="10" type="flex" align="middle" class="doublespaced">
				<el-col :span="3" class="align-right">
						<span v-if="Foas.length > 1">FOAs à modifier:</span>
						<span v-else>FOA à modifier</span>
				</el-col>
				<el-col :span="20">
					<span>{{SelectedFOAList()}}</span>
				</el-col>
			</el-row>
		</div>

		<el-row :gutter="10" type="flex" align="middle" class="doublespaced">
			<el-col :span="4">					
				<span v-if="Foas.length > 1">Mise à jour des états :</span>
				<span v-else>Mise à jour de l'état :</span>
			</el-col>
		</el-row>

		<el-row type="flex" align="middle" :gutter="10" style="width: 100%" class="spaced">
			<!-- Actors -->
			<el-col :span="6">
				<el-select v-model="EditedFoa.State.Actors" filterable multiple placeholder="Acteurs" size="mini" style="width: 100%"
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
								v-model="EditedFoa.State.Date"
								value-format="yyyy-MM-dd"
								:picker-options="{firstDayOfWeek:1, disabledDate(time) { return time.getTime() > Date.now(); }}"
								:disabled="DisableDates()" :clearable="true"
								@change="UpdateStatus()"
				></el-date-picker>
			</el-col>
		
			<!-- Status -->
			<el-col :span="4">
				<el-select v-model="EditedFoa.State.Status" filterable
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
						  v-model="EditedFoa.State.Comment"
				></el-input>
			</el-col>
		</el-row>
	</div>
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

	Foas      []*fmfoa.Foa `js:"Foas"`
	User      *fm.User     `js:"user"`
	Client    string       `js:"client"`
	EditMode  bool         `js:"EditMode"`
	EditedFoa *fmfoa.Foa   `js:"EditedFoa"`
	OnApply   func()       `js:"OnApply"`
}

func NewFoaUpdateModalModel(vm *hvue.VM) *FoaUpdateModalModel {
	fumm := &FoaUpdateModalModel{ModalModel: modal.NewModalModel(vm)}

	fumm.Foas = []*fmfoa.Foa{}
	fumm.User = fm.NewUser()
	fumm.Client = ""
	fumm.EditMode = false
	fumm.EditedFoa = fmfoa.NewFoa()
	fumm.EditedFoa.State.Date = date.TodayAfter(0)
	fumm.OnApply = func() {}

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

func (fumm *FoaUpdateModalModel) SetModel(foa *fmfoa.Foa) {
	fumm.EditedFoa = foa.Clone()
}

func (fumm *FoaUpdateModalModel) Show(selectedFoas *fmfoa.FoaSite, OnApply func()) {
	fumm.EditMode = false
	fumm.Foas = selectedFoas.Foas
	fumm.ModalModel.Show()
	fumm.OnApply = OnApply
}

func (fumm *FoaUpdateModalModel) ShowEdit(foa *fmfoa.Foa, OnApply func()) {
	fumm.EditMode = true
	fumm.EditedFoa = foa.Clone()
	fumm.Foas = []*fmfoa.Foa{foa}
	fumm.ModalModel.Show()
	fumm.OnApply = OnApply
}

func (fumm *FoaUpdateModalModel) Hide() {
	fumm.EditMode = false
	fumm.ModalModel.Hide()
}

//////////////////////////////////////////////////////////////////////////////////////////////
// HTML Methods

func (fumm *FoaUpdateModalModel) ApplyChange(vm *hvue.VM) {
	fumm = FoaUpdateModalModelFromJS(vm.Object)
	if fumm.EditMode {
		fumm.Foas[0].Copy(fumm.EditedFoa)
	} else {
		for _, foa := range fumm.Foas {
			foa.State.Copy(fumm.EditedFoa.State)
		}
	}
	fumm.OnApply()
	fumm.Hide()
}

func (fumm *FoaUpdateModalModel) UpdateStatus(vm *hvue.VM) {
	fumm = FoaUpdateModalModelFromJS(vm.Object)
	if fumm.EditedFoa.State.IsCanceled() {
		fumm.EditedFoa.State.Date = ""
		fumm.EditedFoa.State.Actors = []string{}
		return
	}
	if !tools.Empty(fumm.EditedFoa.State.Date) && len(fumm.EditedFoa.State.Actors) > 0 {
		fumm.EditedFoa.State.Status = foaconst.StateDone
	}
}

func (fumm *FoaUpdateModalModel) DisableDates(vm *hvue.VM) bool {
	fumm = FoaUpdateModalModelFromJS(vm.Object)
	if !fumm.EditedFoa.State.IsCanceled() && len(fumm.EditedFoa.State.Actors) > 0 {
		return false
	}
	return true
}

func (fumm *FoaUpdateModalModel) DisableTeam(vm *hvue.VM) bool {
	fumm = FoaUpdateModalModelFromJS(vm.Object)
	return fumm.EditedFoa.State.IsCanceled()
}

func (fumm *FoaUpdateModalModel) GetActors(vm *hvue.VM) []*elements.ValueLabelDisabled {
	fumm = FoaUpdateModalModelFromJS(vm.Object)
	client := fumm.User.GetClientByName(fumm.Client)
	if client == nil {
		return nil
	}

	res := []*elements.ValueLabelDisabled{}
	for _, actor := range client.Actors {
		res = append(res, actor.GetElementsValueLabelDisabled())
	}
	return res
}

func (fumm *FoaUpdateModalModel) GetStatuses(vm *hvue.VM) []*elements.ValueLabel {
	fumm = FoaUpdateModalModelFromJS(vm.Object)
	return fmfoa.GetStatesValueLabel()
}

func (fumm *FoaUpdateModalModel) SelectedFOAList(vm *hvue.VM) string {
	fumm = FoaUpdateModalModelFromJS(vm.Object)
	res := ""
	last := len(fumm.Foas) - 1
	for i, foa := range fumm.Foas {
		res += foa.Insee + " - " + foa.Ref + " (" + foa.Type + ")"
		if i != last {
			res += ", "
		}
	}
	return res
}
