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
                <i class="fas fa-chart-line icon--left"></i>Productivité des équipes
            </h2>
            <el-radio-group v-if="SiteMode=='Rip'" v-model="GroupMode" size="mini" @change="RefreshStat">
                <el-radio-button label="activity">Par activité</el-radio-button>
                <el-radio-button label="site">Par site</el-radio-button>
                <el-radio-button label="mean">Moyenne</el-radio-button>
            </el-radio-group>
			<el-radio-group v-model="ActiveMode" size="mini" @change="RefreshStat">
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
				<div v-for="(ts, index) in TeamStats" :key="index">
					<h3>{{ts.Team}}</h3>
					<team-productivity-chart :stats="ts"></team-productivity-chart>			
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
                        <div v-for="(ts, index) in RipTeamStats" :key="ts">
                            <h3>{{ts.Team}}</h3>
                            <ripteam-productivity-chart :stats="ts" :colors="SiteColors"></ripteam-productivity-chart>
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
