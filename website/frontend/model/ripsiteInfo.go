package model

import (
	"github.com/gopherjs/gopherjs/js"
	"strings"
)

type RipsiteInfo struct {
	*js.Object

	Id                   int    `js:"Id"`
	Client               string `js:"Client"`
	Ref                  string `js:"Ref"`
	Manager              string `js:"Manager"`
	OrderDate            string `js:"OrderDate"`
	UpdateDate           string `js:"UpdateDate"`
	Status               string `js:"Status"`
	Comment              string `js:"Comment"`
	NbPulling            int    `js:"NbPulling"`
	NbPullingBlocked     int    `js:"NbPullingBlocked"`
	NbPullingDone        int    `js:"NbPullingDone"`
	NbJunction           int    `js:"NbJunction"`
	NbJunctionBlocked    int    `js:"NbJunctionBlocked"`
	NbJunctionDone       int    `js:"NbJunctionDone"`
	NbMeasurement        int    `js:"NbMeasurement"`
	NbMeasurementBlocked int    `js:"NbMeasurementBlocked"`
	NbMeasurementDone    int    `js:"NbMeasurementDone"`
	Search               string `js:"Search"`
}

func NewBERipsiteInfo() *RipsiteInfo {
	return &RipsiteInfo{}
}

func NewRipsiteInfoFromJS(o *js.Object) *RipsiteInfo {
	return &RipsiteInfo{Object: o}
}

func (rsi *RipsiteInfo) TextFiltered(filter string) bool {
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
	return strings.Contains(rsi.Search, filter) == expected
}
