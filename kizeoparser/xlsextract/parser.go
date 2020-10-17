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

	parser.PrintColumnNames()
	return parser, nil
}

func ParseFile(xlsfile string) ([]*PoleRecord, error) {
	parser, err := NewParserFile(xlsfile)
	if err != nil {
		return nil, err
	}
	res := []*PoleRecord{}
	for parser.Next() {
		rec, err := parser.ParseRecord()
		if err != nil {
			return nil, fmt.Errorf("could not parse row %4d: %s\n", parser.rowNum, err)
		}
		res = append(res, rec)
	}
	return res, nil
}
