package foasites

import "github.com/lpuig/ewin/doe/website/frontend/model/foasite/foaconst"

type State struct {
	Status  string
	Actors  []string
	Date    string
	Comment string
}

func NewState() *State {
	return &State{
		Status:  foaconst.StateToDo,
		Actors:  []string{},
		Date:    "",
		Comment: "",
	}
}
