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
	actors := m.GetCurrentUserActors()
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
	// Update actInfos actorId according to updated actors Id
	for i, actInfo := range actInfos {
		if actInfo.ActorId < 0 {
			actInfo.ActorId = acts[i].Id
		}
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
		if actors[i].Company == actors[j].Company {
			return actors[i].Ref < actors[j].Ref
		}
		return actors[i].Company < actors[j].Company
	})
	actorRhs := m.ActorInfos.GetActorHRsByActors(actors, true)
	monthlytimesheet := m.TimeSheets.GetMonthlyTimeSheetFor(date, actors)
	return m.TemplateEngine.GetActorsMonthlyTimeSheetTemplate(writer, actorRhs, m.GenGroupById(), monthlytimesheet, m.DaysOff.days)
}
