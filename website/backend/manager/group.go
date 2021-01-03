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
	return m.Groups.GetGroups()
}

// GroupSizePerDays returns a map of "groupName" and "groupName : actorLabel" -> []int giving number of active actors for each day within given days slice
func (m Manager) GroupSizePerDays(days []string) map[string][]int {
	daysRange := date.DateStringRange{
		Begin: days[0],
		End:   days[len(days)-1],
	}
	res := make(map[string][]int)
	groupNameDict := make(map[int]string)
	for _, group := range m.Groups.GetGroups() {
		res[group.Name] = make([]int, len(days))
		groupNameDict[group.Id] = group.Name
	}

	for _, actor := range m.Actors.GetAllActors() {
		actorLabel := actor.GetLabel()
		if !actor.IsActiveOnDateRange(daysRange) {
			continue
		}
		groupIdPerDay, groupIds := actor.Groups.ActiveGroupPerDay(days)
		for _, groupId := range groupIds {
			res[groupNameDict[groupId]+" : "+actorLabel] = make([]int, len(days))
		}
		for dayNum, day := range days {
			if !actor.IsWorkingOn(day) {
				continue
			}
			grpName := groupNameDict[groupIdPerDay[dayNum]]
			res[grpName][dayNum]++
			res[grpName+" : "+actorLabel][dayNum] = 1
		}
	}
	// force default group activity to 0
	res[groupNameDict[0]] = make([]int, len(days))
	return res
}
