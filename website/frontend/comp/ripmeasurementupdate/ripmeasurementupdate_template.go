package ripmeasurementupdate

const template string = `
<el-table
        :data="filteredMeasurements"
        :row-class-name="TableRowClassName"
        height="100%" :border=true size="mini"
>
    <el-table-column
            label="Noeud"
            width="100px" :resizable="true" :show-overflow-tooltip=true
    >
        <template slot-scope="scope">
            <el-popover placement="bottom-start" width="600"
                        title="Evenements de mesure:"
                        trigger="hover"
            >
                <el-row v-for="(nodename, index) in scope.row.NodeNames" :key="index" :gutter="5">
                    <el-col :span="7">
                        <div>{{index+1}} - {{nodename}}</div>
                    </el-col>
                    <el-col :span="3">
                        <span>{{GetNode(nodename).DistFromPm}} m</span>
                    </el-col>
                    <el-col :span="14">
                        <span>{{GetNode(nodename).Address}}</span>
                    </el-col>
                </el-row>
                <div slot="reference">{{scope.row.DestNodeName}}</div>
            </el-popover>
        </template>
    </el-table-column>
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
