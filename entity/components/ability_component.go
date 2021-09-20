package components

import (
	"github.com/alexglazkov9/survgram/entity"
	"github.com/alexglazkov9/survgram/entity/combat"
)

type AbilityComponent struct {
	BaseComponent `bson:"-" json:"-"`

	Energy  int `bson:"-"`
	Ability combat.Ability
}

func (ac *AbilityComponent) GainEnergy(amt int) {
	ac.Energy += amt
}

func (ac *AbilityComponent) IsReadyToCast() bool {
	return ac.Energy >= ac.Ability.Energy_cost
}

func (ac *AbilityComponent) Clone() entity.IComponent {
	copy := *ac
	return &copy
}
