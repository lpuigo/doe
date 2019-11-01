package polemap

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/leafletmap"
	"github.com/lpuig/ewin/doe/website/frontend/model/polesite"
	"github.com/lpuig/ewin/doe/website/frontend/tools/leaflet"
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
		hvue.MethodsOf(&PoleMap{}),
		hvue.Mounted(func(vm *hvue.VM) {
			pm := newPoleMap(vm)
			pm.VM = vm
			if len(pm.Poles) > 0 {
				pm.AddPoles(pm.Poles)
			}
		}),
		//hvue.BeforeUpdate(func(vm *hvue.VM) {
		//	pm := PoleMapFromJS(vm.Object)
		//	print("polemap beforeUpdate", pm.Poles)
		//	//pm.AddPoles(pm.Poles, "attrib Poteaux")
		//}),
		//hvue.Updated(func(vm *hvue.VM) {
		//	pm := PoleMapFromJS(vm.Object)
		//	print("polemap Updated", pm.Poles)
		//}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type PoleMap struct {
	leafletmap.LeafletMap
	Poles            []*polesite.Pole    `js:"poles"`
	PoleMarkers      []*PoleMarker       `js:"PoleMarkers"`
	PoleMarkersGroup *leaflet.LayerGroup `js:"PoleMarkersGroup"`
	//PoleOverlays map[string]*leaflet.Layer `js:"PoleOverlays"`
}

func PoleMapFromJS(obj *js.Object) *PoleMap {
	return &PoleMap{LeafletMap: leafletmap.LeafletMap{Object: obj}}
}

func newPoleMap(vm *hvue.VM) *PoleMap {
	pm := PoleMapFromJS(vm.Object)
	pm.LeafletMap.VM = vm
	pm.LeafletMap.Init()
	pm.Poles = nil
	pm.PoleMarkers = nil
	pm.PoleMarkersGroup = nil
	//pm.PoleOverlays = make(map[string]*leaflet.Layer)
	return pm
}

// NewPoleMarker creates and returns a new configured PoleMarker for given pole
func (pm *PoleMap) NewPoleMarker(pole *polesite.Pole) *PoleMarker {
	dio := leaflet.DefaultDivIconOptions()
	ico := leaflet.NewDivIcon(dio)
	mOption := leaflet.DefaultMarkerOption()
	mOption.Icon = &ico.Icon
	mOption.Opacity = 0.5
	mOption.Title = pole.GetTitle()

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

	return poleMarker
}

// AddPoles creates PoleMarkers and adds them to Map
func (pm *PoleMap) AddPoles(poles []*polesite.Pole) {
	pm.Poles = poles
	pms := make([]*PoleMarker, len(poles))
	polesLayer := make([]*leaflet.Layer, len(poles))

	for i, pole := range poles {
		poleMarker := pm.NewPoleMarker(pole)
		pms[i] = poleMarker
		polesLayer[i] = &poleMarker.Layer
	}
	pm.PoleMarkers = pms
	pm.PoleMarkersGroup = leaflet.NewLayerGroup(polesLayer)
	pm.PoleMarkersGroup.AddTo(pm.LeafletMap.Map)
	pm.LeafletMap.ControlLayers.AddOverlay(&pm.PoleMarkersGroup.Layer, "Poteaux")
}

// RefreshPoles removes previously existing PoleMarkers from Map, and adds given Poles to map
func (pm *PoleMap) RefreshPoles(poles []*polesite.Pole) {
	pm.LeafletMap.ControlLayers.RemoveLayer(&pm.PoleMarkersGroup.Layer)
	//pm.Map.RemoveLayer(&pm.PoleMarkersGroup.Layer)
	pm.PoleMarkersGroup.Remove()

	pm.AddPoles(poles)
}

// GetPoleMarkerById returns the polemarker bounded to given pole's Id (nil if not found)
func (pm *PoleMap) GetPoleMarkerById(id int) *PoleMarker {
	//layers := pm.PoleMarkersGroup.GetLayers()
	//for i := 0; i < layers.Length(); i++ {
	//	poleMarker := PoleMarkerFromJS(layers.Index(i))
	//	if poleMarker.Pole.Id == id {
	//		return poleMarker
	//	}
	//}

	for _, poleMarker := range pm.PoleMarkers {
		if poleMarker.Pole.Id == id {
			return poleMarker
		}
	}
	return nil
}

func (pm *PoleMap) CenterOnPoles() {
	clat, clong, minlat, minlong, maxlat, maxlong := polesite.GetCenterAndBounds(pm.Poles)
	pm.LeafletMap.Map.Stop()
	pm.LeafletMap.Map.SetView(leaflet.NewLatLng(clat, clong), 12)
	pm.LeafletMap.Map.FitBounds(leaflet.NewLatLng(minlat, minlong), leaflet.NewLatLng(maxlat, maxlong))
}

func (pm *PoleMap) CenterOn(lat, long float64, zoom int) {
	pm.LeafletMap.Map.Stop()
	pm.LeafletMap.Map.SetView(leaflet.NewLatLng(lat, long), zoom)
}

func (pm *PoleMap) GetCenter() *leaflet.LatLng {
	return pm.Map.GetCenter()
}

func (pm *PoleMap) FitBounds(min, max *leaflet.LatLng) {
	pm.LeafletMap.Map.FitBounds(min, max)
}
