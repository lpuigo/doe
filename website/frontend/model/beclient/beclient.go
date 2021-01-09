package beclient

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
)

// type BeClient reflects backend/model/clients.Client struct
type BeClient struct {
	*js.Object
	Id                             int                         `js:"Id"`
	Name                           string                      `js:"Name"`
	Activities                     map[string]CategoryArticles `js:"Activities"`
	Boxes                          map[string]map[string]*Box  `js:"Boxes"`
	OptionMeasurementPricePerFiber bool                        `js:"OptionMeasurementPricePerFiber"`
}

func NewBeClient() *BeClient {
	bc := &BeClient{Object: tools.O()}
	bc.Id = -10000
	bc.Name = ""
	bc.Activities = make(map[string]CategoryArticles)
	bc.Boxes = make(map[string]map[string]*Box)
	bc.OptionMeasurementPricePerFiber = false
	return bc
}

func BeClientFromJS(obj *js.Object) *BeClient {
	return &BeClient{Object: obj}
}

type CategoryArticles map[string][]*Article

type Article struct {
	*js.Object
	Name  string  `js:"Name"`
	Unit  int     `js:"Unit"`
	Price float64 `js:"Price"`
	Work  float64 `js:"Work"`
}

func NewArticle() *Article {
	a := &Article{Object: tools.O()}
	a.Name = ""
	a.Unit = 0
	a.Price = 0
	a.Work = 0
	return a
}

type Box struct {
	*js.Object
	Name  string `js:"Name"`
	Size  int    `js:"Size"`
	Usage string `js:"Usage"`
}

func NewBox() *Box {
	b := &Box{Object: tools.O()}
	b.Name = ""
	b.Size = 0
	b.Usage = ""
	return b
}
