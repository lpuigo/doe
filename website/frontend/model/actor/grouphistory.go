package actor

import (
	"github.com/lpuig/ewin/doe/website/frontend/tools/date"
	"sort"
)

type GroupHistory map[string]int

func NewGroupHistory() GroupHistory {
	return make(GroupHistory)
}

// GetCurrentInfo returns current GroupId and Assignement date (if no group is found, groupId is set to -1)
func (gh GroupHistory) GetCurrentInfo() (currentId int, assignDate string) {
	defaultDate := "0001-01-01"
	assignDate = defaultDate
	currentDate := date.TodayAfter(0)
	for d, id := range gh {
		if d <= currentDate {
			assignDate = d
			currentId = id
		}
	}
	if assignDate == defaultDate {
		currentId = -1
	}
	return
}

type GroupAssign struct {
	Id     int
	Period *date.DateRange
}

func (gh GroupHistory) GetGroupAssignsInRange(dr *date.DateRange) []GroupAssign {
	if len(gh) == 0 {
		return []GroupAssign{}
	}
	res := []GroupAssign{}
	affectationDates := make([]string, len(gh))
	i := 0
	for d, _ := range gh {
		affectationDates[i] = d
		i++
	}
	sort.Strings(affectationDates)
	if dr.End < affectationDates[0] {
		return []GroupAssign{}
	}
	lastAffectation := affectationDates[len(affectationDates)-1]
	if dr.Begin >= lastAffectation {
		return []GroupAssign{{
			Id:     gh[lastAffectation],
			Period: date.NewDateRangeFrom(lastAffectation, "9999-12-31"),
		}}
	}
	endDate := "9999-12-31"
	for i := len(affectationDates) - 1; i >= 0; i-- {
		affDate := affectationDates[i]
		currentRange := date.NewDateRangeFrom(affDate, endDate)
		if currentRange.Overlap(dr) {
			res = append(res, GroupAssign{
				Id:     gh[affDate],
				Period: date.NewDateRangeFrom(affDate, endDate),
			})
		}
		endDate = date.After(affDate, -1)
	}
	return res
}

func (gh GroupHistory) Copy() GroupHistory {
	ngh := NewGroupHistory()
	for assignDate, gId := range gh {
		ngh[assignDate] = gId
	}
	return ngh
}
