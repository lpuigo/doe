package vehicule

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/date"
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

func NewInventoryItem() *InventoryItem {
	ni := &InventoryItem{Object: tools.O()}
	ni.Name = ""
	ni.ReferenceQuantity = 1
	ni.ControledQuantity = 0
	ni.Comment = ""
	return ni
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
	ni.ReferenceDate = date.TodayAfter(0)
	ni.ControledDate = ""
	ni.Comment = ""
	ni.Items = []*InventoryItem{}
	return ni
}

func NewInventoryFromModel(mi *Inventory) *Inventory {
	ni := NewInventory()
	items := make([]*InventoryItem, len(mi.Items))
	for index, item := range mi.Items {
		nii := NewInventoryItem()
		nii.Name = item.Name
		nii.ReferenceQuantity = item.ReferenceQuantity
		items[index] = nii
	}
	ni.Items = items
	return ni
}

func NewInventoryFromControledModel(mi *Inventory) *Inventory {
	ni := NewInventory()
	ni.ReferenceDate = date.After(mi.ControledDate, 1)
	ni.Comment = mi.Comment
	items := make([]*InventoryItem, len(mi.Items))
	for index, item := range mi.Items {
		nii := NewInventoryItem()
		nii.Name = item.Name
		nii.Comment = item.Comment
		nii.ReferenceQuantity = item.ControledQuantity
		items[index] = nii
	}
	ni.Items = items
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
