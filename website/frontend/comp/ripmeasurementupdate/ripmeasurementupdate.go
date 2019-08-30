package ripmeasurementupdate

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/ripstateupdate"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	fmrip "github.com/lpuig/ewin/doe/website/frontend/model/ripsite"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements/message"
	"github.com/lpuig/ewin/doe/website/frontend/tools/json"
	"strconv"
	"strings"
)

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("rip-measurement-update", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		ripstateupdate.RegisterComponent(),
		hvue.Template(template),
		hvue.Props("value", "user", "filter"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewRipMeasurementUpdateModel(vm)
		}),
		hvue.MethodsOf(&RipMeasurementUpdateModel{}),
		hvue.Computed("filteredMeasurements", func(vm *hvue.VM) interface{} {
			rmum := RipMeasurementUpdateModelFromJS(vm.Object)
			return rmum.GetFilteredMeasurements()
		}),
		hvue.Computed("NbTotFiber", func(vm *hvue.VM) interface{} {
			rmum := RipMeasurementUpdateModelFromJS(vm.Object)
			rmum.CalcStats()
			return rmum.MeasurementSummary.NbFiber
		}),
		//hvue.Filter("pct", func(vm *hvue.VM, value *js.Object, args ...*js.Object) interface{} {
		//}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type RipMeasurementUpdateModel struct {
	*js.Object

	Ripsite *fmrip.Ripsite `js:"value"`
	//ReferenceRipsite *fmrip.Ripsite `js:"refRipsite"`
	User   *fm.User `js:"user"`
	Filter string   `js:"filter"`

	UploadVisible      bool               `js:"uploadVisible"`
	MeasurementTeam    string             `js:"measurementTeam"`
	MeasurementSummary *fmrip.Measurement `js:"measurementSummary"`

	VM *hvue.VM `js:"VM"`
}

func NewRipMeasurementUpdateModel(vm *hvue.VM) *RipMeasurementUpdateModel {
	rmum := &RipMeasurementUpdateModel{Object: tools.O()}
	rmum.VM = vm
	rmum.Ripsite = fmrip.NewRisite()
	//rmum.ReferenceWorksite = nil
	rmum.User = nil
	rmum.Filter = ""
	rmum.UploadVisible = false
	rmum.MeasurementTeam = ""
	rmum.MeasurementSummary = fmrip.NewMeasurement()
	return rmum
}

func RipMeasurementUpdateModelFromJS(o *js.Object) *RipMeasurementUpdateModel {
	return &RipMeasurementUpdateModel{Object: o}
}

func (rmum *RipMeasurementUpdateModel) CalcStats() {
	rmum.MeasurementSummary = fmrip.NewMeasurement()
	for _, meas := range rmum.Ripsite.Measurements {
		rmum.MeasurementSummary.NbFiber += meas.NbFiber
		rmum.MeasurementSummary.NbOK += meas.NbOK
		rmum.MeasurementSummary.NbWarn1 += meas.NbWarn1
		rmum.MeasurementSummary.NbWarn2 += meas.NbWarn2
		rmum.MeasurementSummary.NbKO += meas.NbKO
	}
}

func (rmum *RipMeasurementUpdateModel) GetPct(vm *hvue.VM, value int) string {
	rmum = RipMeasurementUpdateModelFromJS(vm.Object)
	if rmum.MeasurementSummary.NbFiber == 0 {
		return " - "
	}
	return strconv.FormatFloat(100.0*float64(value)/float64(rmum.MeasurementSummary.NbFiber), 'f', 1, 64) + "%"
}

func (rmum *RipMeasurementUpdateModel) GetFilteredMeasurements() []*fmrip.Measurement {
	if rmum.Filter == "" {
		return rmum.Ripsite.Measurements
	}
	res := []*fmrip.Measurement{}
	filter := strings.ToLower(rmum.Filter)
	for _, meas := range rmum.Ripsite.Measurements {
		if strings.Contains(strings.ToLower(json.Stringify(meas)), filter) {
			res = append(res, meas)
		}
	}
	return res
}

func (rmum *RipMeasurementUpdateModel) TableRowClassName(rowInfo *js.Object) string {
	meas := &fmrip.Measurement{Object: rowInfo.Get("row")}
	return rmum.MeasurementClassName(meas)
}

func (rmum *RipMeasurementUpdateModel) MeasurementClassName(meas *fmrip.Measurement) string {
	return meas.State.GetRowStyle()
}

func (rmum *RipMeasurementUpdateModel) GetNode(vm *hvue.VM, name string) *fmrip.Node {
	rmum = RipMeasurementUpdateModelFromJS(vm.Object)
	return rmum.Ripsite.Nodes[name]
}

func (rmum *RipMeasurementUpdateModel) GetDestNodeDist(vm *hvue.VM, meas *fmrip.Measurement) int {
	rmum = RipMeasurementUpdateModelFromJS(vm.Object)
	node := rmum.Ripsite.Nodes[meas.DestNodeName]
	return node.DistFromPm
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Upload Methods

func (rmum *RipMeasurementUpdateModel) GetTeams(vm *hvue.VM) []*elements.ValueLabel {
	rmum = RipMeasurementUpdateModelFromJS(vm.Object)
	return rmum.User.GetTeamValueLabelsFor(rmum.Ripsite.Client)
}

func (rmum *RipMeasurementUpdateModel) UploadError(vm *hvue.VM, err, file *js.Object) {
	rmum.UploadVisible = false
	message.ErrorStr(vm, err.String(), false)
}

func (rmum *RipMeasurementUpdateModel) UploadSuccess(vm *hvue.VM, response, file *js.Object) {
	rmum = RipMeasurementUpdateModelFromJS(vm.Object)
	rmum.UploadVisible = false

	for _, meas := range rmum.Ripsite.Measurements {
		omr := response.Get(meas.DestNodeName)
		if omr == js.Undefined {
			continue
		}
		mr := &fmrip.MeasurementReport{Object: omr}
		meas.UpdateWith(mr, rmum.MeasurementTeam)
	}
	rmum.CalcStats()
}

func (rmum *RipMeasurementUpdateModel) BeforeUpload(vm *hvue.VM, file *js.Object) bool {
	if file.Get("type").String() != "application/x-zip-compressed" {
		message.ErrorStr(vm, "Le fichier '"+file.Get("name").String()+"' n'est pas une archive ZIP", false)
		return false
	}
	return true
}
