package worksiteeditmodal

const template string = `
<el-dialog 
		:visible.sync="visible" 
		width="90%"
		:before-close="HideWithControl"
>
	<!-- 
		Modal Title
	-->
    <span slot="title">
		<el-row :gutter="10" type="flex" align="middle">
			<el-col :span="10">
				<h2 v-if="current_worksite" style="margin: 0 0">
					<i class="far fa-edit icon--left"></i>Edition du chantier: <span style="color: #ccebff">{{current_worksite.City}} - {{current_worksite.Ref}}</span>
				</h2>
			</el-col>
			<el-col :span="6">
				<worksite-info :worksite="current_worksite"></worksite-info>
			</el-col>		
		</el-row>
    </span>

	<!-- 
		Modal Body
	-->
	<div style="max-height: 70vh;overflow-x: hidden;overflow-y: auto;padding-right: 6px;">
		<worksite-detail
				:worksite="current_worksite"
				:readonly="false"
		></worksite-detail>

	</div>

	<!-- 
		Body Action Bar
	-->	
	<span slot="footer">
		<el-row :gutter="15">
			<el-col :span="24" style="text-align: right">
				<el-popover
						ref="confirm_delete_popover"
						placement="top"
						width="160"
						v-model="showconfirmdelete"
				>
					<p>Supprimer ce chantier ?</p>
					<div style="text-align: left; margin: 0;">
						<el-button size="mini" type="text" @click="showconfirmdelete = false">Non</el-button>
						<el-button size="mini" type="primary" @click="DeleteWorksite">Oui</el-button>
					</div>
				</el-popover>

				<el-tooltip effect="light" :open-delay="500">
					<div slot="content">Supprimer ce chantier</div>
					<el-button :disabled="isNewWorksite" type="danger" plain size="mini" icon="far fa-trash-alt" v-popover:confirm_delete_popover></el-button>
				</el-tooltip>
				<el-tooltip effect="light" :open-delay="500">
					<div slot="content">Dupliquer ce chantier</div>
					<el-button :disabled="isNewWorksite" type="info" plain size="mini" icon="far fa-clone" @click="Duplicate"></el-button>
				</el-tooltip>
				<el-tooltip effect="light" :open-delay="500">
					<div slot="content">Annuler les changements</div>
					<el-button :disabled="!hasChanged" type="info" plain size="mini" icon="fas fa-undo-alt" @click="UndoChange"></el-button>
				</el-tooltip>
				<el-button @click="Hide" size="mini">Fermer</el-button>
				<el-button :type="hasWarning" plain size="mini" :disabled="!hasChanged" @click="ConfirmChange">
					<span v-if="!isNewWorksite">Enregistrer</span>
					<span v-else>Create New</span>
				</el-button>
			</el-col>
		</el-row>
	</span>
</el-dialog>`
