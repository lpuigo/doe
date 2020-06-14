package actors

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
	id := -1
	curDay := day
	for assignDate, grId := range gh {
		if assignDate <= day && assignDate <= curDay {
			id = grId
			curDay = assignDate
		}
	}
	return id
}
