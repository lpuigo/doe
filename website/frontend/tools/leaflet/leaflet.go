// Package leaflet provides a (currently minimal) wrapper around leaflet.js
// for use with gopherjs. The bindings are currently for leaflet version 1.0.2.
package leaflet

import "github.com/gopherjs/gopherjs/js"

// L is the primary leaflet javascript object.
var L = js.Global.Get("L")

const (
	TileMaxNativeZoomLevel int = 19
	TileMaxZoomLevel       int = 21
)

// OSMTileLayer returns OpenStreetMap standard TileLayer
func OSMTileLayer() *TileLayer {
	tileOption := DefaultTileLayerOptions()
	tileOption.MaxNativeZoom = TileMaxNativeZoomLevel
	tileOption.MaxZoom = TileMaxZoomLevel
	tileOption.Attribution = `&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors`
	url := "https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
	tl := NewTileLayer(url, tileOption)
	print("OSMTileLayer:", tl.Object)
	return tl
}

const MapboxToken string = "pk.eyJ1IjoibGF1cmVudC1wdWlnIiwiYSI6ImNqeDgxazRqYzBmOGEzbnA3Z2lld3Rja2cifQ.Oq6cQfmK3uKYyVQffiIn_Q"

// MapBoxTileLayer returns mapbox standard TileLayer
//  mapbox.streets
//  mapbox.satellite
//  mapbox.outdoors
//  mapbox.light
func MapBoxTileLayer(id string) *TileLayer {
	tileOption := DefaultTileLayerOptions()
	tileOption.MaxNativeZoom = TileMaxNativeZoomLevel
	tileOption.MaxZoom = TileMaxZoomLevel
	tileOption.Attribution = `&copy; <a href="https://www.mapbox.com/about/maps/">Mapbox</a> &copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> <a href="https://apps.mapbox.com/feedback/">Improve this map</a>`
	tileOption.Id = id
	tileOption.AccesToken = MapboxToken
	url := "https://api.tiles.mapbox.com/v4/{id}/{z}/{x}/{y}.png?access_token={accessToken}"
	return NewTileLayer(url, tileOption)
}
