package ripmeasurementupdate

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/ripstateupdate"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	fmrip "github.com/lpuig/ewin/doe/website/frontend/model/ripsite"
	"github.com/lpuig/ewin/doe/website/frontend/model/ripsite/ripconst"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements/message"
	"strconv"
	"strings"
	"time"
)

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("rip-measurement-update", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		ripstateupdate.RegisterComponent(),
		hvue.Template(template),
		hvue.Props("value", "user", "filter", "filtertype"),
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

	Ripsite    *fmrip.Ripsite `js:"value"`
	User       *fm.User       `js:"user"`
	Filter     string         `js:"filter"`
	FilterType string         `js:"filtertype"`

	UploadVisible       bool                 `js:"uploadVisible"`
	AnalysisVisible     bool                 `js:"analysisVisible"`
	MeasurementActors   []string             `js:"measurementActors"`
	MeasurementSummary  *fmrip.Measurement   `js:"measurementSummary"`
	MeasurementAnalysis []*fmrip.Measurement `js:"measurementAnalysis"`

	SizeLimit int      `js:"SizeLimit"`
	VM        *hvue.VM `js:"VM"`
}

func NewRipMeasurementUpdateModel(vm *hvue.VM) *RipMeasurementUpdateModel {
	rmum := &RipMeasurementUpdateModel{Object: tools.O()}
	rmum.VM = vm
	rmum.Ripsite = fmrip.NewRisite()
	rmum.User = nil
	rmum.Filter = ""
	rmum.FilterType = ripconst.FilterValueAll
	rmum.UploadVisible = false
	rmum.AnalysisVisible = false
	rmum.MeasurementActors = []string{}
	rmum.MeasurementSummary = fmrip.NewMeasurement()
	rmum.MeasurementAnalysis = []*fmrip.Measurement{}
	rmum.SetSizeLimit()
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
	if rmum.FilterType == ripconst.FilterValueAll && rmum.Filter == "" {
		return rmum.GetSizeLimitedResult(rmum.Ripsite.Measurements)
	}
	res := []*fmrip.Measurement{}
	expected := strings.ToUpper(rmum.Filter)
	filter := func(p *fmrip.Measurement) bool {
		sis := p.SearchString(rmum.FilterType)
		if sis == "" {
			return false
		}
		return strings.Contains(strings.ToUpper(sis), expected)
	}

	for _, meas := range rmum.Ripsite.Measurements {
		if filter(meas) {
			res = append(res, meas)
		}
	}
	return rmum.GetSizeLimitedResult(res)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Size Related Methods

const (
	sizeLimitDefault int = 15
	sizeLimitTimer       = 300
)

func (rmum *RipMeasurementUpdateModel) GetSizeLimitedResult(res []*fmrip.Measurement) []*fmrip.Measurement {
	if len(res) == rmum.SizeLimit {
		return res
	}
	if len(res) > sizeLimitDefault {
		rmum.ResetSizeLimit(len(res))
		return res[:sizeLimitDefault]
	}
	return res
}

func (rmum *RipMeasurementUpdateModel) SetSizeLimit() {
	rmum.SizeLimit = -1
}

func (rmum *RipMeasurementUpdateModel) ResetSizeLimit(size int) {
	go func() {
		time.Sleep(sizeLimitTimer * time.Millisecond)
		rmum.SizeLimit = size
	}()
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Component Methods

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

func (rmum *RipMeasurementUpdateModel) GetActors(vm *hvue.VM) []*elements.ValueLabelDisabled {
	rmum = RipMeasurementUpdateModelFromJS(vm.Object)
	client := rmum.User.GetClientByName(rmum.Ripsite.Client)
	if client == nil {
		return nil
	}

	res := []*elements.ValueLabelDisabled{}
	for _, actor := range client.Actors {
		res = append(res, actor.GetElementsValueLabelDisabled())
	}
	return res
}

func (rmum *RipMeasurementUpdateModel) UploadError(vm *hvue.VM, err, file *js.Object) {
	rmum.UploadVisible = false
	message.ErrorStr(vm, err.String(), false)
}

func (rmum *RipMeasurementUpdateModel) UploadSuccess(vm *hvue.VM, response, file *js.Object) {
	rmum = RipMeasurementUpdateModelFromJS(vm.Object)
	rmum.UploadVisible = false

	nbOk, nbKo := 0, 0
	for _, meas := range rmum.Ripsite.Measurements {
		omr := response.Get(meas.DestNodeName)
		if omr == js.Undefined {
			nbKo++
			continue
		}
		nbOk++
		mr := &fmrip.MeasurementReport{Object: omr}
		meas.UpdateWith(mr, rmum.MeasurementActors)
	}
	message.SuccesStr(vm, "Mise à jour des mesures :\n"+strconv.Itoa(nbOk)+" mesures mises à jour\n"+strconv.Itoa(nbKo)+" mesures inchangées")
	rmum.CalcStats()
}

func (rmum *RipMeasurementUpdateModel) BeforeUpload(vm *hvue.VM, file *js.Object) bool {
	if file.Get("type").String() != "application/x-zip-compressed" {
		message.ErrorStr(vm, "Le fichier '"+file.Get("name").String()+"' n'est pas une archive ZIP", false)
		return false
	}
	return true
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Analysis Methods

func (rmum *RipMeasurementUpdateModel) Analyse(vm *hvue.VM) {
	rmum = RipMeasurementUpdateModelFromJS(vm.Object)

	measAnalysis := make(map[string]*fmrip.Measurement)
	for _, meas := range rmum.Ripsite.Measurements {
		if !meas.State.IsMeasured() {
			continue
		}
		if tools.Empty(meas.State.Comment) {
			continue
		}
		for _, warning := range meas.ParseComment() {
			nearestNode := rmum.getClosestNode(meas, warning.Dist)
			nodedefect, found := measAnalysis[nearestNode]
			if !found {
				nodedefect = fmrip.NewMeasurement()
				nodedefect.DestNodeName = nearestNode
				measAnalysis[nearestNode] = nodedefect
			}
			switch warning.WarnLvl {
			case "KO Splice":
				nodedefect.NbKO++
			case "Warn2":
				nodedefect.NbWarn2++
			case "Warn1":
				nodedefect.NbWarn1++
			case "KO Connector":
				nodedefect.NbKO++
			}
			destNodeFound := false
			for _, destNode := range nodedefect.NodeNames {
				if destNode == meas.DestNodeName {
					destNodeFound = true
					break
				}
			}
			if !destNodeFound {
				nodedefect.NodeNames = append(nodedefect.NodeNames, meas.DestNodeName)
			}
		}
	}
	res := []*fmrip.Measurement{}
	for _, meas := range measAnalysis {
		res = append(res, meas)
	}
	rmum.MeasurementAnalysis = res
	rmum.Get("measurementAnalysis").Call("sort", func(a, b *fmrip.Measurement) int {
		if a.NbKO > b.NbKO {
			return -1
		}
		if a.NbKO < b.NbKO {
			return 1
		}
		if a.NbWarn2 > b.NbWarn2 {
			return -1
		}
		if a.NbWarn2 < b.NbWarn2 {
			return 1
		}
		if a.NbWarn1 > b.NbWarn1 {
			return -1
		}
		if a.NbWarn1 < b.NbWarn1 {
			return 1
		}
		return 0
	})
	rmum.AnalysisVisible = true
}

func (rmum *RipMeasurementUpdateModel) getClosestNode(meas *fmrip.Measurement, dist float64) string {
	actDist := -1.0
	nearestNode := ""
	for _, node := range meas.NodeNames {
		newDist := float64(rmum.Ripsite.Nodes[node].DistFromPm) - dist
		if newDist < 0 {
			newDist = -newDist
		}
		if !(nearestNode != "" && newDist >= actDist) {
			actDist = newDist
			nearestNode = node
		}
	}
	return nearestNode
}

func (rmum *RipMeasurementUpdateModel) GetAnalysisDestNodes(vm *hvue.VM, meas *fmrip.Measurement) string {
	return strings.Join(meas.NodeNames, ", ")
}
