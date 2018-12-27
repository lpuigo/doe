package model

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/dates"
	"strings"
)

type Rework struct {
	*js.Object
	ControlDate    string   `js:"ControlDate"`
	SubmissionDate string   `js:"SubmissionDate"`
	CompletionDate string   `js:"CompletionDate"`
	Defects        []string `js:"Defects"`
}

func NewRework() *Rework {
	rw := &Rework{Object: tools.O()}
	rw.ControlDate = ""
	rw.SubmissionDate = ""
	rw.CompletionDate = ""
	rw.Defects = []string{}
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
	rw.Defects = []string{}
	for _, d := range orw.Defects {
		rw.Defects = append(rw.Defects, d)
	}
}

func (rw *Rework) SearchInString() string {
	res := "RW_Control:" + rw.ControlDate + "\n"
	res += "RW_Submit:" + date.DateString(rw.SubmissionDate) + "\n"
	res += "RW_Completion:" + date.DateString(rw.CompletionDate) + "\n"
	res += "RW_Defects:" + strings.Join(rw.Defects, ",") + "\n"
	return res
}
