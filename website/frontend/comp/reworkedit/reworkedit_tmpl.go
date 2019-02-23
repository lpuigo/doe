package reworkedit

const template string = `
<el-container 
        style="height: 100%"
>
    <el-header style="height: auto; padding: 5px">
        <el-row :gutter="10" type="flex" align="middle" style="margin-bottom: 10px">
            <el-col :span="3">
                <el-button
                        type="primary" plain icon="fas fa-tools icon--left" size="mini"
                        @click="AddDefect"
                >Ajouter Reprise</el-button>
            </el-col>
            <el-col :offset="2" :span="2">
                <span style="float:right">Contrôle:</span>
            </el-col>
            <el-col :span="4">
                <el-date-picker format="dd/MM/yyyy" placeholder="Contrôle" size="mini"
                                style="width: 100%" type="date"
                                v-model="worksite.Rework.ControlDate"
                                value-format="yyyy-MM-dd"
                                :picker-options="{firstDayOfWeek:1, disabledDate(time) { return time.getTime() > Date.now(); }}"
                                :clearable="false"
                ></el-date-picker>
            </el-col>
            <el-col :span="1">
                <el-button
                        icon="fas fa-angle-double-right" size="mini" circle
                        @click="worksite.Rework.SubmissionDate=worksite.Rework.ControlDate"
                ></el-button>
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
                                :clearable="false"
                ></el-date-picker>
            </el-col>
            <el-col :offset="4" :span="4">
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
                    label="Action"
                    width="60px"
            >
                <template slot-scope="scope">
                    <el-button
                            type="danger" icon="el-icon-delete" circle size="mini"
                            @click="RemoveDefect(scope.$index)"
                    ></el-button>
                </template>
            </el-table-column>
            
            <el-table-column
                    label="PT" prop="PT"
                    width="300px" :resizable="true" :show-overflow-tooltip=true
            >
                <template slot-scope="scope">
                    <el-select v-model="scope.row.PT" filterable placeholder="Choix PT" size="mini" style="width: 100%">
                        <el-option
                                v-for="item in GetPTs()"
                                :key="item.value"
                                :label="item.label"
                                :value="item.value">
                            <!--<span style="float: left">{{ item.value }}</span>-->
                            <!--<span style="float: right; color: #8492a6; font-size: 90%">{{ item.label }}</span>-->
                        </el-option>
                    </el-select>                            
                </template>
            </el-table-column>
            
            <el-table-column
                    label="Soumission"
                    width="150px" :resizable="true" :show-overflow-tooltip=true
            >
                <template slot-scope="scope">
                    <el-date-picker format="dd/MM/yyyy" placeholder="Soumission" size="mini"
                                    style="width: 100%" type="date"
                                    v-model="scope.row.SubmissionDate"
                                    value-format="yyyy-MM-dd"
                                    :picker-options="{firstDayOfWeek:1, disabledDate(time) { return time.getTime() > Date.now(); }}"
                                    :clearable="false"
                    ></el-date-picker>
                </template>
            </el-table-column>
            
            <el-table-column
                    label="Contrôles OK / KO"   
                    width="190px"
            >
                <template slot-scope="scope">
                    <el-row type="flex" align="middle" :gutter="10">
                        <el-col :span="12">
                            <el-input-number 
                                    v-model="scope.row.NbOK" 
                                    controls-position="right" 
                                    :min="0" size="mini" style="width: 80px"
                            ></el-input-number>
                        </el-col>
                        <el-col :span="12">
                            <el-input-number 
                                    v-model="scope.row.NbKO" 
                                    controls-position="right" 
                                    :min="0" size="mini" style="width: 80px"
                            ></el-input-number>
                        </el-col>
                    </el-row>
                </template>  
            </el-table-column>
            
            <el-table-column
                    label="Reprise"   
                    width="100px"
            >
                <template slot-scope="scope">
                    <el-checkbox v-model="scope.row.ToBeFixed"></el-checkbox>
                </template>  
            </el-table-column>
            
            <el-table-column
                    label="Description"                            
            >
                <template slot-scope="scope">
                    <el-input clearable placeholder="Description de la reprise" size="mini" type="textarea" autosize
                              v-model.trim="scope.row.Description"
                    ></el-input>
                </template>  
            </el-table-column>
            
        </el-table>
    </el-main>
</el-container>

`
