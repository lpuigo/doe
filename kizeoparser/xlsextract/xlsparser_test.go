package xlsextract

import "testing"

func TestXlsPoleParser_UnmergeHeaderCells(t *testing.T) {
	xlsfile := `C:\Users\Laurent\Desktop\TEMPORAIRE\Eiffage Signes\2020-10-16 Extract Kizeo\Poteau_Eiffage_Signes_20201016.xlsx`
	parser, err := NewParserFile(xlsfile)
	if err != nil {
		t.Fatalf("NewParserFile returns unexpected %v", err)
	}

	err = parser.UnmergeHeaderCells()
	if err != nil {
		t.Fatalf("UnmergeHeaderCells returns unexpected %v", err)
	}

}
