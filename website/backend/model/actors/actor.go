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
	Vacation  []date.DateStringRange
	Client    []string
	Comment   string
}

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
		Vacation:  []date.DateStringRange{},
		Client:    []string{},
		Comment:   "",
	}
}

func (a *Actor) IsActiveOn(date string) bool {
	if !(a.Period.Begin != "" && a.Period.Begin <= date) {
		return false
	}
	if a.Period.End == "" {
		return true
	}
	return date < a.Period.End
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

func (a *Actor) WorksForClient(client ...string) bool {
	for _, clt := range a.Client {
		for _, cl := range client {
			if clt == cl {
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
		onHoliday := activity.Overlap(vdr)
		if onHoliday.IsEmpty() {
			continue
		}
		comments = append(comments, fmt.Sprintf("En congés du %s au %s", date.ToDDMMYYYY(vdr.Begin), date.ToDDMMYYYY(vdr.End)))
		dur -= onHoliday.Duration()
	}
	if dur < 0 {
		dur = 0
	}
	return dur, strings.Join(comments, "\n")
}
