package foasites

import (
	"fmt"
	"github.com/lpuig/ewin/doe/website/backend/model/bpu"
	"github.com/lpuig/ewin/doe/website/backend/model/items"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/model/foasite/foaconst"
	"strings"
)

type FoaSite struct {
	Id         int
	Client     string
	Ref        string
	Manager    string
	OrderDate  string
	UpdateDate string
	Status     string
	Comment    string

	Foas []*Foa
}

func NewFoaSite() *FoaSite {
	return &FoaSite{
		Id:         -1,
		Client:     "",
		Ref:        "",
		Manager:    "",
		OrderDate:  "",
		UpdateDate: "",
		Status:     "",
		Comment:    "",
		Foas:       []*Foa{},
	}
}

func (fs *FoaSite) AddFoa(f *Foa) {
	f.Id = len(fs.Foas)
	fs.Foas = append(fs.Foas, f)
}

func (fs *FoaSite) GetInfo() *fm.FoaSiteInfo {
	fsi := fm.NewBEFoaSiteInfo()

	fsi.Id = fs.Id
	fsi.Client = fs.Client
	fsi.Ref = fs.Ref
	fsi.Manager = fs.Manager
	fsi.OrderDate = fs.OrderDate
	fsi.UpdateDate = fs.UpdateDate
	fsi.Status = fs.Status
	fsi.Comment = fs.Comment

	fsi.NbFoa, fsi.NbFoaBlocked, fsi.NbFoaDone = fs.GetFoasNumbers()

	var searchBuilder strings.Builder
	fmt.Fprintf(&searchBuilder, "%s:%s,", "Client", strings.ToUpper(fs.Client))
	fmt.Fprintf(&searchBuilder, "%s:%s,", "Ref", strings.ToUpper(fs.Ref))
	fmt.Fprintf(&searchBuilder, "%s:%s,", "Manager", strings.ToUpper(fs.Manager))
	fmt.Fprintf(&searchBuilder, "%s:%s,", "OrderDate", strings.ToUpper(fs.OrderDate))
	fmt.Fprintf(&searchBuilder, "%s:%s,", "Comment", strings.ToUpper(fs.Comment))
	for _, foa := range fs.Foas {
		fmt.Fprintf(&searchBuilder, "%s,", foa.SearchString())
	}
	fsi.Search = searchBuilder.String()

	return fsi
}

// GetFoasNumbers returns total, blocked and done number of Foas
func (fs *FoaSite) GetFoasNumbers() (total, blocked, done int) {
	for _, f := range fs.Foas {
		switch f.State.Status {
		//case poleconst.StateNotSubmitted:
		case foaconst.StateToDo:
			total++
		case foaconst.StateIncident:
			total++
			blocked++
		case foaconst.StateDone, foaconst.StateAttachment:
			total++
			done++
		}
	}
	return
}

type IsFoaSiteVisible func(s *FoaSite) bool

// Itemize returns slice of item pertaining to polesite poles list
func (fs *FoaSite) Itemize(currentBpu *bpu.Bpu) ([]*items.Item, error) {
	res := []*items.Item{}

	for _, foa := range fs.Foas {
		items, err := foa.Itemize(fs.Client, fs.Ref, currentBpu)
		if err != nil {
			return nil, err
		}
		res = append(res, items...)
	}
	return res, nil
}

// AddStat adds Stats into values for given Polesite
func (fs *FoaSite) AddStat(stats items.Stats, sc items.StatContext, currentBpu *bpu.Bpu, showprice bool) error {

	addValue := func(date, serie string, actors []string, value float64) {
		stats.AddStatValue(fs.Ref, fs.Client, date, "", serie, value)
		if sc.ShowTeam && len(actors) > 0 {
			value /= float64(len(actors))
			for _, actName := range actors {
				stats.AddStatValue(fs.Ref, fs.Client+" : "+actName, date, "", serie, value)
			}
		}
	}

	calcItems, err := fs.Itemize(currentBpu)
	if err != nil {
		return fmt.Errorf("error on foa stat itemize for '%s':%s", fs.Ref, err.Error())
	}
	for _, item := range calcItems {
		if !item.Done {
			continue
		}
		actorsName := make([]string, len(item.Actors))
		for i, actId := range item.Actors {
			actorsName[i] = sc.ActorById(actId)
		}
		addValue(sc.DateFor(item.Date), items.StatSerieWork, actorsName, item.Work())
		if showprice {
			addValue(sc.DateFor(item.Date), items.StatSeriePrice, actorsName, item.Price())
		}
	}
	return nil
}
