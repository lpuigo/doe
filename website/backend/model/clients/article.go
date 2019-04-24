package clients

type Article struct {
	Name  string
	Unit  int
	Price float64
}

func (a Article) CalcPrice(qty int) float64 {
	return a.Price * float64(qty/a.Unit)
}
