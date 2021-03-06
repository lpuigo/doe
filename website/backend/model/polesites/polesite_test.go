package polesites

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	"github.com/lpuig/ewin/doe/website/backend/model/polesites/test"
	"github.com/lpuig/ewin/doe/website/frontend/model/polesite/poleconst"
	"github.com/lpuig/ewin/doe/website/frontend/tools/latlong"
)

// =====================================================================================================================
// Utilitary Test functions (PoleSite creation or edition) =============================================================

// TestPolesiteFromXLS : convert an XLSx PoleSite Description to its JSON file
func TestPolesiteFromXLS(t *testing.T) {
	psXlsfile := `C:\Users\Laurent\Google Drive (laurent.puig.ewin@gmail.com)\Sogeca\Chantier Poteaux\2021-06-14 Maj SRO 67\Polesite Sogeca Poteau-SRO 67 (65).xlsx`

	path := filepath.Dir(psXlsfile)
	inFile := filepath.Base(psXlsfile)
	suffix := filepath.Ext(inFile)
	outFile := strings.TrimSuffix(inFile, suffix) + " gps" + suffix

	xf, err := os.Open(psXlsfile)
	if err != nil {
		t.Fatalf("could not open file: %s", err.Error())
	}

	ps, err := FromXLS(xf)
	if err != nil {
		t.Fatalf("FromXLS return unexpected: %s", err.Error())
	}

	PoleSiteConsistency(t, ps, false)

	xfr, err := os.Create(filepath.Join(path, outFile))
	if err != nil {
		t.Fatalf("could not create file: %s", err.Error())
	}
	defer xfr.Close()
	err = ToExportXLS(xfr, ps)
	if err != nil {
		t.Fatalf("ToExportXLS return unexpected: %s", err.Error())
	}

	jsonOutFile := filepath.Join(path, fmt.Sprintf("%06d.json", ps.Id))
	ojf, err := os.Create(jsonOutFile)
	if err != nil {
		t.Fatalf("could not create file: %s", err.Error())
	}
	defer ojf.Close()
	je := json.NewEncoder(ojf)
	err = je.Encode(ps)
	if err != nil {
		t.Fatalf("could not encode polesite: %s", err.Error())
	}
}

// TestPoleSite_AppendXlsToJson : add a XLSx PoleSite description to an already provided PoleSite JSON file.
// New Poles are appended next to already existing ones
func TestPoleSite_AppendXlsToJson(t *testing.T) {
	psDir := `C:\Users\Laurent\Google Drive (laurent.puig.ewin@gmail.com)\Sogetrel\Chantiers Poteaux Ouest\2021-04-13 Ordre traitement appuis`
	psXlsfile := `Polesite Eiffage Signes-Peyrolles.xlsx`
	psJsonFile := `000055.json`

	// Get original PoleSiteRecord
	origPsrFile := filepath.Join(psDir, psJsonFile)
	origPsr, err := NewPoleSiteRecordFromFile(origPsrFile)
	if err != nil {
		t.Fatalf("NewPoleSiteRecordFromFile returned unexpected: %s", err.Error())
	}
	// Get new XLS file and create new PoleSite
	xf, err := os.Open(filepath.Join(psDir, psXlsfile))
	if err != nil {
		t.Fatalf("could not open file: %s", err.Error())
	}

	newPs, err := FromXLS(xf)
	if err != nil {
		t.Fatalf("FromXLS returned unexpected: %s", err.Error())
	}

	// append new Poles to original PoleSiteRecord
	origPsr.AppendPolesFrom(newPs)

	// Controle result polesite consistency (detects poles with same name, position, ids)
	PoleSiteConsistency(t, origPsr.PoleSite, true)

	err = os.Rename(origPsrFile, origPsrFile+".bak")
	if err != nil {
		t.Fatalf("Rename origPsrFile returned unexpected: %s", err.Error())
	}
	// Persist updated PoleSite
	err = origPsr.Persist(psDir)
	if err != nil {
		t.Fatalf("Persist returned unexpected: %s", err.Error())
	}
}

