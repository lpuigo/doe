package xlsextract

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"io"
	"sort"
	"strconv"
	"strings"
)

func WriteXlsReport(file string, recs []*PoleRecord) error {
	rw := reportParser{
		file:   excelize.NewFile(),
		sheet:  "Report",
		curRow: 1,
	}
	rw.file.SetSheetName(rw.file.GetSheetName(0), rw.sheet)

	// sort recs
	sort.Slice(recs, func(i, j int) bool {
		if recs[i].SRO != recs[j].SRO {
			return recs[i].SRO < recs[j].SRO
		}
		return recs[i].Ref < recs[j].Ref
	})

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

func ReadXlsReportFromFile(file string) ([]*PoleRecord, error) {
	xlsFile, err := excelize.OpenFile(file)
	if err != nil {
		return nil, err
	}
	return readXlsReport(xlsFile)
}

func ReadXlsReport(file io.Reader) ([]*PoleRecord, error) {
	xlsFile, err := excelize.OpenReader(file)
	if err != nil {
		return nil, err
	}
	return readXlsReport(xlsFile)
}

func readXlsReport(xlsFile *excelize.File) ([]*PoleRecord, error) {
	rp := reportParser{
		file:   xlsFile,
		sheet:  xlsFile.GetSheetName(0),
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
	rp.file.SetColWidth(rp.sheet, "A", "Z", 15)

	rp.file.SetCellStr(rp.sheet, "A"+strconv.Itoa(rp.curRow), "Récup Image")
	rp.file.SetCellStr(rp.sheet, "B"+strconv.Itoa(rp.curRow), "Date")
	rp.file.SetCellStr(rp.sheet, "C"+strconv.Itoa(rp.curRow), "Heure")
	rp.file.SetCellStr(rp.sheet, "D"+strconv.Itoa(rp.curRow), "SRO")
	rp.file.SetCellStr(rp.sheet, "E"+strconv.Itoa(rp.curRow), "Appui")
	rp.file.SetCellStr(rp.sheet, "F"+strconv.Itoa(rp.curRow), "Commentaire")
	rp.file.SetCellStr(rp.sheet, "G"+strconv.Itoa(rp.curRow), "Latitude")
	rp.file.SetCellStr(rp.sheet, "H"+strconv.Itoa(rp.curRow), "Longitude")
	rp.file.SetCellStr(rp.sheet, "I"+strconv.Itoa(rp.curRow), "Google Maps")
	rp.file.SetCellStr(rp.sheet, "J"+strconv.Itoa(rp.curRow), "Image (link)")

	return nil
}

func (rp *reportParser) writeRecord(rec *PoleRecord) error {
	if rec.ExtractPicture {
		rp.file.SetCellStr(rp.sheet, "A"+strconv.Itoa(rp.curRow), "Oui")
	}
	rp.file.SetCellStr(rp.sheet, "B"+strconv.Itoa(rp.curRow), rec.Date)
	rp.file.SetCellStr(rp.sheet, "C"+strconv.Itoa(rp.curRow), rec.Hour)
	rp.file.SetCellStr(rp.sheet, "D"+strconv.Itoa(rp.curRow), rec.SRO)
	rp.file.SetCellStr(rp.sheet, "E"+strconv.Itoa(rp.curRow), rec.Ref)
	rp.file.SetCellStr(rp.sheet, "F"+strconv.Itoa(rp.curRow), rec.Comment)
	rp.file.SetCellFloat(rp.sheet, "G"+strconv.Itoa(rp.curRow), rec.lat, 7, 64)
	rp.file.SetCellFloat(rp.sheet, "H"+strconv.Itoa(rp.curRow), rec.long, 7, 64)

	gmcoord := "I" + strconv.Itoa(rp.curRow)
	rp.file.SetCellStr(rp.sheet, gmcoord, "maps")
	style, _ := rp.file.NewStyle(`{"font":{"color":"#1265BE","underline":"single"}}`)
	rp.file.SetCellStr(rp.sheet, gmcoord, "maps")
	rp.file.SetCellHyperLink(rp.sheet, gmcoord, getGMAPUrl(rec), "External")
	rp.file.SetCellStyle(rp.sheet, gmcoord, gmcoord, style)

	for i, label := range rec.GetImageLabels() {
		coord := fmt.Sprintf("%c%d", 'J'+i, rp.curRow)
		rp.file.SetCellStr(rp.sheet, coord, label)
		rp.file.SetCellHyperLink(rp.sheet, coord, rec.Images[label], "External")
		rp.file.SetCellStyle(rp.sheet, coord, coord, style)
	}
	return nil
}

func (rp *reportParser) readRecord(cols []string) (*PoleRecord, error) {
	var err error
	if len(cols) < 8 {
		return nil, fmt.Errorf("not enough columns")
	}
	rec := &PoleRecord{
		Images: make(map[string]string),
	}
	rec.ExtractPicture = strings.ToLower(cols[0]) == "oui" // Column A ExtractPicture
	rec.Date = cols[1]                                     // column B Date
	rec.Hour = cols[2]                                     // column C Heure
	rec.SRO = cols[3]                                      // column D SRO
	rec.Ref = cols[4]                                      // column E Appui
	rec.Comment = cols[5]                                  // column F Comment
	// column G Latitude
	rec.lat, err = strconv.ParseFloat(strings.ReplaceAll(cols[6], ",", "."), 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse latitude '%s' at %s", cols[5], "F"+strconv.Itoa(rp.curRow))
	}
	// column H Longitude
	rec.long, err = strconv.ParseFloat(strings.ReplaceAll(cols[7], ",", "."), 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse longitude '%s' at %s", cols[6], "G"+strconv.Itoa(rp.curRow))
	}
	// column J> Image Link
	if !rec.ExtractPicture { // no image to process, exit now
		return rec, nil
	}
	for colnum := 9; colnum < len(cols); colnum++ {
		label := cols[colnum]
		if label == "" {
			continue
		}
		coord := fmt.Sprintf("%c%d", 'J'+colnum-9, rp.curRow)
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
