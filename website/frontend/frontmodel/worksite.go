package frontmodel

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
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
	res := ""
	//res += "Id:" Skipped on purpose
	res += "Ref:" + ws.Ref + "\n"
	res += "OrderDate:" + ws.Ref + "\n"
	res += "City:" + ws.City + "\n"
	res += "Status:" + ws.Status + "\n"
	res += "Pmz:" + ws.Ref + "\n"
	res += "Pa:" + ws.Ref + "\n"
	res += "Comment:" + ws.Ref + "\n"
	res += "Orders:" + ws.Ref + "\n"

	for _, v := range ws.Orders {
		res += v.SearchInString()
	}
	return res
}
