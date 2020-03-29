package polemap

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/leafletmap"
	"github.com/lpuig/ewin/doe/website/frontend/model/polesite"
	"github.com/lpuig/ewin/doe/website/frontend/tools/leaflet"
	"sort"
	"strconv"
	"strings"
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
	Poles            []*polesite.Pole               `js:"poles"`
	PoleMarkers      map[string][]*PoleMarker       `js:"PoleMarkers"`
	PoleMarkersGroup map[string]*leaflet.LayerGroup `js:"PoleMarkersGroup"`
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
	pm.initPoleMarkersGroups()
	//pm.PoleOverlays = make(map[string]*leaflet.Layer)
	return pm
}

func (pm *PoleMap) initPoleMarkersGroups() {
	pm.PoleMarkers = make(map[string][]*PoleMarker)
	pm.PoleMarkersGroup = make(map[string]*leaflet.LayerGroup)
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
	poleMarker.Map = pm
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
	poleMarker.On("mouseover", func(o *js.Object) {
		poleMarker := PoleMarkerFromJS(o.Get("target"))
		poleMarker.Map.ShowPoleInfo(poleMarker)
	})
	poleMarker.On("mouseout", func(o *js.Object) {
		poleMarker := PoleMarkerFromJS(o.Get("target"))
		poleMarker.Map.HidePoleInfo()
	})

	return poleMarker
}

type namedLayers struct {
	Name   string
	Layers []*leaflet.Layer
}

func (pm *PoleMap) updatePoleMarkersGroups() {
	pms := make(map[string][]*PoleMarker)
	polesLayer := make(map[string][]*leaflet.Layer)

	// create group
	for _, pole := range pm.Poles {
		poleMarker := pm.NewPoleMarker(pole)
		pms[pole.State] = append(pms[pole.State], poleMarker)
		polesLayer[pole.State] = append(polesLayer[pole.State], &poleMarker.Layer)
	}
	pm.PoleMarkers = pms

	// set controlGroup and add LayerGroup to map
	pmg := make(map[string]*leaflet.LayerGroup)
	slGroup := []namedLayers{}
	for name, layers := range polesLayer {
		slGroup = append(slGroup, namedLayers{Name: name, Layers: layers})
	}
	sort.Slice(slGroup, func(i, j int) bool {
		return slGroup[i].Name < slGroup[j].Name
	})
	for _, nl := range slGroup {
		nlg := leaflet.NewLayerGroup(nl.Layers)
		pmg[nl.Name] = nlg
		nlg.AddTo(pm.LeafletMap.Map)
		groupName := polesite.PoleStateLabel(nl.Name) + " (" + strconv.Itoa(len(nl.Layers)) + ")"
		pm.LeafletMap.ControlLayers.AddOverlay(&nlg.Layer, groupName)
	}
	pm.PoleMarkersGroup = pmg
}

// AddPoles creates PoleMarkers and adds them to Map
func (pm *PoleMap) AddPoles(poles []*polesite.Pole) {
	pm.Poles = poles
	pm.initPoleMarkersGroups()
	pm.updatePoleMarkersGroups()
}

// RefreshPoles removes previously existing PoleMarkers from Map, and adds given Poles to map
func (pm *PoleMap) RefreshPoles(poles []*polesite.Pole) {
	// remove Poles groups from map and controlLayer
	for _, group := range pm.PoleMarkersGroup {
		pm.LeafletMap.ControlLayers.RemoveLayer(&group.Layer)
		group.Remove()
	}
	pm.AddPoles(poles)
}

// GetPoleMarkerById returns the polemarker bounded to given pole's Id (nil if not found)
func (pm *PoleMap) GetPoleMarkerById(id int) *PoleMarker {
	for _, group := range pm.PoleMarkers {
		for _, poleMarker := range group {
			if poleMarker.Pole.Id == id {
				return poleMarker
			}
		}
	}
	return nil
}

// RefreshGroup refreshs all group of polemarkers
func (pm *PoleMap) GetPoleMarkers() []*PoleMarker {
	res := []*PoleMarker{}
	for _, pms := range pm.PoleMarkers {
		res = append(res, pms...)
	}
	return res
}

// RefreshGroup refreshs all group of polemarkers
func (pm *PoleMap) RefreshGroup() {
	for _, group := range pm.PoleMarkersGroup {
		group.Refresh()
	}
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
	return pm.LeafletMap.Map.GetCenter()
}

func (pm *PoleMap) FitBounds(min, max *leaflet.LatLng) {
	pm.LeafletMap.Map.FitBounds(min, max)
}

func (pm *PoleMap) ShowPoleInfo(poleMarker *PoleMarker) {
	pole := poleMarker.Pole
	html := "<h4>" + pole.GetTitle() + "</h4>"
	html += `<p class="right">` + pole.Material + " " + strconv.Itoa(pole.Height) + "m" + "</p>"
	html += strings.Join(pole.Product, ", ") + "<br />"

	if pole.Comment != "" {
		html += `<p><span class="title">Commentaire: </span>` + pole.Comment + "</p>"
	}
	if pole.DictInfo != "" {
		html += `<p><span class="title">Info DICT: </span>` + pole.DictInfo + "</p>"
	}
	pm.ControlInfo.Update(html)
}

func (pm *PoleMap) HidePoleInfo() {
	pm.ControlInfo.Update("")
}
