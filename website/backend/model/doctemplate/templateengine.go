package doctemplate

import (
	"archive/zip"
	"fmt"
	"github.com/lpuig/ewin/doe/website/backend/model/actors"
	"github.com/lpuig/ewin/doe/website/backend/model/items"
	"github.com/lpuig/ewin/doe/website/backend/model/polesites"
	"github.com/lpuig/ewin/doe/website/backend/tools/xlsx"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/lpuig/ewin/doe/model"
	"github.com/lpuig/ewin/doe/website/backend/model/clients"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	"github.com/lpuig/ewin/doe/website/backend/model/ripsites"
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
func (te *DocTemplateEngine) GetRipsiteXLSAttachement(w io.Writer, site *ripsites.Site, getClient clients.ClientByName, actorById clients.ActorById) error {
	client := getClient(site.Client)
	if client == nil {
		return fmt.Errorf("unknown client '%s'", site.Client)
	}

	its, err := site.Itemize(client.Bpu)
	if err != nil {
		return fmt.Errorf("unable to create items: %s", err.Error())
	}

	return te.GetItemsXLSAttachement(w, its, actorById)
}

// GetRipsiteXLSAttachementName returns the name of the XLSx file pertaining to given Ripsite
func (te *DocTemplateEngine) GetPolesiteXLSAttachementName(site *polesites.PoleSite) string {
	return fmt.Sprintf("ATTACHEMENT %s.xlsx", site.Ref)
}

// GetRipsiteXLSAttachement generates and writes on given writer the attachment data pertaining to given Ripsite
func (te *DocTemplateEngine) GetPolesiteXLSAttachement(w io.Writer, site *polesites.PoleSite, getClient clients.ClientByName, actorById clients.ActorById) error {
	client := getClient(site.Client)
	if client == nil {
		return fmt.Errorf("unknown client '%s'", site.Client)
	}

	its, err := site.Itemize(client.Bpu)
	if err != nil {
		return fmt.Errorf("unable to create items: %s", err.Error())
	}

	return te.GetItemsXLSAttachement(w, its, actorById)
}

// GetRipsiteXLSAttachement generates and writes on given writer the attachment data pertaining to given items
func (te *DocTemplateEngine) GetItemsXLSAttachement(w io.Writer, its []*items.Item, actorById clients.ActorById) error {
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
	xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 0), "Client")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 1), "Site")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 2), "Activité")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 3), "Item")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 4), "Info")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 5), "Code BPU")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 6), "Quantité")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 7), "Prix")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 8), "Quant. Tr.")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 9), "Travail")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 10), "Installé")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 11), "Equipe")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 12), "Date")
	for _, item := range its {
		if !(item.Todo && item.Quantity > 0) {
			continue
		}
		row++
		actors := []string{}
		for _, actorId := range item.Actors {
			if act := actorById(actorId); act != "" {
				actors = append(actors, act)
			}
		}
		xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 0), item.Client)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 1), item.Site)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 2), item.Activity)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 3), item.Name)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 4), item.Info)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 5), item.Article.Name)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 6), item.Quantity)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 7), item.Price())
		xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 8), item.WorkQuantity)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 9), item.Work())
		if item.Done {
			xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 10), "Oui")
			xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 11), strings.Join(actors, "\n"))
			xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 12), item.Date)
		} else {
			xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 10), "")
			xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 11), "")
			xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 12), "")
		}
	}
	xf.UpdateLinkedValue()
	return xf.Write(w)
}

const (
	awhrXLSFile   string = "CRA _COMPANY_ _DATE_.xlsx"
	awhrSheetName string = "CRA"
	awhrRowStart  int    = 4
)

func (te *DocTemplateEngine) GetActorsWorkingHoursRecordXLS(w io.Writer, monthDate string, actors []*actors.Actor) error {
	file := filepath.Join(te.tmplDir, awhrXLSFile)
	xf, err := excelize.OpenFile(file)
	if err != nil {
		return err
	}

	found := xf.GetSheetIndex(awhrSheetName)
	if found == 0 {
		return fmt.Errorf("could not find CRA sheet")
	}

	begDate := date.DateFrom(monthDate).GetMonth()
	endDate := begDate.AddDays(32).GetMonth().AddDays(-1)
	row := awhrRowStart
	for _, week := range date.GetMonthlyWeeksBetween(begDate, endDate) {
		strBeg := week.Begin.ToDDMMYYYY()
		strEnd := week.End.ToDDMMYYYY()
		strWeek := week.ToDateStringRange()
		for _, actor := range actors {
			if !actor.IsActiveOnDateRange(strWeek) {
				continue
			}
			xf.SetCellValue(awhrSheetName, xlsx.RcToAxis(row, 1), strBeg)
			xf.SetCellValue(awhrSheetName, xlsx.RcToAxis(row, 2), strEnd)
			xf.SetCellValue(awhrSheetName, xlsx.RcToAxis(row, 3), actor.LastName)
			xf.SetCellValue(awhrSheetName, xlsx.RcToAxis(row, 4), actor.FirstName)
			xf.SetCellValue(awhrSheetName, xlsx.RcToAxis(row, 5), actor.Role)
			activity, comment := actor.GetActivityInfoFor(strWeek)
			xf.SetCellValue(awhrSheetName, xlsx.RcToAxis(row, 6), activity*7)
			xf.SetCellValue(awhrSheetName, xlsx.RcToAxis(row, 7), 0)
			xf.SetCellValue(awhrSheetName, xlsx.RcToAxis(row, 8), activity)
			xf.SetCellValue(awhrSheetName, xlsx.RcToAxis(row, 9), comment)

			xf.SetCellValue(awhrSheetName, xlsx.RcToAxis(row, 11), actor.Company)

			row++
		}
	}
	return xf.Write(w)
}
