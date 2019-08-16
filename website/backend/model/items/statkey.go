package items

const (
	StatSerieWork  string = "Work"
	StatSeriePrice string = "Price"
)

type StatKey struct {
	Team    string
	Date    string
	Site    string
	Article string
	Serie   string
}

type Stats map[StatKey]float64

func (s Stats) AddStatValue(site, team, date, article, serie string, value float64) {
	s[StatKey{
		Team:    team,
		Date:    date,
		Site:    site,
		Article: article,
		Serie:   serie,
	}] += value
}
