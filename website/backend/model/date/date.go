package date

import "time"

type Date time.Time

const TimeJSLayout string = "2006-01-02"

func DateFrom(d string) Date {
	date, err := time.Parse(TimeJSLayout, d)
	if err != nil {
		panic("could not parse date format :" + err.Error())
	}
	return Date(date)
}

func (d Date) ToTime() time.Time {
	return time.Time(d)
}

func (d Date) String() string {
	return time.Time(d).Format(TimeJSLayout)
}

func (d Date) GetMonday() Date {
	wd := int(d.ToTime().Weekday())
	if wd == 0 {
		wd = 7
	}
	wd--
	return Date(d.ToTime().AddDate(0, 0, -wd))
}

func (d Date) AddDays(n int) Date {
	return Date(d.ToTime().AddDate(0, 0, n))
}

func (d Date) After(d2 Date) bool {
	return d.ToTime().After(time.Time(d2))
}

func (d Date) Before(d2 Date) bool {
	return d.ToTime().Before(time.Time(d2))
}

func (d Date) Equal(d2 Date) bool {
	return d.ToTime().Equal(time.Time(d2))
}

func Today() Date {
	return Date(time.Now().Truncate(24 * time.Hour))
}

func GetMonday(d string) string {
	return DateFrom(d).GetMonday().String()
}
