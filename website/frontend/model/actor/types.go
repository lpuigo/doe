package actor

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

// Type DateComment reflects ewin/doe/website/backend/model/actorinfos.DateComment
type DateComment struct {
	*js.Object

	Date    string `js:"Date"`
	Comment string `js:"Comment"`
}

func NewDateComment() *DateComment {
	nd := &DateComment{Object: tools.O()}
	nd.Date = ""
	nd.Comment = ""
	return nd
}

// Type DateAmountComment reflects ewin/doe/website/backend/model/actorinfos.DateAmountComment
type DateAmountComment struct {
	DateComment
	Amount float64 `js:"Amount"`
}

func NewDateAmountComment() *DateAmountComment {
	nd := &DateAmountComment{DateComment: *NewDateComment()}
	nd.Amount = 0
	return nd
}

// Type Events reflects ewin/doe/website/backend/model/actorinfos.Events
// Events type is dedicated for date-defined event. key is the event name (ex : ProBTP, Caces Nacelle, ...)
type Events map[string]DateComment

func (e Events) Copy() Events {
	ne := Events{}
	for s, comment := range e {
		ne[s] = comment
	}
	return ne
}

// Type EarningHistory reflects ewin/doe/website/backend/model/actorinfos.EarningHistory
// EarningHistory type is dedicated to date-changing earning (like salary) : only most recent entry is applicable. Other are stored for history purpose
type EarningHistory []DateAmountComment

// Type Earnings reflects ewin/doe/website/backend/model/actorinfos.Earnings
// Earnings type is dedicated to date-defined earning (like monthly bonuses)
type Earnings []DateAmountComment
