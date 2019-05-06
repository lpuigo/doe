package ripsite

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

type Troncon struct {
	*js.Object

	Name string `js:"Name"`
	Size int    `js:"Size"`
}

func NewTroncon() *Troncon {
	tr := &Troncon{Object: tools.O()}
	tr.Name = ""
	tr.Size = 0

	return tr
}
