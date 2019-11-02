package timesheets

type ActorsTime struct {
	Hours []int
}

type TimeSheet struct {
	Id          int
	WeekDate    string
	ActorsTimes map[string]*ActorsTime
}

func NewTimeSheet(weekdate string) *TimeSheet {
	ts := &TimeSheet{
		Id:          -1,
		WeekDate:    weekdate,
		ActorsTimes: make(map[string]*ActorsTime),
	}
	return ts
}