// TestPoleSite_MergeXlsToJson : merge a XLSx PoleSite description to an already provided PoleSite JSON file.
// New Poles are appended next to already existing ones. Updated Poles are changed. Pole not providied in XLSx file are deleted from target JSON file
func TestPoleSite_MergeXlsToJson(t *testing.T) {
	psDir := `C:\Users\Laurent\Google Drive (laurent.puig.ewin@gmail.com)\Eiffage\Eiffage Poteau Signes\Chantiers\2021-05-06 Maj Peyrolles`
	psXlsfile := `Polesite Eiffage Signes-Peyrolles v2.xlsx`
	psJsonFile := `000058.json`

	// Get original PoleSiteRecord
	origPsrFile := filepath.Join(psDir, psJsonFile)
	origPsr, err := NewPoleSiteRecordFromFile(origPsrFile)
	if err != nil {
		t.Fatalf("NewPoleSiteRecordFromFile returned unexpected: %s", err.Error())
	}
	// Get new XLS file and create new PoleSite
	xf, err := os.Open(filepath.Join(psDir, psXlsfile))
	if err != nil {
		t.Fatalf("could not open file: %s", err.Error())
	}

	newPs, err := FromXLS(xf)
	if err != nil {
		t.Fatalf("FromXLS returned unexpected: %s", err.Error())
	}

	// Merge Xlsx defeind new Poles with original PoleSiteRecord
	msgs := origPsr.MergeWith(newPs, false)

	for _, msg := range msgs {
		fmt.Printf("%s", msg.String())
	}

	// Controle result polesite consistency (detects poles with same name, position, ids)
	PoleSiteConsistency(t, origPsr.PoleSite, false)

	err = os.Rename(origPsrFile, origPsrFile+".bak")
	if err != nil {
		t.Fatalf("Rename origPsrFile returned unexpected: %s", err.Error())
	}
	// Persist updated PoleSite
	err = origPsr.Persist(psDir)
	if err != nil {
		t.Fatalf("Persist returned unexpected: %s", err.Error())
	}
}

// Polesite generation from directory containing Orange "Fiche Appui" xlsx files
func TestBrowseFicheAppuiXlsxFiles(t *testing.T) {
	faDir := `C:\Users\Laurent\OneDrive\Documents\TEMPORAIRE\Axians Alsace\2021-06-07 SRO 88-005\2021-06-15 Fiches manquantes`
	polesiteId := 62
	polesiteName := "NRO 88-005"
	refPrefix := ""

	ps := &PoleSite{
		Id:         polesiteId,
		Client:     "",
		Ref:        polesiteName,
		Manager:    "",
		OrderDate:  date.Today().String(),
		UpdateDate: "",
		Status:     "",
		Comment:    "",
		Poles:      []*Pole{},
	}

	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() {
			return nil
		}
		// dir & files to skip
		if !strings.HasSuffix(info.Name(), ".xlsx") {
			return nil
		}
		if !strings.HasPrefix(info.Name(), "FicheAppui") {
			return nil
		}
		// dir & files to process
		fmt.Printf("processessing: %s\n", path)
		if err = processFicheAppuiXlsxFile(path, faDir, refPrefix, ps); err != nil {
			fmt.Printf("error occured : %s\n", err.Error())
		}
		return nil
	}

	// Browse Directory for FicheAppui
	err := filepath.Walk(faDir, walkFunc)
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", faDir, err)
		return
	}

	// Create export Xlsx PoleSite File
	resultPsFile := filepath.Join(faDir, polesiteName+".xlsx")
	oxf, err := os.Create(resultPsFile)
	if err != nil {
		t.Fatalf("could not create result xlsx file '%s':%s", resultPsFile, err.Error())
	}

	err = ToExportXLS(oxf, ps)
	if err != nil {
		t.Fatalf("could not save xlsx file '%s':%s", resultPsFile, err.Error())
	}
}

