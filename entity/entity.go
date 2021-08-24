package entity

type Entity interface {
	Update(dt float64)
	GetBasicAttackDamage() int
	IsReadyToAttack() bool
	Attack(targets []*Entity)
	TakeDamage(dmg int) int
}
