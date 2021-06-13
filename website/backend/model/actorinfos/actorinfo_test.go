package actorinfos

import (
	"fmt"
	"testing"
)

func TestNewActorInfosPersister(t *testing.T) {
	dir := `C:\Users\Laurent\Golang\src\github.com\lpuig\ewin\doe\website\backend\model\actorinfos\test`
	aip, err := NewActorInfosPersister(dir)
	if err != nil {
		t.Fatal("NewActorInfosPersister returned unexpected: ", err.Error())
	}
	aip.persister.SetPersistDelay(0)
	err = aip.LoadDirectory()
	if err != nil {
		t.Fatal("LoadDirectory returned unexpected: ", err.Error())
	}
	for _, actorInfoRecord := range aip.actorInfosById {
		for _, s := range actorInfoRecord.PopulateBonuses() {
			t.Log(s)
		}
		aip.persister.MarkDirty(actorInfoRecord)
	}
	//time.Sleep(2*time.Second)
}

func (ai *ActorInfo) PopulateBonuses() (res []string) {
	if len(ai.EarnedBonuses) == 0 {
		return
	}
	res = append(res, fmt.Sprintf("ActorInfo %d", ai.Id))
	for _, bonus := range ai.EarnedBonuses {
		amount := bonus.Amount
		res = append(res, fmt.Sprintf("\tearned %s: %.2f (%s)", bonus.Date, bonus.Amount, bonus.Comment))
		ph := Payment{
			Earning: DateAmountComment{
				DateComment: DateComment{
					Date:    bonus.Date,
					Comment: bonus.Comment,
				},
				Amount: bonus.Amount,
			},
			Payment: make(EarningHistory, 0),
		}
		for _, paidBonus := range ai.PaidBonuses {
			if !(paidBonus.Date == bonus.Date && paidBonus.Comment == bonus.Comment) {
				continue
			}
			ph.Payment = append(ph.Payment, DateAmountComment{
				DateComment: DateComment{
					Date:    paidBonus.Date,
					Comment: paidBonus.Comment,
				},
				Amount: paidBonus.Amount},
			)
			amount -= paidBonus.Amount
			res = append(res, fmt.Sprintf("\t\tpayment %.2f (%s)", paidBonus.Amount, paidBonus.Comment))
		}
		if amount != 0.0 {
			res = append(res, fmt.Sprintf("\tremaining %.2f", amount))
		}
		ai.Bonuses = append(ai.Bonuses, ph)
	}
	return
}
