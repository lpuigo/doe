package worksitetable

const (
	template string = `<!--header-row-class-name="prjptf-light"-->
<!--:default-sort = "{prop: 'client', order: 'ascending'}"-->
<el-table
    :data="filteredWorksites"
    :row-class-name="TableRowClassName"
    :default-sort = "{prop: 'Status', order: 'descending'}"
    @current-change="SetSelectedWorksite"
    @row-dblclick="SelectRow"
	height="100%"
	:border=true
>
    <el-table-column type="expand">
        <template slot-scope="props">
            <worksite-detail
                    :worksite="props.row"
                    :readonly="false"
            ></worksite-detail>
        </template>
    </el-table-column>
    <el-table-column
            label="Statut" prop="Status" sortable :sort-by="['Status', 'Ref']"
            width="100px" :resizable=true :show-overflow-tooltip=true
            :filters="FilterList('Status')"	:filter-method="FilterHandler"	filter-placement="bottom-end" :filtered-value="FilteredStatusValue()"
    ></el-table-column>
    <el-table-column
            label="Ville" prop="City" sortable :sort-by="['City', 'Ref']"
            width="140px" :resizable=true :show-overflow-tooltip=true
            :filters="FilterList('City')"	:filter-method="FilterHandler"	filter-placement="bottom-end"
    ></el-table-column>
    <el-table-column
            label="Référence" prop="Ref" sortable
            width="140px" :resizable=true :show-overflow-tooltip=true
    ></el-table-column>
    <el-table-column
            label="Soumission" prop="OrderDate" sortable :sort-by="['OrderDate', 'Ref']"
            width="120px" :resizable=true :show-overflow-tooltip=true
            align="center"	:formatter="FormatDate"
    ></el-table-column>
    <el-table-column
            label="Nb Cmd" sortable
            width="120px" :resizable=true :show-overflow-tooltip=true
    >
        <template slot-scope="scope">
            <worksite-info
                    :worksite="scope.row"
            ></worksite-info>
        </template>
    </el-table-column>
    <el-table-column
            label="Commentaire" prop="Comment"
            min-width="120px" :resizable=true :show-overflow-tooltip=true
    ></el-table-column>
</el-table>
`
)
