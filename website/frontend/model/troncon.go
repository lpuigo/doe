package model

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/dates"
	"strconv"
)

type Troncon struct {
	*js.Object

	Ref           string `js:"Ref"`
	Pb            *PT    `js:"Pb"`
	NbRacco       int    `js:"NbRacco"`
	NbFiber       int    `js:"NbFiber"`
	NeedSignature bool   `js:"NeedSignature"`
	Signed        bool   `js:"Signed"`
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
	tr.Signed = false
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
	tr.Signed = otr.Signed
	tr.InstallDate = otr.InstallDate
	tr.MeasureDate = otr.MeasureDate
	tr.Comment = otr.Comment
}

func (tr *Troncon) SearchInString() string {
	res := "T_Ref:" + tr.Ref + "\n"
	res += "T_Pb:" + tr.Pb.SearchInString()
	res += "T_NbRacco:" + strconv.Itoa(tr.NbRacco) + "\n"
	res += "T_NbFiber:" + strconv.Itoa(tr.NbFiber) + "\n"
	res += "T_NeedSignature:" + strconv.FormatBool(tr.NeedSignature) + "\n"
	res += "T_Signed:" + strconv.FormatBool(tr.Signed) + "\n"
	res += "T_InstallDate:" + date.DateString(tr.InstallDate) + "\n"
	res += "T_MeasureDate:" + date.DateString(tr.MeasureDate) + "\n"
	res += "T_Comment:" + tr.Comment + "\n"

	return res
}
