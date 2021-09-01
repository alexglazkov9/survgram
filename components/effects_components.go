package components

import (
	"github.com/alexglazkov9/survgram/ability"
	"github.com/alexglazkov9/survgram/entity"
)

type EffectsComponent struct {
	parent entity.Entity

	Effects []ability.Effect
}

func (ec *EffectsComponent) Update(dt float64) {

}
