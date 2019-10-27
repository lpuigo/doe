package poletable

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	ps "github.com/lpuig/ewin/doe/website/frontend/model/polesite"
	"github.com/lpuig/ewin/doe/website/frontend/model/polesite/poleconst"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
	"strings"
)

const template string = `<el-container  style="height: 100%; padding: 0px">
    <el-header style="height: auto; margin-top: 5px">
        <el-row type="flex" align="middle" :gutter="5">
			<el-col :span="2" style="text-align: right"><span>Mode d'affichage:</span></el-col>
			<el-col :span="2">
			  <el-select v-model="mode" placeholder="Select" size="mini" @change="ChangeMode">
				<el-option
				  v-for="item in GetModes()"
				  :key="item.value"
				  :label="item.label"
				  :value="item.value">
				</el-option>
			  </el-select>
			</el-col>
        </el-row>
    </el-header>
    <div style="height: 100%;overflow-x: hidden;overflow-y: auto;padding: 0px 0px; margin-top: 8px">
		<pole-table-creation v-if="mode == 'creation'"
				:user="user"
				:polesite="polesite"
				:filter="filter"
				:filtertype="filtertype"
				@pole-selected="SetSelectedPole"
		></pole-table-creation>
		<pole-table-followup v-if="mode == 'followup'"
				:user="user"
				:polesite="polesite"
				:filter="filter"
				:filtertype="filtertype"
				@pole-selected="SetSelectedPole"
		></pole-table-followup>
		<pole-table-billing v-if="mode == 'billing'"
				:user="user"
				:polesite="polesite"
				:filter="filter"
				:filtertype="filtertype"
				@pole-selected="SetSelectedPole"
		></pole-table-billing>
    </div>
</el-container>
`

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Registration

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("pole-table", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		registerComponentTable("creation"),
		registerComponentTable("followup"),
		registerComponentTable("billing"),
		hvue.Template(template),
		hvue.Props("user", "polesite", "filter", "filtertype", "mode"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewPoleTablesModel(vm)
		}),
		hvue.MethodsOf(&PoleTablesModel{}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type PoleTablesModel struct {
	*js.Object

	Polesite   *ps.Polesite `js:"polesite"`
	User       *fm.User     `js:"user"`
	Filter     string       `js:"filter"`
	FilterType string       `js:"filtertype"`
	Mode       string       `js:"mode"`

	VM *hvue.VM `js:"VM"`
}

func NewPoleTablesModel(vm *hvue.VM) *PoleTablesModel {
	rtm := &PoleTablesModel{Object: tools.O()}
	rtm.Polesite = nil
	rtm.User = fm.NewUser()
	rtm.Filter = ""
	rtm.FilterType = ""
	rtm.Mode = "creation"
	rtm.VM = vm
	return rtm
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Actions related Methods

func (ptm *PoleTablesModel) GetModes() []*elements.ValueLabel {
	return []*elements.ValueLabel{
		elements.NewValueLabel("creation", "Création"),
		elements.NewValueLabel("followup", "Suivi"),
		elements.NewValueLabel("billing", "Facturation"),
	}
}

func (ptm *PoleTablesModel) ChangeMode(vm *hvue.VM, mode string) {
	print("ChangeMode", mode)
	vm.Emit("update:mode", mode)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Event Related Methods

func (ptm *PoleTablesModel) SetSelectedPole(vm *hvue.VM, p *ps.Pole) {
	vm.Emit("pole-selected", p)
}

//func (ptm *PoleTablesModel) AddPole(vm *hvue.VM) {
//	message.InfoStr(vm, "AddPole non encore implémenté", false)
//}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Formatting Related Methods

//func (ptm *PoleTablesModel) TableRowClassName(rowInfo *js.Object) string {
//	p := &ps.Pole{Object: rowInfo.Get("row")}
//	return ps.PoleRowClassName(p.State)
//}
//
//func (ptm *PoleTablesModel) HeaderCellStyle() string {
//	return "background: #a1e6e6;"
//}
//
//func (ptm *PoleTablesModel) EndDate(d string, delay int) string {
//	if d == "" {
//		return ""
//	}
//	return date.DateString(date.After(d, delay))
//}
//
//func (ptm *PoleTablesModel) FormatDate(r, c *js.Object, d string) string {
//	return date.DateString(d)
//}
//
//func (ptm *PoleTablesModel) FormatState(r, c *js.Object, d string) string {
//	return ps.PoleStateLabel(d)
//}
//
//func (ptm *PoleTablesModel) FormatProduct(p *ps.Pole) string {
//	return strings.Join(p.Product, "\n")
//}
//
//func (ptm *PoleTablesModel) FormatType(p *ps.Pole) string {
//	return p.Material + " " + strconv.Itoa(p.Height) + "m"
//}
//
//func (ptm *PoleTablesModel) FormatActors(vm *hvue.VM, p *ps.Pole) string {
//	ptm = &PoleTablesModel{Object: vm.Object}
//	client := ptm.User.GetClientByName(ptm.Polesite.Client)
//	actors := []string{}
//	for _, actId := range p.Actors {
//		actor := client.GetActorBy(actId)
//		if actor != nil {
//			actors = append(actors, actor.LastName)
//		}
//	}
//	return strings.Join(actors, "\n")
//}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Sorting Methods

//func (ptm *PoleTablesModel) SortDate(attrib string) func(obj *js.Object) string {
//	return func(obj *js.Object) string {
//		val := obj.Get(attrib).String()
//		if val=="" {
//			return "9999-12-31"
//		}
//		return val
//	}
//}
//
//func (ptm *PoleTablesModel) SortState(a, b *ps.Pole) int {
//	la := ps.PoleStateLabel(a.State)
//	lb := ps.PoleStateLabel(b.State)
//	if la < lb {
//		return -1
//	}
//	if la == lb {
//		return 0
//	}
//	return 1
//}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Row Filtering Related Methods

func (ptm *PoleTablesModel) GetFilteredPole() []*ps.Pole {
	if ptm.FilterType == poleconst.FilterValueAll && ptm.Filter == "" {
		return ptm.Polesite.Poles
	}

	res := []*ps.Pole{}
	expected := strings.ToUpper(ptm.Filter)
	filter := func(p *ps.Pole) bool {
		sis := p.SearchString(ptm.FilterType)
		if sis == "" {
			return false
		}
		return strings.Contains(strings.ToUpper(sis), expected)
	}
	for _, pole := range ptm.Polesite.Poles {
		if filter(pole) {
			res = append(res, pole)
		}
	}
	return res
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Column Filtering Related Methods

//func (ptm *PoleTablesModel) FilterHandler(value string, p *js.Object, col *js.Object) bool {
//	prop := col.Get("property").String()
//	return p.Get(prop).String() == value
//}
//
//func (ptm *PoleTablesModel) FilterList(vm *hvue.VM, prop string) []*elements.ValText {
//	ptm = &PoleTablesModel{Object: vm.Object}
//	count := map[string]int{}
//	attribs := []string{}
//
//	var translate func(string) string
//	switch prop {
//	case "State":
//		translate = func(val string) string {
//			return ps.PoleStateLabel(val)
//		}
//	default:
//		translate = func(val string) string { return val }
//	}
//
//	for _, psi := range ptm.Polesite.Poles {
//		attrib := psi.Object.Get(prop).String()
//		if _, exist := count[attrib]; !exist {
//			attribs = append(attribs, attrib)
//		}
//		count[attrib]++
//	}
//	sort.Strings(attribs)
//	res := []*elements.ValText{}
//	for _, a := range attribs {
//		fa := a
//		if fa == "" {
//			fa = "Vide"
//		}
//		res = append(res, elements.NewValText(a, translate(fa)+" ("+strconv.Itoa(count[a])+")"))
//	}
//	return res
//}
//
//func (ptm *PoleTablesModel) FilteredStatusValue() []string {
//	res := []string{
//		//poleconst.PsStatusNew,
//		//poleconst.PsStatusInProgress,
//		//poleconst.PsStatusBlocked,
//		//poleconst.PsStatusCancelled,
//		//poleconst.PsStatusDone,
//	}
//	return res
//}
