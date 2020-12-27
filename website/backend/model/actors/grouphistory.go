package actors

import "sort"

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

// ActiveGroupOnDate returns the assigned group" id on given date (-1 if no group assigned)
func (gh GroupHistory) ActiveGroupOnDate(day string) int {
	if len(gh) == 0 {
		return -1
	}
	if len(gh) == 1 {
		for _, i := range gh {
			return i
		}
	}
	dates := make([]string, len(gh))
	i := 0
	for assignDate, _ := range gh {
		dates[i] = assignDate
		i++
	}
	sort.Strings(dates)
	effDate := ""
	for _, assignDate := range dates {
		if day >= assignDate {
			effDate = assignDate
		}
	}
	if effDate == "" {
		effDate = dates[0]
	}
	return gh[effDate]
}

// ActiveGroupPerDay returns slice of active group for each given days (group id 0 is used as default if no group is assigned)
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
