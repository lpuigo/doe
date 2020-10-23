package items

import (
	"fmt"
	"github.com/lpuig/ewin/doe/website/backend/model/bpu"
)

type Item struct {
	Client       string // Client name
	Site         string // Site name
	Activity     string // Racco, Tirage, ...
	Name         string // PTxxx, Cablezzz, ...
	Info         string // BoxType + nbFO ...
	Date         string
	StartDate    string
	AttachDate   string
	Team         string
	Comment      string
	Actors       []string
	Article      *bpu.Article
	Quantity     int
	WorkQuantity int
	DivideBy     int
	Todo         bool
	Done         bool
	Blocked      bool
	Billed       bool
}

func NewItem(client, site, activity, name, info, date, team string, chapter *bpu.Article, quantity, workQuantity int, todo, done, blocked, billed bool) *Item {
	return &Item{
		Client:       client,
		Site:         site,
		Activity:     activity,
		Name:         name,
		Info:         info,
		Date:         date,
		StartDate:    date,
		Team:         team,
		Actors:       []string{},
		Article:      chapter,
		Quantity:     quantity,
		WorkQuantity: workQuantity,
		DivideBy:     1,
		Todo:         todo,
		Done:         done,
		Blocked:      blocked,
		Billed:       billed,
	}
}

func (i *Item) Clone() *Item {
	return &Item{
		Client:       i.Client,
		Site:         i.Site,
		Activity:     i.Activity,
		Name:         i.Name,
		Info:         i.Info,
		Date:         i.Date,
		StartDate:    i.StartDate,
		AttachDate:   i.AttachDate,
		Team:         i.Team,
		Actors:       i.Actors[:],
		Article:      i.Article,
		Quantity:     i.Quantity,
		WorkQuantity: i.WorkQuantity,
		DivideBy:     i.DivideBy,
		Todo:         i.Todo,
		Done:         i.Done,
		Blocked:      i.Blocked,
		Billed:       i.Billed,
	}
}

// SplitByActors returns a slice containing 1 copy of given item per actors
func (i *Item) SplitByActors() []*Item {
	res := make([]*Item, len(i.Actors))
	if len(res) <= 1 {
		return []*Item{i}
	}
	nbActors := i.DivideBy * len(i.Actors)
	for numAct, actId := range i.Actors {
		actItem := i.Clone()
		actItem.Actors = []string{actId}
		actItem.DivideBy = nbActors
		res[numAct] = actItem
	}
	return res
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
	Blocked: %t
	Billed: %t
`, i.Client, i.Site, i.Activity, i.Name, i.Info, i.Date, i.Team, i.Article.Name, i.Quantity, i.Todo, i.Done, i.Blocked, i.Billed)
}

// Price returns the price for the given item
func (i *Item) Price() float64 {
	res := i.Article.Price * float64(i.Quantity)
	if i.DivideBy > 1 {
		return res / float64(i.DivideBy)
	}
	return res
}

// Price returns the price for the given item
func (i *Item) Work() float64 {
	res := i.Article.Work * float64(i.WorkQuantity)
	if i.DivideBy > 1 {
		return res / float64(i.DivideBy)
	}
	return res
}
