package date

import "testing"

func TestGetWeeksBetween(t *testing.T) {
	beg := DateFrom("2019-09-01")
	end := DateFrom("2019-09-30")

	weeks := GetWeeksBetween(beg, end)
	for i, dr := range weeks {
		t.Logf("week %d: %s - %s", i, dr.Begin.String(), dr.End.String())
	}
}

func TestGetMonthlyWeeksBetween(t *testing.T) {
	beg := DateFrom("2019-12-01")
	end := DateFrom("2020-01-31")

	weeks := GetMonthlyWeeksBetween(beg, end)
	for i, dr := range weeks {
		t.Logf("week %d: %s - %s", i, dr.Begin.String(), dr.End.String())
	}
}
