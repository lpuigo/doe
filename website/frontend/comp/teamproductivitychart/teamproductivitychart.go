package teamproductivitychart

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/dates"
)

type TeamProductivityChart struct {
	*js.Object
	VM    *hvue.VM      `js:"VM"`
	Stats *fm.TeamStats `js:"stats"`
}

func NewTeamProductivityChart(vm *hvue.VM) *TeamProductivityChart {
	tpc := &TeamProductivityChart{Object: tools.O()}
	tpc.VM = vm
	tpc.Stats = fm.NewTeamStats()
	return tpc
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Component Template

const template string = `
<div 
		class="issuechart" 
		ref="container" 
		:style="SetStyle()"
></div>
`

//////////////////////////////////////////////////////////////////////////////////////////////
// Component Methods

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("team-productivity-chart", componentOption()...)
}

func componentOption() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.Props("stats"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewTeamProductivityChart(vm)
		}),
		hvue.MethodsOf(&TeamProductivityChart{}),
		hvue.Mounted(func(vm *hvue.VM) {
			tpc := &TeamProductivityChart{Object: vm.Object}
			tpc.setChart()
		}),
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Model Methods

func (tpc *TeamProductivityChart) SetStyle() string {
	return "width:100%; height:250px;"
}

func (tpc *TeamProductivityChart) setChart() {
	ts := tpc.Stats
	startDate := date.JSDate(ts.StartDate)

	chartdesc := js.M{
		"chart": js.M{
			"backgroundColor": "#F7F7F7",
			"type":            "line",
		},
		"title": js.M{
			"text": nil,
		},
		//"xAxis": js.M{
		//	"categories": ts.Dates,
		//	"tickPixelInterval" : 400,
		//},
		"xAxis": js.M{
			"type": "datetime",
			"dateTimeLabelFormats": js.M{
				"day": "%e %b",
			},
		},
		"yAxis": js.M{
			"title": js.M{
				"text": "Days",
			},
		},
		"legend": js.M{
			"layout":        "vertical",
			"align":         "right",
			"verticalAlign": "top",
		},
		"plotOptions": js.M{
			"series": js.M{
				"allowPointSelect": false,
				"pointStart":       startDate,
				"pointInterval":    7 * 24 * 3600 * 1000, // one week
				"marker":           js.M{"enabled": false},
				"animation":        false,
			},
		},
		"series": js.S{
			js.M{
				"name":      "Nb. EL",
				"color":     "#51A825",
				"lineWidth": 5,
				"data":      ts.NbEls,
			},
		},
	}
	js.Global.Get("Highcharts").Call("chart", tpc.VM.Refs("container"), chartdesc)
}
