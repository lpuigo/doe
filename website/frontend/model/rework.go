package model

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/dates"
)

type Defect struct {
	*js.Object
	PT             string `js:"PT"`
	SubmissionDate string `js:"SubmissionDate"`
	Description    string `js:"Description"`
	FixDate        string `js:"FixDate"`
}

func NewDefect() *Defect {
	d := &Defect{Object: tools.O()}
	d.PT = ""
	d.SubmissionDate = ""
	d.Description = ""
	d.FixDate = ""
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
	d.FixDate = od.FixDate
}

func (d *Defect) SearchInString() string {
	res := "D_PT:" + d.PT + "\n"
	res += "D_Submit:" + date.DateString(d.SubmissionDate) + "\n"
	res += "D_Desc:" + d.Description + "\n"
	res += "D_Fix:" + date.DateString(d.FixDate) + "\n"
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
