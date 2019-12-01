package date

import (
	"strings"
	"time"
)

type Date time.Time

type DateAggreg func(string) string

const (
	TimeJSLayout string = "2006-01-02"
	TimeLayout   string = "02/01/2006"
)

func DateFrom(d string) Date {
	checkedDate, err := ParseDate(d)
	if err != nil {
		return Date{}
	}
	return checkedDate
}

func ParseDate(d string) (Date, error) {
	date, err := time.Parse(TimeJSLayout, d)
	if err != nil {
		return Date{}, err
	}
	return Date(date), nil

}

func (d Date) ToTime() time.Time {
	return time.Time(d)
}

// String returns format YYYY-MM-DD date string
func (d Date) String() string {
	return time.Time(d).Format(TimeJSLayout)
}

func (d Date) ToDDMMYYYY() string {
	return time.Time(d).Format(TimeLayout)
}

func (d Date) GetMonday() Date {
	wd := int(d.ToTime().Weekday())
	if wd == 0 {
		wd = 7
	}
	wd--
	return Date(d.ToTime().AddDate(0, 0, -wd))
}

func (d Date) GetMonth() Date {
	t := d.ToTime()
	return Date(time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC))
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

// GetDayNum returns the week day number (0: monday -> 6: Sunday)
func GetDayNum(d string) int {
	wd := int(DateFrom(d).ToTime().Weekday())
	if wd == 0 {
		wd = 7
	}
	return wd - 1
}

func NbDaysBetween(beg, end string) int {
	b := DateFrom(beg)
	e := DateFrom(end)
	return int(float64(e.ToTime().Sub(b.ToTime()) / time.Duration(24*time.Hour)))
}

func GetMonth(d string) string {
	return DateFrom(d).GetMonth().String()
}

func ChangeDDMMYYYYtoYYYYMMDD(d string) string {
	cols := strings.Split(d, "/")
	return cols[2] + "-" + cols[1] + "-" + cols[0]
}

func ToDDMMYYYY(d string) string {
	cols := strings.Split(d, "-")
	return cols[2] + "/" + cols[1] + "/" + cols[0]
}

func GetDateAfter(d string, nbDay int) string {
	return DateFrom(d).AddDays(nbDay).String()
}
