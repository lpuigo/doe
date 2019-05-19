package doctemplate

import (
	"archive/zip"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	"github.com/lpuig/ewin/doe/website/backend/model/worksites"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	attachementFile string = "ATTACHEMENT _REF_CITY.xlsx"
	sheetName       string = "EWIN"
	rowTotal        int    = 21
	colTotal        int    = 9
	colFI           int    = 1
	colPA           int    = 2
	colPB           int    = 3
	colCode         int    = 4
	colAddr         int    = 5
	colCity         int    = 6
	colNbEl         int    = 7
	colEndD         int    = 8
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

// GetAttachmentName returns the name of the WLSx file pertaining to given Worksite
func (te *DocTemplateEngine) GetAttachmentName(ws *worksites.WorkSiteRecord) string {
	return fmt.Sprintf("ATTACHEMENT %s_%s.xlsx", ws.Ref, ws.City)
}

// GetAttachmentXLS generates and writes on given writer the attachment data pertaining to given Worksite
func (te *DocTemplateEngine) GetAttachmentXLS(w io.Writer, ws *worksites.WorkSiteRecord) error {
	file := filepath.Join(te.tmplDir, attachementFile)
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

	line := 26
	totEl := 0
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
			totEl += nbEl
			xf.SetCellInt(sheetName, coord(line, colNbEl), nbEl)
			endDate := ""
			if tr.MeasureDate != "" {
				endDate = date.DateFrom(tr.MeasureDate).ToDDMMAAAA()
			}
			xf.SetCellStr(sheetName, coord(line, colEndD), endDate)
			line++
		}
	}

	xf.SetCellInt(sheetName, coord(rowTotal, colTotal), totEl)
	xf.UpdateLinkedValue()

	return xf.Write(w)
}

// GetDOEArchiveName returns the name of the DOE Struct zip file pertaining to given Worksite
func (te *DocTemplateEngine) GetDOEArchiveName(ws *worksites.WorkSiteRecord) string {
	return fmt.Sprintf("DOE %s_%s.zip", ws.Ref, ws.City)
}

// GetDOEArchiveZIP generates and writes on given writer the DOE Struct zip pertaining to given Worksite
func (te *DocTemplateEngine) GetDOEArchiveZIP(w io.Writer, ws *worksites.WorkSiteRecord) error {
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
