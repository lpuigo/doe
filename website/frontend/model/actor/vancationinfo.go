package actor

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/date"
	"github.com/lpuig/ewin/doe/website/frontend/tools/json"
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

func NewLeavePeriodFrom(beg, end, typ, cmt string) *LeavePeriod {
	lp := NewLeavePeriod()
	lp.Begin = beg
	lp.End = end
	lp.Type = typ
	lp.Comment = cmt
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

func VacationInfoFromJS(obj *js.Object) *VacationInfo {
	return &VacationInfo{Object: obj}
}

func (vi *VacationInfo) Copy() *VacationInfo {
	return VacationInfoFromJS(json.Parse(json.Stringify(vi.Object)))
}

func (vi *VacationInfo) GetNextVacation() *date.DateRange {
	if len(vi.Vacation) == 0 {
		return nil
	}
	today := date.TodayAfter(0)
	vacBegin := ""
	vacEnd := ""
	for _, vacPeriod := range vi.Vacation {
		if vacPeriod.End < today {
			continue
		}
		if vacBegin == "" && vacPeriod.End >= today {
			vacBegin = vacPeriod.Begin
			vacEnd = vacPeriod.End
			continue
		}
		// vacBegin != ""
		if vacPeriod.Begin < vacBegin {
			vacBegin = vacPeriod.Begin
			vacEnd = vacPeriod.End
		}
	}

	if vacBegin == "" {
		return nil
	}
	vdr := date.NewDateRange()
	vdr.Begin = vacBegin
	vdr.End = vacEnd
	return vdr
}
