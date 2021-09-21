package components

import (
	"github.com/alexglazkov9/survgram/entity"
	"github.com/alexglazkov9/survgram/items"
)

type InventoryComponent struct {
	BaseComponent `bson:"-" json:"-"`

	Items []*items.ItemBundle
	Slots int `bson:"slots"`
}

func (ic InventoryComponent) Update(float64) {

}

func (ic *InventoryComponent) AddItems(itms ...items.ItemBundle) {
	for _, newItm := range itms {
		itmExists := false
		for _, itm := range ic.Items {
			if itm.ID == newItm.ID {
				itm.Qty += newItm.Qty
				itmExists = true
				break
			}
		}
		if !itmExists {
			ic.Items = append(ic.Items, &items.ItemBundle{ID: newItm.ID, Qty: newItm.Qty})
		}
	}
}

func (ic *InventoryComponent) GetItems(args ...items.ItemType) []items.ItemBundle {
	if len(args) == 0 {
		itms := make([]items.ItemBundle, len(ic.Items))
		for i, itm := range ic.Items {
			itms[i] = *itm
		}
		return itms
	}

	itms := make([]items.ItemBundle, 0)
	for _, itm := range ic.Items {
		if containsType(args, itm.GetItem().GetType()) {
			itms = append(itms, *itm)
		}
	}
	return itms
}

func (ac *InventoryComponent) Clone() entity.IComponent {
	copy := *ac
	return &copy
}

func containsType(s []items.ItemType, e items.ItemType) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
