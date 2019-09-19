package polesites

import (
	"fmt"
	"github.com/lpuig/ewin/doe/website/backend/model/bpu"
	"github.com/lpuig/ewin/doe/website/backend/model/items"
	"github.com/lpuig/ewin/doe/website/frontend/model/polesite/poleconst"
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
	case poleconst.StateDictToDo:
		return true
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
	//case poleconst.StateDictToDo:
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

func (p *Pole) IsBlocked() bool {
	switch p.State {
	//case poleconst.StateNotSubmitted:
	case poleconst.StateNoGo:
		return true
	//case poleconst.StateDictToDo:
	//case poleconst.StateToDo:
	//case poleconst.StateHoleDone:
	case poleconst.StateIncident:
		return true
	//case poleconst.StateDone:
	//case poleconst.StateCancelled:
	default:
		return false
	}
}

const (
	activityPole    string = "Poteaux"
	catPoleCreation string = "Création"
)

func (p *Pole) Itemize(client, site string, currentBpu *bpu.Bpu) ([]*items.Item, error) {
	res := []*items.Item{}

	poleArticles := currentBpu.GetCategoryArticles(activityPole)

	todo, done, blocked := p.IsTodo(), p.IsDone(), p.IsBlocked()

	article, err := poleArticles.GetArticleFor(catPoleCreation, p.Height)
	if err != nil {
		return nil, fmt.Errorf("can not define pole creation Item: %s", err.Error())
	}

	info := fmt.Sprintf("Création poteau %s %dm", p.Material, p.Height)
	if p.Comment != "" {
		info += fmt.Sprintf("\nCmt: %s", p.Comment)
		//strings.ReplaceAll(info, "\n", "\r\n")
	}

	it := items.NewItem(
		client,
		site,
		activityPole,
		p.Ref,
		info,
		p.Date,
		"",
		article,
		1,
		1,
		todo,
		done,
		blocked,
	)
	it.Actors = p.Actors
	res = append(res, it)

	for _, product := range p.Product {
		article, err := poleArticles.GetArticleFor(product, p.Height)
		if err != nil {
			return nil, fmt.Errorf("can not define pole product item: %s", err.Error())
		}

		it := items.NewItem(
			client,
			site,
			activityPole,
			p.Ref,
			fmt.Sprintf("prestation complémentaire %s", product),
			p.Date,
			"",
			article,
			1,
			1,
			todo,
			done,
			blocked,
		)
		it.Actors = p.Actors
		res = append(res, it)
	}

	return res, nil
}
