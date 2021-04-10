package polesites

import (
	"archive/zip"
	"fmt"
	"io"
	"math"
	"path/filepath"
	"sort"
	"strings"

	"github.com/lpuig/ewin/doe/website/backend/model/bpu"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	"github.com/lpuig/ewin/doe/website/backend/model/extraactivities"
	"github.com/lpuig/ewin/doe/website/backend/model/items"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/model/polesite/poleconst"
)

type UpdatedPoleSite struct {
	Polesite     *PoleSite
	IgnoredPoles []string
}

type PoleSite struct {
	Id         int
	Client     string
	Ref        string
	Manager    string
	OrderDate  string
	UpdateDate string
	Status     string
	Comment    string

	Poles           []*Pole
	ExtraActivities []*extraactivities.ExtraActivity
}

func (ps *PoleSite) GetRef() string {
	return ps.Ref
}

func (ps *PoleSite) GetClient() string {
	return ps.Client
}

func (ps *PoleSite) GetType() string {
	return "polesite"
}

func (ps *PoleSite) GetUpdateDate() string {
	return ps.UpdateDate
}

func (ps *PoleSite) Len() int {
	return len(ps.Poles)
}

func (ps *PoleSite) Swap(i, j int) {
	ps.Poles[i], ps.Poles[j] = ps.Poles[j], ps.Poles[i]
}

func (ps *PoleSite) Less(i, j int) bool {
	return ps.Poles[i].Id < ps.Poles[j].Id
}

// NextPoleId returns Id for a new pole to be added. Assume reciever PoleSite is sorted
func (ps *PoleSite) NextPoleId() int {
	if len(ps.Poles) == 0 {
		return 0
	}
	return ps.Poles[len(ps.Poles)-1].Id + 1
}

func (ps *PoleSite) GetInfo() *fm.PolesiteInfo {
	psi := fm.NewBEPolesiteInfo()

	psi.Id = ps.Id
	psi.Client = ps.Client
	psi.Ref = ps.Ref
	psi.Manager = ps.Manager
	psi.OrderDate = ps.OrderDate
	psi.UpdateDate = ps.UpdateDate
	psi.Status = ps.Status
	psi.Comment = ps.Comment

	psi.NbPole, psi.NbPoleBlocked, psi.NbPoleDone, psi.NbPoleBilled = ps.GetPolesNumbers()

	var searchBuilder strings.Builder
	fmt.Fprintf(&searchBuilder, "%s:%s,", "Client", strings.ToUpper(ps.Client))
	fmt.Fprintf(&searchBuilder, "%s:%s,", "Ref", strings.ToUpper(ps.Ref))
	fmt.Fprintf(&searchBuilder, "%s:%s,", "Manager", strings.ToUpper(ps.Manager))
	fmt.Fprintf(&searchBuilder, "%s:%s,", "OrderDate", strings.ToUpper(ps.OrderDate))
	fmt.Fprintf(&searchBuilder, "%s:%s,", "Comment", strings.ToUpper(ps.Comment))
	for _, pole := range ps.Poles {
		if pole.State == poleconst.StateDeleted {
			continue
		}
		fmt.Fprintf(&searchBuilder, "%s,", pole.SearchString())
	}
	psi.Search = searchBuilder.String()

	return psi
}

// GetPolesNumbers returns total, blocked and done number of Pullings
func (ps *PoleSite) GetPolesNumbers() (total, blocked, done, billed int) {
	for _, p := range ps.Poles {
		switch p.State {
		//case poleconst.StateNotSubmitted:
		//case poleconst.StateNoGo:
		case poleconst.StateDictToDo:
			total++
			blocked++
		case poleconst.StateDaToDo:
			total++
			blocked++
		case poleconst.StateDaExpected:
			total++
			blocked++
		case poleconst.StatePermissionPending:
			total++
			blocked++
		case poleconst.StateToDo, poleconst.StateMarked:
			total++
		case poleconst.StateNoAccess:
			total++
		case poleconst.StateDenseNetwork:
			total++
		case poleconst.StateHoleDone:
			total++
		case poleconst.StateIncident:
			total++
			blocked++
		case poleconst.StateDone:
			total++
			done++
		case poleconst.StateAttachment:
			total++
			billed++
			//case poleconst.StateCancelled:
			//case poleconst.StateDeleted:
		}
	}
	return
}

