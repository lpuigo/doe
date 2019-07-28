package polesites

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	"github.com/lpuig/ewin/doe/website/backend/model/doctemplate"
	"github.com/lpuig/ewin/doe/website/frontend/model/polesite/poleconst"
)

const (
	rowPolesiteHeader int = 0
	rowPolesiteInfo   int = 1
	rowPoleHeader     int = 3
	rowPoleInfo       int = 4
)

func ToXLS(w io.Writer, ps *PoleSite) error {
	xf := excelize.NewFile()
	sheetName := ps.Ref
	xf.SetSheetName(xf.GetSheetName(1), sheetName)

	// Set PoleSite infos
	xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPolesiteHeader, 0), "polesite")
	xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPolesiteHeader, 1), "Id")
	xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPolesiteHeader, 2), "Client")
	xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPolesiteHeader, 3), "Ref")
	xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPolesiteHeader, 4), "Manager")
	xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPolesiteHeader, 5), "OrderDate")
	xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPolesiteHeader, 6), "Status")
	xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPolesiteHeader, 7), "Comment")

	xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPolesiteInfo, 0), "polesite")
	xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPolesiteInfo, 1), ps.Id)
	xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPolesiteInfo, 2), ps.Client)
	xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPolesiteInfo, 3), ps.Ref)
	xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPolesiteInfo, 4), ps.Manager)
	xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPolesiteInfo, 5), date.DateFrom(ps.OrderDate).ToTime())
	xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPolesiteInfo, 6), ps.Status)
	xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPolesiteInfo, 7), ps.Comment)

	// Set Poles infos
	xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleHeader, 0), "pole")
	xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleHeader, 1), "Ref")
	xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleHeader, 2), "City")
	xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleHeader, 3), "Address")
	xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleHeader, 4), "Lat")
	xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleHeader, 5), "Long")
	xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleHeader, 6), "State")
	xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleHeader, 7), "Actors")
	xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleHeader, 8), "Date")
	xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleHeader, 9), "DtRef")
	xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleHeader, 10), "DictRef")
	xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleHeader, 11), "DictInfo")
	xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleHeader, 12), "Height")
	xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleHeader, 13), "Material")
	xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleHeader, 14), poleconst.ProductCoated)
	xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleHeader, 15), poleconst.ProductMoise)
	xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleHeader, 16), poleconst.ProductReplace)
	xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleHeader, 17), poleconst.ProductRemove)

	for i, pole := range ps.Poles {
		xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleInfo+i, 0), "pole")
		xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleInfo+i, 1), pole.Ref)
		xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleInfo+i, 2), pole.City)
		xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleInfo+i, 3), pole.Address)
		xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleInfo+i, 4), pole.Lat)
		xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleInfo+i, 5), pole.Long)
		xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleInfo+i, 6), pole.State)
		xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleInfo+i, 7), "") // Actor ("Pierre, Paul, Jacques")
		xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleInfo+i, 8), "") // Date
		xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleInfo+i, 9), pole.DtRef)
		xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleInfo+i, 10), pole.DictRef)
		xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleInfo+i, 11), pole.DictInfo)
		xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleInfo+i, 12), pole.Height)
		xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleInfo+i, 13), pole.Material)
		products := map[string]int{}
		for _, product := range pole.Product {
			products[product] = 1
		}
		xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleInfo+i, 14), products[poleconst.ProductCoated])
		xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleInfo+i, 15), products[poleconst.ProductMoise])
		xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleInfo+i, 16), products[poleconst.ProductReplace])
		xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleInfo+i, 17), products[poleconst.ProductRemove])
	}

	err := xf.Write(w)
	if err != nil {
		return fmt.Errorf("could not write XLS file:%s", err.Error())
	}
	return nil
}

