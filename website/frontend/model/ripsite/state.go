package ripsite

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

const (
	StateToDo       string = "00 A faire"
	StateInProgress string = "10 En cours"
	StateBlocked    string = "20 Bloqué"
	StateDone       string = "90 Fait"
	StateCancelled  string = "99 Annulé"
)

type State struct {
	*js.Object

	Status    string `js:"Status"`
	Team      string `js:"Team"`
	DateStart string `js:"DateStart"`
	DateEnd   string `js:"DateEnd"`
	Comment   string `js:"Comment"`
}

func NewState() *State {
	s := &State{Object: tools.O()}
	s.Status = StateToDo
	s.Team = ""
	s.DateStart = ""
	s.DateEnd = ""
	s.Comment = ""

	return s
}

func (s *State) IsBlocked() bool {
	return s.Status == StateBlocked
}
