package actorinfo

import "github.com/gopherjs/gopherjs/js"

// Type Actor reflects ewin/doe/website/backend/model/actorinfos.ActorInfo
type ActorInfo struct {
	*js.Object

	Id            int            `js:"Id"`
	ActorId       int            `js:"ActorId"`
	Trainings     Events         `js:"Trainings"`
	Salary        EarningHistory `js:"Salary"`
	EarnedBonuses Earnings       `js:"EarnedBonuses"`
	PaidBonuses   Earnings       `js:"PaidBonuses"`
	TravelSubsidy Earnings       `js:"TravelSubsidy"`
}

func NewActorInfoFromJS(obj *js.Object) *ActorInfo {
	return &ActorInfo{Object: obj}
}
