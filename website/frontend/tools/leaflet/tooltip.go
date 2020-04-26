package leaflet

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

// ToolTip is a leaflet ToolTip: https://leafletjs.com/reference-1.6.0.html#tooltip.
type ToolTip struct {
	Layer
}

// NewMarker creates a new Marker
func NewToolTip(option ToolTypeOption, layer *Layer) *ToolTip {
	return &ToolTip{
		Layer: Layer{
			Object: L.Call("tooltip", option, layer),
		},
	}
}

func ToolTipFromJs(o *js.Object) *ToolTip {
	return &ToolTip{
		Layer: Layer{
			Object: o,
		},
	}
}

type ToolTypeOption struct {
	*js.Object
	Offset      *Point  `js:"offset"`
	Direction   string  `js:"direction"`
	Permanent   bool    `js:"permanent"`
	Sticky      bool    `js:"sticky"`
	Interactive bool    `js:"interactive"`
	Opacity     float64 `js:"opacity"`
}

func DefaultToolTypeOption() ToolTypeOption {
	to := ToolTypeOption{Object: tools.O()}
	to.Offset = NewPoint(0, 0)
	to.Permanent = true
	to.Opacity = 0.9
	return to
}
