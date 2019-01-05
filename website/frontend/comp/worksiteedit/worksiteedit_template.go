package worksiteedit

const template string = `
<div class="worksite-detail">
    <el-row :gutter="10" type="flex" align="middle">
        <el-col :span="3">
            <el-input placeholder="Client" :readonly="readonly" clearable size="mini"
                      v-model="worksite.Client"
            >
                <template slot="prepend">Client:</template>
            </el-input>
        </el-col>
        <el-col :span="4">
            <el-input placeholder="Ville" :readonly="readonly" clearable size="mini"
                      v-model="worksite.City"
            >
                <template slot="prepend">Ville:</template>
            </el-input>
        </el-col>
        <el-col :span="3">
            <el-input placeholder="PA-99999-XXXX" 
                      v-model="worksite.Ref"
                      :readonly="readonly" clearable size="mini"
            >
                <template slot="prepend">Chantier:</template>
            </el-input>
        </el-col>
        <el-col :span="3" style="text-align: right">
            <el-tag :type="StatusType" size="medium" style="width: 100%;text-align: left">{{worksite.Status | FormatStatus}}</el-tag>
        </el-col>
        <el-col :span="3">
            <el-date-picker :readonly="readonly" format="dd/MM/yyyy" placeholder="Soumission" size="mini"
                            style="width: 100%" type="date"
                            v-model="worksite.OrderDate"
                            value-format="yyyy-MM-dd"
                            :picker-options="{firstDayOfWeek:1, disabledDate(time) { return time.getTime() > Date.now(); }}"
                            :clearable="false"
            ></el-date-picker>
        </el-col>
        <el-col :span="3">
            <el-date-picker :readonly="readonly" format="dd/MM/yyyy" placeholder="Envoi DOE" size="mini"
                            style="width: 100%" type="date"
                            v-model="worksite.DoeDate"
                            value-format="yyyy-MM-dd"
                            :picker-options="{firstDayOfWeek:1, disabledDate(time) { return time.getTime() > Date.now(); }}"
                            :clearable="false"
                            :disabled="worksite.Status < '30 DOE'"
            ></el-date-picker>
        </el-col>
        <el-col :span="3">
            <el-date-picker :readonly="readonly" format="dd/MM/yyyy" placeholder="Attachement" size="mini"
                            style="width: 100%" type="date"
                            v-model="worksite.AttachmentDate"
                            value-format="yyyy-MM-dd"
                            :picker-options="{firstDayOfWeek:1, disabledDate(time) { return time.getTime() > Date.now(); }}"
                            :clearable="false"
                            :disabled="worksite.Status < '40 Attachment'"
            ></el-date-picker>
        </el-col>
        <el-col :span="3">
            <el-date-picker :readonly="readonly" format="dd/MM/yyyy" placeholder="Payement" size="mini"
                            style="width: 100%" type="date"
                            v-model="worksite.PaymentDate"
                            value-format="yyyy-MM-dd"
                            :picker-options="{firstDayOfWeek:1, disabledDate(time) { return time.getTime() > Date.now(); }}"
                            :clearable="false"
                            :disabled="worksite.Status < '50 Payment'"
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
	<hr>
    <el-row :gutter="10" type="flex" align="middle">
        <el-col :span="2">
            <el-button type="primary" plain
					   icon="fas fa-sitemap icon--left" 
					   size="mini" style="width: 100%"
					   @click="AddOrder()"
			>Ajouter</el-button>
        </el-col>
    </el-row>
</div>
`
