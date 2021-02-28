package vehicules

type InventoryItem struct {
	Name              string
	ReferenceQuantity int
	ControledQuantity int
	Comment           string
}

type Inventory struct {
	ReferenceDate string
	ControledDate string
	Items         []InventoryItem
	Comment       string
}
