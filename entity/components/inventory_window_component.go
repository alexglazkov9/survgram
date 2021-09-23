package components

import "github.com/alexglazkov9/survgram/entity"

type InventoryWindowComponent struct {
	BaseComponent

	Owner          *entity.Entity
	SelectedItemID *int
	IsSent         bool
}
