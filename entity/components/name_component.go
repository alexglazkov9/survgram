package components

import "github.com/alexglazkov9/survgram/entity"

type NameComponent struct {
	BaseComponent `bson:"-" json:"-"`

	Name string
}

func (nc *NameComponent) Update(dt float64) {

}

func (nc NameComponent) GetName() string {
	return nc.Name
}

func (ac *NameComponent) Clone() entity.IComponent {
	copy := *ac
	return &copy
}
