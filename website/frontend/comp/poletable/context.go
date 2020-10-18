package poletable

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/date"
)

type Context struct {
	*js.Object

	Mode         string `js:"Mode"`
	SelectedPole int    `js:"SelectedPole"`

	AttachmentVisible  bool     `js:"attachmentVisible"`
	AttachmentRange    []string `js:"attachmentRange"`
	AttachmentDate     string   `js:"attachmentDate"`
	AttachmentOverride bool     `js:"attachmentOverride"`

	RefGroupVisible bool `js:"refGroupVisible"`
}

const None int = -100

func NewContext(mode string) *Context {
	c := &Context{Object: tools.O()}
	c.Mode = mode
	c.SelectedPole = None

	c.AttachmentVisible = false
	c.AttachmentDate = date.TodayAfter(0)
	c.AttachmentRange = []string{"", ""}
	c.AttachmentOverride = false

	c.RefGroupVisible = false

	return c
}
