package reworkupdatemodal

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
					<i class="fas fa-tools icon--left"></i>Réalisation de Reprise: <span style="color: #ccebff">{{current_worksite.City}} - {{current_worksite.Ref}}</span>
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
	<div v-loading="loading" style="height: 65vh;overflow-x: hidden;overflow-y: auto;padding-right: 6px;">
        <el-container 
                v-if="HasRework" 
                style="height: 100%"
        >
            <el-header style="height: auto; padding: 5px">
                <el-row :gutter="10" type="flex" align="middle" style="margin-bottom: 10px">
                    <el-col :offset="2" :span="2">
                        <span style="float:right">Contrôle:</span>
                    </el-col>
                    <el-col :span="4">
                        <el-date-picker format="dd/MM/yyyy" placeholder="Contrôle" size="mini"
                                        style="width: 100%" type="date"
                                        v-model="current_worksite.Rework.ControlDate"
                                        value-format="yyyy-MM-dd"
                                        :picker-options="{firstDayOfWeek:1, disabledDate(time) { return time.getTime() > Date.now(); }}"
                                        :disabled="true"
                        ></el-date-picker>
                    </el-col>
                    <el-col :span="2">
                        <span style="float:right">Soumission:</span>
                    </el-col>
                    <el-col :span="4">
                        <el-date-picker format="dd/MM/yyyy" placeholder="Soumission" size="mini"
                                        style="width: 100%" type="date"
                                        v-model="current_worksite.Rework.SubmissionDate"
                                        value-format="yyyy-MM-dd"
                                        :picker-options="{firstDayOfWeek:1, disabledDate(time) { return time.getTime() > Date.now(); }}"
                                        :disabled="true"
                        ></el-date-picker>
                    </el-col>
                    <el-col :offset="8" :span="4">
                        <worksite-status-tag v-model="current_worksite"></worksite-status-tag>
                    </el-col>
                </el-row>
            </el-header>
            <el-main  style="height: 100%; padding: 0px">
                <el-table
                        :data="filteredReworks"
                        :row-class-name="TableRowClassName"
                        height="100%" :border=true size="mini"
                >
                   
                    <el-table-column
                            label="PT" prop="PT"
                            width="150px" :resizable="true" :show-overflow-tooltip=true
                    >
                        <template slot-scope="scope">
                            <div>{{GetTronconRef(scope.row.PT)}}</div>
                            <div>{{GetTronconAddress(scope.row.PT)}}</div>
                        </template>
                    </el-table-column>
                    
                    <el-table-column
                            label="Nb Reprise"
                            width="70px" :resizable="true" :show-overflow-tooltip=true
                    >
                        <template slot-scope="scope">
                            <div>{{scope.row.NbKO}}</div>
                        </template>
                    </el-table-column>
                    
                    <el-table-column
                            label="Installation"
                            width="200px" :resizable="true" :show-overflow-tooltip=true
                    >
                        <template slot-scope="scope">
                            <div>{{GetTronconInstallDate(scope.row.PT)}}</div>
                            <div>{{GetTronconInstallActor(scope.row.PT)}}</div>
                        </template>
                    </el-table-column>

                    <el-table-column
                            label="Description" prop="Description"
                            width="400px" :resizable=true
                    ></el-table-column>
                    
                    <el-table-column
                            label="Reprise"   
                    >
                        <template slot-scope="scope">
                            <el-row type="flex" align="middle" :gutter="10">
                                <el-col :span="12">
                                    <el-autocomplete v-model="scope.row.FixActor"
                                                     :fetch-suggestions="UserSearch"
                                                     placeholder="Equipier"
                                                     prefix-icon="fas fa-user"
                                                     clearable size="mini" style="width: 100%"
                                                     @clear="scope.row.FixDate = ''"
                                    ></el-autocomplete>
                                </el-col>
                                <el-col :span="12">
                                    <el-date-picker format="dd/MM/yyyy" placeholder="Reprise" size="mini"
                                                    style="width: 100%" type="date"
                                                    v-model="scope.row.FixDate"
                                                    value-format="yyyy-MM-dd"
                                                    :picker-options="{firstDayOfWeek:1, disabledDate(time) { return time.getTime() > Date.now(); }}"
                                                    :disabled="!scope.row.FixActor" :clearable="false"
                                    ></el-date-picker>
                                </el-col>
                            </el-row>
                        </template>  
                    </el-table-column>
                    
                </el-table>
            </el-main>
        </el-container>
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
