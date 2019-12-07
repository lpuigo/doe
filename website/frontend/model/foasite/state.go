package foasite

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/model/foasite/foaconst"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

// type State reflects backend/model/foasites.State struct
type State struct {
	*js.Object
	Status  string   `js:"Status"`
	Actors  []string `js:"Actors"`
	Date    string   `js:"Date"`
	Comment string   `js:"Comment"`
}

func NewState() *State {
	ns := &State{Object: tools.O()}
	ns.Status = foaconst.StateToDo
	ns.Actors = []string{}
	ns.Date = ""
	ns.Comment = ""
	return ns
}

// Copy copies attributes of os in s
func (s *State) Copy(os *State) {
	s.Status = os.Status
	s.Actors = os.Actors[:]
	s.Date = os.Date
	s.Comment = os.Comment
}

func (s *State) IsCanceled() bool {
	switch s.Status {
	case foaconst.StateCancelled:
		return true
	default:
		return false
	}
}

func (s *State) IsToDo() bool {
	switch s.Status {
	case foaconst.StateToDo:
		return true
	default:
		return false
	}
}
