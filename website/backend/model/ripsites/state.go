package ripsites

import "github.com/lpuig/ewin/doe/website/frontend/model/ripsite/ripconst"

type State struct {
	Status    string
	Team      string
	DateStart string
	DateEnd   string
	Comment   string
}

func (s State) TodoBlockedDone() (todo, blocked, done bool) {
	if s.ToDo() {
		return true, s.Blocked(), s.Done()
	}
	return
}

func (s State) ToDo() bool {
	return !(s.Status == ripconst.StateCanceled)
}

func (s State) Blocked() bool {
	return s.Status == ripconst.StateBlocked
}

func (s State) Done() bool {
	switch s.Status {
	case ripconst.StateDone:
		return true
	case ripconst.StateWarning1:
		return true
	case ripconst.StateWarning2:
		return true
	default:
		return false
	}
}

func (s State) GetTodoDone() (todo, done bool) {
	todo = s.ToDo()
	done = s.Done()
	return
}

func (s *State) SetDone() {
	s.Status = ripconst.StateDone
}

func (s *State) SetRedo() {
	s.Status = ripconst.StateRedo
}

func (s *State) SetInProgress() {
	s.Status = ripconst.StateInProgress
}
