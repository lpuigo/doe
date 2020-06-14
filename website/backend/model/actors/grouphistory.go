package actors

import "sort"

type GroupHistory map[string]int

func NewGroupHistory() GroupHistory {
	return make(GroupHistory)
}

// ActiveGroupOnDate returns the assigned group" id on given date (-1 if no group found)
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
