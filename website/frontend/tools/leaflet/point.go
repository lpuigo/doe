package leaflet

import (
	"github.com/gopherjs/gopherjs/js"
)

// Point is a leaflet Point: https://leafletjs.com/reference-1.6.0.html#point.
type Point struct {
	*js.Object
	X float64 `js:"x"`
	Y float64 `js:"y"`
}

// NewPoint creates a new Point
func NewPoint(x, y float64) *Point {
	return &Point{Object: L.Call("point", x, y)}
}

func PointFromJs(o *js.Object) *Point {
	return &Point{Object: o}
}
