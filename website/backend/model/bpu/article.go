package bpu

type Article struct {
	Name  string
	Unit  int
	Price float64
	Work  float64
}

func NewArticle() *Article {
	p := &Article{}
	return p
}

func (a Article) CalcPrice(qty int) float64 {
	return a.Price * float64(qty/a.Unit)
}
