package actor

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

// Type Actor reflects ewin/doe/website/backend/model/actorinfos.ActorInfo
type ActorInfo struct {
	*js.Object

	Id            int            `js:"Id"`
	ActorId       int            `js:"ActorId"`
	Trainings     Events         `js:"Trainings"`
	Salary        EarningHistory `js:"Salary"`
	EarnedBonuses Earnings       `js:"EarnedBonuses"`
	PaidBonuses   Earnings       `js:"PaidBonuses"`
	Bonuses       PaymentHistory `js:"Bonuses"`
	TravelSubsidy EarningHistory `js:"TravelSubsidy"`
}

func NewActorInfo() *ActorInfo {
	nai := &ActorInfo{Object: tools.O()}
	nai.Id = -1
	nai.ActorId = -1
	nai.Trainings = Events{}
	nai.Salary = EarningHistory{}
	nai.EarnedBonuses = Earnings{}
	nai.PaidBonuses = Earnings{}
	nai.Bonuses = PaymentHistory{}
	nai.TravelSubsidy = EarningHistory{}
	return nai
}

func NewActorInfoForActor(actor *Actor) *ActorInfo {
	nai := &ActorInfo{Object: tools.O()}
	nai.Id = -1
	nai.ActorId = actor.Id
	nai.Trainings = Events{}
	nai.Salary = EarningHistory{}
	nai.EarnedBonuses = Earnings{}
	nai.PaidBonuses = Earnings{}
	nai.Bonuses = PaymentHistory{}
	nai.TravelSubsidy = EarningHistory{}
	return nai
}

func ActorInfoFromJS(obj *js.Object) *ActorInfo {
	return &ActorInfo{Object: obj}
}

func (ai *ActorInfo) Clone(oai *ActorInfo) {
	ai.Id = oai.Id
	ai.ActorId = oai.ActorId
	ai.Trainings = oai.Trainings.Copy()
	ai.Salary = oai.Salary[:]
	ai.EarnedBonuses = oai.EarnedBonuses[:]
	ai.PaidBonuses = oai.PaidBonuses[:]
	ai.Bonuses = oai.Bonuses[:]
	ai.TravelSubsidy = oai.TravelSubsidy[:]
}
