package actors

import (
	"github.com/lpuig/ewin/doe/website/frontend/model/actor/actorconst"
	"testing"
)

func TestMigrate_Actors(t *testing.T) {
	dir := `C:\Users\Laurent\Golang\src\github.com\lpuig\ewin\doe\website\backend\model\actors\test`
	ap, err := NewActorsPersister(dir)
	if err != nil {
		t.Fatal("NewActorsPersister returned unexpected: ", err.Error())
	}
	ap.persister.SetPersistDelay(0)
	err = ap.LoadDirectory()
	if err != nil {
		t.Fatal("LoadDirectory returned unexpected: ", err.Error())
	}
	for _, ar := range ap.actors {
		ar.PopulateVacationInfo()
		ap.persister.MarkDirty(ar)
	}
	//time.Sleep(2*time.Second)
}

func (a *Actor) PopulateVacationInfo() {
	a.VacInfo.Vacation = make([]LeavePeriod, len(a.Vacation))
	var t string
	if a.Contract != actorconst.ContractTemp {
		t = LeaveTypePaid
	} else {
		t = LeaveTypeUnpaid
	}
	for i, v := range a.Vacation {
		a.VacInfo.Vacation[i] = LeavePeriod{
			DateStringRange: v.DateStringRange,
			Type:            t,
			Comment:         v.Comment,
		}
	}
}
