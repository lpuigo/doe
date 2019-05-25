package clients

import "github.com/lpuig/ewin/doe/website/backend/model/bpu"

type Client struct {
	Id    int
	Name  string
	Teams []Team
	*bpu.Bpu
}

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