// SetUpdateDate set UpdateDate for receiver PoleSite, based on pole date (if no pole date available, use Today' date)
func (ps *PoleSite) SetUpdateDate() {
	updatedate := date.TimeJSMinDate
	for _, pole := range ps.Poles {
		if pole.State == poleconst.StateDeleted {
			continue
		}
		if pole.Date == "" {
			continue
		}
		if pole.Date > updatedate {
			updatedate = pole.Date
		}
	}
	if updatedate == date.TimeJSMinDate {
		updatedate = date.Today().String()
	}
	ps.UpdateDate = updatedate
}

type IsPolesiteVisible func(s *PoleSite) bool

// Itemize returns slice of item pertaining to polesite poles list
func (ps *PoleSite) Itemize(currentBpu *bpu.Bpu, doneOnly bool) ([]*items.Item, error) {
	res := []*items.Item{}

	for _, pole := range ps.Poles {
		if doneOnly && !pole.IsDone() {
			continue
		}
		items, err := pole.Itemize(ps.Client, ps.Ref, currentBpu)
		if err != nil {
			return nil, err
		}
		res = append(res, items...)
	}
	res = append(res, extraactivities.Itemize(ps.ExtraActivities, ps, doneOnly)...)
	return res, nil
}

// ExportName returns the PoleSite XLS export file name
func (ps *PoleSite) ExportName() string {
	return fmt.Sprintf("Polesite %s-%s (%d).xlsx", ps.Client, ps.Ref, ps.Id)
}

// RefExportName returns the PoleSite XLS reference export file name
func (ps *PoleSite) RefExportName() string {
	return fmt.Sprintf("Polesite %s-%s References.xlsx", ps.Client, ps.Ref)
}

// XLSExport returns the PoleSite XLS export
func (ps *PoleSite) XLSExport(w io.Writer) error {
	return ToExportXLS(w, ps)
}

// XLSRefExport returns the PoleSite XLS references export
func (ps *PoleSite) XLSRefExport(w io.Writer) error {
	return ToRefExportXLS(w, ps)
}

// ExportName returns the PoleSite XLS export file name
func (ps *PoleSite) ProgressName() string {
	year, weeknum := date.Today().ToTime().ISOWeek()
	return fmt.Sprintf("Avancement %s-%s (sem%02d-%04d).xlsx", ps.Client, ps.Ref, weeknum, year)
}

// XLSExport returns the PoleSite XLS export
func (ps *PoleSite) XLSProgress(w io.Writer) error {
	return ToProgressXLS(w, ps)
}

// ExportName returns the PoleSite XLS export file name
func (ps *PoleSite) DictZipName() string {
	return fmt.Sprintf("Polesite %s-%s.zip", ps.Client, ps.Ref)
}

// ExportName returns the PoleSite XLS export file name
func (ps *PoleSite) DictZipArchive(w io.Writer) error {
	type cityCount map[string]int

	zw := zip.NewWriter(w)
	path := strings.TrimSuffix(ps.DictZipName(), ".zip")
	makeDir := func(base ...string) string {
		return filepath.Join(base...) + "/"
	}

	// Create sorted List of DICT / citycount in PoleSite
	dicts := map[string]cityCount{}
	for _, pole := range ps.Poles {
		if pole.State == poleconst.StateDeleted {
			continue
		}
		dict := strings.Trim(pole.DictRef, " \t")
		if dict == "" {
			continue
		}
		cc := dicts[pole.DictRef]
		if cc == nil {
			cc = make(cityCount)
			dicts[pole.DictRef] = cc
		}
		cc[pole.City]++
	}
	dictList := make([]string, len(dicts))
	i := 0
	for dict, citycount := range dicts {
		// choose most common city
		selectedCity := ""
		maxcount := -1
		for city, count := range citycount {
			if count > maxcount {
				count = maxcount
				selectedCity = city
			}
		}
		res := dict
		if selectedCity != "" {
			res += " - " + selectedCity
		}
		dictList[i] = res
		i++
	}
	sort.Strings(dictList)

	for _, dict := range dictList {
		_, err := zw.Create(makeDir(path, dict))
		if err != nil {
			return err
		}
	}

	return zw.Close()
}

