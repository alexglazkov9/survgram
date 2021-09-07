package components

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PlayerComponent struct {
	BaseComponent

	TelegramID      int
	ChatID          int64
	CurrentLocation primitive.ObjectID
}

func (pc PlayerComponent) Update(dt float64) {

}

func (pc PlayerComponent) Save(dt float64) {

}
