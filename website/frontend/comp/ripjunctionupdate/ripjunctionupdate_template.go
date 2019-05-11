package ripjunctionupdate

const template string = `
<el-table
        :data="filteredJunctions"
        :row-class-name="TableRowClassName"
        height="100%" :border=true size="mini"
>
    <el-table-column
            label="Noeud"
            width="500px" :resizable="true" :show-overflow-tooltip=true
    >
        <template slot-scope="scope">
            <div>{{GetNodeDesc(scope.row)}}</div>
        </template>
    </el-table-column>
    <el-table-column
            label="Tronçon"
            width="120px" :resizable="true" :show-overflow-tooltip=true
    >
        <template slot-scope="scope">
            <div>{{GetTronconDesc(scope.row)}}</div>
        </template>
    </el-table-column>
    <el-table-column
            label="Etat"
    >
        <template slot-scope="scope">
            <rip-state-update v-model="scope.row.State" :user="user" :client="value.Client"></rip-state-update>
        </template>
    </el-table-column>
</el-table>
`