// CloneInfoForArchive returns a empty PoleSite with duplicated PoleSite data
func (ps *PoleSite) CloneInfoForArchive() *PoleSite {
	nps := &PoleSite{
		Id:         -1,
		Client:     ps.Client,
		Ref:        ps.Ref + " " + date.Today().String(),
		Manager:    ps.Manager,
		OrderDate:  ps.OrderDate,
		UpdateDate: date.Today().String(),
		Status:     poleconst.PsStatusDone,
		Comment:    fmt.Sprintf("Archive de %s au %s", ps.Ref, date.Today().ToDDMMYYYY()),
		Poles:      []*Pole{},
	}
	return nps
}

// MoveCompletedGroupTo moves poles sharing the same ref and all having a terminal status (Attachement or Cancelled)
// from receiver PoleSite to the given PoleSite, and returns archived pole number
func (ps *PoleSite) MoveCompletedGroupTo(aps *PoleSite) int {
	doArchive := make(map[string]bool)

	for _, pole := range ps.Poles {
		if pole.State == poleconst.StateDeleted {
			continue
		}
		doArchivePole := pole.IsArchivable()
		if doArchiveRef, found := doArchive[pole.Ref]; !found {
			doArchive[pole.Ref] = doArchivePole // first encountered pole state sets initial pertaining ref state
		} else {
			if doArchiveRef && !doArchivePole {
				doArchive[pole.Ref] = false
			}
		}
	}

	apid := 0
	unArchivedPoles := []*Pole{}
	for _, pole := range ps.Poles {
		if pole.State == poleconst.StateDeleted {
			continue
		}
		doArch := doArchive[pole.Ref]
		if strings.HasPrefix(pole.Ref, "Dépôt") && pole.IsArchivable() {
			doArch = true
		}
		if doArch {
			pole.Id = apid
			apid++
			aps.Poles = append(aps.Poles, pole)
		} else {
			unArchivedPoles = append(unArchivedPoles, pole)
		}
	}

	if apid > 0 {
		ps.Poles = unArchivedPoles
	}

	return apid
}

// AppendPolesFrom appends Poles from nps PoleSite to receiver Polesite
func (ps *PoleSite) AppendPolesFrom(nps *PoleSite) {
	poleId := -100
	for _, pole := range ps.Poles {
		if pole.Id > poleId {
			poleId = pole.Id
		}
	}

	for _, pole := range nps.Poles {
		poleId++
		npole := *pole
		npole.Id = poleId
		ps.Poles = append(ps.Poles, &npole)
	}
}

