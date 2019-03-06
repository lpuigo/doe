package modal

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
			Modal simple
		</h2>
	</span>

	<!-- 
		Modal Body
		style="height: 100%;"		
	-->
	<div v-loading="loading" style="height: 65vh;">
		Body
	</div>

	<!-- 
		Body Action Bar
	-->	
	<!--<span slot="footer">-->
	<!--</span>-->
</el-dialog>`
