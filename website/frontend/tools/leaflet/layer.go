package leaflet

import "github.com/gopherjs/gopherjs/js"

// TileLayer is a leaflet TileLayer object: http://leafletjs.com/reference-1.0.2.html#tilelayer
type TileLayer struct {
	Layer
}

// NewTileLayer creates a new TileLayer with the specified URL template and
// options.
func NewTileLayer(urlTemplate string, options *TileLayerOptions) *TileLayer {
	return &TileLayer{
		Layer: Layer{
			Object: L.Call("tileLayer", urlTemplate, options),
		},
	}
}

// TileLayerOptions specifies the options for the TileLayer: http://leafletjs.com/reference-1.0.2.html#tilelayer-option
// They need to be initialized with DefaultTileLayerOptions.
type TileLayerOptions struct {
	*js.Object
	MinZoom       int      `js:"minZoom"`
	MaxZoom       int      `js:"maxZoom"`
	MinNativeZoom int      `js:"minNativeZoom"`
	MaxNativeZoom int      `js:"maxNativeZoom"`
	Subdomains    []string `js:"subdomains"`
	ErrorTileURL  string   `js:"errorTileUrl"`
	ZoomOffset    int      `js:"zoomOffset"`
	TMS           bool     `js:"tms"`
	ZoomReverse   bool     `js:"zoomReverse"`
	DetectRetina  bool     `js:"detectRetina"`
	CrossOrigin   bool     `js:"crossOrigin"`
	Id            string   `js:"id"`
	AccesToken    string   `js:"accessToken"`

	Pane        string `js:"pane"`
	Attribution string `js:"attribution"`
}

// DefaultTileLayerOptions returns the default TileLayer options.
func DefaultTileLayerOptions() *TileLayerOptions {
	return &TileLayerOptions{
		Object: js.Global.Get("Object").New(),
	}
}

// Layer is a leaflet layer object: http://leafletjs.com/reference-1.5.0.html#layer.
type Layer struct {
	*js.Object
}

// AddTo add the receiver to the specified Map.
func (l *Layer) AddTo(m *Map) {
	l.Object.Call("addTo", m)
}

// Remove removes the receiver from its current map.
func (l *Layer) Remove() {
	l.Object.Call("remove")
}

// Refresh refreshes the receiver by removing and adding it on its current map.
func (l *Layer) Refresh() {
	curMap := &Map{Object: l.Get("_map")}
	l.Remove()
	l.AddTo(curMap)
}

// CenterOnMap centers layer in its map, with given zoom level.
func (l *Layer) CenterOnMap(zoom int) {
	curMap := &Map{Object: l.Get("_map")}
	curMap.SetView(l.GetLatLong(), zoom)
}

func (l *Layer) GetLatLong() *LatLng {
	return &LatLng{Object: l.Get("_latlng")}
}

func (l *Layer) On(event string, handler func(*js.Object)) {
	l.Call("on", event, handler)
}

// LayerGroup is a leaflet LayerGroup: https://leafletjs.com/reference-1.5.0.html#layergroup.
type LayerGroup struct {
	Layer
}

func NewLayerGroup(layers []*Layer) *LayerGroup {
	return &LayerGroup{
		Layer{Object: L.Call("layerGroup", layers)},
	}
}

// ClearLayers removes all the layers from the group
func (lg *LayerGroup) ClearLayers() {
	lg.Call("clearLayers")
}

// ForEach iterates over the layers of the group
func (lg *LayerGroup) ForEach(f func(l *Layer)) {
	lg.Call("eachLayer", f)
}

// GridLayer is a leaflet GridLayer: http://leafletjs.com/reference-1.0.2.html#gridlayer.
type GridLayer struct {
	Layer
}

// NewGridLayer creates a new GridLayer.
func NewGridLayer() *GridLayer {
	return &GridLayer{
		Layer: Layer{
			Object: L.Call("gridLayer"),
		},
	}
}
