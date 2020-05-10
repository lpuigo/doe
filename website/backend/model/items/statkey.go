package items

const (
	StatSerieWork          string = "Work"
	StatSeriePrice         string = "Price"
	StatSiteProgress       string = "Réal."
	StatSiteProgressTarget string = "Réal. Cible"
	StatSerieWorkTarget    string = "WorkTarget"
	StatSeriePriceTarget   string = "PriceTarget"
)

type StatKey struct {
	Team    string
	Date    string
	Site    string
	Article string
	Serie   string
}
