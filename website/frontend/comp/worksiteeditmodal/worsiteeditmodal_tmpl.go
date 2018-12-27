package worksiteeditmodal

const template string = `
<el-dialog 
		:visible.sync="visible" 
		width="80%"
		:before-close="Hide"
>
	<!-- 
		Modal Header
	-->
    <span slot="title" class="ewin_light">
        <h2 v-if="current_worksite" style="margin: 0 0">
        	<i class="far fa-edit icon--left"></i>Edition du chantier: <span style="color: teal">{{current_worksite.ref}}</span>
        </h2>
    </span>
	<!-- 
		Modal Body
	-->
	<div style="max-height: 60vh;overflow-x: hidden;overflow-y: auto;padding-right: 8px;">
		<!-- 
			Body Content
		-->
		<worksite-detail
				:worksite="current_worksite"
				:readonly="false"
		></worksite-detail>
		<!-- 
			Body Action Bar
		-->
		<el-row :gutter="15" class="form-row">
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
					<el-button :disabled="isNewWorksite" type="danger" plain icon="far fa-trash-alt" v-popover:confirm_delete_popover></el-button>
				</el-tooltip>
				<el-tooltip effect="light" :open-delay="500">
					<div slot="content">Dupliquer ce chantier</div>
					<el-button :disabled="isNewWorksite" type="info" plain icon="far fa-clone" @click="Duplicate"></el-button>
				</el-tooltip>
				<el-button @click="Hide">Fermer</el-button>
				<el-button :type="hasWarning" plain :disabled="!hasChanged" @click="ConfirmChange">
					<span v-if="!isNewWorksite">Enregistrer</span>
					<span v-else>Create New</span>
				</el-button>
			</el-col>
		</el-row>
	</div>
</el-dialog>`
