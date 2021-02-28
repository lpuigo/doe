package vehicule

import "github.com/gopherjs/gopherjs/js"

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

	ReferenceDate string          `js:"ReferenceDate"`
	ControledDate string          `js:"ControledDate"`
	Items         []InventoryItem `js:"Items"`
	Comment       string          `js:"Comment"`
}
