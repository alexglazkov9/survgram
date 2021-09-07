package components

import "github.com/alexglazkov9/survgram/entity"

type BaseComponent struct {
	Parent *entity.Entity `bson:"-" json:"-"`
}

func (bc *BaseComponent) SetParent(e *entity.Entity) {
	bc.Parent = e
}
