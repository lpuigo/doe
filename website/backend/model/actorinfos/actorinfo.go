package actorinfos

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

// Earnings type is dedicated to date-defined earning (like monthly bonuses)
type Earnings []DateAmountComment

type ActorInfo struct {
	Id            int
	ActorId       int
	Trainings     Events
	Salary        EarningHistory
	EarnedBonuses Earnings
	PaidBonuses   Earnings
	TravelSubsidy Earnings
}

func NewActorInfo() *ActorInfo {
	return &ActorInfo{
		Id:            -1,
		ActorId:       -1,
		Trainings:     make(Events),
		Salary:        make(EarningHistory, 0),
		EarnedBonuses: make(Earnings, 0),
		PaidBonuses:   make(Earnings, 0),
		TravelSubsidy: make(Earnings, 0),
	}
}
