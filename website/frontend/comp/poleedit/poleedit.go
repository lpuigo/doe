package poleedit

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/polemap"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/model/polesite"
	"github.com/lpuig/ewin/doe/website/frontend/model/polesite/poleconst"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	date "github.com/lpuig/ewin/doe/website/frontend/tools/date"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements/message"
	"github.com/lpuig/ewin/doe/website/frontend/tools/leaflet"
	"strconv"
	"strings"
)

const template string = `<div>
	<div class="header-menu-container">
		<h1>
			Poteau: {{editedpolemarker.Pole.Ref}}
		</h1>
        <span></span>
        <el-popover placement="bottom" width="360" title="Suppression du poteau"
                    v-model="VisibleDeletePole">
            <p>Confirmez la suppression du poteau {{editedpolemarker.Pole.Ref}} ?</p>
            <div style="text-align: right; margin: 0; margin-top: 10px">
                <el-button size="mini" @click="VisibleDeletePole = false">Annuler</el-button>
                <el-button type="danger" plain size="mini" @click="DeletePole">Confirmer</el-button>
            </div>
            <el-button slot="reference" type="danger" plain class="icon" icon="far fa-trash-alt icon--big" size="mini"></el-button>
        </el-popover>
		<el-button class="icon" icon="fas fa-crosshairs icon--big" @click="CenterOnEdited" size="mini"></el-button>
	</div>
    <!-- Référence -->
    <el-row :gutter="5" type="flex" align="middle" class="spaced">
        <el-col :span="6" class="align-right">Référence:</el-col>
        <el-col :span="18">
            <el-input placeholder="Référence"
                      v-model="editedpolemarker.Pole.Ref" clearable size="mini"
					  @change="UpdateTooltip()"
            ></el-input>
        </el-col>
    </el-row>
    
    <!-- Lat / Long -->
    <el-row :gutter="5" type="flex" align="middle" class="spaced">
        <el-col :span="6" class="align-right">Lat / Long:</el-col>
        <el-col :span="18">
            <el-input v-model="editedlatlong" size="mini" @input="UpdatePoleLatLong"></el-input>
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
    <el-row :gutter="5" type="flex" align="middle" class="spaced">
        <el-col :span="6" class="align-right">Adresse:</el-col>
        <el-col :span="18">
            <el-input placeholder="Adresse"
                      v-model="editedpolemarker.Pole.Address" clearable size="mini"
            ></el-input>
        </el-col>
    </el-row>

    <!-- Etiquette -->
    <el-row :gutter="5" type="flex" align="middle" class="doublespaced">
        <el-col :span="6" class="align-right">Etiquette:</el-col>
        <el-col :span="18">
            <el-input placeholder="Etiquette"
                      v-model="editedpolemarker.Pole.Sticker" clearable size="mini"
            ></el-input>
        </el-col>
    </el-row>

    <!-- Matériau -->
    <el-row :gutter="5" type="flex" align="middle" class="doublespaced">
        <el-col :span="6" class="align-right">Matériau:</el-col>
        <el-col :span="18">
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
    </el-row>

    <!-- Hauteur -->
    <el-row :gutter="5" type="flex" align="middle" class="spaced">
        <el-col :span="6" class="align-right">Hauteur:</el-col>
        <el-col :span="18">
            <el-input-number v-model="editedpolemarker.Pole.Height" size="mini" controls-position="right" :precision="0" :min="6" :max="12" style="width: 100%"
            ></el-input-number>
        </el-col>
    </el-row>

    <!-- DT DICT -->
    <el-row :gutter="5" type="flex" align="middle" class="doublespaced">
        <el-col :span="6" class="align-right">DT:</el-col>
        <el-col :span="18">
            <el-input placeholder="Référence DT"
                      v-model="editedpolemarker.Pole.DtRef" clearable size="mini"
            ></el-input>
        </el-col>
    </el-row>
    <el-row :gutter="5" type="flex" align="middle" class="spaced">
        <el-col :span="6" class="align-right">DICT:</el-col>
        <el-col :span="10">
            <el-input placeholder="Référence DICT"
                      v-model="editedpolemarker.Pole.DictRef" clearable size="mini"
            ></el-input>
        </el-col>
        <el-col :span="8">
            <el-date-picker format="dd/MM/yyyy" placeholder="Date" size="mini"
                            style="width: 100%" type="date"
                            v-model="editedpolemarker.Pole.DictDate"
                            value-format="yyyy-MM-dd"
                            :picker-options="{firstDayOfWeek:1}"
                            :clearable="false"
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

    <!-- Kizeo -->
    <el-row :gutter="5" type="flex" align="middle" class="spaced">
        <el-col :span="6" class="align-right">Kizeo:</el-col>
        <el-col :span="18">
            <el-input placeholder="Référence Kizeo"
                      v-model="editedpolemarker.Pole.Kizeo" clearable size="mini"
            ></el-input>
        </el-col>
    </el-row>

    <!-- Date Aspiratrice -->
    <el-row :gutter="5" type="flex" align="middle" class="doublespaced">
        <el-col :span="6" class="align-right">Aspiratrice:</el-col>
        <el-col :span="18">
            <el-date-picker format="dd/MM/yyyy" placeholder="Date" size="mini"
                            style="width: 100%" type="date"
                            v-model="editedpolemarker.Pole.AspiDate"
                            value-format="yyyy-MM-dd"
                            :picker-options="{firstDayOfWeek:1, disabledDate(time) { return time.getTime() > Date.now(); }}"
                            :clearable="false"
            ></el-date-picker>
        </el-col>
    </el-row>

    <!-- Produits -->
    <el-row :gutter="5" type="flex" align="middle" class="spaced">
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

    <!-- Commentaire -->
    <el-row :gutter="5" type="flex" align="middle" class="doublespaced">
        <el-col :span="6" class="align-right">Commentaire:</el-col>
        <el-col :span="18">
            <el-input type="textarea" :autosize="{ minRows: 1, maxRows: 2}" placeholder="Commentaire"
                      v-model="editedpolemarker.Pole.Comment" clearable size="mini"
            ></el-input>
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
                ></el-option>
            </el-select>
        </el-col>
    </el-row>

    <!-- Date Status -->
    <el-row v-if="ShowDate" :gutter="5" type="flex" align="middle" class="spaced">
        <el-col :span="6" class="align-right">Date:</el-col>
        <el-col :span="18">
            <el-date-picker format="dd/MM/yyyy" placeholder="Date" size="mini"
                            style="width: 100%" type="date"
                            v-model="editedpolemarker.Pole.Date"
                            value-format="yyyy-MM-dd"
                            :picker-options="{firstDayOfWeek:1, disabledDate(time) { return time.getTime() > Date.now(); }}"
                            :clearable="false"
            ></el-date-picker>
        </el-col>
    </el-row>


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
		hvue.Props("editedpolemarker", "user", "polesite"),
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
				return strconv.FormatFloat(pem.EditedPoleMarker.Pole.Lat, 'f', 8, 64) +
					", " +
					strconv.FormatFloat(pem.EditedPoleMarker.Pole.Long, 'f', 8, 64)
			}),
		hvue.Computed(
			"ShowDate",
			func(vm *hvue.VM) interface{} {
				pem := PoleEditModelFromJS(vm.Object)
				return pem.EditedPoleMarker.Pole.State == poleconst.StateDone
			}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type PoleEditModel struct {
	*js.Object
	EditedPoleMarker  *polemap.PoleMarker `js:"editedpolemarker"`
	VisibleDeletePole bool                `js:"VisibleDeletePole"`
	User              *fm.User            `js:"user"`
	Polesite          *polesite.Polesite  `js:"polesite"`

	VM *hvue.VM `js:"VM"`
}

func PoleEditModelFromJS(obj *js.Object) *PoleEditModel {
	return &PoleEditModel{Object: obj}
}

func NewPoleEditModel(vm *hvue.VM) *PoleEditModel {
	pem := &PoleEditModel{Object: tools.O()}
	pem.VisibleDeletePole = false
	pem.User = fm.NewUser()
	pem.Polesite = nil
	pem.VM = vm
	pem.EditedPoleMarker = polemap.DefaultPoleMarker()
	return pem
}

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

func (pem *PoleEditModel) GetStates() []*elements.ValueLabel {
	return polesite.GetStatesValueLabel()
}

func (pem *PoleEditModel) UpdateState(vm *hvue.VM) {
	pem = PoleEditModelFromJS(vm.Object)
	ep := pem.EditedPoleMarker
	ep.UpdateFromState()
	ep.Refresh()
	if ep.Pole.State == poleconst.StateDone {
		ep.Pole.Date = date.TodayAfter(0)
	}
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
	pem.VM.Emit("deletepole", pem.EditedPoleMarker)
	pem.VisibleDeletePole = false
}
