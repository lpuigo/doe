package poletable

const template_creation string = `<el-container  style="height: 100%; padding: 0px">
    <el-table
            :data="filteredPoles"
            :row-class-name="TableRowClassName"
            :default-sort = "{prop: 'Ref', order: 'ascending'}"
            height="100%"
            :border=true size="mini"
            @current-change="SetSelectedPole"
    >
		<el-table-column type="index"></el-table-column>
		<el-table-column
                label="Ref" prop="Ref" sortable :sort-by="PoleRefName"
                width="150px" :resizable=true :show-overflow-tooltip=true
        >
			<template slot-scope="scope">
				<span>{{PoleRefName(scope.row)}}</span>
			</template>
		</el-table-column>

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
                label="Déb DICT" prop="DictDate" sortable :sort-by="SortDate('DictDate')"
                width="92px" :resizable=true
                align="center"	:formatter="FormatDate"
        ></el-table-column>

<!--        <el-table-column-->
<!--                label="Fin DICT" sortable :sort-by="SortDate('DictDate')"-->
<!--                width="100px" :resizable=true-->
<!--                align="center"-->
<!--        >-->
<!--            <template slot-scope="scope">-->
<!--                <span>{{DictEndDate(scope.row.DictDate)}}</span>-->
<!--            </template>-->
<!--        </el-table-column>-->

<!--        <el-table-column-->
<!--                label="Info DICT" prop="DictInfo"-->
<!--                width="100px" :resizable=true :show-overflow-tooltip=true-->
<!--        ></el-table-column>-->

        <el-table-column
                label="Dem. DA" prop="DaQueryDate" sortable :sort-by="SortDate('DaQueryDate')"
                width="90px" :resizable=true
                align="center"	:formatter="FormatDate"
        ></el-table-column>

        <el-table-column
                label="Déb DA" prop="DaStartDate" sortable :sort-by="SortDate('DaStartDate')"
                width="90px" :resizable=true
                align="center"	:formatter="FormatDate"
        ></el-table-column>

        <el-table-column
                label="Fin DA" prop="DaEndDate" sortable :sort-by="SortDate('DaEndDate')"
                width="90px" :resizable=true
                align="center"	:formatter="FormatDate"
        ></el-table-column>

        <el-table-column
                label="Type" prop="Material"
                width="100px" :resizable=true :show-overflow-tooltip=true
        >
            <template slot-scope="scope">
                <span>{{FormatType(scope.row)}}</span>
            </template>
        </el-table-column>

        <el-table-column
                label="Produits"
                width="110px" :resizable=true
        >
            <template slot-scope="scope">
                <span style="white-space: pre">{{FormatProduct(scope.row)}}</span>
            </template>
        </el-table-column>

        <el-table-column
                label="Statut" prop="State" :formatter="FormatState" sortable :sort-method="SortState"
                width="100px" :resizable=true :show-overflow-tooltip=true
                :filters="FilterList('State')"	:filter-method="FilterHandler"	filter-placement="bottom-end" :filtered-value="FilteredStatusValue('creation')"
        ></el-table-column>

		<!--
        <el-table-column
                label="Aspi." prop="AspiDate" sortable
                width="90px" :resizable=true
                align="center"	:formatter="FormatDate"
        ></el-table-column>

        <el-table-column
                label="Date" prop="Date" sortable
                width="100px" :resizable=true
                align="center"	:formatter="FormatDate"
        ></el-table-column>

        <el-table-column
                label="Ref. Kizeo" prop="Kizeo"
                width="80px" :resizable=true :show-overflow-tooltip=true
        ></el-table-column>

        <el-table-column
                label="Acteurs"
                width="90px" :resizable=true :show-overflow-tooltip=true
        >
            <template slot-scope="scope">
                <span style="white-space: pre">{{FormatActors(scope.row)}}</span>
            </template>
        </el-table-column>

        <el-table-column
                label="Attachement" prop="AttachmentDate" sortable
                width="110px" :resizable=true
                align="center"	:formatter="FormatDate"
        ></el-table-column>
		-->
        <el-table-column 
                label="Commentaire" prop="Comment"
                min-width="120px"
        >
            <template slot-scope="scope">
                <span style="white-space: pre">{{FormatComment(scope.row)}}</span>
            </template>
        </el-table-column>
    </el-table>
</el-container>
`

