package xlsextract

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"
)

const (
	kizeoXlsExtractDir  string = `C:\Users\Laurent\OneDrive\Documents\TEMPORAIRE\Fiitelcom\Kizeo Vittel\`
	kizeoXlsExtractName string = `Poteau_Fiitelcom_20210128`
	report              string = `Synthese Kizeo 20210128 all`
	kizeoCreatePoleDir  bool   = true
	kizeoXlsExtractFile string = kizeoXlsExtractDir + kizeoXlsExtractName + ".xlsx"
	kizeoXlsReportFile  string = kizeoXlsExtractDir + report + ".xlsx"
)

func Test_CreateReport(t *testing.T) {
	recs, dupFound, err := ParseFile(kizeoXlsExtractFile)
	if err != nil {
		t.Fatalf("ParseFile returned unexpected: %s\n", err.Error())
	}

	if dupFound {
		t.Logf("Warning: duplicate(s) found in Kizeo Report")
	}

	err = WriteXlsReport(kizeoXlsReportFile, recs)
	if err != nil {
		t.Fatalf("WriteXlsReport returned unexpected: %s\n", err.Error())
	}
}

func Test_CheckDuplicate(t *testing.T) {
	recs, err := ReadXlsReportFromFile(kizeoXlsReportFile)
	if err != nil {
		t.Fatalf("ParseFile returned unexpected: %s\n", err.Error())
	}

	checkDuplicate(recs, t)
}

func Test_ExtractReport(t *testing.T) {
	recs, err := ReadXlsReportFromFile(kizeoXlsReportFile)
	if err != nil {
		t.Fatalf("ParseFile returned unexpected: %s\n", err.Error())
	}

	checkDuplicate(recs, t)

	for _, rec := range recs {
		if !rec.ExtractPicture {
			continue
		}
		dir := filepath.Join(kizeoXlsExtractDir, report, rec.GetSafeSRO())
		if kizeoCreatePoleDir {
			dir = filepath.Join(dir, rec.GetSafeRef())
		}
		ensure(dir, t)
		fmt.Printf("%s: Get %s %s\n", time.Now().Format("2006-01-02 15:04:05.0"), rec.SRO, rec.Ref)
		err := rec.GetAllImages(dir, 3)
		if err != nil {
			t.Fatalf("GetAllImages returned unexpected: %s\n", err.Error())
		}
		err = rec.WriteComment(dir)
		if err != nil {
			t.Fatalf("WriteComment returned unexpected: %s\n", err.Error())
		}
	}
}

func ensure(dir string, t *testing.T) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Printf("%s: creating '%s'\n", time.Now().Format("2006-01-02 15:04:05.0"), dir)
		err := os.MkdirAll(dir, os.ModeDir)
		if err != nil {
			t.Fatalf("unable to create '%s': %s", dir, err.Error())
		}
	}
}

func checkDuplicate(recs []*PoleRecord, t *testing.T) {
	sort.Slice(recs, func(i, j int) bool {
		if recs[i].SRO != recs[j].SRO {
			return recs[i].SRO < recs[j].SRO
		}
		if recs[i].Ref != recs[j].Ref {
			return recs[i].Ref < recs[j].Ref
		}
		dateI := recs[i].Date + " " + recs[i].Hour
		dateJ := recs[j].Date + " " + recs[j].Hour
		return dateI < dateJ
	})

	dictRefs := make(map[string]int)

	duplicateFound := false

	for _, rec := range recs {
		// check for duplicate
		sroref := rec.GetSRORef()
		dictRefs[sroref]++
		nb := dictRefs[sroref]
		if nb > 1 {
			sroref += fmt.Sprintf(" doublon %d", nb-1)
			t.Logf("found %s", sroref)
		}
	}
	if duplicateFound {
		t.Fatalf("Duplicate check failed")
	}
}
