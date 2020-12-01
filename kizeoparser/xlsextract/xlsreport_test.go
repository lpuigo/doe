package xlsextract

import (
	"fmt"
	"testing"
)

func TestWriteXlsReport(t *testing.T) {
	recs, _, err := ParseFile(testXlsFile)
	if err != nil {
		t.Fatalf("ParseFile returned unexpected: %s\n", err.Error())
	}

	err = WriteXlsReport(testXlsReport, recs)
	if err != nil {
		t.Fatalf("WriteXlsReport returned unexpected: %s\n", err.Error())
	}
}

func TestReadXlsReport(t *testing.T) {
	recs, err := ReadXlsReportFromFile(testXlsReport)
	if err != nil {
		t.Fatalf("ReadXlsReport returned unexpected: %s\n", err.Error())
	}

	for i, poleRecord := range recs {
		fmt.Printf("%4d: %s", i, poleRecord.String())
	}
}
