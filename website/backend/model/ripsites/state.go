package ripsites

type State struct {
	Status    string
	Team      string
	DateStart string
	DateEnd   string
	Comment   string
}

const (
	StateToDo       string = "A faire"
	StateInProgress string = "En cours"
	StateBlocked    string = "Bloqué"
	StateDone       string = "Fait"
	StateCancelled  string = "Annulé"
)

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
	return !(s.Status == StateCancelled)
}

func (s State) Blocked() bool {
	return s.Status == StateBlocked
}

func (s State) Done() bool {
	return s.Status == StateDone
}
