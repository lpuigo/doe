package date

import "github.com/lpuig/ewin/doe/website/frontend/tools/date"

type DateStringRange struct {
	Begin string
	End   string
}

// Overlap checks returns DateStringRange with common dates for both DateStringRange (empty DateStringRange if no overlap)
//
// It is assumed that both DateStringRange must be fully defined (ie Begin and End date are not empty)
func (dsr DateStringRange) Overlap(odsr DateStringRange) DateStringRange {
	if dsr.Begin > odsr.End || dsr.End < odsr.Begin {
		return DateStringRange{}
	}
	if odsr.Begin > dsr.Begin {
		dsr.Begin = odsr.Begin
	}
	if odsr.End < dsr.End {
		dsr.End = odsr.End
	}
	return dsr
}

func (dsr DateStringRange) IsEmpty() bool {
	return dsr.Begin == "" && dsr.End == ""
}

func (dsr DateStringRange) Duration() int {
	if dsr.IsEmpty() {
		return 0
	}
	return int(date.NbDaysBetween(dsr.Begin, dsr.End)) + 1
}

type DateRange struct {
	Begin Date
	End   Date
}

// GetWeeksBetween returns slide of DateRange containing all weeks beetween beg and end Dates (beg and end are included into first and last DateRange)
func GetWeeksBetween(beg, end Date) []DateRange {
	res := []DateRange{}

	const (
		WeekDuration    = 7
		PlusWorkingDays = 5 // Monday to Saturday
	)

	for currentBeg := beg.GetMonday(); !currentBeg.After(end); currentBeg = currentBeg.AddDays(WeekDuration) {
		res = append(res, DateRange{
			Begin: currentBeg,
			End:   currentBeg.AddDays(PlusWorkingDays),
		})
	}
	return res
}

func GetMonthlyWeeksBetween(beg, end Date) []DateRange {
	res := []DateRange{}

	for _, dr := range GetWeeksBetween(beg, end) {
		if dr.End.Before(beg) {
			continue
		}
		if dr.Begin.Before(beg) {
			dr.Begin = beg
		}
		if dr.End.After(end) {
			dr.End = end
		}
		month := dr.End.GetMonth()
		if !dr.Begin.GetMonth().Equal(month) {
			res = append(res, DateRange{dr.Begin, month.AddDays(-1)})
			dr.Begin = month
		}
		res = append(res, dr)
	}
	return res
}

func (dr *DateRange) ToDateStringRange() DateStringRange {
	return DateStringRange{
		Begin: dr.Begin.String(),
		End:   dr.End.String(),
	}
}
