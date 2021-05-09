package polemap

import (
	"strconv"

	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/model/polesite"
	"github.com/lpuig/ewin/doe/website/frontend/model/polesite/poleconst"
	"github.com/lpuig/ewin/doe/website/frontend/tools/leaflet"
)

type PoleMarker struct {
	leaflet.Marker
	Pole      *polesite.Pole `js:"Pole"`
	Map       *PoleMap       `js:"Map"`
	Draggable bool           `js:"Draggable"`
	Selected  bool           `js:"Selected"`
	Edited    bool           `js:"Edited"`
	Class     string         `js:"Class"`
	Html      string         `js:"Html"`
}

func PoleMarkerFromJS(obj *js.Object) *PoleMarker {
	return &PoleMarker{Marker: *leaflet.MarkerFromJs(obj)}
}

func NewPoleMarker(option *leaflet.MarkerOptions, pole *polesite.Pole) *PoleMarker {
	np := &PoleMarker{Marker: *leaflet.NewMarker(pole.Lat, pole.Long, option)}
	np.Pole = pole
	np.Map = nil
	np.Draggable = false
	np.Selected = false
	np.Edited = false
	np.Class = ""
	np.Html = pmHtmlPin
	return np
}

func DefaultPoleMarker() *PoleMarker {
	opt := leaflet.DefaultMarkerOption()
	np := &PoleMarker{Marker: *leaflet.NewMarker(0.0, 0.0, opt)}
	np.Pole = polesite.NewPole()
	return np
}

// StartEditMode set the receiver as edited, updates its look and refreshes it
func (pm *PoleMarker) StartEditMode(drag bool) {
	if drag {
		pm.SetDraggable(true)
	}
	pm.Edited = true
	pm.UpdateDivIconHtml()
	pm.Map.PoleLine.SetTarget(pm)
	pm.Refresh()
}

// EndEditMode set the receiver as not edited, updates its look and refreshes it
func (pm *PoleMarker) EndEditMode(refresh bool) {
	pm.Edited = false
	pm.SetDraggable(false)
	pm.UpdateDivIconHtml()
	pm.Map.PoleLine.Reinit()
	if refresh {
		pm.Refresh()
	}
}

// SwitchSelection switched receiver selected status, and updates its html and map Selection list
func (pm *PoleMarker) SwitchSelection() {
	pm.Selected = !pm.Selected
	if pm.Selected {
		pm.Map.AddSelected(pm)
	} else {
		pm.Map.RemoveSelected(pm)
	}
	pm.UpdateDivIconHtml()
	pm.Refresh()
}

// Deselect set the receiver unselected, and updates its look and refreshes it (Map selection list is not updated)
func (pm *PoleMarker) Deselect() {
	pm.Selected = false
	pm.UpdateDivIconHtml()
	pm.Refresh()
}

// Select set the receiver as selected, updates its look and refreshes it (Map selection list is not updated)
func (pm *PoleMarker) Select() {
	pm.Selected = true
	pm.UpdateDivIconHtml()
	pm.Refresh()
}

// SetDraggable sets the marker as draggable and refreshes it
func (pm *PoleMarker) SetDraggable(drag bool) {
	pm.Draggable = drag
	pm.Marker.SetDraggable(pm.Draggable)
	pm.UpdateDivIconHtml()
	pm.Refresh()
}

// Refresh refreshes the look of the PoleMarker receiver on its map (State and groups are not updated)
func (pm *PoleMarker) Refresh() {
	pm.Remove()
	pm.AddTo(pm.Map.Map)
}

// FullRefreshState refreshes the look of the PoleMarker receiver on its map and reinit state groups.
//
// As PoleMarkers are discarded and created back during the process, pointer on newly created PoleMarker bound to the same receiver pole is returned.
func (pm *PoleMarker) FullRefreshState() *PoleMarker {
	pm.Map.RefreshPoles(pm.Map.Poles)
	return pm.Map.GetPoleMarkerById(pm.Pole.Id)
}

// UpdateRefreshGroup updates the polemarker receiver look, update the layer groups and refreshes all PoleMarkers on map
func (pm *PoleMarker) UpdateRefreshGroup() {
	pm.UpdateFromState()
	pm.Map.RefreshPoleMarkersGroups()
}

const (
	pmHtmlPin           string = `<i class="fas fa-map-pin fa-fw fa-3x pole-marker"></i>`
	pmHtmlReplace       string = `<i class="fas fa-arrows-alt-v fa-fw fa-3x pole-marker"></i>`
	pmHtmlTrickyReplace string = `<i class="fas fa-level-down-alt fa-fw fa-3x pole-marker"></i>`
	pmHtmlCreation      string = `<i class="fas fa-long-arrow-alt-down fa-fw fa-3x pole-marker"></i>`
	pmHtmlReplenish     string = `<i class="fas fa-angle-double-down fa-fw fa-3x pole-marker"></i>`
	pmHtmlPlain         string = `<i class="fas fa-map-marker fa-fw fa-3x pole-marker"></i>`
	pmHtmlBolt          string = `<i class="fas fa-bolt fa-fw fa-3x pole-marker"></i>`
	pmHtmlExclam        string = `<i class="fas fa-exclamation fa-fw fa-3x pole-marker"></i>`
	pmHtmlHole          string = `<i class="fas fa-map-marker-alt fa-fw fa-3x pole-marker"></i>`
	pmHtmlOutline       string = `<i class="el-icon-location-outline pole-marker" style="font-size: 3.3em"></i>`

	pmHtmlShadow       string = `<div class="pole_marker_shadow"></div>`
	pmHtmlShadowEdited string = `<div class="pole_marker_shadow edited"></div>`
	pmHtmlDragEdited   string = `<i class="fas fa-expand-arrows-alt fa-fw fa-2x pole-marker drag"></i>`
)

