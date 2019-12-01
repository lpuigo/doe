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

type IsFoasiteVisible func(s *FoaSite) bool

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
