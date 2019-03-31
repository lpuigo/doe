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
            <el-radio-group v-model="ActiveMode" size="mini" @change="RefreshStat">
                <el-radio-button label="week">Hebdo</el-radio-button>
                <el-radio-button label="month">Mensuel</el-radio-button>
            </el-radio-group>
            <span></span>                        
            <span></span>                        
        </div>
	</span>

	<!-- 
		Modal Body
		style="height: 100%;"		
	-->
	<div v-loading="loading" style="height: 65vh;overflow-x: hidden;overflow-y: auto;padding-right: 6px;">
		<div v-if="!loading" v-for="(ts, index) in TeamStats" :key="index"
		>
			<h3>Equipe : {{ts.Team}}</h3>
			<team-productivity-chart :stats="ts"></team-productivity-chart>			
		</div>	
	</div>

	<!-- 
		Body Action Bar
	-->	
	<!--<span slot="footer">-->
	<!--</span>-->
</el-dialog>`
