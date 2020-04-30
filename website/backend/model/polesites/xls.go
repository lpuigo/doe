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
	rowPolesiteHeader int = 1
	rowPolesiteInfo   int = 2
	rowPoleHeader     int = 4
	rowPoleInfo       int = 5
)

const (
	colPoleId int = iota + 1
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
	xf.SetSheetName(xf.GetSheetName(0), sheetName)

	// Set PoleSite infos
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteHeader, 1), "polesite")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteHeader, 2), "Id")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteHeader, 3), "Client")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteHeader, 4), "Ref")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteHeader, 5), "Manager")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteHeader, 6), "OrderDate")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteHeader, 7), "Status")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteHeader, 8), "Comment")

	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteInfo, 1), "polesite")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteInfo, 2), ps.Id)
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteInfo, 3), ps.Client)
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteInfo, 4), ps.Ref)
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteInfo, 5), ps.Manager)
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteInfo, 6), date.DateFrom(ps.OrderDate).ToTime())
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteInfo, 7), ps.Status)
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteInfo, 8), ps.Comment)

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
	sheetName := xf.GetSheetName(0)
	//
	// Read PoleSite Header & Info
	//
	if err := checkValue(xf, sheetName, xlsx.RcToAxis(rowPolesiteHeader, 1), "polesite"); err != nil {
		return nil, err
	}
	if err := checkValue(xf, sheetName, xlsx.RcToAxis(rowPolesiteInfo, 1), "polesite"); err != nil {
		return nil, err
	}

	var ps *PoleSite = &PoleSite{}
	ps.Id, err = getCellInt(xf, sheetName, xlsx.RcToAxis(rowPolesiteInfo, 2))
	if err != nil {
		return nil, err
	}
	ps.Client, _ = xf.GetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteInfo, 3))
	ps.Ref, _ = xf.GetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteInfo, 4))
	ps.Manager, _ = xf.GetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteInfo, 5))
	ps.OrderDate, err = getCellDate(xf, sheetName, xlsx.RcToAxis(rowPolesiteInfo, 6))
	if err != nil {
		return nil, err
	}
	ps.Status, _ = xf.GetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteInfo, 7))
	ps.Comment, _ = xf.GetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteInfo, 8))

	//
	// Read Poles Header & Info
	//
	if err := checkValue(xf, sheetName, xlsx.RcToAxis(rowPoleHeader, 1), "pole"); err != nil {
		return nil, err
	}
	productKeys := map[int]string{}
	for _, col := range []int{colPoleProduct, colPoleProduct + 1, colPoleProduct + 2, colPoleProduct + 3, colPoleProduct + 4} {
		productKeys[col], _ = xf.GetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, col))
	}
	if err := checkValue(xf, sheetName, xlsx.RcToAxis(rowPoleInfo, colPoleId), "pole"); err != nil {
		return ps, nil
	}

	cellCoord := func(col, row int) string {
		return sheetName + "!" + xlsx.RcToAxis(row, col)
	}

	id := 0

	rows, err := xf.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	line := 0
	row := []string{}
	getCol := func(col int) string {
		if col > len(row) {
			return ""
		}
		return row[col-1]
	}

	for line, row = range rows {
		if line+1 < rowPoleInfo {
			continue
		}
		if getCol(colPoleId) != "pole" {
			continue
		}

		lat, errlat := strconv.ParseFloat(getCol(colPoleLat), 64)
		long, errlong := strconv.ParseFloat(getCol(colPoleLong), 64)
		geomsg := ""
		if errlat != nil && errlong != nil {
			// Perform Geoloc Search
			addr := getCol(colPoleAddress)
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
				return nil, fmt.Errorf("%s: could not parse latitude '%s': %s", cellCoord(colPoleLat, line), getCol(colPoleLat), errlat.Error())
			}
			if errlong != nil {
				return nil, fmt.Errorf("%s: could not parse longitude '%s': %s", cellCoord(colPoleLong, line), getCol(colPoleLong), errlong.Error())
			}
		}
		state := getCol(colPoleState)
		if state == "" {
			state = poleconst.StateNotSubmitted
		}
		height := 8
		if getCol(colPoleHeight) != "" {
			height, err = strconv.Atoi(getCol(colPoleHeight))
			if err != nil {
				return nil, fmt.Errorf("%s: could not parse height '%s': %s", cellCoord(colPoleHeight, line), getCol(colPoleHeight), err.Error())
			}
		}
		pdate, err := parseDate(getCol(colPoleDate))
		if err != nil {
			return nil, fmt.Errorf("%s: %s", cellCoord(colPoleDate, line), err.Error())
		}
		adate, err := parseDate(getCol(colPoleAttachmentDate))
		if err != nil {
			return nil, fmt.Errorf("%s: %s", cellCoord(colPoleAttachmentDate, line), err.Error())
		}
		aspdate, err := parseDate(getCol(colPoleAspiDate))
		if err != nil {
			return nil, fmt.Errorf("%s: %s", cellCoord(colPoleAspiDate, line), err.Error())
		}
		dddate, err := parseDate(getCol(colPoleDictDate))
		if err != nil {
			return nil, fmt.Errorf("%s: %s", cellCoord(colPoleDictDate, line), err.Error())
		}
		actors := []string{}
		if getCol(colPoleActors) != "" {
			actors = strings.Split(getCol(colPoleActors), ",")
			for i, actor := range actors {
				actors[i] = strings.Trim(actor, " ")
			}
		}
		products := []string{}
		for _, col := range []int{colPoleProduct, colPoleProduct + 1, colPoleProduct + 2, colPoleProduct + 3, colPoleProduct + 4} {
			if getCol(col) == "1" {
				products = append(products, productKeys[col])
			}
		}

		comment := getCol(colPoleComment)
		if geomsg != "" {
			comment += "\n" + geomsg
		}

		pole := &Pole{
			Id:             id,
			Ref:            getCol(colPoleRef),
			City:           getCol(colPoleCity),
			Address:        getCol(colPoleAddress),
			Lat:            lat,
			Long:           long,
			State:          state,
			Actors:         actors,
			Date:           pdate,
			AttachmentDate: adate,
			Sticker:        getCol(colPoleSticker),
			DtRef:          getCol(colPoleDtRef),
			DictRef:        getCol(colPoleDictRef),
			DictDate:       dddate,
			DictInfo:       getCol(colPoleDictInfo),
			Height:         height,
			Material:       getCol(colPoleMaterial),
			AspiDate:       aspdate,
			Kizeo:          getCol(colPoleKizeo),
			Comment:        comment,
			Product:        products,
		}

		ps.Poles = append(ps.Poles, pole)
		id++
	}

	return ps, nil
}

