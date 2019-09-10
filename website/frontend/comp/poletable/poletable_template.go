package poletable

const template string = `<el-main  style="height: 100%; padding: 0px">
    <el-table
            :data="filteredPoles"
            :row-class-name="TableRowClassName"
            :default-sort = "{prop: 'Ref', order: 'ascending'}"
            height="100%"
            :border=true size="mini"
			highlight-current-row
			@current-change="SetSelectedPole"
    >
        <!--  :sort-method="SortStatus" :sort-by="['Status', 'Client', 'City', 'Ref']"  -->
        <el-table-column
                fixed label="Ref" prop="Ref" sortable
                width="100px" :resizable=true :show-overflow-tooltip=true
        ></el-table-column>

        <el-table-column
                label="Ville" prop="City" sortable :sort-by="['City', 'Ref']"
                width="100px" :resizable=true :show-overflow-tooltip=true
                :filters="FilterList('City')"	:filter-method="FilterHandler"	filter-placement="bottom-end"
        ></el-table-column>

        <el-table-column
                label="Adresse" prop="Address"
                width="160px" :resizable=true :show-overflow-tooltip=true
        ></el-table-column>

        <el-table-column
                label="DT" prop="DtRef" sortable
                width="130px" :resizable=true :show-overflow-tooltip=true
        ></el-table-column>

        <el-table-column
                label="DICT" prop="DictRef" sortable
                width="120px" :resizable=true :show-overflow-tooltip=true
        ></el-table-column>

        <el-table-column
                label="Déb.Trx" prop="DictDate" sortable
                width="90px" :resizable=true
                align="center"	:formatter="FormatDate"
        ></el-table-column>

        <el-table-column
                label="Info DICT" prop="DictInfo"
                width="100px" :resizable=true :show-overflow-tooltip=true
        ></el-table-column>

        <el-table-column
                label="Hauteur" prop="Height"
                width="80px" :resizable=true :show-overflow-tooltip=true
        ></el-table-column>

        <el-table-column
                label="Matière" prop="Material"
                width="80px" :resizable=true :show-overflow-tooltip=true
        ></el-table-column>

        <el-table-column
                label="Aspi." prop="AspiDate" sortable
                width="90px" :resizable=true
                align="center"	:formatter="FormatDate"
        ></el-table-column>

        <el-table-column
                label="Produits"
                width="90px" :resizable=true
        >
            <template slot-scope="scope">
                <span style="white-space: pre">{{FormatProduct(scope.row)}}</span>
            </template>
        </el-table-column>

        <el-table-column
                label="Acteurs"
                width="90px" :resizable=true :show-overflow-tooltip=true
        >
            <template slot-scope="scope">
                <span style="white-space: pre">{{FormatActors(scope.row)}}</span>
            </template>
        </el-table-column>

        <el-table-column
                label="Ref. Kizeo" prop="Kizeo"
                width="80px" :resizable=true :show-overflow-tooltip=true
        ></el-table-column>

        <el-table-column
                label="Statut" prop="State" :formatter="FormatState" sortable :sort-method="SortState"
                width="100px" :resizable=true :show-overflow-tooltip=true
                :filters="FilterList('State')"	:filter-method="FilterHandler"	filter-placement="bottom-end" :filtered-value="FilteredStatusValue()"
        ></el-table-column>

        <el-table-column
                label="Date" prop="Date" sortable
                width="90px" :resizable=true
                align="center"	:formatter="FormatDate"
        ></el-table-column>

<!--
        <el-table-column
                label="Référence" sortable
                width="120px" :resizable=true :show-overflow-tooltip=true
        >        
            <template slot-scope="scope">
                <div class="header-menu-container">
                    <span @click="OpenPolesite(scope.row.Id)" class="link">{{scope.row.Ref}}</span>
                    <a v-if="user.Permissions.Invoice" :href="AttachmentUrl(scope.row.Id)"><i class="link fas fa-file-excel"></i></a>
                </div>
            </template>
        </el-table-column>
-->
        <el-table-column
                label="Commentaire" prop="Comment"
                min-width="120px" :resizable=true
        ></el-table-column>
    </el-table>		
</el-main>
`
