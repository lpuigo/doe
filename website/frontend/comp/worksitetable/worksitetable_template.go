package worksitetable

const (
	template string = `<!--header-row-class-name="prjptf-light"-->
<!--:default-sort = "{prop: 'client', order: 'ascending'}"-->
<el-container style="height: 100%">
	<el-header style="height: auto; padding: 5px">
		<el-row type="flex" :gutter="10" align="center">
			<el-col :span="20">
				<el-button type="primary" plain
						   icon="fas map-marker-alt icon--left"
						   size="mini"
						   @click="AddWorksite"
				>Nouveau Chantier</el-button>
			</el-col>
			<el-col :span="4">
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
				:data="filteredWorksites"
				:row-class-name="TableRowClassName"
				:default-sort = "{prop: 'Status', order: 'descending'}"
				height="100%"
				:border=true
				@row-dblclick="SetSelectedWorksite"
		>
            <!--  :sort-method="SortStatus" :sort-by="['Status', 'Client', 'City', 'Ref']"  -->
			<el-table-column
					label="Statut" prop="Status" :formatter="FormatStatus"
                    sortable :sort-method="SortStatus"
					width="100px" :resizable=true :show-overflow-tooltip=true
					:filters="FilterList('Status')"	:filter-method="FilterHandler"	filter-placement="bottom-end" :filtered-value="FilteredStatusValue()"
			></el-table-column>
			<el-table-column
					label="Client" prop="Client" sortable :sort-by="['Client', 'City', 'Ref']"
					width="140px" :resizable=true :show-overflow-tooltip=true
					:filters="FilterList('Client')"	:filter-method="FilterHandler"	filter-placement="bottom-end"
			></el-table-column>
			<el-table-column
					label="Ville" prop="City" sortable :sort-by="['City', 'Ref']"
					width="140px" :resizable=true :show-overflow-tooltip=true
					:filters="FilterList('City')"	:filter-method="FilterHandler"	filter-placement="bottom-end"
			></el-table-column>
			<el-table-column
					label="Référence" sortable
					width="140px" :resizable=true :show-overflow-tooltip=true
			>        
				<template slot-scope="scope">
                    <span @click="SetSelectedWorksite(scope.row)" class="link">{{scope.row.Ref}}</span>
				</template>
			</el-table-column>
			<el-table-column
					label="Soumission" prop="OrderDate" sortable :sort-by="['OrderDate', 'Ref']"
					width="120px" :resizable=true :show-overflow-tooltip=true
					align="center"	:formatter="FormatDate"
			></el-table-column>
			<el-table-column
					label="Nb Cmd"
					width="150px" :resizable=true :show-overflow-tooltip=true
			>
				<template slot-scope="scope">
					<worksite-info :worksite="scope.row"></worksite-info>
				</template>
			</el-table-column>
			<el-table-column
					label="Commentaire" prop="Comment"
					min-width="120px" :resizable=true
			></el-table-column>
		</el-table>		
	</el-main>
</el-container>


`
)
