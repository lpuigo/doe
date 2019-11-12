package date

import (
	"strings"
	"time"
)

const TimeJSLayout string = "2006-01-02"

func JSDate(s string) int64 {
	res := time.Time{}
	res, _ = time.Parse(TimeJSLayout, s)
	return res.Unix() * 1000
}

func New(s string) time.Time {
	t, err := time.Parse(TimeJSLayout, s)
	if err != nil {
		return time.Time{}
	}
	return t
}

func NbDaysBetween(beg, end string) float64 {
	if beg == end {
		return 0
	}
	b := New(beg)
	e := New(end)
	return float64(e.Sub(b) / time.Duration(24*time.Hour))
}

func MinMax(date ...string) (min, max string) {
	min = "9999"
	max = "0000"
	for _, d := range date {
		if d == "" {
			continue
		}
		if d >= max {
			max = d
		}
		if d <= min {
			min = d
		}
	}
	return
}

func After(s string, d int) string {
	t := New(s).Add(time.Duration(d*24) * time.Hour)
	return t.Format(TimeJSLayout)
}

func TodayAfter(d int) string {
	t := time.Now().Truncate(24 * time.Hour).Add(time.Duration(d*24) * time.Hour)
	return t.Format(TimeJSLayout)
}

// DateString convert Date (js format YYYY-MM-DD) to DD/MM/YYYY
func DateString(v string) string {
	if strings.Contains(v, "-") {
		d := strings.Split(v, "-")
		return d[2] + "/" + d[1] + "/" + d[0]
	}
	return "-"
}

func Day(v string) string {
	if strings.Contains(v, "-") {
		d := strings.Split(v, "-")
		return d[2]
	}
	return "-"
}

func DayMonth(v string) string {
	if strings.Contains(v, "-") {
		d := strings.Split(v, "-")
		return d[2] + "/" + d[1]
	}
	return "-"
}

func MonthYear(v string) string {
	if strings.Contains(v, "-") {
		d := strings.Split(v, "-")
		return d[1] + "/" + d[0]
	}
	return "-"
}

func GetFirstOfMonth(v string) string {
	if strings.Contains(v, "-") {
		d := strings.Split(v, "-")
		return d[0] + "-" + d[1] + "-01"
	}
	return "-"
}

func GetMonday(v string) string {
	d := New(v)
	daynum := (int(d.Weekday()) + 6) % 7
	return d.Truncate(24 * time.Hour).Add(time.Duration(-daynum*24) * time.Hour).Format(TimeJSLayout)
}
