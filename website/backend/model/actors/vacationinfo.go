package actors

import "github.com/lpuig/ewin/doe/website/backend/model/date"

const (
	LeaveTypePaid          string = "Congés Payés"
	LeaveTypeUnpaid        string = "Congés Sans Solde"
	LeaveTypeSick          string = "Congés Maladie"
	LeaveTypePublicHoliday string = "Jour Férié"
	LeaveTypeInjury        string = "Accident Du Travail"
)

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
