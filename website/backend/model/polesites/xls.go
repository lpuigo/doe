package polesites

import (
	"fmt"
	"github.com/lpuig/ewin/doe/website/backend/tools/nominatim"
	"github.com/lpuig/ewin/doe/website/backend/tools/xlsx"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	"github.com/lpuig/ewin/doe/website/frontend/model/polesite/poleconst"
)

const (
	rowPolesiteHeader int = 0
	rowPolesiteInfo   int = 1
	rowPoleHeader     int = 3
	rowPoleInfo       int = 4
)

const (
	colPoleId int = iota
	colPoleRef
	colPoleCity
	colPoleAddress
	colPoleLat
	colPoleLong
	colPoleState
	colPoleActors
	colPoleDate
	colPoleAttachmentDate
	colPoleSticker
	colPoleDtRef
	colPoleDictRef
	colPoleDictDate
	colPoleDictInfo
	colPoleHeight
	colPoleMaterial
	colPoleAspiDate
	colPoleKizeo
	colPoleComment
	colPoleProduct
)

func ToXLS(w io.Writer, ps *PoleSite) error {
	xf := excelize.NewFile()
	sheetName := ps.Ref
	xf.SetSheetName(xf.GetSheetName(1), sheetName)

	// Set PoleSite infos
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteHeader, 0), "polesite")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteHeader, 1), "Id")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteHeader, 2), "Client")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteHeader, 3), "Ref")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteHeader, 4), "Manager")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteHeader, 5), "OrderDate")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteHeader, 6), "Status")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteHeader, 7), "Comment")

	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteInfo, 0), "polesite")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteInfo, 1), ps.Id)
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteInfo, 2), ps.Client)
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteInfo, 3), ps.Ref)
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteInfo, 4), ps.Manager)
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteInfo, 5), date.DateFrom(ps.OrderDate).ToTime())
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteInfo, 6), ps.Status)
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteInfo, 7), ps.Comment)

	// Set Poles infos
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleId), "pole")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleRef), "Ref")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleCity), "City")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleAddress), "Address")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleLat), "Lat")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleLong), "Long")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleState), "State")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleActors), "Actors")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleDate), "Date")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleAttachmentDate), "AttachmentDate")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleSticker), "Sticker")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleDtRef), "DtRef")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleDictRef), "DictRef")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleDictDate), "DictDate")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleDictInfo), "DictInfo")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleHeight), "Height")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleMaterial), "Material")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleAspiDate), "AspiDate")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleKizeo), "Kizeo")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleComment), "Comment")

	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleProduct+0), poleconst.ProductCoated)
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleProduct+1), poleconst.ProductMoise)
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleProduct+2), poleconst.ProductCouple)
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleProduct+3), poleconst.ProductReplace)
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleProduct+4), poleconst.ProductRemove)

	for i, pole := range ps.Poles {
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleInfo+i, colPoleId), "pole")
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleInfo+i, colPoleRef), pole.Ref)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleInfo+i, colPoleCity), pole.City)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleInfo+i, colPoleAddress), pole.Address)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleInfo+i, colPoleLat), pole.Lat)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleInfo+i, colPoleLong), pole.Long)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleInfo+i, colPoleState), pole.State)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleInfo+i, colPoleActors), "")      // Actor ("Pierre, Paul, Jacques")
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleInfo+i, colPoleDate), pole.Date) // Date
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleInfo+i, colPoleAttachmentDate), pole.AttachmentDate)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleInfo+i, colPoleSticker), pole.Sticker)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleInfo+i, colPoleDtRef), pole.DtRef)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleInfo+i, colPoleDictRef), pole.DictRef)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleInfo+i, colPoleDictDate), pole.DictDate)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleInfo+i, colPoleDictInfo), pole.DictInfo)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleInfo+i, colPoleHeight), pole.Height)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleInfo+i, colPoleMaterial), pole.Material)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleInfo+i, colPoleAspiDate), pole.AspiDate)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleInfo+i, colPoleKizeo), pole.Kizeo)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleInfo+i, colPoleComment), pole.Comment)
		products := map[string]string{}
		for _, product := range pole.Product {
			products[product] = "1"
		}
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleInfo+i, colPoleProduct+0), products[poleconst.ProductCoated])
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleInfo+i, colPoleProduct+1), products[poleconst.ProductMoise])
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleInfo+i, colPoleProduct+2), products[poleconst.ProductCouple])
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleInfo+i, colPoleProduct+3), products[poleconst.ProductReplace])
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleInfo+i, colPoleProduct+4), products[poleconst.ProductRemove])
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
	//
	// Read PoleSite Header & Info
	//
	if err := checkValue(xf, sheetName, xlsx.RcToAxis(rowPolesiteHeader, 0), "polesite"); err != nil {
		return nil, err
	}
	if err := checkValue(xf, sheetName, xlsx.RcToAxis(rowPolesiteInfo, 0), "polesite"); err != nil {
		return nil, err
	}

	var ps *PoleSite = &PoleSite{}
	ps.Id, err = getCellInt(xf, sheetName, xlsx.RcToAxis(rowPolesiteInfo, 1))
	if err != nil {
		return nil, err
	}
	ps.Client = xf.GetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteInfo, 2))
	ps.Ref = xf.GetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteInfo, 3))
	ps.Manager = xf.GetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteInfo, 4))
	ps.OrderDate, err = getCellDate(xf, sheetName, xlsx.RcToAxis(rowPolesiteInfo, 5))
	if err != nil {
		return nil, err
	}
	ps.Status = xf.GetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteInfo, 6))
	ps.Comment = xf.GetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteInfo, 7))

	//
	// Read Poles Header & Info
	//
	if err := checkValue(xf, sheetName, xlsx.RcToAxis(rowPoleHeader, 0), "pole"); err != nil {
		return nil, err
	}
	productKeys := map[int]string{}
	for _, col := range []int{colPoleProduct, colPoleProduct + 1, colPoleProduct + 2, colPoleProduct + 3, colPoleProduct + 4} {
		productKeys[col] = xf.GetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, col))
	}
	if err := checkValue(xf, sheetName, xlsx.RcToAxis(rowPoleInfo, colPoleId), "pole"); err != nil {
		return ps, nil
	}

	cellCoord := func(col, row int) string {
		return sheetName + "!" + xlsx.RcToAxis(row, col)
	}

	id := 0
	for line, row := range xf.GetRows(sheetName) {
		if line < rowPoleInfo {
			continue
		}
		if row[colPoleId] != "pole" {
			continue
		}

		lat, errlat := strconv.ParseFloat(row[colPoleLat], 64)
		long, errlong := strconv.ParseFloat(row[colPoleLong], 64)
		geomsg := ""
		if errlat != nil && errlong != nil {
			// Perform Geoloc Search
			addr := row[colPoleAddress]
			res := []nominatim.Geoloc{}
			if addr == "" {
				goto GeolocDone
			}
			res, err = nominatim.GeolocSearch(addr)
			if err != nil {
				geomsg = "ERR Geoloc:" + err.Error()
			}
			if len(res) == 0 {
				geomsg = "Geoloc not found"
				goto GeolocDone
			}
			lat, long, err = res[0].GetLatLong()
			if err != nil {
				geomsg = "ERR Geoloc:" + err.Error()
			}
		GeolocDone:
		} else {
			if errlat != nil {
				return nil, fmt.Errorf("%s: could not parse latitude '%s': %s", cellCoord(colPoleLat, line), row[colPoleLat], errlat.Error())
			}
			if errlong != nil {
				return nil, fmt.Errorf("%s: could not parse longitude '%s': %s", cellCoord(colPoleLong, line), row[colPoleLong], errlong.Error())
			}
		}
		state := row[colPoleState]
		if state == "" {
			state = poleconst.StateNotSubmitted
		}
		height := 8
		if row[colPoleHeight] != "" {
			height, err = strconv.Atoi(row[colPoleHeight])
			if err != nil {
				return nil, fmt.Errorf("%s: could not parse height '%s': %s", cellCoord(colPoleHeight, line), row[colPoleHeight], err.Error())
			}
		}
		pdate, err := parseDate(row[colPoleDate])
		if err != nil {
			return nil, fmt.Errorf("%s: %s", cellCoord(colPoleDate, line), err.Error())
		}
		adate, err := parseDate(row[colPoleAttachmentDate])
		if err != nil {
			return nil, fmt.Errorf("%s: %s", cellCoord(colPoleAttachmentDate, line), err.Error())
		}
		aspdate, err := parseDate(row[colPoleAspiDate])
		if err != nil {
			return nil, fmt.Errorf("%s: %s", cellCoord(colPoleAspiDate, line), err.Error())
		}
		dddate, err := parseDate(row[colPoleDictDate])
		if err != nil {
			return nil, fmt.Errorf("%s: %s", cellCoord(colPoleDictDate, line), err.Error())
		}
		actors := []string{}
		if row[colPoleActors] != "" {
			actors = strings.Split(row[colPoleActors], ",")
			for i, actor := range actors {
				actors[i] = strings.Trim(actor, " ")
			}
		}
		products := []string{}
		for _, col := range []int{colPoleProduct, colPoleProduct + 1, colPoleProduct + 2, colPoleProduct + 3, colPoleProduct + 4} {
			if row[col] == "1" {
				products = append(products, productKeys[col])
			}
		}

		comment := row[colPoleComment]
		if geomsg != "" {
			comment += "\n" + geomsg
		}

		pole := &Pole{
			Id:             id,
			Ref:            row[colPoleRef],
			City:           row[colPoleCity],
			Address:        row[colPoleAddress],
			Lat:            lat,
			Long:           long,
			State:          state,
			Actors:         actors,
			Date:           pdate,
			AttachmentDate: adate,
			Sticker:        row[colPoleSticker],
			DtRef:          row[colPoleDtRef],
			DictRef:        row[colPoleDictRef],
			DictDate:       dddate,
			DictInfo:       row[colPoleDictInfo],
			Height:         height,
			Material:       row[colPoleMaterial],
			AspiDate:       aspdate,
			Kizeo:          row[colPoleKizeo],
			Comment:        comment,
			Product:        products,
		}

		ps.Poles = append(ps.Poles, pole)
		id++
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

func parseDate(source string) (string, error) {
	pdate := ""
	if source != "" {
		tdate, err := time.Parse("2006-01-02", source)
		if err != nil {
			tdate, err = time.Parse("01-02-06", source)
			if err != nil {
				tdate, err = time.Parse("1/2/06 15:04", source)
				if err != nil {
					return "", fmt.Errorf("could not parse date '%s': %s", source, err.Error())
				}
			}
		}
		pdate = date.Date(tdate).String()
	}
	return pdate, nil
}
