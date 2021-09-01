package enemy

import (
	"github.com/alexglazkov9/survgram/components"
)

const ATTACK_SPEED float64 = 1.5

type Enemy struct {
	components.NameComponent
	components.AttackComponent
	components.HealthComponent
	components.AbilityComponent
}

func (e *Enemy) Update(dt float64) {
	e.AttackComponent.Update(dt)
}
