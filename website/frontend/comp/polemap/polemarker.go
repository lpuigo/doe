package polemap

import (
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
}

func PoleMarkerFromJS(obj *js.Object) *PoleMarker {
	return &PoleMarker{Marker: *leaflet.MarkerFromJs(obj)}
}

func NewPoleMarker(option *leaflet.MarkerOptions, pole *polesite.Pole) *PoleMarker {
	np := &PoleMarker{Marker: *leaflet.NewMarker(pole.Lat, pole.Long, option)}
	np.Pole = pole
	np.Map = nil
	np.Draggable = false
	return np
}

func DefaultPoleMarker() *PoleMarker {
	opt := leaflet.DefaultMarkerOption()
	np := &PoleMarker{Marker: *leaflet.NewMarker(0.0, 0.0, opt)}
	np.Pole = polesite.NewPole()
	return np
}

func (pm *PoleMarker) StartEditMode(drag bool) {
	pm.SetOpacity(poleconst.OpacitySelected)
	if drag {
		pm.SetDraggable(true)
	}
	pm.Map.PoleLine.SetTarget(pm)
	pm.Refresh()
}

func (pm *PoleMarker) EndEditMode(refresh bool) {
	pm.SetOpacity(poleconst.OpacityNormal)
	pm.SetDraggable(false)
	pm.Map.PoleLine.Reinit()
	if refresh {
		pm.Refresh()
	}
}

func (pm *PoleMarker) SetDraggable(drag bool) {
	pm.Draggable = drag
	pm.Marker.SetDraggable(pm.Draggable)
	pm.Refresh()
}

func (pm *PoleMarker) Refresh() {
	pm.Remove()
	pm.AddTo(pm.Map.Map)
}

func (pm *PoleMarker) RefreshState() *PoleMarker {
	pm.Map.RefreshPoles(pm.Map.Poles)
	return pm.Map.GetPoleMarkerById(pm.Pole.Id)
}

const (
	pmHtmlPin     string = `<i class="fas fa-map-pin fa-3x"></i>`
	pmHtmlPlain   string = `<i class="fas fa-map-marker fa-3x"></i>`
	pmHtmlBolt    string = `<i class="fas fa-bolt fa-3x"></i>`
	pmHtmlExclam  string = `<i class="fas fa-exclamation fa-3x"></i>`
	pmHtmlHole    string = `<i class="fas fa-map-marker-alt fa-3x"></i>`
	pmHtmlOutline string = `<i class="el-icon-location-outline" style="font-size: 3.3em"></i>`
)

func (pm *PoleMarker) UpdateFromState() {
	var html, class string

	switch pm.Pole.State {
	case poleconst.StateNotSubmitted:
		html = pmHtmlPin
		class = ""
	case poleconst.StateNoGo:
		html = pmHtmlOutline
		class = "red"
	case poleconst.StatePermissionPending:
		html = pmHtmlOutline
		class = "orange"
	case poleconst.StateToDo:
		html = pmHtmlPlain
		class = "blue"
	case poleconst.StateNoAccess:
		html = pmHtmlExclam
		class = "blue"
	case poleconst.StateDenseNetwork:
		html = pmHtmlBolt
		class = "blue"
	case poleconst.StateHoleDone:
		html = pmHtmlHole
		class = "orange"
	case poleconst.StateIncident:
		html = pmHtmlHole
		class = "red"
	case poleconst.StateDone:
		html = pmHtmlPlain
		class = "green"
	case poleconst.StateAttachment:
		html = pmHtmlPlain
		class = "darkgreen"
	case poleconst.StateCancelled:
		html = pmHtmlPlain
		class = ""
	default:
		html = pmHtmlPin
		class = "red"
	}

	pm.UpdateDivIconClassname(class)
	pm.UpdateDivIconHtml(html)
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
