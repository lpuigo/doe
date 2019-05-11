package ripsiteupdatemodal

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
				<h2 v-if="current_ripsite" style="margin: 0 0">
					<i class="far fa-edit icon--left"></i>Edition du chantier: <span style="color: #ccebff">{{current_ripsite.Ref}}</span>
				</h2>
			</el-col>
			<el-col :span="6">
				<ripsite-info :ripsite="current_ripsite"></ripsite-info>
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
	<el-container v-loading="loading" style="height: 65vh;overflow-x: hidden;overflow-y: auto;padding: 6px 6px;">
        <el-header style="height: auto; padding: 0px 0px">
            <el-row :gutter="10" style="margin-bottom: 10px">
<!--                GroupButton for Pulling / Junction / Measurement-->
                <el-col :span="4">
                    <el-radio-group v-model="ActivityMode" size="mini">
                        <el-radio-button v-if="current_ripsite.Pullings.length > 0" label="Pulling">Tirage: {{current_ripsite.Pullings.length}}</el-radio-button>
                        <el-radio-button v-if="current_ripsite.Junctions.length > 0" label="Junction">Racco: {{current_ripsite.Junctions.length}}</el-radio-button>
                        <el-radio-button v-if="current_ripsite.Measurements.length > 0" label="Measurement">Mesure: {{current_ripsite.Measurements.length}}</el-radio-button>
                    </el-radio-group>
                </el-col>
                <el-col :offset="1" :span="2" >
                    <span style="float:right; text-align: right">Commentaire dossier:</span>
                </el-col>
                <el-col :span="12">
                    <el-input type="textarea" autosize placeholder="Commentaire sur le chantier" size="mini"
                              v-model="current_ripsite.Comment"
                    ></el-input>
                </el-col>
                <el-col :offset="1" :span="4">
                    <el-input
                            placeholder="filtre"
                            prefix-icon="el-icon-search"
                            v-model="filter"
                            size="mini"	clearable
                    ></el-input>
                </el-col>
            </el-row>
        </el-header>
        <el-main style="height: 100%; padding: 0px">
            <rip-pulling-update v-if="ActivityMode == 'Pulling'" v-model="current_ripsite" :user="User" :filter="filter"></rip-pulling-update>
            <rip-junction-update v-if="ActivityMode == 'Junction'" v-model="current_ripsite" :user="User" :filter="filter"></rip-junction-update>
            <rip-measurement-update v-if="ActivityMode == 'Measurement'" v-model="current_ripsite" :user="User" :filter="filter"></rip-measurement-update>
        </el-main>

	</el-container>
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
						<el-button :loading="saving" size="mini" type="primary" @click="DeleteRipsite">Oui</el-button>
					</div>
				</el-popover>

				<el-tooltip effect="light" :open-delay="500">
					<div slot="content">Supprimer ce chantier</div>
					<el-button :loading="saving" :disabled="isNewRipsite" type="danger" plain size="mini" icon="far fa-trash-alt" v-popover:confirm_delete_popover></el-button>
				</el-tooltip>
				
<!--				<el-tooltip effect="light" :open-delay="500">-->
<!--					<div slot="content">Dupliquer ce chantier</div>-->
<!--					<el-button :loading="saving" :disabled="isNewRipsite" type="info" plain size="mini" icon="far fa-clone" @click="Duplicate"></el-button>-->
<!--				</el-tooltip>-->
				
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
