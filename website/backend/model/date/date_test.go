package date

import "testing"

func TestGetMonth(t *testing.T) {
	date := "2019-03-25"
	resdate := "2019-03-01"
	if GetMonth(date) != resdate {
		t.Errorf("GetMonth(%s) returns unexpected '%s'", date, resdate)
	}
}

func TestNbDaysBetween(t *testing.T) {
	for _, tc := range []struct {
		date1, date2 string
		diff         int
	}{
		{"2019-10-28", "2019-10-28", 0},
		{"2019-10-28", "2019-11-01", 4},
		{"2019-11-01", "2019-10-28", -4},
		{"2019-11-01", "2019-12-31", 60},
		{"2019-01-01", "2020-01-01", 365},
	} {
		res := NbDaysBetween(tc.date1, tc.date2)
		if res != tc.diff {
			t.Errorf("NbDaysBetween(%s, %s) returns unexpected %d (expected %d)", tc.date1, tc.date2, res, tc.diff)
		}
	}
}
