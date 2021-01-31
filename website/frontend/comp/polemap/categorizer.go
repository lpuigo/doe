package polemap

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/model/polesite"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

type Categorizer struct {
	*js.Object

	CheckDA          bool                            `js:"CheckDA"`
	GroupName        func(marker *PoleMarker) string `js:"GroupName"`
	PoleMarkerVisual func(marker *PoleMarker)        `js:"PoleMarkerVisual"`
	GroupLabel       func(name string) string        `js:"GroupLabel"`
}

func NewCategorizer() *Categorizer {
	c := &Categorizer{Object: tools.O()}
	c.CheckDA = true
	c.SetByState()
	return c
}

func (c *Categorizer) SetByState() {
	c.GroupName = func(poleMarker *PoleMarker) string {
		return poleMarker.Pole.State
	}
	c.GroupLabel = func(name string) string {
		return polesite.PoleStateLabel(name)
	}
	c.PoleMarkerVisual = func(marker *PoleMarker) {
		html, class := marker.visualByState()
		marker.Html = html
		marker.Class = class
	}
}

func (c *Categorizer) SetByAge() {
	c.GroupName = func(poleMarker *PoleMarker) string {
		return poleMarker.Pole.GetPermissionDateRange(c.CheckDA)
	}
	c.GroupLabel = func(name string) string {
		return polesite.GetGroupNameByAge(name)
	}
	c.PoleMarkerVisual = func(marker *PoleMarker) {
		html, class := marker.visualByAge(c.GroupName(marker))
		marker.Html = html
		marker.Class = class
	}
}
