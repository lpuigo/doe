package frontmodel

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

type Troncon struct {
	*js.Object

	Ref           string `js:"Ref"`
	Pb            *PT    `js:"Pb"`
	NbRacco       int    `js:"NbRacco"`
	NbFiber       int    `js:"NbFiber"`
	NeedSignature bool   `js:"NeedSignature"`
	InstallDate   string `js:"InstallDate"`
	MeasureDate   string `js:"MeasureDate"`
	Comment       string `js:"Comment"`
}

func NewTroncon() *Troncon {
	tr := &Troncon{Object: tools.O()}
	tr.Ref = ""
	tr.Pb = NewPT()
	tr.NbRacco = 0
	tr.NbFiber = 0
	tr.NeedSignature = false
	tr.InstallDate = ""
	tr.MeasureDate = ""
	tr.Comment = ""
	return tr
}

func (tr *Troncon) Clone() *Troncon {
	ntr := &Troncon{Object: tools.O()}
	ntr.Copy(tr)
	return ntr
}

func (tr *Troncon) Copy(otr *Troncon) {
	tr.Ref = otr.Ref
	tr.Pb = otr.Pb.Clone()
	tr.NbRacco = otr.NbRacco
	tr.NbFiber = otr.NbFiber
	tr.NeedSignature = otr.NeedSignature
	tr.InstallDate = otr.InstallDate
	tr.MeasureDate = otr.MeasureDate
	tr.Comment = otr.Comment
}
