package actors

import "github.com/lpuig/ewin/doe/website/backend/model/date"

type Actor struct {
	Id        int
	Ref       string
	FirstName string
	LastName  string
	State     string
	Period    date.DateRange
	Company   string
	Contract  string
	Role      string
	Vacation  []date.DateRange
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
		Period:    date.DateRange{},
		Company:   company,
		Role:      "",
		Vacation:  []date.DateRange{},
		Client:    []string{},
		Comment:   "",
	}
}

func (a *Actor) IsActiveOn(date string) bool {
	if a.Period.End == "" {
		return true
	}
	return date < a.Period.End
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
