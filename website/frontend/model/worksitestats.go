package model

import "github.com/gopherjs/gopherjs/js"

type WorksiteStats struct {
	*js.Object

	StartDate string   `js:"StartDate"`
	Teams     []string `js:"Teams"`
	NbEls     [][]int  `js:"NbEls"`
}

func NewBEWorksiteStats() *WorksiteStats {
	return &WorksiteStats{}
}

func NewWorksiteStatsFromJs(o *js.Object) *WorksiteStats {
	return &WorksiteStats{Object: o}
}
