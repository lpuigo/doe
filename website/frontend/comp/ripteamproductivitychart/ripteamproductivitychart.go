package ripteamproductivitychart

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	rs "github.com/lpuig/ewin/doe/website/frontend/model/ripsite"
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
	return hvue.Component("ripteam-productivity-chart", componentOption()...)
}

func componentOption() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.Props("stats", "colors"),
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
	VM     *hvue.VM      `js:"VM"`
	Stats  *rs.TeamStats `js:"stats"`
	Colors SiteColorMap  `js:"colors"`
}

func NewTeamProductivityChart(vm *hvue.VM) *TeamProductivityChart {
	tpc := &TeamProductivityChart{Object: tools.O()}
	tpc.VM = vm
	tpc.Stats = rs.NewTeamStats()
	tpc.Colors = nil
	return tpc
}

func (tpc *TeamProductivityChart) SetStyle() string {
	return "width:100%; height:250px;"
}

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
			//	"type": "datetime",
			//	"dateTimeLabelFormats": js.M{
			//		"day": "%e %b",
			//	},
			"categories": ts.Dates,
			//"tickPixelInterval" : 400,
		},
		"yAxis": js.S{
			js.M{
				"labels": js.M{
					"format": "{value} h",
				},
				"title": js.M{
					"text": "Heures",
				},
			},
			js.M{
				"labels": js.M{
					"format": "{value}€",
				},
				"title": js.M{
					"text": "Revenus",
				},
				"opposite": true,
			},
		},
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
				"stacking": "normal",
			},
			"column": js.M{
				//"pointPadding": 0.1,
				//"borderWidth":  0,
				//"groupPadding": 0,
				"borderRadius": 2,
				"shadow":       false,
			},
		},
		"series": tpc.getSeries(),
	}
	js.Global.Get("Highcharts").Call("chart", tpc.VM.Refs("container"), chartdesc)
}

func (tpc *TeamProductivityChart) getSeries() []js.M {
	res := []js.M{}
	res = append(res, newSerie("column", "Travail", "work", "", " h", tpc.Colors["Work"], 0, 0.2, tpc.Stats.Values["Work"])...)
	if len(tpc.Stats.Values["Price"]) > 0 {
		res = append(res, newSerie("column", "€", "price", "", " €", tpc.Colors["Price"], 1, 0, tpc.Stats.Values["Price"])...)
	}
	return res
}
