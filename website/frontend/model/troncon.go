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
	Blockage      bool   `js:"Blockage"`
	NeedSignature bool   `js:"NeedSignature"`
	Signed        bool   `js:"Signed"`
	InstallDate   string `js:"InstallDate"`
	InstallActor  string `js:"InstallActor"`
	MeasureDate   string `js:"MeasureDate"`
	MeasureActor  string `js:"MeasureActor"`
	Comment       string `js:"Comment"`
}

func NewTroncon() *Troncon {
	tr := &Troncon{Object: tools.O()}
	tr.Ref = ""
	tr.Pb = NewPT()
	tr.NbRacco = 3
	tr.NbFiber = 6
	tr.Blockage = false
	tr.NeedSignature = false
	tr.Signed = false
	tr.InstallDate = ""
	tr.InstallActor = ""
	tr.MeasureDate = ""
	tr.MeasureActor = ""
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
	tr.Blockage = otr.Blockage
	tr.NeedSignature = otr.NeedSignature
	tr.Signed = otr.Signed
	tr.InstallDate = otr.InstallDate
	tr.InstallActor = otr.InstallActor
	tr.MeasureDate = otr.MeasureDate
	tr.MeasureActor = otr.MeasureActor
	tr.Comment = otr.Comment
}

func (tr *Troncon) SearchInString() string {
	res := "T_Ref:" + tr.Ref + "\n"
	res += "T_Pb:" + tr.Pb.SearchInString()
	res += "T_NbRacco:" + strconv.Itoa(tr.NbRacco) + "\n"
	res += "T_NbFiber:" + strconv.Itoa(tr.NbFiber) + "\n"
	res += "T_Blockage:" + strconv.FormatBool(tr.Blockage) + "\n"
	res += "T_NeedSignature:" + strconv.FormatBool(tr.NeedSignature) + "\n"
	res += "T_Signed:" + strconv.FormatBool(tr.Signed) + "\n"
	res += "T_InstallDate:" + date.DateString(tr.InstallDate) + "\n"
	res += "T_InstallActor:" + tr.InstallActor + "\n"
	res += "T_MeasureDate:" + date.DateString(tr.MeasureDate) + "\n"
	res += "T_MeasureActor:" + tr.MeasureActor + "\n"
	res += "T_Comment:" + tr.Comment + "\n"

	return res
}

func (tr *Troncon) IsCompleted() bool {
	if tr.Blockage {
		return true
	}
	if !tools.Empty(tr.InstallDate) && !tools.Empty(tr.MeasureDate) {
		return true
	}
	return false
}

func (tr *Troncon) IsFilledIn() bool {
	return !tools.Empty(tr.Ref) && tr.Pb.IsFilledIn()
}
