package combat

type AttackType int

const (
	PhysicalAttack AttackType = iota
	MagicalAttack
)

type AttackDetails struct {
	AttackType AttackType
	Damage     int
	Target     string
	Attacker   string
}
