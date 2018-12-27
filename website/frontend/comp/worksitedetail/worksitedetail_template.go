package worksitedetail

const template string = `
<div class="worksite-detail">
    <el-row :gutter="10">
        <el-col :span="6">
            <el-input v-if="readonly" placeholder="Statut" :readonly="true" clearable size="mini"
                      v-model="worksite.Ref"
            ></el-input>            
            <el-select v-else placeholder="Statut" size="mini" style="width: 100%"
                       v-model="worksite.Status">
                <el-option label="Nouveau" value="New"></el-option>
                <el-option label="En cours" value="InProgress"></el-option>
                <el-option label="Terminé" value="Done"></el-option>
                <el-option label="A Reprendre" value="Rework"></el-option>
            </el-select>
        </el-col>
		<el-col :span="6">
			<el-input placeholder="Ville" :readonly="readonly" clearable size="mini"
					  v-model="worksite.City"
			></el-input>
		</el-col>        
		<el-col :span="6">
            <el-input placeholder="PA-99999-XXXX" :readonly="readonly" clearable size="mini"
                      v-model="worksite.Ref"
            ></el-input>
        </el-col>

        <el-col :span="6">
            <el-date-picker :readonly="readonly" format="dd/MM/yyyy" placeholder="Soumission" size="mini"
                            style="width: 100%" type="date"
                            v-model="worksite.OrderDate"
                            value-format="yyyy-MM-dd"
                            :picker-options="{firstDayOfWeek:1, disabledDate(time) { return time.getTime() > Date.now(); }}"
                            :clearable="false"
            ></el-date-picker>
        </el-col>
    </el-row>
    <el-row :gutter="10">
        <el-col :span="24">
            <el-input :readonly="readonly" clearable placeholder="Commentaire sur le dossier" size="mini" type="textarea" autosize
                      v-model="worksite.Comment"
            ></el-input>
        </el-col>
    </el-row>
    <!-- 
        Attributes about PMZ & PA
    -->
	<el-row :gutter="10">
		<el-col :offset="2" :span="22">
			<pt-edit title="PMZ" v-model="worksite.Pmz" :readonly="readonly"></pt-edit>
		</el-col>
	</el-row>
	<el-row :gutter="10">
		<el-col :offset="2" :span="22">
			<pt-edit title="PA" v-model="worksite.Pa" :readonly="readonly"></pt-edit>
		</el-col>
	</el-row>
    <!-- 
        Attributes about Orders
    -->
    <el-row :gutter="10" type="flex" align="middle">
        <el-col :span="2">
            <el-button type="primary" plain
					   icon="fas fa-sitemap icon--left" 
					   size="mini" style="width: 100%"
					   @click="AddOrder()"
			>Ajouter</el-button>
        </el-col>
        <el-col :span="2">
            <span>Commandes:</span>
        </el-col>
    </el-row>
	<div v-for="(order, index) in worksite.Orders" :key="index">
		<hr>
		<el-row :gutter="10">
			<el-col :span="2">
				<el-button type="danger" plain
						   icon="fas fa-sitemap icon--left"
						   size="mini" style="width: 100%"
						   :disabled="worksite.Orders.length<=1"
						   @click="DeleteOrder(index)"
				>Supprimer</el-button>
			</el-col>
			<el-col :span="22">
				<!-- 
					Attributes about each Order 
				-->
				<order-edit v-model="order" :readonly="readonly"></order-edit>
			</el-col>
		</el-row>
	</div>
</div>
`
