package teamproductivitymodal

const template string = `
<el-dialog
		:visible.sync="visible" 
		width="90%"
		:before-close="Hide"
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
                <el-radio-button label="activity">Par activité</el-radio-button>
                <el-radio-button label="site">Par site</el-radio-button>
                <el-radio-button label="mean">Moyenne</el-radio-button>
            </el-radio-group>
			<el-radio-group v-model="PeriodMode" size="mini" @change="RefreshStat">
                <el-radio-button v-if="SiteMode!='Orange'" label="day">Jour</el-radio-button>
                <el-radio-button label="week">Hebdo</el-radio-button>
                <el-radio-button label="month">Mensuel</el-radio-button>
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
	<div v-loading="loading" style="height: 65vh;overflow-x: hidden;overflow-y: auto;padding-right: 6px;">
		<div v-if="!loading" style="height: 100%">
			<div v-if="SiteMode == 'Orange'">
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
			<el-container v-else style="height: 100%">
                <el-aside width="200px" style="height: 100%">
                    <div v-for="(val, site) in RipStats.Sites" :key="site" style="margin-top: 8px">
                        <el-checkbox 
                                border size="small" 
                                v-model="SelectedSites[site]" 
                                @change="CheckSitesChange"
                                style="width: 100%"
                        >{{site}}<i class="fas fa-circle icon--right" :style="SiteCircleStyle(site)"></i></el-checkbox>
                    </div>
                </el-aside>
                <el-main style="height: 100%">
                    <div style="height: 100%;overflow-x: hidden;overflow-y: auto;padding-right: 6px;">
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
                    </div>
                </el-main>
			</el-container>
		</div>
	</div>

	<!-- 
		Body Action Bar
	-->	
	<!--<span slot="footer">-->
	<!--</span>-->
</el-dialog>`
