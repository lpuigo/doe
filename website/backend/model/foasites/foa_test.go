package foasites

import (
	"encoding/json"
	"fmt"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	"os"
	"path/filepath"
	"testing"
)

const (
	testXLSFile    string = `C:\Users\Laurent\Google Drive (laurent.puig.ewin@gmail.com)\Axians\Axians Moselle\Chantier FOA\2020-01-15 CCCM_MAX\TABLEAU_SUIVI_FOA_CCCM_MAX.xlsx`
	testFoaSiteId  int    = 3
	testFoaSiteRef string = "CCCM_MAX"
	//testXLSFile string = `C:\Users\Laurent\Google Drive (laurent.puig.ewin@gmail.com)\Axians\Axians Moselle\Chantier FOA\2019-12-20 CCSMS_SAR\TABLEAU SUIVI FOA CCSMS_SAR 2.xlsx`
	//testXLSFile string = `C:\Users\Laurent\Google Drive (laurent.puig.ewin@gmail.com)\Axians\Axians Moselle\Chantier FOA\2019-11-07 CCPHVA_TRE\TABLEAU SUIVI.xlsx`
	//testXLSFile string = `C:\Users\Laurent\Google Drive (laurent.puig.ewin@gmail.com)\Axians\Axians Moselle\Chantier FOA\2019-11-28 CCSMS_SAR\TABLEAU_SUIVI_FOA.xlsx`
)

func TestNewFoaSiteFromXLS(t *testing.T) {
	f, err := os.Open(testXLSFile)
	if err != nil {
		t.Fatalf("could not open: %v", err)
	}
	defer f.Close()

	nfs, err := NewFoaSiteFromXLS(f)
	if err != nil {
		t.Fatalf("NewFoaSiteFromXLS returned unexpected error: %v", err)
	}
	nfs.Id = testFoaSiteId
	nfs.Client = "Axians Moselle Foa"
	nfs.Ref = testFoaSiteRef
	nfs.Manager = "BERNARD Kevin"
	nfs.Status = "20 InProgress"
	nfs.OrderDate = date.GetMonday(date.Today().String())

	dir := filepath.Dir(testXLSFile)
	file := fmt.Sprintf("%06d.json", nfs.Id)
	jsonFile, err := os.Create(filepath.Join(dir, file))
	if err != nil {
		t.Fatalf("Create result JSON file returned unexpected error: %v", err)
	}
	defer jsonFile.Close()

	jsonEnconder := json.NewEncoder(jsonFile)
	jsonEnconder.SetIndent("", "\t")
	err = jsonEnconder.Encode(nfs)
	if err != nil {
		t.Fatalf("FoaSite.Encode returned unexpected error: %v", err)
	}
}
