package components

type NameComponent struct {
	BaseComponent

	Name string
}

func (nc *NameComponent) Update(dt float64) {

}

func (nc NameComponent) GetName() string {
	return nc.Name
}
