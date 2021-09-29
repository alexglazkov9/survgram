package components

import (
	"github.com/alexglazkov9/survgram/entity"
)

/** Component represents status of the activity it is attached to.
When `IsComplete` is true, parent system will dispose of this activity */
type ActivityStatusComponent struct {
	BaseComponent `bson:"-" json:"-"`

	IsComplete bool
}

func (asc *ActivityStatusComponent) Clone() entity.IComponent {
	copy := *asc
	return &copy
}
