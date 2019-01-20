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
	<div v-loading="loading" style="height: 65vh;overflow-x: hidden;overflow-y: auto;padding-right: 6px;">
        <el-container 
                v-if="HasRework" 
                style="height: 100%"
        >
            <el-header style="height: auto; padding: 5px">
                <el-row :gutter="10" type="flex" align="middle">
                    <el-col :span="3">
                        <el-button
                                type="success" plain icon="fas fa-tools icon--left" size="mini"
                                @click="AddDefect"
                        >Ajouter Reprise</el-button>
                    </el-col>
                    <el-col :offset="2" :span="2">
                        <span style="float:right">Contrôle:</span>
                    </el-col>
                    <el-col :span="4">
                        <el-date-picker format="dd/MM/yyyy" placeholder="Contrôle" size="mini"
                                        style="width: 100%" type="date"
                                        v-model="current_worksite.Rework.ControlDate"
                                        value-format="yyyy-MM-dd"
                                        :picker-options="{firstDayOfWeek:1, disabledDate(time) { return time.getTime() > Date.now(); }}"
                                        :clearable="false"
                        ></el-date-picker>
                    </el-col>
                    <el-col :span="1">
                        <el-button
                                icon="fas fa-angle-double-right" size="mini" circle
                                @click="current_worksite.Rework.SubmissionDate=current_worksite.Rework.ControlDate"
                        ></el-button>
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
                                        :clearable="false"
                        ></el-date-picker>
                    </el-col>
                </el-row>
            </el-header>
            <el-main  style="height: 100%; padding: 0px">
                <!--:span-method="OrderSpanMethod"-->
                <el-table
                        :data="filteredReworks"
                        :row-class-name="TableRowClassName"
                        height="100%" :border=true size="mini"
                >
                    <el-table-column
                            label="Action"
                            width="60px"
                    >
                        <template slot-scope="scope">
                            <el-button
                                    type="danger" icon="el-icon-delete" circle size="mini"
                                    @click="RemoveDefect(scope.$index)"
                            ></el-button>
                        </template>
                    </el-table-column>
                    <el-table-column
                            label="PT" prop="PT"
                            width="250px" :resizable="true" :show-overflow-tooltip=true
                    >
                        <template slot-scope="scope">
                            <el-select v-model="scope.row.PT" filterable placeholder="Choix PT" size="mini" style="width: 100%">
                                <el-option
                                        v-for="item in GetPTs()"
                                        :key="item.value"
                                        :label="item.label"
                                        :value="item.value">
                                    <!--<span style="float: left">{{ item.value }}</span>-->
                                    <!--<span style="float: right; color: #8492a6; font-size: 90%">{{ item.label }}</span>-->
                                </el-option>
                            </el-select>                            
                        </template>
                    </el-table-column>
                    <el-table-column
                            label="Soumission"
                            width="150px" :resizable="true" :show-overflow-tooltip=true
                    >
                        <template slot-scope="scope">
                            <el-date-picker format="dd/MM/yyyy" placeholder="Soumission" size="mini"
                                            style="width: 100%" type="date"
                                            v-model="scope.row.SubmissionDate"
                                            value-format="yyyy-MM-dd"
                                            :picker-options="{firstDayOfWeek:1, disabledDate(time) { return time.getTime() > Date.now(); }}"
                                            :clearable="false"
                            ></el-date-picker>
                        </template>
                    </el-table-column>
                    <el-table-column
                            label="Description"                            
                    >
                        <!--<template slot="header" slot-scope="scope">-->
                            <!--<el-input-->
                                    <!--v-model="search"-->
                                    <!--size="mini"-->
                                    <!--placeholder="Type to search"/>-->
                        <!--</template>-->
                        <template slot-scope="scope">
                            <el-input clearable placeholder="Description de la reprise" size="mini" type="textarea" autosize
                                      v-model.trim="scope.row.Description"
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
