package actor

type GroupHistory map[string]int

func NewGroupHistory() GroupHistory {
	return make(GroupHistory)
}

// GetCurrentInfo returns current GroupId and Assignement date (if no group is found, groupId is set to -1)
func (gh GroupHistory) GetCurrentInfo() (currentId int, assignDate string) {
	defaultDate := "0001-01-01"
	assignDate = defaultDate
	for d, id := range gh {
		if d > assignDate {
			assignDate = d
			currentId = id
		}
	}
	if assignDate == defaultDate {
		currentId = -1
	}
	return
}

func (gh GroupHistory) Copy() GroupHistory {
	ngh := NewGroupHistory()
	for assignDate, gId := range gh {
		ngh[assignDate] = gId
	}
	return ngh
}