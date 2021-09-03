package components

import (
	"github.com/alexglazkov9/survgram/items"
)

type InventoryComponent struct {
	BaseComponent

	Items []items.IItem
	Slots int
}

func (ic InventoryComponent) Update(float64) {

}

func (ic *InventoryComponent) AddItems(items ...items.IItem) {
	ic.Items = append(ic.Items, items...)
}
