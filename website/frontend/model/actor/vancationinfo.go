package actor

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/date"
)

// Type LeavePeriod reflects ewin/doe/website/backend/model/actors.LeavePeriod
type LeavePeriod struct {
	date.DateRange
	Type    string `js:"Type"`
	Comment string `js:"Comment"`
}

func NewLeavePeriod() *LeavePeriod {
	lp := &LeavePeriod{DateRange: *date.NewDateRange()}
	lp.Type = ""
	lp.Comment = ""
	return lp
}

// Type VacationInfo reflects ewin/doe/website/backend/model/actors.VacationInfo
type VacationInfo struct {
	*js.Object

	Vacation      []*LeavePeriod `js:"Vacation"`
	EarnedDays    float64        `js:"EarnedDays"`
	AvailableDays float64        `js:"AvailableDays"`
	TakenDays     float64        `js:"TakenDays"`
}

func NewVacationInfo() *VacationInfo {
	vi := &VacationInfo{Object: tools.O()}
	vi.Vacation = []*LeavePeriod{}
	vi.EarnedDays = 0
	vi.AvailableDays = 0
	vi.TakenDays = 0
	return vi
}
