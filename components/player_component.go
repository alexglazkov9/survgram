package components

import (
	"github.com/alexglazkov9/survgram/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PlayerComponent struct {
	parent entity.Entity

	ID              primitive.ObjectID `bson:"_id"`
	TelegramID      int
	ChatID          int64
	CurrentLocation primitive.ObjectID
}
