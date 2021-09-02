package interfaces

import "github.com/alexglazkov9/survgram/entity"

type IEffect interface {
	Apply(target *entity.Entity) bool
	Update(dt float64)
}
