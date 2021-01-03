package model

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
	"strconv"
)

type Team struct {
	*js.Object

	Name     string `js:"Name"`
	Members  string `js:"Members"`
	IsActive bool   `js:"IsActive"`
}

type Actor struct {
	*js.Object

	Id        int    `js:"Id"`
	LastName  string `js:"LastName"`
	FirstName string `js:"FirstName"`
	Role      string `js:"Role"`
	Active    bool   `js:"Active"`
	Assigned  bool   `js:"Assigned"`
}

func (a *Actor) GetRef() string {
	ext := ""
	if !a.Active {
		ext = " (parti)"
	} else if !a.Assigned {
		ext = " (réaffecté)"
	} else {
		ext = " (" + a.Role + ")"
	}
	return a.LastName + " " + a.FirstName + ext
}

func (a *Actor) GetElementsValueLabelDisabled() *elements.ValueLabelDisabled {
	active := a.Active && a.Assigned
	return elements.NewValueLabelDisabled(strconv.Itoa(a.Id), a.GetRef(), !active)
}

//res := []*elements.ValueLabelDisabled{}
//for _, actor := range client.Actors {
//res = append(res, elements.NewValueLabelDisabled(strconv.Itoa(actor.Id), actor.GetRef(), !actor.Active))
//}
//return res
