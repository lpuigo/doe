package manager

import (
	"encoding/json"
	"github.com/lpuig/ewin/doe/website/backend/model/actors"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	"github.com/lpuig/ewin/doe/website/backend/model/groups"
	"io"
)

func (m Manager) GetGroups(writer io.Writer) error {
	grps := m.Groups.GetGroups()
	return json.NewEncoder(writer).Encode(grps)
}

func (m Manager) UpdateGroups(updatedGroups []*groups.Group) error {
	return m.Groups.UpdateGroups(updatedGroups)
}

func (m Manager) GenGroupByName() groups.GroupByName {
	groupByName := make(map[string]*groups.Group)
	for _, group := range m.Groups.GetGroups() {
		groupByName[group.Name] = group
	}
	return func(groupName string) *groups.Group {
		return groupByName[groupName]
	}
}

func (m Manager) GenActorsByGroupId() actors.ActorsByGroupId {
	today := date.Today().String()
	actorsByGroupId := make(map[int][]*actors.Actor)
	for _, actor := range m.Actors.GetAllActors() {
		groupId := actor.Groups.ActiveGroupOnDate(today)
		if groupId == -1 {
			continue
		}
		actorsList := actorsByGroupId[groupId]
		actorsByGroupId[groupId] = append(actorsList, actor)
	}
	return func(groupId int) []*actors.Actor {
		return actorsByGroupId[groupId]
	}
}

func (m Manager) GenGroupById() groups.GroupById {
	groupById := make(map[int]*groups.Group)
	for _, group := range m.Groups.GetGroups() {
		groupById[group.Id] = group
	}
	return func(groupId int) *groups.Group {
		return groupById[groupId]
	}
}

func (m Manager) GroupSizeOnMonth(days []string) map[string][]int {
	monthRange := date.DateStringRange{
		Begin: days[0],
		End:   days[len(days)-1],
	}
	res := make(map[string][]int)
	for _, group := range m.Groups.GetGroups() {
		res[group.Name] = make([]int, len(days))
	}

	for _, actor := range m.Actors.GetAllActors() {
		if !actor.IsActiveOnDateRange(monthRange) {
			continue
		}
		actorActivity := make([]int, len(days))
		actorGroups := actor.Groups.ActiveGroupPerDay(days)
		for dayNum, day := range days {
			if !actor.IsWorkingOn(day) {
				continue
			}
			actorActivity[dayNum] = 1
			gr := m.Groups.GetById(actorGroups[dayNum])
			if gr == nil {
				continue
			}
			res[gr.Name][dayNum]++
		}
		res[actor.GetLabel()] = actorActivity
	}
	return res
}
