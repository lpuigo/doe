package foasites

import (
	"encoding/json"
	"os"
	"testing"
)

const (
	testXLSFile string = `C:\Users\Laurent\Google Drive (laurent.puig.ewin@gmail.com)\Axians\Axians Moselle\Chantier FOA\2019-11-07 CCPHVA_TRE\TABLEAU SUIVI.xlsx`
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
	jsonEnconder := json.NewEncoder(os.Stdout)
	jsonEnconder.SetIndent("", "\t")
	err = jsonEnconder.Encode(nfs)
	if err != nil {
		t.Fatalf("FoaSite.Endode returned unexpected error: %v", err)
	}
}
