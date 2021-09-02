package components

import (
	"log"

	"github.com/alexglazkov9/survgram/entity"
	"github.com/alexglazkov9/survgram/interfaces"
)

type EffectsComponent struct {
	Parent *entity.Entity `bson:"-"`

	Effects []interfaces.IEffect
}

func (ec *EffectsComponent) Update(dt float64) {
	for _, effect := range ec.Effects {
		effect.Update(dt)
		remove := effect.Apply(ec.Parent)
		if remove {
			ec.RemoveEffect(effect)
		}
	}
}

func (ec *EffectsComponent) RemoveEffect(effect interfaces.IEffect) {
	log.Println(len(ec.Effects))
	for i, eff := range ec.Effects {
		log.Println(eff)
		log.Println(effect)
		if eff == effect {
			ec.Effects[i] = ec.Effects[len(ec.Effects)-1]
			ec.Effects = ec.Effects[:len(ec.Effects)-1]
			return
		}
	}
}

func (ec *EffectsComponent) ApplyEffect(effects []interfaces.IEffect) {
	ec.Effects = append(ec.Effects, effects...)
}
