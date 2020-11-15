package polesites

import (
	"archive/zip"
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strings"

	"github.com/lpuig/ewin/doe/kizeoparser/xlsextract"
	"github.com/lpuig/ewin/doe/website/backend/model/bpu"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
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

	Poles []*Pole
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
		case poleconst.StateToDo:
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
		}
	}
	return
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
	return res, nil
}

// AddStat adds Stats into values for given Polesite
//func (ps *PoleSite) AddStat(stats items.Stats, sc items.StatContext,
//	actorById clients.ActorById, currentBpu *bpu.Bpu, showprice bool) error {
//
//	addValue := func(date, serie string, actors []string, value float64) {
//		stats.AddStatValue(ps.Ref, ps.Client, date, "", serie, value)
//		if sc.ShowTeam && len(actors) > 0 {
//			value /= float64(len(actors))
//			for _, actName := range actors {
//				stats.AddStatValue(ps.Ref, ps.Client+" : "+actName, date, "", serie, value)
//			}
//		}
//	}
//
//	calcItems, err := ps.Itemize(currentBpu)
//	if err != nil {
//		return fmt.Errorf("error on polesite stat itemize for '%s':%s", ps.Ref, err.Error())
//	}
//	for _, item := range calcItems {
//		if !item.Done {
//			continue
//		}
//		actorsName := make([]string, len(item.Actors))
//		for i, actId := range item.Actors {
//			actorsName[i] = actorById(actId)
//		}
//		addValue(sc.DateFor(item.Date), items.StatSerieWork, actorsName, item.Work())
//		if showprice {
//			addValue(sc.DateFor(item.Date), items.StatSeriePrice, actorsName, item.Price())
//		}
//	}
//	return nil
//}

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
	zw := zip.NewWriter(w)

	path := strings.TrimSuffix(ps.DictZipName(), ".zip")

	makeDir := func(base ...string) string {
		return filepath.Join(base...) + "/"
	}

	// Create sorted List of DICT in PoleSite
	dicts := map[string]int{}
	for _, pole := range ps.Poles {
		dict := strings.Trim(pole.DictRef, " \t")
		if dict == "" {
			continue
		}
		dicts[pole.DictRef]++
	}
	dictList := make([]string, len(dicts))
	i := 0
	for dict, _ := range dicts {
		dictList[i] = dict
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
	origPoles := make(map[int]*Pole)
	for _, pole := range ps.Poles {
		origPoles[pole.Id] = pole
	}
	updtPoles := []*Pole{}
	for _, updtPole := range ups.Poles {
		origPole, exists := origPoles[updtPole.Id]
		if !exists {
			// its a new pole => let's add it with new timestamp
			updtPole.TimeStamp = timeStamp
			updtPoles = append(updtPoles, updtPole)
			continue
		}
		if updtPole.IsEqual(origPole) {
			// unchanged => let's add it unchanged
			updtPoles = append(updtPoles, origPole)
			continue
		}
		// some change occured
		if updtPole.TimeStamp != origPole.TimeStamp {
			// updated pole is outdated compared to current pole info, let's ignore update and keep origPole
			updtPoles = append(updtPoles, origPole)
			ignoreList = append(ignoreList, updtPole.ExtendedRef())
			continue
		}
		updtPole.TimeStamp = timeStamp
		updtPoles = append(updtPoles, updtPole)
	}
	resPoleSite.Poles = updtPoles
	*ps = resPoleSite //shallow copy of updated PoleSite on receiver PoleSite

	return ignoreList
}

type KizeoReport struct {
	NbUpdate   int
	UnknownRef []string
}

// UpdateFromKizeoPoleRecords updates kizeo data pertaining to poles described in given PoleRecord slice.
// Its retuns a KizeoReport describing number of updated pole and list of non matching Pole Ref (eg. ref not found in reciever PoleSite)
func (ps *PoleSite) UpdateFromKizeoPoleRecords(prs []*xlsextract.PoleRecord) *KizeoReport {
	timeStamp := date.Now().TimeStamp()

	report := &KizeoReport{
		NbUpdate:   0,
		UnknownRef: []string{},
	}
	poleDict := make(map[string]*Pole)
	for _, pole := range ps.Poles {
		poleDict[pole.ExtendedRef()] = pole
	}
	for _, pr := range prs {
		ref := pr.SRO + " " + pr.Ref
		pole, found := poleDict[ref]
		if !found {
			report.UnknownRef = append(report.UnknownRef, ref)
			continue
		}
		pole.Kizeo = "OK " + pr.Date + " " + pr.Hour
		pole.TimeStamp = timeStamp
		report.NbUpdate++
	}

	return report
}

// PoleSite Consistency Control

type ConsistencyMsg struct {
	Category string
	Msg      string
	Poles    []*Pole
}

func (m ConsistencyMsg) String() string {
	res := m.Category
	res += ":" + m.Msg + "\n"
	for _, pole := range m.Poles {
		res += fmt.Sprintf("\tPole (%3d) %s : %s\n", pole.Id, pole.ExtendedRef(), pole.State)
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
	sort.Slice(ps.Poles, func(i, j int) bool {
		return ps.Poles[i].Id < ps.Poles[j].Id
	})

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
	type geoFence struct {
		lat, long int
	}
	posDict := make(map[geoFence][]*Pole)
	duplicateFound := false
	for _, pole := range ps.Poles {
		geofence := geoFence{
			lat:  int(pole.Lat * prec),
			long: int(pole.Long * prec),
		}
		if posDict[geofence] == nil {
			posDict[geofence] = []*Pole{pole}
		} else {
			posDict[geofence] = append(posDict[geofence], pole)
			duplicateFound = true
		}
	}
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
