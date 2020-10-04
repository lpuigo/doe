package leafletmap

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/leaflet"
)

const template string = `
<div id="LeafLetMap" style="height: 100%"></div>
`

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Registration

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("leaflet-map", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.MethodsOf(&LeafletMap{}),
		//hvue.Computed("progressPct", func(vm *hvue.VM) interface{} {
		//	wspb := &WorksiteProgressBarModel{Object: vm.Object}
		//	return wspb.ProgressPct()
		//}),
		hvue.Mounted(func(vm *hvue.VM) {
			llm := NewLeafletMap(vm)
			llm.Init()
		}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type LeafletMap struct {
	*js.Object

	Map           *leaflet.Map           `js:"Map"`
	ControlLayers *leaflet.ControlLayers `js:"ControlLayers"`
	ControlInfo   *leaflet.ControlInfo   `js:"ControlInfo"`

	VM *hvue.VM `js:"VM"`
}

func NewLeafletMap(vm *hvue.VM) *LeafletMap {
	llm := &LeafletMap{Object: tools.O()}
	llm.VM = vm
	llm.ControlLayers = nil
	llm.ControlInfo = nil
	return llm
}

func LeafletMapFromJS(o *js.Object) *LeafletMap {
	return &LeafletMap{Object: o}
}

func (llm *LeafletMap) Init() {
	mapOption := leaflet.DefaultMapOptions()

	llm.Map = leaflet.NewMap("LeafLetMap", mapOption)
	osmlayer := leaflet.OSMTileLayer()
	mplayer := leaflet.MapBoxTileLayer("mapbox/outdoors-v11")
	satlayer := leaflet.MapBoxTileLayer("mapbox/satellite-streets-v11")

	baseMaps := js.M{
		"Plan":        osmlayer,
		"Plan MapBox": mplayer,
		"Satellite":   satlayer,
	}
	osmlayer.AddTo(llm.Map)

	llm.ControlLayers = leaflet.NewControlLayers(baseMaps, js.M{})
	llm.ControlLayers.AddTo(llm.Map)

	leaflet.NewControlScale().AddTo(llm.Map)

	llm.ControlInfo = leaflet.NewControlInfo()
	llm.ControlInfo.AddTo(llm.Map)

	llm.SetView(leaflet.NewLatLng(48, 5), 6)

}

func (llm *LeafletMap) SetView(center *leaflet.LatLng, zoom int) {
	llm.Map.SetView(center, zoom)
}

func (llm *LeafletMap) FitBounds(bound1, bound2 *leaflet.LatLng) {
	llm.Map.FitBounds(bound1, bound2)
}
