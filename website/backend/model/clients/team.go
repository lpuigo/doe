package clients

type Team struct {
	Name     string
	Members  string
	IsActive bool
}

func MakeTeam(name, members string) Team {
	return Team{
		Name:     name,
		Members:  members,
		IsActive: false,
	}
}
