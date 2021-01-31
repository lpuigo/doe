package polemap

import (
	"sort"
	"strconv"
	"strings"

	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/leafletmap"
	"github.com/lpuig/ewin/doe/website/frontend/model/polesite"
	"github.com/lpuig/ewin/doe/website/frontend/model/polesite/poleconst"
	"github.com/lpuig/ewin/doe/website/frontend/tools/date"
	"github.com/lpuig/ewin/doe/website/frontend/tools/leaflet"
)

const template string = `
<!--<div id="LeafLetMap" style="height: 100%"></div>-->
<div style="height: 100%">
	<div v-if="NbSelected > 0"  class="map-button">
		<el-button type="warning" plain  @click="HandleClearSelected()" size="mini">Déselectionner {{NbSelected}} appui(s)</el-button>
	</div>
	<div id="LeafLetMap" style="height: 100%"></div>
<!--	<div v-if="SelectedPoleMarkers.length > 0" class="map-button">-->
<!--		<el-button type="primary" plain  @click="ClearSelected" size="mini">Déselectionner {{SelectedPoleMarkers.length}} appui(s)</el-button>-->
<!--	</div>-->
</div>
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
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewPoleMap(vm)
		}),
		hvue.Mounted(func(vm *hvue.VM) {
			pm := PoleMapFromJS(vm.Object)
			pm.Init()
		}),
		hvue.Computed("NbSelected", func(vm *hvue.VM) interface{} {
			pm := PoleMapFromJS(vm.Object)
			return len(pm.SelectedPoleMarkers)
		}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type PoleMap struct {
	*leafletmap.LeafletMap
	Categorizer         *Categorizer                   `js:"Categorizer"`
	PoleLine            *PoleLine                      `js:"PoleLine"`
	Poles               []*polesite.Pole               `js:"poles"`
	PoleMarkers         []*PoleMarker                  `js:"PoleMarkers"`
	PoleMarkersByState  map[string][]*PoleMarker       `js:"PoleMarkersByState"`
	PoleMarkersGroup    map[string]*leaflet.LayerGroup `js:"PoleMarkersGroup"`
	SelectedPoleMarkers []*PoleMarker                  `js:"SelectedPoleMarkers"`
}

func PoleMapFromJS(obj *js.Object) *PoleMap {
	return &PoleMap{LeafletMap: leafletmap.LeafletMapFromJS(obj)}
}

func NewPoleMap(vm *hvue.VM) *PoleMap {
	pm := &PoleMap{LeafletMap: leafletmap.NewLeafletMap(vm)}
	//pm.LeafletMap.VM = vm
	//pm.LeafletMap.Init()
	pm.Categorizer = NewCategorizer()
	pm.PoleLine = nil
	pm.Poles = nil
	pm.clearPoleMarkersGroups()
	return pm
}

func (pm *PoleMap) Init() {
	pm.LeafletMap.Init()
	pm.PoleLine = NewPoleLine(pm)
}

// clearPoleMarkersGroups clears all PoleMarker containers
func (pm *PoleMap) clearPoleMarkersGroups() {
	pm.PoleMarkers = []*PoleMarker{}
	pm.PoleMarkersByState = make(map[string][]*PoleMarker)
	pm.PoleMarkersGroup = make(map[string]*leaflet.LayerGroup)
	pm.SelectedPoleMarkers = []*PoleMarker{}
}

// NewPoleMarker creates and returns a new configured PoleMarker for given pole
func (pm *PoleMap) NewPoleMarker(pole *polesite.Pole) *PoleMarker {
	dio := leaflet.DefaultDivIconOptions()
	ico := leaflet.NewDivIcon(dio)
	mOption := leaflet.DefaultMarkerOption()
	mOption.Icon = &ico.Icon
	mOption.Opacity = 0.5
	//mOption.Title = pole.GetTitle()

	poleMarker := NewPoleMarker(mOption, pole)
	poleMarker.Map = pm
	//poleMarker.BindPopup(pole.Ref)
	poleMarker.UpdateFromState()
	poleMarker.On("click", func(event *js.Object) {
		poleMarker := PoleMarkerFromJS(event.Get("sourceTarget"))
		if event.Get("originalEvent").Get("shiftKey").Bool() {
			poleMarker.SwitchSelection()
			return
		}
		poleMarker.Map.VM.Emit("marker-click", poleMarker, event)
	})
	poleMarker.On("contextmenu", func(event *js.Object) {
		poleMarker := PoleMarkerFromJS(event.Get("sourceTarget"))
		poleMarker.Map.PoleLine.SetPivot(poleMarker)
		poleMarker.Map.PoleLine.Draw()
	})
	poleMarker.On("move", func(event *js.Object) {
		poleMarker := PoleMarkerFromJS(event.Get("target"))
		poleMarker.Map.RefreshPivotLine(event)
	})
	poleMarker.On("dragend", func(event *js.Object) {
		poleMarker := PoleMarkerFromJS(event.Get("target"))
		if poleMarker.Map.PoleLine.IsDrawn {
			poleMarker.SetLatLng(poleMarker.Map.PoleLine.RoundedPos)
			poleMarker.Refresh()
			return
		}
		poleMarker.Pole.Lat, poleMarker.Pole.Long = poleMarker.GetLatLong().ToFloats()
	})
	poleMarker.On("mouseover", func(event *js.Object) {
		poleMarker := PoleMarkerFromJS(event.Get("target"))
		poleMarker.Map.ShowPoleInfo(poleMarker)
	})
	poleMarker.On("mouseout", func(event *js.Object) {
		poleMarker := PoleMarkerFromJS(event.Get("target"))
		poleMarker.Map.HidePoleInfo()
	})

	return poleMarker
}

type namedLayers struct {
	Name   string
	Layers []*leaflet.Layer
}

// createPoleMarkers resets PoleMarkersGroups and creates new PoleMarkers for attached Poles
//
// Deleted Poles are ignored
func (pm *PoleMap) createPoleMarkers() {
	pm.clearPoleMarkersGroups()
	pms := []*PoleMarker{}

	for _, pole := range pm.Poles {
		if pole.State == poleconst.StateDeleted {
			// skip deleted poles
			continue
		}
		pms = append(pms, pm.NewPoleMarker(pole))
	}
	pm.PoleMarkers = pms
}

// setPoleMarkersGroups resets PoleMarkersByState and PoleMarkersGroup according to attached PoleMarkers
func (pm *PoleMap) setPoleMarkersGroups() {
	pms := make(map[string][]*PoleMarker)
	polesLayer := make(map[string][]*leaflet.Layer)

	// create group
	for _, poleMarker := range pm.PoleMarkers {
		groupName := pm.Categorizer.GroupName(poleMarker)
		pms[groupName] = append(pms[groupName], poleMarker)
		polesLayer[groupName] = append(polesLayer[groupName], &poleMarker.Layer)
	}
	pm.PoleMarkersByState = pms

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
		groupName := pm.Categorizer.GroupLabel(nl.Name) + " (" + strconv.Itoa(len(nl.Layers)) + ")"
		pm.LeafletMap.ControlLayers.AddOverlay(&nlg.Layer, groupName)
	}
	pm.PoleMarkersGroup = pmg
}

// SetPoles creates PoleMarkers and adds them to Map (pre-existing markers and groups are deleted)
func (pm *PoleMap) SetPoles(poles []*polesite.Pole) {
	pm.Poles = poles
	pm.createPoleMarkers()
	pm.setPoleMarkersGroups()
}

// RefreshPoles removes existing PoleMarkers from Map, and attach given Poles to map (PoleMarkers and Groups are rebuilt)
func (pm *PoleMap) RefreshPoles(poles []*polesite.Pole) {
	// remove Poles groups from map and controlLayer
	for _, group := range pm.PoleMarkersGroup {
		pm.LeafletMap.ControlLayers.RemoveLayer(&group.Layer)
		group.Remove()
	}
	pm.SetPoles(poles)
}

// RefreshPoleMarkersGroups updates PoleMarkersGroups according to existing PoleMarkers (PoleMarkers Groups are reseted)
func (pm *PoleMap) RefreshPoleMarkersGroups() {
	// remove Poles groups from map and controlLayer
	for _, group := range pm.PoleMarkersGroup {
		pm.LeafletMap.ControlLayers.RemoveLayer(&group.Layer)
		group.Remove()
	}
	pm.setPoleMarkersGroups()
}

// GetPoleMarkerById returns the polemarker bounded to given pole's Id (nil if not found)
func (pm *PoleMap) GetPoleMarkerById(id int) *PoleMarker {
	for _, poleMarker := range pm.PoleMarkers {
		if poleMarker.Pole.Id == id {
			return poleMarker
		}
	}
	return nil
}

// RefreshGroup refreshs all group of polemarkers
func (pm *PoleMap) GetPoleMarkers() []*PoleMarker {
	return pm.PoleMarkers[:]
}

// RefreshGroup refreshs all group of polemarkers
func (pm *PoleMap) RefreshGroup() {
	for _, group := range pm.PoleMarkersGroup {
		group.Refresh()
	}
}

func (pm *PoleMap) CenterOnPoleMarkers(pms []*PoleMarker) {
	if len(pms) == 0 {
		return
	}
	poles := make([]*polesite.Pole, len(pms))
	for i, poleMarker := range pms {
		poles[i] = poleMarker.Pole
	}
	_, _, minlat, minlong, maxlat, maxlong := polesite.GetCenterAndBounds(poles)
	pm.LeafletMap.Map.Stop()
	pm.LeafletMap.Map.FitBounds(leaflet.NewLatLng(minlat, minlong), leaflet.NewLatLng(maxlat, maxlong))
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
	html := `<h4><span class="header">` + polesite.PoleStateLabel(pole.State) + ": </span>" + pole.GetTitle() + "</h4>"
	html += `<p class="right">` + pole.Material + " " + strconv.Itoa(pole.Height) + "m" + "</p>"
	html += strings.Join(pole.Product, ", ") + "<br />"

	if pole.Comment != "" {
		html += `<p><span class="title">Commentaire: </span>` + pole.Comment + "</p>"
	}
	if pole.DictInfo != "" {
		html += `<p><span class="title">Info DICT: </span>` + pole.DictInfo + "</p>"
	}
	if pole.IsToDo() {
		if pole.State > poleconst.StateDictToDo {
			html += `<br /><p><span class="title">DICT: </span> du ` + date.DateString(pole.DictDate) + " au " + date.DateString(date.After(pole.DictDate, poleconst.DictValidityDuration)) + "</p>"
		}
		if pole.State > poleconst.StateDaToDo {
			daInfo := ""
			if pole.State == poleconst.StateDaExpected {
				daInfo = "demandé le " + date.DateString(pole.DaQueryDate)
			} else {
				daInfo = "du " + date.DateString(pole.DaStartDate) + " au " + date.DateString(pole.DaEndDate)
			}
			html += `<p><span class="title">DA: </span>` + daInfo + "</p>"
		}
		if !(!pole.IsDone() && !pole.IsAttachment()) {
			html += `<br /><p><span class="title">Réalisé:</span> le ` + date.DateString(pole.Date) + "</p>"
			if pole.IsAttachment() {
				html += `<p><span class="title">Facturé:</span> le ` + date.DateString(pole.AttachmentDate) + "</p>"
			}
		}
	}
	pm.ControlInfo.Update(html)
}

func (pm *PoleMap) HidePoleInfo() {
	pm.ControlInfo.Update("")
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Selected Poles related methods

func (pm *PoleMap) NbSelected() int {
	return len(pm.SelectedPoleMarkers)
}

func (pm *PoleMap) SelectedIndex(p *PoleMarker) int {
	for i, poleMarker := range pm.SelectedPoleMarkers {
		if poleMarker.Pole.Id == p.Pole.Id {
			return i
		}
	}
	return -1
}

func (pm *PoleMap) AddSelected(p *PoleMarker) {
	if pm.SelectedIndex(p) >= 0 {
		return
	}
	//pm.SelectedPoleMarkers = append(pm.SelectedPoleMarkers, p)
	pm.Get("SelectedPoleMarkers").Call("push", p)
}

func (pm *PoleMap) RemoveSelected(p *PoleMarker) {
	index := pm.SelectedIndex(p)
	if index < 0 {
		return
	}
	pm.Get("SelectedPoleMarkers").Call("splice", index, 1)
}

func (pm *PoleMap) HandleClearSelected(vm *hvue.VM) {
	pm = PoleMapFromJS(vm.Object)
	pm.ClearSelected()
}

func (pm *PoleMap) ClearSelected() {
	for _, poleMarker := range pm.SelectedPoleMarkers {
		poleMarker.Deselect()
	}
	pm.SelectedPoleMarkers = []*PoleMarker{}
}

// SelectByFilter resets current SelectionList, and select all Pole mathcing given filter
func (pm *PoleMap) SelectByFilter(filter func(*PoleMarker) bool) {
	pm.ClearSelected()
	spms := []*PoleMarker{}
	for _, poleMarker := range pm.GetPoleMarkers() {
		if !poleMarker.Edited && filter(poleMarker) {
			poleMarker.Select()
			spms = append(spms, poleMarker)
		}
	}
	pm.SelectedPoleMarkers = spms
}

// ApplyOnSelected applies action on all selected PoleMarkers
func (pm *PoleMap) ApplyOnSelected(action func(*PoleMarker)) {
	for _, poleMarker := range pm.SelectedPoleMarkers {
		action(poleMarker)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// PoleLine related methods

func (pm *PoleMap) DisablePoleLine() {
	pm.PoleLine.Reinit()
}

func (pm *PoleMap) RefreshPivotLine(moveEvent *js.Object) {
	latlng := leaflet.LatLngFromJS(moveEvent.Get("latlng"))
	pm.PoleLine.RefreshTargetPos(latlng)
}
