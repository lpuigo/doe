package leaflet

import "github.com/gopherjs/gopherjs/js"

// LatLng specifies a point in latitude and longitude
type LatLng struct {
	*js.Object
}

func LatLngFromJS(o *js.Object) *LatLng {
	return &LatLng{Object: o}
}

// NewLatLng returns a new LatLng object.
func NewLatLng(lat, lng float64) *LatLng {
	return &LatLng{
		Object: L.Call("latLng", lat, lng),
	}
}

func (ll *LatLng) ToFloats() (lat, long float64) {
	lat = ll.Get("lat").Float()
	long = ll.Get("lng").Float()
	return
}

func (ll *LatLng) DistanceTo(dest *LatLng) float64 {
	return ll.Call("distanceTo", dest).Float()
}

// LatLng specifies a point in latitude and longitude
type LatLngBounds struct {
	*js.Object
}

// NewLatLng returns a new LatLng object.
func NewLatLngBounds(bound1, bound2 *LatLng) *LatLngBounds {
	return &LatLngBounds{
		Object: L.Call("latLngBounds", bound1, bound2),
	}
}
