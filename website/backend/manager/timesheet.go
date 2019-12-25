package manager

import (
	"encoding/json"
	"io"
)

func (m Manager) GetTimeSheet(writer io.Writer, weekdate string) error {
	clientsNames := m.GetCurrentUserClientsName()
	actors := m.Actors.GetActorsByClient(false, clientsNames...)
	timesheet, err := m.TimeSheets.GetTimeSheetFor(weekdate, actors)
	if err != nil {
		return err
	}
	return json.NewEncoder(writer).Encode(timesheet)
}
