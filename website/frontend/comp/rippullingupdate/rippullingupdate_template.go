package rippullingupdate

const template string = `<el-table
        :border=true
        :data="filteredPullings"
        :row-class-name="TableRowClassName" height="100%" size="mini"
>
    <el-table-column
            :resizable="true" :show-overflow-tooltip=true
            label="Cable" width="200px"
    >
        <template slot-scope="scope">
            <el-popover placement="bottom-start" width="700"
                        title="Tronçons traversés"
                        trigger="hover"
            >
                <el-row v-for="(chunk, index) in scope.row.Chuncks" :key="index" :gutter="5">
                    <el-col :span="5">
                        <span>{{index+1}} - {{chunk.TronconName}}</span>                        
                    </el-col>
                    <el-col :span="8">
                        <div>{{chunk.StartingNodeName}}<i
                                class="icon--right fas fa-arrow-right icon--left"></i>{{chunk.EndingNodeName}}
                        </div>                        
                    </el-col>
                    <el-col :span="11">
                        <span>Lov.: {{chunk.LoveDist}}m, Sout.: {{chunk.UndergroundDist}}m, Aér.: {{chunk.AerialDist}}m, Faç.: {{chunk.BuildingDist}}m, </span>
                    </el-col>
                </el-row>
                <span slot="reference">{{scope.row.CableName}}</span>
            </el-popover>
        </template>
    </el-table-column>
    <el-table-column
            :resizable="true" :show-overflow-tooltip=true 
            label="Troncon" width="100px"
    >
        <template slot-scope="scope">
            <div>{{GetFirstPullingChunk(scope.row).TronconName}}</div>
        </template>
    </el-table-column>
    <el-table-column
            :resizable="true" :show-overflow-tooltip=true 
            label="PT Départ / Arrivée" width="200px"
    >
        <template slot-scope="scope">
            <div>{{GetFirstPullingChunk(scope.row).StartingNodeName}}<i
                    class="icon--right fas fa-arrow-right icon--left"></i>{{GetLastPullingChunk(scope.row).EndingNodeName}}
            </div>
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
            :resizable="true"
            :show-overflow-tooltip=true label="Distance" width="100px"
    >
        <template slot-scope="scope">
            <pulling-distances-info v-model="scope.row"></pulling-distances-info>
        </template>
    </el-table-column>
    <el-table-column
            label="Etat"
    >
        <template slot-scope="scope">
            <rip-state-update :client="value.Client" :user="user" v-model="scope.row.State"></rip-state-update>
        </template>
    </el-table-column>
</el-table>
`
