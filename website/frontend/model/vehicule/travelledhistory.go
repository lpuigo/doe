package vehicule

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/date"
)

// Type TravelledHistory reflects ewin/doe/website/backend/model/vehicules.ActorHistory
type TravelledHistory struct {
	*js.Object

	Date string `js:"Date"`
	Kms  int    `js:"Kms"`
}

func NewTravelledHistory() *TravelledHistory {
	nah := &TravelledHistory{Object: tools.O()}
	nah.Date = date.TodayAfter(0)
	nah.Kms = -1
	return nah
}

func (ah TravelledHistory) Copy() *TravelledHistory {
	nah := NewTravelledHistory()
	nah.Date = ah.Date
	nah.Kms = ah.Kms
	return nah
}

func CompareTravelledHistory(a, b TravelledHistory) int {
	if a.Date > b.Date {
		return -1
	}
	if a.Date == b.Date {
		return 0
	}
	return 1
}
