package ripjunctionupdate

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/rippullingdistinfo"
	"github.com/lpuig/ewin/doe/website/frontend/comp/ripstateupdate"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	fmrip "github.com/lpuig/ewin/doe/website/frontend/model/ripsite"
	"github.com/lpuig/ewin/doe/website/frontend/model/ripsite/ripconst"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/date"
	"strconv"
	"strings"
)

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("rip-junction-update", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		rippullingdistinfo.RegisterComponent(),
		ripstateupdate.RegisterComponent(),
		hvue.Template(template),
		hvue.Props("value", "user", "filter", "filtertype"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewRipJunctionUpdateModel(vm)
		}),
		hvue.MethodsOf(&RipJunctionUpdateModel{}),
		hvue.Computed("filteredJunctions", func(vm *hvue.VM) interface{} {
			rpum := RipJunctionUpdateModelFromJS(vm.Object)
			return rpum.GetFilteredJunctions()
		}),
		//hvue.Filter("FormatTronconRef", func(vm *hvue.VM, value *js.Object, args ...*js.Object) interface{} {
		//	rpum := RipJunctionUpdateModelFromJS(vm.Object)
		//	t := &fm.Troncon{Object: value}
		//	return rpum.GetFormatTronconRef(t)
		//}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type RipJunctionUpdateModel struct {
	*js.Object

	Ripsite *fmrip.Ripsite `js:"value"`
	//ReferenceRipsite *fmrip.Ripsite `js:"refRipsite"`
	User       *fm.User `js:"user"`
	Filter     string   `js:"filter"`
	FilterType string   `js:"filtertype"`

	VM *hvue.VM `js:"VM"`
}

func NewRipJunctionUpdateModel(vm *hvue.VM) *RipJunctionUpdateModel {
	rpum := &RipJunctionUpdateModel{Object: tools.O()}
	rpum.VM = vm
	rpum.Ripsite = fmrip.NewRisite()
	//rpum.ReferenceWorksite = nil
	rpum.User = nil
	rpum.Filter = ""
	rpum.FilterType = ripconst.FilterValueAll
	return rpum
}

func RipJunctionUpdateModelFromJS(o *js.Object) *RipJunctionUpdateModel {
	return &RipJunctionUpdateModel{Object: o}
}

func (rjum *RipJunctionUpdateModel) GetFilteredJunctions() []*fmrip.Junction {
	if rjum.FilterType == ripconst.FilterValueAll && rjum.Filter == "" {
		return rjum.Ripsite.Junctions
	}
	res := []*fmrip.Junction{}
	expected := strings.ToUpper(rjum.Filter)
	filter := func(p *fmrip.Junction) bool {
		sis := p.SearchString(rjum.FilterType)
		if sis == "" {
			return false
		}
		return strings.Contains(strings.ToUpper(sis), expected)
	}

	for _, junction := range rjum.Ripsite.Junctions {
		if filter(junction) {
			res = append(res, junction)
		}
	}
	return res
}

func (rjum *RipJunctionUpdateModel) TableRowClassName(rowInfo *js.Object) string {
	junction := &fmrip.Junction{Object: rowInfo.Get("row")}
	return junction.State.GetRowStyle()
}

func (rjum *RipJunctionUpdateModel) GetNbFiber(vm *hvue.VM, junction *fmrip.Junction) int {
	rjum = RipJunctionUpdateModelFromJS(vm.Object)
	return junction.GetNbFiber()
}

func (rjum *RipJunctionUpdateModel) GetNode(vm *hvue.VM, junction *fmrip.Junction) *fmrip.Node {
	rjum = RipJunctionUpdateModelFromJS(vm.Object)
	return rjum.Ripsite.Nodes[junction.NodeName]
}

func (rjum *RipJunctionUpdateModel) GetNodeDesc(vm *hvue.VM, junction *fmrip.Junction) string {
	rjum = RipJunctionUpdateModelFromJS(vm.Object)
	node := rjum.Ripsite.Nodes[junction.NodeName]
	return node.Name + " - " + node.Ref
}

func (rjum *RipJunctionUpdateModel) GetNodeType(vm *hvue.VM, junction *fmrip.Junction) string {
	rjum = RipJunctionUpdateModelFromJS(vm.Object)
	node := rjum.Ripsite.Nodes[junction.NodeName]
	return node.Type + ": " + node.BoxType
}

func (rjum *RipJunctionUpdateModel) GetTronconDesc(vm *hvue.VM, junction *fmrip.Junction) string {
	rjum = RipJunctionUpdateModelFromJS(vm.Object)
	node := rjum.Ripsite.Nodes[junction.NodeName]
	troncon := rjum.Ripsite.Troncons[node.TronconInName]
	if troncon == nil {
		return node.TronconInName
	}
	return troncon.Name + " ( " + strconv.Itoa(troncon.Size) + "FO)"
}

func (rjum *RipJunctionUpdateModel) GetActors(vm *hvue.VM, junction *fmrip.Junction) string {
	rjum = RipJunctionUpdateModelFromJS(vm.Object)
	client := rjum.User.GetClientByName(rjum.Ripsite.Client)
	if client == nil {
		return ""
	}

	res := []string{}
	for _, actId := range junction.State.Actors {
		actor := client.GetActorBy(actId)
		if actor == nil {
			continue
		}
		res = append(res, actor.GetRef())
	}
	return strings.Join(res, "\n")
}

func (rjum *RipJunctionUpdateModel) FormatDate(r, c *js.Object, d string) string {
	return date.DateString(d)
}

func (rjum *RipJunctionUpdateModel) FormatStatus(r, c *js.Object, d string) string {
	return fmrip.GetStatusLabel(d)
}

func (rjum *RipJunctionUpdateModel) GetNodeAttr(vm *hvue.VM, col string) func(a *fmrip.Junction) string {
	rjum = RipJunctionUpdateModelFromJS(vm.Object)
	return func(a *fmrip.Junction) string {
		node := rjum.Ripsite.Nodes[a.NodeName]
		return node.Get(col).String()
	}
}
