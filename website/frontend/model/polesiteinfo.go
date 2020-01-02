package model

import (
	"github.com/gopherjs/gopherjs/js"
	"strings"
)

type PolesiteInfo struct {
	*js.Object

	Id            int    `js:"Id"`
	Client        string `js:"Client"`
	Ref           string `js:"Ref"`
	Manager       string `js:"Manager"`
	OrderDate     string `js:"OrderDate"`
	UpdateDate    string `js:"UpdateDate"`
	Status        string `js:"Status"`
	Comment       string `js:"Comment"`
	NbPole        int    `js:"NbPole"`
	NbPoleBlocked int    `js:"NbPoleBlocked"`
	NbPoleDone    int    `js:"NbPoleDone"`
	NbPoleBilled  int    `js:"NbPoleBilled"`
	Search        string `js:"Search"`
}

func NewBEPolesiteInfo() *PolesiteInfo {
	return &PolesiteInfo{}
}

func NewPolesiteInfoFromJS(o *js.Object) *PolesiteInfo {
	return &PolesiteInfo{Object: o}
}

func (psi *PolesiteInfo) TextFiltered(filter string) bool {
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

	return strings.Contains(psi.Search, strings.ToUpper(filter)) == expected
}
