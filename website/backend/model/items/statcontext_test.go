package items

import (
	"github.com/lpuig/ewin/doe/website/backend/model/date"
	"testing"
)

func TestNewStatContextMonth(t *testing.T) {
	st, err := NewStatContext("month")
	if err != nil {
		t.Fatalf("NewStatContext returned unexpected %s", err.Error())
	}

	today := "2020-02-01"
	for i := 0; i < 31; i++ {
		d := date.GetDateAfter(today, i)
		sd := date.GetMonth(date.GetDateAfter(date.GetMonth(d), (1-st.MaxVal)*30))
		t.Logf("%s -> %s\n", d, sd)
	}
}
