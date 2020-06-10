package poleedit

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/polemap"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/model/polesite"
	"github.com/lpuig/ewin/doe/website/frontend/model/polesite/poleconst"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/autocomplete"
	"github.com/lpuig/ewin/doe/website/frontend/tools/date"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements/message"
	"github.com/lpuig/ewin/doe/website/frontend/tools/latlong"
	"github.com/lpuig/ewin/doe/website/frontend/tools/leaflet"
	"sort"
	"strconv"
	"strings"
)

const template string = `<div>
	<div class="header-menu-container spaced">
		<h1>Appui: <span style="font-size: 1.5em">{{editedpolemarker.Pole.Ref}} {{editedpolemarker.Pole.Sticker}}</span></h1>
        
        <el-tooltip content="Zoomer sur la carte" placement="bottom" effect="light" open-delay="500">
            <el-button class="icon" icon="fas fa-crosshairs icon--big" @click="CenterOnEdited" size="mini"></el-button>
        </el-tooltip>
	</div>
	<!--	Tools -->
	<div class="header-menu-container" style="margin-bottom: 5px">
		<div></div>
		
		<div>
			<el-popover placement="bottom" width="360" title="Suppression de l'appui"
						v-model="VisibleDeletePole">
				<p>Confirmez la suppression de l'appui {{editedpolemarker.Pole.Ref}} {{editedpolemarker.Pole.Sticker}}?</p>
				<div style="text-align: right; margin: 0; margin-top: 10px">
					<el-button size="mini" @click="VisibleDeletePole = false">Annuler</el-button>
					<el-button type="danger" plain size="mini" @click="DeletePole">Confirmer</el-button>
				</div>
				<el-tooltip slot="reference" content="Supprimer cet appui" placement="bottom" effect="light" open-delay="500">
					<el-button type="danger" plain class="icon" icon="far fa-trash-alt icon--big" size="mini"></el-button>
				</el-tooltip>
			</el-popover>
			
			<el-popover placement="bottom" width="360" title="Duplication de l'appui"
						v-model="VisibleDuplicatePole">
				<el-row :gutter="5" type="flex" align="middle" class="spaced">
					<el-col :span="6">
						<el-checkbox v-model="DuplicateContext.DoIncrement">Incrément</el-checkbox>
					</el-col>
					<el-col v-if="DuplicateContext.DoIncrement" :span="18">
						<el-input-number placeholder="Incrément"
										 v-model="DuplicateContext.Increment" controls-position="right" :min="-20" :max="20" clearable size="mini"
						></el-input-number>
					</el-col>
				</el-row>
				<p>Confirmez la duplication de l'appui {{editedpolemarker.Pole.Ref}} {{editedpolemarker.Pole.Sticker}}?</p>
				<div style="text-align: right; margin: 0; margin-top: 10px">
					<el-button size="mini" @click="VisibleDuplicatePole = false">Annuler</el-button>
					<el-button type="warning" plain size="mini" @click="DuplicatePole">Confirmer</el-button>
				</div>
				<el-tooltip slot="reference" content="Dupliquer cet appui" placement="bottom" effect="light" open-delay="500">
					<el-button type="warning" plain class="icon" icon="far fa-clone icon--big" size="mini"></el-button>
				</el-tooltip>	
			</el-popover>
		</div>
	</div>

    <!--	Référence -->
    <el-row :gutter="5" type="flex" align="middle" class="spaced">
        <el-col :span="6" class="align-right">Référence:</el-col>
        <el-col :span="15">
            <el-input placeholder="Référence"
                      v-model="editedpolemarker.Pole.Ref" clearable size="mini"
            ></el-input>
        </el-col>
		<el-col :span="3">
			<el-tooltip content="Selectionner les autres appuis du groupe" placement="bottom" effect="light" open-delay="500">
				<el-button type="info" plain style="width: 100%" class="icon" icon="fas fa-eye-dropper" size="mini" @click="SelectByRef"></el-button>
			</el-tooltip>	
		</el-col>
    </el-row>
    
    <!-- Etiquette -->
    <el-row :gutter="5" type="flex" align="middle" class="spaced">
        <el-col :span="6" class="align-right">Appui (Etiquette):</el-col>
        <el-col :span="18">
            <el-input placeholder="Etiquette"
                      v-model="editedpolemarker.Pole.Sticker" clearable size="mini"
            ></el-input>
        </el-col>
    </el-row>

    <!-- Commentaire -->
    <el-row :gutter="5" type="flex" align="middle" class="spaced">
        <el-col :span="6" class="align-right">Commentaire:</el-col>
        <el-col :span="18">
            <el-input type="textarea" :autosize="{ minRows: 1, maxRows: 2}" placeholder="Commentaire"
                      v-model="editedpolemarker.Pole.Comment" clearable size="mini"
            ></el-input>
        </el-col>
    </el-row>

    <el-collapse v-model="chapters" @change="ChapterChange">
		<!-- ======================== Chapter Address ===================================-->
        <el-collapse-item name="address">
            <template slot="title">
                <h2 class="title">Adresse: <a v-if="editedlatlong != ''" :href="GetGMAPUrl(editedlatlong)" rel="noopener noreferrer" target="_blank">(GMap)</a><span class="blue"> {{editedpolemarker.Pole.Address}}</span></h1>
            </template>
            <!-- Déplacable -->
            <el-row :gutter="5" type="flex" align="middle" class="spaced">
                <el-col :offset="6" :span="18">
                    <el-checkbox v-model="editedpolemarker.Draggable" @change="UpdatePoleDrag()">Marqueur déplaçable</el-checkbox>
                </el-col>
            </el-row>
        
            <!-- Lat / Long -->
            <el-row :gutter="5" type="flex" align="middle" class="spaced">
                <el-col :span="6" class="align-right">Lat / Long:</el-col>
                <el-col :span="15">
                    <el-input v-model="editedlatlong" size="mini" @input="UpdatePoleLatLong"></el-input>
                </el-col>
                <el-col :span="3">
                    <el-popover title="Latitude / Longitude" placement="right"
                                trigger="click" width="300" v-model="VisibleLatLong">
                        <el-row :gutter="5" type="flex" align="middle" class="spaced">
                            <el-col :span="4" class="align-right">Lat:</el-col>
                            <el-col :span="20">
                                <el-input placeholder="Latitude"
                                          v-model="EditedLat" size="mini"
                                          @change="UpdatePoleDegLat"
                                ></el-input>
                            </el-col>
                        </el-row>
                        <el-row :gutter="5" type="flex" align="middle" class="spaced">
                            <el-col :span="4" class="align-right">Long:</el-col>
                            <el-col :span="20">
                                <el-input placeholder="Longitude"
                                          v-model="EditedLong" size="mini"
                                          @change="UpdatePoleDegLong"
                                ></el-input>
                            </el-col>
                        </el-row>
                        <el-button slot="reference" type="info" plain style="width: 100%" class="icon" icon="far fa-edit" size="mini" :disabled="VisibleLatLong"></el-button>
                    </el-popover>
                </el-col>
            </el-row>
            
            <!-- Ville -->
            <el-row :gutter="5" type="flex" align="middle" class="spaced">
                <el-col :span="6" class="align-right">Ville:</el-col>
                <el-col :span="18">
					<el-autocomplete placeholder="Ville"
							v-model="editedpolemarker.Pole.City" clearable size="mini"
							:fetch-suggestions="GetCities" :trigger-on-focus="false" 
							@change="UpdatePoleCity"	
					></el-autocomplete>
                </el-col>
            </el-row>
        
            <!-- Adresse -->
            <el-row :gutter="5" type="flex" align="middle">
                <el-col :span="6" class="align-right">Adresse:</el-col>
                <el-col :span="18">
                    <el-input placeholder="Adresse"
                              v-model="editedpolemarker.Pole.Address" clearable size="mini"
                    ></el-input>
                </el-col>
            </el-row>
        </el-collapse-item>

		<!-- ======================== Chapter Work to do ===================================-->
        <el-collapse-item name="work">
            <template slot="title">
                <h2 class="title">Travaux: <span class="blue">{{editedpolemarker.Pole.Material}} {{editedpolemarker.Pole.Height}}m</span></h1>
            </template>
            <!-- Matériau & Hauteur-->
            <el-row :gutter="5" type="flex" align="middle" class="spaced">
                <el-col :span="6" class="align-right">Matériau:</el-col>
                <el-col :span="10">
                    <el-select v-model="editedpolemarker.Pole.Material" filterable size="mini" style="width: 100%"
                               @clear=""
                    >
                        <el-option
                                v-for="item in GetMaterials()"
                                :key="item.value"
                                :label="item.label"
                                :value="item.value"
                        ></el-option>
                    </el-select>
                </el-col>
                <el-col :span="8">
                    <el-input-number v-model="editedpolemarker.Pole.Height" size="mini" controls-position="right" :precision="0" :min="6" :max="15" style="width: 100%"
                    ></el-input-number>
                </el-col>
            </el-row>
        
            <!-- Produits -->
            <el-row :gutter="5" type="flex" align="middle">
                <el-col :span="6" class="align-right">Produits:</el-col>
                <el-col :span="18">
                    <el-select v-model="editedpolemarker.Pole.Product" multiple placeholder="Produits" size="mini" style="width: 100%"
                               @clear=""
                               @change="UpdateProduct"
                    >
                        <el-option
                                v-for="item in GetProducts()"
                                :key="item.value"
                                :label="item.label"
                                :value="item.value"
                        ></el-option>
                    </el-select>
                </el-col>
            </el-row>
        </el-collapse-item>

		<!-- ======================== Chapter DICT DT DA ===================================-->
        <el-collapse-item name="dict">
            <template slot="title">
                <h2 class="title">DT, DICT et DA&nbsp;<a v-if="editedpolemarker.Pole.DictRef != ''" :href="GetDICTUrl()" rel="noopener noreferrer" target="_blank">(Carte DICT.fr)</a></h1>
            </template>
            <!-- DT -->
            <el-row :gutter="5" type="flex" align="middle" class="spaced">
				<el-col :span="3">
					<el-tooltip content="Appliquer les infos DT / DICT à la selection" placement="bottom" effect="light" open-delay="500">
						<el-button type="warning" plain style="width: 100%" class="icon" icon="fas fa-clone" size="mini" @click="ApplyDict"></el-button>
					</el-tooltip>	
				</el-col>
                <el-col :span="3" class="align-right">DT:</el-col>
                <el-col :span="10">
                    <el-input placeholder="Référence DT"
                              v-model.trim="editedpolemarker.Pole.DtRef" clearable size="mini"
                    ></el-input>
                </el-col>
				<el-col :offset="5" :span="3">
					<el-tooltip content="Selectionner les autres appuis ayant la même DT / DICT" placement="bottom" effect="light" open-delay="500">
						<el-button type="info" plain style="width: 100%" class="icon" icon="fas fa-eye-dropper" size="mini" @click="SelectByDict"></el-button>
					</el-tooltip>	
				</el-col>
            </el-row>
            <!-- DICT -->
            <el-row :gutter="5" type="flex" align="middle" class="spaced">
                <el-col :span="6" class="align-right">DICT:</el-col>
                <el-col :span="10">
                    <el-input placeholder="Référence DICT"
                              v-model.trim="editedpolemarker.Pole.DictRef" clearable size="mini"
							  @change="CheckPermissions()"
                    ></el-input>
                </el-col>
                <el-col :span="8">
                    <el-date-picker format="dd/MM/yyyy" placeholder="Début" size="mini"
                                    style="width: 100%" type="date"
                                    v-model="editedpolemarker.Pole.DictDate"
                                    value-format="yyyy-MM-dd"
                                    :picker-options="{firstDayOfWeek:1}"
									@change="CheckPermissions()"
                    ></el-date-picker>
                </el-col>
            </el-row>
            <el-row :gutter="5" type="flex" align="middle" class="spaced">
                <el-col :span="6" class="align-right">Info DICT:</el-col>
                <el-col :span="18">
                    <el-input type="textarea" :autosize="{ minRows: 1, maxRows: 2}" placeholder="Information DICT"
                              v-model="editedpolemarker.Pole.DictInfo" clearable size="mini"
                    ></el-input>
                </el-col>
            </el-row>
            <!-- DA -->
            <el-row :gutter="5" type="flex" align="middle" class="spaced">
				<el-col :span="3">
					<el-tooltip content="Appliquer les infos DA à la selection" placement="bottom" effect="light" open-delay="500">
						<el-button type="warning" plain style="width: 100%" class="icon" icon="fas fa-clone" size="mini" @click="ApplyDa"></el-button>
					</el-tooltip>	
				</el-col>
                <el-col :span="3" class="align-right">DA:</el-col>
                <el-col :span="9">
                    <el-date-picker format="dd/MM/yyyy" placeholder="Demande" size="mini"
                                    style="width: 100%" type="date"
                                    v-model="editedpolemarker.Pole.DaQueryDate"
                                    value-format="yyyy-MM-dd"
                                    :picker-options="{firstDayOfWeek:1}"
									@change="CheckPermissions()"
                    ></el-date-picker>
                </el-col>
            </el-row>
            <el-row :gutter="5" type="flex" align="middle">
                <el-col :offset="6" :span="9">
                    <el-date-picker format="dd/MM/yyyy" placeholder="Début" size="mini"
                                    style="width: 100%" type="date"
                                    v-model="editedpolemarker.Pole.DaStartDate"
                                    value-format="yyyy-MM-dd"
                                    :picker-options="{firstDayOfWeek:1}"
									@change="CheckPermissions()"
                    ></el-date-picker>
                </el-col>
                <el-col :span="9">
                    <el-date-picker format="dd/MM/yyyy" placeholder="Fin" size="mini"
                                    style="width: 100%" type="date"
                                    v-model="editedpolemarker.Pole.DaEndDate"
                                    value-format="yyyy-MM-dd"
                                    :picker-options="{firstDayOfWeek:1}"
									@change="CheckPermissions()"
                    ></el-date-picker>
                </el-col>
            </el-row>
        </el-collapse-item>

		<!-- ======================== Chapter State ===================================-->
        <el-collapse-item name="state">
            <template slot="title">
                <h2 class="title">Etat: <span class="blue">{{FormatState(editedpolemarker.Pole)}}</span></h1>
            </template>
            <!-- Status -->
            <el-row :gutter="5" type="flex" align="middle" class="spaced">
                <el-col :span="6" class="align-right">Status:</el-col>
                <el-col :span="18">
                    <el-select v-model="editedpolemarker.Pole.State" filterable size="mini" style="width: 100%"
                               @clear=""
                               @change="UpdateState()"
                    >
                        <el-option
                                v-for="item in GetStates()"
                                :key="item.value"
                                :label="item.label"
                                :value="item.value"
								:disabled="item.disabled"
                        ></el-option>
                    </el-select>
                </el-col>
            </el-row>

            <!-- Actors -->
            <el-row :gutter="5" type="flex" align="middle" class="spaced">
                <el-col :span="6" class="align-right">Acteurs:</el-col>
                <el-col :span="18">
                    <el-select v-model="editedpolemarker.Pole.Actors" multiple placeholder="Acteurs" size="mini" style="width: 100%"
                               @clear=""
                               @change="UpdateActors"
                    >
                        <el-option
                                v-for="item in GetActors()"
                                :key="item.value"
                                :label="item.label"
                                :value="item.value"
                        >
                            <!--
                            <span style="float: left">{{ item.value }}</span>
                            <span style="float: right; color: #8492a6; font-size: 0.9em; margin-right: 15px">{{ item.label }}</span>
                            -->
                        </el-option>
                    </el-select>
                </el-col>
            </el-row>

            <!-- Date Aspiratrice -->
            <el-row :gutter="5" type="flex" align="middle" class="spaced">
                <el-col :span="6" class="align-right">Aspiratrice:</el-col>
                <el-col :span="18">
                    <el-date-picker format="dd/MM/yyyy" placeholder="Date" size="mini"
                                    style="width: 100%" type="date"
                                    v-model="editedpolemarker.Pole.AspiDate"
                                    value-format="yyyy-MM-dd"
                    ></el-date-picker>
                </el-col>
            </el-row>
        
            <!-- Kizeo -->
            <el-row :gutter="5" type="flex" align="middle" class="spaced">
                <el-col :span="6" class="align-right">Kizeo:</el-col>
                <el-col :span="18">
                    <el-input placeholder="Référence Kizeo"
                              v-model="editedpolemarker.Pole.Kizeo" clearable size="mini"
                    ></el-input>
                </el-col>
            </el-row>
        
            <!-- Date -->
            <el-row v-if="ShowDate" :gutter="5" type="flex" align="middle" class="spaced">
                <el-col :span="6" class="align-right">Réalisation:</el-col>
                <el-col :span="18">
                    <el-date-picker format="dd/MM/yyyy" placeholder="Date" size="mini"
                                    style="width: 100%" type="date"
                                    v-model="editedpolemarker.Pole.Date"
                                    value-format="yyyy-MM-dd"
                                    :picker-options="{firstDayOfWeek:1, disabledDate(time) { return time.getTime() > Date.now(); }}"
                    ></el-date-picker>
                </el-col>
            </el-row>
        
            <!-- AttachmentDate -->
            <el-row v-if="ShowAttachmentDate" :gutter="5" type="flex" align="middle">
                <el-col :span="6" class="align-right">Attachement:</el-col>
                <el-col :span="18">
                    <el-date-picker format="dd/MM/yyyy" placeholder="Date" size="mini"
                                    style="width: 100%" type="date"
                                    v-model="editedpolemarker.Pole.AttachmentDate"
                                    value-format="yyyy-MM-dd"
                                    :picker-options="{firstDayOfWeek:1, disabledDate(time) { return time.getTime() > Date.now(); }}"
                    ></el-date-picker>
                </el-col>
            </el-row>
        </el-collapse-item>
        
    </el-collapse>
</div>
`

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Registration

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("pole-edit", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.Props("editedpolemarker", "user", "polesite", "chapters"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewPoleEditModel(vm)
		}),
		//hvue.Mounted(func(vm *hvue.VM) {
		//	pem := PoleEditModelFromJS(vm.Object)
		//}),
		hvue.MethodsOf(&PoleEditModel{}),
		hvue.Computed(
			"editedlatlong",
			func(vm *hvue.VM) interface{} {
				pem := PoleEditModelFromJS(vm.Object)
				pem.setDegLatLong()
				return strconv.FormatFloat(pem.EditedPoleMarker.Pole.Lat, 'f', 8, 64) +
					", " +
					strconv.FormatFloat(pem.EditedPoleMarker.Pole.Long, 'f', 8, 64)
			}),
		hvue.Computed(
			"ShowDate",
			func(vm *hvue.VM) interface{} {
				pem := PoleEditModelFromJS(vm.Object)
				return pem.EditedPoleMarker.Pole.State == poleconst.StateDone || pem.EditedPoleMarker.Pole.State == poleconst.StateAttachment
			}),
		hvue.Computed(
			"ShowAttachmentDate",
			func(vm *hvue.VM) interface{} {
				pem := PoleEditModelFromJS(vm.Object)
				if !pem.User.HasPermissionInvoice() {
					return false
				}
				stateShow := pem.EditedPoleMarker.Pole.State == poleconst.StateDone || pem.EditedPoleMarker.Pole.State == poleconst.StateAttachment
				return stateShow && !tools.Empty(pem.EditedPoleMarker.Pole.Date)
			}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type PoleEditModel struct {
	*js.Object
	EditedPoleMarker     *polemap.PoleMarker `js:"editedpolemarker"`
	VisibleDeletePole    bool                `js:"VisibleDeletePole"`
	VisibleLatLong       bool                `js:"VisibleLatLong"`
	VisibleDuplicatePole bool                `js:"VisibleDuplicatePole"`

	DuplicateContext *DuplicateContext `js:"DuplicateContext"`

	ActiveChapter []string `js:"chapters"`

	User     *fm.User           `js:"user"`
	Polesite *polesite.Polesite `js:"polesite"`

	EditedLat  string `js:"EditedLat"`
	EditedLong string `js:"EditedLong"`

	VM *hvue.VM `js:"VM"`
}

func PoleEditModelFromJS(obj *js.Object) *PoleEditModel {
	return &PoleEditModel{Object: obj}
}

func NewPoleEditModel(vm *hvue.VM) *PoleEditModel {
	pem := &PoleEditModel{Object: tools.O()}
	pem.VisibleDeletePole = false
	pem.VisibleLatLong = false
	pem.VisibleDuplicatePole = false

	pem.DuplicateContext = NewDuplicateContext()

	pem.ActiveChapter = []string{}

	pem.User = fm.NewUser()
	pem.Polesite = nil

	pem.EditedLat = ""
	pem.EditedLong = ""

	pem.VM = vm
	pem.EditedPoleMarker = polemap.DefaultPoleMarker()
	return pem
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Methods

func (pem *PoleEditModel) GetActors(vm *hvue.VM) []*elements.ValueLabel {
	pem = PoleEditModelFromJS(vm.Object)
	client := pem.User.GetClientByName(pem.Polesite.Client)
	if client == nil {
		return nil
	}
	res := []*elements.ValueLabel{}
	for _, actor := range client.Actors {
		if !actor.Active {
			// skip inactive actors
			continue
		}
		ref := actor.GetRef()
		res = append(res, elements.NewValueLabel(strconv.Itoa(actor.Id), ref))
	}
	return res
}

func (pem *PoleEditModel) UpdateActors(vm *hvue.VM) {
	pem = PoleEditModelFromJS(vm.Object)
	client := pem.User.GetClientByName(pem.Polesite.Client)
	if client == nil {
		return
	}
	actors := make(map[string]string)
	for _, actor := range client.Actors {
		actors[strconv.Itoa(actor.Id)] = actor.GetRef()
	}
	pem.EditedPoleMarker.Pole.Get("Actors").Call("sort", func(a, b string) int {
		// check if actors are not known
		if actors[a] == "" && actors[b] == "" {
			return 0
		}
		if !(actors[a] != "" && actors[b] != "") {
			return 1
		}
		// compare known actors
		if actors[a] < actors[b] {
			return -1
		}
		return 1
	})
	pem.EditedPoleMarker.Pole.Get("Product").Call("sort")
}

func (pem *PoleEditModel) FormatDate(d string) string {
	if d == "" {
		return ""
	}
	return date.DateString(d)
}

func (pem *PoleEditModel) FormatState(p *polesite.Pole) string {
	res := polesite.PoleStateLabel(p.State)
	switch p.State {
	case poleconst.StateAttachment:
		return res + " (Fait le " + date.DateString(p.Date) + ")"
	case poleconst.StateDone:
		return res + " le " + date.DateString(p.Date)
	default:
		return res
	}
}

func (pem *PoleEditModel) GetStates(vm *hvue.VM) []*elements.ValueLabelDisabled {
	pem = PoleEditModelFromJS(vm.Object)
	return polesite.GetStatesValueLabel(pem.User.HasPermissionInvoice())
}

func (pem *PoleEditModel) UpdateProduct(vm *hvue.VM) {
	pem = PoleEditModelFromJS(vm.Object)
	ep := pem.EditedPoleMarker.Pole
	ep.CheckProductConsistency()
	ep.Get("Product").Call("sort")
	if ep.IsInStateToBeChecked() {
		ep.SetState(poleconst.StateToDo)
	}
	pem.UpdateState(vm)
}

// CheckPermissions updates EditedPoleMarker according to DICT and DA Dates
func (pem *PoleEditModel) CheckPermissions(vm *hvue.VM) {
	pem = PoleEditModelFromJS(vm.Object)
	ep := pem.EditedPoleMarker
	currentState := ep.Pole.State
	ep.Pole.UpdateState()
	if ep.Pole.State != currentState {
		pem.UpdateState(vm)
	}
}

// UpdateState
func (pem *PoleEditModel) UpdateState(vm *hvue.VM) {
	pem = PoleEditModelFromJS(vm.Object)
	ep := pem.EditedPoleMarker
	//ep.Pole.SetState(ep.Pole.State)
	ep.Pole.UpdateState()
	ep.UpdateFromState()
	ep.Map.RefreshPoleMarkersGroups()
}

func (pem *PoleEditModel) GetMaterials() []*elements.ValueLabel {
	return polesite.GetMaterialsValueLabel()
}

func (pem *PoleEditModel) GetProducts() []*elements.ValueLabel {
	return polesite.GetProductsValueLabel()
}

func (pem *PoleEditModel) UpdateTooltip(vm *hvue.VM) {
	pem = PoleEditModelFromJS(vm.Object)
	pem.EditedPoleMarker.UpdateTitle()
	pem.EditedPoleMarker.Refresh()
}

func (pem *PoleEditModel) UpdatePoleDegLat(vm *hvue.VM, value string) {
	pem = PoleEditModelFromJS(vm.Object)
	val, err := latlong.DegToDec(value)
	if err != nil {
		pem.errMsgDegLatLong(value)
		return
	}
	pem.EditedPoleMarker.Pole.Lat = val
	pem.EditedPoleMarker.UpdateMarkerLatLng()
	pem.EditedPoleMarker.CenterOnMap(20)

}

func (pem *PoleEditModel) GetCities(vm *hvue.VM, value string, callback *js.Object) {
	pem = PoleEditModelFromJS(vm.Object)
	if len(value) < 3 {
		return
	}
	cities := make(map[string]bool)
	for _, pole := range pem.Polesite.Poles {
		cities[pole.City] = true
	}
	matchingCities := []*autocomplete.Result{}
	matchvalue := strings.ToLower(value)
	for cityName, _ := range cities {
		if cityName != value && strings.Contains(strings.ToLower(cityName), matchvalue) {
			matchingCities = append(matchingCities, autocomplete.NewResult(cityName))
		}
	}
	sort.Slice(matchingCities, func(i, j int) bool {
		return matchingCities[i].Value < matchingCities[j].Value
	})
	callback.Invoke(matchingCities)
}

func (pem *PoleEditModel) UpdatePoleCity(vm *hvue.VM, value string) {
	pem = PoleEditModelFromJS(vm.Object)
	pem.EditedPoleMarker.Pole.CheckInfoConsistency()
}

func (pem *PoleEditModel) UpdatePoleDegLong(vm *hvue.VM, value string) {
	pem = PoleEditModelFromJS(vm.Object)
	val, err := latlong.DegToDec(value)
	if err != nil {
		pem.errMsgDegLatLong(value)
		return
	}
	pem.EditedPoleMarker.Pole.Long = val
	pem.EditedPoleMarker.UpdateMarkerLatLng()
	pem.EditedPoleMarker.CenterOnMap(20)

}

func (pem *PoleEditModel) errMsgDegLatLong(value string) {
	msg := "impossible de lire la valeur: '" + value + "'\n\n"
	msg += "format attendu: 12°34'56,789"
	message.ErrorStr(pem.VM, msg, true)
}

// setDegLatLong sets EditedLat and EditedLong value from edited pole lat/long
func (pem *PoleEditModel) setDegLatLong() {
	pem.EditedLat = latlong.DecToDeg(pem.EditedPoleMarker.Pole.Lat)
	pem.EditedLong = latlong.DecToDeg(pem.EditedPoleMarker.Pole.Long)
}

func (pem *PoleEditModel) UpdatePoleDrag(vm *hvue.VM) {
	pem = PoleEditModelFromJS(vm.Object)
	pem.EditedPoleMarker.SetDraggable(pem.EditedPoleMarker.Draggable)
}

func (pem *PoleEditModel) UpdatePoleLatLong(vm *hvue.VM, value string) {
	pem = PoleEditModelFromJS(vm.Object)
	lls := strings.Split(value, ",")

	errmsg := func() {
		msg := "impossible de lire la latitude/longitude: '" + value + "'\n\n"
		msg += "format attendu: 12.123456, 12.123456"
		message.ErrorStr(pem.VM, msg, true)
	}

	if len(lls) < 2 {
		errmsg()
		return
	}
	slat := strings.Trim(lls[0], " ")
	slong := strings.Trim(lls[1], " ")
	lat, err := strconv.ParseFloat(slat, 64)
	if err != nil {
		errmsg()
		return
	}
	long, err := strconv.ParseFloat(slong, 64)
	if err != nil {
		errmsg()
		return
	}
	pem.EditedPoleMarker.SetLatLng(leaflet.NewLatLng(lat, long))
	pem.EditedPoleMarker.CenterOnMap(20)
}

func (pem *PoleEditModel) CenterOnEdited(vm *hvue.VM) {
	pem = PoleEditModelFromJS(vm.Object)
	vm.Emit("center-on-pole", pem.EditedPoleMarker.Pole)
	//pem.EditedPoleMarker.CenterOnMap(20)
}

func (pem *PoleEditModel) DeletePole(vm *hvue.VM) {
	pem = PoleEditModelFromJS(vm.Object)
	pem.VisibleDeletePole = false
	pem.VM.Emit("delete-pole", pem.EditedPoleMarker)
}

func (pem *PoleEditModel) DuplicatePole(vm *hvue.VM) {
	pem = PoleEditModelFromJS(vm.Object)
	pem.DuplicateContext.Model = pem.EditedPoleMarker
	pem.VisibleDuplicatePole = false
	pem.VM.Emit("duplicate-pole", pem.DuplicateContext)
}

func (pem *PoleEditModel) ChapterChange(vm *hvue.VM, val *js.Object) {
	vm.Emit("update:chapters", val)
}

func (pem *PoleEditModel) GetGMAPUrl(vm *hvue.VM, gps string) string {
	pem = PoleEditModelFromJS(vm.Object)
	//if !(pem.EditedPoleMarker.Object != nil && pem.EditedPoleMarker.Pole.DictRef != "") {
	//	return ""
	//}
	url := "http://maps.google.com/maps?q=" + gps
	return url
}

func (pem *PoleEditModel) GetDICTUrl(vm *hvue.VM) string {
	pem = PoleEditModelFromJS(vm.Object)
	if !(pem.EditedPoleMarker.Object != nil && pem.EditedPoleMarker.Pole.DictRef != "") {
		return ""
	}
	url := "https://apps.sogelink.fr/declaration/?tmstp=1574760154241#!/suivi/carto/302791299?idAgence=86986&filtre=" + pem.EditedPoleMarker.Pole.DictRef
	return url
}

func (pem *PoleEditModel) ClearSelection(vm *hvue.VM) {
	pem = PoleEditModelFromJS(vm.Object)
	pem.EditedPoleMarker.Map.ClearSelected()
}

func (pem *PoleEditModel) SelectByRef(vm *hvue.VM) {
	pem = PoleEditModelFromJS(vm.Object)
	editedPoleMarker := pem.EditedPoleMarker
	editedPoleMarker.Map.SelectByFilter(func(p *polemap.PoleMarker) bool {
		if p.Pole.Ref == editedPoleMarker.Pole.Ref {
			return true
		}
		return false
	})
	centerOn := append(editedPoleMarker.Map.SelectedPoleMarkers, editedPoleMarker)
	editedPoleMarker.Map.CenterOnPoleMarkers(centerOn)
}

func (pem *PoleEditModel) SelectByDict(vm *hvue.VM) {
	pem = PoleEditModelFromJS(vm.Object)
	editedPoleMarker := pem.EditedPoleMarker
	crit := ""
	if !tools.Empty(editedPoleMarker.Pole.DictRef) {
		crit = editedPoleMarker.Pole.DictRef
	} else if !tools.Empty(editedPoleMarker.Pole.DtRef) {
		crit = editedPoleMarker.Pole.DtRef
	}
	if crit == "" {
		editedPoleMarker.Map.ClearSelected()
		return
	}
	editedPoleMarker.Map.SelectByFilter(func(p *polemap.PoleMarker) bool {
		if !(p.Pole.DictRef != crit && p.Pole.DtRef != crit) {
			return true
		}
		return false
	})
	centerOn := append(editedPoleMarker.Map.SelectedPoleMarkers, editedPoleMarker)
	editedPoleMarker.Map.CenterOnPoleMarkers(centerOn)
}

func (pem *PoleEditModel) ApplyDict(vm *hvue.VM) {
	pem = PoleEditModelFromJS(vm.Object)
	editedPoleMarker := pem.EditedPoleMarker
	editedPoleMarker.Map.ApplyOnSelected(func(p *polemap.PoleMarker) {
		p.Pole.DictRef = editedPoleMarker.Pole.DictRef
		p.Pole.DictDate = editedPoleMarker.Pole.DictDate
		p.Pole.DictInfo = editedPoleMarker.Pole.DictInfo
		p.Pole.UpdateState()
		p.UpdateFromState()
	})
	editedPoleMarker.Map.RefreshPoleMarkersGroups()
}

func (pem *PoleEditModel) ApplyDa(vm *hvue.VM) {
	pem = PoleEditModelFromJS(vm.Object)
	editedPoleMarker := pem.EditedPoleMarker
	editedPoleMarker.Map.ApplyOnSelected(func(p *polemap.PoleMarker) {
		p.Pole.DaQueryDate = editedPoleMarker.Pole.DaQueryDate
		p.Pole.DaStartDate = editedPoleMarker.Pole.DaStartDate
		p.Pole.DaEndDate = editedPoleMarker.Pole.DaEndDate
		p.Pole.UpdateState()
		p.UpdateFromState()
	})
	editedPoleMarker.Map.RefreshPoleMarkersGroups()
}
