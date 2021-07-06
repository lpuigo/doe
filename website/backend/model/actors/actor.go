package actors

import (
	"fmt"
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	"strings"
)

type Actor struct {
	Id        int
	Ref       string
	FirstName string
	LastName  string
	State     string
	Period    date.DateStringRange
	Company   string
	Contract  string
	Role      string
	Vacation  []date.DateStringRangeComment
	VacInfo   VacationInfo
	Client    []string
	Groups    GroupHistory
	Comment   string
}

// ActorsByGroupId is a getter function to retrieve actors list by GroupId. returns nil if group's id not found
type ActorsByGroupId func(groupId int) []*Actor

func NewActor(firstName, lastName, company string) *Actor {
	return &Actor{
		Id:        0,
		Ref:       lastName + " " + firstName,
		FirstName: firstName,
		LastName:  lastName,
		State:     "",
		Period:    date.DateStringRange{},
		Company:   company,
		Role:      "",
		Vacation:  []date.DateStringRangeComment{},
		Client:    []string{},
		Groups:    NewGroupHistory(),
		Comment:   "",
	}
}

func (a *Actor) CheckConsistency() {
	if a.Groups == nil {
		a.Groups = NewGroupHistory()
		a.Groups[a.Period.Begin] = 0
	}
}

func (a *Actor) IsActiveOn(date string) bool {
	if !(a.Period.Begin != "" && a.Period.Begin <= date) {
		return false
	}
	if a.Period.End == "" {
		return true
	}
	return date <= a.Period.End
}

func (a *Actor) IsActiveOnWeek(weekDate string) bool {
	if !(a.Period.Begin != "" && date.GetMonday(a.Period.Begin) <= weekDate) {
		return false
	}
	if a.Period.End == "" {
		return true
	}
	return weekDate <= date.GetMonday(a.Period.End)
}

func (a *Actor) IsActiveOnDateRange(dr date.DateStringRange) bool {
	if a.Period.Begin == "" {
		return false
	}
	if a.Period.Begin > dr.End {
		return false
	}
	if a.Period.End != "" && a.Period.End < dr.Begin {
		return false
	}
	return true
}

// IsWorkingOn returns true if actor is active and is not on holiday for given date
func (a *Actor) IsWorkingOn(date string) bool {
	if !a.IsActiveOn(date) {
		return false
	}
	for _, vacation := range a.Vacation {
		if vacation.OverlapDate(date) {
			return false
		}
	}
	return true
}

// WorksForClient returns true if actor works for one of given clients (true if clients slice is empty)
func (a *Actor) WorksForClient(clients ...string) bool {
	if len(clients) == 0 {
		return true
	}
	for _, actorClient := range a.Client {
		for _, client := range clients {
			if actorClient == client {
				return true
			}

		}
	}
	return false
}

// GetActivityInfoFor returns nb days of activity for current actor on given daterange, and pertaining comment (leave or vacation dates ...)
//
// It is assumed that actor is active on given daterange
func (a *Actor) GetActivityInfoFor(dr date.DateStringRange) (int, string) {
	activity := a.Period
	if activity.End == "" {
		activity.End = dr.End
	}
	comments := []string{}

	if a.Period.Begin >= dr.Begin {
		comments = append(comments, fmt.Sprintf("Début d'activité le %s", date.ToDDMMYYYY(a.Period.Begin)))
	}
	if a.Period.End != "" && a.Period.End < dr.End {
		comments = append(comments, fmt.Sprintf("Fin d'activité le %s", date.ToDDMMYYYY(a.Period.End)))
	}

	activity = activity.Overlap(dr)
	if activity.IsEmpty() {
		return 0, "Pas de présence sur la semaine"
	}
	if date.GetDayNum(activity.End) == 5 { // saturday
		activity.End = date.GetDateAfter(activity.End, -1)
	}
	dur := activity.Duration()
	for _, vdr := range a.Vacation {
		onHoliday := activity.Overlap(vdr.DateStringRange)
		if onHoliday.IsEmpty() {
			continue
		}
		cmt := fmt.Sprintf("En congés du %s au %s", date.ToDDMMYYYY(vdr.Begin), date.ToDDMMYYYY(vdr.End))
		if vdr.Comment != "" {
			cmt += fmt.Sprintf(" (%s)", vdr.Comment)
		}
		comments = append(comments, cmt)
		dur -= onHoliday.Duration()
	}
	if dur < 0 {
		dur = 0
	}
	return dur, strings.Join(comments, "\n")
}

func (a *Actor) GetLabel() string {
	return "(" + a.Role + ") " + a.Ref
}
