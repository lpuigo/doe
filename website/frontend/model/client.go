package model

import (
	"github.com/gopherjs/gopherjs/js"
	"strconv"
	"strings"
)

type Client struct {
	*js.Object

	Name     string   `js:"Name"`
	Teams    []*Team  `js:"Teams"`
	Actors   []*Actor `js:"Actors"`
	Articles []string `js:"Articles"`
}

// GetActorBy tries to search client actors by id (if id is castable to int) or by lastname (return nil if not found)
func (c *Client) GetActorBy(id string) *Actor {
	if actID, err := strconv.Atoi(id); err == nil {
		for _, actor := range c.Actors {
			if actor.Id == actID {
				return actor
			}
		}
		return nil
	}
	actLastName := strings.ToUpper(id)
	for _, actor := range c.Actors {
		if strings.ToUpper(actor.LastName) == actLastName {
			return actor
		}
	}
	return nil
}
