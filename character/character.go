package character

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Character sturct
type Character struct {
	ID              primitive.ObjectID `bson:"_id"`
	TelegramID      int
	ChatID          int64
	Name            string
	HealthPoints    int
	CurrentLocation primitive.ObjectID
}

// New - Creates a new character
func New(telegramID int, chatID int64, name string, current_location primitive.ObjectID) *Character {
	return &Character{ID: primitive.NewObjectID(), TelegramID: telegramID, ChatID: chatID, Name: name, HealthPoints: 10, CurrentLocation: current_location}
}

// // New - Creates a new character
// func (c Character) Save() {
// 	database.GetInstance().UpdateCharacter(c)
// }
