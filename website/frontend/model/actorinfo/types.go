package actorinfo

import "github.com/gopherjs/gopherjs/js"

// Type DateComment reflects ewin/doe/website/backend/model/actorinfos.DateComment
type DateComment struct {
	*js.Object

	Date    string
	Comment string
}

// Type DateAmountComment reflects ewin/doe/website/backend/model/actorinfos.DateAmountComment
type DateAmountComment struct {
	DateComment
	Amount float64
}

// Type Events reflects ewin/doe/website/backend/model/actorinfos.Events
// Events type is dedicated for date-defined event. key is the event name (ex : ProBTP, Caces Nacelle, ...)
type Events map[string]DateComment

// Type EarningHistory reflects ewin/doe/website/backend/model/actorinfos.EarningHistory
// EarningHistory type is dedicated to date-changing earning (like salary) : only most recent entry is applicable. Other are stored for history purpose
type EarningHistory []DateAmountComment

// Type Earnings reflects ewin/doe/website/backend/model/actorinfos.Earnings
// Earnings type is dedicated to date-defined earning (like monthly bonuses)
type Earnings []DateAmountComment
