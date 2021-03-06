package ripjunctionupdate

const template_flat string = `
<el-table
        :data="filteredJunctions"
        :row-class-name="TableRowClassName"
        height="100%" :border=true size="mini"
		@row-dblclick="SetSelectedState"
>
    <el-table-column
            label="Noeud" prop="NodeName" sortable
            width="150px" :resizable="true" :show-overflow-tooltip=true
    >
        <template slot-scope="scope">
            <el-popover placement="bottom-start" width="400"
                        title="Opérations à réaliser:"
                        trigger="hover"
            >
                <el-row v-for="(ope, index) in scope.row.Operations" :key="index" :gutter="5">
                    <el-col :span="14">
                        <div v-if="ope.TronconName">{{index+1}} - {{ope.Type}}<i class="icon--right fas fa-arrow-right icon--left"></i>{{ope.TronconName}}</div>
                        <div v-else>{{index+1}} - {{ope.Type}}</div>
                    </el-col>
                    <el-col :span="10">
                        <span>{{ope.NbFiber}} fibre(s)</span>
                    </el-col>
                </el-row>
                <div slot="reference" class="header-menu-container">
                    <span>{{scope.row.NodeName}}</span>
                    <span>{{GetNode(scope.row).Ref}}</span>
                </div>
            </el-popover>
        </template>
    </el-table-column>

    <el-table-column
            label="Adresse" sortable :sort-by="GetNodeAttr('Address')"
            width="250px" :resizable="true" :show-overflow-tooltip=true
    >
        <template slot-scope="scope"> 
            <div>{{GetNode(scope.row).Address}}</div>
        </template>	
    </el-table-column>
    
    <el-table-column
            label="Type" sortable :sort-by="GetNodeType"
            width="150px" :resizable="true" :show-overflow-tooltip=true
    >
        <template slot-scope="scope">
            <div>{{GetNodeType(scope.row)}}</div>
        </template>
    </el-table-column>
    
    <el-table-column
            label="Tronçon" sortable :sort-by="GetNodeAttr('TronconInName')"
            width="120px" :resizable="true" :show-overflow-tooltip=true
    >
        <template slot-scope="scope">
            <div>{{GetTronconDesc(scope.row)}}</div>
        </template>
    </el-table-column>
    
    <el-table-column
            label="Etat" prop="State.Status" sortable
            width="120px" :resizable="true" :formatter="FormatStatus"
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
            label="Début" prop="State.DateStart" sortable
            width="120px" :resizable="true" :formatter="FormatDate"
    ></el-table-column>
    
    <el-table-column
            label="Fin" prop="State.DateEnd" sortable
            width="120px" :resizable="true" :formatter="FormatDate"
    ></el-table-column>
    
    <el-table-column
            label="Commentaire" prop="State.Comment" sortable
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

const template string = `
<el-table
        :data="filteredJunctions"
        :row-class-name="TableRowClassName"
        height="100%" :border=true size="mini"
>
    <el-table-column
            label="Noeud"
            width="300px" :resizable="true" :show-overflow-tooltip=true
    >
        <template slot-scope="scope">
			<div class="header-menu-container">
				<span>{{GetNodeDesc(scope.row)}}</span>
				<span :class="GetNodeTypeClass(scope.row)">{{GetNodeType(scope.row)}}
					<el-popover placement="bottom" width="400"
								title="Opérations à réaliser:"
								trigger="hover"
					>
						<div style="margin-bottom: 8px">{{GetNode(scope.row).Address}}</div>
						<el-row v-for="(ope, index) in scope.row.Operations" :key="index" :gutter="5">
							<el-col :span="14">
								<div v-if="ope.TronconName">{{index+1}} - {{ope.Type}}<i class="icon--right fas fa-arrow-right icon--left"></i>{{ope.TronconName}}</div>
								<div v-else>{{index+1}} - {{ope.Type}}</div>
							</el-col>
							<el-col :span="10">
								<span>{{ope.NbFiber}} fibre(s)</span>
							</el-col>
						</el-row>
						<i slot="reference" class="fas fa-info-circle icon--right"></i>
					</el-popover>
				</span>
			</div>
        </template>
    </el-table-column>

    <el-table-column
            label="Tronçon"
            width="140px" :resizable="true" :show-overflow-tooltip=true
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

const template_tree string = `
<el-table
        :data="filteredJunctionsTree"
        :row-class-name="TableRowClassName"
        height="100%" :border=true size="mini"
		row-key="NodeName" default-expand-all
>
    <el-table-column
            label="Noeud"
            width="500px" :resizable="true" :show-overflow-tooltip=true
    >
        <template slot-scope="scope">
            <el-popover placement="bottom-start" width="400"
                        title="Opérations à réaliser:"
                        trigger="hover"
            >
                <el-row v-for="(ope, index) in scope.row.Operations" :key="index" :gutter="5">
                    <el-col :span="14">
                        <div v-if="ope.TronconName">{{index+1}} - {{ope.Type}}<i class="icon--right fas fa-arrow-right icon--left"></i>{{ope.TronconName}}</div>
                        <div v-else>{{index+1}} - {{ope.Type}}</div>
                    </el-col>
                    <el-col :span="10">
                        <span>{{ope.NbFiber}} fibre(s)</span>
                    </el-col>
                </el-row>
                <span slot="reference">{{GetNodeDesc(scope.row)}}</span>
            </el-popover>
<!--			<span>{{GetNodeDesc(scope.row)}}</span>-->
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
