package doctemplate

import (
	"archive/zip"
	"fmt"
	"github.com/lpuig/ewin/doe/website/backend/model/actors"
	"github.com/lpuig/ewin/doe/website/backend/model/items"
	"github.com/lpuig/ewin/doe/website/backend/model/timesheets"
	"github.com/lpuig/ewin/doe/website/backend/tools/xlsx"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/lpuig/ewin/doe/model"
	"github.com/lpuig/ewin/doe/website/backend/model/clients"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
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

	rowRipHeader int = 9

	colFI       int = 2
	colPA       int = 3
	colPB       int = 4
	colCode     int = 5
	colAddr     int = 6
	colCity     int = 7
	colNbEl     int = 8
	colEndD     int = 9
	colCEMPrice int = 9
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
		xf.SetCellValue(sheetName, xlsx.RcToAxis(cemLine[article.Name], colCEMPrice), article.Price)
	}

	line := rowIMB
	for _, ord := range ws.Orders {
		for _, tr := range ord.Troncons {
			xf.SetCellStr(sheetName, xlsx.RcToAxis(line, colFI), ord.Ref)
			xf.SetCellStr(sheetName, xlsx.RcToAxis(line, colPA), ws.Ref)
			xf.SetCellStr(sheetName, xlsx.RcToAxis(line, colPB), tr.Pb.Ref)
			xf.SetCellStr(sheetName, xlsx.RcToAxis(line, colCode), tr.Article)
			xf.SetCellStr(sheetName, xlsx.RcToAxis(line, colAddr), tr.Pb.Address)
			xf.SetCellStr(sheetName, xlsx.RcToAxis(line, colCity), ws.City)
			nbEl := tr.NbRacco
			if tr.Blockage {
				nbEl = 0
			}
			xf.SetCellInt(sheetName, xlsx.RcToAxis(line, colNbEl), nbEl)
			endDate := ""
			if tr.MeasureDate != "" {
				endDate = date.DateFrom(tr.MeasureDate).ToDDMMYYYY()
			}
			xf.SetCellStr(sheetName, xlsx.RcToAxis(line, colEndD), endDate)
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

// GetSiteXLSAttachement generates and writes on given writer the attachment data pertaining to given FoaSite
func (te *DocTemplateEngine) GetItemizableSiteXLSAttachement(w io.Writer, site items.ItemizableSite, getClient clients.ClientByName, actorById clients.ActorById) error {
	client := getClient(site.GetClient())
	if client == nil {
		return fmt.Errorf("unknown client '%s'", site.GetClient())
	}

	its, err := site.Itemize(client.Bpu)
	if err != nil {
		return fmt.Errorf("unable to create items: %s", err.Error())
	}

	return te.GetItemsXLSAttachement(w, its, actorById)
}

// GetItemsXLSAttachement generates and writes on given writer the attachment data pertaining to given items
func (te *DocTemplateEngine) GetItemsXLSAttachement(w io.Writer, its []*items.Item, actorById clients.ActorById) error {
	file := filepath.Join(te.tmplDir, xlsRipsiteAttachementFile)
	xf, err := excelize.OpenFile(file)
	if err != nil {
		return err
	}

	found := xf.GetSheetIndex(sheetName)
	if found == -1 {
		return fmt.Errorf("could not find EWIN sheet")
	}

	row := rowRipHeader
	xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 1), "Client")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 2), "Site")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 3), "Activité")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 4), "Item")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 5), "Info")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 6), "Code BPU")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 7), "Quantité")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 8), "Prix")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 9), "Quant. Tr.")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 10), "Travail")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 11), "Installé")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 12), "Equipe")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 13), "Date")
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
		sort.Strings(actors)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 1), item.Client)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 2), item.Site)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 3), item.Activity)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 4), item.Name)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 5), item.Info)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 6), item.Article.Name)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 7), item.Quantity)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 8), item.Price())
		xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 9), item.WorkQuantity)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 10), item.Work())
		switch {
		case item.Billed:
			xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 11), "Attachement")
			xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 12), strings.Join(actors, "\n"))
			xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 13), date.DateFrom(item.Date).ToTime())
			if attachDate, err := date.ParseDate(item.AttachDate); err == nil {
				xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 14), attachDate.ToTime())
			}
		case item.Done:
			xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 11), "Fait")
			xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 12), strings.Join(actors, "\n"))
			xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 13), date.DateFrom(item.Date).ToTime())
			xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 14), "")
		default:
			xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 11), "")
			xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 12), "")
			xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 13), "")
			xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 14), "")
		}
		xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 15), item.Comment)

	}
	xf.UpdateLinkedValue()
	return xf.Write(w)
}

