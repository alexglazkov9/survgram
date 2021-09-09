package activities

type ActivityType string

const (
	FISHING   ActivityType = "fishing"
	BATTLE    ActivityType = "battle"
	GATHERING ActivityType = "gathering"
)

type ActivityConfig struct {
	Type           ActivityType
	ActivityChance float64
	SpawnChances   map[int]float64
}
