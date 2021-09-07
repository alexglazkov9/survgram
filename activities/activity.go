package activities

import "github.com/alexglazkov9/survgram/entity"

type IActivity interface {
	Start()
	Update(dt float64)
	IsComplete() bool
	GetDescription() string
	AddPlayer(player *entity.Entity)
}
