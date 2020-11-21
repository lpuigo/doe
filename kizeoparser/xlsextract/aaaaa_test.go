package xlsextract

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const (
	kizeoXlsExtractDir  string = `C:\Users\Laurent\Desktop\TEMPORAIRE\Fiitelcom\Kizeo Vittel\`
	kizeoXlsExtractName string = `Poteau_Fiitelcom_20201121`
	report              string = `Synthese Kizeo Vittel`
	kizeoCreatePoleDir  bool   = true
	kizeoXlsExtractFile string = kizeoXlsExtractDir + kizeoXlsExtractName + ".xlsx"
	kizeoXlsReportFile  string = kizeoXlsExtractDir + report + ".xlsx"
)

func Test_CreateReport(t *testing.T) {
	recs, err := ParseFile(kizeoXlsExtractFile)
	if err != nil {
		t.Fatalf("ParseFile returned unexpected: %s\n", err.Error())
	}

	err = WriteXlsReport(kizeoXlsReportFile, recs)
	if err != nil {
		t.Fatalf("WriteXlsReport returned unexpected: %s\n", err.Error())
	}
}

func Test_ExtractReport(t *testing.T) {
	recs, err := ReadXlsReportFromFile(kizeoXlsReportFile)
	if err != nil {
		t.Fatalf("ParseFile returned unexpected: %s\n", err.Error())
	}

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
