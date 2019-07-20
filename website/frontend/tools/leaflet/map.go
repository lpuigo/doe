package leaflet

import "github.com/gopherjs/gopherjs/js"

// Map is a leaflet map object: http://leafletjs.com/reference-1.0.2.html#map
type Map struct {
	*js.Object
}

// NewMap creates a new map in the div specified by divID.
func NewMap(divID string, options *MapOptions) *Map {
	return &Map{
		Object: L.Call("map", divID),
	}
}

// DefaultMapOptions returns the default Map options.
func DefaultMapOptions() *MapOptions {
	return &MapOptions{
		Object: js.Global.Get("Object").New(),
	}
}

// MapOptions specify the options for a map:
// http://leafletjs.com/reference-1.0.2.html#map.
// They need to be initialized with DefaultMapOptions.
type MapOptions struct {
	Object       *js.Object
	PreferCanvas bool `js:"preferCanvas"`
}

// SetView sets the center and zoom level of the map.
func (m *Map) SetView(center *LatLng, zoom int) {
	m.Object.Call("setView", center, zoom)
}

// SetView sets the center and zoom level of the map.
func (m *Map) SetZoom(zoom int) {
	m.Object.Call("setZoom", zoom)
}

// FitBounds sets a map view that contains the given geographical bounds with the maximum zoom level possible.
func (m *Map) FitBounds(min, max *LatLng) {
	m.Object.Call("fitBounds", NewLatLngBounds(min, max))
}

// CreatePane creates a new Pane with the given name:
// http://leafletjs.com/reference-1.0.2.html#map-createpane
func (m *Map) CreatePane(name string) *Pane {
	return &Pane{Object: m.Object.Call("createPane", name)}
}

// Pane is a leaflet pane.
type Pane struct {
	*js.Object
}

// SetZIndex sets the Z index of the pane.
func (p *Pane) SetZIndex(index int) {
	p.Object.Get("style").Set("zIndex", index)
}
