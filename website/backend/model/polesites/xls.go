package polesites

import (
	"fmt"
	"github.com/lpuig/ewin/doe/website/backend/tools/nominatim"
	"github.com/lpuig/ewin/doe/website/backend/tools/xlsx"
	"io"
	"sort"
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
	colPoleSticker
	colPoleCity
	colPoleAddress
	colPoleLat
	colPoleLong
	colPoleState
	colPoleActors
	colPoleDate
	colPoleAttachmentDate
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

func ToExportXLS(w io.Writer, ps *PoleSite) error {
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
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleSticker), "Sticker")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleCity), "City")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleAddress), "Address")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleLat), "Lat")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleLong), "Long")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleState), "State")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleActors), "Actors")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleDate), "Date")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleAttachmentDate), "AttachmentDate")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleDtRef), "DtRef")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleDictRef), "DictRef")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleDictDate), "DictDate")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleDictInfo), "DictInfo")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleHeight), "Height")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleMaterial), "Material")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleAspiDate), "AspiDate")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleKizeo), "Kizeo")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleComment), "Comment")

	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleProduct+0), poleconst.ProductCreation)
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleProduct+1), poleconst.ProductReplace)
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleProduct+2), poleconst.ProductStraighten)
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleProduct+3), poleconst.ProductCouple)
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleProduct+4), poleconst.ProductMoise)
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleProduct+5), poleconst.ProductHauban)
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleProduct+6), poleconst.ProductCoated)
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleProduct+7), poleconst.ProductRemove)
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleProduct+8), poleconst.ProductNoAccess)
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, colPoleProduct+9), poleconst.ProductDenseNetwork)

	for i, pole := range ps.Poles {
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleInfo+i, colPoleId), "pole")
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleInfo+i, colPoleRef), pole.Ref)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleInfo+i, colPoleSticker), pole.Sticker)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleInfo+i, colPoleCity), pole.City)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleInfo+i, colPoleAddress), pole.Address)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleInfo+i, colPoleLat), pole.Lat)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleInfo+i, colPoleLong), pole.Long)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleInfo+i, colPoleState), pole.State)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleInfo+i, colPoleActors), "")      // Actor ("Pierre, Paul, Jacques")
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleInfo+i, colPoleDate), pole.Date) // Date
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleInfo+i, colPoleAttachmentDate), pole.AttachmentDate)
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
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleInfo+i, colPoleProduct+0), products[poleconst.ProductCreation])
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleInfo+i, colPoleProduct+1), products[poleconst.ProductReplace])
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleInfo+i, colPoleProduct+2), products[poleconst.ProductStraighten])
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleInfo+i, colPoleProduct+3), products[poleconst.ProductCouple])
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleInfo+i, colPoleProduct+4), products[poleconst.ProductMoise])
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleInfo+i, colPoleProduct+5), products[poleconst.ProductHauban])
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleInfo+i, colPoleProduct+6), products[poleconst.ProductCoated])
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleInfo+i, colPoleProduct+7), products[poleconst.ProductRemove])
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleInfo+i, colPoleProduct+8), products[poleconst.ProductNoAccess])
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPoleInfo+i, colPoleProduct+9), products[poleconst.ProductDenseNetwork])
	}

	err := xf.Write(w)
	if err != nil {
		return fmt.Errorf("could not write XLS file:%s", err.Error())
	}
	return nil
}

func ToRefExportXLS(w io.Writer, ps *PoleSite) error {
	refSet := make(map[string]int)
	for _, pole := range ps.Poles {
		refSet[pole.Ref]++
	}

	refs := make([]string, len(refSet))
	i := 0
	for ref, _ := range refSet {
		refs[i] = ref
		i++
	}
	sort.Strings(refs)

	xf := excelize.NewFile()
	sheetName := ps.Ref
	xf.SetSheetName(xf.GetSheetName(0), sheetName)

	// Set PoleSite infos
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteHeader, 1), "Client")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteHeader, 2), "polesite")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteHeader, 3), "Ref")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteHeader, 4), "Nb Poles")

	for i, refName := range refs {
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteInfo+i, 1), ps.Client)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteInfo+i, 2), ps.Ref)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteInfo+i, 3), refName)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(rowPolesiteInfo+i, 4), refSet[refName])
	}

	err := xf.Write(w)
	if err != nil {
		return fmt.Errorf("could not write XLS file:%s", err.Error())
	}
	return nil
}

const (
	colProgressPoleRef int = iota + 1
	colProgressPoleSticker
	colProgressPoleCity
	colProgressPoleAddress
	colProgressPoleDate
	colProgressPoleHeight
	colProgressPoleMaterial
	colProgressPoleProduct
)

const (
	colProgressPoleRefWidth      float64 = 18
	colProgressPoleStickerWidth  float64 = 12
	colProgressPoleCityWidth     float64 = 25
	colProgressPoleAddressWidth  float64 = 40
	colProgressPoleDateWidth     float64 = 12
	colProgressPoleHeightWidth   float64 = 8
	colProgressPoleMaterialWidth float64 = 18
	colProgressPoleProductWidth  float64 = 60
)

