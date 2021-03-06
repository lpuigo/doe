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
			<el-col :span="12">
				<h2 v-if="current_worksite" style="margin: 0 0">
					<i class="far fa-edit icon--left"></i>Edition du chantier: <span style="color: #ccebff">{{current_worksite.City}} - {{current_worksite.Ref}}</span>
				</h2>
			</el-col>
			<el-col :span="6">
				<worksite-info :worksite="current_worksite"></worksite-info>
			</el-col>	
            <el-col :offset="2" :span="1">
                <h2 style="margin: 0 0"><a :href="Attachment()"><i class="link fas fa-file-excel"></i></a></h2>
            </el-col>
		</el-row>
    </span>

	<!-- 
		Modal Body
		style="height: 100%;"
		
	-->
	<div v-loading="loading" style="height: 65vh;">
        <el-tabs v-if="!loading" 
                 v-model="activeTabName" 
                 tab-position="left" type="border-card"
                 :stretch="true"
                 style="height: 100%;"
        >
            <el-tab-pane v-if="user.Permissions.Create" label="Création" name="Create">
                <worksite-edit style="height: 65vh;overflow-x: hidden;overflow-y: auto;padding-right: 6px;"
                               :worksite="current_worksite"
                               :readonly="false"
                               :user="user"
                ></worksite-edit>
            </el-tab-pane>
            <el-tab-pane v-if="user.Permissions.Update" label="Maj" name="Update">
                <worksite-update style="height: 65vh;overflow-x: hidden;overflow-y: auto;padding-right: 6px;"
                               :worksite="current_worksite"
							   :user="user"
                ></worksite-update>
            </el-tab-pane>
        </el-tabs>
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
						<el-button :loading="saving" size="mini" type="text" @click="showconfirmdelete = false">Non</el-button>
						<el-button :loading="saving" size="mini" type="primary" @click="DeleteWorksite">Oui</el-button>
					</div>
				</el-popover>

				<el-tooltip effect="light" :open-delay="500">
					<div slot="content">Supprimer ce chantier</div>
					<el-button :loading="saving" :disabled="isNewWorksite" type="danger" plain size="mini" icon="far fa-trash-alt" v-popover:confirm_delete_popover></el-button>
				</el-tooltip>
				
				<el-tooltip effect="light" :open-delay="500">
					<div slot="content">Dupliquer ce chantier</div>
					<el-button :loading="saving" :disabled="isNewWorksite" type="info" plain size="mini" icon="far fa-clone" @click="Duplicate"></el-button>
				</el-tooltip>
				
				<el-tooltip effect="light" :open-delay="500">
					<div slot="content">Annuler les changements</div>
					<el-button :loading="saving" :disabled="!hasChanged" type="info" plain size="mini" icon="fas fa-undo-alt" @click="UndoChange"></el-button>
				</el-tooltip>
				
				<el-button @click="Hide" size="mini">Fermer</el-button>
				
				<el-button :loading="saving" :type="hasWarning" plain size="mini" :disabled="!hasChanged" @click="ConfirmChange"
				>Enregistrer</el-button>
			</el-col>
		</el-row>
	</span>
</el-dialog>`
