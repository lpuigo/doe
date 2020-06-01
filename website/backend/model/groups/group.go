package groups

type Group struct {
	Id               int
	Name             string
	ActorDailyWork   float64
	ActorDailyIncome float64
}

//type ClientByName func(clientName string) *Group
//type IsTeamVisible func(ClientTeam) bool
//type TeamNameByMember func(string) string
//type ActorById func(id string) string
//type ActorInfoById func(id string) []string

//type IsActorVisible func(ClientActor) string

func NewGroup(name string) *Group {
	return &Group{
		Id:   0,
		Name: name,
	}
}
