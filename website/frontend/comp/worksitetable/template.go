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
            <p>PMZ: {{ props.row.Pmz.Ref }}</p>
            <p>PA: {{ props.row.Pa.Ref }}</p>
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

<!--    
    
    <el-table-column 
            label="Client"	prop="client"	width="160px" sortable :sort-by="['client','name']" 
            :resizable=true :show-overflow-tooltip=true
    >
        <template slot-scope="scope">
            <i v-if="HasAuditInfo(scope.row)"  class="fas fa-info-circle link icon--left" @click="ShowTableProjectAudit(scope.row)"></i><span>{{scope.row.client}}</span>
        </template>
    </el-table-column>

    <el-table-column
            label="Project Name"	width="200px"
			:resizable=true :show-overflow-tooltip=true
    >
        <template slot-scope="scope">
            <i v-if="scope.row.hasStat" 
            		class="fas fa-chart-line icon--left link"
					@click="ShowTableProjectStat(scope.row)"
            ></i><span>{{scope.row.name}}</span>
        </template>
	</el-table-column>

    <el-table-column
            label="Comment" min-width="120px" sortable :sort-by="['risk', 'client','name']"
		    :resizable=false
    >
        <template slot-scope="scope">
            <i :class="RiskIconClass(scope.row.risk)"></i><span>{{scope.row.comment}}</span>
        </template>
	</el-table-column>

    <el-table-column 
            label="KickOff"	prop="milestones.Kickoff"	width="100px"	sortable    :sort-by="['milestones.Kickoff', 'client','name']"
		    :resizable=false    align="center"	:formatter="FormatDate"
    ></el-table-column>

    <el-table-column 
            label="UAT"	prop="milestones.UAT"	width="100px"	sortable    :sort-by="['milestones.UAT', 'client','name']"
		    :resizable=false    align="center"	:formatter="FormatDate"
    ></el-table-column>

    <el-table-column 
            label="RollOut"	prop="milestones.RollOut"	width="100px"	sortable    :sort-by="['milestones.RollOut', 'client','name']"
		    :resizable=false    align="center"	:formatter="FormatDate"
    ></el-table-column>

    <el-table-column 
            label="WorkLoad"	width="120px"
		    :resizable=false	align="center"
    >
        <template slot-scope="scope">
            <project-progress-bar :project="scope.row"></project-progress-bar>
        </template>
	</el-table-column>

    <el-table-column
            label="PS"	prop="lead_ps"	width="120px" sortable :sort-by="['lead_ps', 'client','name']"
		    :resizable=false :show-overflow-tooltip=true
 		    :filters="FilterList('lead_ps')"	:filter-method="FilterHandler"	filter-placement="bottom-end"
    ></el-table-column>

    <el-table-column
            label="Lead Dev"	prop="lead_dev"	width="120px" sortable :sort-by="['lead_dev', 'client','name']"
		    :resizable=false :show-overflow-tooltip=true
 		    :filters="FilterList('lead_dev')"	:filter-method="FilterHandler"	filter-placement="bottom-end"
    ></el-table-column>

    <el-table-column
            label="Type"	prop="type"	width="80px"
		    :resizable=false
 		    :filters="FilterList('type')"	:filter-method="FilterHandler"	filter-placement="bottom-end"
    ></el-table-column>

    <el-table-column
            label="Status"	prop="status"	width="120px" sortable :sort-by="['status', 'client','name']"
		    :resizable=false :show-overflow-tooltip=true
 		    :filters="FilterList('status')"	:filter-method="FilterHandler"	filter-placement="bottom-end" :filtered-value="FilteredValue()"
	></el-table-column>
-->
</el-table>
`
)
