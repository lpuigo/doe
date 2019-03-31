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
	InvoiceDate    string
	InvoiceName    string
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
	wsi.Client = ws.Client
	wsi.Ref = ws.Ref
	wsi.OrderDate = ws.OrderDate
	wsi.DoeDate = ws.DoeDate
	wsi.AttachmentDate = ws.AttachmentDate
	wsi.InvoiceDate = ws.InvoiceDate
	wsi.InvoiceName = ws.InvoiceName
	wsi.PaymentDate = ws.PaymentDate
	wsi.City = ws.City
	wsi.Status = ws.Status
	ws.inspectForInfo(wsi)

	if ws.Rework != nil && len(ws.Rework.Defects) > 0 {
		wsi.Inspected = true
		for _, d := range ws.Rework.Defects {
			if d.ToBeFixed {
				wsi.NbRework += 1
				if d.FixDate != "" {
					wsi.NbReworkDone += 1
				}
			}
		}
	}
	return wsi
}

// setInvoiceAmount sets Invoice Amount in given WorksiteInfo (only if DoeDate is set)
func (ws *Worksite) setInvoiceAmount(wsi *fm.WorksiteInfo) {
	if ws.DoeDate == "" || ws.DoeDate == "null" {
		return
	}
	// fast exit if whole worksite blocked
	if wsi.NbElBlocked == wsi.NbElTotal {
		return
	}
	// TODO To be updated for different CEM amount
	const CEM42Amount float64 = 70
	wsi.InvoiceAmount = float64(wsi.NbElMeasured) * CEM42Amount
}

func (ws *Worksite) inspectForInfo(wsi *fm.WorksiteInfo) {
	wsi.Comment = ws.Comment
	searchPt := func(t string, p PT) string {
		return fmt.Sprintf("%s:%s, PT:%s, Address:%s, ", t, p.Ref, p.RefPt, p.Address)
	}
	wsi.Search = fmt.Sprintf("Chantier:%s, Ville:%s, ", ws.Ref, ws.City)
	wsi.Search += searchPt("PMZ", ws.Pmz) + searchPt("PA", ws.Pa)
	for _, o := range ws.Orders {
		lf := "\n"
		wsi.NbOrder += 1
		wsi.Search += fmt.Sprintf("Order:%s, ", o.Ref)
		if o.Comment != "" {
			if wsi.Comment == "" {
				lf = ""
			}
			wsi.Comment += fmt.Sprintf("%s%s: %s", lf, o.Ref, o.Comment)
		}
		for _, t := range o.Troncons {
			lf := "\n"
			wsi.NbTroncon += 1
			wsi.NbElTotal += t.NbRacco
			if t.Blockage {
				wsi.NbElBlocked += t.NbRacco
			}
			if !t.Blockage && t.InstallDate != "" && t.InstallDate != "null" {
				wsi.NbElInstalled += t.NbRacco
			}
			if !t.Blockage && t.MeasureDate != "" && t.MeasureDate != "null" {
				wsi.NbElMeasured += t.NbRacco
			}
			if t.Comment != "" {
				if wsi.Comment == "" {
					lf = ""
				}
				wsi.Comment += fmt.Sprintf("%s%s (%s): %s", lf, t.Pb.Ref, t.Pb.RefPt, t.Comment)
			}
			wsi.Search += searchPt("PB", t.Pb)
		}
	}
	ws.setInvoiceAmount(wsi)
}

type StatKey struct {
	Team string
	Date string
	Mes  string
}

type ClientTeam struct {
	Client string
	Team   string
}

type IsWSVisible func(ws *Worksite) bool
type IsTeamVisible func(ClientTeam) bool
type DateAggreg func(string) string

const (
	NbElsInstalled string = "Installed"
	NbElsBlocked   string = "Blocked"
	NbElsMeasured  string = "Measured"
	NbElsDOE       string = "DOE"
)

// AddStat adds nb of El installed per date (in map[date]nbEl) by visible Teams
func (ws *Worksite) AddStat(nbels map[StatKey]int, dateFor DateAggreg, isTeamVisible IsTeamVisible) {
	nbDOE := 0
	teamDOE := ""
	for _, o := range ws.Orders {
		for _, t := range o.Troncons {
			if !isTeamVisible(ClientTeam{Client: ws.Client, Team: t.InstallActor}) {
				continue
			}
			// NbElsInstalled
			if !t.Blockage && t.InstallDate != "" {
				key := StatKey{
					Team: t.InstallActor,
					Date: dateFor(t.InstallDate),
					Mes:  NbElsInstalled,
				}
				nbels[key] += t.NbRacco
			}
			// NbElsMeasured
			if !t.Blockage && t.MeasureDate != "" {
				key := StatKey{
					Team: t.MeasureActor,
					Date: dateFor(t.MeasureDate),
					Mes:  NbElsMeasured,
				}
				nbels[key] += t.NbRacco
				nbDOE += t.NbRacco
				teamDOE = t.MeasureActor
			}
			// NbElsBlocked
			if t.Blockage {
				d := dateFor(ws.OrderDate)
				if t.InstallDate != "" {
					d = dateFor(t.InstallDate)
				}
				key := StatKey{
					Team: t.MeasureActor,
					Date: d,
					Mes:  NbElsBlocked,
				}
				nbels[key] += t.NbRacco
			}
		}
	}

	// NbElsDOE
	if ws.DoeDate != "" {
		key := StatKey{
			Team: teamDOE,
			Date: dateFor(ws.DoeDate),
			Mes:  NbElsDOE,
		}
		nbels[key] += nbDOE
	}
}
