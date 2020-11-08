package rippullingupdate

const template_flat string = `<el-table
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
const template string = `<el-table
        :border=true
        :data="filteredPullings"
        :row-class-name="TableRowClassName" height="100%" size="mini"
>
    <el-table-column
            :resizable="true" :show-overflow-tooltip=true
            label="Cable" width="350px"
    >
        <template slot-scope="scope">
			<div class="header-menu-container">
				<span>{{scope.row.CableName}}</span>
				<span :class="GetPullingTypeClass(scope.row)">{{GetFirstPullingChunk(scope.row).TronconName}}
					<el-popover placement="bottom-start" width="700"
								title="Tronçons traversés"
								trigger="hover"
					>
						<div>Boitier de départ: <span class="pt-name">{{GetFirstPullingChunk(scope.row).StartingNodeName}}</span></div>
						<el-row v-for="(chunk, index) in scope.row.Chuncks" :key="index" :gutter="5">
							<el-col :span="6">
								<span>{{index+1}} - {{chunk.TronconName}}</span>                        
							</el-col>
							<el-col :span="8">
								<div>
									<i class="icon--right fas fa-arrow-right icon--left"></i>
									<span class="pt-name">{{chunk.EndingNodeName}}</span>
								</div>                        
							</el-col>
							<el-col :span="10">
								<span v-if="chunk.LoveDist > 0">Lov.: {{chunk.LoveDist}}m </span>
								<span v-if="chunk.UndergroundDist > 0">Sout.: {{chunk.UndergroundDist}}m </span>
								<span v-if="chunk.AerialDist > 0" class="pulling-aerial">Aér.: {{chunk.AerialDist}}m </span>
								<span v-if="chunk.BuildingDist > 0" class="pulling-aerial">Faç.: {{chunk.BuildingDist}}m</span>
							</el-col>
						</el-row>
						<i slot="reference" class="fas fa-info-circle icon--right"></i>
					</el-popover>
				</span>
			</div>
        </template>
    </el-table-column>

    <el-table-column
            :resizable="true"
            :show-overflow-tooltip=true label="Distance" width="250px"
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
