package ripsite

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/model/ripsite/ripconst"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
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
	s.Status = ripconst.StateToDo
	s.Team = ""
	s.DateStart = ""
	s.DateEnd = ""
	s.Comment = ""

	return s
}

func (s *State) IsBlocked() bool {
	return s.Status == ripconst.StateBlocked
}

func GetStateStatusesValueLabel() []*elements.ValueLabel {
	return []*elements.ValueLabel{
		elements.NewValueLabel(ripconst.StateToDo, "A faire"),
		elements.NewValueLabel(ripconst.StateInProgress, "En cours"),
		elements.NewValueLabel(ripconst.StateBlocked, "Bloqué"),
		elements.NewValueLabel(ripconst.StateDone, "Fait"),
		elements.NewValueLabel(ripconst.StateCancelled, "Annulé"),
	}
}
