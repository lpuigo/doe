package manager

import (
	"encoding/json"
	"io"

	"github.com/lpuig/ewin/doe/website/backend/model/date"
	"github.com/lpuig/ewin/doe/website/backend/model/groups"
)

func (m Manager) GetGroups(writer io.Writer) error {
	grps := m.Groups.GetGroups()
	return json.NewEncoder(writer).Encode(grps)
}

func (m Manager) UpdateGroups(updatedGroups []*groups.Group) error {
	return m.Groups.UpdateGroups(updatedGroups)
}

// GetCurrentUserVisibleGroups returns a slice of groups.Group visible by current user
func (m Manager) GetCurrentUserVisibleGroups() []*groups.Group {
	res := []*groups.Group{}
	getGroupById := m.GenGroupById()
	if len(m.CurrentUser.Groups) > 0 { // Current User has a restricted group visibility
		for _, grpId := range m.CurrentUser.Groups {
			grp := getGroupById(grpId)
			if grp == nil {
				continue
			}
			res = append(res, grp)
		}
		return res
	}
	// Current user has a global group visibility
	for _, grp := range m.Groups.GetGroups() {
		res = append(res, grp)
	}
	return res
}

// GroupSizePerDays returns a map of groupName -> []int giving number of active actors for groupName for each day within given days slice
func (m Manager) GroupSizePerDays(days []string) map[string][]int {
	daysRange := date.DateStringRange{
		Begin: days[0],
		End:   days[len(days)-1],
	}
	res := make(map[string][]int)
	for _, group := range m.Groups.GetGroups() {
		res[group.Name] = make([]int, len(days))
	}

	for _, actor := range m.Actors.GetAllActors() {
		if !actor.IsActiveOnDateRange(daysRange) {
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
