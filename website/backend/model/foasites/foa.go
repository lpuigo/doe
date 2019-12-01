package foasites

import (
	"fmt"
	"github.com/lpuig/ewin/doe/website/backend/model/bpu"
	"github.com/lpuig/ewin/doe/website/backend/model/items"
	"github.com/lpuig/ewin/doe/website/frontend/model/foasite/foaconst"
	"strings"
)

type Foa struct {
	Id    int
	Ref   string
	Insee string
	Type  string

	State *State
}

func (f *Foa) SearchString() string {
	var searchBuilder strings.Builder
	fmt.Fprintf(&searchBuilder, "foa%s:%s,", "Ref", strings.ToUpper(f.Ref))
	fmt.Fprintf(&searchBuilder, "foa%s:%s,", "Insee", strings.ToUpper(f.Insee))
	fmt.Fprintf(&searchBuilder, "foa%s:%s,", "Type", strings.ToUpper(f.Type))
	fmt.Fprintf(&searchBuilder, "foa%s:%s,", "Date", strings.ToUpper(f.State.Date))
	return searchBuilder.String()
}

func (p *Foa) IsTodo() bool {
	switch p.State.Status {
	case foaconst.StateToDo:
		return true
	case foaconst.StateIncident:
		return true
	case foaconst.StateDone:
		return true
	default:
		return false
	}
}

func (p *Foa) IsDone() bool {
	switch p.State.Status {
	case foaconst.StateDone:
		return true
	default:
		return false
	}
}

func (p *Foa) IsBlocked() bool {
	switch p.State.Status {
	case foaconst.StateIncident:
		return true
	default:
		return false
	}
}

const (
	activityFoa     string = "Foa"
	catFoaInventory string = "Inventaire"
)

func (f *Foa) ExtendedRef() string {
	ref := f.Ref
	if f.Insee != "" {
		ref = f.Insee + " " + f.Ref
	}
	return ref
}

func (f *Foa) Itemize(client, site string, currentBpu *bpu.Bpu) ([]*items.Item, error) {
	res := []*items.Item{}

	foaArticles := currentBpu.GetCategoryArticles(activityFoa)

	todo, done, blocked := f.IsTodo(), f.IsDone(), f.IsBlocked()

	article, err := foaArticles.GetArticleFor(catFoaInventory, 1)
	if err != nil {
		return nil, fmt.Errorf("can not define foa inventory Item: %s", err.Error())
	}

	info := fmt.Sprintf("Inventaire chambre %s %s", f.Ref, f.Insee)
	if f.State.Comment != "" {
		info += fmt.Sprintf("\nCmt: %s", f.State.Comment)
	}

	ref := f.ExtendedRef()

	it := items.NewItem(
		client,
		site,
		activityFoa,
		ref,
		info,
		f.State.Date,
		"",
		article,
		1,
		1,
		todo,
		done,
		blocked,
	)
	it.Actors = f.State.Actors
	res = append(res, it)

	return res, nil
}
