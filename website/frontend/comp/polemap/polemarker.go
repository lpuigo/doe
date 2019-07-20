package polemap

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/mappoc/frontend/leaflet"
	"github.com/lpuig/ewin/mappoc/frontend/mappoc/model"
)

type PoleMarker struct {
	leaflet.Marker
	Pole *model.Pole `js:"Pole"`
}

const (
	MarkerOpacityDefault  float64 = 0.5
	MarkerOpacitySelected float64 = 1
)

func PoleMarkerFromJS(obj *js.Object) *PoleMarker {
	return &PoleMarker{Marker: *leaflet.MarkerFromJs(obj)}
}

func NewPoleMarker(option *leaflet.MarkerOptions, pole *model.Pole) *PoleMarker {
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
	case model.PoleStateNotSubmitted:
		html = pmHtmlPin
		class = ""
	case model.PoleStateToDo:
		html = pmHtmlPlain
		class = "blue"
	case model.PoleStateHoleDone:
		html = pmHtmlHole
		class = "orange"
	case model.PoleStateIncident:
		html = pmHtmlHole
		class = "red"
	case model.PoleStateDone:
		html = pmHtmlPlain
		class = "green"
	case model.PoleStateCancelled:
		html = pmHtmlPlain
		class = ""
	}

	pm.UpdateDivIconClassname(class)
	pm.UpdateDivIconHtml(html)
}

func (pm *PoleMarker) UpdateTitle() {
	pm.Marker.UpdateToolTip(pm.Pole.Ref)
}
