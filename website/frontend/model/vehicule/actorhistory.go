package vehicule

import "github.com/gopherjs/gopherjs/js"

// Type ActorHistory reflects ewin/doe/website/backend/model/vehicules.ActorHistory
type ActorHistory struct {
	*js.Object

	Date    string `js:"Date"`
	ActorId int    `js:"ActorId"`
}
