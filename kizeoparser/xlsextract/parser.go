package xlsextract

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
)

func NewParserFile(xlsfile string) (*xlsPoleParser, error) {
	xlsf, err := excelize.OpenFile(xlsfile)
	if err != nil {
		return nil, err
	}

	sheet := xlsf.GetSheetName(0)
	parser, err := newXlsParser(xlsf, sheet)
	if err != nil {
		return nil, err
	}

	err = parser.ParseHeader()
	if err != nil {
		return nil, err
	}

	//parser.PrintColumnNames()
	return parser, nil
}

func ParseFile(xlsfile string) ([]*PoleRecord, bool, error) {
	parser, err := NewParserFile(xlsfile)
	if err != nil {
		return nil, false, err
	}
	dictRefs := make(map[string]int)
	res := []*PoleRecord{}
	dupFound := false
	for parser.Next() {
		rec, err := parser.ParseRecord()
		if err != nil {
			return nil, false, fmt.Errorf("could not parse row %4d: %s\n", parser.rowNum, err)
		}

		// check for duplicate
		sroref := rec.GetSRORef()
		dictRefs[sroref]++
		nb := dictRefs[sroref]
		if nb > 1 {
			rec.Ref += fmt.Sprintf(" doublon %d", nb-1)
			dupFound = true
		}
		res = append(res, rec)
	}
	return res, dupFound, nil
}
