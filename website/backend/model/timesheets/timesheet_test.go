package timesheets

import "testing"

func TestTimeSheet_Merge(t *testing.T) {
	ts := NewMonthlyTimeSheetForActorsIds("2019-11-01", []int{1, 2})
	ots := []*TimeSheet{}
	for i, tc := range []struct {
		weekdate string
		ids      []int
	}{
		{"2019-10-28", []int{0, 1}},
		{"2019-11-04", []int{0, 1, 2}},
		{"2019-11-25", []int{3, 2}},
		{"2019-12-02", []int{2, 3}},
	} {
		ot := NewTimeSheetForActorsIds(tc.weekdate, tc.ids)
		for _, at := range ot.ActorsTimes {
			for ih, _ := range at.Hours {
				at.Hours[ih] = i + 1
			}
		}
		ots = append(ots, ot)

		ts.Merge(ot)
	}

	for id, at := range ts.ActorsTimes {
		t.Logf("%d: %v", id, at)
	}
}
