package items

import (
	"fmt"
	"github.com/lpuig/ewin/doe/website/backend/model/bpu"
)

type Item struct {
	Client       string // Site name
	Site         string // Site name
	Activity     string // Racco, Tirage, ...
	Name         string // PTxxx, Cablezzz, ...
	Info         string // BoxType + nbFO ...
	Date         string
	Team         string
	Actors       []string
	Article      *bpu.Article
	Quantity     int
	WorkQuantity int
	Todo         bool
	Done         bool
}

func NewItem(client, site, activity, name, info, date, team string, chapter *bpu.Article, quantity, workQuantity int, todo, done bool) *Item {
	return &Item{
		Client:       client,
		Site:         site,
		Activity:     activity,
		Name:         name,
		Info:         info,
		Date:         date,
		Team:         team,
		Article:      chapter,
		Quantity:     quantity,
		WorkQuantity: workQuantity,
		Todo:         todo,
		Done:         done,
	}
}

func (i *Item) String() string {
	return fmt.Sprintf(`Client: %s, Site: %s, Activity: %s Name: %s
	Info: %s
	Date: %s
	Team: %s
	Article: %s
	Quantity: %d
	Todo: %t
	Done: %t
`, i.Client, i.Site, i.Activity, i.Name, i.Info, i.Date, i.Team, i.Article.Name, i.Quantity, i.Todo, i.Done)
}

// Price returns the price for the given item
func (i *Item) Price() float64 {
	return i.Article.Price * float64(i.Quantity)
}

// Price returns the price for the given item
func (i *Item) Work() float64 {
	return i.Article.Work * float64(i.WorkQuantity)
}
