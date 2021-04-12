package date

import (
	"strings"
	"time"
)

type Date time.Time

type DateAggreg func(string) string

const (
	TimeStampLayout      string = "2006-01-02 15:04:05.000"
	TimeStampShortLayout string = "2006-01-02 150405"
	TimeJSLayout         string = "2006-01-02"
	TimeLayout           string = "02/01/2006"

	TimeJSMinDate string = "0000-01-01"
	TimeJSMaxDate string = "9999-12-31"
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

// String returns format YYYY-MM-DD HH:MM:SS.mmm date string
func (d Date) TimeStamp() string {
	return time.Time(d).Format(TimeStampLayout)
}

// String returns format YYYY-MM-DD HHMMSS date string
func (d Date) TimeStampShort() string {
	return time.Time(d).Format(TimeStampShortLayout)
}

// TODDMMYYYY returns format DD/MM/YYYY date string
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

func (d Date) GetWeekDay() int {
	return int(d.ToTime().Weekday())
}

func (d Date) IsSaturdaySunday() bool {
	wd := d.GetWeekDay()
	return wd == 0 || wd == 6
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

func Now() Date {
	return Date(time.Now())
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
