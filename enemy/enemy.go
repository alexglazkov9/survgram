package enemy

import (
	"math/rand"

	"github.com/alexglazkov9/survgram/entity"
)

const ATTACK_SPEED float64 = 1.5

type Enemy struct {
	Name         string
	Health       int
	AttackDamage int

	nextAttackTime float64
}

func (e *Enemy) Update(dt float64) {
	e.nextAttackTime += dt
}

func (e Enemy) GetBasicAttackDamage() int {
	return e.AttackDamage
}

func (e *Enemy) TakeDamage(dmg int) int {
	e.Health -= dmg
	return dmg
}

func (e Enemy) IsReadyToAttack() bool {
	return e.nextAttackTime >= ATTACK_SPEED
}

func (e Enemy) Attack(targets []*entity.Entity) {
	target := targets[rand.Intn(len(targets))]

	(*target).TakeDamage(e.AttackDamage)
}