const (
	awhrXLSFile   string = "CRA _COMPANY_ _DATE_.xlsx"
	awhrSheetName string = "CRA"
	awhrRowStart  int    = 5
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
			xf.SetCellValue(awhrSheetName, xlsx.RcToAxis(row, 2), strBeg)
			xf.SetCellValue(awhrSheetName, xlsx.RcToAxis(row, 3), strEnd)
			xf.SetCellValue(awhrSheetName, xlsx.RcToAxis(row, 4), actor.LastName)
			xf.SetCellValue(awhrSheetName, xlsx.RcToAxis(row, 5), actor.FirstName)
			xf.SetCellValue(awhrSheetName, xlsx.RcToAxis(row, 6), actor.Role)
			activity, comment := actor.GetActivityInfoFor(strWeek)
			xf.SetCellValue(awhrSheetName, xlsx.RcToAxis(row, 7), activity*7)
			xf.SetCellValue(awhrSheetName, xlsx.RcToAxis(row, 8), 0)
			xf.SetCellValue(awhrSheetName, xlsx.RcToAxis(row, 9), activity)
			xf.SetCellValue(awhrSheetName, xlsx.RcToAxis(row, 10), comment)

			xf.SetCellValue(awhrSheetName, xlsx.RcToAxis(row, 12), actor.Company)
			xf.SetCellValue(awhrSheetName, xlsx.RcToAxis(row, 13), strings.Join(actor.Client, ", "))

			row++
		}
	}
	return xf.Write(w)
}

const (
	amtsXLSFile   string = "CRA MONTH _COMPANY_ _DATE_.xlsx"
	amtsSheetName string = "CRA"
	amtsRowDate   int    = 3
	amtsColDate   int    = 5
	amtsRowStart  int    = 9

	amtsColStart      int = 2
	amtsColMonthStart int = 5
	amtsColHours      int = 36
	amtsColExtraHours int = 37
	amtsColActingDays int = 38
	amtsColComment    int = 39
	amtsColCompany    int = 41
	amtsColClient     int = 42
	amtsColId         int = 43

	amtsColWEHeaderTpl       int = 44
	amtsColWEHoursTpl        int = 45
	amtsColHolidaysHeaderTpl int = 46
	amtsColHolidaysHoursTpl  int = 47
	amtsColOffHoursTpl       int = 48

	amtsNormalDay int = 0
	amtsWeekEnd   int = 1
	amtsHoliday   int = 2
	amtsNextMonth int = 3
)

