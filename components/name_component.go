package components

import "github.com/alexglazkov9/survgram/entity"

type NameComponent struct {
	parent entity.Entity

	Name string
}

func (nc NameComponent) GetName() string {
	return nc.Name
}
