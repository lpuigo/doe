package leaflet

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

// Marker is a leaflet Marker: https://leafletjs.com/reference-1.5.0.html#marker.
type Marker struct {
	Layer
}

// NewMarker creates a new Marker
func NewMarker(lat, long float64, option *MarkerOptions) *Marker {
	return &Marker{
		Layer: Layer{
			Object: L.Call("marker", NewLatLng(lat, long), option),
		},
	}
}

func MarkerFromJs(o *js.Object) *Marker {
	return &Marker{
		Layer: Layer{
			Object: o,
		},
	}
}

func (m *Marker) BindPopup(content string) {
	m.Call("bindPopup", content)
}

func (m *Marker) UpdateDivIconClassname(newclass string) {
	m.Object.Get("options").Get("icon").Get("options").Set("className", newclass)
}

func (m *Marker) UpdateDivIconHtml(newhtml string) {
	m.Object.Get("options").Get("icon").Get("options").Set("html", newhtml)
}

func (m *Marker) UpdateToolTip(text string) {
	m.Object.Get("options").Set("title", text)
}

func (m *Marker) SetDraggable(drag bool) {
	if drag {
		m.Object.Get("dragging").Call("enable")
		return
	}
	m.Object.Get("dragging").Call("disable")
	m.Object.Get("options").Set("draggable", false)
}

func (m *Marker) SetOpacity(op float64) {
	m.Object.Get("options").Set("opacity", op)
}

// SetLatLng returns the current geographical position of the marker.
func (m *Marker) GetLatLng() *LatLng {
	return LatLngFromJS(m.Call("getLatLng"))
}

// SetLatLng changes the marker geographical position to the given point.
func (m *Marker) SetLatLng(latlng *LatLng) {
	m.Call("setLatLng", latlng)
}

type MarkerOptions struct {
	*js.Object
	Icon                *Icon   `js:"icon"`
	Keyboard            bool    `js:"keyboard"`
	Title               string  `js:"title"`
	Alt                 string  `js:"alt"`
	ZIndexOffset        float64 `js:"zIndexOffset"`
	Opacity             float64 `js:"opacity"`
	RiseOnHover         bool    `js:"riseOnHover"`
	RiseOffset          bool    `js:"riseOffset"`
	Pane                string  `js:"pane"`
	BubblingMouseEvents bool    `js:"bubblingMouseEvents"`
	Draggable           bool    `js:"draggable"`
	AutoPan             bool    `js:"autoPan"`
	//AutoPanPadding Point `js:"autoPanPadding"`
	AutoPanSpeed float64 `js:"autoPanSpeed"`
}

func DefaultMarkerOption() *MarkerOptions {
	mo := &MarkerOptions{Object: tools.O()}
	mo.Keyboard = false
	return mo
}

// Icon is a leaflet Icon object: https://leafletjs.com/reference-1.5.0.html#icon.
type Icon struct {
	*js.Object
}

// DivIcon is a leaflet DivIcon object: https://leafletjs.com/reference-1.5.0.html#divicon.
type DivIcon struct {
	Icon
}

func NewDivIcon(options *DivIconOptions) *DivIcon {
	return &DivIcon{
		Icon: Icon{Object: L.Call("divIcon", options)},
	}
}

type DivIconOptions struct {
	*js.Object
	Html          string `js:"html"`
	BgPos         js.S   `js:"bgPos"`
	IconSize      js.S   `js:"iconSize"`
	IconAnchor    js.S   `js:"iconAnchor"`
	PopupAnchor   js.S   `js:"popupAnchor"`
	TooltipAnchor js.S   `js:"tooltipAnchor"`
	ClassName     string `js:"className"`
}

func DefaultDivIconOptions() *DivIconOptions {
	dio := &DivIconOptions{Object: tools.O()}
	dio.Html = `<i class="fas fa-map-marker-alt fa-3x"></i>`
	dio.IconSize = js.S{27, 37}
	dio.IconAnchor = js.S{13, 37}
	dio.PopupAnchor = js.S{0, -37}
	dio.ClassName = ""
	return dio
}
