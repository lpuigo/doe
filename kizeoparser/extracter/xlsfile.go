package extracter

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/lpuig/ewin/doe/kizeoparser/api"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type xlscontext struct {
	file   *excelize.File
	sheet  string
	curRow int
}

func (xc *xlscontext) writeHeader() {
	xc.file.SetColWidth(xc.sheet, "A", "Z", 15)

	xc.file.SetCellStr(xc.sheet, "A"+strconv.Itoa(xc.curRow), "Récup Image")
	xc.file.SetCellStr(xc.sheet, "B"+strconv.Itoa(xc.curRow), "Date")
	xc.file.SetCellStr(xc.sheet, "C"+strconv.Itoa(xc.curRow), "Heure")
	xc.file.SetCellStr(xc.sheet, "D"+strconv.Itoa(xc.curRow), "SRO")
	xc.file.SetCellStr(xc.sheet, "E"+strconv.Itoa(xc.curRow), "Appui")
	xc.file.SetCellStr(xc.sheet, "F"+strconv.Itoa(xc.curRow), "Commentaire")
	xc.file.SetCellStr(xc.sheet, "G"+strconv.Itoa(xc.curRow), "Latitude")
	xc.file.SetCellStr(xc.sheet, "H"+strconv.Itoa(xc.curRow), "Longitude")
	xc.file.SetCellStr(xc.sheet, "I"+strconv.Itoa(xc.curRow), "Google Maps")
	xc.file.SetCellStr(xc.sheet, "J"+strconv.Itoa(xc.curRow), "FormId")
	xc.file.SetCellStr(xc.sheet, "K"+strconv.Itoa(xc.curRow), "recordId")
	xc.file.SetCellStr(xc.sheet, "L"+strconv.Itoa(xc.curRow), "Image (uuid)")
}

const (
	linkPrefix string = "http://localhost/"
)

func (xc *xlscontext) writeRecord(rec *api.SearchData) {
	updatedate, updatehour := rec.GetDateHour()
	sro, ref := rec.GetSroRef()
	curRowNum := strconv.Itoa(xc.curRow)

	if rec.ExtractData {
		xc.file.SetCellStr(xc.sheet, "A"+curRowNum, "api")
	}
	xc.file.SetCellStr(xc.sheet, "B"+curRowNum, updatedate)
	xc.file.SetCellStr(xc.sheet, "C"+curRowNum, updatehour)
	xc.file.SetCellStr(xc.sheet, "D"+curRowNum, sro)
	xc.file.SetCellStr(xc.sheet, "E"+curRowNum, ref)
	xc.file.SetCellStr(xc.sheet, "F"+curRowNum, rec.Comment)
	xc.file.SetCellStr(xc.sheet, "G"+curRowNum, rec.Geoloc.Lat)
	xc.file.SetCellStr(xc.sheet, "H"+curRowNum, rec.Geoloc.Long)

	style, _ := xc.file.NewStyle(`{"font":{"color":"#1265BE","underline":"single"}}`)
	gmcoord := "I" + curRowNum
	xc.file.SetCellStr(xc.sheet, gmcoord, rec.Geoloc.Lat+", "+rec.Geoloc.Long)
	xc.file.SetCellHyperLink(xc.sheet, gmcoord, getGMAPUrl(rec), "External")
	xc.file.SetCellStyle(xc.sheet, gmcoord, gmcoord, style)

	xc.file.SetCellStr(xc.sheet, "J"+curRowNum, rec.FormID)
	xc.file.SetCellStr(xc.sheet, "K"+curRowNum, rec.ID)

	imgNames := []string{}
	for imgName, _ := range rec.Pictures {
		imgNames = append(imgNames, imgName)
	}
	sort.Strings(imgNames)
	for numImg, imgName := range imgNames {
		coord, _ := excelize.CoordinatesToCellName(12+numImg, xc.curRow) // col L and beyond
		xc.file.SetCellStr(xc.sheet, coord, imgName)
		xc.file.SetCellHyperLink(xc.sheet, coord, linkPrefix+rec.Pictures[imgName], "External")
		xc.file.SetCellStyle(xc.sheet, coord, coord, style)
	}
}

func (xc *xlscontext) readRecord(cols []string) (*api.SearchData, error) {
	if len(cols) < 10 {
		return nil, fmt.Errorf("not enough columns")
	}
	rec := &api.SearchData{
		Pictures: make(map[string]string),
	}
	rec.ExtractData = strings.ToLower(cols[0]) == "api" // Column A ExtractPicture
	rec.UpdateTime = cols[1] + " " + cols[2]            // column B Date & C Hour
	rec.SummarySubtitle = cols[3] + "|" + cols[4]       // column D SRO & E Appui
	rec.Comment = cols[5]                               // column F Comment
	rec.Geoloc.Lat = cols[6]                            // column G Latitude
	rec.Geoloc.Long = cols[7]                           // column H Longitude
	rec.FormID = cols[9]                                // column J Form ID
	rec.ID = cols[10]                                   // column K Record ID

	for colnum := 11; colnum < len(cols); colnum++ {
		label := cols[colnum]
		if label == "" {
			continue
		}
		coord, _ := excelize.CoordinatesToCellName(colnum+1, xc.curRow) // col L and beyond
		isLink, link, err := xc.file.GetCellHyperLink(xc.sheet, coord)
		if !isLink || err != nil {
			return nil, fmt.Errorf("could not get link for image '%s' at %s", label, coord)
		}
		rec.Pictures[label] = strings.TrimPrefix(link, linkPrefix)
	}

	return rec, nil
}

func writeXlsFormsFile(file, formName string, form []*api.SearchData) error {
	xc := xlscontext{
		file:   excelize.NewFile(),
		sheet:  formName,
		curRow: 1,
	}
	xc.file.SetSheetName(xc.file.GetSheetName(0), xc.sheet)

	// sort recs
	SortSearchDatasBySroRef(form)

	// Write Header
	xc.writeHeader()
	// Write Record
	for i, rec := range form {
		xc.curRow = i + 2
		xc.writeRecord(rec)
	}

	err := xc.file.SaveAs(file)
	if err != nil {
		return fmt.Errorf("could not save xls report as '%s': %s\n", filepath.Base(file), err.Error())
	}
	return nil
}

func readXlsFormsFile(file string, forms map[int][]*api.SearchData) error {
	xlsFile, err := excelize.OpenFile(file)
	if err != nil {
		return err
	}
	xc := xlscontext{
		file:   xlsFile,
		sheet:  xlsFile.GetSheetName(0),
		curRow: 1,
	}
	rows, err := xc.file.Rows(xc.sheet)
	if err != nil {
		return err
	}

	// consume header
	if !rows.Next() {
		return nil
	}
	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	// consume records
	for rows.Next() {
		xc.curRow++
		cols, err = rows.Columns()
		if err != nil {
			return err
		}
		rec, err := xc.readRecord(cols)
		if err != nil {
			return err
		}
		fId, err := strconv.Atoi(rec.FormID)
		if err != nil {
			return fmt.Errorf("could not read '%s' as a formId (int)", rec.FormID)
		}
		forms[fId] = append(forms[fId], rec)
	}
	return nil
}
