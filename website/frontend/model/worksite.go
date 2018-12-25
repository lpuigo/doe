package model

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/dates"
	"strings"
)

type Worksite struct {
	*js.Object

	Id        int      `js:"Id"`
	Ref       string   `js:"Ref"`
	OrderDate string   `js:"OrderDate"`
	City      string   `js:"City"`
	Status    string   `js:"Status"`
	Pmz       *PT      `js:"Pmz"`
	Pa        *PT      `js:"Pa"`
	Comment   string   `js:"Comment"`
	Orders    []*Order `js:"Orders"`
}

func NewWorkSite() *Worksite {
	ws := &Worksite{Object: tools.O()}
	ws.Id = 0
	ws.Ref = ""
	ws.OrderDate = ""
	ws.City = ""
	ws.Status = ""
	ws.Pmz = NewPT()
	ws.Pa = NewPT()
	ws.Comment = ""
	ws.Orders = []*Order{}

	return ws
}

func WorksiteFromJS(o *js.Object) *Worksite {
	ws := &Worksite{Object: o}
	return ws
}

func (ws *Worksite) Clone() *Worksite {
	nws := &Worksite{Object: tools.O()}
	nws.Copy(ws)
	return nws
}

func (ws *Worksite) Copy(ows *Worksite) {
	ws.Id = ows.Id
	ws.Ref = ows.Ref
	ws.OrderDate = ows.OrderDate
	ws.City = ows.City
	ws.Status = ows.Status
	ws.Pmz = ows.Pmz.Clone()
	ws.Pa = ows.Pa.Clone()
	ws.Comment = ows.Comment
	ws.Orders = make([]*Order, len(ows.Orders))
	for i, o := range ows.Orders {
		ws.Orders[i] = o.Clone()
	}
}

func (ws *Worksite) TextFiltered(filter string) bool {
	expected := true
	if filter == "" {
		return true
	}
	if strings.HasPrefix(filter, `\`) {
		if len(filter) > 1 { // prevent from filtering all when only '\' is entered
			expected = false
		}
		filter = filter[1:]
	}
	return ws.Contains(filter) == expected

}

func (ws *Worksite) Contains(str string) bool {
	if str == "" {
		return true
	}
	return strings.Contains(strings.ToLower(ws.SearchInString()), strings.ToLower(str))
}

func (ws *Worksite) SearchInString() string {
	//res += "Id:" Skipped on purpose
	res := "Ref:" + ws.Ref + "\n"
	res += "OrderDate:" + date.DateString(ws.OrderDate) + "\n"
	res += "City:" + ws.City + "\n"
	res += "Status:" + ws.Status + "\n"
	res += "Pmz:" + ws.Pmz.SearchInString()
	res += "Pa:" + ws.Pa.SearchInString()
	res += "Comment:" + ws.Ref + "\n"

	for _, o := range ws.Orders {
		res += o.SearchInString()
	}
	return res
}

func (ws *Worksite) GetInfo() (nbCommand, nbTroncon, nbRacco int) {
	nbCommand = len(ws.Orders)
	for _, o := range ws.Orders {
		nbTroncon += len(o.Troncons)
		for _, t := range o.Troncons {
			nbRacco += t.NbRacco
		}
	}
	return
}

func (ws *Worksite) DeleteOrder(i int) {
	orders := []*Order{}
	for j, o := range ws.Orders {
		if i == j {
			continue
		}
		orders = append(orders, o)
	}
	ws.Orders = orders
}
