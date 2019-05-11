package rippullingupdate

const template string = `
<el-table
        :data="filteredPullings"
        :row-class-name="TableRowClassName"
        height="100%" :border=true size="mini"
>
    <el-table-column
            label="Cable" prop="CableName"
            width="200px" :resizable="true" :show-overflow-tooltip=true
    ></el-table-column>
    <el-table-column
            label="Troncon"
            width="100px" :resizable="true" :show-overflow-tooltip=true
    >
        <template slot-scope="scope">
            <div>{{GetFirstPullingChunk(scope.row).TronconName}}</div>
        </template>
    </el-table-column>
    <el-table-column
            label="PT Départ / Arrivée"
            width="160px" :resizable="true" :show-overflow-tooltip=true
    >
        <template slot-scope="scope">
            <div>{{GetFirstPullingChunk(scope.row).StartingNodeName}}<i class="icon--right fas fa-arrow-right icon--left"></i>{{GetLastPullingChunk(scope.row).EndingNodeName}}</div>
        </template>
    </el-table-column>
<!--                <el-table-column-->
<!--                        label="Adresse Départ"-->
<!--                        width="240px" :resizable="true" :show-overflow-tooltip=true-->
<!--                >-->
<!--                    <template slot-scope="scope">-->
<!--                        <div>{{GetNode(GetFirstPullingChunk(scope.row).StartingNodeName).Address}}</div>-->
<!--                    </template>-->
<!--                </el-table-column>-->
<!--                <el-table-column-->
<!--                        label="PT Arrivée"-->
<!--                        width="100px" :resizable="true" :show-overflow-tooltip=true-->
<!--                >-->
<!--                    <template slot-scope="scope">-->
<!--                        <div>{{GetLastPullingChunk(scope.row).EndingNodeName}}</div>-->
<!--                    </template>-->
<!--                </el-table-column>-->
<!--                <el-table-column-->
<!--                        label="Adresse Arrivée"-->
<!--                        width="240px" :resizable="true" :show-overflow-tooltip=true-->
<!--                >-->
<!--                    <template slot-scope="scope">-->
<!--                        <div>{{GetNode(GetLastPullingChunk(scope.row).EndingNodeName).Address}}</div>-->
<!--                    </template>-->
<!--                </el-table-column>-->
    <el-table-column
            label="Distance"
            width="100px" :resizable="true" :show-overflow-tooltip=true
    >
        <template slot-scope="scope">
            <pulling-distances-info v-model="scope.row"></pulling-distances-info>
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
