package interfaces

import "github.com/alexglazkov9/survgram/entity"

type IDatabase interface {
	UpdateCharacter(c *entity.Entity) bool
	GetAllCharacters() []*entity.Entity
	AddCharacter(c *entity.Entity) bool
	GetStartLocation() ILocation
}
