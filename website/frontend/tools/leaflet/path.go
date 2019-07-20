package leaflet

import "github.com/gopherjs/gopherjs/js"

// Path is a leaflet path object: http://leafletjs.com/reference-1.0.2.html#path.
type Path struct {
	Layer
}

// SetStyle sets the style of the receiver:
// http://leafletjs.com/reference-1.0.2.html#path-setstyle.
func (p *Path) SetStyle(style *PathOptions) {
	p.Object.Call("setStyle", style)
}

// PathOptions specify the options for a path:
// http://leafletjs.com/reference-1.0.2.html#path-option.
// They need to be initialized with DefaultPathOptions.
type PathOptions struct {
	Object      *js.Object
	Stroke      bool    `js:"stroke"`
	Color       string  `js:"color"`
	Weight      int     `js:"weight"`
	Opacity     float64 `js:"opacity"`
	LineCap     string  `js:"lineCap"`
	LineJoin    string  `js:"lineJoin"`
	DashArray   string  `js:"dashArray"`
	DashOffset  string  `js:"dashOffset"`
	Fill        bool    `js:"fill"`
	FillColor   string  `js:"fillColor"`
	FillOpacity float64 `js:"fillOpacity"`
	FillRule    string  `js:"fillRule"`
}

// DefaultPathOptions returns the default TileLayer options.
func DefaultPathOptions() *PathOptions {
	return &PathOptions{
		Object: js.Global.Get("Object").New(),
	}
}

// Polyline is a leaflet polyline object: http://leafletjs.com/reference-1.0.2.html#polyline.
type Polyline struct {
	Path
}

// Polygon is a leaflet polygon object: http://leafletjs.com/reference-1.0.2.html#polygon.
type Polygon struct {
	Polyline
}

// NewPolygon creates a new polygon.
func NewPolygon(latlngs []*LatLng) *Polygon {
	return &Polygon{
		Polyline: Polyline{
			Path: Path{
				Layer: Layer{
					Object: L.Call("polygon", latlngs),
				},
			},
		},
	}
}
