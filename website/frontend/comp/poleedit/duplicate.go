package poleedit

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/comp/polemap"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"strconv"
)

type DuplicateContext struct {
	*js.Object
	Model       *polemap.PoleMarker `js:"Model"`
	DoIncrement bool                `js:"DoIncrement"`
	Increment   int                 `js:"Increment"`
	Distance    float64             `js:"Distance"`
}

func NewDuplicateContext() *DuplicateContext {
	dp := &DuplicateContext{Object: tools.O()}
	dp.DoIncrement = false
	dp.Increment = 5
	dp.Distance = 35.0

	return dp
}

func DuplicateContextFromJS(o *js.Object) *DuplicateContext {
	return &DuplicateContext{Object: o}
}

func (dc *DuplicateContext) NewName() string {
	if !dc.DoIncrement {
		return dc.Model.Pole.Sticker + " copy"
	}
	stickernum, err := strconv.ParseInt(dc.Model.Pole.Sticker, 10, 32)
	if err != nil {
		return dc.Model.Pole.Sticker + " copy"
	}
	return strconv.Itoa(int(stickernum) + dc.Increment)
}
