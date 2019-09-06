package actor

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/model/actor/actorconst"
	"github.com/lpuig/ewin/doe/website/frontend/tools/date"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
)

// Type Actor reflects ewin/doe/website/backend/model/actors.Actor
type Actor struct {
	*js.Object

	Id        int               `js:"Id"`
	Ref       string            `js:"Ref"`
	FirstName string            `js:"FirstName"`
	LastName  string            `js:"LastName"`
	State     string            `js:"State"`
	Period    *date.DateRange   `js:"Period"`
	Company   string            `js:"Company"`
	Contract  string            `js:"Contract"`
	Role      string            `js:"Role"`
	Vacation  []*date.DateRange `js:"Vacation"`
	Client    []string          `js:"Client"`
	Comment   string            `js:"Comment"`
}

func NewActorFromJS(obj *js.Object) *Actor {
	return &Actor{Object: obj}
}

func (a *Actor) SearchString(filter string) string {
	searchItem := func(prefix, typ, value string) string {
		if value == "" {
			return ""
		}
		if filter != actorconst.FilterValueAll && filter != typ {
			return ""
		}
		return prefix + typ + value
	}

	res := searchItem("", actorconst.FilterValueCompany, a.Company)
	res += searchItem("", actorconst.FilterValueName, a.Ref)
	res += searchItem("", actorconst.FilterValueComment, a.Comment)
	for _, clt := range a.Client {
		res += searchItem(",", actorconst.FilterValueClient, clt)
	}
	return res
}

func GetFilterTypeValueLabel() []*elements.ValueLabel {
	return []*elements.ValueLabel{
		elements.NewValueLabel(actorconst.FilterValueAll, actorconst.FilterLabelAll),
		elements.NewValueLabel(actorconst.FilterValueCompany, actorconst.FilterLabelCompany),
		elements.NewValueLabel(actorconst.FilterValueName, actorconst.FilterLabelName),
		elements.NewValueLabel(actorconst.FilterValueClient, actorconst.FilterLabelClient),
		elements.NewValueLabel(actorconst.FilterValueComment, actorconst.FilterLabelComment),
	}
}
