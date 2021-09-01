package components

import (
	"github.com/alexglazkov9/survgram/ability"
	"github.com/alexglazkov9/survgram/entity"
)

type AbilityComponent struct {
	parent entity.Entity `bson:"-"`

	energy  int
	ability ability.Ability
}

func (ac *AbilityComponent) GainEnergy(amt int) {
	ac.energy += amt
}

func (ac *AbilityComponent) CastAbility() {

}

func (ac *AbilityComponent) IsReadyToCast() bool {
	return ac.energy >= ac.ability.Energy_cost
}
