package components

import (
	"math/rand"

	"github.com/alexglazkov9/survgram/entity"
	"github.com/alexglazkov9/survgram/entity/combat"
)

type AbilityComponent struct {
	BaseComponent `bson:"-" json:"-"`

	energy  int `bson:"-"`
	Ability combat.Ability
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
