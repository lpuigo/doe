package actors

import "github.com/lpuig/ewin/doe/website/backend/model/date"

type Actor struct {
	Id        int
	Ref       string
	FirstName string
	LastName  string
	Period    date.DateRange
	Company   string
	Contract  string
	Role      string
	Vacation  []date.DateRange

	Client string
}

func NewActor(firstName, lastName, company string) *Actor {
	return &Actor{
		Id:        0,
		Ref:       lastName + " " + firstName,
		FirstName: firstName,
		LastName:  lastName,
		Period:    date.DateRange{},
		Company:   company,
		Role:      "",
		Vacation:  []date.DateRange{},
		Client:    "",
	}
}
