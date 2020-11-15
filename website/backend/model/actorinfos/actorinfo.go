package actorinfos

import "github.com/lpuig/ewin/doe/website/backend/model/date"

type DateComment struct {
	Date    string
	Comment string
}

type DateAmountComment struct {
	DateComment
	Amount float64
}

// Events type is dedicated for date-defined event. key is the event name (ex : ProBTP, Caces Nacelle, ...)
type Events map[string]DateComment

// EarningHistory type is dedicated to date-changing earning (like salary) : only most recent entry is applicable. Other are stored for history purpose
type EarningHistory []DateAmountComment

// GetApplicableByDate returns applicable DateAmountComment for given date. Found boolean must be checked (true if applicable info exists, false otherwise)
func (eh EarningHistory) GetApplicableByDate(d string) (dac DateAmountComment, found bool) {
	for _, dac := range eh {
		if d >= dac.Date {
			return dac, true
		}
	}
	return DateAmountComment{}, false
}

// Earnings type is dedicated to date-defined earning (like monthly bonuses)
type Earnings []DateAmountComment

// GetByDate returns existing DateAmountComment for given date. Found boolean must be checked (true if info exists, false otherwise)
func (e Earnings) GetByDate(d string) (dac DateAmountComment, found bool) {
	lookupDate := date.GetMonth(d)
	for _, dac := range e {
		if lookupDate == date.GetMonth(dac.Date) {
			return dac, true
		}
	}
	return DateAmountComment{}, false
}

type ActorInfo struct {
	Id            int
	ActorId       int
	Trainings     Events
	Salary        EarningHistory
	EarnedBonuses Earnings
	PaidBonuses   Earnings
	TravelSubsidy EarningHistory
}

func NewActorInfo() *ActorInfo {
	return &ActorInfo{
		Id:            -1,
		ActorId:       -1,
		Trainings:     make(Events),
		Salary:        make(EarningHistory, 0),
		EarnedBonuses: make(Earnings, 0),
		PaidBonuses:   make(Earnings, 0),
		TravelSubsidy: make(EarningHistory, 0),
	}
}
