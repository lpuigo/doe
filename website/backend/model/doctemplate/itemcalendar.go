package doctemplate

import (
	"sort"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/lpuig/ewin/doe/website/backend/model/clients"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	"github.com/lpuig/ewin/doe/website/backend/model/items"
	"github.com/lpuig/ewin/doe/website/backend/tools/xlsx"
)

type ItemCalendar struct {
	*items.Item
	ActorNames []string
}

// AddItemsCalendar adds a Sheet to given XLSx File with Items Calendar
func AddItemsCalendar(xf *excelize.File, its []*items.Item, actorById clients.ActorNameById) {
	sheetName := "Calendrier"
	xf.NewSheet(sheetName)

	var startDate, endDate string = "9999-12-31", "0000-00-00"
	actorCols := make(map[string]int)
	ics := []*ItemCalendar{}
	for _, item := range its {
		if !item.Done {
			continue
		}
		ic := &ItemCalendar{
			Item:       item,
			ActorNames: make([]string, len(item.Actors)),
		}
		if startDate > item.StartDate {
			startDate = item.StartDate
		}
		if endDate < item.Date {
			endDate = item.Date
		}
		for actNum, actorId := range item.Actors {
			actName := "_ Non Défini"
			if act := actorById(actorId); act != "" {
				actName = act
			}
			ic.ActorNames[actNum] = actName
			actorCols[actName] = 1
		}
		ics = append(ics, ic)
	}

	sort.Slice(ics, func(i, j int) bool {
		if ics[i].Date != ics[j].Date {
			return ics[i].Date < ics[j].Date
		}
		if ics[i].Client != ics[j].Client {
			return ics[i].Client < ics[j].Client
		}
		if ics[i].Site != ics[j].Site {
			return ics[i].Site < ics[j].Site
		}
		return ics[i].Name < ics[j].Name
	})
	actors := make([]string, len(actorCols))
	i := 0
	for actorName, _ := range actorCols {
		actors[i] = actorName
		i++
	}
	sort.Strings(actors)

	nbDays := date.NbDaysBetween(startDate, endDate)
	dateCols := make(map[string]int)

	// Header
	row := 1
	xf.SetCellStr(sheetName, xlsx.RcToAxis(row, 10), startDate)

	row++
	xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 1), "Client")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 2), "Site")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 3), "Activité")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 4), "Item")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 5), "Info")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 6), "Code BPU")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 7), "Quantité")
	xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 8), "Travail")
	for colNum := 0; colNum <= nbDays; colNum++ {
		col := 10 + colNum
		curDate := date.GetDateAfter(startDate, colNum)
		xf.SetCellStr(sheetName, xlsx.RcToAxis(row, col), curDate[8:10])
		dateCols[curDate] = col
	}
	for actorNum, actorName := range actors {
		col := 11 + nbDays + actorNum
		actorCols[actorName] = col
		xf.SetCellStr(sheetName, xlsx.RcToAxis(row, col), actorName)
	}

	// Content
	for _, item := range ics {
		row++
		xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 1), item.Client)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 2), item.Site)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 3), item.Activity)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 4), item.Name)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 5), item.Info)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 6), item.Article.Name)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 7), item.WorkQuantity)
		xf.SetCellValue(sheetName, xlsx.RcToAxis(row, 8), item.Work())
		work := item.Work()
		nbDays := date.NbDaysBetween(item.StartDate, item.Date)
		dayWork := work / float64(nbDays+1)

		for nbDay := 0; nbDay <= nbDays; nbDay++ {
			curDate := date.GetDateAfter(item.StartDate, nbDay)
			xf.SetCellFloat(sheetName, xlsx.RcToAxis(row, dateCols[curDate]), dayWork, 1, 64)
		}

		for _, actor := range item.ActorNames {
			xf.SetCellFloat(sheetName, xlsx.RcToAxis(row, actorCols[actor]), work/float64(len(item.ActorNames)), 1, 64)
		}
	}
}
