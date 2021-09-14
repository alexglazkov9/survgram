package components

import (
	"github.com/alexglazkov9/survgram/entity"
	"github.com/alexglazkov9/survgram/items"
)

type InventoryComponent struct {
	BaseComponent `bson:"-" json:"-"`

	Items []items.IItem
	Slots int `bson:"slots"`
}

func (ic InventoryComponent) Update(float64) {

}

func (ic *InventoryComponent) AddItems(items ...items.IItem) {
	ic.Items = append(ic.Items, items...)
}

func (ac *InventoryComponent) Clone() entity.IComponent {
	copy := *ac
	return &copy
}
