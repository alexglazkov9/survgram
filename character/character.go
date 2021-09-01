package character

import (
	"github.com/alexglazkov9/survgram/components"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const ATTACK_SPEED float64 = 1.5

// Character sturct
type Character struct {
	components.PlayerComponent
	components.NameComponent
	components.AttackComponent
	components.HealthComponent
	components.AbilityComponent
}

// New - Creates a new character
func New(telegramID int, chatID int64, name string, current_location primitive.ObjectID) *Character {
	return &Character{
		NameComponent: components.NameComponent{
			Name: name,
		},
		PlayerComponent: components.PlayerComponent{
			ID:              primitive.NewObjectID(),
			TelegramID:      telegramID,
			ChatID:          chatID,
			CurrentLocation: current_location,
		},
		AttackComponent: components.AttackComponent{
			AttackDamage: 5,
		},
		HealthComponent: components.HealthComponent{
			MaxHealthPoints: 100,
			HealthPoints:    100,
		},
	}
}

func (c *Character) Update(dt float64) {
	c.AttackComponent.Update(dt)
}