// UpdateFromState updates the look of PoleMarker receiver depending on its map' Categorizer state and its data (no refresh undertaken nor layer group updated)
func (pm *PoleMarker) UpdateFromState() {
	pm.Map.Categorizer.PoleMarkerVisual(pm)
	pm.UpdateDivIconClassname(pm.Class)
	pm.UpdateDivIconHtml()
}

func (pm *PoleMarker) visualByState() (html, class string) {
	switch pm.Pole.State {
	case poleconst.StateNotSubmitted:
		//html = pmHtmlPin
		class = ""
	case poleconst.StateNoGo:
		html = pmHtmlOutline
		class = "red"
	case poleconst.StatePermissionPending:
		class = "purple"
	case poleconst.StateDictToDo:
		class = "darkred"
	case poleconst.StateDaToDo:
		class = "red"
	case poleconst.StateDaExpected:
		class = "orange"
	case poleconst.StateToDo:
		class = "blue"
	case poleconst.StateMarked:
		class = "lightblue"
	case poleconst.StateNoAccess:
		class = "darkblue"
	case poleconst.StateDenseNetwork:
		class = "darkblue"
	case poleconst.StateHoleDone:
		html = pmHtmlHole
		class = "blue"
	case poleconst.StateIncident:
		html = pmHtmlExclam
		class = "blue"
	case poleconst.StateDone:
		class = "green"
	case poleconst.StateAttachment:
		class = "darkgreen"
	case poleconst.StateCancelled:
		class = "grey"
	default:
		class = "darkred"
	}

	//if pm.Pole.HasProduct(poleconst.ProductReplenishment) || pm.Pole.HasProduct(poleconst.ProductFarReplenishment) {
	//	html = pmHtmlReplenish
	//}

	if html == "" {
		html = pmHtmlPin
		switch {
		case pm.Pole.HasProduct(poleconst.ProductReplenishment):
			html = pmHtmlReplenish
		case pm.Pole.HasProduct(poleconst.ProductFarReplenishment):
			html = pmHtmlReplenish
		case pm.Pole.HasProduct(poleconst.ProductTrickyReplace):
			html = pmHtmlTrickyReplace
		case pm.Pole.HasProduct(poleconst.ProductReplace):
			html = pmHtmlReplace
		case pm.Pole.HasProduct(poleconst.ProductCreation):
			html = pmHtmlCreation
			//case pm.Pole.HasProduct(poleconst.ProductCoated):
			//	html = pmHtmlCreation
		}
	}
	if pm.Pole.Priority > 0 && !pm.Pole.IsAlreadyDone() {
		html += "<div class=\"pole_marker_priority\">P" + strconv.Itoa(pm.Pole.Priority) + "</div>"
	}
	return
}

func (pm *PoleMarker) visualByAge(groupName string) (html, class string) {
	// as only class are modified, call visual by state to retreive polemarker html
	html, class = pm.visualByState()
	// set class trivial case
	switch groupName {
	case "BO":
		class = "red"
		return
	case "NS":
		class = "grey"
		return
	case "ERR":
		class = "darkred"
		return
	}
	// set class for non trivial case
	nbWeek, _ := strconv.Atoi(groupName)
	if nbWeek <= 0 {
		// pole is ready until N weeks
		nbWeek = -nbWeek
		switch nbWeek {
		case 0, 1:
			class = "lightblue"
		case 2, 3:
			class = "blue"
		default:
			class = "darkblue"
		}
		return
	}
	// pole will be ready in N weeks
	switch nbWeek {
	case 1:
		class = "purple"
	case 2:
		class = "orange"
	default:
		class = "darkorange"
	}
	return
}

// UpdateDivIconHtml updates the look and the shadow of receiver (no refresh undertaken)
func (pm *PoleMarker) UpdateDivIconHtml() {
	switch {
	case pm.Draggable:
		pm.Marker.UpdateDivIconHtml(pm.Html + pmHtmlDragEdited)
	case pm.Edited:
		pm.Marker.UpdateDivIconHtml(pm.Html + pmHtmlShadowEdited)
		pm.SetOpacity(poleconst.OpacitySelected)
	case pm.Selected:
		pm.Marker.UpdateDivIconHtml(pm.Html + pmHtmlShadow)
	default:
		pm.Marker.UpdateDivIconHtml(pm.Html)
		pm.SetOpacity(poleconst.OpacityNormal)
	}
}

func (pm *PoleMarker) UpdateTitle() {
	pm.Marker.UpdateToolTip(pm.Pole.GetTitle())
}

func (pm *PoleMarker) UpdateMarkerLatLng() {
	pm.Marker.SetLatLng(leaflet.NewLatLng(pm.Pole.Lat, pm.Pole.Long))
}

func (pm *PoleMarker) SetLatLng(latlng *leaflet.LatLng) {
	pm.Pole.Lat, pm.Pole.Long = latlng.ToFloats()
	pm.Marker.SetLatLng(latlng)
}
