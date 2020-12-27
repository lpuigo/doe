package main

import (
	"strconv"
	"strings"

	"github.com/lpuig/ewin/doe/website/frontend/comp/poleinfoupdate"

	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/poleedit"
	"github.com/lpuig/ewin/doe/website/frontend/comp/polemap"
	"github.com/lpuig/ewin/doe/website/frontend/comp/poletable"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/model/polesite"
	"github.com/lpuig/ewin/doe/website/frontend/model/polesite/poleconst"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements/message"
	"github.com/lpuig/ewin/doe/website/frontend/tools/json"
	"github.com/lpuig/ewin/doe/website/frontend/tools/leaflet"
	"github.com/lpuig/ewin/doe/website/frontend/tools/nominatim"
	"honnef.co/go/js/xhr"
)

//go:generate bash ./makejs.sh

func main() {
	mpm := NewMainPageModel()

	hvue.NewVM(
		poleinfoupdate.RegisterComponent(),
		polemap.RegisterComponent(),
		poleedit.RegisterComponent(),
		poletable.RegisterComponent(),
		hvue.El("#polesites_app"),
		hvue.DataS(mpm),
		hvue.MethodsOf(mpm),
		hvue.Mounted(func(vm *hvue.VM) {
			mpm := &MainPageModel{Object: vm.Object}
			tools.BeforeUnloadConfirmation(mpm.PreventLeave)
			mpm.GetUserSession(func() {
				mpm.LoadPolesite(false)
			})
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
		hvue.Computed("IsDirty", func(vm *hvue.VM) interface{} {
			mpm := &MainPageModel{Object: vm.Object}
			mpm.Dirty = (mpm.Reference != json.Stringify(mpm.Polesite))
			return mpm.Dirty
		}),
		hvue.Computed("ShowMap", func(vm *hvue.VM) interface{} {
			mpm := &MainPageModel{Object: vm.Object}
			if mpm.ActiveMode != "Map" {
				return "display: none;"
			}
			return ""
		}),
		hvue.Computed("ShowTable", func(vm *hvue.VM) interface{} {
			mpm := &MainPageModel{Object: vm.Object}
			if mpm.ActiveMode != "Table" {
				return "display: none;"
			}
			return ""
		}),
		hvue.Computed("IsSearchAddressMsg", func(vm *hvue.VM) interface{} {
			mpm := &MainPageModel{Object: vm.Object}
			return mpm.SearchAddressMsg != ""
		}),
	)

	js.Global.Set("mpm", mpm)
}

type MainPageModel struct {
	*js.Object

	VM   *hvue.VM `js:"VM"`
	User *fm.User `js:"User"`

	ActiveMode string             `js:"ActiveMode"`
	Filter     string             `js:"Filter"`
	FilterType string             `js:"FilterType"`
	Polesite   *polesite.Polesite `js:"Polesite"`
	//PolesGroup         *leaflet.LayerGroup `js:"PolesGroup"`

	SelectedPoleMarker *polemap.PoleMarker `js:"SelectedPoleMarker"`
	IsPoleSelected     bool                `js:"IsPoleSelected"`
	ActiveChapter      []string            `js:"ActiveChapter"`
	TableContext       *poletable.Context  `js:"TableContext"`

	SearchAddress    string `js:"SearchAddress"`
	SearchAddressMsg string `js:"SearchAddressMsg"`

	VisibleSearchLoc    bool   `js:"VisibleSearchLoc"`
	VisibleTools        bool   `js:"VisibleTools"`
	VisibleToolsChapter string `js:"VisibleToolsChapter"`

	KizeoComplete bool         `js:"KizeoComplete"`
	kizeoReport   *kizeoReport `js:"KizeoReport"`

	ImportPoleComplete bool              `js:"ImportPoleComplete"`
	ImportPoleReport   *ImportPoleReport `js:"ImportPoleReport"`

	Reference string `js:"Reference"`
	Dirty     bool   `js:"Dirty"`

	//ColumnSelector *poletable.ColumnSelector `js:"ColumnSelector"`
}

func NewMainPageModel() *MainPageModel {
	mpm := &MainPageModel{Object: tools.O()}
	mpm.VM = nil
	mpm.User = fm.NewUser()
	mpm.ActiveMode = "Map"
	mpm.Filter = ""
	mpm.FilterType = poleconst.FilterValueAll
	mpm.Polesite = polesite.NewPolesite()
	//mpm.PolesGroup = nil

	mpm.SelectedPoleMarker = nil
	mpm.IsPoleSelected = false
	mpm.ActiveChapter = []string{}
	mpm.TableContext = poletable.NewContext("billing")
	mpm.SearchAddress = ""
	mpm.SearchAddressMsg = ""
	mpm.VisibleSearchLoc = false
	mpm.VisibleTools = false
	mpm.VisibleToolsChapter = "1"
	mpm.KizeoComplete = false
	mpm.kizeoReport = NewKizeoReport()

	mpm.ImportPoleComplete = false
	mpm.ImportPoleReport = NewImportPoleReport()

	mpm.Reference = ""
	mpm.Dirty = false
	//mpm.ColumnSelector = poletable.DefaultColumnSelector()
	return mpm
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// User Management Methods

func (mpm *MainPageModel) GetUserSession(callback func()) {
	go mpm.callGetUser(callback)
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Tools Methods

func (mpm *MainPageModel) GetPoleMap() *polemap.PoleMap {
	return polemap.PoleMapFromJS(mpm.VM.Refs("MapEwin"))
}

func (mpm *MainPageModel) PreventLeave() bool {
	return mpm.User.HasPermissionUpdate() && mpm.Dirty
}

// GetPoleMarkerById returns the PoleMarker associated with given Id (or nil if not found)
func (mpm *MainPageModel) GetPoleMarkerById(id int) *polemap.PoleMarker {
	return mpm.GetPoleMap().GetPoleMarkerById(id)
}

// initMap init Poles Array in PoleMap component
func (mpm *MainPageModel) initMap() {
	mpm.GetPoleMap().SetPoles(mpm.Polesite.Poles)
}

// UpdateMap updates current Poles Array in PoleMap component
func (mpm *MainPageModel) UpdateMap() {
	mpm.initMap()
	if mpm.IsFilterClear() {
		mpm.CenterMapOnPoles()
		return
	}
	mpm.ApplyFilter(mpm.VM)
}

// RefreshMap refreshes current Poles Array in PoleMap component
func (mpm *MainPageModel) RefreshMap() {
	mpm.GetPoleMap().RefreshPoles(mpm.Polesite.Poles)
}

// CenterMapOnPoles centers PoleMap component to show all poles
func (mpm *MainPageModel) CenterMapOnPoles() {
	mpm.GetPoleMap().CenterOnPoles()
}

// GetMapCenter returns PoleMap center location
func (mpm *MainPageModel) GetMapCenter() *leaflet.LatLng {
	return mpm.GetPoleMap().GetCenter()
}

// CenterMapOnLatLong centers PoleMap component on lat long position
func (mpm *MainPageModel) CenterMapOnLatLong(lat, long float64) {
	mpm.GetPoleMap().CenterOn(lat, long, poleconst.ZoomLevelOnPole)
}

func (mpm *MainPageModel) SelectPole(pm *polemap.PoleMarker, drag bool) {
	mpm.SelectedPoleMarker = pm
	mpm.TableContext.SelectedPoleId = pm.Pole.Id
	pm.StartEditMode(drag)
	mpm.IsPoleSelected = true

	pm.CenterOnMap(poleconst.ZoomLevelOnPole)
}

func (mpm *MainPageModel) UnSelectPole(refresh bool) {
	if mpm.IsPoleSelected {
		mpm.GetPoleMap().DisablePoleLine()
		mpm.CloseEditPole()
	}
}

func (mpm *MainPageModel) DictZipArchiveURL() string {
	return "/api/polesites/" + strconv.Itoa(mpm.Polesite.Id) + "/dictzip"
}

func (mpm *MainPageModel) ProgressXlsURL() string {
	return "/api/polesites/" + strconv.Itoa(mpm.Polesite.Id) + "/progress"
}

func (mpm *MainPageModel) RefExportXlsURL() string {
	return "/api/polesites/" + strconv.Itoa(mpm.Polesite.Id) + "/refexport"
}

func (mpm *MainPageModel) KizeoReportURL() string {
	return "/api/polesites/" + strconv.Itoa(mpm.Polesite.Id) + "/kizeo"
}

func (mpm *MainPageModel) ImportPoleURL() string {
	return "/api/polesites/" + strconv.Itoa(mpm.Polesite.Id) + "/import"
}

// ApplyFilterOnMap applies current Filter-Type and Filter value to Poles Markers and Map Region
func (mpm *MainPageModel) ApplyFilterOnMap() {
	poleMap := mpm.GetPoleMap()
	//defer poleMap.PoleMarkersGroup.Refresh()
	defer poleMap.RefreshGroup()

	if mpm.FilterType == poleconst.FilterValueAll && mpm.Filter == "" {
		for _, poleMarker := range poleMap.GetPoleMarkers() {
			poleMarker.SetOpacity(poleconst.OpacityNormal)
		}
		//mpm.PolesGroup.ForEach(func(l *leaflet.Layer) {
		//	pm := polemap.PoleMarkerFromJS(l.Object)
		//	pm.SetOpacity(poleconst.OpacityNormal)
		//})
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

	expected := strings.ToUpper(strings.Trim(mpm.Filter, " "))
	filter := func(pm *polemap.PoleMarker) bool {
		sis := pm.Pole.SearchString(mpm.FilterType)
		if sis == "" {
			return false
		}
		return strings.Contains(strings.ToUpper(sis), expected)
	}
	found := false
	for _, poleMarker := range poleMap.GetPoleMarkers() {
		if filter(poleMarker) {
			minmax(poleMarker)
			poleMarker.SetOpacity(poleconst.OpacityFiltered)
			found = true
		} else {
			poleMarker.SetOpacity(poleconst.OpacityBlur)
		}
	}
	//mpm.PolesGroup.ForEach(func(l *leaflet.Layer) {
	//	pm := polemap.PoleMarkerFromJS(l.Object)
	//	if filter(pm) {
	//		minmax(pm)
	//		pm.SetOpacity(poleconst.OpacityFiltered)
	//		found = true
	//	} else {
	//		pm.SetOpacity(poleconst.OpacityBlur)
	//	}
	//})
	if found {
		poleMap.FitBounds(leaflet.NewLatLng(minLat, minLong), leaflet.NewLatLng(maxLat, maxLong))
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// UI Methods

func (mpm *MainPageModel) GetFilterType() []*elements.ValueLabel {
	return polesite.GetFilterTypeValueLabel()
}

// GetPoleMarker returns the PoleMarker associated with given Pole (or nil if not found)
func (mpm *MainPageModel) GetPoleMarker(pole *polesite.Pole) *polemap.PoleMarker {
	return mpm.GetPoleMap().GetPoleMarkerById(pole.Id)
}

func (mpm *MainPageModel) LoadPolesite(refresh bool) {
	opsid := tools.GetURLSearchParam("psid")
	if opsid == nil {
		print("psid undefined")
		return
	}
	if opsid.String() == "" {
		print("psid empty")
		return
	}

	callback := mpm.UpdateMap
	if refresh {
		mpm.UnSelectPole(false)
		callback = mpm.RefreshMap
	}
	go mpm.callGetPolesite(opsid.Int(), callback)
}

func (mpm *MainPageModel) SavePolesite(vm *hvue.VM) {
	mpm = &MainPageModel{Object: vm.Object}
	mpm.UnSelectPole(false)
	go mpm.callUpdatePolesite(mpm.Polesite, mpm.RefreshMap)
}

// MarkerClick handles marker-click PoleMap events
func (mpm *MainPageModel) MarkerClick(poleMarkerObj, event *js.Object) {
	pm := polemap.PoleMarkerFromJS(poleMarkerObj)
	mpm.UnSelectPole(true)
	mpm.SelectPole(pm, event.Get("originalEvent").Get("ctrlKey").Bool())
}

// TablePoleSelected handles selected pole via PoleTable Component
func (mpm *MainPageModel) TablePoleSelected(context *poletable.Context) {
	pm := mpm.GetPoleMarkerById(context.SelectedPoleId)
	if pm == nil {
		return
	}
	mpm.UnSelectPole(true)
	mpm.SelectPole(pm, false)
}

// HandleTablePolesiteUpdate handles selected pole via PoleTable Component
func (mpm *MainPageModel) HandleTablePolesiteUpdate(msg string) {
	message.NotifySuccess(mpm.VM, "Mise à jour", msg)
	mpm.RefreshMap()

}

// HandleArchiveRefsGroup handles Archive Refs Group command via PoleTable Component
func (mpm *MainPageModel) HandleArchiveRefsGroup() {
	go mpm.callArchiveRefsGroup(mpm.Polesite, mpm.RefreshMap)
}

// CenterOnPole handles center-on-pole via PoleTable Component
func (mpm *MainPageModel) CenterOnPole(p *polesite.Pole) {
	if mpm.ActiveMode != "Map" {
		mpm.ActiveMode = "Map"
		//TODO Map Display None
		//mpm.initMap()
	}
	pm := mpm.GetPoleMarker(p)
	mpm.UnSelectPole(true)
	mpm.SelectPole(pm, false)
}

// SwitchActiveMode handles ActiveMode change
func (mpm *MainPageModel) SwitchActiveMode(vm *hvue.VM) {
	mpm = &MainPageModel{Object: vm.Object}
	switch mpm.ActiveMode {
	case "Map":
		mpm.ApplyFilterOnMap()
	case "Info":
		mpm.CloseEditPole()
	}
}

//
func (mpm *MainPageModel) CloseEditPole() {
	if !mpm.IsPoleSelected {
		return
	}
	mpm.IsPoleSelected = false
	mpm.SelectedPoleMarker.EndEditMode(true)
	mpm.SelectedPoleMarker = nil
	mpm.TableContext.SelectedPoleId = poletable.None
}

func (mpm *MainPageModel) UpdateSearchLocation(vm *hvue.VM) {
	mpm.SearchAddressMsg = ""
}

func (mpm *MainPageModel) SearchLocation(vm *hvue.VM) {
	n := nominatim.NewNominatim(vm)
	callback := func() {
		if !(n.Found && n.Err == "") {
			print("SearchLocation not found", n.Object)
			if n.Response.Length() == 0 {
				mpm.SearchAddressMsg = "Aucune localisation trouvée"
				return
			}
			mpm.SearchAddressMsg = n.Err
			return
		}
		mpm.SearchAddressMsg = ""
		mpm.CenterMapOnLatLong(n.Lat, n.Long)
		mpm.VisibleSearchLoc = false
	}
	n.SearchAdress(mpm.SearchAddress, callback)
}

// AddPole add given Pole to Polesite, and select new pertaining PoleMarker
func (mpm *MainPageModel) AddPole(newPole *polesite.Pole) {
	mpm.Polesite.AddNewPole(newPole)
	mpm.RefreshMap()

	newPoleMarker := mpm.GetPoleMarker(newPole)
	if newPoleMarker == nil {
		message.ErrorStr(mpm.VM, "Impossible de sélectionner le nouveau poteau", false)
		return
	}
	mpm.SelectPole(newPoleMarker, true)
}

// AddPoleList add given Pole Slice to Polesite
func (mpm *MainPageModel) AddPoleList(poles []*polesite.Pole) {
	for _, pole := range poles {
		mpm.Polesite.AddNewPole(pole)
	}
	mpm.RefreshMap()
}

//
func (mpm *MainPageModel) CreatePole(vm *hvue.VM) {
	mpm = &MainPageModel{Object: vm.Object}

	newPole := polesite.NewPole()
	newPole.Lat, newPole.Long = mpm.GetMapCenter().ToFloats()
	newPole.State = poleconst.StateNotSubmitted
	newPole.AddProduct(poleconst.ProductCreation)
	mpm.AddPole(newPole)
}

//
func (mpm *MainPageModel) DuplicatePole(vm *hvue.VM, dc *js.Object) {
	mpm = &MainPageModel{Object: vm.Object}
	duplicateContext := poleedit.DuplicateContextFromJS(dc)
	//pmToDuplicate := polemap.PoleMarkerFromJS(dc)
	newPole := duplicateContext.Model.Pole.Duplicate(duplicateContext.NewName(), 0.0001)

	mpm.AddPole(newPole)
}

//
func (mpm *MainPageModel) DeletePole(vm *hvue.VM, pm *js.Object) {
	mpm = &MainPageModel{Object: vm.Object}
	mpm.CloseEditPole()

	pmToDelete := polemap.PoleMarkerFromJS(pm)

	if !mpm.Polesite.DeletePole(pmToDelete.Pole) {
		print("DeletePole failed", pmToDelete.Object)
		return
	}
	//mpm.UnSelectPole(false)
	mpm.RefreshMap()
}

//
func (mpm *MainPageModel) ClearFilter(vm *hvue.VM) {
	mpm = &MainPageModel{Object: vm.Object}
	mpm.Filter = ""
	mpm.FilterType = ""
	mpm.ApplyFilter(vm)
}

// IsFilterClear returns true if filter is clear (no filter criteria)
func (mpm *MainPageModel) IsFilterClear() bool {
	return mpm.FilterType == poleconst.FilterValueAll && mpm.Filter == ""
}

// ApplyFilter refreshes the map with current filter/filter-type values. ApplyFilter is also called by the computed value ShowMap when Map is activated
func (mpm *MainPageModel) ApplyFilter(vm *hvue.VM) {
	mpm = &MainPageModel{Object: vm.Object}
	// TODO Map Display None
	if mpm.ActiveMode != "Map" {
		// No Op if Map is not active
		return
	}
	mpm.ApplyFilterOnMap()
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Kizeo Upload Methods

type kizeoReport struct {
	*js.Object
	NbUpdate   int      `js:"NbUpdate"`
	UnknownRef []string `js:"UnknownRef"`
}

func NewKizeoReport() *kizeoReport {
	kr := &kizeoReport{Object: tools.O()}
	kr.NbUpdate = 0
	kr.UnknownRef = []string{}
	return kr
}

type kizeoRefs struct {
	*js.Object
	Refs map[string]string `js:"Refs"`
}

func (mpm *MainPageModel) KizeoUploadError(vm *hvue.VM, err, file *js.Object) {
	mpm = &MainPageModel{Object: vm.Object}
	mpm.VisibleTools = false
	mpm.KizeoComplete = false
	message.ErrorStr(vm, err.String(), true)
}

func (mpm *MainPageModel) KizeoUploadSuccess(vm *hvue.VM, response, file *js.Object) {
	mpm = &MainPageModel{Object: vm.Object}
	refs := &kizeoRefs{Object: response}
	//mpm.kizeoReport = &kizeoReport{Object: response}
	mpm.kizeoReport = mpm.MatchKizeoRefs(refs)
	if mpm.kizeoReport.NbUpdate == 0 {
		message.WarningStr(vm, "Kizeo : aucune correspondance trouvée")
		return
	}
	mpm.KizeoComplete = true
}

func (mpm *MainPageModel) MatchKizeoRefs(refs *kizeoRefs) *kizeoReport {
	report := NewKizeoReport()
	poleDict := make(map[string]*polesite.Pole)
	for _, pole := range mpm.Polesite.Poles {
		poleRef := pole.Ref + "|" + pole.Sticker
		poleDict[poleRef] = pole
	}

	for poleTitle, info := range refs.Refs {
		if pole, found := poleDict[poleTitle]; found {
			pole.Kizeo = "OK " + info
			report.NbUpdate++
		} else {
			report.UnknownRef = append(report.UnknownRef, poleTitle)
		}
	}
	return report
}

func (mpm *MainPageModel) GetUnmatchingKizeoRefs(max int) []string {
	nbKo := len(mpm.kizeoReport.UnknownRef)
	nbKoExemples := nbKo
	offset := 0
	if nbKo > max {
		nbKoExemples = max
		offset = 1
	}
	res := make([]string, nbKoExemples)
	for i := 0; i < nbKoExemples-offset; i++ {
		res[i] = mpm.kizeoReport.UnknownRef[i]
	}
	if nbKoExemples < nbKo {
		res[nbKoExemples-1] = "..."
	}
	return res
}

func (mpm *MainPageModel) BeforeUpload(vm *hvue.VM, file *js.Object) bool {
	if file.Get("type").String() != "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
		message.ErrorStr(vm, "Le fichier '"+file.Get("name").String()+"' n'est pas un document Xlsx", false)
		return false
	}
	return true
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Import Pole Upload Methods

type ImportPoleReport struct {
	*js.Object
	Errors []string         `js:"Errors"`
	Poles  []*polesite.Pole `js:"Poles"`
}

func NewImportPoleReport() *ImportPoleReport {
	ipr := &ImportPoleReport{Object: tools.O()}
	ipr.Errors = []string{}
	ipr.Poles = []*polesite.Pole{}
	return ipr
}

func (mpm *MainPageModel) ImportPoleUploadError(vm *hvue.VM, err, file *js.Object) {
	mpm = &MainPageModel{Object: vm.Object}
	mpm.VisibleTools = false
	mpm.ImportPoleComplete = false
	message.ErrorStr(vm, err.String(), true)
}

func (mpm *MainPageModel) ImportPoleUploadSuccess(vm *hvue.VM, response, file *js.Object) {
	mpm = &MainPageModel{Object: vm.Object}
	mpm.ImportPoleReport = &ImportPoleReport{Object: response}
	if len(mpm.ImportPoleReport.Errors) == 1 && len(mpm.ImportPoleReport.Poles) == 0 {
		message.WarningStr(vm, "Import d'appui : "+mpm.ImportPoleReport.Errors[0])
		return
	}
	mpm.AddPoleList(mpm.ImportPoleReport.Poles)

	mpm.ImportPoleComplete = true
}

func (mpm *MainPageModel) GetImportPoleReportErrorsRefs(max int) []string {
	nbKo := len(mpm.ImportPoleReport.Errors)
	nbKoExemples := nbKo
	offset := 0
	if nbKo > max {
		nbKoExemples = max
		offset = 1
	}
	res := make([]string, nbKoExemples)
	for i := 0; i < nbKoExemples-offset; i++ {
		res[i] = mpm.ImportPoleReport.Errors[i]
	}
	if nbKoExemples < nbKo {
		res[nbKoExemples-1] = "..."
	}
	return res
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// WS call Methods

func (mpm *MainPageModel) callGetUser(callback func()) {
	req := xhr.NewRequest("GET", "/api/login")
	req.Timeout = tools.LongTimeOut
	req.ResponseType = xhr.JSON
	err := req.Send(nil)
	if err != nil {
		message.ErrorStr(mpm.VM, "Oups! "+err.Error(), true)
		return
	}
	if req.Status != tools.HttpOK {
		message.ErrorRequestMessage(mpm.VM, req)
		return
	}
	mpm.User.Copy(fm.UserFromJS(req.Response))
	if mpm.User.Name == "" {
		mpm.User = fm.NewUser()
		return
	}
	mpm.User.Connected = true
	callback()
}

func (mpm *MainPageModel) callGetPolesite(psid int, callback func()) {
	req := xhr.NewRequest("GET", "/api/polesites/"+strconv.Itoa(psid))
	req.Timeout = tools.LongTimeOut
	req.ResponseType = xhr.JSON
	err := req.Send(nil)
	if err != nil {
		message.ErrorStr(mpm.VM, "Oups! "+err.Error(), true)
		return
	}
	if req.Status != tools.HttpOK {
		message.ErrorRequestMessage(mpm.VM, req)
		return
	}
	mpm.subLoadPolesite(polesite.PolesiteFromJS(req.Response), callback)
}

func (mpm *MainPageModel) subLoadPolesite(ps *polesite.Polesite, callback func()) {
	mpm.Polesite = ps
	mpm.Reference = json.Stringify(ps)
	mpm.Polesite.CheckPolesStatus()
	newTitle := mpm.Polesite.Ref + " - " + mpm.Polesite.Client
	js.Global.Get("document").Set("title", newTitle)
	callback()
}

func (mpm *MainPageModel) callUpdatePolesite(ups *polesite.Polesite, callback func()) {
	//defer func() {}()
	req := xhr.NewRequest("PUT", "/api/polesites/"+strconv.Itoa(ups.Id))
	req.Timeout = tools.TimeOut
	req.ResponseType = xhr.JSON
	err := req.Send(json.Stringify(ups))
	if err != nil {
		message.ErrorStr(mpm.VM, "Oups! "+err.Error(), true)
		return
	}
	if req.Status != tools.HttpOK {
		message.ErrorRequestMessage(mpm.VM, req)
		return
	}
	UpdtPoleSite := polesite.UpdatedFromJS(req.Response)
	mpm.subLoadPolesite(UpdtPoleSite.Polesite, callback)

	msg := "Chantier sauvegardé"
	nbIgnored := len(UpdtPoleSite.IgnoredPoles)
	switch nbIgnored {
	case 0:
		message.SuccesStr(mpm.VM, msg)
	case 1:
		msg += ". Appui '" + UpdtPoleSite.IgnoredPoles[0] + "' ignoré"
		message.WarningStr(mpm.VM, msg)
	default:
		msg += ". " + strconv.Itoa(nbIgnored) + " appuis ignorés: " + strings.Join(UpdtPoleSite.IgnoredPoles, ", ")
		message.WarningStr(mpm.VM, msg)
	}
}

func (mpm *MainPageModel) callArchiveRefsGroup(ups *polesite.Polesite, callback func()) {
	if mpm.Dirty {
		message.ConfirmWarning(mpm.VM, "Sauvegarder des modifications avant d'archiver les groupes finalisés ?", func() { mpm.SavePolesite(mpm.VM) })
		return
	}

	req := xhr.NewRequest("GET", "/api/polesites/"+strconv.Itoa(ups.Id)+"/archivecompleted")
	req.Timeout = tools.TimeOut
	req.ResponseType = xhr.JSON
	err := req.Send(json.Stringify(ups))
	if err != nil {
		message.ErrorStr(mpm.VM, "Oups! "+err.Error(), true)
		return
	}
	if req.Status != tools.HttpOK {
		message.ErrorRequestMessage(mpm.VM, req)
		return
	}
	mpm.callGetPolesite(ups.Id, func() {
		message.SuccesStr(mpm.VM, "Groupes finalisés archivés")
		callback()
	})
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Data Management Methods

func (mpm *MainPageModel) DetectDuplicate() {
	mpm.Polesite.DetectDuplicate()
}

func (mpm *MainPageModel) DetectMissingDAValidation() {
	if mpm.Polesite.DetectMissingDAValidation() {
		mpm.RefreshMap()
	}
}

func (mpm *MainPageModel) DetectProductInconsistency() {
	mpm.Polesite.DetectProductInconsistency()
	mpm.GetPoleMap().RefreshPoles(mpm.Polesite.Poles)
}
