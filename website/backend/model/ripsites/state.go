package ripsites

import "github.com/lpuig/ewin/doe/website/frontend/model/ripsite/ripconst"

type State struct {
	Status    string
	Team      string
	DateStart string
	DateEnd   string
	Comment   string
}

func (s State) TotalBlockedDone() (total, blocked, done int) {
	if s.ToDo() {
		total++
		if s.Blocked() {
			blocked++
		}
		if s.Done() {
			done++
		}
	}
	return
}

func (s State) ToDo() bool {
	return !(s.Status == ripconst.StateCancelled)
}

func (s State) Blocked() bool {
	return s.Status == ripconst.StateBlocked
}

func (s State) Done() bool {
	return s.Status == ripconst.StateDone
}
