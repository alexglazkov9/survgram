package interfaces

type Battler interface {
	Updateable
	Damageable
	Attackable

	GainEnergy(amt int)
	GetName() string
	GetHealth() int
	GetMaxHealth() int
	IsAlive() bool
	CastAbility()
	IsReadyToCast() bool
}
