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

func newXlsParser(file *excelize.File, sheet string, rows *excelize.Rows) *xlsPoleParser {
	return &xlsPoleParser{
		file:     file,
		sheet:    sheet,
		colNames: make(map[int]string),
		rows:     rows,
		rowNum:   0,
	}
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

func (xpp *xlsPoleParser) ParseHeader() error {
	// Parse Header
	if !xpp.Next() {
		return fmt.Errorf("could not read header row: %s\n", xpp.rows.Error())
	}
	headerCols, err := xpp.rows.Columns()
	if err != nil {
		return fmt.Errorf("could not read header columns: %s\n", err.Error())
	}

	for i, headerCol := range headerCols {
		colString, err := excelize.ColumnNumberToName(i + 1)
		if err != nil {
			return fmt.Errorf("could not create column number for column %d\n", i+1)
		}
		//fmt.Printf("%s: '%s'\n", colString, headerCol)
		_ = colString
		if headerCol == "" {
			continue
		}
		xpp.colNames[i+1] = headerCol
	}
	return nil
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
		switch colName {
		case "Date de réponse":
			datefloat, err := strconv.ParseFloat(value, 64)
			if err != nil {
				// Date is strinig format, parse it
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
		case "SRO":
			sro = value
		case "Référence Poteau":
			ref = value
		case "Localisation GPS Poteau":
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
		case "Commentaire":
			comment = value
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
