package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/poleedit"
	"github.com/lpuig/ewin/doe/website/frontend/comp/polemap"
	"github.com/lpuig/ewin/doe/website/frontend/model/polesite"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements/message"
	"github.com/lpuig/ewin/doe/website/frontend/tools/leaflet"
	"honnef.co/go/js/xhr"
	"strconv"
	"strings"
)

//go:generate bash ./makejs.sh

func main() {
	mpm := NewMainPageModel()

	hvue.NewVM(
		polemap.RegisterComponent(),
		poleedit.RegisterComponent(),
		hvue.El("#polesites_app"),
		hvue.DataS(mpm),
		hvue.MethodsOf(mpm),
		hvue.Mounted(func(vm *hvue.VM) {
			mpm := &MainPageModel{Object: vm.Object}
			//mpm.Poles = []*polesite.Pole{}
			tools.BeforeUnloadConfirmation(mpm.PreventLeave)
			mpm.LoadPole()
		}),
		hvue.Computed("Title", func(vm *hvue.VM) interface{} {
			mpm := &MainPageModel{Object: vm.Object}
			if mpm.Polesite.Object == js.Undefined {
				return ""
			}
			if tools.Empty(mpm.Polesite.Client) && tools.Empty(mpm.Polesite.Ref) {
				return ""
			}
			return mpm.Polesite.Client + " / " + mpm.Polesite.Ref
		}),
	)

	js.Global.Set("mpm", mpm)
}

type MainPageModel struct {
	*js.Object

	VM       *hvue.VM           `js:"VM"`
	Filter   string             `js:"Filter"`
	Polesite *polesite.Polesite `js:"Polesite"`
	//Poles              []*polesite.Pole    `js:"Poles"`
	PolesGroup         *leaflet.LayerGroup `js:"PolesGroup"`
	SelectedPoleMarker *polemap.PoleMarker `js:"SelectedPoleMarker"`
	IsPoleSelected     bool                `js:"IsPoleSelected"`
	Dirty              bool                `js:"Dirty"`
}

func NewMainPageModel() *MainPageModel {
	mpm := &MainPageModel{Object: tools.O()}
	mpm.VM = nil
	mpm.Filter = ""
	mpm.Polesite = polesite.NewPolesite()
	//mpm.Poles = []*polesite.Pole{}
	mpm.PolesGroup = nil
	mpm.SelectedPoleMarker = nil
	mpm.IsPoleSelected = false
	mpm.Dirty = false
	return mpm
}

func (mpm *MainPageModel) LoadPole() {
	opsid := tools.GetURLSearchParam("psid")
	if opsid == nil {
		print("psid undefined")
		return
	}
	if opsid.String() == "" {
		print("psid empty")
		return
	}
	go mpm.callGetPolesite(opsid.Int())
}

func (mpm *MainPageModel) PreventLeave() bool {
	return mpm.Dirty
}

// UpdateMap updates current Poles Arrays in PoleMap component
func (mpm *MainPageModel) UpdateMap() {
	pm := polemap.PoleMapFromJS(mpm.VM.Refs("MapEwin"))
	mpm.PolesGroup = pm.AddPoles(mpm.Polesite.Poles, "Poteaux")
}

// CenterMapOnPoles centers PoleMap component to show all poles
func (mpm *MainPageModel) CenterMapOnPoles() {
	pm := polemap.PoleMapFromJS(mpm.VM.Refs("MapEwin"))
	pm.CenterOnPoles()
}

// MarkerClick handles marker-click PoleMap events
func (mpm *MainPageModel) MarkerClick(poleMarkerObj, event *js.Object) {
	pm := polemap.PoleMarkerFromJS(poleMarkerObj)
	mpm.SelectPole(pm)
}

func (mpm *MainPageModel) SelectPole(pm *polemap.PoleMarker) {
	if mpm.IsPoleSelected {
		mpm.SelectedPoleMarker.EndEditMode()
	}
	mpm.SelectedPoleMarker = pm
	pm.StartEditMode()
	mpm.IsPoleSelected = true

	pm.CenterOnMap(20)
}

func (mpm *MainPageModel) UnSelectPole() {
	mpm.SelectedPoleMarker.EndEditMode()
	mpm.SelectedPoleMarker = nil
	mpm.IsPoleSelected = false
}

//
func (mpm *MainPageModel) SwitchPoleState(pm *polemap.PoleMarker) {
	pm.Pole.SwitchState()
	pm.UpdateFromState()
	pm.Refresh()
}

//
func (mpm *MainPageModel) ApplyFilter(vm *hvue.VM) {
	mpm = &MainPageModel{Object: vm.Object}
	defer mpm.PolesGroup.Refresh()

	if mpm.Filter == "" {
		mpm.PolesGroup.ForEach(func(l *leaflet.Layer) {
			pm := polemap.PoleMarkerFromJS(l.Object)
			pm.SetOpacity(0.50)
		})
		//mpm.CenterMapOnPoles()
		return
	}

	var minLat, minLong, maxLat, maxLong float64 = 500, 500, 0, 0
	min := func(lat, long float64) {
		if lat < minLat {
			minLat = lat
		}
		if long < minLong {
			minLong = long
		}
	}

	max := func(lat, long float64) {
		if lat > maxLat {
			maxLat = lat
		}
		if long > maxLong {
			maxLong = long
		}
	}

	minmax := func(pm *polemap.PoleMarker) {
		lat, long := pm.GetLatLong().ToFloats()
		min(lat, long)
		max(lat, long)
	}

	expected := strings.ToUpper(mpm.Filter)
	filter := func(pm *polemap.PoleMarker) bool {
		return strings.Contains(strings.ToUpper(pm.Pole.Ref), expected)
	}
	found := false
	mpm.PolesGroup.ForEach(func(l *leaflet.Layer) {
		pm := polemap.PoleMarkerFromJS(l.Object)
		if filter(pm) {
			minmax(pm)
			pm.SetOpacity(0.80)
			found = true
		} else {
			pm.SetOpacity(0.20)
		}
	})
	if found {
		pm := polemap.PoleMapFromJS(mpm.VM.Refs("MapEwin"))
		pm.LeafletMap.Map.FitBounds(leaflet.NewLatLng(minLat, minLong), leaflet.NewLatLng(maxLat, maxLong))
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// WS call Methods

func (mpm *MainPageModel) callGetPolesite(psid int) {
	req := xhr.NewRequest("GET", "/api/polesites/"+strconv.Itoa(psid))
	req.Timeout = tools.LongTimeOut
	req.ResponseType = xhr.JSON
	defer func() {
		mpm.UpdateMap()
	}()
	err := req.Send(nil)
	if err != nil {
		message.ErrorStr(mpm.VM, "Oups! "+err.Error(), true)
		return
	}
	if req.Status != tools.HttpOK {
		message.ErrorRequestMessage(mpm.VM, req)
		return
	}
	mpm.Polesite = polesite.PolesiteFromJS(req.Response)
}
