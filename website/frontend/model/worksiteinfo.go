package model

import "github.com/gopherjs/gopherjs/js"

type WorksiteInfo struct {
	*js.Object

	Id             int    `js:"Id"`
	Client         string `js:"Client"`
	Ref            string `js:"Ref"`
	OrderDate      string `js:"OrderDate"`
	DoeDate        string `js:"DoeDate"`
	AttachmentDate string `js:"AttachmentDate"`
	PaymentDate    string `js:"PaymentDate"`
	City           string `js:"City"`
	Status         string `js:"Status"`
	Comment        string `js:"Comment"`
	NbOrder        int    `js:"NbOrder"`
	NbTroncon      int    `js:"NbTroncon"`
	NbElTotal      int    `js:"NbElTotal"`
	NbElBlocked    int    `js:"NbElBlocked"`
	NbElInstalled  int    `js:"NbElInstalled"`
	NbElMeasured   int    `js:"NbElMeasured"`
	Inspected      bool   `js:"Inspected"`
	NbRework       int    `js:"NbRework"`
	NbReworkDone   int    `js:"NbReworkDone"`
	Search         string `js:"Search"`
}

func NewBEWorksiteInfo() *WorksiteInfo {
	return &WorksiteInfo{}
}

func NewWorksiteInfoFromJs(o *js.Object) *WorksiteInfo {
	return &WorksiteInfo{Object: o}
}
