package interfaces

import "github.com/alexglazkov9/survgram/attacks"

type Attackable interface {
	IsReadyToAttack() bool
	Attack(targets []Battler) attacks.AttackDetails
}
