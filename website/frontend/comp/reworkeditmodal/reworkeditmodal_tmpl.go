package reworkeditmodal

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
					<i class="fas fa-tools icon--left"></i>Création de Reprise: <span style="color: #ccebff">{{current_worksite.City}} - {{current_worksite.Ref}}</span>
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
	<div v-loading="loading" style="height: 65vh;">
        <el-tabs v-if="!loading" 
                 v-model="activeTabName" 
                 tab-position="left" type="border-card"
                 :stretch="true"
                 style="height: 100%;"
        >
            <el-tab-pane v-if="user.Permissions.Create" label="Création" name="Create">
                <rework-edit v-if="HasRework"
                             style="height: 65vh;overflow-x: hidden;overflow-y: auto;padding-right: 6px;"
                             :worksite="current_worksite"
                             :user="user"
                ></rework-edit>
            </el-tab-pane>
            <el-tab-pane v-if="user.Permissions.Update" label="Maj" name="Update">
                <rework-update v-if="HasRework"
                             style="height: 65vh;overflow-x: hidden;overflow-y: auto;padding-right: 6px;"
                             :worksite="current_worksite"
                             :user="user"
                ></rework-update>
            </el-tab-pane>
        </el-tabs>
	</div>

	<!-- 
		Body Action Bar
	-->	
	<span slot="footer">
		<el-row :gutter="15">
			<el-col :span="24" style="text-align: right">
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
