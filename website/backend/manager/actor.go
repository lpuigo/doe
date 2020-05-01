package manager

import (
	"encoding/json"
	"fmt"
	"github.com/lpuig/ewin/doe/website/backend/model/actorinfos"
	"github.com/lpuig/ewin/doe/website/backend/model/actors"
	"io"
	"sort"
)

func (m Manager) GetActors(writer io.Writer) error {
	clientsNames := m.GetCurrentUserClientsName()
	actors := m.Actors.GetActorsByClient(false, clientsNames...)
	actorsHrs := m.ActorInfos.GetActorHRsByActors(actors, m.CurrentUser.HasPermissionHR())
	return json.NewEncoder(writer).Encode(actorsHrs)
}

func (m Manager) UpdateActors(updatedActors []*actorinfos.ActorHr) error {
	acts := make([]*actors.Actor, len(updatedActors))
	actInfos := make([]*actorinfos.ActorInfo, len(updatedActors))
	for i, actHr := range updatedActors {
		acts[i] = actHr.Actor
		actInfos[i] = actHr.Info
	}

	err := m.Actors.UpdateActors(acts)
	if err != nil {
		return err
	}
	if !m.CurrentUser.HasPermissionHR() {
		// If User has no HR permission, ignore returned ActorInfo data
		return nil
	}
	return m.ActorInfos.UpdateActorInfos(actInfos)
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