// UpdateWith updates receiver PoleSite with information from given PoleSite by checking update on each poles.
// Outdated poles are ignored. Updated pole' timestamp is set to current time.
// Ignored Poles References List is return.
func (ps *PoleSite) UpdateWith(ups *PoleSite) []string {
	timeStamp := date.Now().TimeStamp()
	ignoreList := []string{}

	resPoleSite := *ups //shallow copy of updated PoleSite
	updatedPoles := make(map[int]*Pole)
	for _, updatedPole := range ups.Poles {
		updatedPoles[updatedPole.Id] = updatedPole
	}

	// browse existing poles
	sort.Sort(ps)
	for opId, origPole := range ps.Poles {
		updatedPole, exists := updatedPoles[origPole.Id]
		if !exists {
			// origPole is not described in ups => keep it as is
			continue
		}
		if updatedPole.IsEqual(origPole) {
			// origPole unchanged => keep it as is
			delete(updatedPoles, updatedPole.Id)
			continue
		}
		// updatedPole has changed compared to origPole
		if updatedPole.TimeStamp != origPole.TimeStamp {
			// updated pole is outdated compared to current pole info, let's ignore update and keep origPole
			ignoreList = append(ignoreList, updatedPole.ExtendedRef())
			delete(updatedPoles, updatedPole.Id)
			continue
		}
		// updatedPole has changed and is legit => update timstamp and replace origPole
		updatedPole.TimeStamp = timeStamp
		ps.Poles[opId] = updatedPole
		delete(updatedPoles, updatedPole.Id)
	}

	// browse remaining (unvisited) updatePoles (should be new poles (ie with id < 0))
	if len(updatedPoles) > 0 {
		newPoles := make([]*Pole, len(updatedPoles))
		npId := 0
		for updatedPoleId, updatedPole := range updatedPoles {
			if updatedPoleId > 0 {
				ignoreList = append(ignoreList, "Unexpected "+updatedPole.ExtendedRef())
				continue
			}
			if updatedPole.State == poleconst.StateDeleted {
				// ignore new added deleted pole
				continue
			}
			newPoles[npId] = updatedPole
			npId++
		}
		newPoles = newPoles[:npId] // remove unused preallocated newPoles entries
		// sort new added poles in creation order (descending id)
		resPoleSite.Poles = newPoles
		sort.Sort(sort.Reverse(&resPoleSite))
		for _, newPole := range newPoles {
			newPole.TimeStamp = timeStamp
			newPole.Id = ps.NextPoleId()
			ps.Poles = append(ps.Poles, newPole)
		}
	}
	resPoleSite.Poles = ps.Poles

	*ps = resPoleSite //shallow copy of updated PoleSite on receiver PoleSite

	return ignoreList
}

