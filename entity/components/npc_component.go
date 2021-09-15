package components

import (
	"github.com/alexglazkov9/survgram/entity"
	"github.com/alexglazkov9/survgram/items"
)

type NPCComponent struct {
	BaseComponent `bson:"-" json:"-"`

	PossibleLoot []items.ItemBundle
}

func (nc NPCComponent) Update(float64) {

}

func (ac *NPCComponent) Clone() entity.IComponent {
	copy := *ac
	return &copy
}
