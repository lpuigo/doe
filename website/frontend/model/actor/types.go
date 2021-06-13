package actor

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/date"
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

func DateAmountCommentFromJs(o *js.Object) *DateAmountComment {
	return &DateAmountComment{DateComment: DateComment{Object: o}}
}

func CompareDateAmountComment(a, b DateAmountComment) int {
	if a.Date > b.Date {
		return -1
	}
	if a.Date == b.Date {
		return 0
	}
	return 1
}

func (dac *DateAmountComment) Copy() *DateAmountComment {
	res := NewDateAmountComment()
	res.Date = dac.Date
	res.Comment = dac.Comment
	res.Amount = dac.Amount
	return res
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

// Type Earnings reflects ewin/doe/website/backend/model/actorinfos.Earnings
// Earnings type is dedicated to date-defined earning (like monthly bonuses)
type Earnings []DateAmountComment

// Type EarningHistory reflects ewin/doe/website/backend/model/actorinfos.EarningHistory
// EarningHistory type is dedicated to date-changing earning (like salary) : only most recent entry is applicable. Other are stored for history purpose
type EarningHistory []DateAmountComment

func (eh EarningHistory) CurrentDateAmountComment() *DateAmountComment {
	if len(eh) == 0 {
		return nil
	}
	var currentDac DateAmountComment
	// search for current applicable entry
	switch len(eh) {
	case 1:
		currentDac = eh[0]
	default:
		today := date.TodayAfter(0)
		for _, dac := range eh {
			if dac.Date > today {
				continue
			}
			currentDac = dac
			break
		}
	}
	return &currentDac
}

// Type PaymentHistory reflects ewin/doe/website/backend/model/actorinfos.PaymentHistory
// PaymentHistory type is dedicated to date-defined earning and payment (like monthly bonuses)
type PaymentHistory []Payment

// Type Payment reflects ewin/doe/website/backend/model/actorinfos.Payment
// Payment type is dedicated to date-defined earning and payment (like monthly bonuses)
type Payment struct {
	*js.Object
	Earning *DateAmountComment `js:"Earning"`
	Payment EarningHistory     `js:"Payment"`
}

func NewPayment() *Payment {
	p := &Payment{Object: tools.O()}
	p.Earning = NewDateAmountComment()
	p.Payment = make(EarningHistory, 0)
	return p
}

func ComparePayment(a, b Payment) int {
	if a.Earning.Date > b.Earning.Date {
		return -1
	}
	if a.Earning.Date == b.Earning.Date {
		return 0
	}
	return 1
}

func (p *Payment) IsPaid() bool {
	if len(p.Payment) == 0 {
		return false
	}
	var amount float64
	for _, pm := range p.Payment {
		amount += pm.Amount
	}
	return amount >= p.Earning.Amount
}
