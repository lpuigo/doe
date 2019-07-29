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
	xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleHeader, 9), "AttachmentDate")
	xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleHeader, 10), "AttachmentDate")
	xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleHeader, 11), "DtRef")
	xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleHeader, 12), "DictRef")
	xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleHeader, 13), "DictInfo")
	xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleHeader, 14), "Height")
	xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleHeader, 15), "Material")
	xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleHeader, 16), "AspiDate")
	xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleHeader, 17), "Kizeo")
	xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleHeader, 18), "Comment")
	xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleHeader, 19), poleconst.ProductCoated)
	xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleHeader, 20), poleconst.ProductMoise)
	xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleHeader, 21), poleconst.ProductReplace)
	xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleHeader, 22), poleconst.ProductRemove)

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
		xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleInfo+i, 9), pole.AttachmentDate)
		xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleInfo+i, 10), pole.Sticker)
		xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleInfo+i, 11), pole.DtRef)
		xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleInfo+i, 12), pole.DictRef)
		xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleInfo+i, 13), pole.DictInfo)
		xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleInfo+i, 14), pole.Height)
		xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleInfo+i, 15), pole.Material)
		xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleInfo+i, 16), pole.AspiDate)
		xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleInfo+i, 17), pole.Kizeo)
		xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleInfo+i, 18), pole.Comment)
		products := map[string]string{}
		for _, product := range pole.Product {
			products[product] = "1"
		}
		xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleInfo+i, 19), products[poleconst.ProductCoated])
		xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleInfo+i, 20), products[poleconst.ProductMoise])
		xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleInfo+i, 21), products[poleconst.ProductReplace])
		xf.SetCellValue(sheetName, doctemplate.RcToAxis(rowPoleInfo+i, 22), products[poleconst.ProductRemove])
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
	for _, col := range []int{19, 20, 21, 22} {
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
		height, err := strconv.Atoi(row[14])
		if err != nil {
			return nil, fmt.Errorf("could not parse height '%s' row %d: %s", row[12], line+1, err.Error())
		}
		pdate, err := parseDate(row[8], line)
		if err != nil {
			return nil, err
		}
		adate, err := parseDate(row[9], line)
		if err != nil {
			return nil, err
		}
		aspdate, err := parseDate(row[16], line)
		if err != nil {
			return nil, err
		}
		actors := []string{}
		if row[7] != "" {
			actors = strings.Split(row[7], ",")
			for i, actor := range actors {
				actors[i] = strings.Trim(actor, " ")
			}
		}
		products := []string{}
		for _, col := range []int{19, 20, 21, 22} {
			if row[col] == "1" {
				products = append(products, productKeys[col])
			}
		}

		pole := &Pole{
			Ref:            row[1],   // row 1
			City:           row[2],   // row 2
			Address:        row[3],   // row 3
			Lat:            lat,      // row 4
			Long:           long,     // row 5
			State:          row[6],   // row 6
			Actors:         actors,   // row 7
			Date:           pdate,    // row 8
			AttachmentDate: adate,    // row 9
			Sticker:        row[10],  // row 10
			DtRef:          row[11],  // row 11
			DictRef:        row[12],  // row 12
			DictInfo:       row[13],  // row 13
			Height:         height,   // row 14
			Material:       row[15],  // row 15
			AspiDate:       aspdate,  // row 16
			Kizeo:          row[17],  // row 17
			Comment:        row[18],  // row 18
			Product:        products, // row 19-22
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

func parseDate(source string, line int) (string, error) {
	pdate := ""
	if source != "" {
		tdate, err := time.Parse("01-02-06", source)
		if err != nil {
			tdate, err = time.Parse("1/2/06 15:04", source)
			if err != nil {
				return "", fmt.Errorf("could not parse date '%s' row %d: %s", source, line+1, err.Error())
			}
		}
		pdate = date.Date(tdate).String()
	}
	return pdate, nil
}
