package poleedit

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/comp/polemap"
	"github.com/lpuig/ewin/doe/website/frontend/model/polesite"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/latlong"
	"strconv"
)

type DuplicateContext struct {
	*js.Object
	Model       *polemap.PoleMarker `js:"Model"`
	DoIncrement bool                `js:"DoIncrement"`
	Increment   int                 `js:"Increment"`
	RedoComment string              `js:"RedoComment"`
	Distance    float64             `js:"Distance"`
	NewPole     *polesite.Pole      `js:"NewPole"`
}

func NewDuplicateContext() *DuplicateContext {
	dp := &DuplicateContext{Object: tools.O()}
	dp.DoIncrement = false
	dp.Increment = 5
	dp.RedoComment = ""
	dp.Distance = 2
	dp.NewPole = polesite.NewPole()

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

func (dc *DuplicateContext) Duplicate(model *polemap.PoleMarker) {
	dc.Model = model
	dc.NewPole = model.Pole.Duplicate(dc.NewName(), 10*latlong.GpsMeter)
}

func (dc *DuplicateContext) Redo(model *polemap.PoleMarker) {
	dc.Model = model
	newName := model.Pole.Sticker + "_Repr"

	dc.NewPole = model.Pole.Redo(newName, dc.RedoComment, 2*latlong.GpsMeter)
}
