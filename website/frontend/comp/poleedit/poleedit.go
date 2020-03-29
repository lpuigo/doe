package poleedit

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/polemap"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/model/polesite"
	"github.com/lpuig/ewin/doe/website/frontend/model/polesite/poleconst"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/date"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements/message"
	"github.com/lpuig/ewin/doe/website/frontend/tools/latlong"
	"github.com/lpuig/ewin/doe/website/frontend/tools/leaflet"
	"strconv"
	"strings"
)

const template string = `<div>
	<div class="header-menu-container">
		<h1 class="blue">
			Poteau: <span style="font-size: 1.5em">{{editedpolemarker.Pole.Ref}} {{editedpolemarker.Pole.Sticker}}</span>
		</h1>
        <span></span>
        
        <el-popover placement="bottom" width="360" title="Suppression du poteau"
                    v-model="VisibleDeletePole">
            <p>Confirmez la suppression du poteau {{editedpolemarker.Pole.Ref}} {{editedpolemarker.Pole.Sticker}}?</p>
            <div style="text-align: right; margin: 0; margin-top: 10px">
                <el-button size="mini" @click="VisibleDeletePole = false">Annuler</el-button>
                <el-button type="danger" plain size="mini" @click="DeletePole">Confirmer</el-button>
            </div>
            <el-button slot="reference" type="danger" plain class="icon" icon="far fa-trash-alt icon--big" size="mini"></el-button>
        </el-popover>
        
        <el-popover placement="bottom" width="360" title="Duplication du poteau"
                    v-model="VisibleDuplicatePole">
            <p>Confirmez la duplication du poteau {{editedpolemarker.Pole.Ref}} {{editedpolemarker.Pole.Sticker}}?</p>
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
            <div style="text-align: right; margin: 0; margin-top: 10px">
                <el-button size="mini" @click="VisibleDuplicatePole = false">Annuler</el-button>
                <el-button type="warning" plain size="mini" @click="DuplicatePole">Confirmer</el-button>
            </div>
            <el-button slot="reference" type="warning" plain class="icon" icon="far fa-clone icon--big" size="mini"></el-button>
        </el-popover>

        <el-tooltip content="Zoomer sur la carte" placement="bottom" effect="light" open-delay="500">
            <el-button class="icon" icon="fas fa-crosshairs icon--big" @click="CenterOnEdited" size="mini"></el-button>
        </el-tooltip>
	</div>
    <!-- Référence -->
    <el-row :gutter="5" type="flex" align="middle" class="spaced">
        <el-col :span="6" class="align-right">Référence:</el-col>
        <el-col :span="18">
            <el-input placeholder="Référence"
                      v-model="editedpolemarker.Pole.Ref" clearable size="mini"
					  @input="UpdateTooltip()"
            ></el-input>
        </el-col>
    </el-row>
    
    <!-- Etiquette -->
    <el-row :gutter="5" type="flex" align="middle" class="spaced">
        <el-col :span="6" class="align-right">Appui (Etiquette):</el-col>
        <el-col :span="18">
            <el-input placeholder="Etiquette"
                      v-model="editedpolemarker.Pole.Sticker" clearable size="mini"
					  @input="UpdateTooltip()"
            ></el-input>
        </el-col>
    </el-row>

    <!-- Commentaire -->
    <el-row :gutter="5" type="flex" align="middle" class="doublespaced">
        <el-col :span="6" class="align-right">Commentaire:</el-col>
        <el-col :span="18">
            <el-input type="textarea" :autosize="{ minRows: 1, maxRows: 2}" placeholder="Commentaire"
                      v-model="editedpolemarker.Pole.Comment" clearable size="mini"
            ></el-input>
        </el-col>
    </el-row>

    <el-collapse v-model="chapters" @change="ChapterChange">
        <el-collapse-item name="1">
            <template slot="title">
                <h1 class="title">Adresse: <a v-if="editedlatlong != ''" :href="GetGMAPUrl(editedlatlong)" rel="noopener noreferrer" target="_blank">(GMap)</a><span class="blue"> {{editedpolemarker.Pole.Address}}</span></h1>
            </template>
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
                    <el-input placeholder="Ville"
                              v-model="editedpolemarker.Pole.City" clearable size="mini"
                    ></el-input>
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

        <el-collapse-item name="2">
            <template slot="title">
                <h1 class="title">Travaux: <span class="blue">{{editedpolemarker.Pole.Material}} {{editedpolemarker.Pole.Height}}m</span></h1>
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
                               @change=""
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

        <el-collapse-item name="3">
            <template slot="title">
                <h1 class="title">DT, DICT et DA&nbsp;<a v-if="editedpolemarker.Pole.DictRef != ''" :href="GetDICTUrl()" rel="noopener noreferrer" target="_blank">(Carte DICT.fr)</a></h1>
            </template>
            <!-- DT -->
            <el-row :gutter="5" type="flex" align="middle" class="spaced">
                <el-col :span="6" class="align-right">DT:</el-col>
                <el-col :span="18">
                    <el-input placeholder="Référence DT"
                              v-model="editedpolemarker.Pole.DtRef" clearable size="mini"
                    ></el-input>
                </el-col>
            </el-row>
            <!-- DICT -->
            <el-row :gutter="5" type="flex" align="middle" class="spaced">
                <el-col :span="6" class="align-right">DICT:</el-col>
                <el-col :span="10">
                    <el-input placeholder="Référence DICT"
                              v-model="editedpolemarker.Pole.DictRef" clearable size="mini"
                    ></el-input>
                </el-col>
                <el-col :span="8">
                    <el-date-picker format="dd/MM/yyyy" placeholder="Début" size="mini"
                                    style="width: 100%" type="date"
                                    v-model="editedpolemarker.Pole.DictDate"
                                    value-format="yyyy-MM-dd"
                                    :picker-options="{firstDayOfWeek:1}"
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
            <el-row :gutter="5" type="flex" align="middle">
                <el-col :span="6" class="align-right">DA:</el-col>
                <el-col :span="9">
                    <el-date-picker format="dd/MM/yyyy" placeholder="Début" size="mini"
                                    style="width: 100%" type="date"
                                    v-model="editedpolemarker.Pole.DaStartDate"
                                    value-format="yyyy-MM-dd"
                                    :picker-options="{firstDayOfWeek:1}"
                    ></el-date-picker>
                </el-col>
                <el-col :span="9">
                    <el-date-picker format="dd/MM/yyyy" placeholder="Fin" size="mini"
                                    style="width: 100%" type="date"
                                    v-model="editedpolemarker.Pole.DaEndDate"
                                    value-format="yyyy-MM-dd"
                                    :picker-options="{firstDayOfWeek:1}"
                    ></el-date-picker>
                </el-col>
            </el-row>
        </el-collapse-item>

        <el-collapse-item name="4">
            <template slot="title">
                <h1 class="title">Etat: <span class="blue">{{FormatState(editedpolemarker.Pole)}}</span></h1>
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
                               @change=""
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
                                    @change="UpdateState()"
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

func (pem *PoleEditModel) UpdateState(vm *hvue.VM) {
	pem = PoleEditModelFromJS(vm.Object)
	ep := pem.EditedPoleMarker
	if ep.Pole.State == poleconst.StateDone {
		if tools.Empty(ep.Pole.Date) {
			ep.Pole.Date = date.TodayAfter(0)
		}
		if !tools.Empty(ep.Pole.AttachmentDate) {
			ep.Pole.State = poleconst.StateAttachment
		}
	}
	ep.UpdateFromState()
	ep.Refresh()
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
