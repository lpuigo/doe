package invoicetable

const (
	template string = `<!--header-row-class-name="prjptf-light"-->
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
				:data="filteredWorksites"
				:row-class-name="TableRowClassName"
				:default-sort = "{prop: 'Status', order: 'descending'}"
				height="100%"
				:border=true size="mini"
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
					label="Référence" prop="Ref" sortable :sort-by="['City', 'Ref']"
					width="110px" :resizable=true :show-overflow-tooltip=true
			></el-table-column>
			<el-table-column
					label="Nb Cmd"
					width="150px" :resizable=true :show-overflow-tooltip=true
			>
				<template slot-scope="scope">
					<worksiteinfo-info v-model="scope.row"></worksiteinfo-info>
				</template>
			</el-table-column>
			<el-table-column
					label="Montant" prop="InvoiceAmount" sortable :sort-by="['InvoiceAmount', 'Ref']"
					width="110px" :resizable=true :show-overflow-tooltip=true align="right"
			></el-table-column>
            <el-table-column
                    label="DOE" prop="DoeDate" sortable :sort-by="['DoeDate', 'Ref']"
                    width="110px" :resizable=true :show-overflow-tooltip=true
                    align="center"	:formatter="FormatDate"
            ></el-table-column>
            <el-table-column
                    label="Attachement" prop="AttachmentDate" sortable :sort-by="['AttachmentDate', 'Ref']"
                    width="110px" :resizable=true :show-overflow-tooltip=true
                    align="center"	:formatter="FormatDate"
            ></el-table-column>
            <el-table-column
                    label="Facture" prop="InvoiceName" sortable :sort-by="['InvoiceName', 'Ref']"
                    width="110px" :resizable=true :show-overflow-tooltip=true
            ></el-table-column>
			<el-table-column
					label="Date Fact." prop="InvoiceDate" sortable :sort-by="['InvoiceDate', 'Ref']"
					width="110px" :resizable=true :show-overflow-tooltip=true
					align="center"	:formatter="FormatDate"
			></el-table-column>
			<el-table-column
					label="Date Pay." prop="PaymentDate" sortable :sort-by="['PaymentDate', 'Ref']"
					width="110px" :resizable=true :show-overflow-tooltip=true
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
)
