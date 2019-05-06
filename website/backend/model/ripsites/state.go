package ripsites

import ferip "github.com/lpuig/ewin/doe/website/frontend/model/ripsite"

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
	return !(s.Status == ferip.StateCancelled)
}

func (s State) Blocked() bool {
	return s.Status == ferip.StateBlocked
}

func (s State) Done() bool {
	return s.Status == ferip.StateDone
}
