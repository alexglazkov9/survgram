package components

import "github.com/alexglazkov9/survgram/entity"

type PlayerActivityComponent struct {
	BaseComponent

	Type string
}

func (pac *PlayerActivityComponent) Clone() entity.IComponent {
	copy := *pac
	return &copy
}
