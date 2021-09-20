package components

import (
	"github.com/alexglazkov9/survgram/entity"
)

type ActivityStatusComponent struct {
	BaseComponent `bson:"-" json:"-"`

	IsComplete bool
}

func (asc *ActivityStatusComponent) Clone() entity.IComponent {
	copy := *asc
	return &copy
}
