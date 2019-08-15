package polesites

import (
	"fmt"
	"github.com/lpuig/ewin/doe/website/backend/model/bpu"
	"github.com/lpuig/ewin/doe/website/backend/model/clients"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	"github.com/lpuig/ewin/doe/website/backend/model/items"
	"github.com/lpuig/ewin/doe/website/frontend/model/polesite/poleconst"
	"sort"
	"strconv"
	"strings"
)

type Pole struct {
	Id             int
	Ref            string
	City           string
	Address        string
	Sticker        string
	Lat            float64
	Long           float64
	State          string
	Date           string
	Actors         []string
	DtRef          string
	DictRef        string
	DictDate       string
	DictInfo       string
	Height         int
	Material       string
	AspiDate       string
	Kizeo          string
	Comment        string
	Product        []string
	AttachmentDate string
}

func (p *Pole) SearchString() string {
	var searchBuilder strings.Builder
	fmt.Fprintf(&searchBuilder, "pole%s:%s,", "Ref", strings.ToUpper(p.Ref))
	fmt.Fprintf(&searchBuilder, "pole%s:%s,", "City", strings.ToUpper(p.City))
	fmt.Fprintf(&searchBuilder, "pole%s:%s,", "Address", strings.ToUpper(p.Address))
	fmt.Fprintf(&searchBuilder, "pole%s:%s,", "DtRef", strings.ToUpper(p.DtRef))
	fmt.Fprintf(&searchBuilder, "pole%s:%s,", "DictRef", strings.ToUpper(p.DictRef))
	fmt.Fprintf(&searchBuilder, "pole%s:%s,", "Height", strconv.Itoa(p.Height)+"M")
	fmt.Fprintf(&searchBuilder, "pole%s:%s,", "Material", strings.ToUpper(p.Material))
	for _, key := range p.Product {
		fmt.Fprintf(&searchBuilder, "poleProduct:%s,", strings.ToUpper(key))
	}
	for _, actor := range p.Actors {
		fmt.Fprintf(&searchBuilder, "poleActor:%s,", strings.ToUpper(actor))
	}
	fmt.Fprintf(&searchBuilder, "pole%s:%s,", "DictInfo", strings.ToUpper(p.DictInfo))
	fmt.Fprintf(&searchBuilder, "pole%s:%s,", "Date", strings.ToUpper(p.Date))
	fmt.Fprintf(&searchBuilder, "pole%s:%s,", "AspiDate", strings.ToUpper(p.AspiDate))
	fmt.Fprintf(&searchBuilder, "pole%s:%s,", "AttachDate", strings.ToUpper(p.AttachmentDate))
	return searchBuilder.String()
}

func (p *Pole) IsTodo() bool {
	switch p.State {
	//case poleconst.StateNotSubmitted:
	//case poleconst.StateNoGo:
	case poleconst.StateToDo:
		return true
	case poleconst.StateHoleDone:
		return true
	case poleconst.StateIncident:
		return true
	case poleconst.StateDone:
		return true
	//case poleconst.StateCancelled:
	default:
		return false
	}
}

func (p *Pole) IsDone() bool {
	switch p.State {
	//case poleconst.StateNotSubmitted:
	//case poleconst.StateNoGo:
	//case poleconst.StateToDo:
	//case poleconst.StateHoleDone:
	//case poleconst.StateIncident:
	case poleconst.StateDone:
		return true
	//case poleconst.StateCancelled:
	default:
		return false
	}
}

const (
	activityPole    string = "Poteaux"
	catPoleCreation string = "Création"
)

func (p *Pole) Itemize(currentBpu *bpu.Bpu, actorById clients.ActorById) ([]*items.Item, error) {
	res := []*items.Item{}

	poleArticles := currentBpu.GetCategoryArticles(activityPole)

	todo, done := p.IsTodo(), p.IsDone()

	actors := []string{}
	for _, actorId := range p.Actors {
		actors = append(actors, actorById(actorId))
	}
	sort.Strings(actors)

	article, err := poleArticles.GetArticleFor(catPoleCreation, p.Height)
	if err != nil {
		return nil, fmt.Errorf("can not define pole creation Item: %s", err.Error())
	}

	it := items.NewItem(
		activityPole,
		p.Ref,
		fmt.Sprintf("Création poteau %s", p.Ref),
		p.Date,
		strings.Join(actors, ", "),
		article,
		1,
		1,
		todo,
		done,
	)
	it.Actors = actors
	res = append(res, it)

	for _, product := range p.Product {
		article, err := poleArticles.GetArticleFor(product, p.Height)
		if err != nil {
			return nil, fmt.Errorf("can not define pole product Item: %s", err.Error())
		}

		it := items.NewItem(
			activityPole,
			p.Ref,
			fmt.Sprintf("prestation complémentaire %s pour %s", product, p.Ref),
			p.Date,
			strings.Join(actors, ", "),
			article,
			1,
			1,
			todo,
			done,
		)
		it.Actors = actors
		res = append(res, it)
	}

	return res, nil
}

// AddStat adds nb of El installed per date (in map[date]nbEl) by visible Client & Client : Teams
func (p *Pole) AddStat(values map[items.StatKey]float64, dateFor date.DateAggreg, isActorVisible clients.IsTeamVisible,
	currentBpu *bpu.Bpu, teamName clients.TeamNameByMember, showprice bool) error {

	//TODO to implement

	return nil
}