func FromXLS(r io.Reader) (*PoleSite, error) {
	xf, err := excelize.OpenReader(r)
	if err != nil {
		return nil, err
	}
	sheetName := xf.GetSheetName(1)
	// Read PoleSite Header & Info
	if err := checkValue(xf, sheetName, doctemplate.RcToAxis(rowPolesiteHeader, 0), "polesite"); err != nil {
		return nil, err
	}
	if err := checkValue(xf, sheetName, doctemplate.RcToAxis(rowPolesiteInfo, 0), "polesite"); err != nil {
		return nil, err
	}

	var ps *PoleSite = &PoleSite{}
	ps.Id, err = getCellInt(xf, sheetName, doctemplate.RcToAxis(rowPolesiteInfo, 1))
	if err != nil {
		return nil, err
	}
	ps.Client = xf.GetCellValue(sheetName, doctemplate.RcToAxis(rowPolesiteInfo, 2))
	ps.Ref = xf.GetCellValue(sheetName, doctemplate.RcToAxis(rowPolesiteInfo, 3))
	ps.Manager = xf.GetCellValue(sheetName, doctemplate.RcToAxis(rowPolesiteInfo, 4))
	ps.OrderDate, err = getCellDate(xf, sheetName, doctemplate.RcToAxis(rowPolesiteInfo, 5))
	if err != nil {
		return nil, err
	}
	ps.Status = xf.GetCellValue(sheetName, doctemplate.RcToAxis(rowPolesiteInfo, 6))
	ps.Comment = xf.GetCellValue(sheetName, doctemplate.RcToAxis(rowPolesiteInfo, 7))

	// Read Poles Header & Info
	if err := checkValue(xf, sheetName, doctemplate.RcToAxis(rowPoleHeader, 0), "pole"); err != nil {
		return nil, err
	}
	productKeys := map[int]string{}
	for _, col := range []int{14, 15, 16, 17} {
		productKeys[col] = xf.GetCellValue(sheetName, doctemplate.RcToAxis(rowPoleHeader, col))
	}
	if err := checkValue(xf, sheetName, doctemplate.RcToAxis(rowPoleInfo, 0), "pole"); err != nil {
		return ps, nil
	}

	for line, row := range xf.GetRows(sheetName) {
		if line < rowPoleInfo {
			continue
		}
		if row[0] != "pole" {
			continue
		}

		lat, err := strconv.ParseFloat(row[4], 64)
		if err != nil {
			return nil, fmt.Errorf("could not parse latitude '%s' row %d: %s", row[4], line+1, err.Error())
		}
		long, err := strconv.ParseFloat(row[5], 64)
		if err != nil {
			return nil, fmt.Errorf("could not parse longitude '%s' row %d: %s", row[5], line+1, err.Error())
		}
		height, err := strconv.Atoi(row[12])
		if err != nil {
			return nil, fmt.Errorf("could not parse height '%s' row %d: %s", row[12], line+1, err.Error())
		}
		pdate := ""
		if row[8] != "" {
			tdate, err := time.Parse("01-02-06", row[8])
			if err != nil {
				tdate, err = time.Parse("1/2/06 15:04", row[8])
				if err != nil {
					return nil, fmt.Errorf("could not parse date '%s' row %d: %s", row[8], line+1, err.Error())
				}
			}
			pdate = date.Date(tdate).String()
		}
		actors := []string{}
		if row[7] != "" {
			actors = strings.Split(row[7], ",")
			for i, actor := range actors {
				actors[i] = strings.Trim(actor, " ")
			}
		}

		pole := &Pole{
			Ref:      row[1],
			City:     row[2],
			Address:  row[3],
			Lat:      lat,
			Long:     long,
			State:    row[6],
			Date:     pdate,
			Actors:   actors,
			DtRef:    row[9],
			DictRef:  row[10],
			DictInfo: row[11],
			Height:   height,     // row 12
			Material: row[13],    // row 13
			Product:  []string{}, // row 10
		}
		for _, col := range []int{14, 15, 16, 17} {
			if row[col] == "1" {
				pole.Product = append(pole.Product, productKeys[col])
			}
		}

		ps.Poles = append(ps.Poles, pole)
	}

	return ps, nil
}

func checkValue(xf *excelize.File, sheetname, axis, value string) error {
	foundValue := xf.GetCellValue(sheetname, axis)
	if foundValue != value {
		return fmt.Errorf("misformated XLS file: cell %s!%s should contain '%s' (found '%s' instead)",
			sheetname, axis,
			foundValue, value,
		)
	}
	return nil
}

func getCellInt(xf *excelize.File, sheetname, axis string) (int, error) {
	val, err := strconv.Atoi(xf.GetCellValue(sheetname, axis))
	if err != nil {
		return 0, fmt.Errorf("misformated XLS file: cell %s!%s should contain int value", sheetname, axis)
	}
	return val, nil
}

func getCellFloat(xf *excelize.File, sheetname, axis string) (float64, error) {
	val, err := strconv.ParseFloat(xf.GetCellValue(sheetname, axis), 64)
	if err != nil {
		return 0, fmt.Errorf("misformated XLS file: cell %s!%s should contain float value", sheetname, axis)
	}
	return val, nil
}

func getCellDate(xf *excelize.File, sheetname, axis string) (string, error) {
	foundValue := xf.GetCellValue(sheetname, axis)
	foundDate, err := time.Parse("01-02-06", foundValue)
	if err != nil {
		foundDate, err = time.Parse("1/2/06 15:04", foundValue)
		if err != nil {
			return "", fmt.Errorf("misformated XLS file: cell %s!%s should contain date value ('%s' found instead): %s", sheetname, axis, foundValue, err.Error())
		}
	}
	return date.Date(foundDate).String(), nil
}
