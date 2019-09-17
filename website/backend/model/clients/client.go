package clients

import "github.com/lpuig/ewin/doe/website/backend/model/bpu"

type Client struct {
	Id    int
	Name  string
	Teams []Team
	*bpu.Bpu
}

type ClientTeam struct {
	Client string
	Team   string
}

//type ClientActor struct {
//	Client string
//	Actor  string
//}

type ClientByName func(clientName string) *Client
type IsTeamVisible func(ClientTeam) bool
type TeamNameByMember func(string) string
type ActorById func(id string) string
type ActorInfoById func(id string) []string

//type IsActorVisible func(ClientActor) string

func NewClient(name string) *Client {
	return &Client{
		Id:    0,
		Name:  name,
		Teams: []Team{},
		Bpu:   bpu.NewBpu(),
	}
}

func (c Client) GetOrangeArticleNames() []string {
	res := []string{}
	for _, a := range c.GetOrangeArticles() {
		res = append(res, a.Name)
	}
	return res
}

func (c Client) GetOrangeArticles() []*bpu.Article {
	if c.Bpu == nil {
		return nil
	}
	ca := c.Bpu.GetCategoryArticles("Orange")
	if ca == nil {
		return []*bpu.Article{}
	}
	return ca.GetArticles("El")
}

func (c Client) GenTeamNameByMember() TeamNameByMember {
	teamName := map[string]string{}
	activeMembers := map[string]string{}
	for _, team := range c.Teams {
		if team.IsActive {
			activeMembers[team.Name] = team.Members
		}
		teamName[team.Members] = team.Name
	}
	for member, team := range teamName {
		teamName[member] = team + " (" + activeMembers[team] + ")"
	}
	return func(member string) string {
		team, found := teamName[member]
		if !found {
			return "Unknown"
		}
		return team
	}
}
