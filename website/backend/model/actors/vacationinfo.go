package actors

import "github.com/lpuig/ewin/doe/website/backend/model/date"

type LeavePeriod struct {
	date.DateStringRange
	Type    string
	Comment string
}

type VacationInfo struct {
	Vacation      []LeavePeriod
	EarnedDays    float64
	AvailableDays float64
	TakenDays     float64
}
