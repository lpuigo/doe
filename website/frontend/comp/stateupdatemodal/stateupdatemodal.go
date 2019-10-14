package stateupdatemodal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/modal"
	"github.com/lpuig/ewin/doe/website/frontend/comp/ripstateupdate"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	fmrip "github.com/lpuig/ewin/doe/website/frontend/model/ripsite"
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
	<el-container v-loading="loading" style="padding: 6px 6px;">
            <rip-state-update style="width: 100%" v-model="State" :user="user" :client="client"></rip-state-update>
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
				<el-button @click="Hide" size="mini">Fermer</el-button>
				
				<!--
				<el-button :loading="saving" :type="hasWarning" plain size="mini" :disabled="!hasChanged" @click="ConfirmChange"
				>Enregistrer</el-button>
				-->
			</el-col>
		</el-row>
	</span>
</el-dialog>`

type StateUpdateModalModel struct {
	*modal.ModalModel

	State  *fmrip.State `js:"State"`
	User   *fm.User     `js:"user"`
	Client string       `js:"client"`
}

func NewStateUpdateModalModel(vm *hvue.VM) *StateUpdateModalModel {
	summ := &StateUpdateModalModel{ModalModel: modal.NewModalModel(vm)}

	summ.State = fmrip.NewState()
	summ.User = fm.NewUser()
	summ.Client = ""

	return summ
}

func StateUpdateModalModelFromJS(o *js.Object) *StateUpdateModalModel {
	return &StateUpdateModalModel{ModalModel: &modal.ModalModel{Object: o}}
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Component Methods

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("state-update-modal", componentOption()...)
}

func componentOption() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		ripstateupdate.RegisterComponent(),
		hvue.Props("user", "client"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewStateUpdateModalModel(vm)
		}),
		hvue.MethodsOf(&StateUpdateModalModel{}),
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Modal Methods

func (summ *StateUpdateModalModel) Show(state *fmrip.State) {
	summ.State = state
	summ.ModalModel.Show()
}

func (summ *StateUpdateModalModel) HideWithControl() {
	summ.Hide()
}
