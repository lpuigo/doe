package model

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/dates"
	"strings"
)

const (
	WsStatusNew            string = "00 New"
	WsStatusFormInProgress string = "10 FormInProgress"
	WsStatusInProgress     string = "20 InProgress"
	WsStatusDOE            string = "30 DOE"
	WsStatusAttachment     string = "40 Attachment"
	WsStatusPayment        string = "50 Payment"
	WsStatusRework         string = "80 Rework"
	WsStatusBlocked        string = "98 Blocked"
	WsStatusDone           string = "99 Done"
)

type Worksite struct {
	*js.Object

	Id             int      `js:"Id"`
	Client         string   `js:"Client"`
	Ref            string   `js:"Ref"`
	OrderDate      string   `js:"OrderDate"`
	DoeDate        string   `js:"DoeDate"`
	AttachmentDate string   `js:"AttachmentDate"`
	PaymentDate    string   `js:"PaymentDate"`
	City           string   `js:"City"`
	Status         string   `js:"Status"`
	Pmz            *PT      `js:"Pmz"`
	Pa             *PT      `js:"Pa"`
	Comment        string   `js:"Comment"`
	Orders         []*Order `js:"Orders"`
	Rework         *Rework  `js:"Rework"`
	Dirty          bool     `js:"Dirty"`
}

func NewWorkSite() *Worksite {
	ws := &Worksite{Object: tools.O()}
	ws.Id = -1
	ws.Client = ""
	ws.Ref = ""
	ws.OrderDate = ""
	ws.DoeDate = ""
	ws.AttachmentDate = ""
	ws.PaymentDate = ""
	ws.City = ""
	ws.Status = WsStatusNew
	ws.Pmz = NewPT()
	ws.Pa = NewPT()
	ws.Comment = ""
	ws.Orders = []*Order{}
	ws.Rework = NewRework()
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
	ws.AttachmentDate = ows.AttachmentDate
	ws.PaymentDate = ows.PaymentDate
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
	res += "AttachmentDate:" + date.DateString(ws.AttachmentDate) + "\n"
	res += "PaymentDate:" + date.DateString(ws.PaymentDate) + "\n"
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

func (ws *Worksite) InstallDates() string {
	min := "9999-99-99"
	max := "0000-00-00"
	for _, o := range ws.Orders {
		for _, t := range o.Troncons {
			if t.InstallDate < min {
				min = t.InstallDate
			}
			if t.InstallDate > max {
				max = t.InstallDate
			}
		}
	}
	res := ""
	if min != "9999-99-99" {
		res += date.DateString(min)
	}
	res += " - "
	if max != "0000-00-00" {
		res += date.DateString(max)
	}
	return res
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

func (ws *Worksite) IsBlocked() bool {
	_, _, nbAivailEl, _ := ws.GetInfo()
	return nbAivailEl == 0
}

func (ws *Worksite) NeedRework() bool {
	if ws.Rework == nil {
		return false
	}
	return ws.Rework.NeedRework()
}

func (ws *Worksite) WorksiteStatusLabel() string {
	return WorksiteStatusLabel(ws.Status)
}

func WorksiteStatusLabel(status string) string {
	switch status {
	case WsStatusNew:
		return "Nouveau"
	case WsStatusFormInProgress:
		return "Saisie en cours"
	case WsStatusInProgress:
		return "Réal. En cours"
	case WsStatusDOE:
		return "DOE à faire"
	case WsStatusAttachment:
		return "Attachement attendu"
	case WsStatusPayment:
		return "Paiement attendu"
	case WsStatusRework:
		return "Reprise à faire"
	case WsStatusBlocked:
		return "Bloqué"
	case WsStatusDone:
		return "Terminé"
	default:
		return "<" + status + ">"
	}
}

func (ws *Worksite) UpdateStatus() {
	if !ws.IsDefined() {
		ws.Status = WsStatusNew
		return
	}
	if !ws.IsFilledIn() {
		ws.Status = WsStatusFormInProgress
		return
	}
	if !ws.OrdersCompleted() {
		ws.DoeDate = ""
		ws.Status = WsStatusInProgress
		return
	}
	if ws.IsBlocked() {
		ws.Status = WsStatusBlocked
		return
	}
	if tools.Empty(ws.DoeDate) {
		ws.Status = WsStatusDOE
		return
	}
	if ws.NeedRework() {
		ws.Status = WsStatusRework
		return
	}
	if tools.Empty(ws.AttachmentDate) {
		ws.Status = WsStatusAttachment
		return
	}
	if tools.Empty(ws.PaymentDate) {
		ws.Status = WsStatusPayment
		return
	}
	ws.Status = WsStatusDone
}

func (ws *Worksite) GetPtByName(refpt string) *Troncon {
	for _, o := range ws.Orders {
		for _, tr := range o.Troncons {
			if tr.Pb.RefPt == refpt {
				return tr
			}
		}
	}
	return nil
}

func WorksiteIsUpdatable(value string) bool {
	switch value {
	//case WsStatusNew:
	//	return true
	//case WsStatusFormInProgress:
	//	return true
	case WsStatusInProgress:
		return true
	case WsStatusDOE:
		return true
	//case WsStatusAttachment:
	//	return true
	//case WsStatusPayment:
	//	return true
	case WsStatusRework:
		return true
	case WsStatusBlocked:
		return true
		//case WsStatusDone:
		//	return true
	}
	return false
}

func WorksiteMustRework(value string) bool {
	switch value {
	//case WsStatusNew:
	//	return true
	//case WsStatusFormInProgress:
	//	return true
	//case WsStatusInProgress:
	//	return true
	//case WsStatusDOE:
	//	return true
	//case WsStatusAttachment:
	//	return true
	//case WsStatusPayment:
	//	return true
	case WsStatusRework:
		return true
		//case WsStatusBlocked:
		//	return true
		//case WsStatusDone:
		//	return true
	}
	return false
}

func WorksiteIsReworkable(value string) bool {
	switch value {
	//case WsStatusNew:
	//	return true
	//case WsStatusFormInProgress:
	//	return true
	//case WsStatusInProgress:
	//	return true
	//case WsStatusDOE:
	//	return true
	case WsStatusAttachment:
		return true
	case WsStatusPayment:
		return true
	case WsStatusRework:
		return true
	//case WsStatusBlocked:
	//	return true
	case WsStatusDone:
		return true
	}
	return false
}
