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
                <div slot="reference">{{GetNodeDesc(scope.row)}}</div>
            </el-popover>
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
