package polemap

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/leaflet"
	"math"
	"strconv"
)

type PoleLine struct {
	*js.Object

	PivotPole   *PoleMarker `js:"PivotPole"`
	IsPivotSet  bool        `js:"IsPivotSet"`
	TargetPole  *PoleMarker `js:"TargetPole"`
	IsTargetSet bool        `js:"IsTargetSet"`

	PolyLine   *leaflet.Polyline `js:"PolyLine"`
	IsDrawn    bool              `js:"IsDrawn"`
	RoundedPos *leaflet.LatLng   `js:"RoundedPos"`

	Map *PoleMap `js:"Map"`
}

func PoleLineFromJS(obj *js.Object) *PoleLine {
	return &PoleLine{Object: obj}
}

func NewPoleLine(polemap *PoleMap) *PoleLine {
	pl := &PoleLine{Object: tools.O()}
	pl.Map = polemap

	pathOpt := leaflet.DefaultPathOptions()
	pathOpt.Weight = 2
	pathOpt.DashArray = "10, 5"
	pathOpt.DashOffset = "0"
	pathOpt.Opacity = 0.6
	pl.PolyLine = leaflet.NewPolyline(nil, pathOpt)
	pl.PolyLine.AddTo(pl.Map.Map)
	topt := leaflet.DefaultToolTypeOption()
	topt.Offset = leaflet.NewPoint(15, -15)
	pl.IsDrawn = false
	pl.RoundedPos = nil
	pl.Reinit()
	return pl
}

func (pl *PoleLine) Reinit() {
	if pl.IsTargetSet {
		pl.TargetPole.UnBindTooltip()
	}
	pl.PivotPole = nil
	pl.IsPivotSet = false
	pl.TargetPole = nil
	pl.IsTargetSet = false
	pl.Erase()
}

func (pl *PoleLine) SetPivot(p *PoleMarker) {
	pl.PivotPole = p
	pl.IsPivotSet = true
}

func (pl *PoleLine) SetTarget(p *PoleMarker) {
	//if pl.IsTargetSet {
	//	if pl.TargetPole.Pole.Id != p.Pole.Id {
	//		// A new pole is targeted -> update the tooltip
	//		pl.TargetPole.UnBindTooltip()
	//	} else {
	//		// same pole is selected => it is a No Op
	//		return
	//	}
	//}
	pl.TargetPole = p
	pl.IsTargetSet = true
}

func (pl *PoleLine) Erase() {
	if !pl.IsDrawn {
		return
	}
	pl.PolyLine.SetLatLngs(nil)
	pl.IsDrawn = false
}

func (pl *PoleLine) Draw() {
	if !(pl.IsPivotSet && pl.IsTargetSet) {
		return
	}
	if pl.IsDrawn {
		pl.PolyLine.SetLatLngs(nil)
	}
	latlongs := []*leaflet.LatLng{
		pl.PivotPole.GetLatLng(),
		pl.TargetPole.GetLatLng(),
	}
	pl.PolyLine.SetLatLngs(latlongs)

	pl.TargetPole.UnBindTooltip()
	pl.BindToolTip()
	pl.IsDrawn = true
}

func (pl *PoleLine) RefreshTargetPos(targetLatLng *leaflet.LatLng) {
	if !(pl.IsPivotSet && pl.IsTargetSet && pl.IsDrawn) {
		return
	}
	roundedTargetLatLng, dist, _ := pl.GetLatLngAtRoundedDistance(targetLatLng)
	latlongs := []*leaflet.LatLng{
		pl.PivotPole.GetLatLng(),
		roundedTargetLatLng,
	}
	pl.RoundedPos = roundedTargetLatLng
	pl.PolyLine.SetLatLngs(latlongs)
	//pl.TargetPole.SetLatLng(roundedTargetLatLng)
	pl.TargetPole.SetTooltipContent(pl.GetDistanceString(dist))
}

func (pl *PoleLine) BindToolTip() {
	topt := leaflet.DefaultToolTypeOption()
	topt.Offset = leaflet.NewPoint(0, 5)
	topt.Direction = "bottom"
	topt.Opacity = 0.75
	pl.TargetPole.BindTooltip(pl.GetDistance(pl.TargetPole.GetLatLng()), topt)
}

func (pl *PoleLine) GetDistanceString(dist float64) string {
	return strconv.FormatFloat(dist, 'f', 2, 64) + " m"
}

func (pl *PoleLine) GetDistance(targetLatLng *leaflet.LatLng) string {
	dist := pl.PivotPole.GetLatLng().DistanceTo(targetLatLng)
	return pl.GetDistanceString(dist)
}

func (pl *PoleLine) GetLatLngAtRoundedDistance(targetLatLng *leaflet.LatLng) (*leaflet.LatLng, float64, float64) {
	dist := pl.PivotPole.GetLatLng().DistanceTo(targetLatLng)
	floorDist := math.Round(dist)
	ratio := floorDist / dist
	slat, slng := pl.PivotPole.GetLatLng().ToFloats()
	tlat, tlng := targetLatLng.ToFloats()
	dlat, dlng := (tlat-slat)*ratio, (tlng-slng)*ratio
	return leaflet.NewLatLng(slat+dlat, slng+dlng), dist, floorDist
}
