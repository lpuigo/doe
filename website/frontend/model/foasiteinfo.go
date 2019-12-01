package model

import (
	"github.com/gopherjs/gopherjs/js"
	"strings"
)

type FoaSiteInfo struct {
	*js.Object

	Id           int    `js:"Id"`
	Client       string `js:"Client"`
	Ref          string `js:"Ref"`
	Manager      string `js:"Manager"`
	OrderDate    string `js:"OrderDate"`
	UpdateDate   string `js:"UpdateDate"`
	Status       string `js:"Status"`
	Comment      string `js:"Comment"`
	NbFoa        int    `js:"NbFoa"`
	NbFoaBlocked int    `js:"NbFoaBlocked"`
	NbFoaDone    int    `js:"NbFoaDone"`
	Search       string `js:"Search"`
}

func NewBEFoaSiteInfo() *FoaSiteInfo {
	return &FoaSiteInfo{}
}

func NewFoaSiteInfoFromJS(o *js.Object) *FoaSiteInfo {
	return &FoaSiteInfo{Object: o}
}

func (fsi *FoaSiteInfo) TextFiltered(filter string) bool {
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

	return strings.Contains(fsi.Search, strings.ToUpper(filter)) == expected
}
