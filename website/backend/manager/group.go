package manager

import (
	"encoding/json"
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

func (m Manager) GetGroupByName() groups.GroupByName {
	groupByName := make(map[string]*groups.Group)
	for _, group := range m.Groups.GetGroups() {
		groupByName[group.Name] = group
	}
	return func(groupName string) *groups.Group {
		return groupByName[groupName]
	}
}

func (m Manager) GroupSizeOnMonth(days []string) map[string][]int {
	month := date.DateStringRange{
		Begin: days[0],
		End:   days[len(days)-1],
	}
	res := make(map[string][]int)
	for _, group := range m.Groups.GetGroups() {
		res[group.Name] = make([]int, len(days))
	}

	for _, actor := range m.Actors.GetAllActors() {
		if !actor.IsActiveOnDateRange(month) {
			continue
		}
		for dayNum, day := range days {
			if !actor.IsActiveOn(day) { // TODO manage actors holidays
				continue
			}

			groupId := actor.Groups.ActiveGroupOnDate(day)
			if groupId == -1 {
				continue
			}
			gr := m.Groups.GetById(groupId)
			if gr == nil {
				continue
			}
			res[gr.Name][dayNum]++
		}
	}
	return res
}
