package worksite

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/dates"
	"strconv"
)

type Defect struct {
	*js.Object
	PT             string `js:"PT"`
	SubmissionDate string `js:"SubmissionDate"`
	Description    string `js:"Description"`
	NbOK           int    `js:"NbOK"`
	NbKO           int    `js:"NbKO"`
	ToBeFixed      bool   `js:"ToBeFixed"`
	FixDate        string `js:"FixDate"`
	FixActor       string `js:"FixActor"`
}

func NewDefect() *Defect {
	d := &Defect{Object: tools.O()}
	d.PT = ""
	d.SubmissionDate = ""
	d.Description = ""
	d.NbOK = 0
	d.NbKO = 0
	d.ToBeFixed = false
	d.FixDate = ""
	d.FixActor = ""
	return d
}

func (d *Defect) Clone() *Defect {
	nd := NewDefect()
	nd.Copy(d)
	return nd
}

func (d *Defect) Copy(od *Defect) {
	d.PT = od.PT
	d.SubmissionDate = od.SubmissionDate
	d.Description = od.Description
	d.NbOK = od.NbOK
	d.NbKO = od.NbKO
	d.ToBeFixed = od.ToBeFixed
	d.FixDate = od.FixDate
	d.FixActor = od.FixActor
}

func (d *Defect) SearchInString() string {
	res := "D_PT:" + d.PT + "\n"
	res += "D_Submit:" + date.DateString(d.SubmissionDate) + "\n"
	res += "D_Desc:" + d.Description + "\n"
	res += "D_NbOK:" + strconv.Itoa(d.NbOK) + "\n"
	res += "D_NbKO:" + strconv.Itoa(d.NbKO) + "\n"
	res += "D_ToBeFixed:" + strconv.FormatBool(d.ToBeFixed) + "\n"
	res += "D_FixDate:" + date.DateString(d.FixDate) + "\n"
	res += "D_FixActor:" + d.FixActor + "\n"
	return res
}

type Rework struct {
	*js.Object
	ControlDate    string    `js:"ControlDate"`
	SubmissionDate string    `js:"SubmissionDate"`
	CompletionDate string    `js:"CompletionDate"`
	Defects        []*Defect `js:"Defects"`
}

func NewRework() *Rework {
	rw := &Rework{Object: tools.O()}
	rw.ControlDate = ""
	rw.SubmissionDate = ""
	rw.CompletionDate = ""
	rw.Defects = []*Defect{}
	return rw
}

func (rw *Rework) Clone() *Rework {
	nrw := &Rework{Object: tools.O()}
	nrw.Copy(rw)
	return nrw
}

func (rw *Rework) Copy(orw *Rework) {
	rw.ControlDate = orw.ControlDate
	rw.SubmissionDate = orw.SubmissionDate
	rw.CompletionDate = orw.CompletionDate
	rw.Defects = []*Defect{}
	for _, d := range orw.Defects {
		rw.Defects = append(rw.Defects, d.Clone())
	}
}

func (rw *Rework) SearchInString() string {
	res := "RW_Control:" + date.DateString(rw.ControlDate) + "\n"
	res += "RW_Submit:" + date.DateString(rw.SubmissionDate) + "\n"
	res += "RW_Completion:" + date.DateString(rw.CompletionDate) + "\n"
	for _, d := range rw.Defects {
		res += d.SearchInString()
	}
	return res
}

func (rw *Rework) NeedRework() bool {
	if len(rw.Defects) == 0 {
		return false
	}
	for _, d := range rw.Defects {
		if d.ToBeFixed && tools.Empty(d.FixDate) {
			return true
		}
	}
	return false
}
