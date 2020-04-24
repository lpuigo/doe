package xlsextract

import (
	"fmt"
	"testing"
)

func TestPoleRecord_GetImage(t *testing.T) {
	parser, err := NewParserFile(testXlsFile)
	if err != nil {
		t.Fatalf("ParseFile returned unexpected: %s\n", err.Error())
	}

	if !parser.Next() {
		t.Fatalf("could not get row %4d: %s\n", parser.GetRowNum(), parser.Error())
	}
	rec, err := parser.ParseRecord()
	if err != nil {
		t.Logf("could not parse row %4d: %s\n", parser.GetRowNum(), err)
	}

	for i, label := range rec.GetImageLabels() {
		fmt.Printf("%02d: Get Image '%s'\n", i, label)
		err := rec.GetImage(testXlsExtractDir, label)
		if err != nil {
			t.Logf("GetImage returned unexpected: %s\n", err.Error())
		}
	}
}

func TestPoleRecord_GetAllImages(t *testing.T) {
	parser, err := NewParserFile(testXlsFile)
	if err != nil {
		t.Fatalf("ParseFile returned unexpected: %s\n", err.Error())
	}

	if !parser.Next() {
		t.Fatalf("could not get row %4d: %s\n", parser.GetRowNum(), parser.Error())
	}
	rec, err := parser.ParseRecord()
	if err != nil {
		t.Logf("could not parse row %4d: %s\n", parser.GetRowNum(), err)
	}

	err = rec.GetAllImages(testXlsExtractDir, 3)
	if err != nil {
		t.Logf("could not parse row %4d: %s\n", parser.GetRowNum(), err)
	}
}
