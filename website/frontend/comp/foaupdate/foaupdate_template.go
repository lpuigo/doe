package foaupdate

const template string = `<el-container style="height: 100%">
    <el-header style="height: auto; padding: 0px">
        <el-row type="flex" align="middle">
            <el-col :span="2">
				<span style="margin-left: 5px">Sélection: {{SelectedFoas.Foas.length}} / {{filteredJunctions.length}}</span>
            </el-col>
            <el-col :span="10">
				<el-button-group>
					<el-button type="primary" icon="far fa-edit icon--left" size="mini"
							:disabled="SelectedFoas.Foas.length == 0" 
							@click="EditSelectedFoas()"
					>Editer: {{SelectedFoas.Foas.length}}</el-button>
					<el-button type="primary" icon="fas fa-plus icon--left" size="mini" @click="">Ajouter</el-button>
				</el-button-group>
            </el-col>
        </el-row>
    </el-header>
    <div style="height: 100%;overflow-x: hidden;overflow-y: auto;padding: 0px 0px; margin-top: 8px">
		<el-table
				ref="foaTable"
				:data="filteredJunctions"
				:row-class-name="TableRowClassName"
				height="100%" :border=true size="mini"
				@row-dblclick="EditFoa"
				@selection-change="HandleSelectionChange"
		>
			<el-table-column
					type="selection"
					width="55" align="center">
			</el-table-column>
		
			<el-table-column
					label="Insee" prop="Insee" sortable
					width="150px" align="center" :resizable="true" :show-overflow-tooltip=true
			></el-table-column>
			
			<el-table-column
					label="Référence" prop="Ref" sortable
					width="150px" :resizable="true" :show-overflow-tooltip=true
			></el-table-column>
		
			<el-table-column
					label="Type" prop="Type" sortable
					width="150px" align="center" :resizable="true" :show-overflow-tooltip=true
			></el-table-column>
			
			<el-table-column
					label="Etat" prop="State.Status" sortable
					width="120px" align="center" :resizable="true" :formatter="FormatStatus"
			></el-table-column>
			<!--
				:filters="FilterList('State.Status')"	:filter-method="FilterHandler"	filter-placement="bottom-end"
			-->
		
			<el-table-column
					label="Acteurs"
					width="250px" :resizable="true"
			>
				<template slot-scope="scope">
					<div>{{GetActors(scope.row)}}</div>
				</template>
			</el-table-column>
		
			<el-table-column
					label="Date" prop="State.Date" sortable
					width="120px" align="center" :resizable="true" :formatter="FormatDate"
			></el-table-column>
			
			<el-table-column
					label="Commentaire" prop="State.Comment"
			></el-table-column>
			
			<!--
			<el-table-column
					label="Etat"
			>
				<template slot-scope="scope">
					<rip-state-update v-model="scope.row.State" :user="user" :client="value.Client"></rip-state-update>
				</template>
			</el-table-column>
			-->
		</el-table>
	</div>
</el-container>
`
