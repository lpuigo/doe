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
        <h3>Archive des <a href="/api/worksites/archive">Chantiers</a></h3>
	</div>

	<!-- 
		Body Action Bar
	-->	
	<!--<span slot="footer">-->
	<!--</span>-->
</el-dialog>`
