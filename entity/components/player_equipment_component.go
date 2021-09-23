package components

type PlayerEquipmentComponent struct {
	BaseComponent `bson:"-" json:"-"`

	Weapon   *int
	Head     *int
	Body     *int
	Hands    *int
	Legs     *int
	Feet     *int
	Ring     *int
	Amulet   *int
	Belt     *int
	Backpack *int
	Tool     *int
}
