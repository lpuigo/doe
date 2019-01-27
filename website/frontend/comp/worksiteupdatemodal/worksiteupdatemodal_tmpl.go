package worksiteupdatemodal

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
					<i class="fas fa-wrench icon--left"></i>Mise à jour du chantier: <span style="color: #ccebff">{{current_worksite.City}} - {{current_worksite.Ref}}</span>
				</h2>
			</el-col>
			<el-col :span="6">
				<worksite-info :worksite="current_worksite"></worksite-info>
			</el-col>		
            <el-col :offset="2" :span="2">
                <h2 style="margin: 0 0"><a :href="DOEArchive(current_worksite.Id)"><i class="link fas fa-file-archive icon--left"></i>DOE Zip</a></h2>
            </el-col>
        </el-row>
    </span>

	<!-- 
		Modal Body
	-->
	<div v-loading="loading" style="height: 65vh;overflow-x: hidden;overflow-y: auto;padding-right: 6px;">
        <el-container style="height: 100%">
            <el-header style="height: auto; padding: 5px">
                <el-row :gutter="10" style="margin-bottom: 10px">
                    <el-col :span="3">
                        <worksite-status-tag v-model="current_worksite"></worksite-status-tag>
                    </el-col>
                    <el-col :offset="1" :span="2" >
                        <span style="float:right; text-align: right">Commentaire dossier:</span>
                    </el-col>
                    <el-col :span="13">
                        <el-input clearable placeholder="Commentaire sur le dossier" size="mini" type="textarea" autosize
                                  v-model.trim="current_worksite.Comment"
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
            <el-main  style="height: 100%; padding: 0px">
                <!--:span-method="OrderSpanMethod"-->
                <el-table
                        :data="filteredTroncons"
                        :row-class-name="TableRowClassName"
                        :span-method="OrderSpanMethod"
                        height="100%" :border=true size="mini"
                >
                    <el-table-column
                            label="Commande" prop="Order"
                            width="100px" :resizable="true" :show-overflow-tooltip=true
                    ></el-table-column>
                    <el-table-column
                            label="Tronçon"
                            width="140px" :resizable="true" :show-overflow-tooltip=true
                    >
                        <template slot-scope="scope">
                            <div>{{scope.row | FormatTronconRef}}</div>
                            <div>{{scope.row.Pb.Address}}</div>
                        </template>
                    </el-table-column>
                    <el-table-column
                            label="Nb EL" prop="NbRacco"
                            width="50px" align="center"
                    ></el-table-column>
                    <el-table-column
                            label="Status"
                            width="120px" :resizable="true" :show-overflow-tooltip=true
                    >
                        <template slot-scope="scope">
                            <troncon-status-tag v-model="scope.row"></troncon-status-tag>
                        </template>
                    </el-table-column>
                    <el-table-column
                            label="Installation"
                            width="300px" min-width="250px" :resizable="true"
                    >
                        <template slot-scope="scope">
                            <el-row type="flex" align="middle" :gutter="10">
                                <el-col :span="12">
                                    <el-autocomplete v-model="scope.row.InstallActor"
                                                     :fetch-suggestions="UserSearch"
                                                     placeholder="Equipier"
                                                     prefix-icon="fas fa-user"
                                                     clearable size="mini" style="width: 100%"
                                                     @clear="scope.row.InstallDate = ''"
                                    ></el-autocomplete>
                                </el-col>
                                <el-col :span="12">
                                    <el-date-picker format="dd/MM/yyyy" placeholder="Install." size="mini"
                                                    style="width: 100%" type="date"
                                                    v-model="scope.row.InstallDate"
                                                    value-format="yyyy-MM-dd"
                                                    :picker-options="{firstDayOfWeek:1, disabledDate(time) { return time.getTime() > Date.now(); }}"
                                                    :disabled="!scope.row.InstallActor" :clearable="false"
                                    ></el-date-picker>
                                </el-col>
                            </el-row>
                        </template>
                    </el-table-column>
                    <el-table-column
                            label="Mesure"
                            width="300px" min-width="250px" :resizable="true"
                    >
                        <template slot-scope="scope">
                            <el-row type="flex" align="middle" :gutter="10">
                                <el-col :span="12">
                                    <el-autocomplete v-model="scope.row.MeasureActor"
                                                     :fetch-suggestions="UserSearch"
                                                     placeholder="Equipier"
                                                     prefix-icon="fas fa-user"
                                                     clearable size="mini" style="width: 100%"
                                                     :disabled="!scope.row.InstallDate"
                                                     @clear="scope.row.MeasureDate = ''"
                                    ></el-autocomplete>
                                </el-col>
                                <el-col :span="12">
                                    <el-date-picker format="dd/MM/yyyy" placeholder="Mesure" size="mini"
                                                    style="width: 100%" type="date"
                                                    v-model="scope.row.MeasureDate"
                                                    value-format="yyyy-MM-dd"
                                                    :picker-options="{firstDayOfWeek:1, disabledDate(time) { return time.getTime() > Date.now(); }}"
                                                    :disabled="!scope.row.MeasureActor" :clearable="false"
                                    ></el-date-picker>
                                </el-col>
                            </el-row>
                        </template>
                    </el-table-column>
                    <el-table-column
                            label="Bloquage"
                            width="140px" :resizable="true"
                    >
                        <template slot-scope="scope">
                            <div>
                                <el-checkbox 
                                        v-model="scope.row.Blockage"
                                        size="mini" :disabled="scope.row.NeedSignature && !scope.row.Signed"
                                >Bloquage</el-checkbox>
                                <div v-if="scope.row.NeedSignature">
                                    <el-checkbox 
                                            v-model="scope.row.Signed"
                                            size="mini" @change="CheckSignature(scope.row)"
                                    >Sign. convention</el-checkbox>
                                </div>
                            </div>
                        </template>
                    </el-table-column>
                    <el-table-column
                            label="Commentaire" prop="Comment"
                    >
                        <template slot-scope="scope">
                            <el-input clearable placeholder="Commentaire sur tronçon" size="mini" type="textarea" autosize
                                      v-model.trim="scope.row.Comment"
                            ></el-input>
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
