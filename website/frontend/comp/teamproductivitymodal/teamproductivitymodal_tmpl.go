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
		<h2 style="margin: 0 0">
			<i class="fas fa-chart-line icon--left"></i>Productivité des équipes
		</h2>
	</span>

	<!-- 
		Modal Body
		style="height: 100%;"		
	-->
	<div v-loading="loading" style="height: 65vh;overflow-x: hidden;overflow-y: auto;padding-right: 6px;">
		<pre>{{Stats}}</pre>
	</div>

	<!-- 
		Body Action Bar
	-->	
	<!--<span slot="footer">-->
	<!--</span>-->
</el-dialog>`
