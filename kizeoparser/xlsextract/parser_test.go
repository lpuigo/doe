package xlsextract

import (
	"fmt"
	"testing"
)

func TestParseFile(t *testing.T) {
	res, dupFound, err := ParseFile(testXlsFile)
	if err != nil {
		t.Fatalf("ParseFile returned unexpected: %s\n", err.Error())
	}

	if dupFound {
		t.Logf("Warning: duplicate(s) found in Kizeo Report")
	}

	for i, poleRecord := range res {
		fmt.Printf("%4d: %s", i, poleRecord.String())
	}
}

func TestNewParserFile(t *testing.T) {
	parser, err := NewParserFile(testXlsFile)
	if err != nil {
		t.Fatalf("ParseFile returned unexpected: %s\n", err.Error())
	}

	for parser.Next() {
		rec, err := parser.ParseRecord()
		if err != nil {
			t.Logf("could not parse row %4d: %s\n", parser.rowNum, err)
		}
		fmt.Printf("%4d: %s", parser.GetRowNum(), rec.String())
	}
}
