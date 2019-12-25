package manager

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

func (m Manager) GetActors(writer io.Writer) error {
	clientsNames := m.GetCurrentUserClientsName()
	actors := m.Actors.GetActorsByClient(false, clientsNames...)
	return json.NewEncoder(writer).Encode(actors)
}

func (m Manager) GetActorsWorkingHoursRecordXLSName(monthDate string) string {
	return fmt.Sprintf("CRA %s.xlsx", monthDate)
}

func (m Manager) GetActorsWorkingHoursRecordXLS(writer io.Writer, date string) error {
	actors := m.Actors.GetAllActors()
	sort.Slice(actors, func(i, j int) bool {
		return actors[i].Ref < actors[j].Ref
	})
	return m.TemplateEngine.GetActorsWorkingHoursRecordXLS(writer, date, actors)
}

func (m Manager) GetActorsMonthlyTimeSheetsXLS(writer io.Writer, date string) error {
	actors := m.Actors.GetAllActors()
	sort.Slice(actors, func(i, j int) bool {
		return actors[i].Ref < actors[j].Ref
	})
	monthlytimesheet := m.TimeSheets.GetMonthlyTimeSheetFor(date, actors)
	return m.TemplateEngine.GetActorsMonthlyTimeSheetTemplate(writer, actors, monthlytimesheet, m.DaysOff.days)
}
