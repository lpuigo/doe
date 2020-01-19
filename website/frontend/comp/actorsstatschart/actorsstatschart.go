package actorsstatschart

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	rs "github.com/lpuig/ewin/doe/website/frontend/model/ripsite"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/date"
)

//////////////////////////////////////////////////////////////////////////////////////////////
// Component Template

const template string = `
<div 
		class="statchart" 
		ref="container" 
		:style="SetStyle()"
></div>
`

//////////////////////////////////////////////////////////////////////////////////////////////
// Component Methods

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("actors-stats-chart", componentOption()...)
}

func componentOption() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.Props("stats", "colors"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewActorsStatsChart(vm)
		}),
		hvue.MethodsOf(&ActorsStatsChart{}),
		hvue.Mounted(func(vm *hvue.VM) {
			asc := &ActorsStatsChart{Object: vm.Object}
			asc.SetChart(asc.Stats)
		}),
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Tools Functions

//////////////////////////////////////////////////////////////////////////////////////////////
// Model Methods

type ActorsStatsChart struct {
	*js.Object
	VM    *hvue.VM      `js:"VM"`
	Stats *rs.TeamStats `js:"stats"`
}

func NewActorsStatsChart(vm *hvue.VM) *ActorsStatsChart {
	asc := &ActorsStatsChart{Object: tools.O()}
	asc.VM = vm
	asc.Stats = rs.NewTeamStats()
	return asc
}

func ActorsStatsChartFromJS(o *js.Object) *ActorsStatsChart {
	return &ActorsStatsChart{Object: o}
}

func (asc *ActorsStatsChart) SetStyle() string {
	return "width:100%; height:450px;"
}

func (asc *ActorsStatsChart) SetChart(ts *rs.TeamStats) {
	asc.Stats = ts
	chartdesc := js.M{
		"chart": js.M{
			"backgroundColor": "#F7F7F7",
			"type":            "column",
		},
		"title": js.M{
			"text": nil,
		},
		"xAxis": js.M{
			//	"type": "datetime",
			//	"dateTimeLabelFormats": js.M{
			//		"day": "%e %b",
			//	},
			"categories": date.ConvertDates(ts.Dates),
			//"tickPixelInterval" : 400,
		},
		"yAxis": asc.getAxis(),
		"legend": js.M{
			"enabled": false,
			//"layout":        "vertical",
			//"align":         "right",
			//"verticalAlign": "top",
		},
		"tooltip": js.M{
			//"shared":      true,
			//"pointFormat": "<b>{series.name}:</b> {point.y:.1f}",
			"valueDecimals": 1,
		},
		"plotOptions": js.M{
			"series": js.M{
				"allowPointSelect": true,
				//"pointStart":       startDate,
				//"pointInterval":    7 * 24 * 3600 * 1000, // one week
				"marker":    js.M{"enabled": false},
				"animation": false,
				//"grouping":  true,
				"stacking":     "normal",
				"groupPadding": 0.08,
			},
			"column": js.M{
				//"pointPadding": 0.1,
				//"borderWidth":  0,
				//"groupPadding": 0,
				"borderRadius": 2,
				"shadow":       false,
			},
		},
		"series": asc.getSeries(),
	}
	js.Global.Get("Highcharts").Call("chart", asc.VM.Refs("container"), chartdesc)
}

func (asc *ActorsStatsChart) getAxis() []js.M {
	res := []js.M{}
	if len(asc.Stats.Values["employees"]) > 0 {
		res = append(res, js.M{
			"labels": js.M{
				"format": "{value}",
			},
			"title": js.M{
				"text": "Acteurs",
			},
			"opposite": false,
		})
	}
	return res
}

func (asc *ActorsStatsChart) getSeries() []js.M {
	res := []js.M{}
	res = append(res, newSerie("line", "Acteurs", "employees", "", "",
		"#888888", 0,
		0.05,
		asc.Stats.Values["employees"])...)
	res = append(res, newSerie("line", "Présents", "acting", "", "",
		"#67C23A", 0,
		0.05,
		asc.Stats.Values["acting"])...)
	return res
}

func newSerie(style, name, stack, tooltipPrefix, tooltipSuffix, color string, yAxis int, pointpadding float64, data map[string][]float64) []js.M {
	//s := &serie{Object: tools.O()}
	//s.Name = name
	//s.Color = color
	//s.Data = data
	res := []js.M{}
	for key, value := range data {
		res = append(res, js.M{
			//"name":         name + " " + key,
			"name":         key,
			"type":         style,
			"color":        color,
			"yAxis":        yAxis,
			"data":         value,
			"pointPadding": pointpadding,
			"stack":        stack,
			"tooltip": js.M{
				"valuePrefix": tooltipPrefix,
				"valueSuffix": tooltipSuffix,
			},
		})
	}
	return res
}
