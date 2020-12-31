package ripteamproductivitychart

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	rs "github.com/lpuig/ewin/doe/website/frontend/model/ripsite"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

//////////////////////////////////////////////////////////////////////////////////////////////
// Component Template

const template string = `<div class="statchart" ref="container" :style="SetStyle()"></div>`

//////////////////////////////////////////////////////////////////////////////////////////////
// Component Methods

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("ripteam-productivity-chart", componentOption()...)
}

func componentOption() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.Props("stats", "colors", "heigth", "mode"),
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

const (
	ModeProductivity string = "productivity"
	ModeProgress     string = "progress"
)

//////////////////////////////////////////////////////////////////////////////////////////////
// Model Methods

type TeamProductivityChart struct {
	*js.Object
	VM     *hvue.VM      `js:"VM"`
	Stats  *rs.TeamStats `js:"stats"`
	Colors SiteColorMap  `js:"colors"`
	Heigth string        `js:"heigth"`
	Mode   string        `js:"mode"`
}

func NewTeamProductivityChart(vm *hvue.VM) *TeamProductivityChart {
	tpc := &TeamProductivityChart{Object: tools.O()}
	tpc.VM = vm
	tpc.Stats = rs.NewTeamStats()
	tpc.Colors = nil
	tpc.Heigth = "250px"
	tpc.Mode = ModeProductivity // or ModeProgress
	return tpc
}

func (tpc *TeamProductivityChart) SetStyle(vm *hvue.VM) string {
	tpc = &TeamProductivityChart{Object: vm.Object}
	return "width:100%; height:" + tpc.Heigth + ";"
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
		"yAxis": tpc.getAxis(),
		"legend": js.M{
			"enabled": false,
			//"layout":        "vertical",
			//"align":         "right",
			//"verticalAlign": "top",
		},
		"tooltip": js.M{
			//"shared":      true,
			//"pointFormat": "<b>{series.name}:</b> {point.y:.1f}",
			"footerFormat":  "Total: <b>{point.total:.1f}</b>",
			"valueDecimals": 1,
		},
		"plotOptions": js.M{
			"series": js.M{
				//"allowPointSelect": true,
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
		"series": tpc.getSeries(),
	}
	js.Global.Get("Highcharts").Call("chart", tpc.VM.Refs("container"), chartdesc)
}

func (tpc *TeamProductivityChart) getAxis() []js.M {
	res := []js.M{}
	if !(len(tpc.Stats.Values["Work"]) == 0 &&
		len(tpc.Stats.Values["RoleMeanWork"]) == 0 &&
		len(tpc.Stats.Values["GlobalMeanWork"]) == 0) {
		res = append(res, js.M{
			"labels": js.M{
				"format": "{value}",
			},
			"title": js.M{
				"text": "Points",
			},
		})
	}
	if len(tpc.Stats.Values["NbActorsWork"]) > 0 {
		res = append(res, js.M{
			"labels": js.M{
				"format": "{value}",
			},
			"title": js.M{
				"text": "Nb. Acteurs",
			},
			"opposite": true,
		})
	}
	if len(tpc.Stats.Values["Price"]) > 0 {
		res = append(res, js.M{
			"labels": js.M{
				"format": "{value}€",
			},
			"title": js.M{
				"text": "Revenus",
			},
			"opposite": true,
		})
	}
	return res
}

func (tpc *TeamProductivityChart) getSeries() []js.M {
	res := []js.M{}
	switch tpc.Mode {
	case ModeProgress:
		if len(tpc.Stats.Values["Work"]) > 0 {
			res = append(res, newSerie("line", "Solid", "Travail", "work", "", " Pts.",
				tpc.Colors["Work"], 0,
				0,
				tpc.Stats.Values["Work"])...)
		}
		if len(tpc.Stats.Values["WorkTarget"]) > 0 {
			res = append(res, newSerie("line", "Dash", "Travail", "worktarget", "", " Pts.",
				tpc.Colors["Work"], 0,
				0,
				tpc.Stats.Values["WorkTarget"])...)
		}
		if len(tpc.Stats.Values["Price"]) > 0 {
			res = append(res, newSerie("line", "Solid", "Euros", "price", "", " €",
				tpc.Colors["Price"], 1,
				0,
				tpc.Stats.Values["Price"])...)
		}
		if len(tpc.Stats.Values["PriceTarget"]) > 0 {
			res = append(res, newSerie("line", "Dash", "Euros", "pricetarget", "", " €",
				tpc.Colors["Price"], 1,
				0,
				tpc.Stats.Values["PriceTarget"])...)
		}
	default:
		if len(tpc.Stats.Values["Work"]) > 0 {
			res = append(res, newSerie("column", "Solid", "Travail", "work", "", " Pts.",
				tpc.Colors["Work"], 0,
				0.05,
				tpc.Stats.Values["Work"])...)
		}
		if len(tpc.Stats.Values["Price"]) > 0 {
			res = append(res, newSerie("column", "Solid", "Euros", "price", "", " €",
				tpc.Colors["Price"], 1,
				0,
				tpc.Stats.Values["Price"])...)
		}
	}
	if len(tpc.Stats.Values["RoleMeanWork"]) > 0 {
		res = append(res, newSerie("column", "Solid", "Moyenne Rôle", "rolemean", "", " Pts.",
			tpc.Colors["Price"], 0,
			0.15,
			tpc.Stats.Values["RoleMeanWork"])...)
	}
	if len(tpc.Stats.Values["GlobalMeanWork"]) > 0 {
		res = append(res, newSerie("column", "Solid", "Moyenne globale", "globalmean", "", " Pts.",
			tpc.Colors["Price"], 0,
			0.15,
			tpc.Stats.Values["GlobalMeanWork"])...)
	}
	if len(tpc.Stats.Values["NbActorsWork"]) > 0 {
		res = append(res, newSerie("column", "Solid", "Nb Acteurs", "nbactors", "", "",
			tpc.Colors["Work"], 1,
			0.3,
			tpc.Stats.Values["NbActorsWork"])...)
	}
	return res
}