const template_followup string = `<el-container  style="height: 100%; padding: 0px">
    <el-table
            :data="filteredPoles"
            :row-class-name="TableRowClassName"
            :default-sort = "{prop: 'Ref', order: 'ascending'}"
            height="100%"
            :border=true size="mini"
            @current-change="SetSelectedPole"
    >
		<el-table-column type="index"></el-table-column>
        <el-table-column
                label="Ref" prop="Ref" sortable :sort-by="PoleRefName"
                width="150px" :resizable=true :show-overflow-tooltip=true
        >
			<template slot-scope="scope">
				<span>{{PoleRefName(scope.row)}}</span>
			</template>
		</el-table-column>

        <el-table-column
                label="Ville" prop="City" sortable :sort-by="['City', 'Ref']"
                width="100px" :resizable=true :show-overflow-tooltip=true
                :filters="FilterList('City')"	:filter-method="FilterHandler"	filter-placement="bottom-end"
        ></el-table-column>

        <el-table-column
                label="Adresse" prop="Address"
                width="160px" :resizable=true :show-overflow-tooltip=true
        ></el-table-column>

		<!--
        <el-table-column
                label="DT" prop="DtRef" sortable
                width="130px" :resizable=true :show-overflow-tooltip=true
        ></el-table-column>
		-->

        <el-table-column
                label="DICT" prop="DictRef" sortable
                width="120px" :resizable=true :show-overflow-tooltip=true
        ></el-table-column>

        <el-table-column
                label="Déb DICT" prop="DictDate" sortable :sort-by="SortDate('DictDate')"
                width="92px" :resizable=true
                align="center"	:formatter="FormatDate"
        ></el-table-column>

        <el-table-column
                label="Fin DICT" sortable :sort-by="SortDate('DictDate')"
                width="90px" :resizable=true
                align="center"	:formatter="FormatDate"
        >
            <template slot-scope="scope">
                <span>{{DictEndDate(scope.row.DictDate)}}</span>
            </template>
        </el-table-column>

<!--        <el-table-column-->
<!--                label="Info DICT" prop="DictInfo"-->
<!--                width="100px" :resizable=true :show-overflow-tooltip=true-->
<!--        ></el-table-column>-->

        <el-table-column
                label="Dem. DA" prop="DaQueryDate" sortable :sort-by="SortDate('DaQueryDate')"
                width="90px" :resizable=true
                align="center"	:formatter="FormatDate"
        ></el-table-column>

        <el-table-column
                label="Déb DA" prop="DaStartDate" sortable :sort-by="SortDate('DaStartDate')"
                width="90px" :resizable=true
                align="center"	:formatter="FormatDate"
        ></el-table-column>

        <el-table-column
                label="Fin DA" prop="DaEndDate" sortable :sort-by="SortDate('DaEndDate')"
                width="90px" :resizable=true
                align="center"	:formatter="FormatDate"
        ></el-table-column>

        <el-table-column
                label="Type" prop="Material"
                width="80px" :resizable=true :show-overflow-tooltip=true
        >
            <template slot-scope="scope">
                <span>{{FormatType(scope.row)}}</span>
            </template>
        </el-table-column>

        <el-table-column
                label="Produits"
                width="110px" :resizable=true
        >
            <template slot-scope="scope">
                <span style="white-space: pre">{{FormatProduct(scope.row)}}</span>
            </template>
        </el-table-column>

        <el-table-column
                label="Aspi." prop="AspiDate" sortable
                width="90px" :resizable=true
                align="center"	:formatter="FormatDate"
        ></el-table-column>

        <el-table-column
                label="Statut" prop="State" :formatter="FormatState" sortable :sort-method="SortState"
                width="100px" :resizable=true :show-overflow-tooltip=true
                :filters="FilterList('State')"	:filter-method="FilterHandler"	filter-placement="bottom-end" :filtered-value="FilteredStatusValue('followup')"
        ></el-table-column>

        <el-table-column
                label="Date" prop="Date" sortable
                width="100px" :resizable=true
                align="center"	:formatter="FormatDate"
        ></el-table-column>

		<!--
        <el-table-column
                label="Ref. Kizeo" prop="Kizeo"
                width="80px" :resizable=true :show-overflow-tooltip=true
        ></el-table-column>

        <el-table-column
                label="Acteurs"
                width="90px" :resizable=true :show-overflow-tooltip=true
        >
            <template slot-scope="scope">
                <span style="white-space: pre">{{FormatActors(scope.row)}}</span>
            </template>
        </el-table-column>

        <el-table-column
                label="Attachement" prop="AttachmentDate" sortable
                width="110px" :resizable=true
                align="center"	:formatter="FormatDate"
        ></el-table-column>
		-->

        <el-table-column 
                label="Commentaire" prop="Comment"
                min-width="120px"
        >
            <template slot-scope="scope">
                <span style="white-space: pre">{{FormatComment(scope.row)}}</span>
            </template>
        </el-table-column>
    </el-table>
</el-container>
`

