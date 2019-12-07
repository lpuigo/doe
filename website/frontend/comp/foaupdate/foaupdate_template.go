package foaupdate

const template string = `
<el-table
        :data="filteredJunctions"
        :row-class-name="TableRowClassName"
        height="100%" :border=true size="mini"
		@row-dblclick="SetSelectedState"
		@selection-change="HandleSelectionChange"
>
    <el-table-column
			type="selection"
			width="55" align="center">
    </el-table-column>

    <el-table-column
            label="Référence" prop="Ref" sortable
            width="150px" :resizable="true" :show-overflow-tooltip=true
    ></el-table-column>

    <el-table-column
            label="Insee" prop="Insee" sortable
            width="150px" align="center" :resizable="true" :show-overflow-tooltip=true
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
`
