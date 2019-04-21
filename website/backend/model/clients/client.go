package clients

type Client struct {
	Id       int
	Name     string
	Teams    []Team
	Articles []Article
}

func (c Client) GetArticleNames() []string {
	res := []string{}
	for _, a := range c.Articles {
		res = append(res, a.Name)
	}
	return res
}

func NewClient(name string) *Client {
	return &Client{
		Id:    0,
		Name:  name,
		Teams: []Team{},
	}
}
