package foasites

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"io"
)

const (
	// index are 0 based here
	rowStart int = 3
	colRef   int = 1
	colInsee int = 2
	colType  int = 3
)

// NewFoaSiteFromXLS reads an FollowUp XLSx from input reader and returns a new filled-in FoaSite
func NewFoaSiteFromXLS(r io.Reader) (*FoaSite, error) {
	xf, err := excelize.OpenReader(r)
	if err != nil {
		return nil, err
	}
	sheetName := xf.GetSheetName(1)
	if sheetName == "" {
		return nil, fmt.Errorf("could not get XLS sheetname")
	}

	nfs := NewFoaSite()

	rows, err := xf.GetRows(sheetName)
	if err != nil {
		return nil, err
	}
	for _, row := range rows[rowStart:] {
		ref, insee, typ := row[colRef], row[colInsee], row[colType]
		if ref == "" && insee == "" && typ == "" {
			// we meet an empty line ... assuming end of data
			break
		}
		nFoa := NewFoa()
		nFoa.Ref, nFoa.Insee, nFoa.Type = ref, insee, typ
		nfs.AddFoa(nFoa)
	}
	return nfs, nil
}
