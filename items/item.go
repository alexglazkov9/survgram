package items

import (
	"fmt"
	"strings"
)

type ItemType string

const (
	WEAPON     ItemType = "weapon"
	EQUIPMENT  ItemType = "equipment"
	RESOURCE   ItemType = "resource"
	RECIPE     ItemType = "recipe"
	BACKPACK   ItemType = "backpack"
	IMPROVER   ItemType = "improver"
	CONSUMABLE ItemType = "consumable"
	TOOL       ItemType = "tool"
)

type Rarity string

const (
	COMMON    Rarity = "common"
	UNCOMMON  Rarity = "uncommon"
	RARE      Rarity = "rare"
	EPIC      Rarity = "epic"
	LEGENDARY Rarity = "legendary"
)

type EquipmentType string

const (
	HEAD   EquipmentType = "head"
	BELT   EquipmentType = "belt"
	BODY   EquipmentType = "body"
	HANDS  EquipmentType = "hand"
	LEGS   EquipmentType = "legs"
	FEET   EquipmentType = "feet"
	RING   EquipmentType = "ring"
	AMULET EquipmentType = "amulet"
)

type DamageType string

const (
	MAGICAL  DamageType = "magical"
	PHYSICAL DamageType = "physical"
)

type IItem interface {
	GetID() int
	GetName() string
	GetType() ItemType
	GetTier() int
	GetRarity() string
	GetDescription() string
	GetFormattedItem(include_markdown bool) string
}

type ItemBundle struct {
	ID  int
	Qty int
}

func (ib ItemBundle) GetItem() IItem {
	return GetItemCollection().GetItemById(ib.ID)
}

/* BaseItem */

type BaseItem struct {
	ID          int
	Name        string
	Type        ItemType
	Tier        int
	Rarity      Rarity
	Description string
}

func (i BaseItem) GetID() int {
	return i.ID
}

func (i BaseItem) GetName() string {
	return i.Name
}

func (i BaseItem) GetType() ItemType {
	return i.Type
}

func (i BaseItem) GetTier() int {
	return i.Tier
}

func (i BaseItem) GetRarity() string {
	var prefix string
	switch i.Rarity {
	case COMMON:
		prefix = "⚪"
	case UNCOMMON:
		prefix = "🟢"
	case RARE:
		prefix = "🔵"
	case EPIC:
		prefix = "🟣"
	case LEGENDARY:
		prefix = "🟠"
	}
	return fmt.Sprintf("%s%s", prefix, strings.Title(string(i.Rarity)))
}

func (i BaseItem) GetDescription() string {
	return i.Description
}

func (i BaseItem) GetFormattedItem(include_markdown bool) string {
	var text string
	if include_markdown {
		text = fmt.Sprintf("*%s* - %s\n`%s` - T%d\n\n_%s_", i.GetName(), strings.Title(string(i.GetType())), i.GetRarity(), i.GetTier(), i.GetDescription())
	} else {
		text = fmt.Sprintf("%s - %s\n%s - T%d\n\n%s", i.GetName(), strings.Title(string(i.GetType())), i.GetRarity(), i.GetTier(), i.GetDescription())
	}
	return text
}

/* Resource */
type Resource struct {
	BaseItem
}

/* Weapon */
type Weapon struct {
	BaseItem

	BaseDamage      int
	BaseAttackSpeed float64

	Intellect int
	Strength  int
	Agility   int

	MinLvlReq int
}

func (i Weapon) GetFormattedItem(include_markdown bool) string {
	var text string
	if include_markdown {
		text = fmt.Sprintf("*INT:* %d\n*AGI:* %d\n*STR:* %d\n", i.Intellect, i.Agility, i.Strength)
	} else {
		text = fmt.Sprintf("INT: %d\nAGI: %d\nSTR: %d\n", i.Intellect, i.Agility, i.Strength)
	}
	return i.BaseItem.GetFormattedItem(include_markdown) + text
}

/* Recipe */
type Recipe struct {
	BaseItem

	Ingridients []ItemBundle
	Output      int
}

/* Equipment */
type Equipment struct {
	BaseItem

	EquipmentType EquipmentType

	Armor     int
	ArmorType DamageType

	Intellect int
	Strength  int
	Agility   int

	MinLvlReq int
}

func (i Equipment) GetFormattedItem(include_markdown bool) string {
	var text string
	if include_markdown {
		text = fmt.Sprintf("*INT:* %d\n*AGI:* %d\n*STR:* %d\n", i.Intellect, i.Agility, i.Strength)
	} else {
		text = fmt.Sprintf("INT: %d\nAGI: %d\nSTR: %d\n", i.Intellect, i.Agility, i.Strength)
	}
	return i.BaseItem.GetFormattedItem(include_markdown) + text
}

func (i Equipment) GetEquipmentType() EquipmentType {
	return i.EquipmentType
}

/* Tool */
type Tool struct {
	BaseItem
}

/* Backpack */
type Backpack struct {
	BaseItem
}

/* Improver */
type Improver struct {
	BaseItem
}

/* Consumable */
type Consumable struct {
	BaseItem
}
