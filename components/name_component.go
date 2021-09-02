package components

import "github.com/alexglazkov9/survgram/entity"

type NameComponent struct {
	Parent *entity.Entity `bson:"-"`

	Name string
}

func (nc *NameComponent) Update(dt float64) {

}

func (nc NameComponent) GetName() string {
	return nc.Name
}
