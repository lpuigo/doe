package model

import (
	"fmt"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
)

type Worksite struct {
	Id             int
	Client         string
	Ref            string
	OrderDate      string
	DoeDate        string
	AttachmentDate string
	PaymentDate    string
	City           string
	Status         string
	Pmz            PT
	Pa             PT
	Comment        string
	Orders         []Order
	Rework         *Rework
}

func MakeWorksite(ref, orderdate string, pmz, pa PT, order ...Order) Worksite {
	return Worksite{Ref: ref, Pmz: pmz, Pa: pa, Orders: order, OrderDate: orderdate}
}

func (w Worksite) FileName() string {
	return w.OrderDate + "_" + w.Ref
}

func (ws *Worksite) GetInfo() *fm.WorksiteInfo {
	wsi := fm.NewBEWorksiteInfo()
	wsi.Id = ws.Id

	wsi.Id = ws.Id
	wsi.Client = ws.Client
	wsi.Ref = ws.Ref
	wsi.OrderDate = ws.OrderDate
	wsi.DoeDate = ws.DoeDate
	wsi.AttachmentDate = ws.AttachmentDate
	wsi.PaymentDate = ws.PaymentDate
	wsi.City = ws.City
	wsi.Status = ws.Status
	ws.inspectForInfo(wsi)

	wsi.Inspected = false
	wsi.NbRework = 0
	wsi.NbReworkDone = 0
	return wsi
}

func (ws *Worksite) inspectForInfo(wsi *fm.WorksiteInfo) {
	wsi.Comment = ws.Comment
	searchPt := func(t string, p PT) string {
		return fmt.Sprintf("%s:%s, PT:%s, Address:%s, ", t, p.Ref, p.RefPt, p.Address)
	}
	wsi.Search = searchPt("PMZ", ws.Pmz) + searchPt("PA", ws.Pa)
	for _, o := range ws.Orders {
		wsi.NbOrder += 1
		wsi.Search += fmt.Sprintf("Order:%s, ", o.Ref)
		if o.Comment != "" {
			wsi.Comment += fmt.Sprintf("\n%s: %s", o.Ref, o.Comment)
		}
		for _, t := range o.Troncons {
			wsi.NbTroncon += 1
			wsi.NbElTotal += t.NbRacco
			if t.Blockage {
				wsi.NbElBlocked += t.NbRacco
			}
			if t.InstallDate != "" && t.InstallDate != "null" {
				wsi.NbElInstalled += t.NbRacco
			}
			if t.MeasureDate != "" && t.MeasureDate != "null" {
				wsi.NbElMeasured += t.NbRacco
			}
			if t.Comment != "" {
				wsi.Comment += fmt.Sprintf("\n%s (%s): %s", t.Pb.Ref, t.Pb.RefPt, t.Comment)
			}
			wsi.Search += searchPt("PB", t.Pb)
		}
	}
}
