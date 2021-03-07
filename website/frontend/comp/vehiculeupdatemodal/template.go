package vehiculeupdatemodal

const template string = `<el-dialog
        :before-close="HideWithControl"
        :visible.sync="visible" :close-on-click-modal="false"
        width="70%" top="5vh"
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
								<el-select placeholder="Acteur" size="mini"
										   v-model="scope.row.ActorId"
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
