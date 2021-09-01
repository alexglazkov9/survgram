package interfaces

type IBattle interface {
	GetAllPlayers(args ...bool) []Battler
	GetAllEnemies(args ...bool) []Battler
}
