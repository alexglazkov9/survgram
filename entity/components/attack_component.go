package components

import (
	"math/rand"

	"github.com/alexglazkov9/survgram/entity"
	"github.com/alexglazkov9/survgram/entity/combat"
)

const ATTACK_SPEED float64 = 1.5

type AttackComponent struct {
	BaseComponent `bson:"-" json:"-"`

	AttackDamage int

	nextAttackTimer float64
}

func (ac *AttackComponent) Update(dt float64) {
	ac.nextAttackTimer += dt
}

func (ac AttackComponent) GetBasicAttackDamage() int {
	return ac.AttackDamage
}

func (ac *AttackComponent) Attack(targets []*entity.Entity) combat.AttackDetails {
	ac.nextAttackTimer = 0

	var attackDetails combat.AttackDetails
	target := targets[rand.Intn(len(targets))]
	attackDetails.Damage = target.GetComponent("HealthComponent").(*HealthComponent).ApplyDamage(ac.GetBasicAttackDamage())
	attackDetails.AttackType = combat.PhysicalAttack
	attackDetails.Target = target.GetComponent("NameComponent").(*NameComponent).GetName()

	return attackDetails
}

func (ac *AttackComponent) GetAttackComponent() *AttackComponent {
	return ac
}

func (ac AttackComponent) IsReadyToAttack() bool {
	return ac.nextAttackTimer >= ATTACK_SPEED
}
