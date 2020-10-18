package leaflet

import (
	"github.com/gopherjs/gopherjs/js"
)

// Control is a leaflet Control object: https://leafletjs.com/reference-1.5.0.html#control.
type Control struct {
	*js.Object
}

// AddTo add the receiver to the specified Map.
func (c *Control) AddTo(m *Map) {
	c.Object.Call("addTo", m)
}

// ControlLayers is a leaflet Control.Layers object: https://leafletjs.com/reference-1.5.0.html#control-layers.
type ControlLayers struct {
	Control
}

func NewControlLayers(baseLayer js.M, overlays js.M) *ControlLayers {
	return &ControlLayers{
		Control{L.Get("control").Call("layers", baseLayer, overlays)},
	}
}

// AddBaseLayer adds a base layer (radio button entry) with the given name to the control.
func (cl *ControlLayers) AddBaseLayer(layer *Layer, name string) {
	cl.Call("addBaseLayer", layer, name)
}

// AddOverlay adds an overlay (checkbox entry) with the given name to the control.
func (cl *ControlLayers) AddOverlay(layer *Layer, name string) {
	cl.Call("addOverlay", layer, name)
}

// RemoveLayer Remove the given layer from the control.
func (cl *ControlLayers) RemoveLayer(layer *Layer) {
	cl.Call("removeLayer", layer)
}

// ControlScale is a leaflet Control.Scale object: https://leafletjs.com/reference-1.6.0.html#control-scale.
type ControlScale struct {
	Control
}

func NewControlScale() *ControlScale {
	scaleOpt := js.M{
		"metric":         true,
		"imperial":       false,
		"updateWhenIdle": true,
	}
	return &ControlScale{
		Control{L.Get("control").Call("scale", scaleOpt)},
	}
}

// ControlInfo is a leaflet Control dedicated to display info
type ControlInfo struct {
	Control
	Div *js.Object `js:"_div"`
}

func NewControlInfo() *ControlInfo {
	ci := &ControlInfo{Control: Control{L.Call("control")}}
	ci.Div = nil
	ci.Set("onAdd", ci.onAdd)
	return ci
}

func (ci *ControlInfo) onAdd(m *Map) *js.Object {
	ci.Div = L.Get("DomUtil").Call("create", "div", "control-info")
	//ci.Set("_div", L.Get("DomUtil").Call("create", "div", "control-info"))
	ci.Update("")
	return ci.Get("_div")
}

func (ci *ControlInfo) Update(html string) {
	if html == "" {
		L.Get("DomUtil").Call("addClass", ci.Div, "hidden")
		return
	}
	L.Get("DomUtil").Call("removeClass", ci.Div, "hidden")
	ci.Div.Set("innerHTML", html)
}

// ControlCurrentPos is a Leaflet Control dedicated to get current user position
type ControlCurrentPos struct {
	Control
	Div *js.Object `js:"_div"`
}

func NewControlCurrentPos() *ControlCurrentPos {
	opt := js.M{
		"position": "topleft",
	}
	ccp := &ControlCurrentPos{Control: Control{L.Call("control", opt)}}
	ccp.Div = nil
	ccp.Set("onAdd", ccp.onAdd)
	return ccp
}

func (ccp *ControlCurrentPos) onAdd(m *Map) *js.Object {
	ccp.Div = L.Get("DomUtil").Call("create", "div", "control-currentpos")
	ccp.Div.Set("innerHTML", `<div><i class="fas fa-map-marker-alt icon--big"></i></div>`)
	ccp.Div.Call("addEventListener", "click", func(event *js.Object) {
		//print("ControlCurrentPos Click")
		event.Call("stopPropagation")
		GetCurrentPosition(func(pos *LatLng) {
			m.SetView(pos, 18)
		})
	})
	return ccp.Get("_div")
}

func (ccp *ControlCurrentPos) Update(html string) {
	ccp.Div.Set("innerHTML", "<div><span>Pos</span></div>")
}
