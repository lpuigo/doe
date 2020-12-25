package adminmodal

const template string = `<el-dialog
		:visible.sync="visible" 
		:before-close="HideWithControl" :close-on-click-modal="false"
		width="80%" top="5vh"
>
	<!-- 	Modal Title	-->
	<span slot="title">
		<h2 style="margin: 0 0">
			<i class="fas fa-wrench icon--left"></i>Administration
		</h2>
	</span>

	<!-- 
		Modal Body
		style="height: 100%;"		
	-->
    <el-tabs type="border-card" tab-position="left" style="height: 70vh">
		<!-- ========================================== Admin Tab ================================================= -->
		<el-tab-pane label="Maintenance" lazy=true style="height: 75vh; padding: 5px 25px; overflow-x: hidden;overflow-y: auto;">
			<el-button type="primary" @click="ReloadData">Rechargement des données</el-button>
			<h4>&nbsp;</h4>
			<h3>Archive des <a href="/api/worksites/archive">Chantiers Orange</a></h3>
			<h3>Archive des <a href="/api/ripsites/archive">Chantiers RIP</a></h3>
			<h3>Archive des <a href="/api/polesites/archive">Chantiers Poteaux</a></h3>
			<h3>Archive des <a href="/api/foasites/archive">Chantiers FOAs</a></h3>
			<h4>&nbsp;</h4>
			<h3>Archive des <a href="/api/clients/archive">Clients</a></h3>
			<h3>Archive des <a href="/api/actors/archive">Equipes</a></h3>
			<h3>Archive des <a href="/api/groups/archive">Groupes</a></h3>
			<h3>Archive des <a href="/api/actorinfos/archive">Infos RH Equipes</a></h3>
			<h3>Archive des <a href="/api/timesheet/archive">Pointages Horaire</a></h3>
		</el-tab-pane>

		<!-- ========================================== Users Tab ================================================= -->
		<el-tab-pane label="Utilisateurs" lazy=true style="height: 75vh; padding: 5px 25px; overflow-x: hidden;overflow-y: auto;">
			<el-table
					:border=false
					:data="filteredUsers"
					:default-sort = "{prop: 'Name', order: 'ascending'}"
					:row-class-name="TableRowClassName" height="90%" size="mini"
			>
				<!--	Edit User-->
				<el-table-column type="expand" width="40px" >
					<template slot-scope="scope">
						<!--	User Name & Password -->
						<el-row :gutter="5" type="flex" align="middle" class="spaced">
							<el-col :span="2" class="align-right">Nom:</el-col>
							<el-col :span="8">
								<el-input v-model="scope.row.Name" size="mini"></el-input>
							</el-col>
							<el-col :span="2" class="align-right">MdP:</el-col>
							<el-col :span="8">
								<el-input v-model="scope.row.Password" size="mini"></el-input>
							</el-col>
						</el-row>
						<!--	User Permission -->
						<el-row :gutter="5" type="flex" align="middle" class="spaced">
							<el-col :span="2" class="align-right">Permissions:</el-col>
							<el-col :span="6" >
								<p><el-switch v-model="scope.row.Permissions.Update" active-text="Peut faire des mises à jour"></el-switch></p>
								<p><el-switch v-model="scope.row.Permissions.Create" active-text="Peut créer des éléments"></el-switch></p>
							</el-col>
							<el-col :span="6" >
								<p><el-switch v-model="scope.row.Permissions.HR" active-text="Accès aux infos RH"></el-switch></p>
								<p><el-switch v-model="scope.row.Permissions.Invoice" active-text="Accès aux infos financières"></el-switch></p>
							</el-col>
							<el-col :span="6" >
								<p><el-switch v-model="scope.row.Permissions.Admin" active-text="Accès aux fonctions d'administration"></el-switch></p>
								<p><el-switch v-model="scope.row.Permissions.Review" active-text="Accès restreint en lecture seulement"></el-switch></p>
							</el-col>
						</el-row>
						<el-row :gutter="5" type="flex" align="middle" class="doublespaced">
							<el-col :span="2" class="align-right">Groupes:</el-col>
							<el-col :span="18">
								<el-select multiple placeholder="Tous les Groupes" size="mini"
										   v-model="scope.row.Groups"
										   style="width: 100%"
								>
									<el-option v-for="item in GetGroupList()"
											   :key="item.value"
											   :label="item.label"
											   :value="item.value"
									>
									</el-option>
								</el-select>
							</el-col>
						</el-row>
						<el-row v-if="scope.row.Groups.length > 0" :gutter="5" type="flex" align="middle" class="spaced">
							<el-col :span="3" class="align-right">Clients visibles:</el-col>
							<el-col :span="18">
								<span>{{ClientList(scope.row)}}</span>
							</el-col>
						</el-row>
						<el-row :gutter="5" type="flex" align="middle" class="doublespaced">
							<el-col :span="2" class="align-right">Clients:</el-col>
							<el-col :span="18">
								<el-select multiple placeholder="Tous les Clients" size="mini"
										   v-model="scope.row.Clients"
										   @change="UpdateUserClients(scope.row)"
										   :disabled="scope.row.Groups.length > 0"
										   style="width: 100%"
								>
									<el-option v-for="item in GetClientList()"
											   :key="item.value"
											   :label="item.label"
											   :value="item.value"
									>
									</el-option>
								</el-select>
							</el-col>
						</el-row>
					</template>
				</el-table-column>

    			<!--	Index   -->
				<el-table-column
						label="N°" width="40px" align="right"
						type="index"
						index=1 
				></el-table-column>
			
				<!--	Actions   -->
				<el-table-column label="" width="80px">
					<template slot="header" slot-scope="scope">
						<el-button type="success" plain icon="fas fa-users fa-fw" size="mini" @click="AddNewUser()"></el-button>
					</template>
				</el-table-column>
				
				<!--	User Name   -->
				<el-table-column
						:resizable="true" :show-overflow-tooltip=true 
						prop="Name" label="Utilisateur" width="210px"
						sortable :sort-by="['Name']"
				></el-table-column>
				<!-- :filters="FilterList('Name')" :filter-method="FilterHandler"	filter-placement="bottom-end"-->
				
				<!--	groups -->   
				<el-table-column
						:resizable="true" :show-overflow-tooltip=true 
						prop="Groups" label="Groupes" width="300px"
						:filters="FilterList('Group')" :filter-method="FilterHandler"	filter-placement="bottom-end"
				>
					<template slot-scope="scope">
						<span>{{GroupList(scope.row)}}</span>
					</template>
				</el-table-column>
			
				<!--	clients -->   
				<el-table-column
						:resizable="true" :show-overflow-tooltip=true 
						prop="Clients" label="Clients" width="300px"
						:filters="FilterList('Client')" :filter-method="FilterHandler"	filter-placement="bottom-end"
				>
					<template slot-scope="scope">
						<span>{{ClientList(scope.row)}}</span>
					</template>
				</el-table-column>
			
				<!--	Permissions   -->
				<el-table-column
						:resizable="true"
						label="Permissions"
				>
					<template slot-scope="scope">
						<span>
							<i class="fas fa-edit icon--medium icon--left" :class="{ 'icon--disabled': !scope.row.Permissions.Update }"></i>
							<i class="fas fa-plus-circle icon--medium icon--left" :class="{ 'icon--disabled': !scope.row.Permissions.Create }"></i>
							<i class="fas fa-id-card icon--medium icon--left" :class="{ 'icon--disabled': !scope.row.Permissions.HR }"></i>
							<i class="fas fa-euro-sign icon--medium icon--left" :class="{ 'icon--disabled': !scope.row.Permissions.Invoice }"></i>
							<i class="fas fa-tools icon--medium icon--left" :class="{ 'icon--disabled': !scope.row.Permissions.Admin }"></i>
							<i class="fas fa-eye icon--medium" :class="{ 'icon--disabled': !scope.row.Permissions.Review }"></i>
						</span>
					</template>
				</el-table-column>
			</el-table>
		</el-tab-pane>
    </el-tabs>

    <!-- 
        Modal Footer Action Bar
    -->
    <span slot="footer">
		<el-row :gutter="15">
			<el-col :span="24" style="text-align: right">
				<el-tooltip :open-delay="500" effect="light">
					<div slot="content">Annuler les changements</div>
					<el-button :disabled="!hasChanged" @click="UndoChange" icon="fas fa-undo-alt" plain size="mini"
                               type="info"></el-button>
				</el-tooltip>
				
				<el-button @click="Hide" size="mini">Fermer</el-button>
				
				<el-button :disabled="!hasChanged" type="success" @click="ConfirmChange" plain size="mini">Enregistrer</el-button>
			</el-col>
		</el-row>
	</span>
</el-dialog>`
