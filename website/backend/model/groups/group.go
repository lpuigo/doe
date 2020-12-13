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

// GroupByName is a getter function to retrieve group by name. returns nil if group's name not found
type GroupByName func(groupName string) *Group

// GroupById is a getter function to retrieve group by id. returns nil if group's id not found
type GroupById func(groupId int) *Group
