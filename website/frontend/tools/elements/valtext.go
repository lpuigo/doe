package elements

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

type ValText struct {
	*js.Object
	Value string `js:"value"`
	Text  string `js:"text"`
}

func NewValText(val, text string) *ValText {
	vt := &ValText{Object: tools.O()}
	vt.Value = val
	vt.Text = text
	return vt
}

func IsInValTextList(value string, vtl []*ValText) bool {
	for _, vt := range vtl {
		if vt.Value == value {
			return true
		}
	}
	return false
}

func NewValTextList(list *js.Object) []*ValText {
	res := []*ValText{}
	objlist := list.Interface().([]interface{})
	for _, o := range objlist {
		res = append(res, o.(*ValText))
	}
	return res
}

type ValueLabel struct {
	*js.Object
	Value string `js:"value"`
	Label string `js:"label"`
}

func NewValueLabel(value, label string) *ValueLabel {
	vl := &ValueLabel{Object: tools.O()}
	vl.Value = value
	vl.Label = label
	return vl
}

type ValueLabelDisabled struct {
	*js.Object
	Value    string `js:"value"`
	Label    string `js:"label"`
	Disabled bool   `js:"disabled"`
}

func NewValueLabelDisabled(value, label string, disabled bool) *ValueLabelDisabled {
	vl := &ValueLabelDisabled{Object: tools.O()}
	vl.Value = value
	vl.Label = label
	vl.Disabled = disabled
	return vl
}
