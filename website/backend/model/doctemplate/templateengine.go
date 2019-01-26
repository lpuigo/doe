package doctemplate

import (
	"fmt"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	"github.com/lpuig/ewin/doe/website/backend/model/worksites"
	"github.com/tealeg/xlsx"
	"io"
	"os"
	"path/filepath"
)

const (
	attachementFile string = "ATTACHEMENT _REF_CITY.xlsx"
	rowTotal        int    = 20
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
	xf, err := xlsx.OpenFile(file)
	if err != nil {
		return err
	}

	xs, found := xf.Sheet["EWIN"]
	if !found {
		return fmt.Errorf("could not find EWIN sheet")
	}

	line := 25
	totEl := 0
	for _, ord := range ws.Orders {
		for _, tr := range ord.Troncons {
			xs.Cell(line, colFI).Value = ord.Ref
			xs.Cell(line, colPA).Value = ws.Ref
			xs.Cell(line, colPB).Value = tr.Pb.Ref
			xs.Cell(line, colCode).Value = "CEM42"
			xs.Cell(line, colAddr).Value = tr.Pb.Address
			xs.Cell(line, colCity).Value = ws.City
			nbEl := tr.NbRacco
			if tr.Blockage {
				nbEl = 0
			}
			totEl += nbEl
			xs.Cell(line, colNbEl).SetInt(nbEl)
			endDate := ""
			if tr.MeasureDate != "" {
				endDate = date.DateFrom(tr.MeasureDate).ToDDMMAAAA()
			}
			xs.Cell(line, colEndD).Value = endDate
			line++
		}
	}

	xs.Cell(rowTotal, colTotal).SetInt(totEl)

	return xf.Write(w)
}
