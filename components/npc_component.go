package components

import (
	"github.com/alexglazkov9/survgram/items"
)

type NPCComponent struct {
	BaseComponent

	PossibleLoot []items.IItem
}

func (nc NPCComponent) Update(float64) {

}
