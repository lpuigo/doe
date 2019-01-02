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
	Client    string   `js:"Client"`
	Ref       string   `js:"Ref"`
	OrderDate string   `js:"OrderDate"`
	DoeDate   string   `js:"DoeDate"`
	City      string   `js:"City"`
	Status    string   `js:"Status"`
	Pmz       *PT      `js:"Pmz"`
	Pa        *PT      `js:"Pa"`
	Comment   string   `js:"Comment"`
	Orders    []*Order `js:"Orders"`
	Rework    *Rework  `js:"Rework"`
	Dirty     bool     `js:"Dirty"`
}

func NewWorkSite() *Worksite {
	ws := &Worksite{Object: tools.O()}
	ws.Id = -1
	ws.Client = ""
	ws.Ref = ""
	ws.OrderDate = ""
	ws.DoeDate = ""
	ws.City = ""
	ws.Status = "New"
	ws.Pmz = NewPT()
	ws.Pa = NewPT()
	ws.Comment = ""
	ws.Orders = []*Order{}
	ws.Rework = nil
	ws.Dirty = true

	return ws
}

func WorksiteFromJS(o *js.Object) *Worksite {
	ws := &Worksite{Object: o}
	return ws
}

func (ws *Worksite) HasRework() bool {
	if ws.Rework == nil {
		return false
	}
	if ws.Rework.Object == nil {
		return false
	}
	if ws.Rework.Object == js.Undefined {
		return false
	}
	return true
	//return ws.Rework != nil && ws.Rework.Object != nil && ws.Rework.Object != js.Undefined
}

func (ws *Worksite) Clone() *Worksite {
	nws := &Worksite{Object: tools.O()}
	nws.Copy(ws)
	return nws
}

func (ws *Worksite) Copy(ows *Worksite) {
	ws.Id = ows.Id
	ws.Client = ows.Client
	ws.Ref = ows.Ref
	ws.OrderDate = ows.OrderDate
	ws.DoeDate = ows.DoeDate
	ws.City = ows.City
	ws.Status = ows.Status
	ws.Pmz = ows.Pmz.Clone()
	ws.Pa = ows.Pa.Clone()
	ws.Comment = ows.Comment
	if ows.HasRework() {
		ws.Rework = ows.Rework.Clone()
	}
	ws.Dirty = false // ows.Dirty
	ws.Orders = []*Order{}
	for _, o := range ows.Orders {
		ws.Orders = append(ws.Orders, o.Clone())
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
	res := "Client:" + ws.Client + "\n"
	res += "Ref:" + ws.Ref + "\n"
	res += "OrderDate:" + date.DateString(ws.OrderDate) + "\n"
	res += "DoeDate:" + date.DateString(ws.DoeDate) + "\n"
	res += "City:" + ws.City + "\n"
	res += "Status:" + ws.Status + "\n"
	res += "Pmz:" + ws.Pmz.SearchInString()
	res += "Pa:" + ws.Pa.SearchInString()
	res += "Comment:" + ws.Comment + "\n"

	for _, o := range ws.Orders {
		res += o.SearchInString()
	}
	if ws.HasRework() {
		res += ws.Rework.SearchInString()
	}

	return res
}

func (ws *Worksite) GetInfo() (nbCommand, nbTroncon, nbAvailRacco, nbRacco int) {
	nbCommand = len(ws.Orders)
	for _, o := range ws.Orders {
		nbTroncon += len(o.Troncons)
		for _, t := range o.Troncons {
			nbRacco += t.NbRacco
			if !t.Blockage {
				nbAvailRacco += t.NbRacco
			}
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

func (ws *Worksite) AddOrder() {
	order := NewOrder()
	order.AddTroncon()
	ws.Orders = append(ws.Orders, order)
}

func (ws *Worksite) AddRework() {
	ws.Rework = NewRework()
}

func (ws *Worksite) OrdersCompleted() bool {
	for _, o := range ws.Orders {
		if !o.IsCompleted() {
			return false
		}
	}
	return true
}

func (ws *Worksite) IsDefined() bool {
	return !tools.Empty(ws.Client) &&
		!tools.Empty(ws.City) &&
		!tools.Empty(ws.Ref) &&
		!tools.Empty(ws.OrderDate) &&
		ws.Pmz.IsFilledIn() &&
		ws.Pa.IsFilledIn()
}

func (ws *Worksite) IsFilledIn() bool {
	for _, o := range ws.Orders {
		if !o.IsFilledIn() {
			return false
		}
	}
	return true
}

func WorksiteStatusLabel(value string) string {
	switch value {
	case "New":
		return "Nouveau"
	case "FormInProgress":
		return "Saisie en cours"
	case "InProgress":
		return "Réal. En cours"
	case "DOE":
		return "DOE à faire"
	case "Done":
		return "Terminé"
	case "Rework":
		return "A Reprendre"
	default:
		return "<" + value + ">"
	}
}