func checkValue(xf *excelize.File, sheetname, axis, value string) error {
	foundValue, err := xf.GetCellValue(sheetname, axis)
	if err != nil {
		return err
	}
	if foundValue != value {
		return fmt.Errorf("misformated XLS file: cell %s!%s should contain '%s' (found '%s' instead)",
			sheetname, axis,
			foundValue, value,
		)
	}
	return nil
}

func getCellInt(xf *excelize.File, sheetname, axis string) (int, error) {
	foundValue, err := xf.GetCellValue(sheetname, axis)
	if err != nil {
		return 0, err
	}
	val, err := strconv.Atoi(foundValue)
	if err != nil {
		return 0, fmt.Errorf("misformated XLS file: cell %s!%s should contain int value", sheetname, axis)
	}
	return val, nil
}

func getCellFloat(xf *excelize.File, sheetname, axis string) (float64, error) {
	foundValue, err := xf.GetCellValue(sheetname, axis)
	if err != nil {
		return 0, err
	}
	val, err := strconv.ParseFloat(foundValue, 64)
	if err != nil {
		return 0, fmt.Errorf("misformated XLS file: cell %s!%s should contain float value", sheetname, axis)
	}
	return val, nil
}

func getCellDate(xf *excelize.File, sheetname, axis string) (string, error) {
	foundValue, _ := xf.GetCellValue(sheetname, axis)
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
