package components

type NameComponent struct {
	BaseComponent `bson:"-" json:"-"`

	Name string
}

func (nc *NameComponent) Update(dt float64) {

}

func (nc NameComponent) GetName() string {
	return nc.Name
}
