package xlsextract

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"strconv"
	"strings"
	"time"
)

type xlsPoleParser struct {
	file     *excelize.File
	sheet    string
	colNames map[int]string
	rows     *excelize.Rows
	rowNum   int
}

func newXlsParser(file *excelize.File, sheet string) (*xlsPoleParser, error) {
	xpp := &xlsPoleParser{
		file:     file,
		sheet:    sheet,
		colNames: make(map[int]string),
		rows:     nil,
		rowNum:   0,
	}
	err := xpp.UnmergeHeaderCells()
	if err != nil {
		return nil, err
	}
	rows, err := file.Rows(sheet)
	if err != nil {
		return nil, err
	}
	xpp.rows = rows

	return xpp, nil
}

func (xpp *xlsPoleParser) Next() bool {
	xpp.rowNum++
	return xpp.rows.Next()
}

func (xpp *xlsPoleParser) Error() error {
	return xpp.rows.Error()
}

func (xpp *xlsPoleParser) GetRowNum() int {
	return xpp.rowNum
}

func (xpp *xlsPoleParser) UnmergeHeaderCells() error {
	mergedCells, err := xpp.file.GetMergeCells(xpp.sheet)
	if err != nil {
		return err
	}
	for _, mergedCell := range mergedCells {
		value := mergedCell.GetCellValue()
		startAxis, endAxis := mergedCell.GetStartAxis(), mergedCell.GetEndAxis()
		startColumn, startRow, err := excelize.CellNameToCoordinates(startAxis)
		if err != nil {
			return fmt.Errorf("could not get coordinates for startAxis '%s': %s\n", endAxis, err.Error())
		}
		// check if mergedCells is in header
		if startRow != 1 {
			//fmt.Printf("Skiping Merged Cell not on Header : %s %s -> %s\n", value, startAxis, endAxis)
			continue
		}

		endColumn, _, err := excelize.CellNameToCoordinates(endAxis)
		if err != nil {
			return fmt.Errorf("could not get coordinates for endAxis '%s': %s\n", endAxis, err.Error())
		}
		//fmt.Printf("Merged Cell: %s %s -> %s\n", value, startAxis, endAxis)
		err = xpp.file.UnmergeCell(xpp.sheet, startAxis, endAxis)
		if err != nil {
			return fmt.Errorf("could not unmerge cells '%s' to '%s': %s\n", startAxis, endAxis, err.Error())
		}
		for col := startColumn; col <= endColumn; col++ {
			axis, err := excelize.CoordinatesToCellName(col, startRow)
			if err != nil {
				return fmt.Errorf("could not get axis for coordinates col=%d row=%d: %s\n", col, startRow, err.Error())
			}
			err = xpp.file.SetCellDefault(xpp.sheet, axis, value)
			if err != nil {
				return fmt.Errorf("could not set value for cell '%s:%s': %s\n", xpp.sheet, axis, err.Error())
			}
			//fmt.Printf("\t %s = %s\n", axis, value)
		}
	}
	return nil
}

func (xpp *xlsPoleParser) ParseHeader() error {
	// Parse Header
	if !xpp.Next() {
		return fmt.Errorf("could not read header row: %s\n", xpp.rows.Error())
	}
	headerCols, err := xpp.rows.Columns()
	if err != nil {
		return fmt.Errorf("could not read header columns: %s\n", err.Error())
	}

	previousValue := ""
	value := ""
	chapterNumber := -1
	numberInChapter := 0
	images := false
	for i, headerCol := range headerCols {
		value = strings.ReplaceAll(headerCol, ":", "")
		value = strings.ReplaceAll(value, "/", "_")
		value = strings.ReplaceAll(value, "  ", " ")
		value = strings.Trim(value, " ")
		colString, err := excelize.ColumnNumberToName(i + 1)
		if err != nil {
			return fmt.Errorf("could not create column number for column %d\n", i+1)
		}
		//fmt.Printf("%s: '%s'\n", colString, headerCol)
		_ = colString
		if headerCol == "" {
			continue
		}
		if images {
			if headerCol != previousValue {
				chapterNumber++
				previousValue = headerCol
				numberInChapter = 1
			} else {
				numberInChapter++
			}
			value = fmt.Sprintf("%s %s %d", string(rune('A'+chapterNumber)), value, numberInChapter)
		}
		xpp.colNames[i+1] = value
		if !images && strings.HasPrefix(value, "Localisation") {
			images = true
		}
	}
	return nil
}

