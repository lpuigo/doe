package invoiceupdate

const template string = `
<div>
    <el-container style="height: 100%">
        <el-header style="height: auto; padding: 5px">
            <el-row :gutter="10" style="margin-bottom: 10px">
                <el-col :span="3">
                    <worksite-status-tag v-model="worksite"></worksite-status-tag>
                </el-col>
                <el-col :offset="1" :span="3" >
                    <span style="float:right; text-align: right">Commentaire dossier:</span>
                </el-col>
                <el-col :span="17">
                    <el-input clearable placeholder="Commentaire sur le dossier" size="mini" type="textarea" autosize
                              v-model.trim="worksite.Comment"
                    ></el-input>                    
                </el-col>		
            </el-row>
        </el-header>
        <el-main  class="invoice-detail" style="height: 100%; padding: 5px">
            <el-row :gutter="10" type="flex" align="middle">
                <el-col :span="4">
                    <span style="float:right; text-align: right">Attachement signé:</span>
                </el-col>
                <el-col :span="4">
                    <el-date-picker format="dd/MM/yyyy" placeholder="Attachement" size="mini"
                                    style="width: 100%" type="date"
                                    v-model="worksite.AttachmentDate"
                                    value-format="yyyy-MM-dd"
                                    :picker-options="{firstDayOfWeek:1, disabledDate(time) { return time.getTime() > Date.now(); }}"
                                    :clearable="false"
                                    :disabled="IsDisabled('AttachmentDate')"
                    ></el-date-picker>
                </el-col>
            </el-row>
            <el-row :gutter="10" type="flex" align="middle">
                <el-col :span="4">
                    <span style="float:right; text-align: right">Emission Facture:</span>
                </el-col>
                <el-col :span="4">
                    <el-date-picker format="dd/MM/yyyy" placeholder="Facturation" size="mini"
                                    style="width: 100%" type="date"
                                    v-model="worksite.InvoiceDate"
                                    value-format="yyyy-MM-dd"
                                    :picker-options="{firstDayOfWeek:1, disabledDate(time) { return time.getTime() > Date.now(); }}"
                                    :clearable="false"
                                    :disabled="IsDisabled('InvoiceDate')"
                    ></el-date-picker>
                </el-col>
                <el-col :span="10">
                    <el-input clearable placeholder="Nom de la facture" size="mini"
                              v-model.trim="worksite.InvoiceName"
                              :disabled="IsDisabled('InvoiceName')"
                    ></el-input>
                </el-col>
            </el-row>
            <el-row :gutter="10" type="flex" align="middle">
                <el-col :span="4">
                    <span style="float:right; text-align: right">Paiement:</span>
                </el-col>
                <el-col :span="4">
                    <el-date-picker format="dd/MM/yyyy" placeholder="Paiement" size="mini"
                                    style="width: 100%" type="date"
                                    v-model="worksite.PaymentDate"
                                    value-format="yyyy-MM-dd"
                                    :picker-options="{firstDayOfWeek:1, disabledDate(time) { return time.getTime() > Date.now(); }}"
                                    :clearable="false"
                                    :disabled="IsDisabled('PaymentDate')"
                    ></el-date-picker>
                </el-col>
            </el-row>
        </el-main>
    </el-container>
</div>
`
