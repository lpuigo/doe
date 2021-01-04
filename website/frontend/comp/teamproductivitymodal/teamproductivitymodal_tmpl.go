package teamproductivitymodal

const template string = `
<el-dialog
		:visible.sync="visible" 
		width="90%"
		:before-close="Hide"
		top="8vh"
>
	<!-- 
		Modal Title
	-->
	<span slot="title">
        <div class="header-menu-container">
            <h2 style="margin: 0 0">
                <i class="fas fa-chart-line icon--left"></i>Productivité des équipes {{ModeName()}}
            </h2>
            <el-radio-group v-if="SiteMode=='Orange'" v-model="InfoMode" size="mini" @change="RefreshStat">
                <el-radio-button label="prod">Production</el-radio-button>
                <el-radio-button label="stock">Stock</el-radio-button>
            </el-radio-group>
            <el-radio-group v-if="SiteMode=='Rip'" v-model="GroupMode" size="mini" @change="RefreshStat">
                <el-radio-button v-if="!user.Permissions.Review" label="activity">Par activité</el-radio-button>
                <el-radio-button label="client">Par client</el-radio-button>
                <el-radio-button v-if="!user.Permissions.Review" label="actor">Par acteur</el-radio-button>
                <el-radio-button v-if="!user.Permissions.Review" label="group">Par équipe</el-radio-button>
           </el-radio-group>
            <el-radio-group v-if="SiteMode=='Poles'" v-model="PoleGroupMode" size="mini" @change="RefreshStat">
                <el-radio-button label="client">Par client</el-radio-button>
                <el-radio-button v-if="!user.Permissions.Review" label="actor">Par acteur</el-radio-button>
                <el-radio-button v-if="!user.Permissions.Review" label="group">Par équipe</el-radio-button>
            </el-radio-group>
			<el-radio-group v-model="PeriodMode" size="mini" @change="RefreshStat">
                <el-radio-button v-if="SiteMode!='Orange'" label="day">Jour</el-radio-button>
                <el-radio-button label="week">Hebdo</el-radio-button>
                <el-radio-button label="month">Mensuel</el-radio-button>
                <el-radio-button v-if="!user.Permissions.Review && SiteMode!='Orange'" label="progress">Progression</el-radio-button>
            </el-radio-group>

            <a v-if="SiteMode=='Rip' && user.Permissions.Invoice" :href="GetActorsActivity()"><i class="far fa-file-excel icon--big"></i></a>
            <span v-else></span>                        
            <span></span>                        
        </div>
	</span>

	<!-- 
		Modal Body
		style="height: 100%;"		
	-->
	<div v-loading="loading" style="height: 75vh">
		<div v-if="!loading" style="height: 100%">
			<!--	Productivity Stat for ORANGE -->
			<div v-if="SiteMode == 'Orange'" style="height: 100%; overflow-x: hidden;overflow-y: auto;padding-right: 6px;">
				<div v-for="(ts, index) in GetClientOrangeTeams()" :key="index">
					<h3>{{ts.Team}}</h3>
					<team-productivity-chart :stats="ts"></team-productivity-chart>			
					<div v-if="ts.HasTeams" style="margin-top: 5px; padding-left: 5px ;border-left: 5px solid darkgrey">
						<el-switch v-model="ts.ShowTeams" active-text="Détail des équipes"></el-switch>
						<div v-if="ts.ShowTeams">
							<div v-for="(cts, index) in GetSubOrangeTeams(ts.Team)" :key="cts">
								<h4>{{cts.Team}}</h4>
								<team-productivity-chart :stats="cts"></team-productivity-chart>
							</div>
						</div>
					</div>
				</div>	
			</div>
			<div v-else style="height: 100%">
				<!--	Progress Stat  -->
				<div v-if="PeriodMode == 'progress'" style="height: 100%">
					<el-row type="flex" align="middle" :gutter="10" style="margin-bottom: 5px">
						<el-col :span="4" class="align-right">Progression du mois :</el-col>
						<el-col :span="8">
							<el-date-picker
										v-model="Month" type="month" :clearable="false" size="mini"
										value-format="yyyy-MM-dd" format="dd/MM/yyyy"
										:picker-options="{disabledDate(time) { return time.getTime() > Date.now(); }}"
										placeholder="Choisir un mois"
										@change="RefreshStat"
							></el-date-picker>				
						</el-col>
					</el-row>
					<div style="height: calc(100% - 35px); overflow-x: hidden;overflow-y: auto;padding-right: 6px;">
						<div v-for="(ts, index) in GetClientTeams()" :key="ts">
							<h3>{{ts.Team}}</h3>
							<ripteam-productivity-chart :stats="ts" :colors="SiteColors" heigth="300px" mode="progress"></ripteam-productivity-chart>
							<div v-if="ts.HasTeams" style="margin-top: 5px; padding-left: 5px ;border-left: 5px solid darkgrey">
								<el-switch v-model="ts.ShowTeams" active-text="Détail des acteurs"></el-switch>
								<div v-if="ts.ShowTeams">
									<div v-for="(cts, index) in GetSubTeams(ts.Team)" :key="cts">
										<h4>{{cts.Team}}</h4>
										<ripteam-productivity-chart :stats="cts" :colors="SiteColors" heigth="200px" mode="progress"></ripteam-productivity-chart>
									</div>
								</div>
							</div>
						</div>
					</div>
				</div>
				<!--	Productivity Stat for FOA/RIP/POLE -->
				<el-container v-else style="height: 100%">
					<el-aside width="200px" style="height: 100%">
						<div v-for="(val, site) in RipStats.Sites" :key="site" style="margin-top: 8px">
							<el-checkbox 
									border size="mini" 
									v-model="SelectedSites[site]" 
									@change="CheckSitesChange"
									style="width: 100%"
							>
								<div class="header-menu-container" style="width: 150px">
									<span>{{site}}</span>
									<i class="fas fa-circle icon--right" :style="SiteCircleStyle(site)"></i>
								</div>
							</el-checkbox>
						</div>
					</el-aside>
					<el-main style="height: 100%;overflow-x: hidden;overflow-y: auto;padding-right: 6px;">
						<div v-for="(ts, index) in GetClientTeams()" :key="ts">
							<h3>{{ts.Team}}</h3>
							<ripteam-productivity-chart :stats="ts" :colors="SiteColors" heigth="250px"></ripteam-productivity-chart>
							<div v-if="ts.HasTeams" style="margin-top: 5px; padding-left: 5px ;border-left: 5px solid darkgrey">
								<el-switch v-model="ts.ShowTeams" active-text="Détail des acteurs"></el-switch>
								<div v-if="ts.ShowTeams">
									<div v-for="(cts, index) in GetSubTeams(ts.Team)" :key="cts">
										<h4>{{cts.Team}}</h4>
										<ripteam-productivity-chart :stats="cts" :colors="SiteColors" heigth="180px"></ripteam-productivity-chart>
									</div>
								</div>
							</div>
						</div>
					</el-main>
				</el-container>
			</div>
		</div>
	</div>

	<!-- 
		Body Action Bar
	-->	
	<!--<span slot="footer">-->
	<!--</span>-->
</el-dialog>`
