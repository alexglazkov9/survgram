package components

import (
	"math/rand"

	"github.com/alexglazkov9/survgram/ability"
	"github.com/alexglazkov9/survgram/entity"
)

type AbilityComponent struct {
	Parent *entity.Entity `bson:"-"`

	energy  int `bson:"-"`
	Ability ability.Ability
}

func (ac AbilityComponent) Update(dt float64) {

}

func (ac *AbilityComponent) GainEnergy(amt int) {
	ac.energy += amt
}

func (ac *AbilityComponent) CastAbility(allies []*entity.Entity, enemies []*entity.Entity) {
	target := enemies[rand.Intn(len(enemies))]
	target.GetComponent("EffectsComponent").(*EffectsComponent).ApplyEffect(ac.Ability.Effects)
	ac.energy -= ac.Ability.Energy_cost
}

func (ac *AbilityComponent) IsReadyToCast() bool {
	return ac.energy >= ac.Ability.Energy_cost
}