func (te *DocTemplateEngine) GetActorsMonthlyTimeSheetTemplate(w io.Writer, actors []*actors.Actor, monthTimeSheet *timesheets.TimeSheet, publicHolidays map[string]string) error {
	file := filepath.Join(te.tmplDir, amtsXLSFile)
	xf, err := excelize.OpenFile(file)
	if err != nil {
		return err
	}

	styleHeaderWeekEnd, _ := xf.GetCellStyle(amtsSheetName, xlsx.RcToAxis(amtsRowStart-1, amtsColWEHeaderTpl))
	styleHoursWeekEnd, _ := xf.GetCellStyle(amtsSheetName, xlsx.RcToAxis(amtsRowStart-1, amtsColWEHoursTpl))
	styleHeaderHolidays, _ := xf.GetCellStyle(amtsSheetName, xlsx.RcToAxis(amtsRowStart-1, amtsColHolidaysHeaderTpl))
	styleHoursHolidays, _ := xf.GetCellStyle(amtsSheetName, xlsx.RcToAxis(amtsRowStart-1, amtsColHolidaysHoursTpl))
	styleHoursOff, _ := xf.GetCellStyle(amtsSheetName, xlsx.RcToAxis(amtsRowStart-1, amtsColOffHoursTpl))

	found := xf.GetSheetIndex(amtsSheetName)
	if found == 0 {
		return fmt.Errorf("could not find CRA sheet")
	}
	dates := make([]string, 31) // dates of month
	dayAttr := make([]int, 31)  // attr for days of month : 0 open, 1 WE, 2 NextMonth
	var endOfMonth bool
	for colNum := 0; colNum < 31; colNum++ {
		curDate := date.GetDateAfter(monthTimeSheet.WeekDate, colNum)
		if !endOfMonth && colNum > 27 && date.GetMonth(curDate) > monthTimeSheet.WeekDate {
			endOfMonth = true
		}
		celladdr := xlsx.RcToAxis(amtsRowStart-1, amtsColDate+colNum)
		if !endOfMonth {
			dates[colNum] = curDate
			switch {
			case publicHolidays[curDate] != "":
				dayAttr[colNum] = amtsHoliday
				_ = xf.SetCellStyle(amtsSheetName, celladdr, celladdr, styleHeaderHolidays)
			case date.GetDayNum(curDate) > 4:
				dayAttr[colNum] = amtsWeekEnd
				_ = xf.SetCellStyle(amtsSheetName, celladdr, celladdr, styleHeaderWeekEnd)
			}
		} else {
			dayAttr[colNum] = amtsNextMonth
			_ = xf.SetCellStr(amtsSheetName, celladdr, "")
		}
	}

	_ = xf.SetCellValue(amtsSheetName, xlsx.RcToAxis(amtsRowDate, amtsColDate), date.DateFrom(monthTimeSheet.WeekDate).ToTime())
	var countHour, countExtraHour, countActiveDays int
	row := amtsRowStart
	for _, actor := range actors {
		ats, exists := monthTimeSheet.ActorsTimes[actor.Id]
		if !exists {
			continue
		}
		_ = xf.SetCellStr(amtsSheetName, xlsx.RcToAxis(row, amtsColStart+0), actor.LastName)
		_ = xf.SetCellStr(amtsSheetName, xlsx.RcToAxis(row, amtsColStart+1), actor.FirstName)
		_ = xf.SetCellStr(amtsSheetName, xlsx.RcToAxis(row, amtsColStart+2), actor.Role)

		_ = xf.SetCellStr(amtsSheetName, xlsx.RcToAxis(row, amtsColCompany), actor.Company)
		_ = xf.SetCellStr(amtsSheetName, xlsx.RcToAxis(row, amtsColClient), strings.Join(actor.Client, ", "))
		_ = xf.SetCellInt(amtsSheetName, xlsx.RcToAxis(row, amtsColId), actor.Id)

		countHour = 0
		countExtraHour = 0
		countActiveDays = 0

	hoursLoop:
		for dayNum, hour := range ats.Hours {
			celladdr := xlsx.RcToAxis(row, amtsColMonthStart+dayNum)
			dayStatus := dayAttr[dayNum]
			if dayStatus == amtsNextMonth {
				// End of the month reached ... done for this actor
				break hoursLoop
			}
			currentDate := dates[dayNum]
			if !(currentDate >= actor.Period.Begin && !(actor.Period.End != "" && currentDate > actor.Period.End)) {
				// actor not active ... no hours allowed then
				_ = xf.SetCellStyle(amtsSheetName, celladdr, celladdr, styleHoursOff)
				continue hoursLoop
			}
		vacationLoop:
			for _, holiday := range actor.Vacation {
				if holiday.OverlapDate(currentDate) {
					dayStatus = amtsHoliday
					// skip any others vacation
					break vacationLoop
				}
			}

			switch dayStatus {
			case amtsHoliday:
				// in Holidays ... no hours allowed then
				_ = xf.SetCellStyle(amtsSheetName, celladdr, celladdr, styleHoursHolidays)
				continue hoursLoop
			case amtsWeekEnd:
				_ = xf.SetCellStyle(amtsSheetName, celladdr, celladdr, styleHoursWeekEnd)
				if hour > 0 {
					// Week End ... hours count as extra
					countExtraHour += hour
					countActiveDays++
				} else {
					continue hoursLoop
				}
			default:
				// normal workind day
				// extra hours when more than 7 hours
				nh := hour
				eh := 0
				if nh > 7 {
					nh = 7
					eh = hour - 7
				}
				countHour += nh
				countExtraHour += eh
				countActiveDays++
			}
			_ = xf.SetCellInt(amtsSheetName, celladdr, hour)
		}

		_ = xf.SetCellInt(amtsSheetName, xlsx.RcToAxis(row, amtsColHours), countHour)
		_ = xf.SetCellInt(amtsSheetName, xlsx.RcToAxis(row, amtsColExtraHours), countExtraHour)
		_ = xf.SetCellInt(amtsSheetName, xlsx.RcToAxis(row, amtsColActingDays), countActiveDays)

		row++
	}

	return xf.Write(w)
}
