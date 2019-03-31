package date

import "testing"

func TestGetMonth(t *testing.T) {
	date := "2019-03-25"
	resdate := "2019-03-01"
	if GetMonth(date) != resdate {
		t.Errorf("GetMonth(%s) returns unexpected '%s'", date, resdate)
	}
}