func ToProgressXLS(w io.Writer, ps *PoleSite) error {
	progressPoles := []*Pole{}
	for _, pole := range ps.Poles {
		if !pole.IsDone() {
			continue
		}
		progressPoles = append(progressPoles, pole)
	}

	sort.Slice(progressPoles, func(i, j int) bool {
		if progressPoles[i].Date != progressPoles[j].Date {
			return progressPoles[i].Date < progressPoles[j].Date
		}
		if progressPoles[i].Ref != progressPoles[j].Ref {
			return progressPoles[i].Ref < progressPoles[j].Ref
		}
		return progressPoles[i].Sticker < progressPoles[j].Sticker
	})

	xf := excelize.NewFile()
	sheetName := ps.Ref
	xf.SetSheetName(xf.GetSheetName(0), sheetName)

	getColName := func(col int) string {
		colName, _ := excelize.ColumnNumberToName(col)
		return colName
	}

	// Set Cols width & Format
	colName := getColName(colProgressPoleRef)
	xf.SetColWidth(sheetName, colName, colName, colProgressPoleRefWidth)
	colName = getColName(colProgressPoleSticker)
	xf.SetColWidth(sheetName, colName, colName, colProgressPoleStickerWidth)
	colName = getColName(colProgressPoleCity)
	xf.SetColWidth(sheetName, colName, colName, colProgressPoleCityWidth)
	colName = getColName(colProgressPoleAddress)
	xf.SetColWidth(sheetName, colName, colName, colProgressPoleAddressWidth)
	colName = getColName(colProgressPoleDate)
	exp := "dd/mm/yyyy;@"
	style, _ := xf.NewStyle(&excelize.Style{CustomNumFmt: &exp})
	xf.SetColWidth(sheetName, colName, colName, colProgressPoleDateWidth)
	xf.SetColStyle(sheetName, colName, style)
	colName = getColName(colProgressPoleHeight)
	xf.SetColWidth(sheetName, colName, colName, colProgressPoleHeightWidth)
	colName = getColName(colProgressPoleMaterial)
	xf.SetColWidth(sheetName, colName, colName, colProgressPoleMaterialWidth)
	colName = getColName(colProgressPoleProduct)
	xf.SetColWidth(sheetName, colName, colName, colProgressPoleProductWidth)

	row := 1
	// Set Poles infos
	xf.SetCellValue(sheetName, xlsx.RcToAxis(row, colProgressPoleRef), "POI")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(row, colProgressPoleSticker), "Appui")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(row, colProgressPoleCity), "Ville")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(row, colProgressPoleAddress), "Adresse")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(row, colProgressPoleDate), "Date")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(row, colProgressPoleHeight), "Hauteur")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(row, colProgressPoleMaterial), "Materiau")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(row, colProgressPoleProduct), "Prestations")
	row++
	for _, pole := range progressPoles {
		xf.SetCellValue(sheetName, xlsx.RcToAxis(row, colProgressPoleRef), pole.Ref)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(row, colProgressPoleSticker), pole.Sticker)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(row, colProgressPoleCity), pole.City)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(row, colProgressPoleAddress), pole.Address)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(row, colProgressPoleDate), date.DateFrom(pole.Date).ToTime()) // Date
		xf.SetCellValue(sheetName, xlsx.RcToAxis(row, colProgressPoleHeight), pole.Height)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(row, colProgressPoleMaterial), pole.Material)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(row, colProgressPoleProduct), strings.Join(pole.Product, ", "))
		row++
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
	ps.Poles = []*Pole{}

	//
	// Read Poles Header & Info
	//
	if err := checkValue(xf, sheetName, xlsx.RcToAxis(rowPoleHeader, 1), "pole"); err != nil {
		return nil, err
	}
	productKeys := map[int]string{}
	for _, col := range []int{colPoleProduct, colPoleProduct + 1, colPoleProduct + 2, colPoleProduct + 3, colPoleProduct + 4, colPoleProduct + 5, colPoleProduct + 6, colPoleProduct + 7, colPoleProduct + 8, colPoleProduct + 9} {
		productKeys[col], _ = xf.GetCellValue(sheetName, xlsx.RcToAxis(rowPoleHeader, col))
	}
	//if err := checkValue(xf, sheetName, xlsx.RcToAxis(rowPoleInfo, colPoleId), "pole"); err != nil {
	//	return ps, nil
	//}

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

	timeStamp := date.Now().TimeStamp()

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
			state = poleconst.StateDictToDo
		}

		comment := getCol(colPoleComment)
		if geomsg != "" {
			comment += "\n" + geomsg
		}

		products := []string{}
		height := 0
		material := getCol(colPoleMaterial)
		if material != "" && getCol(colPoleHeight) == "" {
			// Decode CAPFT info
			lowerComment := strings.ToLower(comment)
			if !strings.Contains(lowerComment, "redressement") && !strings.Contains(lowerComment, "renforcement") {
				products = append(products, poleconst.ProductCreation)
				if strings.Contains(lowerComment, "remplacement") {
					products = append(products, poleconst.ProductReplace)
				}
			}
			material, height = DecodeCAPFTPoleInfo(material, &products)
		} else {
			if getCol(colPoleHeight) != "" {
				height, err = strconv.Atoi(getCol(colPoleHeight))
				if err != nil {
					return nil, fmt.Errorf("%s: could not parse height '%s': %s", cellCoord(colPoleHeight, line), getCol(colPoleHeight), err.Error())
				}
			}
			for col := colPoleProduct; col < colPoleProduct+len(productKeys); col++ {
				if getCol(col) == "1" {
					products = append(products, productKeys[col])
				}
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

		pole := &Pole{
			Id:             id,
			Ref:            getCol(colPoleRef),
			City:           strings.Title(strings.ToLower(getCol(colPoleCity))),
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
			Material:       material,
			AspiDate:       aspdate,
			Kizeo:          getCol(colPoleKizeo),
			Comment:        comment,
			Product:        products,
			TimeStamp:      timeStamp,
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
