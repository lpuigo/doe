package actor

import (
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/date"
)

type CalendarSeeker struct {
	nbDates   int
	startDate string
	endDate   string

	dateIndex map[string]int

	nbEmployees []float64
	nbActing    []float64
}

func (cs *CalendarSeeker) NbEmployees() []float64 {
	return cs.nbEmployees
}

func NewCalendarSeeker(dates []string) *CalendarSeeker {
	cs := &CalendarSeeker{
		nbDates:     len(dates),
		startDate:   dates[0],
		endDate:     dates[len(dates)-1],
		dateIndex:   map[string]int{},
		nbEmployees: make([]float64, len(dates)),
		nbActing:    make([]float64, len(dates)),
	}
	// dateIndex is set with 1 offset (value from 1 to nbDates), so that lookup on "out of range date" returns 0
	for index, day := range dates {
		cs.dateIndex[day] = index + 1
	}

	return cs
}

func (cs *CalendarSeeker) GetDateIndex(d string) (int, bool) {
	if tools.Empty(d) {
		return 0, false
	}
	return cs.getMondayIndex(date.GetMonday(d))
}

func (cs *CalendarSeeker) getMondayIndex(monday string) (int, bool) {
	pos := cs.dateIndex[monday]
	if pos > 0 {
		return pos - 1, true
	}
	return 0, false
}

func (cs *CalendarSeeker) Append(actr *Actor) {
	// check if actors Period matches Calendar Timelap
	if !(actr.Period.Begin != "" && date.GetMonday(actr.Period.Begin) <= cs.endDate) {
		return
	}
	if actr.Period.End != "" && date.GetMonday(actr.Period.End) < cs.startDate {
		return
	}
	// Set Actor Start and Leaving weeks
	var startPos, endPos int
	var inRange bool
	startPos, inRange = cs.GetDateIndex(actr.Period.Begin)
	if !inRange {
		cs.nbEmployees[0]++
	} else {
		cs.nbEmployees[startPos]++
	}

	endPos, inRange = cs.GetDateIndex(actr.Period.End)
	if inRange && endPos < cs.nbDates-1 {
		cs.nbEmployees[endPos+1]--
	}

	// Calc Acting contribution (actor's vacations)
	var actualBegin, actualEnd string
	var vacStart, vacEnd string
	var posVacStart, posVacEnd int
	var vacStartInRange, vacEndInRange bool
	for _, vac := range actr.Vacation {
		if !(vac.Begin != "" && vac.End != "") { // skip Vacation not properly populated
			continue
		}
		actualBegin = vac.Begin
		actualEnd = vac.End
		if actr.Period.End != "" {
			if actualBegin > actr.Period.End { // skip vacation starting after actor's end date
				continue
			}
			if actr.Period.End < actualEnd {
				actualEnd = actr.Period.End
			}
		}

		vacStart, vacEnd = date.GetMonday(actualBegin), date.GetMonday(actualEnd)
		if !(vacStart <= cs.endDate && vacEnd >= cs.startDate) {
			continue
		}

		// calc first week of vacation
		posVacStart, vacStartInRange = cs.getMondayIndex(vacStart)
		if vacStartInRange {
			if vacStart != actualBegin {
				vacStartDay := date.NbDaysBetween(vacStart, actualBegin)
				if vacStartDay < 5.0 {
					cs.nbActing[posVacStart] -= (5.0 - vacStartDay) / 5.0
				}
			} else {
				cs.nbActing[posVacStart]--
			}
		} else {
			posVacStart = -1
		}

		// calc last week of vacation
		posVacEnd, vacEndInRange = cs.getMondayIndex(vacEnd)
		if vacEndInRange {
			vacEndDay := date.NbDaysBetween(vacEnd, actualEnd) + 1
			if posVacEnd == posVacStart {
				if vacEndDay < 5.0 {
					cs.nbActing[posVacEnd] += (5 - vacEndDay) / 5.0
				}
			} else {
				if vacEndDay < 5.0 {
					cs.nbActing[posVacEnd] -= vacEndDay / 5.0
				} else {
					cs.nbActing[posVacEnd]--
				}
			}
		} else {
			posVacEnd = cs.nbDates
		}

		/// then fill weeks between beg and end
		for i := posVacStart + 1; i < posVacEnd; i++ {
			cs.nbActing[i]--
		}
	}
}

func (cs *CalendarSeeker) CalcStats() (nbEmployees, nbActing []float64) {
	for i := 0; i < cs.nbDates; i++ {
		if i > 0 {
			cs.nbEmployees[i] += cs.nbEmployees[i-1]
		}
		cs.nbActing[i] += cs.nbEmployees[i]
	}
	return cs.nbEmployees, cs.nbActing
}