// MergeWith appends Poles from nps PoleSite to receiver Polesite.
//
// Poles found in nps already declared in ps are compared and updated according to following rules
//
// New Poles in nps (not declared in ps) are appened
//
// Deleted Poles from nps (that is exists only in ps) are canceled (with NoGoClient status)
func (ps *PoleSite) MergeWith(nps *PoleSite) []*ConsistencyMsg {
	sort.Sort(ps)
	nextID := ps.Poles[len(ps.Poles)-1].Id + 1
	getNextId := func() int {
		nextID++
		return nextID - 1
	}
	poleDistance := func(p1, p2 *Pole) float64 {
		dlat := (p1.Lat - p2.Lat) * 100000.0
		dlong := (p1.Long - p2.Long) * 100000.0
		return math.Sqrt(dlat*dlat + dlong*dlong)
	}

	psDict := ps.getExtendedRefDict()
	psGf := geoFencer{}
	psGf.SetPrecision(1)
	psGfDict, _ := psGf.GetGeoFenceDict(ps)

	npsDict := nps.getExtendedRefDict()

	resMsg := []*ConsistencyMsg{}
	resPoles := []*Pole{}

	// create sorted ExtRef keys from npsDict
	npsExtRefList := make([]string, len(npsDict))
	i := 0
	for extRef, _ := range npsDict {
		npsExtRefList[i] = extRef
		i++
	}
	sort.Strings(npsExtRefList)

	// browse npsDict for existing and new poles
	for _, nExtRef := range npsExtRefList {
		npole := npsDict[nExtRef]
		pole, found := psDict[nExtRef]
		if !found { // npole from nps is a new pole, add it
			npole.Id = getNextId()
			resPoles = append(resPoles, npole)
			// Check for coordinate consistency
			npoleGf := psGf.GetGeoFence(npole)
			if nearPoles, foundNearPoles := psGfDict[npoleGf]; foundNearPoles {
				resMsg = append(resMsg, &ConsistencyMsg{
					Category: "Warning",
					Msg:      "Added pole " + nExtRef + " next to existing",
					Poles:    nearPoles, // poles near added npole
				})
				// copy DICT and DA data to added pole
				if npole.DictRef == "" {
					npole.DictInfo = nearPoles[0].DictInfo
					npole.DictRef = nearPoles[0].DictRef
					npole.DictDate = nearPoles[0].DictDate
				}
				if npole.DaQueryDate == "" {
					npole.DaQueryDate = nearPoles[0].DaQueryDate
					npole.DaStartDate = nearPoles[0].DaStartDate
					npole.DaEndDate = nearPoles[0].DaEndDate
					npole.DaValidation = nearPoles[0].DaValidation
				}
				continue
			}
			resMsg = append(resMsg, &ConsistencyMsg{
				Category: "Info",
				Msg:      "Added pole " + nExtRef,
				Poles:    []*Pole{}, // add current and previous pole
			})
			continue
		}
		// npole from nps is also declared in ps
		delete(psDict, nExtRef) // remove its ref from psDict

		if pole.IsEquivalent(npole) {
			// no significant change, keep pole as is
			resPoles = append(resPoles, pole)
			continue
		}

		// npole brings some change in from of pole
		changedInfo := ""
		if pole.City == "" && npole.City != "" {
			pole.City = npole.City
			changedInfo += " City"
		}
		if npole.Address == "" && pole.Address != npole.Address {
			pole.Address = npole.Address
			changedInfo += " Address"
		}
		if pole.Material != npole.Material {
			pole.Material = npole.Material
			pole.Comment = npole.Comment
			changedInfo += " Material"
		}
		if pole.Height != npole.Height {
			pole.Height = npole.Height
			pole.Comment = npole.Comment
			changedInfo += " Height"
		}
		if npole.DtRef != "" && pole.DtRef != npole.DtRef {
			if pole.DictRef != "" {
				pole.Comment += "\nRéférence de DT mise a jour: " + pole.DtRef + " => " + npole.DtRef + ". DICT à renouveller"
			}
			changedInfo += " DT"
			pole.DtRef = npole.DtRef
		}
		distance := poleDistance(pole, npole)
		if distance > 3.0 { // pole position as change > 3m
			pole.Long = npole.Long
			pole.Lat = npole.Lat
			changedInfo += " Position"
			if pole.DictRef != "" {
				pole.Comment += fmt.Sprintf("\nAppui déplacé de %.1fm, vérifier l'emprise de la DICT", distance)
			}
		}
		if pole.State != npole.State {
			if pole.IsDone() {
				if !npole.IsTodo() {
					pole.Comment += fmt.Sprint("\nAppui annulé par le client après réalisation")
					resMsg = append(resMsg, &ConsistencyMsg{
						Category: "Warning",
						Msg:      "Pole " + nExtRef + "was cancelled after being done",
						Poles:    []*Pole{}, // add current and previous pole
					})
				}
			} else {
				if (pole.IsTodo() && !npole.IsTodo()) || (!pole.IsTodo() && npole.IsTodo()) {
					changedInfo += " Status (" + pole.State + " => " + npole.State + ")"
					pole.State = npole.State
				}
			}
		}

		resPoles = append(resPoles, pole)
		if changedInfo != "" {
			resMsg = append(resMsg, &ConsistencyMsg{
				Category: "Info",
				Msg:      "Updated pole " + nExtRef + ":" + changedInfo,
				Poles:    []*Pole{}, // add current and previous pole
			})
		}
	}

	// create sorted ExtRef keys from npsDict
	psExtRefList := make([]string, len(psDict))
	i = 0
	for extRef, _ := range psDict {
		psExtRefList[i] = extRef
		i++
	}
	sort.Strings(psExtRefList)

	// browse psDict for deleted poles
	for _, extRef := range psExtRefList {
		pole := psDict[extRef]
		resMsg = append(resMsg, &ConsistencyMsg{
			Category: "Info",
			Msg:      "Removed pole " + extRef,
			Poles:    []*Pole{pole}, // add current and previous pole
		})
	}

	ps.Poles = resPoles
	sort.Sort(ps)

	return resMsg
}

func (ps *PoleSite) getExtendedRefDict() map[string]*Pole {
	res := make(map[string]*Pole)
	for _, pole := range ps.Poles {
		if pole.State == poleconst.StateDeleted {
			continue
		}
		res[pole.ExtendedRef()] = pole
	}
	return res
}

// PoleSite Consistency Control

type ConsistencyMsg struct {
	Category string
	Msg      string
	Poles    []*Pole
}

