package components

type ItemType string

const (
	MELEE   ItemType = "melee"
	RANGED  ItemType = "ranged"
	MAGICAL ItemType = "ranged"

	BODY     ItemType = "body"
	FEET     ItemType = "feet"
	HANDS    ItemType = "hands"
	HEAD     ItemType = "head"
	LEGS     ItemType = "legs"
	NECKLACE ItemType = "necklace"
	RING     ItemType = "ring"

	RESOURCE ItemType = "resource"
	LOG      ItemType = "log"
	ORE      ItemType = "ore"
	BERRY    ItemType = "berry"
	FISH     ItemType = "fish"

	RECIPE ItemType = "recipe"
)

type Rarity string

const (
	COMMON    Rarity = "common"
	UNCOMMON  Rarity = "uncommon"
	RARE      Rarity = "rare"
	EPIC      Rarity = "epic"
	LEGENDARY Rarity = "legendary"
)

type BaseItemComponent struct {
	BaseComponent

	ID      int
	Name    string
	Type    ItemType
	Tier    int
	Rarirty Rarity
}
