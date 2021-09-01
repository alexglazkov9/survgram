package ability

type Effect struct {
	Effect_type string
	Damage      int
}

type Ability struct {
	name        string
	tier        int
	Energy_cost int
	Effects     []Effect
}
