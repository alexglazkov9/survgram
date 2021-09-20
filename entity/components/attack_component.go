package components

import (
	"github.com/alexglazkov9/survgram/entity"
)

const ATTACK_SPEED float64 = 1.5

type AttackComponent struct {
	BaseComponent `bson:"-" json:"-"`

	AttackDamage int

	NextAttackTimer float64
}

func (ac AttackComponent) GetBasicAttackDamage() int {
	return ac.AttackDamage
}

func (ac AttackComponent) IsReadyToAttack() bool {
	return ac.NextAttackTimer >= ATTACK_SPEED
}

func (ac *AttackComponent) Clone() entity.IComponent {
	copy := *ac
	return &copy
}
