package adminmodal

const template string = `
<el-dialog
		:visible.sync="visible" 
		width="60%"
		:before-close="Hide"
>
	<!-- 
		Modal Title
	-->
	<span slot="title">
		<h2 style="margin: 0 0">
			<i class="fas fa-wrench icon--left"></i>Administration
		</h2>
	</span>

	<!-- 
		Modal Body
		style="height: 100%;"		
	-->
	<div v-loading="loading" style="height: 65vh;overflow-x: hidden;overflow-y: auto;padding: 6px 20px;">
        <h3>Archive des <a href="/api/worksites/archive">Chantiers Orange</a></h3>
        <h3>Archive des <a href="/api/ripsites/archive">Chantiers RIP</a></h3>
        <h3>Archive des <a href="/api/polesites/archive">Chantiers Poteaux</a></h3>
		<h4>&nbsp;</h4>
        <el-button type="primary" @click="ReloadData">Rechargement des données</el-button>
	</div>

	<!-- 
		Body Action Bar
	-->	
	<!--<span slot="footer">-->
	<!--</span>-->
</el-dialog>`
