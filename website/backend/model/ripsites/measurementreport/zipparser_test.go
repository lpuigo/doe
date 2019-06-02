package measurementreport

import (
	"fmt"
	"os"
	"testing"
)

const (
	testZipFile string = "test/PM9 V2.zip"
)

func TestParseZipMeasurementFiles(t *testing.T) {
	f, err := os.Open(testZipFile)
	if err != nil {
		t.Fatalf("Open file returned: %v", err)
	}

	fstat, err := f.Stat()
	if err != nil {
		t.Fatalf("Stat returned: %v", err)
	}

	mr, err := ParseZipMeasurementFiles(f, fstat.Size())
	if err != nil {
		t.Fatalf("ParseZipMeasurementFiles returned: %v", err)
	}
	fmt.Printf("Find %d measurments\n", len(mr))
	for _, m := range mr {
		fmt.Printf("%s\n", m.String())
	}
}
