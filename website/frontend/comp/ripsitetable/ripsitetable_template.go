package ripsitetable

const template string = `<!--header-row-class-name="prjptf-light"-->
<!--:default-sort = "{prop: 'client', order: 'ascending'}"-->
<el-container style="height: 100%">
	<el-header style="height: auto; padding: 5px">
		<el-row type="flex" :gutter="10" align="middle">
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
				:data="filteredRipsites"
				:row-class-name="TableRowClassName"
				:default-sort = "{prop: 'OrderDate', order: 'descending'}"
				height="100%"
				:border=true size="mini"
				@row-dblclick="SetSelectedRipsite"
		>
            <!--  :sort-method="SortStatus" :sort-by="['Status', 'Client', 'City', 'Ref']"  -->
			<el-table-column
					label="Statut" prop="Status" :formatter="FormatStatus"
                    sortable :sort-method="SortStatus"
					width="100px" :resizable=true :show-overflow-tooltip=true
					:filters="FilterList('Status')"	:filter-method="FilterHandler"	filter-placement="bottom-end" :filtered-value="FilteredStatusValue()"
			></el-table-column>

			<el-table-column
					label="Client" prop="Client" sortable :sort-by="['Client', 'Ref']"
					width="120px" :resizable=true :show-overflow-tooltip=true
					:filters="FilterList('Client')"	:filter-method="FilterHandler"	filter-placement="bottom-end"
			></el-table-column>

			<el-table-column
					label="Référence" sortable
					width="120px" :resizable=true :show-overflow-tooltip=true
			>        
				<template slot-scope="scope">
                    <div class="header-menu-container">
                        <span @click="SetSelectedRipsite(scope.row)" class="link">{{scope.row.Ref}}</span>
						<a v-if="user.Permissions.Invoice" :href="AttachmentUrl(scope.row.Id)"><i class="link fas fa-file-excel"></i></a>
                    </div>
				</template>
			</el-table-column>

            <el-table-column
                    label="CAff" prop="Manager" sortable :sort-by="['Manager', 'Ref']"
                    width="120px" :resizable=true :show-overflow-tooltip=true
                    :filters="FilterList('Manager')"	:filter-method="FilterHandler"	filter-placement="bottom-end"
            ></el-table-column>

			<!--
			<el-table-column
					label="Info"
					width="240px" :resizable=true :show-overflow-tooltip=true
			>
				<template slot-scope="scope">
					<ripsiteinfo-info v-model="scope.row"></ripsiteinfo-info>
				</template>
			</el-table-column>
			-->

            <el-table-column
                    label="Soumission" prop="OrderDate" sortable :sort-by="['OrderDate', 'Ref']"
                    width="110px" :resizable=true :show-overflow-tooltip=true
                    align="center"	:formatter="FormatDate"
            ></el-table-column>

            <el-table-column
                    label="Nb. Points" width="130px" :resizable=true align="center"
            >
                <template slot-scope="scope">
                    <ripsiteinfo-progress-bar :total="scope.row.NbPoints" :blocked="scope.row.NbPointsBlocked" :done="scope.row.NbPointsDone"></ripsiteinfo-progress-bar>
                </template>
            </el-table-column>

            <el-table-column
                    label="Tirage" width="130px" :resizable=true align="center"
            >
                <template slot-scope="scope">
                    <ripsiteinfo-progress-bar :total="scope.row.NbPulling" :blocked="scope.row.NbPullingBlocked" :done="scope.row.NbPullingDone"></ripsiteinfo-progress-bar>
                </template>
            </el-table-column>

            <el-table-column
                    label="Racordement" width="130px" :resizable=true align="center"
            >
                <template slot-scope="scope">
                    <ripsiteinfo-progress-bar :total="scope.row.NbJunction" :blocked="scope.row.NbJunctionBlocked" :done="scope.row.NbJunctionDone"></ripsiteinfo-progress-bar>
                </template>
            </el-table-column>

            <el-table-column
                    label="Maj" prop="UpdateDate" sortable
                    width="90px" :resizable=true :show-overflow-tooltip=true
                    align="center"	:formatter="FormatDate"
            ></el-table-column>

			<el-table-column
					label="Commentaire" prop="Comment"
					min-width="120px" :resizable=true
			></el-table-column>
		</el-table>		
	</el-main>
</el-container>
`
