package components

import "github.com/alexglazkov9/survgram/items"

type HideoutComponent struct {
	BaseComponent

	Storage      []*items.ItemBundle
	StorageSlots int
}
