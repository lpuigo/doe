package vehiculeupdatemodal

const template string = `<el-dialog
        :before-close="HideWithControl"
        :visible.sync="visible" :close-on-click-modal="false"
        width="85%" top="5vh"
>
    <!-- 
        Modal Title
    -->
    <span slot="title">
		<el-row :gutter="10" align="middle" type="flex">
			<el-col :span="12">
				<h2 style="margin: 0 0" v-if="current_vehicule">
					<i class="far fa-edit icon--left"></i>Edition du véhicule : <span style="color: #ccebff">{{current_vehicule.Type}} {{current_vehicule.Immat}}</span>
				</h2>
			</el-col>
		</el-row>
    </span>
    <!-- 
        Modal Body
        style="height: 100%;"
    -->
    <el-tabs type="border-card" tab-position="left" style="height: 75vh">
		<!-- ===================================== Vehicule Tab ======================================================= -->
		<el-tab-pane v-if="user.Permissions.Update" label="Vehicule" lazy=true style="height: 75vh; padding: 5px 25px; overflow-x: hidden;overflow-y: auto;">
	        <!-- Immat & Type -->
			<el-row :gutter="10" align="middle" class="doublespaced" type="flex">
				<el-col :span="4" class="align-right">Immatriculation :</el-col>
				<el-col :span="8">
					<el-input @change=""
							  clearable placeholder="AB-123-CD" size="mini"
							  v-model="current_vehicule.Immat"
					></el-input>
  				</el-col>
	
				<el-col :span="4" class="align-right">Type :</el-col>
				<el-col :span="8">
  					<el-select v-model="current_vehicule.Type" placeholder="Select" size="mini" style="width: 100%">
						<el-option
							v-for="item in GetVehiculeType()"
							:key="item.value"
							:label="item.label"
							:value="item.value">
						</el-option>
  					</el-select>
  				</el-col>
			</el-row>

			<!-- Company & Model -->
			<el-row :gutter="10" align="middle" class="spaced" type="flex">
				<el-col :span="4" class="align-right">Compagnie :</el-col>
				<el-col :span="8">
					<el-input @change="CheckCompany"
							  clearable placeholder="Compagnie" size="mini"
							  v-model="current_vehicule.Company"
					></el-input>
				</el-col>

				<el-col :span="4" class="align-right">Modèle :</el-col>
				<el-col :span="8">
					<el-input @change=""
							  clearable placeholder="modèle" size="mini"
							  v-model="current_vehicule.Model"
					></el-input>
				</el-col>
			</el-row>
	        
	        <!-- Service Dates -->
			<el-row :gutter="10" align="middle" class="doublespaced" type="flex">
				<el-col :span="4" class="align-right">Mise en Service :</el-col>
				<el-col :span="8">
					<el-date-picker v-model="current_vehicule.ServiceDate"
									type="date" :clearable="false"
									format="dd/MM/yyyy" value-format="yyyy-MM-dd"
									:picker-options="{firstDayOfWeek:1, disabledDate(time) { return time.getTime() > Date.now(); }}"
									placeholder="Date" size="mini" style="width: 100%"
									@change=""
					></el-date-picker>
				</el-col>

				<el-col :span="4" class="align-right">Fin de Service :</el-col>
				<el-col :span="8">
					<el-date-picker v-model="current_vehicule.EndServiceDate"
									type="date" :clearable="false"
									format="dd/MM/yyyy" value-format="yyyy-MM-dd"
									:picker-options="{firstDayOfWeek:1, disabledDate(time) { return time.getTime() > Date.now(); }}"
									placeholder="Date" size="mini" style="width: 100%"
									@change=""
					></el-date-picker>
				</el-col>
			</el-row>
			
	        <!-- FuelCard & TravelledKms -->
			<el-row :gutter="10" align="middle" class="doublespaced" type="flex">
				<el-col :span="4" class="align-right">Carte Carburant :</el-col>
				<el-col :span="8">
					<el-input v-model="current_vehicule.FuelCard" 
							  clearable placeholder="info carte" size="mini"
							  @change=""
					></el-input>
				</el-col>

				<el-col :span="4" class="align-right">Kilométrage :</el-col>
				<el-col :span="8">
					<el-table 
							:data="current_vehicule.TravelledKms"
							max-height="160" size="mini" border
					>
						<el-table-column label="" width="80">
							<template slot="header" slot-scope="scope">
								<el-button type="success" plain icon="fas fa-plus fa-fw" size="mini" @click="AddTravelledKms()"></el-button>
							</template>
							<template slot-scope="scope">
								<el-button type="danger" plain icon="far fa-trash-alt fa-fw" size="mini" @click="RemoveTravelledKms(scope.$index)"></el-button>
							</template>
						</el-table-column>
	
						<el-table-column label="Date" width="180">
							<template slot-scope="scope">
								<el-date-picker :picker-options="{firstDayOfWeek:1}" 
											placeholder="Date" size="mini" style="width: 100%"
											type="date" format="dd/MM/yyyy" value-format="yyyy-MM-dd"
											v-model="scope.row.Date"
											@change="UpdateTravelledKms"
								></el-date-picker>
							</template>
						</el-table-column>
	
						<el-table-column label="Kms au compteur">
							<template slot-scope="scope">
								<el-input-number v-model="scope.row.Kms" 
										:min="0" :step="1000" size="mini" style="width: 100%;"
										@change=""
								></el-input-number>
							</template>
						</el-table-column>
					</el-table>
				</el-col>
			</el-row>
			
			<!-- Comment -->
			<el-row :gutter="10" align="top" class="doublespaced" type="flex">
				<el-col :span="4" class="align-right">Commentaire :</el-col>
				<el-col :span="20">
					<el-input type="textarea" :autosize="{ minRows: 2, maxRows: 5}" placeholder="Commentaire"
							  v-model="current_vehicule.Comment" clearable size="mini"
					></el-input>
				</el-col>
			</el-row>		
		
			<!-- Actor History -->
			<el-row :gutter="10" align="middle" class="doublespaced" type="flex">
				<el-col :span="4" class="align-right">Affectation :</el-col>
				<el-col :span="20">
					<el-table 
							:data="current_vehicule.InCharge"
							max-height="160" size="mini" border
					>
						<el-table-column label="" width="80">
							<template slot="header" slot-scope="scope">
								<el-button type="success" plain icon="fas fa-user fa-fw" size="mini" @click="AddInCharge()"></el-button>
							</template>
							<template slot-scope="scope">
								<el-button type="danger" plain icon="far fa-trash-alt fa-fw" size="mini" @click="RemoveInCharge(scope.$index)"></el-button>
							</template>
						</el-table-column>
	
						<el-table-column label="Depuis le" width="180">
							<template slot-scope="scope">
								<el-date-picker :picker-options="{firstDayOfWeek:1}" 
											placeholder="Date" size="mini" style="width: 100%"
											type="date" format="dd/MM/yyyy" value-format="yyyy-MM-dd"
											v-model="scope.row.Date"
											@change="UpdateInCharge"
								></el-date-picker>
							</template>
						</el-table-column>
	
						<el-table-column label="Acteur">
							<template slot-scope="scope">
								<el-select v-model="scope.row.ActorId"
										placeholder="Acteur" size="mini" filterable
										@change="UpdateInCharge" style="width: 100%"
								>
									<el-option v-for="item in GetActors()"
											   :key="item.value"
											   :label="item.label"
											   :value="item.value"
									>
									</el-option>
								</el-select>
							</template>
						</el-table-column>
					</el-table>
				</el-col>
			</el-row>
		</el-tab-pane>	

		<!-- ===================================== Inventory Tab ======================================================= -->
		<el-tab-pane v-if="user.Permissions.Update" label="Inventaire" lazy=true style="height: 75vh; padding: 5px 25px; overflow-x: hidden;overflow-y: auto;">
	        <!-- Inventory list & add inventory -->
			<el-row :gutter="10" align="middle" class="spaced" type="flex">
				<el-col :span="4" class="align-right">Inventaires :</el-col>
				<el-col :span="8">
  					<el-select v-model="InventoryNum" 
  							placeholder="Select" size="mini" style="width: 100%" 
  							@change="UpdateInventoryNum()"
							:disabled="current_vehicule.Inventories.length==0"
  					>
						<el-option
							v-for="item in GetInventoryDates()"
							:key="item.value"
							:label="item.label"
							:value="item.value">
						</el-option>
  					</el-select>
  				</el-col>
	
				<el-col :span="4" class="align-right">Actions :</el-col>
				<el-col :span="4" v-if="current_vehicule.Inventories.length==0">
					<el-popover placement="bottom" title="Nouvel inventaire sur le modele d'un autre véhicule :"
								trigger="click"
								width="400"
								v-model="AddInventoryModelVisible"
					>
						<div>
							<el-switch style="margin-bottom: 8px;"
								v-model="AddInventoryModelSameType"
								active-text="Même type"
								inactive-text="Tout type">
							</el-switch>
							<el-row :gutter="10" align="middle" class="spaced" type="flex">
								<el-col :span="8" class="align-right">Modèle :</el-col>
								<el-col :span="16">
									<el-select v-model="AddInventoryModelVehicId" filterable
											   size="mini"
											   placeholder="Inventaire Véhicule"
											   style="width: 100%;"
									>
										<el-option
												v-for="item in GetInventoryModelVehiculeId()"
												:key="item.value"
												:label="item.label"
												:value="item.value"
										>
										</el-option>
									</el-select>
								</el-col>
							</el-row>
							<div v-if="AddInventoryModelVehicId >= 0" >
								<el-row >
									<el-col :span="20">Article</el-col>
									<el-col :span="4">Quantité</el-col>
								</el-row>
								<div style="max-height: 45vh; overflow-x: hidden;overflow-y: auto;">
									<el-row v-for="item in GetInventoryModelItems()"
											:key="item.Name"										
									>
										<el-col :span="20">{{item.Name}}</el-col>
										<el-col :span="4">{{item.ReferenceQuantity}}</el-col>
									</el-row>
								</div>
							</div>
						</div>
						<el-row :gutter="10" align="middle" class="doublespaced" type="flex" >
							<el-col :offset="18" :span="6">
								<el-button type="success" plain size="mini" @click="AddInventoryModel">Valider</el-button>
							</el-col>
						</el-row>
						<el-button slot="reference" type="success" plain size="mini">Modele d'inventaire</el-button>
					</el-popover>
  				</el-col>
				<el-col :span="4" v-if="user.Permissions.Admin">
					<el-button type="warning" plain size="mini" @click="DeleteInventory" :disabled="current_vehicule.Inventories.length==0">Supprimer l'inventaire</el-button>
  				</el-col>
			</el-row>

	        <!-- Current Inventory -->
			<div v-if="InventoryNum >= 0">
				<el-divider content-position="left">Inventaire du {{FormatDate(currentInventory.ReferenceDate)}}</el-divider>
				<!-- Dates -->
				<el-row :gutter="10" align="middle" class="doublespaced" type="flex">
					<el-col :span="2" class="align-right">Date inventaire :</el-col>
					<el-col :span="4">
						<el-date-picker :picker-options="{firstDayOfWeek:1}" 
									placeholder="Date" size="mini" style="width: 100%"
									type="date" format="dd/MM/yyyy" value-format="yyyy-MM-dd"
									v-model="currentInventory.ReferenceDate"
									@change=""
									:disabled="Control"
						></el-date-picker>
					</el-col>
		
					<el-col :span="3" class="align-right">
						<el-checkbox v-model="Control" size="mini" :disabled="InventoryNum == 0 && !user.Permissions.Admin">Date de contrôle :</el-checkbox>
					</el-col>
					<el-col :span="4">
						<el-date-picker :picker-options="{firstDayOfWeek:1}" 
									placeholder="Date" size="mini" style="width: 100%"
									type="date" format="dd/MM/yyyy" value-format="yyyy-MM-dd"
									v-model="currentInventory.ControledDate"
									@change=""
									:disabled="!Control"
						></el-date-picker>
					</el-col>
					<el-col v-if="Control && InventoryNum == 0" :span="4">
						<el-button :disabled="current_vehicule.Inventories.length==0" type="warning" @click="ValidateInventoryControl" plain size="mini">Valider le contrôle</el-button>
					</el-col>
				</el-row>
				<!-- Comment -->
				<el-row :gutter="10" align="top" class="doublespaced" type="flex">
					<el-col :span="2" class="align-right">Commentaire :</el-col>
					<el-col :span="22">
						<el-input type="textarea" :autosize="{ minRows: 2, maxRows: 5}" placeholder="Commentaire"
								  v-model="currentInventory.Comment" clearable size="mini"
						></el-input>
					</el-col>
				</el-row>	
					
				<!-- Items -->
				<el-row :gutter="10" align="middle" class="doublespaced" type="flex">
					<el-col :span="2" class="align-right">Matériel :</el-col>
					<el-col :span="22">
						<el-table 
								:data="currentInventory.Items"
								height="calc(75vh - 250px)" size="mini" border
						>
							<el-table-column label="" width="68">
								<template slot="header" slot-scope="scope">
									<el-button type="success" plain icon="fas fa-plus fa-fw" size="mini" @click="AddItem()" :disabled="Control"></el-button>
								</template>
								<template slot-scope="scope">
									<el-button type="danger" plain icon="far fa-trash-alt fa-fw" size="mini" @click="RemoveItem(scope.$index)" :disabled="Control"></el-button>
								</template>
							</el-table-column>
		
							<el-table-column label="Articles" width="1030">
								<template slot-scope="scope">
									<el-row :gutter="10" align="middle" type="flex">
										<el-col :span="10" >
											<el-input v-model="scope.row.Name"
												  clearable placeholder="Article" size="mini"
												  @change="" :disabled="Control"
											></el-input>
										</el-col>
										
										<el-col :span="12" >
											<el-input v-model="scope.row.Comment"
												  clearable placeholder="Commentaire" size="mini"
												  @change=""
											></el-input>
										</el-col>

										<el-col v-if="Control" :span="2">Qte : {{scope.row.ReferenceQuantity}}</el-col>
									</el-row>
								</template>
							</el-table-column>
		
							<el-table-column v-if="!Control" label="Qte Référence">
								<template slot-scope="scope">
									<el-input-number v-model="scope.row.ReferenceQuantity" 
											controls-position="right" :min="0" size="mini"
											@change=""
									></el-input-number>
								</template>
							</el-table-column>
		
							<el-table-column v-if="Control" label="Qte Contrôle">
								<template slot-scope="scope">
									<div class="header-menu-container on-hover">
										<el-input-number v-model="scope.row.ControledQuantity" 
												controls-position="right" :min="0" size="mini"
												@change=""
										></el-input-number>
										<el-tooltip content="Valider la quantité" placement="bottom" effect="light" open-delay="500">
											<el-button type="success" plain class="icon" icon="fas fa-check fa-fw" size="mini" @click="UpdateControlQuantity(scope.row)"></el-button>
										</el-tooltip>	
									</div>									
								</template>
							</el-table-column>

						</el-table>
					</el-col>
				</el-row>
			</div>
		</el-tab-pane>	

		<!-- ===================================== Event Tab ======================================================= -->
		<el-tab-pane v-if="user.Permissions.Update" label="Evenements" lazy=true style="height: 75vh; padding: 5px 25px; overflow-x: hidden;overflow-y: auto;">
			<!-- Event History -->
			<el-row :gutter="10" align="middle" class="spaced" type="flex">
				<el-col :span="2" class="align-right">Evenements :</el-col>
				<el-col :span="22">
					<el-table 
							:data="current_vehicule.Events"
							height="calc(75vh - 50px)" size="mini" border
					>
						<el-table-column label="" width="68">
							<template slot="header" slot-scope="scope">
								<el-button type="success" plain icon="fas fa-plus fa-fw" size="mini" @click="AddEvent" :disabled="false"></el-button>
							</template>
							<template slot-scope="scope">
								<el-button type="danger" plain icon="far fa-trash-alt fa-fw" size="mini" @click="RemoveEvent(scope.$index)" :disabled="!user.Permissions.Admin"></el-button>
							</template>
						</el-table-column>
	
						<el-table-column label="Début" width="180">
							<template slot-scope="scope">
								<el-date-picker :picker-options="{firstDayOfWeek:1}" 
											placeholder="Date" size="mini" style="width: 100%"
											type="date" format="dd/MM/yyyy" value-format="yyyy-MM-dd"
											v-model="scope.row.StartDate"
											@change="UpdateEvent"
								></el-date-picker>
							</template>
						</el-table-column>
	
						<el-table-column label="Fin" width="180">
							<template slot-scope="scope">
								<el-date-picker :picker-options="{firstDayOfWeek:1}" 
											placeholder="Date" size="mini" style="width: 100%"
											type="date" format="dd/MM/yyyy" value-format="yyyy-MM-dd"
											v-model="scope.row.EndDate"
											@change=""
								></el-date-picker>
							</template>
						</el-table-column>
	
						<el-table-column label="Type" width="150">
							<template slot-scope="scope">
								<el-select v-model="scope.row.Type"
										placeholder="Type" size="mini" filterable
										@change="" style="width: 100%"
								>
									<el-option v-for="item in GetEventTypes()"
											   :key="item.value"
											   :label="item.label"
											   :value="item.value"
									>
									</el-option>
								</el-select>
							</template>
						</el-table-column>
	
						<el-table-column label="Commentaire">
							<template slot-scope="scope">
								<el-input type="textarea" :autosize="{ minRows: 1, maxRows: 5}" placeholder="Commentaire"
										  v-model="scope.row.Comment" clearable size="mini"
								></el-input>
							</template>
						</el-table-column>
					</el-table>
				</el-col>
			</el-row>
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
				
				<el-button :disabled="!hasChanged" type="success" @click="ConfirmChange" plain size="mini">Valider</el-button>
			</el-col>
		</el-row>
	</span>
</el-dialog>`
