package vehicule

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/json"
)

// Type InventoryItem reflects ewin/doe/website/backend/model/vehicules.InventoryItem
type InventoryItem struct {
	*js.Object

	Name              string `js:"Name"`
	ReferenceQuantity int    `js:"ReferenceQuantity"`
	ControledQuantity int    `js:"ControledQuantity"`
	Comment           string `js:"Comment"`
}

// Type Inventory reflects ewin/doe/website/backend/model/vehicules.Inventory
type Inventory struct {
	*js.Object

	ReferenceDate string           `js:"ReferenceDate"`
	ControledDate string           `js:"ControledDate"`
	Items         []*InventoryItem `js:"Items"`
	Comment       string           `js:"Comment"`
}

func InventoryFromJS(obj *js.Object) *Inventory {
	return &Inventory{Object: obj}
}

func (i *Inventory) Copy() *Inventory {
	return InventoryFromJS(json.Parse(json.Stringify(i.Object)))
}

func NewInventory() *Inventory {
	ni := &Inventory{Object: tools.O()}
	ni.ReferenceDate = ""
	ni.ControledDate = ""
	ni.Comment = ""
	ni.Items = []*InventoryItem{}
	return ni
}

func CompareInventoryDate(a, b Inventory) int {
	if a.ReferenceDate > b.ReferenceDate {
		return -1
	}
	if a.ReferenceDate == b.ReferenceDate {
		return 0
	}
	return 1
}
