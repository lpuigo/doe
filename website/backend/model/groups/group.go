package groups

type Group struct {
	Id               int
	Name             string
	Clients          []string
	ActorDailyWork   float64
	ActorDailyIncome float64
}

func NewGroup(name string) *Group {
	return &Group{
		Id:      0,
		Name:    name,
		Clients: []string{},
	}
}

type GroupByName func(groupName string) *Group
