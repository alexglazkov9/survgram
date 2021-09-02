package components

import (
	"github.com/alexglazkov9/survgram/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PlayerComponent struct {
	Parent *entity.Entity `bson:"-"`

	TelegramID      int
	ChatID          int64
	CurrentLocation primitive.ObjectID
}

func (pc PlayerComponent) Update(dt float64) {

}
