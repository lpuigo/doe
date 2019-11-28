package manager

import "github.com/lpuig/ewin/doe/website/backend/config"

type Calendar struct {
	file string
	days map[string]string
}

func NewCalendar(file string) *Calendar {
	return &Calendar{
		file: file,
		days: make(map[string]string),
	}
}

func (c Calendar) IsOff(day string) bool {
	_, found := c.days[day]
	return found
}

func (c *Calendar) SetOff(day, decr string) {
	c.days[day] = decr
}

func (c *Calendar) Reload() error {
	return config.SetFromFile(c.file, &c.days)
}

func (c *Calendar) NbDays() int {
	return len(c.days)
}

func (c *Calendar) GetDays() map[string]string {
	return c.days
}
