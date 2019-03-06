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

func TodayAfter(d int) string {
	t := time.Now().Truncate(24 * time.Hour).Add(time.Duration(d*24) * time.Hour)
	return t.Format(TimeJSLayout)
}

func DateString(v string) string {
	if strings.Contains(v, "-") {
		d := strings.Split(v, "-")
		return d[2] + "/" + d[1] + "/" + d[0]
	}
	return "-"
}
