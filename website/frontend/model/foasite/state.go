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
