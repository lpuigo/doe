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
			<el-col :span="10">
				<h2 v-if="current_worksite" style="margin: 0 0">
					<i class="fas fa-wrench icon--left"></i>Mise à jour du chantier: <span style="color: #ccebff">{{current_worksite.City}} - {{current_worksite.Ref}}</span>
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
	<div v-loading="loading" style="min-height: 35vh;height: 65vh;overflow-x: hidden;overflow-y: auto;padding-right: 6px;">
        <el-container style="height: 100%">
            <el-header style="height: auto; padding: 5px">
                <el-row :span="20">
                    <el-col :offset="20" :span="4">
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
                <el-table
                        :data="filteredTroncons"
                        :row-class-name="TableRowClassName"
                        height="100%" :border=true size="mini"
                >
                    <el-table-column
                            label="Tronçon"
                            :resizable="true" :show-overflow-tooltip=true
                    >
                        <tempalte slot-scope="scope">
                            <span>{{scope.row | FormatTronconRef}}</span>
                        </tempalte>
                    </el-table-column>
                    <el-table-column
                            label="Status"
                            :resizable="true" :show-overflow-tooltip=true
                    >
                        <tempalte slot-scope="scope">
                            <span>{{scope.row | FormatStatus}}</span>
                        </tempalte>
                    </el-table-column>
                    <el-table-column
                            label="Installation"
                            :resizable="true"
                    >
                        <tempalte slot-scope="scope">
                            <el-date-picker format="dd/MM/yyyy" placeholder="Date Install." size="mini"
                                            style="width: 100%" type="date"
                                            v-model="scope.row.InstallDate"
                                            value-format="yyyy-MM-dd"
                                            :picker-options="{firstDayOfWeek:1, disabledDate(time) { return time.getTime() > Date.now(); }}"
                                            :clearable="false"
                            ></el-date-picker>
                        </tempalte>
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
