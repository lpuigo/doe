package teamproductivitychart

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/model/worksite"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
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
			tpc.setColumnChart()
		}),
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Model Methods

type TeamProductivityChart struct {
	*js.Object
	VM    *hvue.VM            `js:"VM"`
	Stats *worksite.TeamStats `js:"stats"`
}

func NewTeamProductivityChart(vm *hvue.VM) *TeamProductivityChart {
	tpc := &TeamProductivityChart{Object: tools.O()}
	tpc.VM = vm
	tpc.Stats = worksite.NewTeamStats()
	return tpc
}

func (tpc *TeamProductivityChart) SetStyle() string {
	return "width:100%; height:250px;"
}

//func (tpc *TeamProductivityChart) setChart() {
//	ts := tpc.Stats
//	startDate := date.JSDate(ts.StartDate)
//
//	chartdesc := js.M{
//		"chart": js.M{
//			"backgroundColor": "#F7F7F7",
//			"type":            "line",
//		},
//		"title": js.M{
//			"text": nil,
//		},
//		//"xAxis": js.M{
//		//	"categories": ts.Dates,
//		//	"tickPixelInterval" : 400,
//		//},
//		"xAxis": js.M{
//			"type": "datetime",
//			"dateTimeLabelFormats": js.M{
//				"day": "%e %b",
//			},
//		},
//		"yAxis": js.M{
//			"title": js.M{
//				"text": "Nbs.",
//			},
//		},
//		"legend": js.M{
//			"layout":        "vertical",
//			"align":         "right",
//			"verticalAlign": "top",
//		},
//		"plotOptions": js.M{
//			"series": js.M{
//				"allowPointSelect": false,
//				"pointStart":       startDate,
//				"pointInterval":    7 * 24 * 3600 * 1000, // one week
//				"marker":           js.M{"enabled": false},
//				"animation":        false,
//			},
//		},
//		"series": tpc.getSeries(),
//	}
//	js.Global.Get("Highcharts").Call("chart", tpc.VM.Refs("container"), chartdesc)
//}

func (tpc *TeamProductivityChart) setColumnChart() {
	ts := tpc.Stats
	//startDate := date.JSDate(ts.StartDate)

	chartdesc := js.M{
		"chart": js.M{
			"backgroundColor": "#F7F7F7",
			"type":            "column",
		},
		"title": js.M{
			"text": nil,
		},
		"xAxis": js.M{
			"categories": ts.Dates,
			//"tickPixelInterval" : 400,
		},
		//"xAxis": js.M{
		//	"type": "datetime",
		//	"dateTimeLabelFormats": js.M{
		//		"day": "%e %b",
		//	},
		//},
		"yAxis": js.M{
			"title": js.M{
				"text": "Nb. ELs",
			},
		},
		"legend": js.M{
			"layout":        "vertical",
			"align":         "right",
			"verticalAlign": "top",
		},
		"tooltip": js.M{
			"shared": true,
		},
		"plotOptions": js.M{
			"series": js.M{
				"allowPointSelect": false,
				//"pointStart":       startDate,
				//"pointInterval":    7 * 24 * 3600 * 1000, // one week
				"marker":    js.M{"enabled": false},
				"animation": false,
			},
			"column": js.M{
				//"pointPadding": 0.1,
				//"borderWidth":  0,
				//"groupPadding": 0,
				"borderRadius": 2,
				"shadow":       false,
				"grouping":     true,
			},
		},
		"series": tpc.getSeries(),
	}
	js.Global.Get("Highcharts").Call("chart", tpc.VM.Refs("container"), chartdesc)
}

func (tpc *TeamProductivityChart) getSeries() []interface{} {
	res := []interface{}{}
	switch {
	case len(tpc.Stats.Values["Work"]) > 0:
		if len(tpc.Stats.Values["Price"]) > 0 {
			res = append(res, newSerie("column", "€", "#51A825", 0, tpc.Stats.Values["Price"]))
		}
		res = append(res, newSerie("line", "Travail", "#389eff", 0, tpc.Stats.Values["Work"]))
	default:
		res = append(res, newSerie("column", "Installés", "#51A825", 0, tpc.Stats.Values["Installed"]))
		res = append(res, newSerie("column", "Mesurés", "#29d1cb", 0, tpc.Stats.Values["Measured"]))
		res = append(res, newSerie("column", "Bloqués", "#cc2020", 0.2, tpc.Stats.Values["Blocked"]))
		res = append(res, newSerie("line", "DOE", "#389eff", 0, tpc.Stats.Values["DOE"]))
	}
	return res
}
