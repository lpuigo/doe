package ripmeasurementupdate

const template string = `
<el-table
        :data="filteredMeasurements"
        :row-class-name="TableRowClassName"
        height="100%" :border=true size="mini"
>
    <el-table-column
            label="Noeud" prop="DestNodeName"
            width="100px" :resizable="true" :show-overflow-tooltip=true
    ></el-table-column>
    <el-table-column
            label="Nb Fibre" prop="NbFiber"
            width="70px" :resizable="true" :show-overflow-tooltip=true
    ></el-table-column>
    <el-table-column
            label="Nb Epissure" 
            width="70px" :resizable="true" :show-overflow-tooltip=true
    >
        <template slot-scope="scope">
            <div>{{scope.row.NodeNames.length}}</div>
        </template>
    </el-table-column>
    <el-table-column
            label="Distance" 
            width="70px" :resizable="true" :show-overflow-tooltip=true
    >
        <template slot-scope="scope">
            <div>{{GetDestNodeDist(scope.row)}}m</div>
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
