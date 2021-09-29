package components

import (
	"log"

	"github.com/alexglazkov9/survgram/items"
)

type EquipmentSlot int

const (
	WEAPON_SLOT   EquipmentSlot = 0
	HEAD_SLOT     EquipmentSlot = 1
	BODY_SLOT     EquipmentSlot = 2
	HANDS_SLOT    EquipmentSlot = 3
	LEGS_SLOT     EquipmentSlot = 4
	FEET_SLOT     EquipmentSlot = 5
	RING_SLOT     EquipmentSlot = 6
	AMULET_SLOT   EquipmentSlot = 7
	BELT_SLOT     EquipmentSlot = 8
	BACKPACK_SLOT EquipmentSlot = 9
	TOOL_SLOT     EquipmentSlot = 10
)

type PlayerEquipmentComponent struct {
	BaseComponent `bson:"-" json:"-"`

	Equipment []*ItemID
}

func (pec *PlayerEquipmentComponent) GetAll() []*ItemID {
	return pec.Equipment
}

func (pec *PlayerEquipmentComponent) TryEquip(id ItemID) bool {
	added := false
	log.Println("TYPE")
	log.Println(id.GetItem().GetType())
	switch id.GetItem().GetType() {
	case items.WEAPON:
		pec.Equipment[WEAPON_SLOT] = &id
		added = true
	case items.EQUIPMENT:
		log.Println("I GOT HERE")
		equipment := id.GetItem().(*items.Equipment)
		log.Println(equipment.GetEquipmentType())
		switch equipment.GetEquipmentType() {
		case items.HEAD:
			pec.Equipment[HEAD_SLOT] = &id
			added = true
		case items.BELT:
			pec.Equipment[BELT_SLOT] = &id
			added = true
		case items.BODY:
			pec.Equipment[BODY_SLOT] = &id
			added = true
		case items.HANDS:
			pec.Equipment[HANDS_SLOT] = &id
			added = true
		case items.LEGS:
			pec.Equipment[LEGS_SLOT] = &id
			added = true
		case items.FEET:
			pec.Equipment[FEET_SLOT] = &id
			added = true
		case items.RING:
			pec.Equipment[RING_SLOT] = &id
			added = true
		case items.AMULET:
			pec.Equipment[AMULET_SLOT] = &id
			added = true
		}
	}

	return added
}

type ItemID int

func (id ItemID) GetItem() items.IItem {
	return items.GetItemCollection().GetItemById(int(id))
}
