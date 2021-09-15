package items

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
	COMMON   Rarity = "common"
	UNCOMMON Rarity = "uncommon"
	RARE     Rarity = "rare"
)

type IItem interface {
	GetID() int
	GetName() string
	GetType() ItemType
}

type ItemBundle struct {
	ID  int
	Qty int
}

func (ib ItemBundle) GetItem() IItem {
	return GetItemCollection().GetItemById(ib.ID)
}

type BaseItem struct {
	ID     int
	Name   string
	Type   ItemType
	Tier   int
	Rarity Rarity
}

func (i BaseItem) GetName() string {
	return i.Name
}

func (i BaseItem) GetType() ItemType {
	return i.Type
}

func (i BaseItem) GetID() int {
	return i.ID
}

type Resource struct {
	BaseItem
}

type Weapon struct {
	BaseItem

	BaseDamage      int
	BaseAttackSpeed float64

	Intellect int
	Strength  int
	Agility   int

	MinLvlReq int
}

type Recipe struct {
	BaseItem

	Ingridients []ItemBundle
	Output      int
}
