package polemap

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/model/polesite"
	"github.com/lpuig/ewin/doe/website/frontend/tools/leaflet"
)

type PoleMarker struct {
	leaflet.Marker
	Pole *polesite.Pole `js:"Pole"`
}

const (
	MarkerOpacityDefault  float64 = 0.5
	MarkerOpacitySelected float64 = 1
)

func PoleMarkerFromJS(obj *js.Object) *PoleMarker {
	return &PoleMarker{Marker: *leaflet.MarkerFromJs(obj)}
}

func NewPoleMarker(option *leaflet.MarkerOptions, pole *polesite.Pole) *PoleMarker {
	np := &PoleMarker{Marker: *leaflet.NewMarker(pole.Lat, pole.Long, option)}
	np.Pole = pole
	return np
}

func (pm *PoleMarker) StartEditMode() {
	pm.SetOpacity(MarkerOpacitySelected)
	pm.SetDraggable(true)
	pm.Refresh()
}

func (pm *PoleMarker) EndEditMode() {
	pm.SetOpacity(MarkerOpacityDefault)
	pm.SetDraggable(false)
	pm.Refresh()
}

const (
	pmHtmlPin   string = `<i class="fas fa-map-pin fa-3x"></i>`
	pmHtmlPlain string = `<i class="fas fa-map-marker fa-3x"></i>`
	pmHtmlHole  string = `<i class="fas fa-map-marker-alt fa-3x"></i>`
)

func (pm *PoleMarker) UpdateFromState() {
	var html, class string

	switch pm.Pole.State {
	case polesite.PoleStateNotSubmitted:
		html = pmHtmlPin
		class = ""
	case polesite.PoleStateToDo:
		html = pmHtmlPlain
		class = "blue"
	case polesite.PoleStateHoleDone:
		html = pmHtmlHole
		class = "orange"
	case polesite.PoleStateIncident:
		html = pmHtmlHole
		class = "red"
	case polesite.PoleStateDone:
		html = pmHtmlPlain
		class = "green"
	case polesite.PoleStateCancelled:
		html = pmHtmlPlain
		class = ""
	}

	pm.UpdateDivIconClassname(class)
	pm.UpdateDivIconHtml(html)
}

func (pm *PoleMarker) UpdateTitle() {
	pm.Marker.UpdateToolTip(pm.Pole.Ref)
}
