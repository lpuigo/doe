package poletable

const template_old string = `<el-container  style="height: 100%; padding: 0px">
    <!--
    <el-header style="height: auto; margin-top: 5px">
        <el-row type="flex" align="middle" :gutter="5">
			<el-col :span="2" style="text-align: right"><span>Mode d'affichage:</span></el-col>
			<el-col :span="2">
			  <el-select v-model="columns.Mode" placeholder="Select" @change="ApplyColumnMode(columns)" size="mini">
				<el-option
				  v-for="item in columns.Refs"
				  :key="item"
				  :label="item"
				  :value="item">
				</el-option>
			  </el-select>
			</el-col>
            <el-col :span="3">
				<el-popover
					placement="right"
					title="Affichage des colonnes"
					width="200"
					trigger="click"
				>
					<div v-for="name in columns.Columns">
					    <el-checkbox :label="name" v-model="columns.Show[name]"></el-checkbox>
					</div>
					<el-button slot="reference" plain size="mini">Choix des colonnes</el-button>            
				</el-popover>
            </el-col>
        </el-row>
    </el-header>
    <div style="height: 100%;overflow-x: hidden;overflow-y: auto;padding: 0px 0px; margin-top: 8px">
    -->
		<el-table
				:data="filteredPoles"
				:row-class-name="TableRowClassName"
				:default-sort = "{prop: 'Ref', order: 'ascending'}"
				height="100%"
				:border=true size="mini"
				@current-change="SetSelectedPole"
		>
			<!--  :sort-method="SortStatus" :sort-by="['Status', 'Client', 'City', 'Ref']"  -->
			<el-table-column
					label="Ref" prop="Ref" sortable
					width="100px" :resizable=true :show-overflow-tooltip=true
			></el-table-column>
	
			<el-table-column v-if="columns.Show['Appui']"
					label="Appui" prop="Sticker" sortable
					width="100px" :resizable=true :show-overflow-tooltip=true
			></el-table-column>
	
			<el-table-column v-if="columns.Show['Ville']"
					label="Ville" prop="City" sortable :sort-by="['City', 'Ref']"
					width="100px" :resizable=true :show-overflow-tooltip=true
					:filters="FilterList('City')"	:filter-method="FilterHandler"	filter-placement="bottom-end"
			></el-table-column>
	
			<el-table-column v-if="columns.Show['Adresse']"
					label="Adresse" prop="Address"
					width="160px" :resizable=true :show-overflow-tooltip=true
			></el-table-column>
	
			<el-table-column v-if="columns.Show['DT']"
					label="DT" prop="DtRef" sortable
					width="130px" :resizable=true :show-overflow-tooltip=true
			></el-table-column>
	
			<el-table-column v-if="columns.Show['DICT']"
					label="DICT" prop="DictRef" sortable
					width="120px" :resizable=true :show-overflow-tooltip=true
			></el-table-column>
	
			<el-table-column v-if="columns.Show['Déb.Trx']"
					label="Déb.Trx" prop="DictDate" sortable
					width="90px" :resizable=true
					align="center"	:formatter="FormatDate"
			></el-table-column>
	
			<el-table-column v-if="columns.Show['Info DICT']"
					label="Info DICT" prop="DictInfo"
					width="100px" :resizable=true :show-overflow-tooltip=true
			></el-table-column>
	
			<el-table-column v-if="columns.Show['Aspi.']"
					label="Aspi." prop="AspiDate" sortable
					width="90px" :resizable=true
					align="center"	:formatter="FormatDate"
			></el-table-column>
	
			<el-table-column v-if="columns.Show['Type']"
					label="Type" prop="Material"
					width="80px" :resizable=true :show-overflow-tooltip=true
			>
				<template slot-scope="scope">
					<span>{{FormatType(scope.row)}}</span>
				</template>
			</el-table-column>
	
			<!--
			<el-table-column
					label="Matière" prop="Material"
					width="80px" :resizable=true :show-overflow-tooltip=true
			></el-table-column>
			-->
	
			<el-table-column v-if="columns.Show['Produits']"
					label="Produits"
					width="110px" :resizable=true
			>
				<template slot-scope="scope">
					<span style="white-space: pre">{{FormatProduct(scope.row)}}</span>
				</template>
			</el-table-column>
	
			<el-table-column v-if="columns.Show['Acteurs']"
					label="Acteurs"
					width="90px" :resizable=true :show-overflow-tooltip=true
			>
				<template slot-scope="scope">
					<span style="white-space: pre">{{FormatActors(scope.row)}}</span>
				</template>
			</el-table-column>
	
			<el-table-column v-if="columns.Show['Statut']"
					label="Statut" prop="State" :formatter="FormatState" sortable :sort-method="SortState"
					width="100px" :resizable=true :show-overflow-tooltip=true
					:filters="FilterList('State')"	:filter-method="FilterHandler"	filter-placement="bottom-end" :filtered-value="FilteredStatusValue()"
			></el-table-column>
	
			<el-table-column v-if="columns.Show['Ref. Kizeo']"
					label="Ref. Kizeo" prop="Kizeo"
					width="80px" :resizable=true :show-overflow-tooltip=true
			></el-table-column>
	
			<el-table-column v-if="columns.Show['Date']"
					label="Date" prop="Date" sortable
					width="100px" :resizable=true
					align="center"	:formatter="FormatDate"
			></el-table-column>
	
			<el-table-column v-if="columns.Show['Attachement']"
					label="Attachement" prop="AttachmentDate" sortable
					width="110px" :resizable=true
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
	<!--
    </div>
    -->		
</el-container>
`

const template string = `<el-container  style="height: 100%; padding: 0px">
    <el-table
            :data="filteredPoles"
            :row-class-name="TableRowClassName"
            :default-sort = "{prop: 'Ref', order: 'ascending'}"
            height="100%"
            :border=true size="mini"
            @current-change="SetSelectedPole"
    >
        <!--  :sort-method="SortStatus" :sort-by="['Status', 'Client', 'City', 'Ref']"  -->
        <el-table-column
                label="Ref" prop="Ref" sortable
                width="100px" :resizable=true :show-overflow-tooltip=true
        ></el-table-column>

        <el-table-column
                label="Appui" prop="Sticker" sortable
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
                label="Aspi." prop="AspiDate" sortable
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
                label="Acteurs"
                width="90px" :resizable=true :show-overflow-tooltip=true
        >
            <template slot-scope="scope">
                <span style="white-space: pre">{{FormatActors(scope.row)}}</span>
            </template>
        </el-table-column>

        <el-table-column
                label="Statut" prop="State" :formatter="FormatState" sortable :sort-method="SortState"
                width="100px" :resizable=true :show-overflow-tooltip=true
                :filters="FilterList('State')"	:filter-method="FilterHandler"	filter-placement="bottom-end" :filtered-value="FilteredStatusValue()"
        ></el-table-column>

        <el-table-column
                label="Ref. Kizeo" prop="Kizeo"
                width="80px" :resizable=true :show-overflow-tooltip=true
        ></el-table-column>

        <el-table-column
                label="Date" prop="Date" sortable
                width="100px" :resizable=true
                align="center"	:formatter="FormatDate"
        ></el-table-column>

        <el-table-column
                label="Attachement" prop="AttachmentDate" sortable
                width="110px" :resizable=true
                align="center"	:formatter="FormatDate"
        ></el-table-column>

        <el-table-column 
                label="Commentaire" prop="Comment"
                min-width="120px" :resizable=true
        ></el-table-column>
    </el-table>
</el-container>
`
