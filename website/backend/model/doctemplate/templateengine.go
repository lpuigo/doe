package doctemplate

import (
	"archive/zip"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/lpuig/ewin/chantiersalsace/parsesuivi/xls"
	"github.com/lpuig/ewin/doe/model"
	"github.com/lpuig/ewin/doe/website/backend/model/clients"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	"github.com/lpuig/ewin/doe/website/backend/model/ripsites"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	xlsWorksiteAttachementFile string = "ATTACHEMENT _REF_CITY.xlsx"
	xlsRipsiteAttachementFile  string = "ATTACHEMENT _RIPREF_.xlsx"
	sheetName                  string = "EWIN"

	rowCEM12 int = 13
	rowCEM22 int = 14
	rowCEM32 int = 15
	rowCEM42 int = 16
	rowIMB   int = 21

	rowRipHeader int = 8

	colFI       int = 1
	colPA       int = 2
	colPB       int = 3
	colCode     int = 4
	colAddr     int = 5
	colCity     int = 6
	colNbEl     int = 7
	colEndD     int = 8
	colCEMPrice int = 8
)

type DocTemplateEngine struct {
	tmplDir string
}

func NewDocTemplateEngine(dir string) (*DocTemplateEngine, error) {
	f, err := os.Stat(dir)
	if err != nil {
		return nil, fmt.Errorf("could not find template directory")
	}
	if !f.IsDir() {
		return nil, fmt.Errorf("given path is not a directory")
	}

	return &DocTemplateEngine{tmplDir: dir}, nil
}

// GetWorksiteXLSAttachmentName returns the name of the WLSx file pertaining to given Worksite
func (te *DocTemplateEngine) GetWorksiteXLSAttachmentName(ws *model.Worksite) string {
	return fmt.Sprintf("ATTACHEMENT %s_%s.xlsx", ws.Ref, ws.City)
}

// GetWorksiteXLSAttachment generates and writes on given writer the attachment data pertaining to given Worksite
func (te *DocTemplateEngine) GetWorksiteXLSAttachment(w io.Writer, ws *model.Worksite, getClient func(clientName string) *clients.Client) error {
	client := getClient(ws.Client)
	if client == nil {
		return fmt.Errorf("unknown client '%s'", ws.Client)
	}

	file := filepath.Join(te.tmplDir, xlsWorksiteAttachementFile)
	xf, err := excelize.OpenFile(file)
	if err != nil {
		return err
	}

	coord := func(row, col int) string {
		acol := excelize.ToAlphaString(col)
		return acol + strconv.Itoa(row)
	}

	found := xf.GetSheetIndex(sheetName)
	if found == 0 {
		return fmt.Errorf("could not find EWIN sheet")
	}

	// Set BPU prices
	cemLine := map[string]int{
		"CEM12": rowCEM12,
		"CEM22": rowCEM22,
		"CEM32": rowCEM32,
		"CEM42": rowCEM42,
	}
	for _, article := range client.GetOrangeArticles() {
		xf.SetCellValue(sheetName, coord(cemLine[article.Name], colCEMPrice), article.Price)
	}

	line := rowIMB
	for _, ord := range ws.Orders {
		for _, tr := range ord.Troncons {
			xf.SetCellStr(sheetName, coord(line, colFI), ord.Ref)
			xf.SetCellStr(sheetName, coord(line, colPA), ws.Ref)
			xf.SetCellStr(sheetName, coord(line, colPB), tr.Pb.Ref)
			xf.SetCellStr(sheetName, coord(line, colCode), tr.Article)
			xf.SetCellStr(sheetName, coord(line, colAddr), tr.Pb.Address)
			xf.SetCellStr(sheetName, coord(line, colCity), ws.City)
			nbEl := tr.NbRacco
			if tr.Blockage {
				nbEl = 0
			}
			xf.SetCellInt(sheetName, coord(line, colNbEl), nbEl)
			endDate := ""
			if tr.MeasureDate != "" {
				endDate = date.DateFrom(tr.MeasureDate).ToDDMMYYYY()
			}
			xf.SetCellStr(sheetName, coord(line, colEndD), endDate)
			line++
		}
	}

	xf.UpdateLinkedValue()

	return xf.Write(w)
}

// GetDOEArchiveName returns the name of the DOE Struct zip file pertaining to given Worksite
func (te *DocTemplateEngine) GetDOEArchiveName(ws *model.Worksite) string {
	return fmt.Sprintf("DOE %s_%s.zip", ws.Ref, ws.City)
}