func processFicheAppuiXlsxFile(path, root, refPrefix string, ps *PoleSite) error {
	xf, err := excelize.OpenFile(path)
	if err != nil {
		return err
	}
	pathSteps := strings.Split(strings.ReplaceAll(path, root+"\\", ""), "\\")
	pa := pathSteps[0]
	sheetName := xf.GetSheetName(0)
	sticker, _ := xf.GetCellValue(sheetName, "D3")
	city, _ := xf.GetCellValue(sheetName, "D4")
	address, _ := xf.GetCellValue(sheetName, "D5")
	latDeg, _ := xf.GetCellValue(sheetName, "P5")
	lat, _ := latlong.DegToDec(latDeg)
	longDeg, _ := xf.GetCellValue(sheetName, "P6")
	long, _ := latlong.DegToDec(longDeg)
	ope, _ := xf.GetCellValue(sheetName, "M53")
	target, _ := xf.GetCellValue(sheetName, "M52")

	if target != "" && !CheckCAPFTPoleInfo(target) { // Target is not a Pole type => assume ope and target were switched in the XLS file
		ope, target = target, ope
	}

	products := []string{}
	lowerOpe := strings.ToLower(ope)
	switch {
	case strings.Contains(lowerOpe, "renforc"):
		// no main product, additionnal product will be added by parsing mat column
	case strings.Contains(lowerOpe, "redress"):
		products = append(products, poleconst.ProductStraighten)
	case strings.Contains(lowerOpe, "recal"):
		products = append(products, poleconst.ProductStraighten)
	default:
		products = append(products, poleconst.ProductCreation)
		if strings.Contains(lowerOpe, "remplac") {
			products = append(products, poleconst.ProductReplace)
		}
	}

	mat := ""
	height := 0
	mat, height = DecodeCAPFTPoleInfo(target, &products)

	pole := &Pole{
		Id:             0,
		Ref:            refPrefix + pa,
		City:           strings.Trim(city, " \t"),
		Address:        strings.Trim(address, " \t"),
		Sticker:        strings.Trim(sticker, " \t"),
		Lat:            lat,
		Long:           long,
		State:          poleconst.StateDictToDo,
		Date:           "",
		Actors:         nil,
		DtRef:          "",
		DictRef:        "",
		DictDate:       "",
		DictInfo:       "",
		DaQueryDate:    "",
		DaValidation:   false,
		DaStartDate:    "",
		DaEndDate:      "",
		Height:         height,
		Material:       mat,
		AspiDate:       "",
		Kizeo:          "",
		Comment:        strings.Trim(ope+" "+target, " \t"),
		Product:        products,
		AttachmentDate: "",
		TimeStamp:      "",
	}
	ps.Poles = append(ps.Poles, pole)
	return nil
}

// =====================================================================================================================
// Utilitary Test functions (PoleSite Fixing) ==========================================================================
func PoleSiteConsistency(t *testing.T, ps *PoleSite, abortOnFailed bool) {
	consistencyMsgs := ps.CheckPoleSiteConsistency()
	if len(consistencyMsgs) > 0 {
		for _, msg := range consistencyMsgs {
			t.Logf("%s", msg.String())
		}
		if abortOnFailed {
			t.Fatalf("Consistency Check failed")
		}
		return
	}
	t.Logf("Consistency Check OK")
}

func TestPoleSite_FixInconsistentPoleId(t *testing.T) {
	psDir := `C:\Users\Laurent\Golang\src\github.com\lpuig\ewin\doe\Ressources\Polesites`
	psJsonFile := `000057.json`
	origPsrFile := filepath.Join(psDir, psJsonFile)

	origPsr, err := NewPoleSiteRecordFromFile(origPsrFile)
	if err != nil {
		t.Fatalf("NewPoleSiteRecordFromFile returned unexpected: %s", err.Error())
	}
	ps := origPsr.PoleSite
	sort.Slice(ps.Poles, func(i, j int) bool {
		return ps.Poles[i].Id < ps.Poles[j].Id
	})
	if ps.Poles[0].Id >= 0 {
		return
	}
	// negative Ids found, lets setup Poles Id
	t.Logf("Negative Id found for PoleSite %s\n", psJsonFile)
	for i := 0; i < len(ps.Poles); i++ {
		ps.Poles[i].Id = i
	}
	err = os.Rename(origPsrFile, origPsrFile+".bak")
	if err != nil {
		t.Fatalf("Rename origPsrFile returned unexpected: %s", err.Error())
	}
	// Persist updated PoleSite
	err = origPsr.Persist(psDir)
	if err != nil {
		t.Fatalf("Persist returned unexpected: %s", err.Error())
	}
}

