package actors

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	"github.com/lpuig/ewin/doe/website/backend/tools/xlsx"
	"io"
	"strconv"
	"strings"
)

const (
	colActorImport int = iota
	colActorId
	colActorLastName
	colActorFirstName
	colActorRole
	colActorCompany
	colActorContract
	colActorHireDate
	colActorLeaveDate
	colActorClient
	colActorVacation
)

func FromXLS(r io.Reader) ([]*Actor, error) {
	xf, err := excelize.OpenReader(r)
	if err != nil {
		return nil, err
	}
	sheetName := xf.GetSheetName(1)

	actors := []*Actor{}

	for line, row := range xf.GetRows(sheetName) {
		if line == 0 {
			// skip header
			continue
		}
		if row[colActorImport] == "" {
			// not to be processed, skip
			continue
		}

		id, err := getCellInt(xf, sheetName, xlsx.RcToAxis(line, colActorId))
		if err != nil {
			return nil, err
		}

		getValue := func(col int) string {
			if col >= len(row) {
				return ""
			}
			return strings.Trim(row[col], " \t")
		}

		fName := strings.Title(getValue(colActorFirstName))
		lName := strings.ToUpper(getValue(colActorLastName))

		vacs := []date.DateRange{}
		for c := colActorVacation; c < len(row); c += 2 {
			beg := getValue(c)
			if beg == "" {
				break
			}
			vacs = append(vacs, date.DateRange{
				Begin: beg,
				End:   getValue(c + 1),
			})
		}

		actor := &Actor{
			Id:        id,
			Ref:       lName + " " + fName,
			FirstName: fName,
			LastName:  lName,
			Period: date.DateRange{
				Begin: getValue(colActorHireDate),
				End:   getValue(colActorLeaveDate),
			},
			Company:  strings.ToTitle(getValue(colActorCompany)),
			Contract: getValue(colActorContract),
			Role:     getValue(colActorRole),
			Vacation: vacs,
			Client:   getValue(colActorClient),
		}

		actors = append(actors, actor)
	}
	return actors, nil
}

func getCellInt(xf *excelize.File, sheetname, axis string) (int, error) {
	val, err := strconv.Atoi(xf.GetCellValue(sheetname, axis))
	if err != nil {
		return 0, fmt.Errorf("misformated XLS file: cell %s!%s should contain int value", sheetname, axis)
	}
	return val, nil
}
