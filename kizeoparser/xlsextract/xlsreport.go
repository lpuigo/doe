package xlsextract

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"strconv"
	"strings"
)

func WriteXlsReport(file string, recs []*PoleRecord) error {
	rw := reportParser{
		file:   excelize.NewFile(),
		sheet:  "Report",
		curRow: 1,
	}
	rw.file.SetSheetName(rw.file.GetSheetName(1), rw.sheet)

	// Write Header
	err := rw.writeHeader()
	if err != nil {
		return fmt.Errorf("could not write header: %s\n", err.Error())
	}
	// Write Record
	for i, rec := range recs {
		rw.curRow = i + 2
		err = rw.writeRecord(rec)
		if err != nil {
			return fmt.Errorf("could not write record: %s\n", err.Error())
		}
	}

	err = rw.file.SaveAs(file)
	if err != nil {
		return fmt.Errorf("could not save xls report as '%s': %s\n", file, err.Error())
	}
	return nil
}

func ReadXlsReport(file string) ([]*PoleRecord, error) {
	xlsFile, err := excelize.OpenFile(file)
	if err != nil {
		return nil, err
	}
	rp := reportParser{
		file:   xlsFile,
		sheet:  xlsFile.GetSheetName(1),
		curRow: 1,
	}

	rows, err := rp.file.Rows(rp.sheet)
	if err != nil {
		return nil, err
	}

	// consume header
	if !rows.Next() {
		return nil, nil
	}
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	res := []*PoleRecord{}
	// consume records
	for rows.Next() {
		rp.curRow++
		cols, err = rows.Columns()
		if err != nil {
			return nil, err
		}
		rec, err := rp.readRecord(cols)
		if err != nil {
			return nil, err
		}
		res = append(res, rec)
	}
	return res, nil
}

type reportParser struct {
	file   *excelize.File
	sheet  string
	curRow int
}

func (rp *reportParser) writeHeader() error {
	rp.file.SetColWidth(rp.sheet, "A", "I", 15)

	rp.file.SetCellStr(rp.sheet, "A"+strconv.Itoa(rp.curRow), "Date")
	rp.file.SetCellStr(rp.sheet, "B"+strconv.Itoa(rp.curRow), "Heure")
	rp.file.SetCellStr(rp.sheet, "C"+strconv.Itoa(rp.curRow), "SRO")
	rp.file.SetCellStr(rp.sheet, "D"+strconv.Itoa(rp.curRow), "Appui")
	rp.file.SetCellStr(rp.sheet, "E"+strconv.Itoa(rp.curRow), "Latitude")
	rp.file.SetCellStr(rp.sheet, "F"+strconv.Itoa(rp.curRow), "Longitude")
	rp.file.SetCellStr(rp.sheet, "G"+strconv.Itoa(rp.curRow), "Google Maps")
	rp.file.SetCellStr(rp.sheet, "H"+strconv.Itoa(rp.curRow), "Image (link)")

	return nil
}

func (rp *reportParser) writeRecord(rec *PoleRecord) error {
	rp.file.SetCellStr(rp.sheet, "A"+strconv.Itoa(rp.curRow), rec.Date)
	rp.file.SetCellStr(rp.sheet, "B"+strconv.Itoa(rp.curRow), rec.Hour)
	rp.file.SetCellStr(rp.sheet, "C"+strconv.Itoa(rp.curRow), rec.SRO)
	rp.file.SetCellStr(rp.sheet, "D"+strconv.Itoa(rp.curRow), rec.Ref)
	rp.file.SetCellFloat(rp.sheet, "E"+strconv.Itoa(rp.curRow), rec.lat, 7, 64)
	rp.file.SetCellFloat(rp.sheet, "F"+strconv.Itoa(rp.curRow), rec.long, 7, 64)
	rp.file.SetCellStr(rp.sheet, "G"+strconv.Itoa(rp.curRow), "maps")

	style, _ := rp.file.NewStyle(`{"font":{"color":"#1265BE","underline":"single"}}`)
	gmcoord := "G" + strconv.Itoa(rp.curRow)
	rp.file.SetCellStr(rp.sheet, gmcoord, "maps")
	rp.file.SetCellHyperLink(rp.sheet, gmcoord, getGMAPUrl(rec), "External")
	rp.file.SetCellStyle(rp.sheet, gmcoord, gmcoord, style)

	for i, label := range rec.GetImageLabels() {
		coord := fmt.Sprintf("%c%d", 'H'+i, rp.curRow)
		rp.file.SetCellStr(rp.sheet, coord, label)
		rp.file.SetCellHyperLink(rp.sheet, coord, rec.Images[label], "External")
		rp.file.SetCellStyle(rp.sheet, coord, coord, style)
	}
	return nil
}

func (rp *reportParser) readRecord(cols []string) (*PoleRecord, error) {
	var err error
	if len(cols) < 7 {
		return nil, fmt.Errorf("not enough columns")
	}
	rec := &PoleRecord{
		Images: make(map[string]string),
	}
	rec.Date = cols[0] // column A Date
	rec.Hour = cols[1] // column B Heure
	rec.SRO = cols[2]  // column C SRO
	rec.Ref = cols[3]  // column D Appui
	// column E Latitude
	rec.lat, err = strconv.ParseFloat(strings.ReplaceAll(cols[4], ",", "."), 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse latitude '%s' at %s", cols[4], "E"+strconv.Itoa(rp.curRow))
	}
	// column F Longitude
	rec.long, err = strconv.ParseFloat(strings.ReplaceAll(cols[5], ",", "."), 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse longitude '%s' at %s", cols[5], "F"+strconv.Itoa(rp.curRow))
	}
	// column H> Image Link
	for colnum := 7; colnum < len(cols); colnum++ {
		label := cols[colnum]
		coord := fmt.Sprintf("%c%d", 'H'+colnum-7, rp.curRow)
		isLink, link, err := rp.file.GetCellHyperLink(rp.sheet, coord)
		if !isLink || err != nil {
			return nil, fmt.Errorf("could not get link for image '%s' at %s", label, coord)
		}
		rec.Images[label] = link
	}
	return rec, nil
}

func getGMAPUrl(rec *PoleRecord) string {
	return fmt.Sprintf("http://maps.google.com/maps?q=%+.7f,%%20%+.7f", rec.lat, rec.long)
}