func TestPoleSite_DetectDuplicatedId(t *testing.T) {
	psDir := `C:\Users\Laurent\Golang\src\github.com\lpuig\ewin\doe\Ressources\Polesites`
	psJsonFile := `000057.json`
	origPsrFile := filepath.Join(psDir, psJsonFile)

	origPsr, err := NewPoleSiteRecordFromFile(origPsrFile)
	if err != nil {
		t.Fatalf("NewPoleSiteRecordFromFile returned unexpected: %s", err.Error())
	}
	ps := origPsr.PoleSite

	// detect duplicated pole id
	polesById := make(map[int][]int)
	for i, pole := range ps.Poles {
		bucket, found := polesById[pole.Id]
		if !found {
			bucket = []int{i}
			polesById[pole.Id] = bucket
			continue
		}
		//t.Logf("found another pole id %d\n", pole.Id)
		polesById[pole.Id] = append(polesById[pole.Id], i)
	}

	for id, bucket := range polesById {
		if len(bucket) > 1 {
			//t.Logf("Found %d duplicate(s) for pole id %d\n", len(bucket)-1, id)

			refB := &strings.Builder{}
			poleRef := ps.Poles[bucket[0]]
			refBEncoder := json.NewEncoder(refB)
			refBEncoder.SetIndent("", "\t")
			refBEncoder.Encode(poleRef)
			ref := refB.String()
			for _, i := range bucket[1:] {
				otherB := &strings.Builder{}
				poleOther := ps.Poles[i]
				otherBEncoder := json.NewEncoder(otherB)
				otherBEncoder.SetIndent("", "\t")
				otherBEncoder.Encode(poleOther)
				other := otherB.String()
				if other == ref {
					t.Logf("poleid:%d #%d doublon de #%d => A supprimer\n", id, i, bucket[0])
					continue
				}
				if poleRef.ExtendedRef() == poleOther.ExtendedRef() {
					t.Logf("poleid:%d #%d en delta de #%d => A décider\nRef  :%sClone:%s\n", id, i, bucket[0], ref, other)
				} else {
					t.Logf("poleid:%d #%d différent de #%d => id a modifier\nRef  :%sClone:%s\n", id, i, bucket[0], ref, other)
				}
			}
		}
	}
}

func TestPoleSite_FixOrderedPoleId(t *testing.T) {
	psDir := `C:\Users\Laurent\Golang\src\github.com\lpuig\ewin\doe\Ressources\Polesites`
	psJsonFile := `000057.json`
	origPsrFile := filepath.Join(psDir, psJsonFile)

	origPsr, err := NewPoleSiteRecordFromFile(origPsrFile)
	if err != nil {
		t.Fatalf("NewPoleSiteRecordFromFile returned unexpected: %s", err.Error())
	}
	ps := origPsr.PoleSite
	sort.Slice(ps.Poles, func(i, j int) bool {
		return ps.Poles[i].Id < ps.Poles[j].Id
	})
	// Persist updated PoleSite
	err = origPsr.Persist(psDir)
	if err != nil {
		t.Fatalf("Persist returned unexpected: %s", err.Error())
	}
}

func TestPoleSite_ReOrderPoleId(t *testing.T) {
	psDir := `C:\Users\Laurent\Golang\src\github.com\lpuig\ewin\doe\Ressources\Polesites`
	psJsonFile := `000057.json`
	origPsrFile := filepath.Join(psDir, psJsonFile)

	origPsr, err := NewPoleSiteRecordFromFile(origPsrFile)
	if err != nil {
		t.Fatalf("NewPoleSiteRecordFromFile returned unexpected: %s", err.Error())
	}
	ps := origPsr.PoleSite
	sort.Slice(ps.Poles, func(i, j int) bool {
		if ps.Poles[i].Ref == ps.Poles[j].Ref {
			return ps.Poles[i].Sticker < ps.Poles[j].Sticker
		}
		return ps.Poles[i].Ref < ps.Poles[j].Ref
	})
	for i, pole := range ps.Poles {
		pole.Id = i
		pole.Sticker = strings.Trim(pole.Sticker, " ")
		pole.Ref = strings.Trim(pole.Ref, " ")
	}
	// Persist updated PoleSite
	err = origPsr.Persist(psDir)
	if err != nil {
		t.Fatalf("Persist returned unexpected: %s", err.Error())
	}
}