func (m ConsistencyMsg) String() string {
	res := m.Category
	res += ": " + m.Msg + "\n"
	for _, pole := range m.Poles {
		res += fmt.Sprintf("\tPole (%3d) %s : %s\n", pole.Id, pole.ExtendedRef(), pole.State)
	}

	return res
}

func (m ConsistencyMsg) ShortString() string {
	res := m.Msg
	if len(m.Poles) > 0 {
		poleNames := make([]string, len(m.Poles))
		for i, pole := range m.Poles {
			poleNames[i] = pole.ExtendedRef()
		}
		res += " (" + strings.Join(poleNames, ", ") + ")"
	}
	return res
}

func (ps *PoleSite) CheckPoleSiteConsistency() []*ConsistencyMsg {
	res := []*ConsistencyMsg{}
	res = append(res, ps.checkPoleId()...)
	res = append(res, ps.checkDuplicatedRef()...)
	res = append(res, ps.checkDuplicatedLocation(100000)...)
	return res
}

func (ps *PoleSite) checkPoleId() []*ConsistencyMsg {
	res := []*ConsistencyMsg{}
	if len(ps.Poles) <= 1 {
		return res
	}

	sort.Sort(ps)

	curId := ps.Poles[0].Id
	for numPole, pole := range ps.Poles[1:] {
		if pole.Id == curId {
			res = append(res, &ConsistencyMsg{
				Category: "PoleId",
				Msg:      "Duplicated Pole Id",
				Poles:    []*Pole{pole, ps.Poles[numPole]}, // add current and previous pole
			})
		}
		curId = pole.Id
	}
	return res
}

func (ps *PoleSite) checkDuplicatedRef() []*ConsistencyMsg {
	res := []*ConsistencyMsg{}
	refDict := make(map[string][]*Pole)
	duplicateFound := false
	for _, pole := range ps.Poles {
		if pole.State == poleconst.StateDeleted {
			continue
		}
		ref := pole.ExtendedRef()
		if refDict[ref] == nil {
			refDict[ref] = []*Pole{pole}
		} else {
			refDict[ref] = append(refDict[ref], pole)
			duplicateFound = true
		}
	}
	if !duplicateFound {
		return res
	}
	refs := []string{}
	for ref, poles := range refDict {
		if len(poles) < 2 {
			continue
		}
		refs = append(refs, ref)
	}
	sort.Strings(refs)
	for _, ref := range refs {
		res = append(res, &ConsistencyMsg{
			Category: "Ref",
			Msg:      "Duplicated Pole Ref: " + ref,
			Poles:    refDict[ref],
		})
	}
	return res
}

func (ps *PoleSite) checkDuplicatedLocation(prec float64) []*ConsistencyMsg {
	res := []*ConsistencyMsg{}
	g := geoFencer{prec: prec}
	posDict, duplicateFound := g.GetGeoFenceDict(ps)
	if !duplicateFound {
		return res
	}
	for _, poles := range posDict {
		if len(poles) < 2 {
			continue
		}
		res = append(res, &ConsistencyMsg{
			Category: "Location",
			Msg:      "Duplicated Pole Location",
			Poles:    poles,
		})
	}
	return res
}

// geoFence type and methods -------------------------------------------------------------------------------------------

type geoFence struct {
	lat, long int
}

type geoFencer struct {
	prec float64
}

func (g *geoFencer) SetPrecision(precMeter float64) {
	g.prec = 100000.0 / precMeter
}

func (g geoFencer) GetGeoFence(pole *Pole) geoFence {
	return geoFence{
		lat:  int(pole.Lat * g.prec),
		long: int(pole.Long * g.prec),
	}
}

func (g geoFencer) GetGeoFenceDict(ps *PoleSite) (map[geoFence][]*Pole, bool) {
	posDict := make(map[geoFence][]*Pole)
	duplicateFound := false
	for _, pole := range ps.Poles {
		if pole.State == poleconst.StateDeleted {
			continue
		}
		geofence := g.GetGeoFence(pole)
		if posDict[geofence] == nil {
			posDict[geofence] = []*Pole{pole}
		} else {
			posDict[geofence] = append(posDict[geofence], pole)
			duplicateFound = true
		}
	}
	return posDict, duplicateFound
}
