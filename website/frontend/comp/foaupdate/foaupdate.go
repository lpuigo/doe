package foaupdate

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	fmfoa "github.com/lpuig/ewin/doe/website/frontend/model/foasite"
	"github.com/lpuig/ewin/doe/website/frontend/model/foasite/foaconst"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/date"
	"strings"
)

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("foa-update", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		//rippullingdistinfo.RegisterComponent(),
		//ripstateupdate.RegisterComponent(),
		hvue.Template(template),
		hvue.Props("value", "user", "filter", "filtertype"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewFoaUpdateModel(vm)
		}),
		hvue.MethodsOf(&FoaUpdateModel{}),
		hvue.Computed("filteredJunctions", func(vm *hvue.VM) interface{} {
			rpum := FoaUpdateModelFromJS(vm.Object)
			return rpum.GetFilteredJunctions()
		}),
		//hvue.Filter("FormatTronconRef", func(vm *hvue.VM, value *js.Object, args ...*js.Object) interface{} {
		//	rpum := FoaUpdateModelFromJS(vm.Object)
		//	t := &fm.Troncon{Object: value}
		//	return rpum.GetFormatTronconRef(t)
		//}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type FoaUpdateModel struct {
	*js.Object

	Foasite      *fmfoa.FoaSite `js:"value"`
	SelectedFoas *fmfoa.FoaSite `js:"SelectedFoas"`
	User         *fm.User       `js:"user"`
	Filter       string         `js:"filter"`
	FilterType   string         `js:"filtertype"`

	VM *hvue.VM `js:"VM"`
}

func NewFoaUpdateModel(vm *hvue.VM) *FoaUpdateModel {
	rpum := &FoaUpdateModel{Object: tools.O()}
	rpum.VM = vm
	rpum.Foasite = fmfoa.NewFoaSite()
	rpum.SelectedFoas = fmfoa.NewFoaSite()
	rpum.User = fm.NewUser()
	rpum.Filter = ""
	rpum.FilterType = foaconst.FilterValueAll
	return rpum
}

func FoaUpdateModelFromJS(o *js.Object) *FoaUpdateModel {
	return &FoaUpdateModel{Object: o}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Utils Methods

func (fum *FoaUpdateModel) ClearSelection() {
	fum.SelectedFoas.Foas = []*fmfoa.Foa{}
	fum.VM.Refs("foaTable").Call("clearSelection")
}

func (fum *FoaUpdateModel) SetSelected(f *fmfoa.Foa) {
	fum.SelectedFoas.Foas = []*fmfoa.Foa{f}
	foaTable := fum.VM.Refs("foaTable")
	foaTable.Call("clearSelection")
	foaTable.Call("toggleRowSelection", f, true)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// HTML Methods

func (fum *FoaUpdateModel) HandleSelectionChange(vm *hvue.VM, val *js.Object) {
	fum = FoaUpdateModelFromJS(vm.Object)
	fum.SelectedFoas.Foas = []*fmfoa.Foa{}
	val.Call("forEach", func(foa *fmfoa.Foa) {
		fum.SelectedFoas.Foas = append(fum.SelectedFoas.Foas, foa)
	})
}

func (fum *FoaUpdateModel) AddFoa(vm *hvue.VM) {
	fum.VM.Emit("add-foa")
}

func (fum *FoaUpdateModel) EditFoa(vm *hvue.VM, f *fmfoa.Foa) {
	fum.SetSelected(f)
	fum.VM.Emit("update-state", fum.SelectedFoas, f)
}

func (fum *FoaUpdateModel) EditSelectedFoas(vm *hvue.VM) {
	fum.VM.Emit("update-state", fum.SelectedFoas, nil)
}

func (fum *FoaUpdateModel) GetFilteredJunctions() []*fmfoa.Foa {
	if fum.FilterType == foaconst.FilterValueAll && fum.Filter == "" {
		return fum.Foasite.Foas
	}
	res := []*fmfoa.Foa{}
	expected := strings.ToUpper(fum.Filter)
	filter := func(f *fmfoa.Foa) bool {
		sis := f.SearchString(fum.FilterType)
		if sis == "" {
			return false
		}
		return strings.Contains(strings.ToUpper(sis), expected)
	}

	for _, f := range fum.Foasite.Foas {
		if filter(f) {
			res = append(res, f)
		}
	}
	return res
}

func (fum *FoaUpdateModel) TableRowClassName(rowInfo *js.Object) string {
	f := &fmfoa.Foa{Object: rowInfo.Get("row")}
	return f.GetRowStyle()
}

func (fum *FoaUpdateModel) GetActors(vm *hvue.VM, f *fmfoa.Foa) string {
	fum = FoaUpdateModelFromJS(vm.Object)
	client := fum.User.GetClientByName(fum.Foasite.Client)
	if client == nil {
		return ""
	}

	res := []string{}
	for _, actId := range f.State.Actors {
		actor := client.GetActorBy(actId)
		if actor == nil {
			continue
		}
		res = append(res, actor.GetRef())
	}
	return strings.Join(res, "\n")
}

func (fum *FoaUpdateModel) FormatDate(r, c *js.Object, d string) string {
	return date.DateString(d)
}

func (fum *FoaUpdateModel) FormatStatus(r, c *js.Object, d string) string {
	return fmfoa.FoaStateLabel(d)
}

/*
func (rjum *FoaUpdateModel) FilterHandler(value string, p *js.Object, col *js.Object) bool {
	prop := col.Get("property").String()
	print("FilterHandler", prop, p.Get(prop).String())
	return p.Get(prop).String() == value
}

func (rjum *FoaUpdateModel) FilterList(vm *hvue.VM, prop string) []*elements.ValText {
	rjum = FoaUpdateModelFromJS(vm.Object)
	count := map[string]int{}
	attribs := []string{}

	var translate func(string) string
	switch prop {
	case "State.Status":
		translate = func(val string) string {
			return fmrip.GetStatusLabel(val)
		}
	default:
		translate = func(val string) string { return val }
	}

	for _, junction := range rjum.GetFilteredJunctions() {
		attrib := junction.Object.Get(prop).String()
		if _, exist := count[attrib]; !exist {
			attribs = append(attribs, attrib)
		}
		count[attrib]++
	}
	sort.Strings(attribs)
	res := []*elements.ValText{}
	for _, a := range attribs {
		fa := a
		if fa == "" {
			fa = "Vide"
		}
		res = append(res, elements.NewValText(a, translate(fa)+" ("+strconv.Itoa(count[a])+")"))
	}
	return res
}

*/