// GetDOEArchiveZIP generates and writes on given writer the DOE Struct zip pertaining to given Worksite
func (te *DocTemplateEngine) GetDOEArchiveZIP(w io.Writer, ws *model.Worksite) error {
	zw := zip.NewWriter(w)

	path := strings.TrimSuffix(te.GetDOEArchiveName(ws), ".zip")

	makeDir := func(base ...string) string {
		return filepath.Join(base...) + "/"
	}

	// Create list of Pb's PT refs
	Pbs := []string{}
	for _, ord := range ws.Orders {
		for _, tr := range ord.Troncons {
			trname := tr.Pb.RefPt
			if tr.Blockage {
				trname += " (bloqué)"
			}
			Pbs = append(Pbs, trname)
		}
	}

	dirs := []string{
		"Alveoles utilisees",
		"Changement du synoptique",
		"Convention ajoutee",
		"Photos poteaux",
		"Photos " + ws.Pa.Ref,
	}
	for _, dir := range dirs {
		_, err := zw.Create(makeDir(path, dir))
		if err != nil {
			return err
		}
	}

	mpath := makeDir(path, "Mesures")
	for _, ptname := range Pbs {
		_, err := zw.Create(makeDir(mpath, ptname))
		if err != nil {
			return err
		}
	}

	ppath := makeDir(path, "Photos PBs")
	for _, ptname := range Pbs {
		for _, subd := range []string{"Pose Boitier", "Test laser"} {
			_, err := zw.Create(makeDir(ppath, ptname, subd))
			if err != nil {
				return err
			}
		}
	}

	return zw.Close()
}

// GetRipsiteXLSAttachementName returns the name of the XLSx file pertaining to given Ripsite
func (te *DocTemplateEngine) GetRipsiteXLSAttachementName(site *ripsites.Site) string {
	return fmt.Sprintf("ATTACHEMENT %s.xlsx", site.Ref)
}

// GetRipsiteXLSAttachement generates and writes on given writer the attachment data pertaining to given Ripsite
func (te *DocTemplateEngine) GetRipsiteXLSAttachement(w io.Writer, site *ripsites.Site, getClient func(clientName string) *clients.Client) error {
	client := getClient(site.Client)
	if client == nil {
		return fmt.Errorf("unknown client '%s'", site.Client)
	}

	items, err := site.Itemize(client.Bpu)
	if err != nil {
		return fmt.Errorf("unable to create items: %s", err.Error())
	}

	file := filepath.Join(te.tmplDir, xlsRipsiteAttachementFile)
	xf, err := excelize.OpenFile(file)
	if err != nil {
		return err
	}

	found := xf.GetSheetIndex(sheetName)
	if found == 0 {
		return fmt.Errorf("could not find EWIN sheet")
	}

	row := rowRipHeader
	xf.SetCellValue(sheetName, xls.RcToAxis(row, 0), "Item")
	xf.SetCellValue(sheetName, xls.RcToAxis(row, 1), "Info")
	xf.SetCellValue(sheetName, xls.RcToAxis(row, 2), "Code BPU")
	xf.SetCellValue(sheetName, xls.RcToAxis(row, 3), "Quantité")
	xf.SetCellValue(sheetName, xls.RcToAxis(row, 4), "Prix")
	xf.SetCellValue(sheetName, xls.RcToAxis(row, 5), "Quant. Tr.")
	xf.SetCellValue(sheetName, xls.RcToAxis(row, 6), "Travail")
	xf.SetCellValue(sheetName, xls.RcToAxis(row, 7), "Installé")
	xf.SetCellValue(sheetName, xls.RcToAxis(row, 8), "Equipe")
	xf.SetCellValue(sheetName, xls.RcToAxis(row, 9), "Date")
	for _, item := range items {
		if !(item.Todo && item.Quantity > 0) {
			continue
		}
		row++
		xf.SetCellValue(sheetName, xls.RcToAxis(row, 0), item.Name)
		xf.SetCellValue(sheetName, xls.RcToAxis(row, 1), item.Info)
		xf.SetCellValue(sheetName, xls.RcToAxis(row, 2), item.Article.Name)
		xf.SetCellValue(sheetName, xls.RcToAxis(row, 3), item.Quantity)
		xf.SetCellValue(sheetName, xls.RcToAxis(row, 4), item.Price())
		xf.SetCellValue(sheetName, xls.RcToAxis(row, 5), item.WorkQuantity)
		xf.SetCellValue(sheetName, xls.RcToAxis(row, 6), item.Work())
		if item.Done {
			xf.SetCellValue(sheetName, xls.RcToAxis(row, 7), "Oui")
			xf.SetCellValue(sheetName, xls.RcToAxis(row, 8), item.Team)
			xf.SetCellValue(sheetName, xls.RcToAxis(row, 9), item.Date)
		} else {
			xf.SetCellValue(sheetName, xls.RcToAxis(row, 7), "")
			xf.SetCellValue(sheetName, xls.RcToAxis(row, 8), "")
			xf.SetCellValue(sheetName, xls.RcToAxis(row, 9), "")
		}
	}
	xf.UpdateLinkedValue()
	return xf.Write(w)
}
