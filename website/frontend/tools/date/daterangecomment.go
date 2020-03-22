package date

type DateRangeComment struct {
	DateRange
	Comment string `js:"Comment"`
}

func NewDateRangeComment() *DateRangeComment {
	drc := &DateRangeComment{DateRange: *NewDateRange()}
	drc.Begin = ""
	drc.End = ""
	drc.Comment = ""
	return drc
}

func NewDateRangeCommentFrom(beg, end, cmt string) *DateRangeComment {
	dr := &DateRangeComment{DateRange: *NewDateRange()}
	dr.Begin = beg
	dr.End = end
	dr.Comment = cmt
	return dr
}
