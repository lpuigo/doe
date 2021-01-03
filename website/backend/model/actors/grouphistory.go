package actors

import (
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	"sort"
)

type GroupHistory map[string]int

func NewGroupHistory() GroupHistory {
	return make(GroupHistory)
}

type dategroupid struct {
	date string
	gid  int
}

func (gh GroupHistory) getDateGroupList() []dategroupid {
	res := []dategroupid{}
	if len(gh) == 0 {
		return res
	}
	for date, gid := range gh {
		res = append(res, dategroupid{
			date: date,
			gid:  gid,
		})
	}
	sort.Slice(res, func(i, j int) bool {
		if res[i].date < res[j].date {
			return true
		}
		return false
	})
	return res
}

// ActiveGroupOnDate returns the assigned group' id on given date (-1 if no group assigned)
//
// if day is prior to first assignement, group 0 is returned
func (gh GroupHistory) ActiveGroupOnDate(day string) int {
	if len(gh) == 0 {
		return -1
	}
	if len(gh) == 1 {
		for _, i := range gh {
			return i
		}
	}
	effectiveGrId := -1
	effectiveGrDay := date.TimeJSMinDate
	for assignDay, grId := range gh {
		if assignDay > day {
			// this group is not applicable yet at given day, skip
			continue
		}
		if assignDay < effectiveGrDay { // this group is older than current effective one, skip
			continue
		}
		// current group is candidate
		effectiveGrDay = assignDay
		effectiveGrId = grId
	}
	if effectiveGrId == -1 { // no group was found, apply default group 0
		effectiveGrId = 0
	}
	return effectiveGrId
}

// AssignedGroup returns all assigned group' id (nil if no group assigned)
func (gh GroupHistory) AssignedGroup() []int {
	if len(gh) == 0 {
		return nil
	}
	if len(gh) == 1 {
		for _, i := range gh {
			return []int{i}
		}
	}
	grps := make(map[int]bool)
	for _, grId := range gh {
		grps[grId] = true
	}
	res := make([]int, len(grps))
	i := 0
	for grId, _ := range grps {
		res[i] = grId
		i++
	}
	return res
}

// ActiveGroupPerDay returns slice of active groupId for each given days (group id 0 is used as default if no group is assigned)
func (gh GroupHistory) ActiveGroupPerDay(days []string) []int {
	res := make([]int, len(days))
	if len(gh) == 0 {
		return res
	}
	groups := gh.getDateGroupList()
	currentGroupId := 0
	nextGroupId := 0
	for i, group := range groups {
		nextGroupId = i
		if days[0] < group.date {
			break
		}
		currentGroupId = group.gid
	}
	res[0] = currentGroupId
	for i, day := range days[1:] {
		if nextGroupId < len(groups) && day >= groups[nextGroupId].date {
			currentGroupId = groups[nextGroupId].gid
			nextGroupId++
		}
		res[i+1] = currentGroupId
	}
	return res
}
