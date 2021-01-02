package extraactivity

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/json"
)

type ExtraActivity struct {
	*js.Object
	Name           string   `js:"Name"`
	State          string   `js:"State"`
	NbPoints       float64  `js:"NbPoints"`
	Income         float64  `js:"Income"`
	Date           string   `js:"Date"`
	AttachmentDate string   `js:"AttachmentDate"`
	Actors         []string `js:"Actors"`
	Comment        string   `js:"Comment"`
}

func NewExtraActivity() *ExtraActivity {
	ea := &ExtraActivity{Object: tools.O()}
	ea.Name = ""
	ea.State = ""
	ea.NbPoints = 0.0
	ea.Income = 0.0
	ea.Date = ""
	ea.AttachmentDate = ""
	ea.Actors = []string{}
	ea.Comment = ""
	return ea
}

func ExtraActivityFromJS(o *js.Object) *ExtraActivity {
	return &ExtraActivity{Object: o}
}

func (ea *ExtraActivity) Clone() *ExtraActivity {
	return &ExtraActivity{Object: json.Parse(json.Stringify(ea))}
}
