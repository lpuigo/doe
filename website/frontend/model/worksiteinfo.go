package model

import (
	"github.com/gopherjs/gopherjs/js"
	"strings"
)

type WorksiteInfo struct {
	*js.Object

	Id             int     `js:"Id"`
	Client         string  `js:"Client"`
	Ref            string  `js:"Ref"`
	OrderDate      string  `js:"OrderDate"`
	DoeDate        string  `js:"DoeDate"`
	AttachmentDate string  `js:"AttachmentDate"`
	InvoiceDate    string  `js:"InvoiceDate"`
	InvoiceName    string  `js:"InvoiceName"`
	PaymentDate    string  `js:"PaymentDate"`
	City           string  `js:"City"`
	Status         string  `js:"Status"`
	Comment        string  `js:"Comment"`
	NbOrder        int     `js:"NbOrder"`
	NbTroncon      int     `js:"NbTroncon"`
	NbElTotal      int     `js:"NbElTotal"`
	NbElBlocked    int     `js:"NbElBlocked"`
	NbElInstalled  int     `js:"NbElInstalled"`
	NbElMeasured   int     `js:"NbElMeasured"`
	InvoiceAmount  float64 `js:"InvoiceAmount"`
	Inspected      bool    `js:"Inspected"`
	NbRework       int     `js:"NbRework"`
	NbReworkDone   int     `js:"NbReworkDone"`
	Search         string  `js:"Search"`
}

func NewBEWorksiteInfo() *WorksiteInfo {
	return &WorksiteInfo{}
}

func NewWorksiteInfoFromJs(o *js.Object) *WorksiteInfo {
	return &WorksiteInfo{Object: o}
}

func (wsi *WorksiteInfo) TextFiltered(filter string) bool {
	expected := true
	if filter == "" {
		return true
	}
	if strings.HasPrefix(filter, `\`) {
		if len(filter) == 1 {
			return true
		}
		expected = false
		filter = filter[1:]
	}
	return strings.Contains(wsi.Search, filter) == expected
}

func (wsi *WorksiteInfo) NeedRework() bool {
	if wsi.NbRework > 0 {
		if wsi.NbReworkDone == wsi.NbRework {
			return false
		}
		return true
	}
	return false
}
