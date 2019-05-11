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

func (s *State) IsCanceled() bool {
	return s.Status == ripconst.StateCanceled
}

func (s *State) SetDone() {
	s.Status = ripconst.StateDone
}

func (s *State) SetInProgress() {
	s.Status = ripconst.StateInProgress
}

func (s *State) SetToDo() {
	switch s.Status {
	case ripconst.StateInProgress, ripconst.StateDone:
		s.Status = ripconst.StateToDo
	default:
		// do not change current status
	}
}

func (s *State) UpdateStatus() {
	if tools.Empty(s.Team) {
		s.SetToDo()
		return
	}
	if !tools.Empty(s.DateEnd) {
		if tools.Empty(s.DateStart) {
			s.DateStart = s.DateEnd
		}
		s.Status = ripconst.StateDone
		return
	}
	// DateEnd is Empty
	if tools.Empty(s.DateStart) {
		s.SetToDo()
	} else {
		s.Status = ripconst.StateInProgress
	}
}

func GetStateStatusesValueLabel() []*elements.ValueLabel {
	return []*elements.ValueLabel{
		elements.NewValueLabel(ripconst.StateToDo, "A faire"),
		elements.NewValueLabel(ripconst.StateInProgress, "En cours"),
		elements.NewValueLabel(ripconst.StateBlocked, "Bloqué"),
		elements.NewValueLabel(ripconst.StateDone, "Fait"),
		elements.NewValueLabel(ripconst.StateCanceled, "Annulé"),
	}
}

func (s *State) GetRowStyle() string {
	switch s.Status {
	case ripconst.StateToDo:
		return "ripactivitystatus-row-todo"
	case ripconst.StateInProgress:
		return "ripactivitystatus-row-inprogress"
	case ripconst.StateBlocked:
		return "ripactivitystatus-row-blocked"
	case ripconst.StateDone:
		return "ripactivitystatus-row-done"
	case ripconst.StateCanceled:
		return "ripactivitystatus-row-canceled"
	default:
		return "worksite-row-error"
	}
}
