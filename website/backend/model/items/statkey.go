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
	Graph        string
	Date         string
	StackedSerie string
	Article      string
	Serie        string
}
