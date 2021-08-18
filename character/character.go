package character

import (
	"github.com/alexglazkov9/survgram/location"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Character sturct
type Character struct {
	ID              primitive.ObjectID `bson:"_id"`
	TelegramID      int
	Name            string
	HealthPoints    int
	CurrentLocation location.Location
}

// New - Creates a new character
func New(telegramID int, name string) *Character {
	return &Character{TelegramID: telegramID, Name: name, HealthPoints: 10}
}

// // New - Creates a new character
// func (c Character) Save() {
// 	database.GetInstance().UpdateCharacter(c)
// }
