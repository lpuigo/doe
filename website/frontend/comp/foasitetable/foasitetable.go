package foasitetable

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/ripprogressbar"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	fs "github.com/lpuig/ewin/doe/website/frontend/model/foasite"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/date"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
	"sort"
	"strconv"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Registration

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("foasite-table", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		ripprogressbar.RegisterComponent(),
		hvue.Template(template),
		hvue.Props("foasiteinfos", "user"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewFoaSiteTableModel(vm)
		}),
		hvue.MethodsOf(&FoaSiteTableModel{}),
		hvue.Computed("filteredFoasites", func(vm *hvue.VM) interface{} {
			rtm := &FoaSiteTableModel{Object: vm.Object}
			if rtm.Filter == "" {
				return rtm.Foasiteinfos
			}
			res := []*fm.FoaSiteInfo{}
			for _, psi := range rtm.Foasiteinfos {
				if psi.TextFiltered(rtm.Filter) {
					res = append(res, psi)
				}
			}
			return res
		}),
		hvue.Filter("DateFormat", func(vm *hvue.VM, value *js.Object, args ...*js.Object) interface{} {
			return date.DateString(value.String())
		}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type FoaSiteTableModel struct {
	*js.Object

	Foasiteinfos []*fm.FoaSiteInfo `js:"foasiteinfos"`
	User         *fm.User          `js:"user"`
	Filter       string            `js:"filter"`

	VM *hvue.VM `js:"VM"`
}

func NewFoaSiteTableModel(vm *hvue.VM) *FoaSiteTableModel {
	rtm := &FoaSiteTableModel{Object: tools.O()}
	rtm.Foasiteinfos = []*fm.FoaSiteInfo{}
	rtm.User = fm.NewUser()
	rtm.Filter = ""
	rtm.VM = vm
	return rtm
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Event Related Methods

func (fstm *FoaSiteTableModel) SetSelectedFoaSite(fsi *fm.FoaSiteInfo) {
	fstm.OpenFoaSite(fsi.Id)
}

func (fstm *FoaSiteTableModel) AddFoaSite(vm *hvue.VM) {
	vm.Emit("new_foasite")
}

func (fstm *FoaSiteTableModel) AttachmentUrl(id int) string {
	return "/api/foasites/" + strconv.Itoa(id) + "/attach"
}

func (fstm *FoaSiteTableModel) OpenFoaSite(id int) {
	js.Global.Get("window").Call("open", "foasite.html?fsid="+strconv.Itoa(id))
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Formatting Related Methods

func (fstm *FoaSiteTableModel) TableRowClassName(rowInfo *js.Object) string {
	fsi := &fm.FoaSiteInfo{Object: rowInfo.Get("row")}
	return fs.FoaSiteRowClassName(fsi.Status)
}

func (fstm *FoaSiteTableModel) HeaderCellStyle() string {
	return "background: #a1e6e6;"
}

func (fstm *FoaSiteTableModel) FormatDate(r, c *js.Object, d string) string {
	return date.DateString(d)
}

func (fstm *FoaSiteTableModel) FormatStatus(r, c *js.Object, d string) string {
	return fs.FoaSiteStatusLabel(d)
}

func (fstm *FoaSiteTableModel) SortStatus(a, b *fm.FoaSiteInfo) int {
	la := fs.FoaSiteStatusLabel(a.Status)
	lb := fs.FoaSiteStatusLabel(b.Status)
	if la < lb {
		return -1
	}
	if la == lb {
		return 0
	}
	return 1
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Column Filtering Related Methods

func (fstm *FoaSiteTableModel) FilterHandler(value string, p *js.Object, col *js.Object) bool {
	prop := col.Get("property").String()
	return p.Get(prop).String() == value
}

func (fstm *FoaSiteTableModel) FilterList(vm *hvue.VM, prop string) []*elements.ValText {
	fstm = &FoaSiteTableModel{Object: vm.Object}
	count := map[string]int{}
	attribs := []string{}

	var translate func(string) string
	switch prop {
	case "Status":
		translate = func(val string) string {
			return fs.FoaSiteStatusLabel(val)
		}
	default:
		translate = func(val string) string { return val }
	}

	for _, psi := range fstm.Foasiteinfos {
		attrib := psi.Object.Get(prop).String()
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

func (fstm *FoaSiteTableModel) FilteredStatusValue() []string {
	res := []string{
		//foaconst.PsStatusNew,
		//foaconst.PsStatusInProgress,
		//foaconst.PsStatusBlocked,
		//foaconst.PsStatusCancelled,
		//foaconst.PsStatusDone,
	}
	return res
}
