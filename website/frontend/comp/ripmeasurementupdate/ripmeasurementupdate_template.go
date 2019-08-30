package ripmeasurementupdate

const template string = `<el-container style="height: 100%">
    <el-header style="height: auto; padding: 0px">
        <el-row type="flex" align="middle">
            <el-col :span="4">
                <el-popover placement="right" title="Chargement d'un fichier ZIP de mesures :"
                            trigger="click"
                            width="400"
                            v-model="uploadVisible"
                >
                    <div>
                        <el-select v-model="measurementTeam" clearable filterable
                                   size="mini"
                                   placeholder="Equipe Mesures"
                                   style="width: 360px; margin-bottom: 10px"
                                   @clear="measurementTeam=''"
                        >
                            <el-option
                                    v-for="item in GetTeams()"
                                    :key="item.value"
                                    :label="item.label"
                                    :value="item.value">
                            </el-option>
                        </el-select>
                    </div>
                    <el-upload v-if="measurementTeam.length > 5" 
                               action="/api/ripsites/measurement"
                               drag
                               style="width: 360px"
                               :before-upload="BeforeUpload"
                               :on-success="UploadSuccess"
                               :on-error="UploadError"
                    >
                        <i class="el-icon-upload"></i>
                        <div class="el-upload__text">Drop Zip file here or <em>click to upload</em></div>
                        <!--        <div class="el-upload__tip" slot="tip">.zip files with a size less than 30Mb</div>-->
                    </el-upload>
                    <el-button slot="reference" type="primary" size="mini">Mesures...</el-button>
                </el-popover>
            </el-col>
            <el-col :offset="3" :span="3">
                <span>Nb Fibres total: {{NbTotFiber}}</span>            
            </el-col>
            <el-col :span="3">
                <span>#OK: {{measurementSummary.NbOK}} ({{GetPct(measurementSummary.NbOK)}})</span>            
            </el-col>
            <el-col :span="3">
                <span>#Seuil 1: {{measurementSummary.NbWarn1}} ({{GetPct(measurementSummary.NbWarn1)}})</span>            
            </el-col>
            <el-col :span="3">
                <span>#Seuil 2: {{measurementSummary.NbWarn2}} ({{GetPct(measurementSummary.NbWarn2)}})</span>            
            </el-col>
            <el-col :span="3">
                <span>#KO: {{measurementSummary.NbKO}} ({{GetPct(measurementSummary.NbKO)}})</span>            
            </el-col>
        </el-row>
    </el-header>
    <div style="height: 100%;overflow-x: hidden;overflow-y: auto;padding: 0px 0px; margin-top: 8px">
		<!--
        <el-row v-for="(meas, index) in filteredMeasurements" :key="meas" 
                style="margin-top: 2px; padding: 2px 5px; border-radius: 4px" 
                type="flex" align="middle"
                :class="MeasurementClassName(meas)"
        >
            <el-col :span="2">
                <el-popover placement="bottom-start" title="Evenements de mesure:"
                            trigger="hover"
                            width="600"
                >
                    <el-row :gutter="5" :key="index" v-for="(nodename, index) in meas.NodeNames">
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
                    <div slot="reference">{{meas.DestNodeName}}</div>
                </el-popover>
            </el-col>
            <el-col :span="1">
                <span>{{meas.NbFiber}} fibres</span>
            </el-col>
            <el-col :span="1">
                <span>{{meas.NodeNames.length}} épi.</span>
            </el-col>
            <el-col :span="1">
                <span>{{GetDestNodeDist(meas)}}m</span>
            </el-col>
            <el-col :span="19">
                <rip-state-update :client="value.Client" :user="user" v-model="meas.State"></rip-state-update>
            </el-col>
        </el-row>
		-->
		<el-table
				:data="filteredMeasurements"
				:row-class-name="TableRowClassName"
				height="100%" :border=true size="mini"
		>
			<el-table-column
					label="Pbo"
					width="200px" :resizable="true" :show-overflow-tooltip=true
			>
				<template slot-scope="scope">
					<el-popover placement="bottom-start" title="Evenements de mesure:"
								trigger="hover"
								width="600"
								open-delay="400"
					>
						<el-row :gutter="5" v-for="(nodename, index) in scope.row.NodeNames" :key="index" >
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
						<span slot="reference">{{scope.row.DestNodeName}}</span>
					</el-popover>
				</template>
			</el-table-column>
			<el-table-column
					label="Info"
					width="200px" :resizable="true" :show-overflow-tooltip=true
			>
				<template slot-scope="scope">
					<el-row :gutter="5">
						<el-col :span="8">
							<span>{{scope.row.NbFiber}} fibres</span>
						</el-col>
						<el-col :span="8">
							<span>{{scope.row.NodeNames.length}} épi.</span>
						</el-col>
						<el-col :span="8">
							<span>{{GetDestNodeDist(scope.row)}}m</span>
						</el-col>
					</el-row>
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

    </div>
</el-container>
`
