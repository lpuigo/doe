package polemap

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/mappoc/frontend/leaflet"
	"github.com/lpuig/ewin/mappoc/frontend/mappoc/comp/leafletmap"
	"github.com/lpuig/ewin/mappoc/frontend/mappoc/model"
)

const template string = `
<div id="LeafLetMap" style="height: 100%"></div>
`

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Registration

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("pole-map", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.Props("poles"),
		hvue.MethodsOf(&PoleMap{}),
		//hvue.Computed("progressPct", func(vm *hvue.VM) interface{} {
		//	wspb := &WorksiteProgressBarModel{Object: vm.Object}
		//	return wspb.ProgressPct()
		//}),
		hvue.Mounted(func(vm *hvue.VM) {
			pm := newPoleMap(vm)
			pm.VM = vm
			if len(pm.Poles) > 0 {
				pm.AddPoles(pm.Poles, "Poteaux")
			}
		}),
		hvue.BeforeUpdate(func(vm *hvue.VM) {
			pm := PoleMapFromJS(vm.Object)
			print("polemap beforeUpdate", pm.Poles)
			//pm.AddPoles(pm.Poles, "attrib Poteaux")
		}),
		hvue.Updated(func(vm *hvue.VM) {
			pm := PoleMapFromJS(vm.Object)
			print("polemap Updated", pm.Poles)
		}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type PoleMap struct {
	leafletmap.LeafletMap
	Poles []*model.Pole `js:"poles"`
	//PoleOverlays map[string]*leaflet.Layer `js:"PoleOverlays"`
}

func PoleMapFromJS(obj *js.Object) *PoleMap {
	return &PoleMap{LeafletMap: leafletmap.LeafletMap{Object: obj}}
}

func newPoleMap(vm *hvue.VM) *PoleMap {
	pm := PoleMapFromJS(vm.Object)
	pm.LeafletMap.VM = vm
	pm.LeafletMap.Init()
	pm.LeafletMap.SetView(leaflet.NewLatLng(48, 6), 5)
	pm.Poles = nil
	//pm.PoleOverlays = make(map[string]*leaflet.Layer)
	return pm
}

func (pm *PoleMap) AddPoles(poles []*model.Pole, name string) *leaflet.LayerGroup {
	pm.Poles = poles
	polesLayer := []*leaflet.Layer{}

	//if layer, exist := pm.PoleOverlays[name]; exist {
	//	pm.LeafletMap.ControlLayers.RemoveLayer(layer)
	//	delete(pm.PoleOverlays, name)
	//}

	for _, pole := range poles {
		dio := leaflet.DefaultDivIconOptions()
		ico := leaflet.NewDivIcon(dio)
		mOption := leaflet.DefaultMarkerOption()
		mOption.Icon = &ico.Icon
		mOption.Opacity = 0.5
		mOption.Title = pole.Ref

		//poleMarker := leaflet.NewMarker(pole.Lat, pole.Long, mOption)
		poleMarker := NewPoleMarker(mOption, pole)
		//poleMarker.BindPopup(pole.Ref)
		poleMarker.UpdateFromState()
		poleMarker.On("click", func(o *js.Object) {
			poleMarker := PoleMarkerFromJS(o.Get("sourceTarget"))
			pm.VM.Emit("marker-click", poleMarker, o)
		})
		poleMarker.On("dragend", func(o *js.Object) {
			poleMarker := PoleMarkerFromJS(o.Get("target"))
			poleMarker.Pole.Lat, poleMarker.Pole.Long = poleMarker.GetLatLong().ToFloats()
		})
		polesLayer = append(polesLayer, &poleMarker.Layer)
	}

	polesGroup := leaflet.NewLayerGroup(polesLayer)
	polesGroup.AddTo(pm.LeafletMap.Map)
	pm.LeafletMap.ControlLayers.AddOverlay(&polesGroup.Layer, name)
	//pm.PoleOverlays[name] = &polesGroup.Layer

	pm.CenterOnPoles()
	return polesGroup
}

func (pm *PoleMap) CenterOnPoles() {
	clat, clong, minlat, minlong, maxlat, maxlong := model.GetCenterAndBounds(pm.Poles)
	pm.LeafletMap.Map.SetView(leaflet.NewLatLng(clat, clong), 12)
	pm.LeafletMap.Map.FitBounds(leaflet.NewLatLng(minlat, minlong), leaflet.NewLatLng(maxlat, maxlong))
}

func (pm *PoleMap) CenterOn(lat, long float64, zoom int) {
	pm.LeafletMap.Map.SetView(leaflet.NewLatLng(lat, long), zoom)
}
