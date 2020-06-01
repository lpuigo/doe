package actors

type GroupHistory map[string]int

func NewGroupHistory() GroupHistory {
	return make(GroupHistory)
}
