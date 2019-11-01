package poletable

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

type Context struct {
	*js.Object

	Mode         string `js:"Mode"`
	SelectedPole int    `js:"SelectedPole"`
}

const None int = -100

func NewContext(mode string) *Context {
	c := &Context{Object: tools.O()}
	c.Mode = mode
	c.SelectedPole = None
	return c
}
