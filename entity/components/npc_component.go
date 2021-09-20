package components

import (
	"github.com/alexglazkov9/survgram/items"
)

type LootDropConfig struct {
	PossibleLoot items.ItemBundle
	Chance       float64
}

type NPCComponent struct {
	BaseComponent `bson:"-" json:"-"`

	PossibleLoot []LootDropConfig
}
