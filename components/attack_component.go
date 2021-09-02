package components

import (
	"math/rand"

	"github.com/alexglazkov9/survgram/attacks"
	"github.com/alexglazkov9/survgram/entity"
)

const ATTACK_SPEED float64 = 1.5

type AttackComponent struct {
	Parent *entity.Entity `bson:"-"`

	AttackDamage int

	nextAttackTimer float64
}

func (ac *AttackComponent) Update(dt float64) {
	ac.nextAttackTimer += dt
}

func (ac AttackComponent) GetBasicAttackDamage() int {
	return ac.AttackDamage
}

func (ac *AttackComponent) Attack(targets []*entity.Entity) attacks.AttackDetails {
	ac.nextAttackTimer = 0

	var attackDetails attacks.AttackDetails
	target := targets[rand.Intn(len(targets))]
	attackDetails.Damage = target.GetComponent("HealthComponent").(*HealthComponent).ApplyDamage(ac.GetBasicAttackDamage())
	attackDetails.AttackType = attacks.PhysicalAttack
	attackDetails.Target = target.GetComponent("NameComponent").(*NameComponent).GetName()

	return attackDetails
}

func (ac *AttackComponent) GetAttackComponent() *AttackComponent {
	return ac
}

func (ac AttackComponent) IsReadyToAttack() bool {
	return ac.nextAttackTimer >= ATTACK_SPEED
}
