package polemap

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/model/polesite"
	"github.com/lpuig/ewin/doe/website/frontend/model/polesite/poleconst"
	"github.com/lpuig/ewin/doe/website/frontend/tools/leaflet"
)

type PoleMarker struct {
	leaflet.Marker
	Pole *polesite.Pole `js:"Pole"`
}

func PoleMarkerFromJS(obj *js.Object) *PoleMarker {
	return &PoleMarker{Marker: *leaflet.MarkerFromJs(obj)}
}

func NewPoleMarker(option *leaflet.MarkerOptions, pole *polesite.Pole) *PoleMarker {
	np := &PoleMarker{Marker: *leaflet.NewMarker(pole.Lat, pole.Long, option)}
	np.Pole = pole
	return np
}

func DefaultPoleMarker() *PoleMarker {
	opt := leaflet.DefaultMarkerOption()
	np := &PoleMarker{Marker: *leaflet.NewMarker(0.0, 0.0, opt)}
	np.Pole = polesite.NewPole()
	return np
}

func (pm *PoleMarker) StartEditMode() {
	pm.SetOpacity(poleconst.OpacitySelected)
	pm.SetDraggable(true)
	pm.Refresh()
}

func (pm *PoleMarker) EndEditMode(refresh bool) {
	pm.SetOpacity(poleconst.OpacityNormal)
	pm.SetDraggable(false)
	if refresh {
		pm.Refresh()
	}
}

const (
	pmHtmlPin     string = `<i class="fas fa-map-pin fa-3x"></i>`
	pmHtmlPlain   string = `<i class="fas fa-map-marker fa-3x"></i>`
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
	case poleconst.StateToDo:
		html = pmHtmlPlain
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
		class = "purple"
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
	pm.Marker.UpdateToolTip(pm.Pole.Ref)
}

func (pm *PoleMarker) UpdateMarkerLatLng() {
	pm.Marker.SetLatLng(leaflet.NewLatLng(pm.Pole.Lat, pm.Pole.Long))
}

func (pm *PoleMarker) SetLatLng(latlng *leaflet.LatLng) {
	pm.Pole.Lat, pm.Pole.Long = latlng.ToFloats()
	pm.Marker.SetLatLng(latlng)
}
