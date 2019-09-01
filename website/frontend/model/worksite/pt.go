package worksite

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

type PT struct {
	*js.Object

	Ref     string `js:"Ref"`
	RefPt   string `js:"RefPt"`
	Address string `js:"Address"`
}

func NewPT() *PT {
	pt := &PT{Object: tools.O()}
	pt.Ref = ""
	pt.RefPt = ""
	pt.Address = ""
	return pt
}

func (pt *PT) Clone() *PT {
	npt := &PT{Object: tools.O()}
	npt.Copy(pt)
	return npt
}

func (pt *PT) Copy(opt *PT) {
	pt.Ref = opt.Ref
	pt.RefPt = opt.RefPt
	pt.Address = opt.Address
}

func (pt *PT) SearchInString() string {
	res := "Pt_Ref:" + pt.Ref + "\n"
	res += "Pt_RefPt:" + pt.RefPt + "\n"
	res += "Pt_Address:" + pt.Address + "\n"
	return res
}

func (pt *PT) IsFilledIn() bool {
	return !tools.Empty(pt.Ref) && !tools.Empty(pt.RefPt) && !tools.Empty(pt.Address)
}
