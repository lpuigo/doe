package worksitedetail

const template string = `<div class="worksite-detail">
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
            <el-input placeholder="PA-99999-XXXX" :readonly="readonly" clearable size="mini"
                      v-model="worksite.Ref"
            ></el-input>
        </el-col>
        <el-col :span="6">
            <el-input placeholder="Ville" :readonly="readonly" clearable size="mini"
                      v-model="worksite.City"
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
        Attributes about PMZ 
    -->
    <el-row :gutter="10" type="flex" align="middle">
        <el-col :span="1">
            <span><strong>PMZ:</strong></span>
        </el-col>
        <el-col :span="4">
            <el-input placeholder="PMZ-99999" :readonly="readonly" clearable size="mini"
                      v-model="worksite.Pmz.Ref"
            ></el-input>
        </el-col>
        <el-col :span="4">
            <el-input placeholder="PT-009999" :readonly="readonly" clearable size="mini"
                      v-model="worksite.Pmz.RefPt"
            ></el-input>
        </el-col>
        <el-col :span="15">
            <el-input placeholder="Adresse PMZ" :readonly="readonly" clearable size="mini"
                      v-model="worksite.Pmz.Address"
            >
            </el-input>
        </el-col>
    </el-row>
    <!-- 
        Attributes about PA 
    -->
    <el-row :gutter="10" type="flex" align="middle">
        <el-col :span="1">
            <span><strong>PA:</strong></span>
        </el-col>
        <el-col :span="4">
            <el-input placeholder="PA-99999" :readonly="readonly" clearable size="mini"
                      v-model="worksite.Pa.Ref"
            ></el-input>
        </el-col>
        <el-col :span="4">
            <el-input placeholder="PT-009999" :readonly="readonly" clearable size="mini"
                      v-model="worksite.Pa.RefPt"
            ></el-input>
        </el-col>
        <el-col :span="15">
            <el-input placeholder="Adresse PMZ" :readonly="readonly" clearable size="mini"
                      v-model="worksite.Pa.Address"
            ></el-input>
        </el-col>
    </el-row>
    <!-- 
        Attributes about Orders 
    -->
    <el-row :gutter="10" type="flex" align="middle">
        <el-col :span="2">
            <el-button type="primary" plain icon="fas fa-sitemap icon--left" size="mini" style="width: 100%">Ajouter</el-button>
        </el-col>
        <el-col :span="2">
            <span>Commandes:</span>
        </el-col>
    </el-row>
    <el-row v-for="(o, index) in worksite.Orders" :key="index" :gutter="10">
        <el-col :span="2">
            <el-button type="danger" plain icon="fas fa-sitemap icon--left" size="mini" style="width: 100%">Supprimer</el-button>
        </el-col>
        <el-col :span="22">
            <!-- 
                Attributes about each Order 
            -->
            <el-row :gutter="10">
                <el-col :span="5">
                    <el-input placeholder="F99999jjmmaa" :readonly="readonly" clearable size="mini"
                              v-model="o.Ref"
                    ></el-input>
                </el-col>
                <el-col :span="19">
                    <el-input placeholder="Commentaire sur la commande" :readonly="readonly" clearable size="mini" type="textarea" autosize
                              v-model="o.Comment"
                    ></el-input>
                </el-col>
            </el-row>
            <!-- 
                 Attributes about Troncons 
             -->
            <el-row :gutter="10" type="flex" align="middle">
                <el-col :span="2">
                    <el-button type="primary" plain icon="fas fa-share-alt icon--left" size="mini" style="width: 100%">Ajouter</el-button>
                </el-col>
                <el-col :span="2">
                    <span>Tronçons:</span>
                </el-col>
            </el-row>
            <el-row v-for="(t, index) in o.Troncons" :key="index" :gutter="10">
                <el-col :span="2">
                    <el-button type="danger" plain icon="fas fa-share-alt icon--left" size="mini" style="width: 100%">Supprimer</el-button>
                </el-col>
                <el-col :span="22">
                    <!-- 
                           Attributes about each Troncon 
                    -->
                    <!-- 
                        Attributes about PB 
                    -->
                    <el-row :gutter="10" type="flex" align="middle">
                        <el-col :span="1">
                            <span><strong>PB:</strong></span>
                        </el-col>
                        <el-col :span="3">
                            <el-input placeholder="PB-99999" :readonly="readonly" clearable size="mini"
                                      v-model="t.Pb.Ref"
                            ></el-input>
                        </el-col>
                        <el-col :span="3">
                            <el-input placeholder="PT-009999" :readonly="readonly" clearable size="mini"
                                      v-model="t.Pb.RefPt"
                            ></el-input>
                        </el-col>
                        <el-col :span="7">
                            <el-input placeholder="Adresse PB" :readonly="readonly" clearable size="mini"
                                      v-model="t.Pb.Address"
                            ></el-input>
                        </el-col>
                        <el-col :span="4">
                            <el-switch v-model="t.NeedSignature"
                                       active-color="#ff3200"
                                       active-text="Signature demandée"
                                       inactive-color="#51a825">
                            </el-switch>
                        </el-col>
                    </el-row>
                    <!-- 
                        Attributes about TR 
                    -->
                    <el-row :gutter="10" type="flex" align="middle">
                        <el-col :span="5">
                            <el-input placeholder="TR-99-9999" :readonly="readonly" clearable size="mini"
                                      v-model="t.Ref"
                            ></el-input>
                        </el-col>
                        <el-col :span="3">
                            <el-input-number v-model="t.NbRacco" :min="0" :max="t.NbFiber" :readonly="readonly" size="mini" label="Nb Racco" controls-position="right" style="width: 100%">
                                <template slot="prepend">Nb El</template>
                            </el-input-number>
                        </el-col>
                        <el-col :span="3">
                            <el-input-number v-model="t.NbFiber" :min="6" :step="6" :readonly="readonly" size="mini" label="Nb Fibre" controls-position="right" style="width: 100%">
                                <template slot="prepend">Nb Fibre</template>
                            </el-input-number>
                        </el-col>
                        <el-col :span="4">
                            <el-date-picker :readonly="readonly" format="dd/MM/yyyy" placeholder="Date Installation" size="mini"
                                            style="width: 100%" type="date"
                                            v-model="t.InstallDate"
                                            value-format="yyyy-MM-dd"
                                            :picker-options="{firstDayOfWeek:1, disabledDate(time) { return time.getTime() > Date.now(); }}"
                                            :clearable="false"
                            ></el-date-picker>
                        </el-col>
                        <el-col :span="4">
                            <el-date-picker :readonly="readonly" format="dd/MM/yyyy" placeholder="Date Mesure" size="mini"
                                            style="width: 100%" type="date"
                                            v-model="t.MeasureDate"
                                            value-format="yyyy-MM-dd"
                                            :picker-options="{firstDayOfWeek:1, disabledDate(time) { return time.getTime() > Date.now(); }}"
                                            :clearable="false"
                            ></el-date-picker>
                        </el-col>
                    </el-row>
                </el-col>  
            </el-row>
        </el-col>
    </el-row>
 </div>
`
