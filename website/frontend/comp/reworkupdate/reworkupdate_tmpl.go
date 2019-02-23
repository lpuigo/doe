package reworkupdate

const template string = `
        <el-container 
                style="height: 100%"
        >
            <el-header style="height: auto; padding: 5px">
                <el-row :gutter="10" type="flex" align="middle" style="margin-bottom: 10px">
                    <el-col :offset="2" :span="2">
                        <span style="float:right">Contrôle:</span>
                    </el-col>
                    <el-col :span="4">
                        <el-date-picker format="dd/MM/yyyy" placeholder="Contrôle" size="mini"
                                        style="width: 100%" type="date"
                                        v-model="worksite.Rework.ControlDate"
                                        value-format="yyyy-MM-dd"
                                        :picker-options="{firstDayOfWeek:1, disabledDate(time) { return time.getTime() > Date.now(); }}"
                                        :disabled="true"
                        ></el-date-picker>
                    </el-col>
                    <el-col :span="2">
                        <span style="float:right">Soumission:</span>
                    </el-col>
                    <el-col :span="4">
                        <el-date-picker format="dd/MM/yyyy" placeholder="Soumission" size="mini"
                                        style="width: 100%" type="date"
                                        v-model="worksite.Rework.SubmissionDate"
                                        value-format="yyyy-MM-dd"
                                        :picker-options="{firstDayOfWeek:1, disabledDate(time) { return time.getTime() > Date.now(); }}"
                                        :disabled="true"
                        ></el-date-picker>
                    </el-col>
                    <el-col :offset="8" :span="4">
                        <worksite-status-tag v-model="worksite"></worksite-status-tag>
                    </el-col>
                </el-row>
            </el-header>
            <el-main  style="height: 100%; padding: 0px">
                <!--:row-class-name="TableRowClassName"-->
                <el-table
                        :data="filteredReworks"
                        height="100%" :border=true size="mini"
                >
                   
                    <el-table-column
                            label="PT" prop="PT"
                            width="150px" :resizable="true" :show-overflow-tooltip=true
                    >
                        <template slot-scope="scope">
                            <div>{{GetTronconRef(scope.row.PT)}}</div>
                            <div>{{GetTronconAddress(scope.row.PT)}}</div>
                        </template>
                    </el-table-column>
                    
                    <el-table-column
                            label="Nb Reprise"
                            width="70px" :resizable="true" :show-overflow-tooltip=true
                    >
                        <template slot-scope="scope">
                            <div>{{scope.row.NbKO}}</div>
                        </template>
                    </el-table-column>
                    
                    <el-table-column
                            label="Installation"
                            width="200px" :resizable="true" :show-overflow-tooltip=true
                    >
                        <template slot-scope="scope">
                            <div>{{GetTronconInstallDate(scope.row.PT)}}</div>
                            <div>{{GetTronconInstallActor(scope.row.PT)}}</div>
                        </template>
                    </el-table-column>

                    <el-table-column
                            label="Description" prop="Description"
                            width="400px" :resizable=true
                    ></el-table-column>
                    
                    <el-table-column
                            label="Reprise"   
                    >
                        <template slot-scope="scope">
                            <el-row type="flex" align="middle" :gutter="10">
                                <el-col :span="12">
                                    <el-autocomplete v-model="scope.row.FixActor"
                                                     :fetch-suggestions="UserSearch"
                                                     placeholder="Equipier"
                                                     prefix-icon="fas fa-user"
                                                     clearable size="mini" style="width: 100%"
                                                     @clear="scope.row.FixDate = ''"
                                    ></el-autocomplete>
                                </el-col>
                                <el-col :span="12">
                                    <el-date-picker format="dd/MM/yyyy" placeholder="Reprise" size="mini"
                                                    style="width: 100%" type="date"
                                                    v-model="scope.row.FixDate"
                                                    value-format="yyyy-MM-dd"
                                                    :picker-options="{firstDayOfWeek:1, disabledDate(time) { return time.getTime() > Date.now(); }}"
                                                    :disabled="!scope.row.FixActor" :clearable="false"
                                    ></el-date-picker>
                                </el-col>
                            </el-row>
                        </template>  
                    </el-table-column>
                    
                </el-table>
            </el-main>
        </el-container>

`
