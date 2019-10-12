package ripsite

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/model/ripsite/ripconst"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
)

type State struct {
	*js.Object

	Status    string   `js:"Status"`
	Team      string   `js:"Team"`
	Actors    []string `js:"Actors"`
	DateStart string   `js:"DateStart"`
	DateEnd   string   `js:"DateEnd"`
	Comment   string   `js:"Comment"`
}

func NewState() *State {
	s := &State{Object: tools.O()}
	s.Status = ripconst.StateToDo
	s.Team = ""
	s.Actors = []string{}
	s.DateStart = ""
	s.DateEnd = ""
	s.Comment = ""

	return s
}

func (s *State) IsDone() bool {
	switch s.Status {
	case ripconst.StateDone, ripconst.StateWarning1, ripconst.StateWarning2:
		return true
	default:
		return false
	}
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

func (s *State) SetBlocked() {
	s.Status = ripconst.StateBlocked
}

func (s *State) SetWarning1() {
	s.Status = ripconst.StateWarning1
}

func (s *State) SetWarning2() {
	s.Status = ripconst.StateWarning2
}

func (s *State) SetInProgress() {
	switch s.Status {
	case ripconst.StateDone, ripconst.StateToDo:
		s.Status = ripconst.StateInProgress
	default:
		// do not change current status
	}
}

func (s *State) IsDoable() bool {
	switch s.Status {
	case ripconst.StateCanceled:
		return false
	default:
		return true
	}
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
	if len(s.Actors) == 0 {
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
		s.SetInProgress()
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

func GetStateStatusesWithWarningValueLabel() []*elements.ValueLabel {
	return []*elements.ValueLabel{
		elements.NewValueLabel(ripconst.StateToDo, ripconst.StateLabelToDo),
		elements.NewValueLabel(ripconst.StateInProgress, ripconst.StateLabelInProgress),
		elements.NewValueLabel(ripconst.StateBlocked, ripconst.StateLabelBlocked),
		elements.NewValueLabel(ripconst.StateWarning2, ripconst.StateLabelWarning2),
		elements.NewValueLabel(ripconst.StateWarning1, ripconst.StateLabelWarning1),
		elements.NewValueLabel(ripconst.StateDone, ripconst.StateLabelDone),
		elements.NewValueLabel(ripconst.StateCanceled, ripconst.StateLabelCanceled),
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
	case ripconst.StateWarning2:
		return "ripactivitystatus-row-warning2"
	case ripconst.StateWarning1:
		return "ripactivitystatus-row-warning1"
	case ripconst.StateDone:
		return "ripactivitystatus-row-done"
	case ripconst.StateCanceled:
		return "ripactivitystatus-row-canceled"
	default:
		return "worksite-row-error"
	}
}

func (s *State) GetLabel() string {
	return GetStatusLabel(s.Status)
}

func GetStatusLabel(s string) string {
	switch s {
	case ripconst.StateToDo:
		return ripconst.StateLabelToDo
	case ripconst.StateInProgress:
		return ripconst.StateLabelInProgress
	case ripconst.StateBlocked:
		return ripconst.StateLabelBlocked
	case ripconst.StateWarning2:
		return ripconst.StateLabelWarning2
	case ripconst.StateWarning1:
		return ripconst.StateLabelWarning1
	case ripconst.StateDone:
		return ripconst.StateLabelDone
	case ripconst.StateCanceled:
		return ripconst.StateLabelCanceled
	default:
		return s
	}
}