func (xpp *xlsPoleParser) PrintColumnNames() {
	fmt.Printf("Columns Names :\n")
	for colNum := 1; colNum <= len(xpp.colNames); colNum++ {
		colName, _ := excelize.ColumnNumberToName(colNum)
		fmt.Printf("\t%s : '%s'\n", colName, xpp.colNames[colNum])
	}
}

func (xpp *xlsPoleParser) getPicture(coord string) (link bool, target string, err error) {
	link, target, err = xpp.file.GetCellHyperLink(xpp.sheet, coord)
	//picName, picData, err = xpp.file.GetPicture(xpp.sheet, coord)
	return
}

func (xpp *xlsPoleParser) ParseRecord() (*PoleRecord, error) {
	cols, err := xpp.rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve columns for row %d: %s\n", xpp.rowNum, err.Error())
	}
	var date, hour, sro, ref, comment string
	var long, lat float64
	img := map[string]string{}

	for i, value := range cols {
		coord, err := excelize.CoordinatesToCellName(i+1, xpp.rowNum)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve columns for column %d, row %d: %s\n", i+1, xpp.rowNum, err.Error())
		}
		colName := xpp.colNames[i+1]
		switch {
		case colName == "Date de réponse":
			datefloat, err := strconv.ParseFloat(value, 64)
			if err != nil {
				// Date is string format, parse it
				times := strings.Split(value, " ")
				if len(times) != 2 {
					fmt.Printf("could not get date time info at %s: '%s'\n", coord, value)
					continue
				}
				hour = times[1]
				dates, err := time.Parse("1/02/06", times[0])
				if err != nil {
					return nil, fmt.Errorf("could not parse date at %s: %s, %s\n", coord, times[0], err.Error())
				}
				date = dates.Format("2006-01-02")
				continue
			}
			datetime, err := excelize.ExcelDateToTime(datefloat, false)
			if err != nil {
				fmt.Printf("could not get date time info at %s: '%f'\n", coord, datefloat)
				continue
			}
			date = datetime.Format("2006-01-02")
			hour = datetime.Format("15:04")
		case colName == "SRO":
			sro = value
		case colName == "Référence Poteau":
			ref = value
		case colName == "Localisation GPS Poteau":
			gps := strings.Split(value, "\n")
			if len(gps) != 2 {
				fmt.Printf("could not get gps info at %s\n", coord)
				continue
			}
			slat := strings.TrimPrefix(gps[0], "Latitude : ")
			lat, err = strconv.ParseFloat(slat, 64)
			if err != nil {
				return nil, fmt.Errorf("could not parle GPS Latitude at %s '%s': %s\n", coord, slat, err.Error())
			}
			slong := strings.TrimPrefix(gps[1], "Longitude : ")
			long, err = strconv.ParseFloat(slong, 64)
			if err != nil {
				return nil, fmt.Errorf("could not parle GPS Longitude at %s '%s': %s\n", coord, slong, err.Error())
			}
		case strings.Contains(colName, "Commentaire"):
			comment += strings.ReplaceAll(value, "\n", "\r\n")
		default:
			link, target, err := xpp.getPicture(coord)
			if err != nil {
				fmt.Printf("could not retrieve picture at %s: %s\n", coord, err.Error())
				continue
			}
			if !link {
				continue
			}
			img[colName] = target
		}
	}

	return &PoleRecord{
		Date:    date,
		Hour:    hour,
		SRO:     sro,
		Ref:     ref,
		Comment: comment,
		Images:  img,
		long:    long,
		lat:     lat,
	}, nil
}
