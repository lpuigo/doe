package actor

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/model/actor/actorconst"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/date"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements"
	"github.com/lpuig/ewin/doe/website/frontend/tools/json"
)

// Type Actor reflects ewin/doe/website/backend/model/actors.Actor
type Actor struct {
	*js.Object

	Id        int                      `js:"Id"`
	Ref       string                   `js:"Ref"`
	FirstName string                   `js:"FirstName"`
	LastName  string                   `js:"LastName"`
	State     string                   `js:"State"`
	Period    *date.DateRange          `js:"Period"`
	Company   string                   `js:"Company"`
	Contract  string                   `js:"Contract"`
	Role      string                   `js:"Role"`
	Vacation  []*date.DateRangeComment `js:"Vacation"`
	Client    []string                 `js:"Client"`
	Comment   string                   `js:"Comment"`
}

func NewActor() *Actor {
	na := &Actor{Object: tools.O()}
	na.Id = -1
	na.Ref = ""
	na.FirstName = ""
	na.LastName = ""
	na.State = ""
	na.Period = date.NewDateRange()
	na.Company = ""
	na.Contract = ""
	na.Role = ""
	na.Vacation = []*date.DateRangeComment{}
	na.Client = []string{}
	na.Comment = ""
	return na
}

func NewActorFromJS(obj *js.Object) *Actor {
	return &Actor{Object: obj}
}

func (a *Actor) Copy() *Actor {
	return NewActorFromJS(json.Parse(json.Stringify(a.Object)))
}

func (a *Actor) Clone(oa *Actor) {
	a.Id = oa.Id
	a.Ref = oa.Ref
	a.FirstName = oa.FirstName
	a.LastName = oa.LastName
	a.State = oa.State
	a.Period.Begin = oa.Period.Begin
	a.Period.End = oa.Period.End
	a.Company = oa.Company
	a.Contract = oa.Contract
	a.Role = oa.Role
	a.Vacation = []*date.DateRangeComment{}
	for _, vac := range oa.Vacation {
		a.Vacation = append(a.Vacation, date.NewDateRangeCommentFrom(vac.Begin, vac.End, vac.Comment))
	}
	a.Client = oa.Client[:]
	a.Comment = oa.Comment
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

func (a *Actor) UpdateState() {
	if tools.Empty(a.Period.Begin) {
		a.Period.Begin = ""
	}
	if tools.Empty(a.Period.End) {
		a.Period.End = ""
	}
	today := date.TodayAfter(0)
	holidayPeriod := a.GetNextVacation()
	switch {
	case a.Period.Begin == "" && a.Period.End == "":
		a.State = actorconst.StateDefection
	case a.Period.Begin == "" || a.Period.Begin > today:
		a.State = actorconst.StateCandidate
	case a.Period.End != "" && a.Period.End < today:
		a.State = actorconst.StateGone
	case holidayPeriod == nil:
		a.State = actorconst.StateActive
	case holidayPeriod.Begin > today:
		a.State = actorconst.StateActive
	case holidayPeriod.Begin <= today:
		a.State = actorconst.StateOnHoliday
	default:
		a.State = "Error"
	}
}

// GetNextVacation returns actor's next (or current) vacation
func (a *Actor) GetNextVacation() *date.DateRange {
	if len(a.Vacation) == 0 {
		return nil
	}
	today := date.TodayAfter(0)
	vacBegin := ""
	vacEnd := ""
	for _, vacPeriod := range a.Vacation {
		if vacPeriod.End < today {
			continue
		}
		if vacBegin == "" && vacPeriod.End >= today {
			vacBegin = vacPeriod.Begin
			vacEnd = vacPeriod.End
			continue
		}
		// vacBegin != ""
		if vacPeriod.Begin < vacBegin {
			vacBegin = vacPeriod.Begin
			vacEnd = vacPeriod.End
		}
	}

	if vacBegin == "" {
		return nil
	}
	vdr := date.NewDateRange()
	vdr.Begin = vacBegin
	vdr.End = vacEnd
	return vdr
}

// GetActiveDays returns a slice of 6 int for given weekDate (-1 for inactive, 0 for Holydays, 1 for working day)
func (a *Actor) GetActiveDays(weekDate string, daysOff map[string]string) []int {
	res := make([]int, 6)
outer:
	for i := 0; i < 6; i++ {
		day := date.After(weekDate, i)
		if a.Period.Begin == "" {
			res[i] = -1
			continue outer
		}
		if !(day >= a.Period.Begin && !(a.Period.End != "" && day > a.Period.End)) {
			res[i] = -1
			continue outer
		}
		if daysOff[day] != "" {
			res[i] = 0
			continue outer
		}
		for _, vac := range a.Vacation {
			if day >= vac.Begin && day <= vac.End {
				res[i] = 0
				continue outer
			}
		}
		res[i] = 1
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