// =====================================================================================================================
// Actual Test function ================================================================================================

func TestPoleSiteEncode(t *testing.T) {
	ps := &PoleSite{
		Id:         0,
		Client:     "Axians Moselle",
		Ref:        "Rechicourt",
		Manager:    "Goeffrey Wecker",
		OrderDate:  "2019-06-15",
		UpdateDate: "2019-07-11",
		Status:     "20 InProgress",
		Comment:    "test",
		Poles:      nil,
	}

	pole_Enrobe := map[string]bool{
		"MF22005": true,
		"MF22009": true,
		"MF22053": true,
		"MF22054": true,
		"MF22072": true,
		"MF22106": true,
		"MF22143": true,
		"MF22146": true,
		"MF22154": true,
		"MF22158": true,
		"MF22174": true,
		"MF22176": true,
		"MF22182": true,
		"MF22183": true,
		"MF22185": true,
		"MF22187": true,
		"MF22245": true,
		"MF22246": true,
		"MF22256": true,
		"MF22260": true,
		"MF22798": true,
		"MF22807": true,
		"MF22820": true,
		"MF22821": true,
		"MF22829": true,
		"MF22830": true,
		"MF22834": true,
		"MF22836": true,
		"MF22837": true,
		"MF22839": true,
		"MF22840": true,
		"MF22841": true,
	}
	pole_9m := map[string]bool{
		"MF22219": true,
		"MF22268": true,
		"MF22272": true,
		"MF22104": true,
		"MF22017": true,
		"MF22003": true,
		"MF22659": true,
		"MF22300": true,
		"MF22011": true,
		"MF22663": true,
		"MF22136": true,
		"MF22186": true,
		"MF22140": true,
		"MF22851": true,
		"MF22833": true,
		"MF22991": true,
		"MF22804": true,
		"MF22773": true,
	}

	for _, bp := range test.Poles {
		p := &Pole{
			Ref:      bp.Ref,
			City:     bp.City,
			Address:  "",
			Lat:      bp.Lat,
			Long:     bp.Long,
			State:    bp.State,
			DtRef:    "",
			DictRef:  "",
			Height:   8,
			Material: poleconst.MaterialWood,
			Product:  []string{},
			DictInfo: "",
		}
		if pole_Enrobe[p.Ref] {
			p.Product = append(p.Product, poleconst.ProductCoated)
		}
		if pole_9m[p.Ref] {
			p.Height = 9
		}
		ps.Poles = append(ps.Poles, p)
	}

	json.NewEncoder(os.Stdout).Encode(ps)
}

func LoadPolesiteFromJSON(t *testing.T, jsonFile string) *PoleSite {
	psf, err := os.Open(jsonFile)
	if err != nil {
		t.Fatalf("could not open file '%s':%s", jsonFile, err.Error())
	}
	defer psf.Close()
	var ps PoleSite
	err = json.NewDecoder(psf).Decode(&ps)
	if err != nil {
		t.Fatalf("could not unmarshall Polesite file '%s':%s", jsonFile, err.Error())
	}

	return &ps
}

func Test_ToXLS(t *testing.T) {
	psfile := `C:\Users\Laurent\Golang\src\github.com\lpuig\ewin\doe\Ressources\Polesites\000000.json`
	ps := LoadPolesiteFromJSON(t, psfile)

	resultPsFile := filepath.Join("test", ps.Ref+".xlsx")
	oxf, err := os.Create(resultPsFile)
	if err != nil {
		t.Fatalf("could not create result xlsx file '%s':%s", resultPsFile, err.Error())
	}

	err = ToExportXLS(oxf, ps)
	if err != nil {
		t.Fatalf("could not save xlsx file '%s':%s", resultPsFile, err.Error())
	}
}
