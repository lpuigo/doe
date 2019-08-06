package poleedit

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/polemap"
	"github.com/lpuig/ewin/doe/website/frontend/model/polesite"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
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
		<el-button class="icon" icon="fas fa-crosshairs icon--big" @click="CenterOnEdited" size="mini"></el-button>
	</div>
    <!-- Référence -->
    <el-row :gutter="5" type="flex" align="middle" class="spaced">
        <el-col :span="6">Référence:</el-col>
        <el-col :span="18">
            <el-input placeholder="Référence"
                      v-model="editedpolemarker.Pole.Ref" clearable size="mini"
					  @change="UpdateTooltip()"
            ></el-input>
        </el-col>
    </el-row>
    
    <!-- Lat / Long -->
    <el-row :gutter="5" type="flex" align="middle" class="spaced">
        <el-col :span="6">Lat / Long:</el-col>
        <el-col :span="18">
            <el-input v-model="editedlatlong" size="mini" @input="UpdatePoleLatLong"></el-input>
        </el-col>
    </el-row>
    
    <!-- Ville -->
    <el-row :gutter="5" type="flex" align="middle" class="spaced">
        <el-col :span="6">Ville:</el-col>
        <el-col :span="18">
            <el-input placeholder="Ville"
                      v-model="editedpolemarker.Pole.City" clearable size="mini"
            ></el-input>
        </el-col>
    </el-row>

    <!-- Adresse -->
    <el-row :gutter="5" type="flex" align="middle" class="spaced">
        <el-col :span="6">Adresse:</el-col>
        <el-col :span="18">
            <el-input placeholder="Adresse"
                      v-model="editedpolemarker.Pole.Address" clearable size="mini"
            ></el-input>
        </el-col>
    </el-row>

    <!-- Etiquette -->
    <el-row :gutter="5" type="flex" align="middle" class="doublespaced">
        <el-col :span="6">Etiquette:</el-col>
        <el-col :span="18">
            <el-input placeholder="Etiquette"
                      v-model="editedpolemarker.Pole.Sticker" clearable size="mini"
            ></el-input>
        </el-col>
    </el-row>

    <!-- Matériau -->
    <el-row :gutter="5" type="flex" align="middle" class="doublespaced">
        <el-col :span="6">Matériau:</el-col>
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
        <el-col :span="6">Hauteur:</el-col>
        <el-col :span="18">
            <el-input-number v-model="editedpolemarker.Pole.Height" size="mini" controls-position="right" :precision="0" :min="7" :max="12" style="width: 100%"
            ></el-input-number>
        </el-col>
    </el-row>

    <!-- DT DICT -->
    <el-row :gutter="5" type="flex" align="middle" class="doublespaced">
        <el-col :span="6">DT:</el-col>
        <el-col :span="18">
            <el-input placeholder="Référence DT"
                      v-model="editedpolemarker.Pole.DtRef" clearable size="mini"
            ></el-input>
        </el-col>
    </el-row>
    <el-row :gutter="5" type="flex" align="middle" class="spaced">
        <el-col :span="6">DICT:</el-col>
        <el-col :span="18">
            <el-input placeholder="Référence DICT"
                      v-model="editedpolemarker.Pole.DictRef" clearable size="mini"
            ></el-input>
        </el-col>
    </el-row>
    <el-row :gutter="5" type="flex" align="middle" class="spaced">
        <el-col :span="6">Info DICT:</el-col>
        <el-col :span="18">
            <el-input type="textarea" :autosize="{ minRows: 1, maxRows: 4}" placeholder="Information DICT"
                      v-model="editedpolemarker.Pole.DictInfo" clearable size="mini"
            ></el-input>
        </el-col>
    </el-row>

    <!-- Kizeo -->
    <el-row :gutter="5" type="flex" align="middle" class="spaced">
        <el-col :span="6">Kizeo:</el-col>
        <el-col :span="18">
            <el-input placeholder="Référence Kizeo"
                      v-model="editedpolemarker.Pole.Kizeo" clearable size="mini"
            ></el-input>
        </el-col>
    </el-row>

    <!-- Produits -->
    <el-row :gutter="5" type="flex" align="middle" class="doublespaced">
        <el-col :span="6">Produits:</el-col>
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
    
    <!-- Status -->
    <el-row :gutter="5" type="flex" align="middle" class="doublespaced">
        <el-col :span="6">Status:</el-col>
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
		hvue.Props("editedpolemarker"),
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
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type PoleEditModel struct {
	*js.Object
	EditedPoleMarker *polemap.PoleMarker `js:"editedpolemarker"`
	Product          string              `js:"Product"`

	VM *hvue.VM `js:"VM"`
}

func PoleEditModelFromJS(obj *js.Object) *PoleEditModel {
	return &PoleEditModel{Object: obj}
}

func NewPoleEditModel(vm *hvue.VM) *PoleEditModel {
	pem := &PoleEditModel{Object: tools.O()}
	pem.Product = ""
	pem.VM = vm
	pem.EditedPoleMarker = polemap.DefaultPoleMarker()
	return pem
}

func (pem *PoleEditModel) GetStates() []*elements.ValueLabel {
	return polesite.GetStatesValueLabel()
}

func (pem *PoleEditModel) UpdateState(vm *hvue.VM) {
	pem = PoleEditModelFromJS(vm.Object)
	pem.EditedPoleMarker.UpdateFromState()
	pem.EditedPoleMarker.Refresh()
}

func (pem *PoleEditModel) GetMaterials() []*elements.ValueLabel {
	return polesite.GetMaterialsValueLabel()
}

func (pem *PoleEditModel) GetProducts() []*elements.ValueLabel {
	return polesite.GetProductsValueLabel()
}

func (pem *PoleEditModel) UpdateProduct(vm *hvue.VM) {
	pem = PoleEditModelFromJS(vm.Object)
	pem.EditedPoleMarker.Pole.Get("Product").Set(pem.Product, 1)
	pem.Product = ""
}

func (pem *PoleEditModel) RemoveProduct(vm *hvue.VM, prd string) {
	pem = PoleEditModelFromJS(vm.Object)
	pem.EditedPoleMarker.Pole.Get("Product").Set(prd, 0)
	pem.EditedPoleMarker.Pole.Get("Product").Delete(prd)
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
	pem.EditedPoleMarker.CenterOnMap(20)
}
