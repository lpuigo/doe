package clients

type Client struct {
	Id    int
	Name  string
	Teams []Team
}

func NewClient(name string) *Client {
	return &Client{
		Id:    0,
		Name:  name,
		Teams: []Team{},
	}
}
