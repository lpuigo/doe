package frontmodel

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

type Worksite struct {
	*js.Object

	Id        int      `js:"Id"`
	Ref       string   `js:"Ref"`
	OrderDate string   `js:"OrderDate"`
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
	ws.Pmz = ows.Pmz.Clone()
	ws.Pa = ows.Pa.Clone()
	ws.Comment = ows.Comment
	ws.Orders = make([]*Order, len(ows.Orders))
	for i, o := range ows.Orders {
		ws.Orders[i] = o.Clone()
	}
}
