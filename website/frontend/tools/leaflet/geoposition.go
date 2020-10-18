package leaflet

import "github.com/gopherjs/gopherjs/js"

func GetCurrentPosition(f func(lng *LatLng)) {
	go func() {
		js.Global.Get("navigator").Get("geolocation").Call("getCurrentPosition", func(obj *js.Object) {
			coords := obj.Get("coords")
			lat, long := coords.Get("latitude").Float(), coords.Get("longitude").Float()
			f(NewLatLng(lat, long))
		})
	}()
}
