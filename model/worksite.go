package model

import (
	"fmt"
	"github.com/lpuig/ewin/doe/website/backend/model/clients"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/model/worksite"
	"strings"
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
	UpdateDate     string
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

func (ws *Worksite) GetInfo(priceByClientArticle func(clientName string, articleName string, qty int) (float64, error)) *fm.WorksiteInfo {
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
	wsi.UpdateDate = ws.UpdateDate
	wsi.City = ws.City
	wsi.Status = ws.Status
	ws.inspectForInfo(wsi, priceByClientArticle)

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

func (ws *Worksite) inspectForInfo(wsi *fm.WorksiteInfo, priceByClientArticle func(clientName string, articleName string, qty int) (float64, error)) {
	wsi.Comment = ws.Comment
	searchPt := func(t string, p PT) string {
		return fmt.Sprintf("%s:%s, PT:%s, Address:%s, ", t, p.Ref, p.RefPt, p.Address)
	}
	wsi.Search = fmt.Sprintf("Comment:%s, ", ws.Comment)
	wsi.Search += fmt.Sprintf("Job:%s, City:%s, ", ws.Ref, ws.City)
	wsi.Search += searchPt("PMZ", ws.Pmz) + searchPt("PA", ws.Pa)
	if wsi.InvoiceName != "" {
		wsi.Search += fmt.Sprintf("Invoice:%s, ", ws.InvoiceName)
	}

	calcInvoice := ws.DoeDate != "" && ws.DoeDate != "null"

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
				if calcInvoice {
					price, err := priceByClientArticle(wsi.Client, t.Article, t.NbRacco)
					if err != nil {
						fmt.Printf("Error evaluating Invoice on %s : %v\n", ws.Ref, err)
						price = 0
					}
					wsi.InvoiceAmount += price
				}
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
	wsi.Search = strings.ToLower(wsi.Search)
}

type StatKey struct {
	Team string
	Date string
	Mes  string
}

type IsWSVisible func(ws *Worksite) bool

// AddStat adds nb of El installed per date (in map[date]nbEl) by visible Client & Client : Teams
func (ws *Worksite) AddStat(nbels map[StatKey]int, dateFor date.DateAggreg, isTeamVisible clients.IsTeamVisible, teamName clients.TeamNameByMember, showTeam, calcToDo bool) {
	nbDOE := 0
	teamDOE := ""

	var addNbEls func(client, team, date, measurement string, nbEl int)

	if showTeam {
		addNbEls = func(client, team, date, measurement string, nbEl int) {
			// add client / team info
			nbels[StatKey{
				Team: client + " : " + "Eq. " + teamName(team),
				Date: date,
				Mes:  measurement,
			}] += nbEl
			// add client info
			nbels[StatKey{
				Team: client,
				Date: date,
				Mes:  measurement,
			}] += nbEl
		}
	} else {
		addNbEls = func(client, team, date, measurement string, nbEl int) {
			// add client info
			nbels[StatKey{
				Team: client,
				Date: date,
				Mes:  measurement,
			}] += nbEl
		}
	}

	for _, o := range ws.Orders {
		for _, t := range o.Troncons {
			if calcToDo {
				nbels[StatKey{
					Team: ws.Client,
					Date: dateFor(ws.OrderDate),
					Mes:  worksite.NbElsSumitted,
				}] += t.NbRacco
			}

			if !isTeamVisible(clients.ClientTeam{Client: ws.Client, Team: t.InstallActor}) {
				continue
			}
			// NbElsInstalled for Team & Client
			if !t.Blockage && t.InstallDate != "" {
				addNbEls(ws.Client, t.InstallActor, dateFor(t.InstallDate), worksite.NbElsInstalled, t.NbRacco)
			}
			// NbElsMeasured for Team & Client
			if !t.Blockage && t.MeasureDate != "" {
				addNbEls(ws.Client, t.MeasureActor, dateFor(t.MeasureDate), worksite.NbElsMeasured, t.NbRacco)

				nbDOE += t.NbRacco
				teamDOE = t.MeasureActor
			}
			// NbElsBlocked
			if t.Blockage {
				d := dateFor(ws.OrderDate)
				if t.InstallDate != "" {
					d = dateFor(t.InstallDate)
				}
				addNbEls(ws.Client, t.InstallActor, d, worksite.NbElsBlocked, t.NbRacco)
			}
		}
	}

	// NbElsDOE
	if ws.DoeDate != "" {
		addNbEls(ws.Client, teamDOE, dateFor(ws.DoeDate), worksite.NbElsDOE, nbDOE)
		if calcToDo && ws.AttachmentDate != "" {
			nbels[StatKey{
				Team: ws.Client,
				Date: dateFor(ws.AttachmentDate),
				Mes:  worksite.NbElsBilled,
			}] += nbDOE
		}
	}
}