const template_billing string = `<el-container  style="height: 100%; padding: 0px">
    <el-table
            :data="filteredPoles"
            :row-class-name="TableRowClassName"
            :default-sort = "{prop: 'Ref', order: 'ascending'}"
            height="100%"
            :border=true size="mini"
            @current-change="SetSelectedPole"
    >
		<el-table-column type="index"></el-table-column>
        <el-table-column
                label="Ref" prop="Ref" sortable :sort-by="PoleRefName"
                width="150px" :resizable=true :show-overflow-tooltip=true
        >
			<template slot-scope="scope">
				<span>{{PoleRefName(scope.row)}}</span>
			</template>
		</el-table-column>

        <el-table-column
                label="Ville" prop="City" sortable :sort-by="['City', 'Ref']"
                width="100px" :resizable=true :show-overflow-tooltip=true
                :filters="FilterList('City')"	:filter-method="FilterHandler"	filter-placement="bottom-end"
        ></el-table-column>

        <el-table-column
                label="Adresse" prop="Address"
                width="160px" :resizable=true :show-overflow-tooltip=true
        ></el-table-column>

		<!--
        <el-table-column
                label="DT" prop="DtRef" sortable
                width="130px" :resizable=true :show-overflow-tooltip=true
        ></el-table-column>

        <el-table-column
                label="DICT" prop="DictRef" sortable
                width="120px" :resizable=true :show-overflow-tooltip=true
        ></el-table-column>

        <el-table-column
                label="Déb DICT" prop="DictDate" sortable :sort-by="SortDate('DictDate')"
                width="92px" :resizable=true
                align="center"	:formatter="FormatDate"
        ></el-table-column>

        <el-table-column
                label="Fin DICT" sortable :sort-by="SortDate('DictDate')"
                width="90px" :resizable=true
                align="center"	:formatter="FormatDate"
        >
            <template slot-scope="scope">
                <span>{{DictEndDate(scope.row.DictDate)}}</span>
            </template>
        </el-table-column>

        <el-table-column
                label="Info DICT" prop="DictInfo"
                width="100px" :resizable=true :show-overflow-tooltip=true
        ></el-table-column>

        <el-table-column
                label="Déb DA" prop="DaStartDate" sortable :sort-by="SortDate('DaStartDate')"
                width="90px" :resizable=true
                align="center"	:formatter="FormatDate"
        ></el-table-column>

        <el-table-column
                label="Fin DA" prop="DaEndDate" sortable :sort-by="SortDate('DaEndDate')"
                width="90px" :resizable=true
                align="center"	:formatter="FormatDate"
        ></el-table-column>
		-->
        <el-table-column
                label="Type" prop="Material"
                width="80px" :resizable=true :show-overflow-tooltip=true
        >
            <template slot-scope="scope">
                <span>{{FormatType(scope.row)}}</span>
            </template>
        </el-table-column>

        <el-table-column
                label="Produits"
                width="110px" :resizable=true
        >
            <template slot-scope="scope">
                <span style="white-space: pre">{{FormatProduct(scope.row)}}</span>
            </template>
        </el-table-column>

        <el-table-column
                label="Aspi." prop="AspiDate" sortable
                width="90px" :resizable=true
                align="center"	:formatter="FormatDate"
        ></el-table-column>

        <el-table-column
                label="Statut" prop="State" :formatter="FormatState" sortable :sort-method="SortState"
                width="100px" :resizable=true :show-overflow-tooltip=true
                :filters="FilterList('State')"	:filter-method="FilterHandler"	filter-placement="bottom-end" :filtered-value="FilteredStatusValue('billing')"
        ></el-table-column>

        <el-table-column
                label="Date" prop="Date" sortable
                width="90px" :resizable=true
                align="center"	:formatter="FormatDate"
        ></el-table-column>

        <el-table-column
                label="Ref. Kizeo" prop="Kizeo"
                width="80px" :resizable=true :show-overflow-tooltip=true
        ></el-table-column>

        <el-table-column
                label="Acteurs"
                width="90px" :resizable=true :show-overflow-tooltip=true
        >
            <template slot-scope="scope">
                <span style="white-space: pre">{{FormatActors(scope.row)}}</span>
            </template>
        </el-table-column>

        <el-table-column
                label="Attachement" prop="AttachmentDate" sortable
                width="110px" :resizable=true
                align="center"	:formatter="FormatDate"
        ></el-table-column>

        <el-table-column 
                label="Commentaire" prop="Comment"
                min-width="120px"
        >
            <template slot-scope="scope">
                <span style="white-space: pre">{{FormatComment(scope.row)}}</span>
            </template>
        </el-table-column>
    </el-table>
</el-container>
`
